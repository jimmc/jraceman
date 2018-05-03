package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/domain"
)

type dbSiteRepo struct {
  db *sql.DB
}

func (r *dbSiteRepo) CreateTable() error {
  sql := stdCreateTableSqlFromStruct("site", domain.Site{})
  _, err := r.db.Exec(sql)
  return err
}

func (r *dbSiteRepo) FindByID(ID string) (*domain.Site, error) {
  site := &domain.Site{}
  sql, targets := stdFindByIDSqlFromStruct("site", site)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return site, nil
}

func (r *dbSiteRepo) Save(site *domain.Site) error {
  // TODO - generate an ID if blank
  sql, values := stdInsertSqlFromStruct("site", site)
  res, err := r.db.Exec(sql, values...)
  return requireOneResult(res, err, "Inserted", "site", site.ID)
}

func (r *dbSiteRepo) List(offset, limit int) ([]*domain.Site, error) {
  site := &domain.Site{}
  sites := make([]*domain.Site, 0)
  sql, targets := stdListSqlFromStruct("site", site, offset, limit)
  err := stdQueryAndCollect(r.db, sql, targets, func() {
    siteCopy := domain.Site(*site)
    sites = append(sites, &siteCopy)
  })
  return sites, err
}

func (r *dbSiteRepo) DeleteByID(ID string) error {
  sql := stdDeleteByIDSql("site")
  res, err := r.db.Exec(sql, ID)
  return requireOneResult(res, err, "Deleted", "site", ID)
}

func (r *dbSiteRepo) UpdateByID(ID string, oldSite, newSite *domain.Site, diffs domain.Diffs) error {
  sql, vals := modsToSql("site", diffs.Modified(), ID)
  res, err := r.db.Exec(sql, vals...)
  return requireOneResult(res, err, "Updated", "site", ID)
}

func (r *dbSiteRepo) Export(dbr *Repos, w io.Writer) error {
  site := &domain.Site{}
  return dbr.exportTableFromStruct(w, "site", site)
}
