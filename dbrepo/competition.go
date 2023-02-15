package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/conn"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBCompetitionRepo struct {
  db conn.DB
}

func (r *DBCompetitionRepo) New() interface{} {
  return domain.Competition{}
}

func (r *DBCompetitionRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "competition", domain.Competition{})
}

func (r *DBCompetitionRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "competition", domain.Competition{}, dryrun)
}

func (r *DBCompetitionRepo) FindByID(ID string) (*domain.Competition, error) {
  competition := &domain.Competition{}
  sql, targets := structsql.FindByIDSql("competition", competition)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return competition, nil
}

func (r *DBCompetitionRepo) Save(competition *domain.Competition) (string, error) {
  if (competition.ID == "") {
    competition.ID = structsql.UniqueID(r.db, "competition", "C1")
  }
  return competition.ID, structsql.Insert(r.db, "competition", competition, competition.ID)
}

func (r *DBCompetitionRepo) List(offset, limit int) ([]*domain.Competition, error) {
  competition := &domain.Competition{}
  competitions := make([]*domain.Competition, 0)
  sql, targets := structsql.ListSql("competition", competition, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    competitionCopy := domain.Competition(*competition)
    competitions = append(competitions, &competitionCopy)
  })
  return competitions, err
}

func (r *DBCompetitionRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "competition", ID)
}

func (r *DBCompetitionRepo) UpdateByID(ID string, oldCompetition, newCompetition *domain.Competition, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "competition", diffs.Modified(), ID)
}

func (r *DBCompetitionRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "competition", &domain.Competition{})
}
