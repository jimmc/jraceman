package domain

// SeedingListRepo describes how SeedingList records are loaded and saved.
type SeedingListRepo interface {
  FindByID(ID string) (*SeedingList, error)
  List(offset, limit int) ([]*SeedingList, error)
  Save(*SeedingList) (string, error)
  UpdateByID(ID string, oldSeedingList, newSeedingList *SeedingList, diffs Diffs) error
  DeleteByID(ID string) error
}

// SeedingList describes the order in which to seed specific competitors.
type SeedingList struct {
  ID string
  SeedingPlanID string
  Rank int
  PersonID string
}

// SeedingListMeta provides funcions related to the SeedingList struct.
type SeedingListMeta struct {}

func (m *SeedingListMeta) EntityTypeName() string {
  return "seedinglist"
}

func (m *SeedingListMeta) NewEntity() interface{} {
  return &SeedingList{}
}
