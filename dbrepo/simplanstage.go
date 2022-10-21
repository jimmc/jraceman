package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBSimplanStageRepo struct {
  db *sql.DB
}

func (r *DBSimplanStageRepo) New() interface{} {
  return domain.SimplanStage{}
}

func (r *DBSimplanStageRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "simplanstage", domain.SimplanStage{})
}

func (r *DBSimplanStageRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "simplanstage", domain.SimplanStage{}, dryrun)
}

func (r *DBSimplanStageRepo) FindByID(ID string) (*domain.SimplanStage, error) {
  simplanstage := &domain.SimplanStage{}
  sql, targets := structsql.FindByIDSql("simplanstage", simplanstage)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return simplanstage, nil
}

func (r *DBSimplanStageRepo) Save(simplanstage *domain.SimplanStage) (string, error) {
  if (simplanstage.ID == "") {
    simplanstage.ID = structsql.UniqueID(r.db, "simplanstage", "SS1")
  }
  return simplanstage.ID, structsql.Insert(r.db, "simplanstage", simplanstage, simplanstage.ID)
}

func (r *DBSimplanStageRepo) List(offset, limit int) ([]*domain.SimplanStage, error) {
  simplanstage := &domain.SimplanStage{}
  simplanstages := make([]*domain.SimplanStage, 0)
  sql, targets := structsql.ListSql("simplanstage", simplanstage, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    simplanstageCopy := domain.SimplanStage(*simplanstage)
    simplanstages = append(simplanstages, &simplanstageCopy)
  })
  return simplanstages, err
}

func (r *DBSimplanStageRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "simplanstage", ID)
}

func (r *DBSimplanStageRepo) UpdateByID(ID string, oldSimplanStage, newSimplanStage *domain.SimplanStage, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "simplanstage", diffs.Modified(), ID)
}

func (r *DBSimplanStageRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "simplanstage", &domain.SimplanStage{})
}
