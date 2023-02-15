package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBProgressionRepo struct {
  db conn.DB
}

func (r *DBProgressionRepo) New() interface{} {
  return domain.Progression{}
}

func (r *DBProgressionRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "progression", domain.Progression{})
}

func (r *DBProgressionRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "progression", domain.Progression{}, dryrun)
}

func (r *DBProgressionRepo) FindByID(ID string) (*domain.Progression, error) {
  progression := &domain.Progression{}
  sql, targets := structsql.FindByIDSql("progression", progression)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return progression, nil
}

func (r *DBProgressionRepo) Save(progression *domain.Progression) (string, error) {
  if (progression.ID == "") {
    progression.ID = structsql.UniqueID(r.db, "progression", "Pr1")
  }
  return progression.ID, structsql.Insert(r.db, "progression", progression, progression.ID)
}

func (r *DBProgressionRepo) List(offset, limit int) ([]*domain.Progression, error) {
  progression := &domain.Progression{}
  progressions := make([]*domain.Progression, 0)
  sql, targets := structsql.ListSql("progression", progression, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    progressionCopy := domain.Progression(*progression)
    progressions = append(progressions, &progressionCopy)
  })
  return progressions, err
}

func (r *DBProgressionRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "progression", ID)
}

func (r *DBProgressionRepo) UpdateByID(ID string, oldProgression, newProgression *domain.Progression, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "progression", diffs.Modified(), ID)
}

func (r *DBProgressionRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "progression", &domain.Progression{})
}
