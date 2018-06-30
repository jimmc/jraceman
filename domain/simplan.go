package domain

// SimplanRepo describes how Simplan records are loaded and saved.
type SimplanRepo interface {
  FindByID(ID string) (*Simplan, error)
  List(offset, limit int) ([]*Simplan, error)
  Save(*Simplan) (string, error)
  UpdateByID(ID string, oldSimplan, newSimplan *Simplan, diffs Diffs) error
  DeleteByID(ID string) error
}

// Simplan describes a simple progression plan for an event.
type Simplan struct {
  ID string
  System string
  Plan string
  MinEntries int
  MaxEntries int
}
