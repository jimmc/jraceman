package dbrepo

import (
  "database/sql"

  "github.com/jimmc/jracemango/domain"
)

type dbAreaRepo struct {
  db *sql.DB
}

func (r *dbAreaRepo) CreateTable() error {
  sql := stdCreateTableSqlFromStruct("area", domain.Area{})
  _, err := r.db.Exec(sql)
  return err
}

func (r *dbAreaRepo) FindByID(ID string) (*domain.Area, error) {
  area := &domain.Area{}
  sql, targets := stdFindByIDSqlFromStruct("area", area)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return area, nil
}

func (r *dbAreaRepo) Save(area *domain.Area) error {
  // TODO - generate an ID if blank
  sql, values := stdInsertSqlFromStruct("area", area)
  res, err := r.db.Exec(sql, values...)
  return requireOneResult(res, err, "Inserted", "area", area.ID)
}

func (r *dbAreaRepo) DeleteByID(ID string) error {
  sql := stdDeleteByIDSql("area")
  res, err := r.db.Exec(sql, ID)
  return requireOneResult(res, err, "Deleted", "area", ID)
}

func (r *dbAreaRepo) UpdateByID(ID string, oldArea, newArea *domain.Area, diffs domain.Diffs) error {
  sql, vals := modsToSql("area", diffs.Modified(), ID)
  res, err := r.db.Exec(sql, vals...)
  return requireOneResult(res, err, "Updated", "area", ID)
}
