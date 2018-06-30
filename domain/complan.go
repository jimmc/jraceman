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
