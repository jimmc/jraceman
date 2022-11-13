package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type registrationfeeQuery struct{
  domain.RegistrationFeeMeta
  h *handler
}

func (sc *registrationfeeQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) registrationfee(w http.ResponseWriter, r *http.Request) {
  sq := &registrationfeeQuery{domain.RegistrationFeeMeta{}, h}
  h.stdquery(w, r, sq)
}
