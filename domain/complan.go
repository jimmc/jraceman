package domain

// ComplanRepo describes how Complan records are loaded and saved.
type ComplanRepo interface {
  FindByID(ID string) (*Complan, error)
  List(offset, limit int) ([]*Complan, error)
  Save(*Complan) (string, error)
  UpdateByID(ID string, oldComplan, newComplan *Complan, diffs Diffs) error
  DeleteByID(ID string) error
}

// Complan describes a complex progression plan for an event.
type Complan struct {
  ID string
  System string
  Plan string
  MinEntries int
  MaxEntries int
  PlanOrder int
}

// ComplanMeta provides funcions related to the Complan struct.
type ComplanMeta struct {}

func (m *ComplanMeta) EntityTypeName() string {
  return "complan"
}

func (m *ComplanMeta) EntityGroupName() string {
  return "plan"
}

func (m *ComplanMeta) NewEntity() interface{} {
  return &Complan{}
}
