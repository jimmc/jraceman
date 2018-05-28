package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBCompetitionRepo struct {
  db *sql.DB
}

func (r *DBCompetitionRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "competition", domain.Competition{})
}

func (r *DBCompetitionRepo) FindByID(ID string) (*domain.Competition, error) {
  competition := &domain.Competition{}
  sql, targets := structsql.FindByIDSql("competition", competition)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return competition, nil
}

func (r *DBCompetitionRepo) Save(competition *domain.Competition) error {
  // TODO - generate an ID if blank
  return structsql.Insert(r.db, "competition", competition, competition.ID)
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
