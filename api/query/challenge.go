package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
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

func (sc *challengeQuery) SummaryQuery(format string) string {
  return "select ID, Name || '[' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) challenge(w http.ResponseWriter, r *http.Request) {
  sq := &challengeQuery{h}
  h.stdquery(w, r, sq)
}
