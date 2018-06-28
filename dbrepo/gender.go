package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBGenderRepo struct {
  db *sql.DB
}

func (r *DBGenderRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "gender", domain.Gender{})
}

func (r *DBGenderRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "gender", domain.Gender{}, dryrun)
}

func (r *DBGenderRepo) FindByID(ID string) (*domain.Gender, error) {
  gender := &domain.Gender{}
  sql, targets := structsql.FindByIDSql("gender", gender)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return gender, nil
}

func (r *DBGenderRepo) Save(gender *domain.Gender) (string, error) {
  if gender.ID == "" {
    gender.ID = structsql.UniqueID(r.db, "gender", "G1")
  }
  return gender.ID, structsql.Insert(r.db, "gender", gender, gender.ID)
}

func (r *DBGenderRepo) List(offset, limit int) ([]*domain.Gender, error) {
  gender := &domain.Gender{}
  genders := make([]*domain.Gender, 0)
  sql, targets := structsql.ListSql("gender", gender, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    genderCopy := domain.Gender(*gender)
    genders = append(genders, &genderCopy)
  })
  return genders, err
}

func (r *DBGenderRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "gender", ID)
}

func (r *DBGenderRepo) UpdateByID(ID string, oldGender, newGender *domain.Gender, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "gender", diffs.Modified(), ID)
}

func (r *DBGenderRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "gender", &domain.Gender{})
}
