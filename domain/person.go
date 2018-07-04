package domain

// PersonRepo describes how Person records are loaded and saved.
type PersonRepo interface {
  FindByID(ID string) (*Person, error)
  List(offset, limit int) ([]*Person, error)
  Save(*Person) (string, error)
  UpdateByID(ID string, oldPerson, newPerson *Person, diffs Diffs) error
  DeleteByID(ID string) error
}

// Person describes the information we store about a person.
type Person struct {
  ID string
  FirstName string
  LastName string
  GenderID string
  TeamID string
  Birthday string       // ISO8601 format, may be partial
  Phone string
  Email string
  Membership string
  MembershipExpiration string
}
