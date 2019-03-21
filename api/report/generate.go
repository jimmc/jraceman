package report

import (
  "fmt"
  "log"
  "net/http"
  "strings"

  apihttp "github.com/jimmc/jracemango/api/http"
  "github.com/jimmc/jracemango/dbrepo"
  reportmain "github.com/jimmc/jracemango/report"
)

func (h *handler) generate(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case http.MethodGet:
      rName := r.URL.Query().Get("name")
      rData := r.URL.Query().Get("data")
      rOrderBy := r.URL.Query().Get("orderby")
      h.generateReportForHTTP(w, r, rName, rData, rOrderBy)
    case http.MethodPost:
      // When using a POST, we expect the name and data values as JSON parameters in the body.
      jsonBody, err := apihttp.GetRequestParameters(r)
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }
      rName := apihttp.GetJsonParameter(jsonBody, "name")
      rData := apihttp.GetJsonParameter(jsonBody, "data")
      rOrderBy := apihttp.GetJsonParameter(jsonBody, "orderby")
      h.generateReportForHTTP(w, r, rName, rData, rOrderBy)
    default:
      http.Error(w, "Method must be GET or POST", http.StatusMethodNotAllowed)
  }
}

func (h *handler) generateReportForHTTP(w http.ResponseWriter, r *http.Request, name, data string, orderby string) {
  name = strings.TrimSpace(name)
  if name == "" {
    http.Error(w, "No report name specified", http.StatusBadRequest)
    return
  }
  options := &reportmain.ReportOptions{
    OrderByKey: orderby,
  }
  dbrepos, ok := h.config.DomainRepos.(*dbrepo.Repos)
  if !ok {
    http.Error(w, "Bad database repo", http.StatusInternalServerError)
    return
  }
  db := dbrepos.DB()
  log.Printf("Generating report: %v", name)
  result, err := reportmain.GenerateResults(db, h.config.ReportRoots, name, data, options)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error generating report: %v", err), http.StatusBadRequest)
    return
  }

  apihttp.MarshalAndReply(w, result)
}
