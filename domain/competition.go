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
  ScheduledDuration *int         // duration in seconds
}

// CompetitionMeta provides funcions related to the Competition struct.
type CompetitionMeta struct {}

func (m *CompetitionMeta) EntityTypeName() string {
  return "competition"
}

func (m *CompetitionMeta) EntityGroupName() string {
  return "sport"
}

func (m *CompetitionMeta) NewEntity() interface{} {
  return &Competition{}
}
