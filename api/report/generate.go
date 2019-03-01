package report

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strings"

  apihttp "github.com/jimmc/jracemango/api/http"
  "github.com/jimmc/jracemango/dbrepo"
  // "github.com/jimmc/jracemango/dbrepo/strsql"
  reportmain "github.com/jimmc/jracemango/report"
)

func (h *handler) generate(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case http.MethodGet:
      rName := r.URL.Query().Get("name")
      rData := r.URL.Query().Get("data")
      h.generateReportForHTTP(w, r, rName, rData)
    case http.MethodPost:
      // When using a POST, we expect the name and data values as JSON parameters in the body.
      jsonBody, err := h.getRequestParameters(r)
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }
      rName := h.getJsonParameter(jsonBody, "name")
      rData := h.getJsonParameter(jsonBody, "data")
      h.generateReportForHTTP(w, r, rName, rData)
    default:
      http.Error(w, "Method must be GET or POST", http.StatusMethodNotAllowed)
  }
}

func (h *handler) getRequestParameters(r *http.Request) (map[string]interface{}, error) {
  contentType := r.Header.Get("content-type")
  log.Printf("content-type: %v\n", contentType)
  if contentType != "application/json" {
    return nil, fmt.Errorf("POST requires content-type=application/json")
  }
  decoder := json.NewDecoder(r.Body)
  jsonBody := make(map[string]interface{}, 0)
  if err := decoder.Decode(&jsonBody); err != nil {
    return nil, fmt.Errorf("Error decoding JSON body: %v", err)
  }
  return jsonBody, nil
}

func (h *handler) getJsonParameter(jsonBody map[string]interface{}, name string) string {
  val, ok := jsonBody[name]
  if !ok {
    return ""
  }
  s, ok := val.(string)
  if !ok {
    return ""
  }
  return s
}

func (h *handler) generateReportForHTTP(w http.ResponseWriter, r *http.Request, name, data string) {
  name = strings.TrimSpace(name)
  if name == "" {
    http.Error(w, "No report name specified", http.StatusBadRequest)
    return
  }
  dbrepos, ok := h.config.DomainRepos.(*dbrepo.Repos)
  if !ok {
    http.Error(w, "Bad database repo", http.StatusInternalServerError)
    return
  }
  db := dbrepos.DB()
  log.Printf("Generating report: %v", name)
  refdirs := []string{h.config.ReportRoot}
  result, err := reportmain.GenerateResults(db, refdirs, name, data)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error generating report: %v", err), http.StatusBadRequest)
    return
  }

  apihttp.MarshalAndReply(w, result)
}
