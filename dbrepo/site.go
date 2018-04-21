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
      id string primary key,
      name string,
      street string,
      street2 string,
      city string,
      state string,
      zip string,
      country string,
      phone string,
      fax string);`
  _, err := r.db.Exec(sql);
  return err
}

// For testing, put some data into our table
func (r *dbSiteRepo) Populate() error {
  columns := "INSERT into site(id, name) values ("
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

func (r *dbSiteRepo) FindById(ID string) (*domain.Site, error) {
  rows, err := r.db.Query("select id, name from site where id=?", ID)
    // Placeholder is ? for MySQL,$N for PostgreSQL,
    // SQLite uses either of those, Oracle is :param1
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  var site domain.Site
  if !rows.Next() {
    return nil, nil     // No data and no error
  }
  err = rows.Scan(&site.ID, &site.Name)
  if err != nil {
    return nil, err
  }
  // TODO - could call rows.Next again and return error if it is true
  // (which means more than one matching row)
  err = rows.Err()
  if err != nil {
    return nil, err
  }
  return &site, nil
}

func (r *dbSiteRepo) TestFindById(ID string) (*domain.Site, error) {
  return &domain.Site{
    ID: ID,
    Name: "Site-" + ID,
  }, nil
}

func (r *dbSiteRepo) Save(site *domain.Site) error {
  return nil
}
