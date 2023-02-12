package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/compat"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBComplanRepo struct {
  db compat.DBorTx
}

func (r *DBComplanRepo) New() interface{} {
  return domain.Complan{}
}

func (r *DBComplanRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "complan", domain.Complan{})
}

func (r *DBComplanRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "complan", domain.Complan{}, dryrun)
}

func (r *DBComplanRepo) FindByID(ID string) (*domain.Complan, error) {
  complan := &domain.Complan{}
  sql, targets := structsql.FindByIDSql("complan", complan)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return complan, nil
}

func (r *DBComplanRepo) Save(complan *domain.Complan) (string, error) {
  if (complan.ID == "") {
    complan.ID = structsql.UniqueID(r.db, "complan", "CP1")
  }
  return complan.ID, structsql.Insert(r.db, "complan", complan, complan.ID)
}

func (r *DBComplanRepo) List(offset, limit int) ([]*domain.Complan, error) {
  complan := &domain.Complan{}
  complans := make([]*domain.Complan, 0)
  sql, targets := structsql.ListSql("complan", complan, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    complanCopy := domain.Complan(*complan)
    complans = append(complans, &complanCopy)
  })
  return complans, err
}

func (r *DBComplanRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "complan", ID)
}

func (r *DBComplanRepo) UpdateByID(ID string, oldComplan, newComplan *domain.Complan, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "complan", diffs.Modified(), ID)
}

func (r *DBComplanRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "complan", &domain.Complan{})
}
