package dbrepo

import (
  "testing"

  "github.com/jimmc/jracemango/domain"

  _ "github.com/mattn/go-sqlite3"
)

type siteDiffs struct {}
func (d *siteDiffs) Modified() map[string]interface{} {
  return map[string]interface{}{
    "Name": "Site FOUR",
  }
}

func TestSiteCreateTable(t *testing.T) {
  sql := stdCreateTableSqlFromStruct("site", domain.Site{})
  if got, want := sql, "CREATE TABLE site(id string primary key, name string not null, street string, street2 string, city string, state string, zip string, country string, phone string, fax string);"; got != want {
    t.Errorf("Create site table: got %v, want %v", got, want)
  }
}

func TestSiteHappyPath(t *testing.T) {
  dbr, err := Open("sqlite3::memory:")
  if err != nil {
    t.Fatalf("Error opening in-memory database: %v", err)
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
    t.Errorf("After save, got name %s, expecting %s", got, want)
  }

  newSite.Name = "Site FOUR"
  if err := siteRepo.UpdateByID("S4", site, newSite, &siteDiffs{}); err != nil {
    t.Errorf("Error updating S4: %v", err)
  }
  site, err = siteRepo.FindByID("S4")
  if err != nil {
    t.Fatalf("Error retrieving just-updated site")
  }
  if got, want := site.Name, newSite.Name; got != want {
    t.Errorf("After update, got name %s, expecting %s", got, want)
  }

  if err := siteRepo.DeleteByID("S4"); err != nil {
    t.Errorf("Error deleting S4: %v", err)
  }

  site, err = siteRepo.FindByID("S4")
  if err == nil {
    t.Errorf("Still found S4 after deleting it")
  }
}

// For testing, put some data into our table
func (r *dbSiteRepo) Populate() error {
  columns := "INSERT into site(id, name) values ("
  values := [] string {
    "'S1', 'Site One'",
    "'S2', 'Site Two'",
    "'S3', 'Site Three'",
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
