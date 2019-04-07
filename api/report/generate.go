package report

import (
  "fmt"
  "net/http"
  "strings"

  apihttp "github.com/jimmc/jracemango/api/http"
  "github.com/jimmc/jracemango/dbrepo"
  reportmain "github.com/jimmc/jracemango/report"

  "github.com/golang/glog"
)

func (h *handler) generate(w http.ResponseWriter, r *http.Request) {
  h.generateByName(w, r, "")
}

func (h *handler) generateByName(w http.ResponseWriter, r *http.Request, rName string) {
  switch r.Method {
    case http.MethodGet:
      if rName == "" {
        rName = r.URL.Query().Get("name")
      }
      rData := r.URL.Query().Get("data")
      rOrderBy := r.URL.Query().Get("orderby")
      // WHERE parameters require nesting that is too complex for 
      // reasonable specification when using GET, so we don't allow that.
      options := reportmain.ReportOptions{
        Name: rName,
        Data: rData,
        OrderBy: rOrderBy,
      }
      h.generateReportForHTTP(w, r, rName, &options)
    case http.MethodPost:
      // When using a POST, we expect the JSON data to confirm to the ReportOptions struct.
      options := reportmain.ReportOptions{}
      // When using a POST, we expect the name and data values as JSON parameters in the body.
      if err := apihttp.GetRequestParametersInto(r, &options); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }
      if options.Name == "" {
        options.Name = rName
      }
      h.generateReportForHTTP(w, r, rName, &options)
    default:
      http.Error(w, "Method must be GET or POST", http.StatusMethodNotAllowed)
  }
}

func (h *handler) generateReportForHTTP(w http.ResponseWriter, r *http.Request, name string, options *reportmain.ReportOptions) {
  name = strings.TrimSpace(name)
  if name == "" {
    name = options.Name
  }
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
  glog.Infof("Generating report: %v", name)
  result, err := reportmain.GenerateResults(db, h.config.ReportRoots, name, options)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error generating report: %v", err), http.StatusBadRequest)
    return
  }

  apihttp.MarshalAndReply(w, result)
}
