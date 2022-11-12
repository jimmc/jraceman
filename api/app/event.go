package app

import (
  "fmt"
  "net/http"
  "strings"

  mainapp "github.com/jimmc/jraceman/app"
  apihttp "github.com/jimmc/jraceman/api/http"
  "github.com/jimmc/jraceman/dbrepo"

  "github.com/golang/glog"
)

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
  dbr := h.config.DomainRepos.(*dbrepo.Repos)
  result, err := mainapp.EventRaceInfo(dbr, eventId)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  apihttp.MarshalAndReply(w, result)
}
