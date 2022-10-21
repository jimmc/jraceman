package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type exceptionQuery struct{
  h *handler
}

func (sc *exceptionQuery) EntityTypeName() string {
  return "exception"
}

func (sc *exceptionQuery) NewEntity() interface{} {
  return &domain.Exception{}
}

func (sc *exceptionQuery) SummaryQuery() string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) exception(w http.ResponseWriter, r *http.Request) {
  sq := &exceptionQuery{h}
  h.stdquery(w, r, sq)
}
