package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBSimplanRepo struct {
  db *sql.DB
}

func (r *DBSimplanRepo) New() interface{} {
  return domain.Simplan{}
}

func (r *DBSimplanRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "simplan", domain.Simplan{})
}

func (r *DBSimplanRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "simplan", domain.Simplan{}, dryrun)
}

func (r *DBSimplanRepo) FindByID(ID string) (*domain.Simplan, error) {
  simplan := &domain.Simplan{}
  sql, targets := structsql.FindByIDSql("simplan", simplan)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return simplan, nil
}

func (r *DBSimplanRepo) Save(simplan *domain.Simplan) (string, error) {
  if (simplan.ID == "") {
    simplan.ID = structsql.UniqueID(r.db, "simplan", "SP1")
  }
  return simplan.ID, structsql.Insert(r.db, "simplan", simplan, simplan.ID)
}

func (r *DBSimplanRepo) List(offset, limit int) ([]*domain.Simplan, error) {
  simplan := &domain.Simplan{}
  simplans := make([]*domain.Simplan, 0)
  sql, targets := structsql.ListSql("simplan", simplan, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    simplanCopy := domain.Simplan(*simplan)
    simplans = append(simplans, &simplanCopy)
  })
  return simplans, err
}

func (r *DBSimplanRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "simplan", ID)
}

func (r *DBSimplanRepo) UpdateByID(ID string, oldSimplan, newSimplan *domain.Simplan, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "simplan", diffs.Modified(), ID)
}

func (r *DBSimplanRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "simplan", &domain.Simplan{})
}
