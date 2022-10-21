package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBContextOptionRepo struct {
  db *sql.DB
}

func (r *DBContextOptionRepo) New() interface{} {
  return domain.ContextOption{}
}

func (r *DBContextOptionRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "contextoption", domain.ContextOption{})
}

func (r *DBContextOptionRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "contextoption", domain.ContextOption{}, dryrun)
}

func (r *DBContextOptionRepo) FindByID(ID string) (*domain.ContextOption, error) {
  contextoption := &domain.ContextOption{}
  sql, targets := structsql.FindByIDSql("contextoption", contextoption)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return contextoption, nil
}

func (r *DBContextOptionRepo) Save(contextoption *domain.ContextOption) (string, error) {
  if (contextoption.ID == "") {
    contextoption.ID = structsql.UniqueID(r.db, "contextoption", "CTX1")
  }
  return contextoption.ID, structsql.Insert(r.db, "contextoption", contextoption, contextoption.ID)
}

func (r *DBContextOptionRepo) List(offset, limit int) ([]*domain.ContextOption, error) {
  contextoption := &domain.ContextOption{}
  contextoptions := make([]*domain.ContextOption, 0)
  sql, targets := structsql.ListSql("contextoption", contextoption, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    contextoptionCopy := domain.ContextOption(*contextoption)
    contextoptions = append(contextoptions, &contextoptionCopy)
  })
  return contextoptions, err
}

func (r *DBContextOptionRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "contextoption", ID)
}

func (r *DBContextOptionRepo) UpdateByID(ID string, oldContextOption, newContextOption *domain.ContextOption, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "contextoption", diffs.Modified(), ID)
}

func (r *DBContextOptionRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "contextoption", &domain.ContextOption{})
}
