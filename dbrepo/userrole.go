package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBUserRoleRepo struct {
  db *sql.DB
}

func (r *DBUserRoleRepo) New() interface{} {
  return domain.UserRole{}
}

func (r *DBUserRoleRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "userrole", domain.UserRole{})
}

func (r *DBUserRoleRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "userrole", domain.UserRole{}, dryrun)
}

func (r *DBUserRoleRepo) FindByID(ID string) (*domain.UserRole, error) {
  userrole := &domain.UserRole{}
  sql, targets := structsql.FindByIDSql("userrole", userrole)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return userrole, nil
}

func (r *DBUserRoleRepo) Save(userrole *domain.UserRole) (string, error) {
  if (userrole.ID == "") {
    userrole.ID = structsql.UniqueID(r.db, "userrole", "UR1")
  }
  return userrole.ID, structsql.Insert(r.db, "userrole", userrole, userrole.ID)
}

func (r *DBUserRoleRepo) List(offset, limit int) ([]*domain.UserRole, error) {
  userrole := &domain.UserRole{}
  userroles := make([]*domain.UserRole, 0)
  sql, targets := structsql.ListSql("userrole", userrole, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    userroleCopy := domain.UserRole(*userrole)
    userroles = append(userroles, &userroleCopy)
  })
  return userroles, err
}

func (r *DBUserRoleRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "userrole", ID)
}

func (r *DBUserRoleRepo) UpdateByID(ID string, oldUserRole, newUserRole *domain.UserRole, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "userrole", diffs.Modified(), ID)
}

func (r *DBUserRoleRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "userrole", &domain.UserRole{})
}
