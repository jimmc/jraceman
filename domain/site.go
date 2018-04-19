package domain

// SiteRepo describes how Site records are loaded and saved.
type SiteRepo interface {
  FindById(ID string) (Site, error)
  Save(Site) error
}

// Site describes an event venue.
type Site struct {
  ID string
  Name string
  Street string
  Street2 string
  City string
  State string
  Zip string
  Country string
  Phone string
  Fax string
}
