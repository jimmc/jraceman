package domain

// MeetRepo describes how Meet records are loaded and saved.
type MeetRepo interface {
  FindByID(ID string) (Meet, error)
  Save(Meet) error
}

// Meet describes a sporting meet or regatta.
type Meet struct {
  ID string
  Name string
  ShortName string
  SiteID string
  StartDate string
  EndDate string
  AgeDate string
}
