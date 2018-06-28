package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBAreaRepo struct {
  db *sql.DB
}

func (r *DBAreaRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "area", domain.Area{})
}

func (r *DBAreaRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "area", domain.Area{}, dryrun)
}

func (r *DBAreaRepo) FindByID(ID string) (*domain.Area, error) {
  area := &domain.Area{}
  sql, targets := structsql.FindByIDSql("area", area)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return area, nil
}

func (r *DBAreaRepo) Save(area *domain.Area) (string, error) {
  if (area.ID == "") {
    area.ID = structsql.UniqueID(r.db, "area", "A1")
  }
  return area.ID, structsql.Insert(r.db, "area", area, area.ID)
}

func (r *DBAreaRepo) List(offset, limit int) ([]*domain.Area, error) {
  area := &domain.Area{}
  areas := make([]*domain.Area, 0)
  sql, targets := structsql.ListSql("area", area, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    areaCopy := domain.Area(*area)
    areas = append(areas, &areaCopy)
  })
  return areas, err
}

func (r *DBAreaRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "area", ID)
}

func (r *DBAreaRepo) UpdateByID(ID string, oldArea, newArea *domain.Area, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "area", diffs.Modified(), ID)
}

func (r *DBAreaRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "area", &domain.Area{})
}
