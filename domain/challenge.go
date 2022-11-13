package domain

// ChallengeRepo describes how Challenge records are loaded and saved.
type ChallengeRepo interface {
  FindByID(ID string) (*Challenge, error)
  List(offset, limit int) ([]*Challenge, error)
  Save(*Challenge) (string, error)
  UpdateByID(ID string, oldChallenge, newChallenge *Challenge, diffs Diffs) error
  DeleteByID(ID string) error
}

// Challenge describes a challenge group.
type Challenge struct {
  ID string
  Name string
}

// ChallengeMeta provides funcions related to the Challenge struct.
type ChallengeMeta struct {}

func (m *ChallengeMeta) EntityTypeName() string {
  return "challenge"
}

func (m *ChallengeMeta) EntityGroupName() string {
  return "roster"
}

func (m *ChallengeMeta) NewEntity() interface{} {
  return &Challenge{}
}
