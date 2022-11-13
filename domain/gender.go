package domain

// GenderRepo describes how Gender records are loaded and saved.
type GenderRepo interface {
  FindByID(ID string) (*Gender, error)
  List(offset, limit int) ([]*Gender, error)
  Save(*Gender) (string, error)
  UpdateByID(ID string, oldGender, newGender *Gender, diffs Diffs) error
  DeleteByID(ID string) error
}

// Gender describes a person or an event gender, such as Men or Women, Boy or Girl, Open or Mixed.
type Gender struct {
  ID string
  Name string
}

// GenderMeta provides funcions related to the Gender struct.
type GenderMeta struct {}

func (m *GenderMeta) EntityTypeName() string {
  return "gender"
}

func (m *GenderMeta) NewEntity() interface{} {
  return &Gender{}
}
