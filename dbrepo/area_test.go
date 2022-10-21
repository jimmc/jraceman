package dbrepo_test

import (
  "bytes"
  "testing"

  "github.com/jimmc/jraceman/dbrepo"
  dbtest "github.com/jimmc/jraceman/dbrepo/test"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type areaDiffs struct {}
func (d *areaDiffs) Modified() map[string]interface{} {
  return map[string]interface{}{
    "Name": "Area FOUR",
  }
}

func TestAreaCreateTable(t *testing.T) {
  sql := structsql.CreateTableSql("site", domain.Area{})
  if got, want := sql, "CREATE TABLE site(id string primary key, name string not null, siteid string not null references site(id), lanes int not null, extralanes int not null);"; got != want {
    t.Errorf("Create site table: got %v, want %v", got, want)
  }
}

func TestAreaHappyPath(t *testing.T) {
  dbr, err := dbtest.ReposEmpty()
  if err != nil {
    t.Fatalf("Error opening test database: %v", err)
  }
  defer dbr.Close()
  areaRepo := dbr.Area().(*dbrepo.DBAreaRepo)

  if err := areaRepo.CreateTable(); err != nil {
    t.Fatalf("CreateTable failed: %v", err)
  }
  if err := populateArea(dbr); err != nil {
    t.Fatalf("Populate failed: %v", err)
  }

  areas, err := areaRepo.List(0, 4)
  if err != nil {
    t.Errorf("List failed: %v", err)
  }
  if got, want := len(areas), 3; got != want {
    t.Errorf("List count: got %d, want %d", got, want)
  }

  area, err := areaRepo.FindByID("A4")
  if err == nil {
    t.Errorf("Did not get error as expected from FindByID %s", "A4")
  }

  newArea := &domain.Area{
    ID: "A4",
    Name: "Area Four",
    SiteID: "S1",
  }
  id, err := areaRepo.Save(newArea)
  if err != nil {
    t.Errorf("Error saving new area record")
  }
  if got, want := id, "A4"; got != want {
    t.Errorf("ID after save: got %v, want %v", got, want)
  }

  area, err = areaRepo.FindByID("A4")
  if err != nil {
    t.Fatalf("Error retrieving just-added area")
  }
  if got, want := area.Name, newArea.Name; got != want {
    t.Errorf("After save, got name %s, expecting %s", got, want)
  }

  newArea.Name = "Area FOUR"
  if err := areaRepo.UpdateByID("A4", area, newArea, &areaDiffs{}); err != nil {
    t.Errorf("Error updating A4: %v", err)
  }
  area, err = areaRepo.FindByID("A4")
  if err != nil {
    t.Fatalf("Error retrieving just-updated area")
  }
  if got, want := area.Name, newArea.Name; got != want {
    t.Errorf("After update, got name %s, expecting %s", got, want)
  }

  if err := areaRepo.DeleteByID("A4"); err != nil {
    t.Errorf("Error deleting A4: %v", err)
  }

  area, err = areaRepo.FindByID("A4")
  if err == nil {
    t.Errorf("Still found A4 after deleting it")
  }

  e := ixport.NewExporter(dbr.DB())
  buf := bytes.NewBufferString("")
  if err = areaRepo.Export(e, buf); err != nil {
    t.Errorf("Error exporting")
  }
  expectedExport := `
!table area
!columns "id","name","siteid","lanes","extralanes"
"A1","Area One","S1",9,0
"A2","Area Two","S1",9,2
"A3","Area Three","S2",100,0
`
  if got, want := buf.String(), expectedExport; got != want {
    t.Errorf("Export got %v, want %v", got, want)
  }
}

// For testing, put some data into our table
func populateArea(dbr *dbrepo.Repos) error {
  insertSql := `
INSERT into area(id, name, siteid, lanes, extralanes) values
('A1', 'Area One', 'S1', 9, 0),
('A2', 'Area Two', 'S1', 9, 2),
('A3', 'Area Three', 'S2', 100, 0)
`
  return strsql.ExecMulti(dbr.DB(), insertSql)
}
