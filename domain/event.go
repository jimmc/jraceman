package domain

// EventRepo describes how Event records are loaded and saved.
type EventRepo interface {
  FindByID(ID string) (*Event, error)
  List(offset, limit int) ([]*Event, error)
  Save(*Event) (string, error)
  UpdateByID(ID string, oldEvent, newEvent *Event, diffs Diffs) error
  DeleteByID(ID string) error
}

// Event describes an event such as a race.
type Event struct {
  ID string
  MeetID string
  Number *int
  Name *string
  CompetitionID string
  LevelID string
  GenderID string
  AreaID *string
  SeedingPlanID *string
  ProgressionID *string
  ProgressionState *string
  ScoringSystemID *string
  Scratched bool
  EventComment *string
}
