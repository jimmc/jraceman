package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type registrationQuery struct{
  domain.RegistrationMeta
  h *handler
}

func (sc *registrationQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) registration(w http.ResponseWriter, r *http.Request) {
  sq := &registrationQuery{domain.RegistrationMeta{}, h}
  h.stdquery(w, r, sq)
}
