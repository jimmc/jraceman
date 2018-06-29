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

func (h *handler) stage(w http.ResponseWriter, r *http.Request) {
  sq := &stageQuery{h}
  h.stdquery(w, r, sq)
}
