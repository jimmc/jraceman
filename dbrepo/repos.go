package dbrepo

import (
  "database/sql"
  "fmt"
  "log"
  "strings"

  "github.com/jimmc/jracemango/domain"
)

// Repos implements the domain.Repos interface.
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

// Open opens a database repository.
// The repoPath argument is of the form dbtype:dbinfo,
// such as "ramsql:TestDatabase" or "mysql:user:password@tcp(...)/hello".
// Note, however, that this package does not import any sql drivers;
// the main program should import whatever drivers it wants to use.
func Open(repoPath string) (*Repos, error) {
  colon := strings.Index(repoPath, ":")
  if colon <= 0 {
    return nil, fmt.Errorf("Bad format for repoPath, it must have a DB type followed by a colon")
  }
  dbtype := repoPath[:colon]
  dbloc := repoPath[colon+1:]
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

  return r, err
}

// CreateTables creates all of the tables in a new database.
// This method is not idempotent, it will fail if any of the
// tables already exist.
func (r *Repos) CreateTables() error {
  if err := r.dbSite.CreateTable(); err != nil {
    return fmt.Errorf("error creating Site table: %v", err)
  }

  if err := r.dbArea.CreateTable(); err != nil {
    return fmt.Errorf("error creating Area table: %v", err)
  }

  return nil
}

// Close closes the database.
func (r *Repos) Close() {
  if r.db == nil {
    return
  }
  r.db.Close()
  r.db = nil
}
