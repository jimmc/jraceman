package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/compat"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBComplanRuleRepo struct {
  db compat.DBorTx
}

func (r *DBComplanRuleRepo) New() interface{} {
  return domain.ComplanRule{}
}

func (r *DBComplanRuleRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "complanrule", domain.ComplanRule{})
}

func (r *DBComplanRuleRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "complanrule", domain.ComplanRule{}, dryrun)
}

func (r *DBComplanRuleRepo) FindByID(ID string) (*domain.ComplanRule, error) {
  complanrule := &domain.ComplanRule{}
  sql, targets := structsql.FindByIDSql("complanrule", complanrule)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return complanrule, nil
}

func (r *DBComplanRuleRepo) Save(complanrule *domain.ComplanRule) (string, error) {
  if (complanrule.ID == "") {
    complanrule.ID = structsql.UniqueID(r.db, "complanrule", "CR1")
  }
  return complanrule.ID, structsql.Insert(r.db, "complanrule", complanrule, complanrule.ID)
}

func (r *DBComplanRuleRepo) List(offset, limit int) ([]*domain.ComplanRule, error) {
  complanrule := &domain.ComplanRule{}
  complanrules := make([]*domain.ComplanRule, 0)
  sql, targets := structsql.ListSql("complanrule", complanrule, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    complanruleCopy := domain.ComplanRule(*complanrule)
    complanrules = append(complanrules, &complanruleCopy)
  })
  return complanrules, err
}

func (r *DBComplanRuleRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "complanrule", ID)
}

func (r *DBComplanRuleRepo) UpdateByID(ID string, oldComplanRule, newComplanRule *domain.ComplanRule, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "complanrule", diffs.Modified(), ID)
}

func (r *DBComplanRuleRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "complanrule", &domain.ComplanRule{})
}
