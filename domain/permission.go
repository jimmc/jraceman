package domain

// PermissionRepo describes how Permission records are loaded and saved.
type PermissionRepo interface {
  FindByID(ID string) (*Permission, error)
  List(offset, limit int) ([]*Permission, error)
  Save(*Permission) (string, error)
  UpdateByID(ID string, oldPermission, newPermission *Permission, diffs Diffs) error
  DeleteByID(ID string) error
}

// Permission describes a permission as used by the application.
// See also user, role, userrole, rolepermission.
type Permission struct {
  ID string
  Name string
  Description string
}
