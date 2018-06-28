package domain

// LevelRepo describes how Level records are loaded and saved.
type LevelRepo interface {
  FindByID(ID string) (*Level, error)
  List(offset, limit int) ([]*Level, error)
  Save(*Level) (string, error)
  UpdateByID(ID string, oldLevel, newLevel *Level, diffs Diffs) error
  DeleteByID(ID string) error
}

// Level describes an age level, used in age-based events.
type Level struct {
  ID string
  Name string
  MinEntryAge *int
  MinAge *int
  MaxAge *int
  MaxEntryAge *int
  UseGroupAverage *bool
}
