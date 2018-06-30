package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBLaneOrderRepo struct {
  db *sql.DB
}

func (r *DBLaneOrderRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "laneorder", domain.LaneOrder{})
}

func (r *DBLaneOrderRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "laneorder", domain.LaneOrder{}, dryrun)
}

func (r *DBLaneOrderRepo) FindByID(ID string) (*domain.LaneOrder, error) {
  laneorder := &domain.LaneOrder{}
  sql, targets := structsql.FindByIDSql("laneorder", laneorder)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return laneorder, nil
}

func (r *DBLaneOrderRepo) Save(laneorder *domain.LaneOrder) (string, error) {
  if (laneorder.ID == "") {
    laneorder.ID = structsql.UniqueID(r.db, "laneorder", "LO1")
  }
  return laneorder.ID, structsql.Insert(r.db, "laneorder", laneorder, laneorder.ID)
}

func (r *DBLaneOrderRepo) List(offset, limit int) ([]*domain.LaneOrder, error) {
  laneorder := &domain.LaneOrder{}
  laneorders := make([]*domain.LaneOrder, 0)
  sql, targets := structsql.ListSql("laneorder", laneorder, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    laneorderCopy := domain.LaneOrder(*laneorder)
    laneorders = append(laneorders, &laneorderCopy)
  })
  return laneorders, err
}

func (r *DBLaneOrderRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "laneorder", ID)
}

func (r *DBLaneOrderRepo) UpdateByID(ID string, oldLaneOrder, newLaneOrder *domain.LaneOrder, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "laneorder", diffs.Modified(), ID)
}

func (r *DBLaneOrderRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "laneorder", &domain.LaneOrder{})
}
