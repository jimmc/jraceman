package domain

// LaneRepo describes how Lane records are loaded and saved.
type LaneRepo interface {
  FindByID(ID string) (*Lane, error)
  List(offset, limit int) ([]*Lane, error)
  Save(*Lane) (string, error)
  UpdateByID(ID string, oldLane, newLane *Lane, diffs Diffs) error
  DeleteByID(ID string) error
}

// Lane describes one lane entry for one person in one race.
type Lane struct {
  ID string
  EntryID string
  RaceID string
  Lane *int
  Result *int
  ExceptionID *string
  Place *int
  ScorePlace *int
  Score *float32
}

// LaneMeta provides funcions related to the Lane struct.
type LaneMeta struct {}

func (m *LaneMeta) EntityTypeName() string {
  return "lane"
}

func (m *LaneMeta) EntityGroupName() string {
  return "regatta"
}

func (m *LaneMeta) NewEntity() interface{} {
  return &Lane{}
}
