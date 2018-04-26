package domain

// GenderRepo describes how Gender records are loaded and saved.
type GenderRepo interface {
  FindByID(ID string) (Gender, error)
  Save(Gender) error
}

// Gender is typically male or female, but we allow other values.
type Gender struct {
  ID string
  Name string
}
