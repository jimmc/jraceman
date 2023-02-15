package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBRegistrationFeeRepo struct {
  db conn.DB
}

func (r *DBRegistrationFeeRepo) New() interface{} {
  return domain.RegistrationFee{}
}

func (r *DBRegistrationFeeRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "registrationfee", domain.RegistrationFee{})
}

func (r *DBRegistrationFeeRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "registrationfee", domain.RegistrationFee{}, dryrun)
}

func (r *DBRegistrationFeeRepo) FindByID(ID string) (*domain.RegistrationFee, error) {
  registrationfee := &domain.RegistrationFee{}
  sql, targets := structsql.FindByIDSql("registrationfee", registrationfee)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return registrationfee, nil
}

func (r *DBRegistrationFeeRepo) Save(registrationfee *domain.RegistrationFee) (string, error) {
  if (registrationfee.ID == "") {
    registrationfee.ID = structsql.UniqueID(r.db, "registrationfee", "RegFee1")
  }
  return registrationfee.ID, structsql.Insert(r.db, "registrationfee", registrationfee, registrationfee.ID)
}

func (r *DBRegistrationFeeRepo) List(offset, limit int) ([]*domain.RegistrationFee, error) {
  registrationfee := &domain.RegistrationFee{}
  registrationfees := make([]*domain.RegistrationFee, 0)
  sql, targets := structsql.ListSql("registrationfee", registrationfee, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    registrationfeeCopy := domain.RegistrationFee(*registrationfee)
    registrationfees = append(registrationfees, &registrationfeeCopy)
  })
  return registrationfees, err
}

func (r *DBRegistrationFeeRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "registrationfee", ID)
}

func (r *DBRegistrationFeeRepo) UpdateByID(ID string, oldRegistrationFee, newRegistrationFee *domain.RegistrationFee, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "registrationfee", diffs.Modified(), ID)
}

func (r *DBRegistrationFeeRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "registrationfee", &domain.RegistrationFee{})
}
