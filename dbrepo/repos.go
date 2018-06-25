package dbrepo

import (
  "database/sql"
  "fmt"
  "io"
  "log"
  "strings"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/domain"
)

// Repos implements the domain.Repos interface.
type Repos struct {
  db *sql.DB
  dbArea *DBAreaRepo
  dbCompetition *DBCompetitionRepo
  dbGender *DBGenderRepo
  dbLevel *DBLevelRepo
  dbSite *DBSiteRepo
}

type SectionRepo interface {
  CreateTable() error
  UpgradeTable(dryrun bool) (bool, string, error)
}

type SectionEntry struct {
  Name string
  Section SectionRepo
}

func (r *Repos) SectionEntries() []SectionEntry {
  return []SectionEntry{
    {"area", r.dbArea},
    {"competition", r.dbCompetition},
    {"gender", r.dbGender},
    {"level", r.dbLevel},
    {"site", r.dbSite},
  }
}

func (r *Repos) DB() *sql.DB {
  return r.db
}

func (r *Repos) Area() domain.AreaRepo {
  return r.dbArea
}

func (r *Repos) Competition() domain.CompetitionRepo {
  return r.dbCompetition
}

func (r *Repos) Gender() domain.GenderRepo {
  return r.dbGender
}

func (r *Repos) Level() domain.LevelRepo {
  return r.dbLevel
}

func (r *Repos) Site() domain.SiteRepo {
  return r.dbSite
}

// Open opens a database repository.
// The repoPath argument is of the form dbtype:dbinfo,
// such as "sqlite3:/foo/bar" or "mysql:user:password@tcp(...)/hello".
// Note, however, that the dbrepo package does not import any sql drivers;
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
    dbArea: &DBAreaRepo{db},
    dbCompetition: &DBCompetitionRepo{db},
    dbGender: &DBGenderRepo{db},
    dbLevel: &DBLevelRepo{db},
    dbSite: &DBSiteRepo{db},
  }

  return r, err
}

// Close closes the database.
func (r *Repos) Close() {
  if r.db == nil {
    return
  }
  r.db.Close()
  r.db = nil
}

// CreateTables creates all of the tables in a new database.
// This method is not idempotent, it will fail if any of the
// tables already exist.
func (r *Repos) CreateTables() error {
  if err := r.dbCompetition.CreateTable(); err != nil {
    return fmt.Errorf("error creating Competition table: %v", err)
  }

  if err := r.dbLevel.CreateTable(); err != nil {
    return fmt.Errorf("error creating Level table: %v", err)
  }

  if err := r.dbSite.CreateTable(); err != nil {
    return fmt.Errorf("error creating Site table: %v", err)
  }

  if err := r.dbArea.CreateTable(); err != nil {
    return fmt.Errorf("error creating Area table: %v", err)
  }

  if err := r.dbGender.CreateTable(); err != nil {
    return fmt.Errorf("error creating Gender table: %v", err)
  }

  return nil
}

func (r *Repos) SectionNames() []string {
  sectionEntries := r.SectionEntries();
  sectionNames := make([]string, len(sectionEntries))
  for i, entry := range sectionEntries {
    sectionNames[i] = entry.Name
  }
  return sectionNames
}

// UpgradeSection performs a database upgrade on the named section.
// Section names are defined in SectionEntries().
// If dryrun is true, then upgrade is not performed.
func (r *Repos) UpgradeSection(sectionName string, dryrun bool) (bool, string, error) {
  // We don't call this method very often, and we don't expect more
  // than perhaps 30 sections, so we just do a linear search.
  sectionEntries := r.SectionEntries();
  for _, entry := range sectionEntries {
    if entry.Name != sectionName {
      continue;
    }
    return entry.Section.UpgradeTable(dryrun)
  }
  return false, "", fmt.Errorf("no such section %s", sectionName)
}

// Import reads in the specified text file and loads our tables.
func (r *Repos) Import(in io.Reader) (int, int, int, error) {
  im := ixport.NewImporter(NewRowRepo(r))
  err := im.Import(in)
  insertCount, updateCount, unchangedCount := im.Counts()
  return insertCount, updateCount, unchangedCount, err
}

// Export writes out all of our tables to a text file that can
// be loaded back in later using Import.
func (r *Repos) Export(w io.Writer) error {
  e := ixport.NewExporter(r.db)
  if err := e.ExportHeader(w); err != nil {
    return err
  }

  // The order of output of the tables is important: tables with
  // foreign keys should be after the tables the point to.

  if err := r.dbCompetition.Export(e, w); err != nil {
    return err
  }

  if err := r.dbLevel.Export(e, w); err != nil {
    return err
  }

  if err := r.dbSite.Export(e, w); err != nil {
    return err
  }

  if err := r.dbArea.Export(e, w); err != nil {
    return err
  }

  if err := r.dbGender.Export(e, w); err != nil {
    return err
  }

  return nil
}
