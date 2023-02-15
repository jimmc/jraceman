package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBOptionRepo struct {
  db conn.DB
}

func (r *DBOptionRepo) New() interface{} {
  return domain.Option{}
}

func (r *DBOptionRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "option", domain.Option{})
}

func (r *DBOptionRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "option", domain.Option{}, dryrun)
}

func (r *DBOptionRepo) FindByID(ID string) (*domain.Option, error) {
  // The option table has no ID field, we use Name instead.
  option := &domain.Option{}
  sql := "select name, value from option where name=?"
  if err := r.db.QueryRow(sql, ID).Scan(&option.Name, &option.Value); err != nil {
    return nil, err
  }
  return option, nil
}

func (r *DBOptionRepo) Save(option *domain.Option) (string, error) {
  // The Option table does not have an ID field
  return option.Name, structsql.Insert(r.db, "option", option, option.Name)
}

func (r *DBOptionRepo) List(offset, limit int) ([]*domain.Option, error) {
  option := &domain.Option{}
  options := make([]*domain.Option, 0)
  sql, targets := structsql.ListSql("option", option, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    optionCopy := domain.Option(*option)
    options = append(options, &optionCopy)
  })
  return options, err
}

func (r *DBOptionRepo) DeleteByID(ID string) error {
  sqlstr := "delete from option where Name=?"
  res, err := r.db.Exec(sqlstr, ID)
  return structsql.RequireOneResult(res, err, "Deleted", "option", ID)
}

func (r *DBOptionRepo) UpdateByID(ID string, oldOption, newOption *domain.Option, diffs domain.Diffs) error {
  sqlstr := "update option set Value=? where Name=?"
  res, err := r.db.Exec(sqlstr, newOption.Value, ID)
  return structsql.RequireOneResult(res, err, "Updated", "option", ID)
}

func (r *DBOptionRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "option", &domain.Option{})
}
