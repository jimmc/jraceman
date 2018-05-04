package dbrepo

import (
  "testing"

  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"

  _ "github.com/mattn/go-sqlite3"
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
  dbr, err := Open("sqlite3::memory:")
  if err != nil {
    t.Fatalf("Error opening in-memory database: %v", err)
  }
  defer dbr.Close()
  areaRepo := dbr.Area().(*dbAreaRepo)

  if err := areaRepo.CreateTable(); err != nil {
    t.Fatalf("CreateTable failed: %v", err)
  }
  if err := areaRepo.Populate(); err != nil {
    t.Fatalf("Populate failed: %v", err)
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
  if err := areaRepo.Save(newArea); err != nil {
    t.Errorf("Error saving new area record")
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
}

// For testing, put some data into our table
func (r *dbAreaRepo) Populate() error {
  columns := "INSERT into area(id, name, siteid, lanes, extralanes) values ("
  values := [] string {
    "'A1', 'Area One', 'S1', 9, 0",
    "'A2', 'Area Two', 'S1', 9, 2",
    "'A3', 'Area Three', 'S2', 100, 0",
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
