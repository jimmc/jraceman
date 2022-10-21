package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBComplanStageRepo struct {
  db *sql.DB
}

func (r *DBComplanStageRepo) New() interface{} {
  return domain.ComplanStage{}
}

func (r *DBComplanStageRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "complanstage", domain.ComplanStage{})
}

func (r *DBComplanStageRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "complanstage", domain.ComplanStage{}, dryrun)
}

func (r *DBComplanStageRepo) FindByID(ID string) (*domain.ComplanStage, error) {
  complanstage := &domain.ComplanStage{}
  sql, targets := structsql.FindByIDSql("complanstage", complanstage)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return complanstage, nil
}

func (r *DBComplanStageRepo) Save(complanstage *domain.ComplanStage) (string, error) {
  if (complanstage.ID == "") {
    complanstage.ID = structsql.UniqueID(r.db, "complanstage", "CS1")
  }
  return complanstage.ID, structsql.Insert(r.db, "complanstage", complanstage, complanstage.ID)
}

func (r *DBComplanStageRepo) List(offset, limit int) ([]*domain.ComplanStage, error) {
  complanstage := &domain.ComplanStage{}
  complanstages := make([]*domain.ComplanStage, 0)
  sql, targets := structsql.ListSql("complanstage", complanstage, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    complanstageCopy := domain.ComplanStage(*complanstage)
    complanstages = append(complanstages, &complanstageCopy)
  })
  return complanstages, err
}

func (r *DBComplanStageRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "complanstage", ID)
}

func (r *DBComplanStageRepo) UpdateByID(ID string, oldComplanStage, newComplanStage *domain.ComplanStage, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "complanstage", diffs.Modified(), ID)
}

func (r *DBComplanStageRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "complanstage", &domain.ComplanStage{})
}
