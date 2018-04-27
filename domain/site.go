package domain

// SiteRepo describes how Site records are manipulated in the repository.
type SiteRepo interface {
  FindByID(ID string) (*Site, error)
  Save(*Site) error
  DeleteByID(ID string) error
  UpdateByID(ID string, oldSite, newSite *Site, diffs Diffs) error
}

// Site describes an event venue.
type Site struct {
  ID string
  Name string
  Street *string
  Street2 *string
  City *string
  State *string
  Zip *string
  Country *string
  Phone *string
  Fax *string
}
