package dbrepo

import (
  "database/sql"

  "github.com/jimmc/jracemango/domain"
)

type dbSiteRepo struct {
  db *sql.DB
}

func (r *dbSiteRepo) CreateTable() error {
  sql := `CREATE TABLE site (
      ID string primary key,
      Name string,
      Street string,
      Street2 string,
      City string,
      State string,
      Zip string,
      Country string,
      Phone string,
      Fax string);`
  _, err := r.db.Exec(sql);
  return err
}

// For testing, put some data into our table
func (r *dbSiteRepo) Populate() error {
  columns := "INSERT into site(ID, Name) values ("
  values := [] string {
    "S1, 'Site One'",
    "S2, 'Site Two'",
    "S3, 'Site Three'",
  }
  for _, vv := range values {
    sql := columns + vv + ");"
    _, err := r.db.Exec(sql)
    if err != nil {
      return err
    }
  }
  return nil
}

func (r *dbSiteRepo) FindById(ID string) (domain.Site, error) {
  // TODO - query the database
  return r.TestFindById(ID)
}

func (r *dbSiteRepo) TestFindById(ID string) (domain.Site, error) {
  return domain.Site{
    ID: ID,
    Name: "Site-" + ID,
  }, nil
}

func (r *dbSiteRepo) Save(site domain.Site) error {
  return nil
}
