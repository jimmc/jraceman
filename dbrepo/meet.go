package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBMeetRepo struct {
  db *sql.DB
}

func (r *DBMeetRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "meet", domain.Meet{})
}

func (r *DBMeetRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "meet", domain.Meet{}, dryrun)
}

func (r *DBMeetRepo) FindByID(ID string) (*domain.Meet, error) {
  meet := &domain.Meet{}
  sql, targets := structsql.FindByIDSql("meet", meet)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return meet, nil
}

func (r *DBMeetRepo) Save(meet *domain.Meet) (string, error) {
  if (meet.ID == "") {
    meet.ID = structsql.UniqueID(r.db, "meet", "M1")
  }
  return meet.ID, structsql.Insert(r.db, "meet", meet, meet.ID)
}

func (r *DBMeetRepo) List(offset, limit int) ([]*domain.Meet, error) {
  meet := &domain.Meet{}
  meets := make([]*domain.Meet, 0)
  sql, targets := structsql.ListSql("meet", meet, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    meetCopy := domain.Meet(*meet)
    meets = append(meets, &meetCopy)
  })
  return meets, err
}

func (r *DBMeetRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "meet", ID)
}

func (r *DBMeetRepo) UpdateByID(ID string, oldMeet, newMeet *domain.Meet, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "meet", diffs.Modified(), ID)
}

func (r *DBMeetRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "meet", &domain.Meet{})
}