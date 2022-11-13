package domain

// SimplanRuleRepo describes how SimplanRule records are loaded and saved.
type SimplanRuleRepo interface {
  FindByID(ID string) (*SimplanRule, error)
  List(offset, limit int) ([]*SimplanRule, error)
  Save(*SimplanRule) (string, error)
  UpdateByID(ID string, oldSimplanRule, newSimplanRule *SimplanRule, diffs Diffs) error
  DeleteByID(ID string) error
}

// SimplanRule describes a progression rule in a Simplan.
type SimplanRule struct {
  ID string
  SimplanID string
  FromStageID string
  ToStageID string
  ThruPlace *int
  NextBestTimes *int
}

// SimplanRuleMeta provides funcions related to the SimplanRule struct.
type SimplanRuleMeta struct {}

func (m *SimplanRuleMeta) EntityTypeName() string {
  return "simplanrule"
}

func (m *SimplanRuleMeta) NewEntity() interface{} {
  return &SimplanRule{}
}
