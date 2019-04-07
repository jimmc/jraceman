package dbrepo

import (
  "database/sql"
  "io"
  "strconv"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBEventRepo struct {
  db *sql.DB
}

func (r *DBEventRepo) New() interface{} {
  return domain.Event{}
}

func (r *DBEventRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "event", domain.Event{})
}

func (r *DBEventRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "event", domain.Event{}, dryrun)
}

func (r *DBEventRepo) FindByID(ID string) (*domain.Event, error) {
  event := &domain.Event{}
  sql, targets := structsql.FindByIDSql("event", event)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return event, nil
}

func (r *DBEventRepo) Save(event *domain.Event) (string, error) {
  if (event.ID == "") {
    baseID := event.MeetID + ".EV1"
    if event.Number != nil {
      baseID = event.MeetID + ".E" + strconv.Itoa(*event.Number)
    }
    event.ID = structsql.UniqueID(r.db, "event", baseID)
  }
  return event.ID, structsql.Insert(r.db, "event", event, event.ID)
}

func (r *DBEventRepo) List(offset, limit int) ([]*domain.Event, error) {
  event := &domain.Event{}
  events := make([]*domain.Event, 0)
  sql, targets := structsql.ListSql("event", event, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    eventCopy := domain.Event(*event)
    events = append(events, &eventCopy)
  })
  return events, err
}

func (r *DBEventRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "event", ID)
}

func (r *DBEventRepo) UpdateByID(ID string, oldEvent, newEvent *domain.Event, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "event", diffs.Modified(), ID)
}

func (r *DBEventRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "event", &domain.Event{})
}
