package domain

// ScoringRuleRepo describes how ScoringRule records are loaded and saved.
type ScoringRuleRepo interface {
  FindByID(ID string) (*ScoringRule, error)
  List(offset, limit int) ([]*ScoringRule, error)
  Save(*ScoringRule) (string, error)
  UpdateByID(ID string, oldScoringRule, newScoringRule *ScoringRule, diffs Diffs) error
  DeleteByID(ID string) error
}

// ScoringRule describes one rule in a scoring system.
type ScoringRule struct {
  ID string
  ScoringSystemID string
  Rule string
  Value int
  Points float32
}

// ScoringRuleMeta provides funcions related to the ScoringRule struct.
type ScoringRuleMeta struct {}

func (m *ScoringRuleMeta) EntityTypeName() string {
  return "scoringrule"
}

func (m *ScoringRuleMeta) EntityGroupName() string {
  return "sport"
}

func (m *ScoringRuleMeta) NewEntity() interface{} {
  return &ScoringRule{}
}
