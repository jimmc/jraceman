package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type exceptionQuery struct{
  domain.ExceptionMeta
  h *handler
}

func (sc *exceptionQuery) SummaryQuery(format string) string {
  return "select ID, Name || ' [' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) exception(w http.ResponseWriter, r *http.Request) {
  sq := &exceptionQuery{domain.ExceptionMeta{}, h}
  h.stdquery(w, r, sq)
}
