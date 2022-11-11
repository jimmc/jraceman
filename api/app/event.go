package app

import (
  "fmt"
  "net/http"
  "strings"

  apihttp "github.com/jimmc/jraceman/api/http"
  "github.com/jimmc/jraceman/dbrepo"

  "github.com/golang/glog"
)

type EventInfo struct {
  EntryCount int
  GroupCount int
  GroupSize int
  Summary string
}

func (h *handler) event(w http.ResponseWriter, r *http.Request) {
  // TODO - check authorization
  morePath := strings.TrimPrefix(r.URL.Path, h.apiPrefix("event"))
  glog.Infof("%s morePath='%s'", r.Method, morePath);
  morePathParts := strings.Split(morePath, "/")
  eventId := ""
  if len(morePathParts)>0 {
    eventId = morePathParts[0]
  }
  if eventId == "" {
    http.Error(w, "Event ID must be specified", http.StatusBadRequest)
    return
  }
  action := ""
  if len(morePathParts)>1 {
    action = morePathParts[1]
  }
  switch r.Method {
    case http.MethodGet:
      if action!="" && action!="info" {
        http.Error(w, fmt.Sprintf("Invalid GET action %s", action), http.StatusBadRequest)
        return
      }
      h.eventInfo(w, eventId)
    case http.MethodPost:
      http.Error(w, "POST is not yet implemented", http.StatusBadRequest)
    default:
      http.Error(w, "Method must be GET or POST", http.StatusMethodNotAllowed)
  }
}

func (h *handler) eventInfo(w http.ResponseWriter, eventId string) {
  db := h.config.DomainRepos.(*dbrepo.Repos).DB()
  entryCountQuery :=
        `(SELECT count(1) as entryCount
        FROM entry JOIN event on entry.eventid = event.id
        WHERE event.id=? AND NOT entry.scratched)`
  groupCountQuery :=
        `(SELECT count(distinct groupname) as groupCount
        FROM entry JOIN event on entry.eventid = event.id
        WHERE event.id=? AND NOT entry.scratched)`
  query := "SELECT "+entryCountQuery+" as EntryCount,"+
        groupCountQuery+` as GroupCount,
        COALESCE(competition.groupsize,0) as GroupSize,
        event.Name || ' [' || event.ID || ']' as Summary
        FROM event
        LEFT JOIN competition on event.competitionid = competition.id
        WHERE event.id=?`
  whereVals := make([]interface{}, 3)
  whereVals[0] = eventId
  whereVals[1] = eventId
  whereVals[2] = eventId
  result := &EventInfo{}
  err := db.QueryRow(query, whereVals...).Scan(
    &result.EntryCount, &result.GroupCount, &result.GroupSize, &result.Summary)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error collecting event info: %v", err), http.StatusBadRequest)
    return
  }
  apihttp.MarshalAndReply(w, result)
}
