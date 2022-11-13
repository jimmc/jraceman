package domain

// RolePermissionRepo describes how RolePermission records are loaded and saved.
type RolePermissionRepo interface {
  FindByID(ID string) (*RolePermission, error)
  List(offset, limit int) ([]*RolePermission, error)
  Save(*RolePermission) (string, error)
  UpdateByID(ID string, oldRolePermission, newRolePermission *RolePermission, diffs Diffs) error
  DeleteByID(ID string) error
}

// RolePermission describes which permissions are granted by a role.
type RolePermission struct {
  ID string
  RoleID string
  PermissionID string
}

// RolePermissionMeta provides funcions related to the RolePermission struct.
type RolePermissionMeta struct {}

func (m *RolePermissionMeta) EntityTypeName() string {
  return "rolepermission"
}

func (m *RolePermissionMeta) EntityGroupName() string {
  return "auth"
}

func (m *RolePermissionMeta) NewEntity() interface{} {
  return &RolePermission{}
}
