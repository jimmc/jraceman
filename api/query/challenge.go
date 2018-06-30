package query

import (
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type challengeQuery struct{
  h *handler
}

func (sc *challengeQuery) EntityTypeName() string {
  return "challenge"
}

func (sc *challengeQuery) NewEntity() interface{} {
  return &domain.Challenge{}
}

func (h *handler) challenge(w http.ResponseWriter, r *http.Request) {
  sq := &challengeQuery{h}
  h.stdquery(w, r, sq)
}
