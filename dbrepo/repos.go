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
  dbChallenge *DBChallengeRepo
  dbCompetition *DBCompetitionRepo
  dbComplan *DBComplanRepo
  dbComplanRule *DBComplanRuleRepo
  dbComplanStage *DBComplanStageRepo
  dbException *DBExceptionRepo
  dbGender *DBGenderRepo
  dbLaneOrder *DBLaneOrderRepo
  dbLevel *DBLevelRepo
  dbMeet *DBMeetRepo
  dbPerson *DBPersonRepo
  dbProgression *DBProgressionRepo
  dbSeedingList *DBSeedingListRepo
  dbSeedingPlan *DBSeedingPlanRepo
  dbScoringRule *DBScoringRuleRepo
  dbScoringSystem *DBScoringSystemRepo
  dbSimplan *DBSimplanRepo
  dbSimplanRule *DBSimplanRuleRepo
  dbSimplanStage *DBSimplanStageRepo
  dbSite *DBSiteRepo
  dbStage *DBStageRepo
  dbTeam *DBTeamRepo
}

type TableRepo interface {
  CreateTable() error
  UpgradeTable(dryrun bool) (bool, string, error)
  Export(e *ixport.Exporter, w io.Writer) error
}

type TableEntry struct {
  Name string
  Table TableRepo
}

func (r *Repos) TableEntries() []TableEntry {
  return []TableEntry{
    // The tables in this list are ordered so that tables that are the target
    // of foreign keys are created/updated before the tables that reference them.
    {"competition", r.dbCompetition},
    {"complan", r.dbComplan},
    {"complanstage", r.dbComplanStage},
    {"complanrule", r.dbComplanRule},
    {"site", r.dbSite},
    {"area", r.dbArea},
    {"exception", r.dbException},
    {"level", r.dbLevel},
    {"stage", r.dbStage},
    {"gender", r.dbGender},
    {"simplan", r.dbSimplan},
    {"simplanstage", r.dbSimplanStage},
    {"simplanrule", r.dbSimplanRule},
    {"progression", r.dbProgression},
    {"scoringsystem", r.dbScoringSystem},
    {"scoringrule", r.dbScoringRule},
    {"laneorder", r.dbLaneOrder},
    {"challenge", r.dbChallenge},
    {"team", r.dbTeam},
    {"person", r.dbPerson},
    {"seedingplan", r.dbSeedingPlan},
    {"seedinglist", r.dbSeedingList},
    {"meet", r.dbMeet},
  }
}

func (r *Repos) DB() *sql.DB {
  return r.db
}

func (r *Repos) Area() domain.AreaRepo { return r.dbArea }
func (r *Repos) Challenge() domain.ChallengeRepo { return r.dbChallenge }
func (r *Repos) Competition() domain.CompetitionRepo { return r.dbCompetition }
func (r *Repos) Complan() domain.ComplanRepo { return r.dbComplan }
func (r *Repos) ComplanRule() domain.ComplanRuleRepo { return r.dbComplanRule }
func (r *Repos) ComplanStage() domain.ComplanStageRepo { return r.dbComplanStage }
func (r *Repos) Exception() domain.ExceptionRepo { return r.dbException }
func (r *Repos) Gender() domain.GenderRepo { return r.dbGender }
func (r *Repos) LaneOrder() domain.LaneOrderRepo { return r.dbLaneOrder }
func (r *Repos) Level() domain.LevelRepo { return r.dbLevel }
func (r *Repos) Meet() domain.MeetRepo { return r.dbMeet }
func (r *Repos) Person() domain.PersonRepo { return r.dbPerson }
func (r *Repos) Progression() domain.ProgressionRepo { return r.dbProgression }
func (r *Repos) ScoringRule() domain.ScoringRuleRepo { return r.dbScoringRule }
func (r *Repos) ScoringSystem() domain.ScoringSystemRepo { return r.dbScoringSystem }
func (r *Repos) SeedingList() domain.SeedingListRepo { return r.dbSeedingList }
func (r *Repos) SeedingPlan() domain.SeedingPlanRepo { return r.dbSeedingPlan }
func (r *Repos) Simplan() domain.SimplanRepo { return r.dbSimplan }
func (r *Repos) SimplanRule() domain.SimplanRuleRepo { return r.dbSimplanRule }
func (r *Repos) SimplanStage() domain.SimplanStageRepo { return r.dbSimplanStage }
func (r *Repos) Site() domain.SiteRepo { return r.dbSite }
func (r *Repos) Stage() domain.StageRepo { return r.dbStage }
func (r *Repos) Team() domain.TeamRepo { return r.dbTeam }

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
    dbChallenge: &DBChallengeRepo{db},
    dbCompetition: &DBCompetitionRepo{db},
    dbComplan: &DBComplanRepo{db},
    dbComplanRule: &DBComplanRuleRepo{db},
    dbComplanStage: &DBComplanStageRepo{db},
    dbException: &DBExceptionRepo{db},
    dbGender: &DBGenderRepo{db},
    dbLaneOrder: &DBLaneOrderRepo{db},
    dbLevel: &DBLevelRepo{db},
    dbMeet: &DBMeetRepo{db},
    dbPerson: &DBPersonRepo{db},
    dbProgression: &DBProgressionRepo{db},
    dbScoringRule: &DBScoringRuleRepo{db},
    dbScoringSystem: &DBScoringSystemRepo{db},
    dbSeedingList: &DBSeedingListRepo{db},
    dbSeedingPlan: &DBSeedingPlanRepo{db},
    dbSimplan: &DBSimplanRepo{db},
    dbSimplanRule: &DBSimplanRuleRepo{db},
    dbSimplanStage: &DBSimplanStageRepo{db},
    dbSite: &DBSiteRepo{db},
    dbStage: &DBStageRepo{db},
    dbTeam: &DBTeamRepo{db},
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
  for _, entry := range r.TableEntries() {
    if err := entry.Table.CreateTable(); err != nil {
      return fmt.Errorf("error creating table %s: %v", entry.Name, err)
    }
  }
  return nil
}

func (r *Repos) TableNames() []string {
  tableEntries := r.TableEntries();
  tableNames := make([]string, len(tableEntries))
  for i, entry := range tableEntries {
    tableNames[i] = entry.Name
  }
  return tableNames
}

// UpgradeTable performs a database upgrade on the named table.
// Table names are defined in TableEntries().
// If dryrun is true, then upgrade is not performed.
func (r *Repos) UpgradeTable(tableName string, dryrun bool) (bool, string, error) {
  // We don't call this method very often, and we don't expect more
  // than perhaps 30 tables, so we just do a linear search.
  tableEntries := r.TableEntries();
  for _, entry := range tableEntries {
    if entry.Name != tableName {
      continue;
    }
    return entry.Table.UpgradeTable(dryrun)
  }
  return false, "", fmt.Errorf("no such table %s", tableName)
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

  for _, entry := range r.TableEntries() {
    if err := entry.Table.Export(e, w); err != nil {
      return fmt.Errorf("error creating table %s: %v", entry.Name, err)
    }
  }
  return nil
}
