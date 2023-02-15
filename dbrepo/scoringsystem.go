package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBScoringSystemRepo struct {
  db conn.DB
}

func (r *DBScoringSystemRepo) New() interface{} {
  return domain.ScoringSystem{}
}

func (r *DBScoringSystemRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "scoringsystem", domain.ScoringSystem{})
}

func (r *DBScoringSystemRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "scoringsystem", domain.ScoringSystem{}, dryrun)
}

func (r *DBScoringSystemRepo) FindByID(ID string) (*domain.ScoringSystem, error) {
  scoringsystem := &domain.ScoringSystem{}
  sql, targets := structsql.FindByIDSql("scoringsystem", scoringsystem)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return scoringsystem, nil
}

func (r *DBScoringSystemRepo) Save(scoringsystem *domain.ScoringSystem) (string, error) {
  if (scoringsystem.ID == "") {
    scoringsystem.ID = structsql.UniqueID(r.db, "scoringsystem", "ScS1")
  }
  return scoringsystem.ID, structsql.Insert(r.db, "scoringsystem", scoringsystem, scoringsystem.ID)
}

func (r *DBScoringSystemRepo) List(offset, limit int) ([]*domain.ScoringSystem, error) {
  scoringsystem := &domain.ScoringSystem{}
  scoringsystems := make([]*domain.ScoringSystem, 0)
  sql, targets := structsql.ListSql("scoringsystem", scoringsystem, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    scoringsystemCopy := domain.ScoringSystem(*scoringsystem)
    scoringsystems = append(scoringsystems, &scoringsystemCopy)
  })
  return scoringsystems, err
}

func (r *DBScoringSystemRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "scoringsystem", ID)
}

func (r *DBScoringSystemRepo) UpdateByID(ID string, oldScoringSystem, newScoringSystem *domain.ScoringSystem, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "scoringsystem", diffs.Modified(), ID)
}

func (r *DBScoringSystemRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "scoringsystem", &domain.ScoringSystem{})
}
