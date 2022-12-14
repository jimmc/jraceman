package domain

// StageRepo describes how Stage records are loaded and saved.
type StageRepo interface {
  FindByID(ID string) (*Stage, error)
  List(offset, limit int) ([]*Stage, error)
  Save(*Stage) (string, error)
  UpdateByID(ID string, oldStage, newStage *Stage, diffs Diffs) error
  DeleteByID(ID string) error
}

// Stage describes things like heats and semifinals.
type Stage struct {
  ID string
  Name string
  Number int
  IsFinal bool
}

// StageMeta provides funcions related to the Stage struct.
type StageMeta struct {}

func (m *StageMeta) EntityTypeName() string {
  return "stage"
}

func (m *StageMeta) EntityGroupName() string {
  return "sport"
}

func (m *StageMeta) NewEntity() interface{} {
  return &Stage{}
}
