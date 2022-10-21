package query

import (
  "net/http"

  "github.com/jimmc/jraceman/domain"
)

type eventQuery struct{
  h *handler
}

func (sc *eventQuery) EntityTypeName() string {
  return "event"
}

func (sc *eventQuery) NewEntity() interface{} {
  return &domain.Event{}
}

func (sc *eventQuery) SummaryQuery() string {
  return "select ID, Name || '[' || ID || ']' as summary from " + sc.EntityTypeName()
}

func (h *handler) event(w http.ResponseWriter, r *http.Request) {
  sq := &eventQuery{h}
  h.stdquery(w, r, sq)
}
