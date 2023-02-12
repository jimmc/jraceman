package dbrepo

import (
  "io"

  "github.com/jimmc/jraceman/dbrepo/compat"
  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBSeedingPlanRepo struct {
  db compat.DBorTx
}

func (r *DBSeedingPlanRepo) New() interface{} {
  return domain.SeedingPlan{}
}

func (r *DBSeedingPlanRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "seedingplan", domain.SeedingPlan{})
}

func (r *DBSeedingPlanRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "seedingplan", domain.SeedingPlan{}, dryrun)
}

func (r *DBSeedingPlanRepo) FindByID(ID string) (*domain.SeedingPlan, error) {
  seedingplan := &domain.SeedingPlan{}
  sql, targets := structsql.FindByIDSql("seedingplan", seedingplan)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return seedingplan, nil
}

func (r *DBSeedingPlanRepo) Save(seedingplan *domain.SeedingPlan) (string, error) {
  if (seedingplan.ID == "") {
    seedingplan.ID = structsql.UniqueID(r.db, "seedingplan", "SeedP1")
  }
  return seedingplan.ID, structsql.Insert(r.db, "seedingplan", seedingplan, seedingplan.ID)
}

func (r *DBSeedingPlanRepo) List(offset, limit int) ([]*domain.SeedingPlan, error) {
  seedingplan := &domain.SeedingPlan{}
  seedingplans := make([]*domain.SeedingPlan, 0)
  sql, targets := structsql.ListSql("seedingplan", seedingplan, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    seedingplanCopy := domain.SeedingPlan(*seedingplan)
    seedingplans = append(seedingplans, &seedingplanCopy)
  })
  return seedingplans, err
}

func (r *DBSeedingPlanRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "seedingplan", ID)
}

func (r *DBSeedingPlanRepo) UpdateByID(ID string, oldSeedingPlan, newSeedingPlan *domain.SeedingPlan, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "seedingplan", diffs.Modified(), ID)
}

func (r *DBSeedingPlanRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "seedingplan", &domain.SeedingPlan{})
}
