package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBPersonRepo struct {
  db conn.DB
}

func (r *DBPersonRepo) New() interface{} {
  return domain.Person{}
}

func (r *DBPersonRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "person", domain.Person{})
}

func (r *DBPersonRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "person", domain.Person{}, dryrun)
}

func (r *DBPersonRepo) FindByID(ID string) (*domain.Person, error) {
  person := &domain.Person{}
  sql, targets := structsql.FindByIDSql("person", person)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return person, nil
}

func (r *DBPersonRepo) Save(person *domain.Person) (string, error) {
  if (person.ID == "") {
    person.ID = structsql.UniqueID(r.db, "person", "P1")
  }
  return person.ID, structsql.Insert(r.db, "person", person, person.ID)
}

func (r *DBPersonRepo) List(offset, limit int) ([]*domain.Person, error) {
  person := &domain.Person{}
  persons := make([]*domain.Person, 0)
  sql, targets := structsql.ListSql("person", person, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    personCopy := domain.Person(*person)
    persons = append(persons, &personCopy)
  })
  return persons, err
}

func (r *DBPersonRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "person", ID)
}

func (r *DBPersonRepo) UpdateByID(ID string, oldPerson, newPerson *domain.Person, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "person", diffs.Modified(), ID)
}

func (r *DBPersonRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "person", &domain.Person{})
}
