package dbrepo

import (
  "database/sql"
  "fmt"

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
  var site domain.Site
  err := r.db.QueryRow("select id, name from site where id=?", ID).Scan(&site.ID, &site.Name)
    // Placeholder is ? for MySQL,$N for PostgreSQL,
    // SQLite uses either of those, Oracle is :param1
  if err != nil {
    return nil, err
  }
  return &site, nil
}

func (r *dbSiteRepo) Save(site *domain.Site) error {
  // TODO - generate an ID if blank
  sql := "insert into site(id, name) values(?, ?);"
  stmt, err := r.db.Prepare(sql)        // TODO - do this in an init phase
  if err != nil {
    return err
  }
  defer stmt.Close()
  res, err := stmt.Exec(site.ID, site.Name)
  if err != nil {
    return err
  }
  rowCnt, err := res.RowsAffected()
  if err != nil {
    return err
  }
  if rowCnt != 1 {
    return fmt.Errorf("Row %s[%s] was not properly inserted", "site", site.ID);
  }
  return nil
}

func (r *dbSiteRepo) DeleteById(ID string) error {
  sql := "delete from site where id=?;"
  stmt, err := r.db.Prepare(sql)        // TODO - do this in an init phase
  if err != nil {
    return err
  }
  defer stmt.Close()
  res, err := stmt.Exec(ID)
  if err != nil {
    return err
  }
  rowCnt, err := res.RowsAffected()
  if err != nil {
    return err
  }
  if rowCnt == 0 {
    return fmt.Errorf("No %s rows deleted", "site")
  }
  if rowCnt > 1 {
    return fmt.Errorf("Deleted %d %s rows for ID %s", rowCnt, "site", ID)
  }
  return nil
}
