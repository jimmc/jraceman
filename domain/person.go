package domain

// PersonRepo describes how Person records are loaded and saved.
type PersonRepo interface {
  FindByID(ID string) (Person, error)
  Save(Person) error
}

// Person describes the information we store about a person.
type Person struct {
  ID string
  FirstName string
  LastName string
  GenderID string
  TeamId string
  Birthday string       // ISO8601 format, may be partial
  Phone string
  Email string
  Membership string
  MembershipExpiration string
}
