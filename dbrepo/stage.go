package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBStageRepo struct {
  db *sql.DB
}

func (r *DBStageRepo) New() interface{} {
  return domain.Stage{}
}

func (r *DBStageRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "stage", domain.Stage{})
}

func (r *DBStageRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "stage", domain.Stage{}, dryrun)
}

func (r *DBStageRepo) FindByID(ID string) (*domain.Stage, error) {
  stage := &domain.Stage{}
  sql, targets := structsql.FindByIDSql("stage", stage)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return stage, nil
}

func (r *DBStageRepo) Save(stage *domain.Stage) (string, error) {
  if (stage.ID == "") {
    stage.ID = structsql.UniqueID(r.db, "stage", "S1")
  }
  return stage.ID, structsql.Insert(r.db, "stage", stage, stage.ID)
}

func (r *DBStageRepo) List(offset, limit int) ([]*domain.Stage, error) {
  stage := &domain.Stage{}
  stages := make([]*domain.Stage, 0)
  sql, targets := structsql.ListSql("stage", stage, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    stageCopy := domain.Stage(*stage)
    stages = append(stages, &stageCopy)
  })
  return stages, err
}

func (r *DBStageRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "stage", ID)
}

func (r *DBStageRepo) UpdateByID(ID string, oldStage, newStage *domain.Stage, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "stage", diffs.Modified(), ID)
}

func (r *DBStageRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "stage", &domain.Stage{})
}
