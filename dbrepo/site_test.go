package dbrepo

import (
  "log"
  "testing"

  "github.com/jimmc/jracemango/domain"

  _ "github.com/proullon/ramsql/driver"
)

func TestSiteHappyPath(t *testing.T) {
  dbr, err := Open("ramsql:TestDatabase")
  if err != nil {
    t.Fatalf("Error opening database: %v", err)
  }
  defer dbr.Close()
  siteRepo := dbr.Site().(*dbSiteRepo)

  if err := siteRepo.CreateTable(); err != nil {
    t.Fatalf("CreateTable failed: %v", err)
  }
  if err := siteRepo.Populate(); err != nil {
    t.Fatalf("Populate failed: %v", err)
  }
  site, err := siteRepo.FindByID("S4")
  if err == nil {
    t.Errorf("Did not get error as expected from FindByID %s", "S4")
  }
  log.Printf("For not-found S1, site is %v", site)

  newSite := &domain.Site{
    ID: "S4",
    Name: "Site Four",
  }
  if err := siteRepo.Save(newSite); err != nil {
    t.Errorf("Error saving new site record")
  }

  site, err = siteRepo.FindByID("S4")
  if err != nil {
    t.Fatalf("Error retrieving just-added site")
  }
  if got, want := site.Name, newSite.Name; got != want {
    t.Errorf("Got name %s, expecting %s", got, want)
  }

  // TODO - delete it, make sure FindByID doesn't still find it
}

// For testing, put some data into our table
func (r *dbSiteRepo) Populate() error {
  columns := "INSERT into site(id, name) values ("
  values := [] string {
    "S1, 'Site One'",
    "S2, 'Site Two'",
    "S3, 'Site Three'",
  }
  for _, vv := range values {
    sql := columns + vv + ");"
    _, err := r.db.Exec(sql)
    if err != nil {
      return err
    }
  }
  return nil
}
