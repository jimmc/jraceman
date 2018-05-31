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

func (h *handler) site(w http.ResponseWriter, r *http.Request) {
  sq := &siteQuery{h}
  h.stdquery(w, r, sq)
}
