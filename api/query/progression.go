package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type progressionQuery struct{
  h *handler
}

func (sc *progressionQuery) EntityTypeName() string {
  return "progression"
}

func (sc *progressionQuery) NewEntity() interface{} {
  return &domain.Progression{}
}

func (sc *progressionQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) progression(w http.ResponseWriter, r *http.Request) {
  sq := &progressionQuery{h}
  h.stdquery(w, r, sq)
}
