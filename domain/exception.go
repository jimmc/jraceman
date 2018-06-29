package domain

// ExceptionRepo describes how Exception records are loaded and saved.
type ExceptionRepo interface {
  FindByID(ID string) (*Exception, error)
  List(offset, limit int) ([]*Exception, error)
  Save(*Exception) (string, error)
  UpdateByID(ID string, oldException, newException *Exception, diffs Diffs) error
  DeleteByID(ID string) error
}

// Exception describes a reason for no finish in an race.
type Exception struct {
  ID string
  Name string
  ShortName string
  ResultAllowedRequired int
}
