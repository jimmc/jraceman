package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBRoleRoleRepo struct {
  db conn.DB
}

func (r *DBRoleRoleRepo) New() interface{} {
  return domain.RoleRole{}
}

func (r *DBRoleRoleRepo) UpdateColumnInfos(columnInfos []structsql.ColumnInfo) []structsql.ColumnInfo {
  // We have two foreign key columns to the role table, so we need
  // to special-case our columninfo creation.
  for c, col := range columnInfos {
    if col.Name == "hasroleid" {
      columnInfos[c].FKTable = "role"
    }
  }
  return columnInfos
}

func (r *DBRoleRoleRepo) CreateTable() error {
  return structsql.CreateTableWithUpdater(r.db, "rolerole", domain.RoleRole{}, r)
}

func (r *DBRoleRoleRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTableWithUpdater(r.db, "rolerole", domain.RoleRole{}, dryrun, r)
}

func (r *DBRoleRoleRepo) FindByID(ID string) (*domain.RoleRole, error) {
  rolerole := &domain.RoleRole{}
  sql, targets := structsql.FindByIDSql("rolerole", rolerole)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return rolerole, nil
}

func (r *DBRoleRoleRepo) Save(rolerole *domain.RoleRole) (string, error) {
  if (rolerole.ID == "") {
    rolerole.ID = structsql.UniqueID(r.db, "rolerole", "RR1")
  }
  return rolerole.ID, structsql.Insert(r.db, "rolerole", rolerole, rolerole.ID)
}

func (r *DBRoleRoleRepo) List(offset, limit int) ([]*domain.RoleRole, error) {
  rolerole := &domain.RoleRole{}
  roleroles := make([]*domain.RoleRole, 0)
  sql, targets := structsql.ListSql("rolerole", rolerole, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    roleroleCopy := domain.RoleRole(*rolerole)
    roleroles = append(roleroles, &roleroleCopy)
  })
  return roleroles, err
}

func (r *DBRoleRoleRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "rolerole", ID)
}

func (r *DBRoleRoleRepo) UpdateByID(ID string, oldRoleRole, newRoleRole *domain.RoleRole, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "rolerole", diffs.Modified(), ID)
}

func (r *DBRoleRoleRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "rolerole", &domain.RoleRole{})
}
