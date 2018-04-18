package domain

// AreaRepo describes how Area records are loaded and saved.
type AreaRepo interace {
  FindById(ID string) (Area, error)
  Save(Area) error
}

// Area describes an event area such as a race course.
type Area struct {
  ID string
  Name string
  SiteID string
  Lanes int
  ExtraLanes int
}

func (a *Area) Site() (Site, error) {
  return nil, nil       // TODO - need a SiteRepo
}
