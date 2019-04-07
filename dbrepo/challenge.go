package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/ixport"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type DBChallengeRepo struct {
  db *sql.DB
}

func (r *DBChallengeRepo) New() interface{} {
  return domain.Challenge{}
}

func (r *DBChallengeRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "challenge", domain.Challenge{})
}

func (r *DBChallengeRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "challenge", domain.Challenge{}, dryrun)
}

func (r *DBChallengeRepo) FindByID(ID string) (*domain.Challenge, error) {
  challenge := &domain.Challenge{}
  sql, targets := structsql.FindByIDSql("challenge", challenge)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return challenge, nil
}

func (r *DBChallengeRepo) Save(challenge *domain.Challenge) (string, error) {
  if (challenge.ID == "") {
    challenge.ID = structsql.UniqueID(r.db, "challenge", "CHL1")
  }
  return challenge.ID, structsql.Insert(r.db, "challenge", challenge, challenge.ID)
}

func (r *DBChallengeRepo) List(offset, limit int) ([]*domain.Challenge, error) {
  challenge := &domain.Challenge{}
  challenges := make([]*domain.Challenge, 0)
  sql, targets := structsql.ListSql("challenge", challenge, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    challengeCopy := domain.Challenge(*challenge)
    challenges = append(challenges, &challengeCopy)
  })
  return challenges, err
}

func (r *DBChallengeRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "challenge", ID)
}

func (r *DBChallengeRepo) UpdateByID(ID string, oldChallenge, newChallenge *domain.Challenge, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "challenge", diffs.Modified(), ID)
}

func (r *DBChallengeRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "challenge", &domain.Challenge{})
}
