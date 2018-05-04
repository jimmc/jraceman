package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jracemango/dbrepo/structsql"
  "github.com/jimmc/jracemango/domain"
)

type dbSiteRepo struct {
  db *sql.DB
}

func (r *dbSiteRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "site", domain.Site{})
}

func (r *dbSiteRepo) FindByID(ID string) (*domain.Site, error) {
  site := &domain.Site{}
  sql, targets := structsql.FindByIDSql("site", site)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return site, nil
}

func (r *dbSiteRepo) Save(site *domain.Site) error {
  // TODO - generate an ID if blank
  return structsql.Insert(r.db, "site", site, site.ID)
}

func (r *dbSiteRepo) List(offset, limit int) ([]*domain.Site, error) {
  site := &domain.Site{}
  sites := make([]*domain.Site, 0)
  sql, targets := structsql.ListSql("site", site, offset, limit)
  err := structsql.QueryAndCollect(r.db, sql, targets, func() {
    siteCopy := domain.Site(*site)
    sites = append(sites, &siteCopy)
  })
  return sites, err
}

func (r *dbSiteRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "site", ID)
}

func (r *dbSiteRepo) UpdateByID(ID string, oldSite, newSite *domain.Site, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "site", diffs.Modified(), ID)
}

func (r *dbSiteRepo) Export(dbr *Repos, w io.Writer) error {
  return dbr.exportTableFromStruct(w, "site", &domain.Site{})
}
