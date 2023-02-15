package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBUserRepo struct {
  db conn.DB
}

func (r *DBUserRepo) New() interface{} {
  return domain.User{}
}

func (r *DBUserRepo) UpdateColumnInfos(columnInfos []structsql.ColumnInfo) []structsql.ColumnInfo {
  // We need the username column to be unique.
  for c, col := range columnInfos {
    if col.Name == "username" {
      columnInfos[c].Unique = true
    }
  }
  return columnInfos
}

func (r *DBUserRepo) CreateTable() error {
  return structsql.CreateTableWithUpdater(r.db, "user", domain.User{}, r)
}

func (r *DBUserRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTableWithUpdater(r.db, "user", domain.User{}, dryrun, r)
}

func (r *DBUserRepo) FindByID(ID string) (*domain.User, error) {
  user := &domain.User{}
  sql, targets := structsql.FindByIDSql("user", user)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return user, nil
}

func (r *DBUserRepo) Save(user *domain.User) (string, error) {
  if (user.ID == "") {
    user.ID = structsql.UniqueID(r.db, "user", "U1")
  }
  return user.ID, structsql.Insert(r.db, "user", user, user.ID)
}

func (r *DBUserRepo) List(offset, limit int) ([]*domain.User, error) {
  user := &domain.User{}
  users := make([]*domain.User, 0)
  sql, targets := structsql.ListSql("user", user, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    userCopy := domain.User(*user)
    users = append(users, &userCopy)
  })
  return users, err
}

func (r *DBUserRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "user", ID)
}

func (r *DBUserRepo) UpdateByID(ID string, oldUser, newUser *domain.User, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "user", diffs.Modified(), ID)
}

func (r *DBUserRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "user", &domain.User{})
}
