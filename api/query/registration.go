package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type registrationQuery struct{
  h *handler
}

func (sc *registrationQuery) EntityTypeName() string {
  return "registration"
}

func (sc *registrationQuery) NewEntity() interface{} {
  return &domain.Registration{}
}

func (h *handler) registration(w http.ResponseWriter, r *http.Request) {
  sq := &registrationQuery{h}
  h.stdquery(w, r, sq)
}
