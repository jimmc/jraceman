package domain

// SimplanStageRepo describes how SimplanStage records are loaded and saved.
type SimplanStageRepo interface {
  FindByID(ID string) (*SimplanStage, error)
  List(offset, limit int) ([]*SimplanStage, error)
  Save(*SimplanStage) (string, error)
  UpdateByID(ID string, oldSimplanStage, newSimplanStage *SimplanStage, diffs Diffs) error
  DeleteByID(ID string) error
}

// SimplanStage describes one stage in a Simplan.
type SimplanStage struct {
  ID string
  SimplanID string
  StageID string
  SectionCount int
  FillOrder *string
}
