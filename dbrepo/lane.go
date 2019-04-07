package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBLaneRepo struct {
  db *sql.DB
}

func (r *DBLaneRepo) New() interface{} {
  return domain.Lane{}
}

func (r *DBLaneRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "lane", domain.Lane{})
}

func (r *DBLaneRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "lane", domain.Lane{}, dryrun)
}

func (r *DBLaneRepo) FindByID(ID string) (*domain.Lane, error) {
  lane := &domain.Lane{}
  sql, targets := structsql.FindByIDSql("lane", lane)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return lane, nil
}

func (r *DBLaneRepo) Save(lane *domain.Lane) (string, error) {
  if (lane.ID == "") {
    baseID := lane.RaceID + "-" + lane.EntryID
    lane.ID = structsql.UniqueID(r.db, "lane", baseID)
  }
  return lane.ID, structsql.Insert(r.db, "lane", lane, lane.ID)
}

func (r *DBLaneRepo) List(offset, limit int) ([]*domain.Lane, error) {
  lane := &domain.Lane{}
  lanes := make([]*domain.Lane, 0)
  sql, targets := structsql.ListSql("lane", lane, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    laneCopy := domain.Lane(*lane)
    lanes = append(lanes, &laneCopy)
  })
  return lanes, err
}

func (r *DBLaneRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "lane", ID)
}

func (r *DBLaneRepo) UpdateByID(ID string, oldLane, newLane *domain.Lane, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "lane", diffs.Modified(), ID)
}

func (r *DBLaneRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "lane", &domain.Lane{})
}
