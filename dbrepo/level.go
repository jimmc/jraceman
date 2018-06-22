package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBLevelRepo struct {
  db *sql.DB
}

func (r *DBLevelRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "level", domain.Level{})
}

func (r *DBLevelRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "level", domain.Level{}, dryrun)
}

func (r *DBLevelRepo) FindByID(ID string) (*domain.Level, error) {
  level := &domain.Level{}
  sql, targets := structsql.FindByIDSql("level", level)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return level, nil
}

func (r *DBLevelRepo) Save(level *domain.Level) error {
  // TODO - generate an ID if blank
  return structsql.Insert(r.db, "level", level, level.ID)
}

func (r *DBLevelRepo) List(offset, limit int) ([]*domain.Level, error) {
  level := &domain.Level{}
  levels := make([]*domain.Level, 0)
  sql, targets := structsql.ListSql("level", level, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    levelCopy := domain.Level(*level)
    levels = append(levels, &levelCopy)
  })
  return levels, err
}

func (r *DBLevelRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "level", ID)
}

func (r *DBLevelRepo) UpdateByID(ID string, oldLevel, newLevel *domain.Level, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "level", diffs.Modified(), ID)
}

func (r *DBLevelRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "level", &domain.Level{})
}
