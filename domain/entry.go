package domain

// EntryRepo describes how Entry records are loaded and saved.
type EntryRepo interface {
  FindByID(ID string) (*Entry, error)
  List(offset, limit int) ([]*Entry, error)
  Save(*Entry) (string, error)
  UpdateByID(ID string, oldEntry, newEntry *Entry, diffs Diffs) error
  DeleteByID(ID string) error
}

// Entry describes an entry for a person to an event.
type Entry struct {
  ID string
  PersonID string
  EventID string
  GroupName *string
  Alternate bool
  Scratched bool
}

// EntryMeta provides funcions related to the Entry struct.
type EntryMeta struct {}

func (m *EntryMeta) EntityTypeName() string {
  return "entry"
}

func (m *EntryMeta) EntityGroupName() string {
  return "regatta"
}

func (m *EntryMeta) NewEntity() interface{} {
  return &Entry{}
}
