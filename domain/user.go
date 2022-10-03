package domain

// UserRepo describes how User records are loaded and saved.
type UserRepo interface {
  FindByID(ID string) (*User, error)
  List(offset, limit int) ([]*User, error)
  Save(*User) (string, error)
  UpdateByID(ID string, oldUser, newUser *User, diffs Diffs) error
  DeleteByID(ID string) error
}

// User describes a user of JRaceman.
// Permissions are done via roles. See these tables:
// role, permission, userrole, rolepermission.
type User struct {
  ID string
  Username string
  Cryptword string
}
