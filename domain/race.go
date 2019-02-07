package domain

// RaceRepo describes how Race records are loaded and saved.
type RaceRepo interface {
  FindByID(ID string) (*Race, error)
  List(offset, limit int) ([]*Race, error)
  Save(*Race) (string, error)
  UpdateByID(ID string, oldRace, newRace *Race, diffs Diffs) error
  DeleteByID(ID string) error
}

// Race describes one race in an event.
type Race struct {
  ID string
  EventID string
  StageID *string
  Round *int
  Section *int
  AreaID *string
  Number *int
  ScheduledStart *string        // datetime
  ScheduledDuration *int        // duration
  ActualStart *string          // datetime
  Scratched bool
  RaceComment *string
}
