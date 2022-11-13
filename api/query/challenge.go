package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type challengeQuery struct{
  domain.ChallengeMeta
  h *handler
}

func (sc *challengeQuery) SummaryQuery(format string) string {
  return "select ID, Name || '[' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) challenge(w http.ResponseWriter, r *http.Request) {
  sq := &challengeQuery{domain.ChallengeMeta{}, h}
  h.stdquery(w, r, sq)
}
