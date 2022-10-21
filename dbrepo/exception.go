package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBExceptionRepo struct {
  db *sql.DB
}

func (r *DBExceptionRepo) New() interface{} {
  return domain.Exception{}
}

func (r *DBExceptionRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "exception", domain.Exception{})
}

func (r *DBExceptionRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "exception", domain.Exception{}, dryrun)
}

func (r *DBExceptionRepo) FindByID(ID string) (*domain.Exception, error) {
  exception := &domain.Exception{}
  sql, targets := structsql.FindByIDSql("exception", exception)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return exception, nil
}

func (r *DBExceptionRepo) Save(exception *domain.Exception) (string, error) {
  if (exception.ID == "") {
    exception.ID = structsql.UniqueID(r.db, "exception", "X1")
  }
  return exception.ID, structsql.Insert(r.db, "exception", exception, exception.ID)
}

func (r *DBExceptionRepo) List(offset, limit int) ([]*domain.Exception, error) {
  exception := &domain.Exception{}
  exceptions := make([]*domain.Exception, 0)
  sql, targets := structsql.ListSql("exception", exception, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    exceptionCopy := domain.Exception(*exception)
    exceptions = append(exceptions, &exceptionCopy)
  })
  return exceptions, err
}

func (r *DBExceptionRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "exception", ID)
}

func (r *DBExceptionRepo) UpdateByID(ID string, oldException, newException *domain.Exception, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "exception", diffs.Modified(), ID)
}

func (r *DBExceptionRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "exception", &domain.Exception{})
}
