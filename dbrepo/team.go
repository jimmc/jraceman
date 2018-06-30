package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBTeamRepo struct {
  db *sql.DB
}

func (r *DBTeamRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "team", domain.Team{})
}

func (r *DBTeamRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "team", domain.Team{}, dryrun)
}

func (r *DBTeamRepo) FindByID(ID string) (*domain.Team, error) {
  team := &domain.Team{}
  sql, targets := structsql.FindByIDSql("team", team)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return team, nil
}

func (r *DBTeamRepo) Save(team *domain.Team) (string, error) {
  if (team.ID == "") {
    team.ID = structsql.UniqueID(r.db, "team", "T1")
  }
  return team.ID, structsql.Insert(r.db, "team", team, team.ID)
}

func (r *DBTeamRepo) List(offset, limit int) ([]*domain.Team, error) {
  team := &domain.Team{}
  teams := make([]*domain.Team, 0)
  sql, targets := structsql.ListSql("team", team, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    teamCopy := domain.Team(*team)
    teams = append(teams, &teamCopy)
  })
  return teams, err
}

func (r *DBTeamRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "team", ID)
}

func (r *DBTeamRepo) UpdateByID(ID string, oldTeam, newTeam *domain.Team, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "team", diffs.Modified(), ID)
}

func (r *DBTeamRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "team", &domain.Team{})
}
