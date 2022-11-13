package domain

// ProgressionRepo describes how Progression records are loaded and saved.
type ProgressionRepo interface {
  FindByID(ID string) (*Progression, error)
  List(offset, limit int) ([]*Progression, error)
  Save(*Progression) (string, error)
  UpdateByID(ID string, oldProgression, newProgression *Progression, diffs Diffs) error
  DeleteByID(ID string) error
}

// Progression defines a named progression for use in events.
type Progression struct {
  ID string
  Name string
  Class string
  Parameters *string
}

// ProgressionMeta provides funcions related to the Progression struct.
type ProgressionMeta struct {}

func (m *ProgressionMeta) EntityTypeName() string {
  return "progression"
}

func (m *ProgressionMeta) EntityGroupName() string {
  return "sport"
}

func (m *ProgressionMeta) NewEntity() interface{} {
  return &Progression{}
}
