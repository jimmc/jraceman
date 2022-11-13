package domain

// ComplanStageRepo describes how ComplanStage records are loaded and saved.
type ComplanStageRepo interface {
  FindByID(ID string) (*ComplanStage, error)
  List(offset, limit int) ([]*ComplanStage, error)
  Save(*ComplanStage) (string, error)
  UpdateByID(ID string, oldComplanStage, newComplanStage *ComplanStage, diffs Diffs) error
  DeleteByID(ID string) error
}

// ComplanStage describes one stage in a Complan.
type ComplanStage struct {
  ID string
  ComplanID string
  StageID string
  Round int
  SectionCount int
  FillOrder *string
}

// ComplanStageMeta provides funcions related to the ComplanStage struct.
type ComplanStageMeta struct {}

func (m *ComplanStageMeta) EntityTypeName() string {
  return "complanstage"
}

func (m *ComplanStageMeta) NewEntity() interface{} {
  return &ComplanStage{}
}
