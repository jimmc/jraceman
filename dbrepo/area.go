package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type dbAreaRepo struct {
  db *sql.DB
}

func (r *dbAreaRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "area", domain.Area{})
}

func (r *dbAreaRepo) FindByID(ID string) (*domain.Area, error) {
  area := &domain.Area{}
  sql, targets := structsql.FindByIDSql("area", area)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return area, nil
}

func (r *dbAreaRepo) Save(area *domain.Area) error {
  // TODO - generate an ID if blank
  return structsql.Insert(r.db, "area", area, area.ID)
}

func (r *dbAreaRepo) List(offset, limit int) ([]*domain.Area, error) {
  area := &domain.Area{}
  areas := make([]*domain.Area, 0)
  sql, targets := structsql.ListSql("area", area, offset, limit)
  err := structsql.QueryAndCollect(r.db, sql, targets, func() {
    areaCopy := domain.Area(*area)
    areas = append(areas, &areaCopy)
  })
  return areas, err
}

func (r *dbAreaRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "area", ID)
}

func (r *dbAreaRepo) UpdateByID(ID string, oldArea, newArea *domain.Area, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "area", diffs.Modified(), ID)
}

func (r *dbAreaRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "area", &domain.Area{})
}
