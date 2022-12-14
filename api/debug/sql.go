package debug

import (
  "fmt"
  "net/http"
  "strings"

  apihttp "github.com/jimmc/jraceman/api/http"
  "github.com/jimmc/jraceman/dbrepo"
  "github.com/jimmc/jraceman/dbrepo/strsql"

  "github.com/golang/glog"
)

func (h *handler) sql(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case http.MethodGet:
      sqlStr := r.URL.Query().Get("q")
      h.executeSql(w, r, sqlStr)
    case http.MethodPost:
      // When using a POST, we expect the query value as a JSON parameter in the body.
      jsonBody, err := apihttp.GetRequestParameters(r)
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }
      sqlStr := apihttp.GetJsonStringParameter(jsonBody, "q")
      h.executeSql(w, r, sqlStr)
    default:
      http.Error(w, "Method must be GET or POST", http.StatusMethodNotAllowed)
  }
}

func (h *handler) executeSql(w http.ResponseWriter, r *http.Request, sqlStr string) {
  sqlStr = strings.TrimSpace(sqlStr)
  if sqlStr == "" {
    http.Error(w, "No sql specified", http.StatusBadRequest)
    return
  }
  glog.V(1).Infof("Executing sql: %v", sqlStr)
  dbrepos, ok := h.config.DomainRepos.(*dbrepo.Repos)
  if !ok {
    http.Error(w, "Bad database repo", http.StatusInternalServerError)
    return
  }
  db := dbrepos.DB()
  result, err := strsql.QueryStarAndCollect(db, sqlStr)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error executing sql: %v", err), http.StatusBadRequest)
    return
  }

  apihttp.MarshalAndReply(w, result)
}
