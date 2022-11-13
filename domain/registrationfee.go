package domain

// RegistrationFeeRepo describes how RegistrationFee records are loaded and saved.
type RegistrationFeeRepo interface {
  FindByID(ID string) (*RegistrationFee, error)
  List(offset, limit int) ([]*RegistrationFee, error)
  Save(*RegistrationFee) (string, error)
  UpdateByID(ID string, oldRegistrationFee, newRegistrationFee *RegistrationFee, diffs Diffs) error
  DeleteByID(ID string) error
}

// RegistrationFee describes a registration fee.
type RegistrationFee struct {
  ID string
  MeetID string
  EventCount int
  AmountCharged int
}

// RegistrationFeeMeta provides funcions related to the RegistrationFee struct.
type RegistrationFeeMeta struct {}

func (m *RegistrationFeeMeta) EntityTypeName() string {
  return "registrationfee"
}

func (m *RegistrationFeeMeta) NewEntity() interface{} {
  return &RegistrationFee{}
}
