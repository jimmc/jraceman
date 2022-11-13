package domain

// ScoringSystemRepo describes how ScoringSystem records are loaded and saved.
type ScoringSystemRepo interface {
  FindByID(ID string) (*ScoringSystem, error)
  List(offset, limit int) ([]*ScoringSystem, error)
  Save(*ScoringSystem) (string, error)
  UpdateByID(ID string, oldScoringSystem, newScoringSystem *ScoringSystem, diffs Diffs) error
  DeleteByID(ID string) error
}

// ScoringSystem defines the name of a scoring system.
type ScoringSystem struct {
  ID string
  Name string
}

// ScoringSystemMeta provides funcions related to the ScoringSystem struct.
type ScoringSystemMeta struct {}

func (m *ScoringSystemMeta) EntityTypeName() string {
  return "scoringsystem"
}

func (m *ScoringSystemMeta) NewEntity() interface{} {
  return &ScoringSystem{}
}
