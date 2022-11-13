package domain

// LaneOrderRepo describes how LaneOrder records are loaded and saved.
type LaneOrderRepo interface {
  FindByID(ID string) (*LaneOrder, error)
  List(offset, limit int) ([]*LaneOrder, error)
  Save(*LaneOrder) (string, error)
  UpdateByID(ID string, oldLaneOrder, newLaneOrder *LaneOrder, diffs Diffs) error
  DeleteByID(ID string) error
}

// LaneOrder describes the order that lanes should be assigned.
type LaneOrder struct {
  ID string
  AreaID string
  Lane int
  Ordering int          // Can't use "order", that's a reserved word in SQL
}

// LaneOrderMeta provides funcions related to the LaneOrder struct.
type LaneOrderMeta struct {}

func (m *LaneOrderMeta) EntityTypeName() string {
  return "laneorder"
}

func (m *LaneOrderMeta) NewEntity() interface{} {
  return &LaneOrder{}
}
