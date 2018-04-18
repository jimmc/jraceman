package domain

// TeamRepo describes how Team records are loaded and saved.
type TeamRepo interface {
  FindByID(ID string) (Team, error)
  Save(Team) error
}

// Team describes a record containing information about one team.
type Team struct {
  ID string
  ShortName string
  Name string
  Number int
  ChallengeId string
  NonScoring bool
  Street string
  Street2 string
  City string
  State string
  Zip string
  Country string
  Phone string
  Fax string
}
