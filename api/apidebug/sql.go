package apidebug

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "strings"

  "github.com/jimmc/jracemango/dbrepo"
  "github.com/jimmc/jracemango/dbrepo/strsql"
)

func (h *handler) sql(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case http.MethodGet:
      sqlStr := r.URL.Query().Get("q")
      h.executeSql(w, r, sqlStr)
    case http.MethodPost:
      sqlBytes, err := ioutil.ReadAll(r.Body)
      if err != nil {
        http.Error(w, "Bad POST body", http.StatusBadRequest)
        return
      }
      sqlStr := string(sqlBytes)
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
  dbrepos, ok := h.config.DomainRepos.(*dbrepo.Repos)
  if !ok {
    http.Error(w, "Bad database repo", http.StatusInternalServerError)
    return
  }
  db := dbrepos.DB()
  rows, err := strsql.QueryStarAndCollect(db, sqlStr)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error executing sql: %v", err), http.StatusBadRequest)
    return
  }

  // TODO - the rest of this function can probably be refactored into
  // a common helper function.
  b, err := json.MarshalIndent(rows, "", "  ")
  if err != nil {
    http.Error(w, fmt.Sprintf("Failed to marshall json results: %v", err), http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write(b)
}
