package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type siteQuery struct{
  domain.SiteMeta
  h *handler
}

func (sc *siteQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) site(w http.ResponseWriter, r *http.Request) {
  sq := &siteQuery{domain.SiteMeta{}, h}
  h.stdquery(w, r, sq)
}
