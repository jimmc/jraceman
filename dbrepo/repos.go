package dbrepo

import (
  "database/sql"
  "fmt"
  "io"
  "strings"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/domain"

  "github.com/golang/glog"
)

// Repos implements the domain.Repos interface.
type Repos struct {
  db *sql.DB
  tableMap map[string]TableRepo
  dbArea *DBAreaRepo
  dbChallenge *DBChallengeRepo
  dbCompetition *DBCompetitionRepo
  dbComplan *DBComplanRepo
  dbComplanRule *DBComplanRuleRepo
  dbComplanStage *DBComplanStageRepo
  dbContextOption *DBContextOptionRepo
  dbEntry *DBEntryRepo
  dbEvent *DBEventRepo
  dbException *DBExceptionRepo
  dbGender *DBGenderRepo
  dbLane *DBLaneRepo
  dbLaneOrder *DBLaneOrderRepo
  dbLevel *DBLevelRepo
  dbMeet *DBMeetRepo
  dbOption *DBOptionRepo
  dbPerson *DBPersonRepo
  dbProgression *DBProgressionRepo
  dbRace *DBRaceRepo
  dbRegistration *DBRegistrationRepo
  dbRegistrationFee *DBRegistrationFeeRepo
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
  New() interface{}     // Returns a new instance of the domain struct for this table.
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
    {"option", r.dbOption},
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
    {"registrationfee", r.dbRegistrationFee},
    {"registration", r.dbRegistration},
    {"event", r.dbEvent},
    {"entry", r.dbEntry},
    {"race", r.dbRace},
    {"lane", r.dbLane},
    {"contextoption", r.dbContextOption},
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
func (r *Repos) ContextOption() domain.ContextOptionRepo { return r.dbContextOption }
func (r *Repos) Entry() domain.EntryRepo { return r.dbEntry }
func (r *Repos) Event() domain.EventRepo { return r.dbEvent }
func (r *Repos) Exception() domain.ExceptionRepo { return r.dbException }
func (r *Repos) Gender() domain.GenderRepo { return r.dbGender }
func (r *Repos) Lane() domain.LaneRepo { return r.dbLane }
func (r *Repos) LaneOrder() domain.LaneOrderRepo { return r.dbLaneOrder }
func (r *Repos) Level() domain.LevelRepo { return r.dbLevel }
func (r *Repos) Meet() domain.MeetRepo { return r.dbMeet }
func (r *Repos) Option() domain.OptionRepo { return r.dbOption }
func (r *Repos) Person() domain.PersonRepo { return r.dbPerson }
func (r *Repos) Progression() domain.ProgressionRepo { return r.dbProgression }
func (r *Repos) Race() domain.RaceRepo { return r.dbRace }
func (r *Repos) Registration() domain.RegistrationRepo { return r.dbRegistration }
func (r *Repos) RegistrationFee() domain.RegistrationFeeRepo { return r.dbRegistrationFee }
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
  glog.Infof("Opening dbrepo type %s at %s", dbtype, dbloc)
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
    dbContextOption: &DBContextOptionRepo{db},
    dbEntry: &DBEntryRepo{db},
    dbEvent: &DBEventRepo{db},
    dbException: &DBExceptionRepo{db},
    dbGender: &DBGenderRepo{db},
    dbLane: &DBLaneRepo{db},
    dbLaneOrder: &DBLaneOrderRepo{db},
    dbLevel: &DBLevelRepo{db},
    dbMeet: &DBMeetRepo{db},
    dbOption: &DBOptionRepo{db},
    dbPerson: &DBPersonRepo{db},
    dbProgression: &DBProgressionRepo{db},
    dbRace: &DBRaceRepo{db},
    dbRegistration: &DBRegistrationRepo{db},
    dbRegistrationFee: &DBRegistrationFeeRepo{db},
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

  tableMap := make(map[string]TableRepo)
  for _, entry := range r.TableEntries() {
    tableMap[entry.Name] = entry.Table
  }
  r.tableMap = tableMap

  return r, err
}

func (r *Repos) TableRepo(name string) (TableRepo, error) {
  table, ok := r.tableMap[name]
  if !ok {
    return nil, fmt.Errorf("no such table %q", name)
  }
  return table, nil
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
func (r *Repos) Import(in io.Reader) (ixport.ImporterCounts, error) {
  rr, err := NewRowRepoWithTx(r)
  if err != nil {
    return ixport.ImporterCounts{}, err
  }
  var committed bool
  defer func() {
    if !committed {
      rr.Rollback()
    }
  }()
  im := ixport.NewImporter(rr)
  err = im.Import(in)
  if err == nil {
    err = rr.Commit()
    committed = true
  }
  return im.Counts(), err
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
