package dbrepo

import (
  "database/sql"

  "github.com/jimmc/jracemango/domain"
)

type dbAreaRepo struct {
  db *sql.DB
}

func (r *dbAreaRepo) FindById(ID string) (domain.Area, error) {
  return domain.Area{
    ID: ID,
    Name: "Area-" + ID,
  }, nil
}
func (r *dbAreaRepo) Save(site domain.Area) error {
  return nil
}
