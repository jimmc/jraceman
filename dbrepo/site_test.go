package dbrepo_test

import (
  "bytes"
  "testing"

  "github.com/jimmc/jracemango/dbrepo"
  "github.com/jimmc/jracemango/dbrepo/dbtest"
  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type siteDiffs struct {}
func (d *siteDiffs) Modified() map[string]interface{} {
  return map[string]interface{}{
    "Name": "Site FOUR",
  }
}

func TestSiteCreateTable(t *testing.T) {
  sql := structsql.CreateTableSql("site", domain.Site{})
  if got, want := sql, "CREATE TABLE site(id string primary key, name string not null, street string, street2 string, city string, state string, zip string, country string, phone string, fax string);"; got != want {
    t.Errorf("Create site table: got %v, want %v", got, want)
  }
}

func TestSiteHappyPath(t *testing.T) {
  dbr, err := dbtest.ReposEmpty()
  if err != nil {
    t.Fatalf("Error opening test database: %v", err)
  }
  defer dbr.Close()
  siteRepo := dbr.Site().(*dbrepo.DBSiteRepo)

  if err := siteRepo.CreateTable(); err != nil {
    t.Fatalf("CreateTable failed: %v", err)
  }
  if err := populateSite(dbr); err != nil {
    t.Fatalf("Populate failed: %v", err)
  }

  sites, err := siteRepo.List(0, 4)
  if err != nil {
    t.Errorf("List failed: %v", err)
  }
  if got, want := len(sites), 3; got != want {
    t.Errorf("List count: got %d, want %d", got, want)
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

  e := ixport.NewExporter(dbr.DB())
  buf := bytes.NewBufferString("")
  if err = siteRepo.Export(e, buf); err != nil {
    t.Errorf("Error exporting")
  }
  expectedExport := `
!table site
!columns "id","name","street","street2","city","state","zip","country","phone","fax"
"S1","Site One",null,null,null,null,null,null,null,null
"S2","Site Two",null,null,null,null,null,null,null,null
"S3","Site Three",null,null,null,null,null,null,null,null
`
  if got, want := buf.String(), expectedExport; got != want {
    t.Errorf("Export got %v, want %v", got, want)
  }
}

// For testing, put some data into our table
func populateSite(dbr *dbrepo.Repos) error {
  insertSql := `
INSERT into site(id, name) values
('S1', 'Site One'),
('S2', 'Site Two'),
('S3', 'Site Three')
`
  return strsql.ExecMulti(dbr.DB(), insertSql)
}
