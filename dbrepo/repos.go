package dbrepo

import (
  "database/sql"
  "log"
  // _ "github.com/go-sql-driver/mysql"
  _ "github.com/proullon/ramsql/driver"

  "github.com/jimmc/jracemango/domain"
)

type Repos struct {
  db *sql.DB
  dbArea *dbAreaRepo
  dbSite *dbSiteRepo
}

func (r *Repos) Area() domain.AreaRepo {
  return r.dbArea
}

func (r *Repos) Site() domain.SiteRepo {
  return r.dbSite
}

func Open() (*Repos, error) {
  dbtype := "ramsql"
  dbloc := "TestDatabase"
  // dbtype := "mysql"
  // dbloc := "user:password@tcp(127.0.0.1:3306)/hello"
  log.Printf("Opening dbrepo type %s at %s", dbtype, dbloc)
  db, err := sql.Open(dbtype, dbloc)
  if err != nil {
    return nil, err
  }

  // Open the database for real
  err = db.Ping()
  if err != nil {
    return nil, err
  }

  r := &Repos{
    db: db,
    dbArea: &dbAreaRepo{db},
    dbSite: &dbSiteRepo{db},
  }

  // TODO - for testing, with ramsql, create tables
  err = r.CreateTables()
  if err == nil {
    r.dbSite.Populate()
  }

  return r, err
}

func (r *Repos) CreateTables() error {
  err := r.dbSite.CreateTable()
  log.Printf("Site.CreateTable error result: %v", err)
  return err
}

func (r *Repos) Close() {
  if r.db == nil {
    return
  }
  r.db.Close()
  r.db = nil
}
