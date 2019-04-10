package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type stageQuery struct{
  h *handler
}

func (sc *stageQuery) EntityTypeName() string {
  return "stage"
}

func (sc *stageQuery) NewEntity() interface{} {
  return &domain.Stage{}
}

func (sc *stageQuery) SummaryQuery() string {
  return "select Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) stage(w http.ResponseWriter, r *http.Request) {
  sq := &stageQuery{h}
  h.stdquery(w, r, sq)
}
