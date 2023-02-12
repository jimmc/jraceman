package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/compat"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBRoleRepo struct {
  db compat.DBorTx
}

func (r *DBRoleRepo) New() interface{} {
  return domain.Role{}
}

func (r *DBRoleRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "role", domain.Role{})
}

func (r *DBRoleRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "role", domain.Role{}, dryrun)
}

func (r *DBRoleRepo) FindByID(ID string) (*domain.Role, error) {
  role := &domain.Role{}
  sql, targets := structsql.FindByIDSql("role", role)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return role, nil
}

func (r *DBRoleRepo) Save(role *domain.Role) (string, error) {
  if (role.ID == "") {
    role.ID = structsql.UniqueID(r.db, "role", "R1")
  }
  return role.ID, structsql.Insert(r.db, "role", role, role.ID)
}

func (r *DBRoleRepo) List(offset, limit int) ([]*domain.Role, error) {
  role := &domain.Role{}
  roles := make([]*domain.Role, 0)
  sql, targets := structsql.ListSql("role", role, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    roleCopy := domain.Role(*role)
    roles = append(roles, &roleCopy)
  })
  return roles, err
}

func (r *DBRoleRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "role", ID)
}

func (r *DBRoleRepo) UpdateByID(ID string, oldRole, newRole *domain.Role, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "role", diffs.Modified(), ID)
}

func (r *DBRoleRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "role", &domain.Role{})
}
