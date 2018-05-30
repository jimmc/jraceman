package debug

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strings"

  apihttp "github.com/jimmc/jracemango/api/http"
  "github.com/jimmc/jracemango/dbrepo"
  "github.com/jimmc/jracemango/dbrepo/strsql"
)

func (h *handler) sql(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case http.MethodGet:
      sqlStr := r.URL.Query().Get("q")
      h.executeSql(w, r, sqlStr)
    case http.MethodPost:
      // When using a POST, we expect the query as a parameter.
      sqlStr, err := h.getRequestParameter(r, "q")
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }
      h.executeSql(w, r, sqlStr)
    default:
      http.Error(w, "Method must be GET or POST", http.StatusMethodNotAllowed)
  }
}

func (h *handler) getRequestParameter(r *http.Request, name string) (string, error) {
    contentType := r.Header.Get("content-type")
    log.Printf("content-type: %v\n", contentType)
    if contentType == "application/json" {
      decoder := json.NewDecoder(r.Body)
      jsonBody := make(map[string]interface{}, 0)
      if err := decoder.Decode(&jsonBody); err != nil {
        return "", fmt.Errorf("Error decoding JSON body: %v", err)
      }
      val, ok := jsonBody[name]
      if ok {
        return val.(string), nil
      } else {
        return "", nil
      }
    } else {
      return r.FormValue("q"), nil
    }
}

func (h *handler) executeSql(w http.ResponseWriter, r *http.Request, sqlStr string) {
  sqlStr = strings.TrimSpace(sqlStr)
  if sqlStr == "" {
    http.Error(w, "No sql specified", http.StatusBadRequest)
    return
  }
  log.Printf("Executing sql: %v", sqlStr)
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
