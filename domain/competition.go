package domain

// CompetitionRepo describes how Competition records are loaded and saved.
type CompetitionRepo interface {
  FindByID(ID string) (Competition, error)
  Save(Competition) error
}

// Competition describes a class of events, which are typically then
// divided up according to level.
type Competition struct {
  ID string
  Name string
  GroupSize int
  MaxAlternates int
}
