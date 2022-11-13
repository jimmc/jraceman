package domain

// ContextOptionRepo describes how ContextOption records are loaded and saved.
type ContextOptionRepo interface {
  FindByID(ID string) (*ContextOption, error)
  List(offset, limit int) ([]*ContextOption, error)
  Save(*ContextOption) (string, error)
  UpdateByID(ID string, oldContextOption, newContextOption *ContextOption, diffs Diffs) error
  DeleteByID(ID string) error
}

// ContextOption describes an option tied to some context.
type ContextOption struct {
  ID string
  Name string
  Value string
  Host *string
  WebContext *string
  MeetID *string
}

// ContextOptionMeta provides funcions related to the ContextOption struct.
type ContextOptionMeta struct {}

func (m *ContextOptionMeta) EntityTypeName() string {
  return "contextoption"
}

func (m *ContextOptionMeta) NewEntity() interface{} {
  return &ContextOption{}
}
