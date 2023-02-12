package dbrepo

import (
  "io"
  "strconv"

  "github.com/jimmc/jraceman/dbrepo/compat"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBRaceRepo struct {
  db compat.DBorTx
}

func (r *DBRaceRepo) New() interface{} {
  return domain.Race{}
}

func (r *DBRaceRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "race", domain.Race{})
}

func (r *DBRaceRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "race", domain.Race{}, dryrun)
}

func (r *DBRaceRepo) FindByID(ID string) (*domain.Race, error) {
  race := &domain.Race{}
  sql, targets := structsql.FindByIDSql("race", race)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return race, nil
}

func (r *DBRaceRepo) Save(race *domain.Race) (string, error) {
  if (race.ID == "") {
    baseID := race.EventID
    if race.StageID != nil {
      baseID += "." + *race.StageID
    }
    if race.Section != nil {
      baseID += "." + strconv.Itoa(*race.Section)
    }
    race.ID = structsql.UniqueID(r.db, "race", baseID)
  }
  return race.ID, structsql.Insert(r.db, "race", race, race.ID)
}

func (r *DBRaceRepo) List(offset, limit int) ([]*domain.Race, error) {
  race := &domain.Race{}
  races := make([]*domain.Race, 0)
  sql, targets := structsql.ListSql("race", race, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    raceCopy := domain.Race(*race)
    races = append(races, &raceCopy)
  })
  return races, err
}

func (r *DBRaceRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "race", ID)
}

func (r *DBRaceRepo) UpdateByID(ID string, oldRace, newRace *domain.Race, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "race", diffs.Modified(), ID)
}

func (r *DBRaceRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "race", &domain.Race{})
}
