package domain

// UserRoleRepo describes how UserRole records are loaded and saved.
type UserRoleRepo interface {
  FindByID(ID string) (*UserRole, error)
  List(offset, limit int) ([]*UserRole, error)
  Save(*UserRole) (string, error)
  UpdateByID(ID string, oldUserRole, newUserRole *UserRole, diffs Diffs) error
  DeleteByID(ID string) error
}

// UserRole specified which roles a user has.
// See also user, role, permission, rolepermission.
type UserRole struct {
  ID string
  UserID string
  RoleID string
}

// UserRoleMeta provides funcions related to the UserRole struct.
type UserRoleMeta struct {}

func (m *UserRoleMeta) EntityTypeName() string {
  return "userrole"
}

func (m *UserRoleMeta) EntityGroupName() string {
  return "auth"
}

func (m *UserRoleMeta) NewEntity() interface{} {
  return &UserRole{}
}
