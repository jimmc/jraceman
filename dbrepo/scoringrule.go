package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBScoringRuleRepo struct {
  db *sql.DB
}

func (r *DBScoringRuleRepo) New() interface{} {
  return domain.ScoringRule{}
}

func (r *DBScoringRuleRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "scoringrule", domain.ScoringRule{})
}

func (r *DBScoringRuleRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "scoringrule", domain.ScoringRule{}, dryrun)
}

func (r *DBScoringRuleRepo) FindByID(ID string) (*domain.ScoringRule, error) {
  scoringrule := &domain.ScoringRule{}
  sql, targets := structsql.FindByIDSql("scoringrule", scoringrule)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return scoringrule, nil
}

func (r *DBScoringRuleRepo) Save(scoringrule *domain.ScoringRule) (string, error) {
  if (scoringrule.ID == "") {
    scoringrule.ID = structsql.UniqueID(r.db, "scoringrule", "ScR1")
  }
  return scoringrule.ID, structsql.Insert(r.db, "scoringrule", scoringrule, scoringrule.ID)
}

func (r *DBScoringRuleRepo) List(offset, limit int) ([]*domain.ScoringRule, error) {
  scoringrule := &domain.ScoringRule{}
  scoringrules := make([]*domain.ScoringRule, 0)
  sql, targets := structsql.ListSql("scoringrule", scoringrule, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    scoringruleCopy := domain.ScoringRule(*scoringrule)
    scoringrules = append(scoringrules, &scoringruleCopy)
  })
  return scoringrules, err
}

func (r *DBScoringRuleRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "scoringrule", ID)
}

func (r *DBScoringRuleRepo) UpdateByID(ID string, oldScoringRule, newScoringRule *domain.ScoringRule, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "scoringrule", diffs.Modified(), ID)
}

func (r *DBScoringRuleRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "scoringrule", &domain.ScoringRule{})
}
