package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type optionQuery struct{
  domain.OptionMeta
  h *handler
}

func (sc *optionQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) option(w http.ResponseWriter, r *http.Request) {
  sq := &optionQuery{domain.OptionMeta{}, h}
  h.stdquery(w, r, sq)
}
