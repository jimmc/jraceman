package domain

// RegistrationRepo describes how Registration records are loaded and saved.
type RegistrationRepo interface {
  FindByID(ID string) (*Registration, error)
  List(offset, limit int) ([]*Registration, error)
  Save(*Registration) (string, error)
  UpdateByID(ID string, oldRegistration, newRegistration *Registration, diffs Diffs) error
  DeleteByID(ID string) error
}

// Registration describes a registration for a competitor.
type Registration struct {
  ID string
  MeetID string
  PersonID string
  AmountCharged int
  Surcharge int
  Discount int
  AmountPaid int
  // Due int - virtual field: (AmountCharged + Surcharge) - (Discount + AmountPaid)
  WaiverSigned bool
  PaymentNotes *string
}

// RegistrationMeta provides funcions related to the Registration struct.
type RegistrationMeta struct {}

func (m *RegistrationMeta) EntityTypeName() string {
  return "registration"
}

func (m *RegistrationMeta) NewEntity() interface{} {
  return &Registration{}
}
