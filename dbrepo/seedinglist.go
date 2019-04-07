package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBSeedingListRepo struct {
  db *sql.DB
}

func (r *DBSeedingListRepo) New() interface{} {
  return domain.SeedingList{}
}

func (r *DBSeedingListRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "seedinglist", domain.SeedingList{})
}

func (r *DBSeedingListRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "seedinglist", domain.SeedingList{}, dryrun)
}

func (r *DBSeedingListRepo) FindByID(ID string) (*domain.SeedingList, error) {
  seedinglist := &domain.SeedingList{}
  sql, targets := structsql.FindByIDSql("seedinglist", seedinglist)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return seedinglist, nil
}

func (r *DBSeedingListRepo) Save(seedinglist *domain.SeedingList) (string, error) {
  if (seedinglist.ID == "") {
    seedinglist.ID = structsql.UniqueID(r.db, "seedinglist", "SeedL1")
  }
  return seedinglist.ID, structsql.Insert(r.db, "seedinglist", seedinglist, seedinglist.ID)
}

func (r *DBSeedingListRepo) List(offset, limit int) ([]*domain.SeedingList, error) {
  seedinglist := &domain.SeedingList{}
  seedinglists := make([]*domain.SeedingList, 0)
  sql, targets := structsql.ListSql("seedinglist", seedinglist, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    seedinglistCopy := domain.SeedingList(*seedinglist)
    seedinglists = append(seedinglists, &seedinglistCopy)
  })
  return seedinglists, err
}

func (r *DBSeedingListRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "seedinglist", ID)
}

func (r *DBSeedingListRepo) UpdateByID(ID string, oldSeedingList, newSeedingList *domain.SeedingList, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "seedinglist", diffs.Modified(), ID)
}

func (r *DBSeedingListRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "seedinglist", &domain.SeedingList{})
}
