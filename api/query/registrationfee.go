package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type registrationfeeQuery struct{
  h *handler
}

func (sc *registrationfeeQuery) EntityTypeName() string {
  return "registrationfee"
}

func (sc *registrationfeeQuery) NewEntity() interface{} {
  return &domain.RegistrationFee{}
}

func (sc *registrationfeeQuery) SummaryQuery(format string) string {
  return "select ID, '[' || ID || '] ' as summary from " + sc.EntityTypeName()
}

func (h *handler) registrationfee(w http.ResponseWriter, r *http.Request) {
  sq := &registrationfeeQuery{h}
  h.stdquery(w, r, sq)
}
