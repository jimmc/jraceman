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
  _, err := r.db.Exec(sql)
  return err
}

// Placeholder is ? for MySQL,$N for PostgreSQL,
// SQLite uses either of those, Oracle is :param1

func (r *dbSiteRepo) FindByID(ID string) (*domain.Site, error) {
  var site domain.Site
  err := r.db.QueryRow("select id, name from site where id=?", ID).Scan(&site.ID, &site.Name)
  if err != nil {
    return nil, err
  }
  return &site, nil
}

func (r *dbSiteRepo) Save(site *domain.Site) error {
  // TODO - generate an ID if blank
  sql := "insert into site(id, name) values(?, ?);"
  res, err := r.db.Exec(sql, site.ID, site.Name)
  return requireOneResult(res, err, "Inserted", "site", site.ID)
}

func (r *dbSiteRepo) DeleteByID(ID string) error {
  sql := "delete from site where id=?;"
  res, err := r.db.Exec(sql, ID)
  return requireOneResult(res, err, "Deleted", "site", ID)
}

func (r *dbSiteRepo) UpdateByID(ID string, oldSite, newSite *domain.Site, diffs domain.Diffs) error {
  sql, vals := modsToSql("site", diffs.Modified(), ID)
  res, err := r.db.Exec(sql, vals...)
  return requireOneResult(res, err, "Updated", "site", ID)
}
