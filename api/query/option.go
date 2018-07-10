package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
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

func (h *handler) option(w http.ResponseWriter, r *http.Request) {
  sq := &optionQuery{h}
  h.stdquery(w, r, sq)
}
