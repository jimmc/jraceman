package domain

// RoleRoleRepo describes how RoleRole records are loaded and saved.
type RoleRoleRepo interface {
  FindByID(ID string) (*RoleRole, error)
  List(offset, limit int) ([]*RoleRole, error)
  Save(*RoleRole) (string, error)
  UpdateByID(ID string, oldRoleRole, newRoleRole *RoleRole, diffs Diffs) error
  DeleteByID(ID string) error
}

// RoleRole describes which additional roles are granted by a role.
type RoleRole struct {
  ID string
  RoleID string         // If a user has this role...
  HasRoleID string      //   then he also gets this role (recursively).
}
