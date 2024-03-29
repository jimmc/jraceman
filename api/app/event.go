package app

import (
  "context"
  "fmt"
  "net/http"
  "strings"

  mainapp "github.com/jimmc/jraceman/app"
  apihttp "github.com/jimmc/jraceman/api/http"

  "github.com/golang/glog"
)

func (h *handler) event(w http.ResponseWriter, r *http.Request) {
  // TODO - check authorization
  morePath := strings.TrimPrefix(r.URL.Path, h.apiPrefix("event"))
  glog.V(1).Infof("%s morePath='%s'", r.Method, morePath);
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
      if action!="" && action!="races" {
        http.Error(w, fmt.Sprintf("Invalid GET action %s", action), http.StatusBadRequest)
        return
      }
      h.eventRaces(w, eventId)
    case http.MethodPost:
      switch action {
        case "createraces":
          params, err := apihttp.GetRequestParameters(r)
          if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
          }
          h.createRaces(r.Context(), w, eventId, params)
          return
        default:
          http.Error(w, fmt.Sprintf("Invalid POST action %s", action), http.StatusBadRequest)
          return
      }
    default:
      http.Error(w, "Method must be GET or POST", http.StatusMethodNotAllowed)
  }
}

func (h *handler) eventRaces(w http.ResponseWriter, eventId string) {
  eventRacesRepo := h.config.DomainRepos.EventRaces()
  result, err := eventRacesRepo.EventRaceInfo(eventId)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  apihttp.MarshalAndReply(w, result)
}

func (h *handler) createRaces(ctx context.Context, w http.ResponseWriter, eventId string, params map[string]interface{}) {
  glog.V(2).Infof("createRaces params=%#v", params)
  laneCount := apihttp.GetJsonIntParameter(params, "laneCount", -1)
  dryRun := apihttp.GetJsonBoolParameter(params, "dryRun", true)
  allowDeleteLanes := apihttp.GetJsonBoolParameter(params, "allowDeleteLanes", false)
  r := h.config.DomainRepos
  result, err := mainapp.EventCreateRaces(ctx, r, eventId, laneCount, dryRun, allowDeleteLanes)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  apihttp.MarshalAndReply(w, result)
}
