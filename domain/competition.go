package domain

// CompetitionRepo describes how Competition records are loaded and saved.
type CompetitionRepo interface {
  FindByID(ID string) (*Competition, error)
  List(offset, limit int) ([]*Competition, error)
  Save(*Competition) (string, error)
  UpdateByID(ID string, oldCompetition, newCompetition *Competition, diffs Diffs) error
  DeleteByID(ID string) error
}

// Competition describes a class of events, which are typically then
// divided up according to level.
type Competition struct {
  ID string
  Name string
  GroupSize *int
  MaxAlternates *int
  Duration *int         // duration in seconds
}
