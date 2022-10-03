package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBPermissionRepo struct {
  db *sql.DB
}

func (r *DBPermissionRepo) New() interface{} {
  return domain.Permission{}
}

func (r *DBPermissionRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "permission", domain.Permission{})
}

func (r *DBPermissionRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "permission", domain.Permission{}, dryrun)
}

func (r *DBPermissionRepo) FindByID(ID string) (*domain.Permission, error) {
  permission := &domain.Permission{}
  sql, targets := structsql.FindByIDSql("permission", permission)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return permission, nil
}

func (r *DBPermissionRepo) Save(permission *domain.Permission) (string, error) {
  if (permission.ID == "") {
    permission.ID = structsql.UniqueID(r.db, "permission", "PM1")
  }
  return permission.ID, structsql.Insert(r.db, "permission", permission, permission.ID)
}

func (r *DBPermissionRepo) List(offset, limit int) ([]*domain.Permission, error) {
  permission := &domain.Permission{}
  permissions := make([]*domain.Permission, 0)
  sql, targets := structsql.ListSql("permission", permission, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    permissionCopy := domain.Permission(*permission)
    permissions = append(permissions, &permissionCopy)
  })
  return permissions, err
}

func (r *DBPermissionRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "permission", ID)
}

func (r *DBPermissionRepo) UpdateByID(ID string, oldPermission, newPermission *domain.Permission, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "permission", diffs.Modified(), ID)
}

func (r *DBPermissionRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "permission", &domain.Permission{})
}
