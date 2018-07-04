package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBRegistrationRepo struct {
  db *sql.DB
}

func (r *DBRegistrationRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "registration", domain.Registration{})
}

func (r *DBRegistrationRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "registration", domain.Registration{}, dryrun)
}

func (r *DBRegistrationRepo) FindByID(ID string) (*domain.Registration, error) {
  registration := &domain.Registration{}
  sql, targets := structsql.FindByIDSql("registration", registration)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return registration, nil
}

func (r *DBRegistrationRepo) Save(registration *domain.Registration) (string, error) {
  if (registration.ID == "") {
    registration.ID = structsql.UniqueID(r.db, "registration", "Reg1")
  }
  return registration.ID, structsql.Insert(r.db, "registration", registration, registration.ID)
}

func (r *DBRegistrationRepo) List(offset, limit int) ([]*domain.Registration, error) {
  registration := &domain.Registration{}
  registrations := make([]*domain.Registration, 0)
  sql, targets := structsql.ListSql("registration", registration, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    registrationCopy := domain.Registration(*registration)
    registrations = append(registrations, &registrationCopy)
  })
  return registrations, err
}

func (r *DBRegistrationRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "registration", ID)
}

func (r *DBRegistrationRepo) UpdateByID(ID string, oldRegistration, newRegistration *domain.Registration, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "registration", diffs.Modified(), ID)
}

func (r *DBRegistrationRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "registration", &domain.Registration{})
}
