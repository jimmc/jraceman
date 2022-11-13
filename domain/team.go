package domain

// TeamRepo describes how Team records are loaded and saved.
type TeamRepo interface {
  FindByID(ID string) (*Team, error)
  List(offset, limit int) ([]*Team, error)
  Save(*Team) (string, error)
  UpdateByID(ID string, oldTeam, newTeam *Team, diffs Diffs) error
  DeleteByID(ID string) error
}

// Team describes a record containing information about one team.
type Team struct {
  ID string
  ShortName string
  Name string
  Number *int
  ChallengeID *string
  NonScoring bool
  Street *string
  Street2 *string
  City *string
  State *string
  Zip *string
  Country *string
  Phone *string
  Fax *string
}

// TeamMeta provides funcions related to the Team struct.
type TeamMeta struct {}

func (m *TeamMeta) EntityTypeName() string {
  return "team"
}

func (m *TeamMeta) EntityGroupName() string {
  return "roster"
}

func (m *TeamMeta) NewEntity() interface{} {
  return &Team{}
}
