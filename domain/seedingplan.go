package domain

// SeedingPlanRepo describes how SeedingPlan records are loaded and saved.
type SeedingPlanRepo interface {
  FindById(ID string) (SeedingPlan, error)
  Save(SeedingPlan) error
}

// SeedingPlan describes how a the entries in an event get seeded into lanes.
type SeedingPlan struct {
  ID string
  Name string
}
