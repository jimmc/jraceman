package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type siteQuery struct{
  h *handler
}

func (sc *siteQuery) EntityTypeName() string {
  return "site"
}

func (sc *siteQuery) NewEntity() interface{} {
  return &domain.Site{}
}

func (sc *siteQuery) SummaryQuery() string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) site(w http.ResponseWriter, r *http.Request) {
  sq := &siteQuery{h}
  h.stdquery(w, r, sq)
}
