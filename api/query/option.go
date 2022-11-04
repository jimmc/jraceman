package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type optionQuery struct{
  h *handler
}

func (sc *optionQuery) EntityTypeName() string {
  return "option"
}

func (sc *optionQuery) NewEntity() interface{} {
  return &domain.Option{}
}

func (sc *optionQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) option(w http.ResponseWriter, r *http.Request) {
  sq := &optionQuery{h}
  h.stdquery(w, r, sq)
}
