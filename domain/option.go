package domain

// OptionRepo describes how Option records are loaded and saved.
type OptionRepo interface {
  FindByID(ID string) (*Option, error)
  List(offset, limit int) ([]*Option, error)
  Save(*Option) (string, error)
  UpdateByID(ID string, oldOption, newOption *Option, diffs Diffs) error
  DeleteByID(ID string) error
}

// Option describes an optional named value
type Option struct {
  Name string
  Value *string
}

// OptionMeta provides funcions related to the Option struct.
type OptionMeta struct {}

func (m *OptionMeta) EntityTypeName() string {
  return "option"
}

func (m *OptionMeta) NewEntity() interface{} {
  return &Option{}
}
