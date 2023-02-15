package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBEntryRepo struct {
  db conn.DB
}

func (r *DBEntryRepo) New() interface{} {
  return domain.Entry{}
}

func (r *DBEntryRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "entry", domain.Entry{})
}

func (r *DBEntryRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "entry", domain.Entry{}, dryrun)
}

func (r *DBEntryRepo) FindByID(ID string) (*domain.Entry, error) {
  entry := &domain.Entry{}
  sql, targets := structsql.FindByIDSql("entry", entry)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return entry, nil
}

func (r *DBEntryRepo) Save(entry *domain.Entry) (string, error) {
  if (entry.ID == "") {
    meetID := "Z"
    if entry.EventID != "" {
      // TODO - get the MeetID from this event
    }
    entry.ID = structsql.UniqueID(r.db, "entry", meetID + ".EN1")
  }
  return entry.ID, structsql.Insert(r.db, "entry", entry, entry.ID)
}

func (r *DBEntryRepo) List(offset, limit int) ([]*domain.Entry, error) {
  entry := &domain.Entry{}
  entrys := make([]*domain.Entry, 0)
  sql, targets := structsql.ListSql("entry", entry, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    entryCopy := domain.Entry(*entry)
    entrys = append(entrys, &entryCopy)
  })
  return entrys, err
}

func (r *DBEntryRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "entry", ID)
}

func (r *DBEntryRepo) UpdateByID(ID string, oldEntry, newEntry *domain.Entry, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "entry", diffs.Modified(), ID)
}

func (r *DBEntryRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "entry", &domain.Entry{})
}
