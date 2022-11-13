package domain

// RoleRepo describes how Role records are loaded and saved.
type RoleRepo interface {
  FindByID(ID string) (*Role, error)
  List(offset, limit int) ([]*Role, error)
  Save(*Role) (string, error)
  UpdateByID(ID string, oldRole, newRole *Role, diffs Diffs) error
  DeleteByID(ID string) error
}

// Role describes a role as in the RBAC model.
// See also user, permission, userrole, rolepermission.
type Role struct {
  ID string
  Name string
  Description string
}

// RoleMeta provides funcions related to the Role struct.
type RoleMeta struct {}

func (m *RoleMeta) EntityTypeName() string {
  return "role"
}

func (m *RoleMeta) NewEntity() interface{} {
  return &Role{}
}
