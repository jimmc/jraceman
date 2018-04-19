package domain

// LevelRepo describes how Level records are loaded and saved.
type LevelRepo interface {
  FindById(ID string) (Level, error)
  Save(Level) error
}

// Level describes an age bracket for an event.
type Level struct {
  ID string
  Name string
  MinAge int
  MaxAge int
  MinEntryAge int
  MaxEntryAge int
  UseGroupAverage bool
}
