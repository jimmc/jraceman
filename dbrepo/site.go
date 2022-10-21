package dbrepo

import (
  "database/sql"
  "io"

  "github.com/jimmc/jraceman/dbrepo/ixport"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"
  "github.com/jimmc/jraceman/domain"
)

type DBSiteRepo struct {
  db *sql.DB
}

func (r *DBSiteRepo) New() interface{} {
  return domain.Site{}
}

func (r *DBSiteRepo) CreateTable() error {
  return structsql.CreateTable(r.db, "site", domain.Site{})
}

func (r *DBSiteRepo) UpgradeTable(dryrun bool) (bool, string, error) {
  return structsql.UpgradeTable(r.db, "site", domain.Site{}, dryrun)
}

func (r *DBSiteRepo) FindByID(ID string) (*domain.Site, error) {
  site := &domain.Site{}
  sql, targets := structsql.FindByIDSql("site", site)
  if err := r.db.QueryRow(sql, ID).Scan(targets...); err != nil {
    return nil, err
  }
  return site, nil
}

func (r *DBSiteRepo) Save(site *domain.Site) (string, error) {
  if (site.ID == "") {
    site.ID = structsql.UniqueID(r.db, "site", "SI1")
  }
  return site.ID, structsql.Insert(r.db, "site", site, site.ID)
}

func (r *DBSiteRepo) List(offset, limit int) ([]*domain.Site, error) {
  site := &domain.Site{}
  sites := make([]*domain.Site, 0)
  sql, targets := structsql.ListSql("site", site, offset, limit)
  err := strsql.QueryAndCollect(r.db, sql, targets, func() {
    siteCopy := domain.Site(*site)
    sites = append(sites, &siteCopy)
  })
  return sites, err
}

func (r *DBSiteRepo) DeleteByID(ID string) error {
  return structsql.DeleteByID(r.db, "site", ID)
}

func (r *DBSiteRepo) UpdateByID(ID string, oldSite, newSite *domain.Site, diffs domain.Diffs) error {
  return structsql.UpdateByID(r.db, "site", diffs.Modified(), ID)
}

func (r *DBSiteRepo) Export(e *ixport.Exporter, w io.Writer) error {
  return e.ExportTableFromStruct(w, "site", &domain.Site{})
}
