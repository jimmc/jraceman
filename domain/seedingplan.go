package domain

// SeedingPlanRepo describes how SeedingPlan records are loaded and saved.
type SeedingPlanRepo interface {
  FindByID(ID string) (*SeedingPlan, error)
  List(offset, limit int) ([]*SeedingPlan, error)
  Save(*SeedingPlan) (string, error)
  UpdateByID(ID string, oldSeedingPlan, newSeedingPlan *SeedingPlan, diffs Diffs) error
  DeleteByID(ID string) error
}

// SeedingPlan describes how the entries in an event get seeded into lanes.
type SeedingPlan struct {
  ID string
  Name string
  SeedingOrder string
}
