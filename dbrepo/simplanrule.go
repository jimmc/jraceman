package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBSimplanRuleRepo struct {
  db conn.DB
}

func (r *DBSimplanRuleRepo) UpdateColumnInfos(columnInfos []structsql.ColumnInfo) []structsql.ColumnInfo {
  // We have two foreign key columns to the stage table, so we need
  // to special-case our columninfo creation.
  for c, col := range columnInfos {
    if col.Name == "fromstageid" || col.Name == "tostageid" {
      columnInfos[c].FKTable = "stage"
    }
  }
  return columnInfos
}

func (r *DBSimplanRuleRepo) New() interface{} {
  return domain.SimplanRule{}
}

func (r *DBSimplanRuleRepo) CreateTable() error {
  return structsql.CreateTableWithUpdater(r.db, "simplanrule", domain.SimplanRule{}, r)
}

func (r *DBSimplanRuleRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTableWithUpdater(r.db, "simplanrule", domain.SimplanRule{}, dryrun, r)
}

func (r *DBSimplanRuleRepo) FindByID(ID string) (*domain.SimplanRule, error) {
  simplanrule := &domain.SimplanRule{}
  sql, targets := structsql.FindByIDSql("simplanrule", simplanrule)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return simplanrule, nil
}

func (r *DBSimplanRuleRepo) Save(simplanrule *domain.SimplanRule) (string, error) {
  if (simplanrule.ID == "") {
    simplanrule.ID = structsql.UniqueID(r.db, "simplanrule", "SR1")
  }
  return simplanrule.ID, structsql.Insert(r.db, "simplanrule", simplanrule, simplanrule.ID)
}

func (r *DBSimplanRuleRepo) List(offset, limit int) ([]*domain.SimplanRule, error) {
  simplanrule := &domain.SimplanRule{}
  simplanrules := make([]*domain.SimplanRule, 0)
  sql, targets := structsql.ListSql("simplanrule", simplanrule, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    simplanruleCopy := domain.SimplanRule(*simplanrule)
    simplanrules = append(simplanrules, &simplanruleCopy)
  })
  return simplanrules, err
}

func (r *DBSimplanRuleRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "simplanrule", ID)
}

func (r *DBSimplanRuleRepo) UpdateByID(ID string, oldSimplanRule, newSimplanRule *domain.SimplanRule, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "simplanrule", diffs.Modified(), ID)
}

func (r *DBSimplanRuleRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "simplanrule", &domain.SimplanRule{})
}
