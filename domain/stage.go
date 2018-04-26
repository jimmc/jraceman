package domain

// StageRepo describes how Stage records are loaded and saved.
type StageRepo interface {
  FindByID(ID string) (Stage, error)
  Save(Stage) error
}

// Stage describes things like heats and semifinals.
type Stage struct {
  ID string
  Name string
  Number int
  IsFinal bool
}
