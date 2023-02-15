package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBRolePermissionRepo struct {
  db conn.DB
}

func (r *DBRolePermissionRepo) New() interface{} {
  return domain.RolePermission{}
}

func (r *DBRolePermissionRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "rolepermission", domain.RolePermission{})
}

func (r *DBRolePermissionRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "rolepermission", domain.RolePermission{}, dryrun)
}

func (r *DBRolePermissionRepo) FindByID(ID string) (*domain.RolePermission, error) {
  rolepermission := &domain.RolePermission{}
  sql, targets := structsql.FindByIDSql("rolepermission", rolepermission)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return rolepermission, nil
}

func (r *DBRolePermissionRepo) Save(rolepermission *domain.RolePermission) (string, error) {
  if (rolepermission.ID == "") {
    rolepermission.ID = structsql.UniqueID(r.db, "rolepermission", "RPM1")
  }
  return rolepermission.ID, structsql.Insert(r.db, "rolepermission", rolepermission, rolepermission.ID)
}

func (r *DBRolePermissionRepo) List(offset, limit int) ([]*domain.RolePermission, error) {
  rolepermission := &domain.RolePermission{}
  rolepermissions := make([]*domain.RolePermission, 0)
  sql, targets := structsql.ListSql("rolepermission", rolepermission, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    rolepermissionCopy := domain.RolePermission(*rolepermission)
    rolepermissions = append(rolepermissions, &rolepermissionCopy)
  })
  return rolepermissions, err
}

func (r *DBRolePermissionRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "rolepermission", ID)
}

func (r *DBRolePermissionRepo) UpdateByID(ID string, oldRolePermission, newRolePermission *domain.RolePermission, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "rolepermission", diffs.Modified(), ID)
}

func (r *DBRolePermissionRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "rolepermission", &domain.RolePermission{})
}
