package domain

// AreaRepo describes how Area records are loaded and saved.
type AreaRepo interface {
  FindByID(ID string) (*Area, error)
  List(offset, limit int) ([]*Area, error)
  Save(*Area) (string, error)
  UpdateByID(ID string, oldArea, newArea *Area, diffs Diffs) error
  DeleteByID(ID string) error
}

// Area describes an event area such as a race course.
type Area struct {
  ID string
  Name string
  SiteID string
  Lanes int
  ExtraLanes int
}
