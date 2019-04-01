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
      h.generateReportForHTTP(w, r, rName, rData, rOrderBy, nil)
    case http.MethodPost:
      // When using a POST, we expect the name and data values as JSON parameters in the body.
      jsonBody, err := apihttp.GetRequestParameters(r)
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }
      if rName == "" {
        rName = apihttp.GetJsonStringParameter(jsonBody, "name")
      }
      rData := apihttp.GetJsonStringParameter(jsonBody, "data")
      rOrderBy := apihttp.GetJsonStringParameter(jsonBody, "orderby")
      rWhere, ok := jsonBody["where"]
      if !ok {
        rWhere = nil
      }
      h.generateReportForHTTP(w, r, rName, rData, rOrderBy, rWhere)
    default:
      http.Error(w, "Method must be GET or POST", http.StatusMethodNotAllowed)
  }
}

func (h *handler) generateReportForHTTP(w http.ResponseWriter, r *http.Request, name, data string, orderby string, where interface{}) {
  name = strings.TrimSpace(name)
  if name == "" {
    http.Error(w, "No report name specified", http.StatusBadRequest)
    return
  }
  options, err := optionsFromParameters(orderby, where)
  if err != nil {
    http.Error(w, "Invalid options", http.StatusBadRequest)
    return
  }
  dbrepos, ok := h.config.DomainRepos.(*dbrepo.Repos)
  if !ok {
    http.Error(w, "Bad database repo", http.StatusInternalServerError)
    return
  }
  db := dbrepos.DB()
  glog.Infof("Generating report: %v", name)
  result, err := reportmain.GenerateResults(db, h.config.ReportRoots, name, data, options)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error generating report: %v", err), http.StatusBadRequest)
    return
  }

  apihttp.MarshalAndReply(w, result)
}

func OptionsFromParametersForTesting(orderby string, where interface{}) (*reportmain.ReportOptions, error) {
  return optionsFromParameters(orderby, where)
}

func optionsFromParameters(orderby string, where interface{}) (*reportmain.ReportOptions, error) {
  whereValues, err := whereMapFromParameters(where)
  if err != nil {
    return nil, err
  }
  options := &reportmain.ReportOptions{
    OrderByKey: orderby,
    WhereValues: whereValues,
  }
  return options, nil
}

func WhereMapFromParametersForTesting(where interface{}) (map[string]reportmain.WhereValue, error) {
  return whereMapFromParameters(where)
}

func whereMapFromParameters(where interface{}) (map[string]reportmain.WhereValue, error) {
  if where == nil {
    return map[string]reportmain.WhereValue{}, nil
  }
  whereMap, ok := where.(map[string]interface{})
  if !ok {
    return nil, fmt.Errorf("invalid 'where' options, must be map[string]interface, but is %T", where)
  }
  r := map[string]reportmain.WhereValue{}
  for k, v := range whereMap {
    vals, ok := v.(map[string]interface{})
    if !ok {
      return nil, fmt.Errorf("invalid value for where field %q, must be map[string]interface, but is %T",
          k, v)
    }
    opv, ok := vals["op"]
    if !ok {
      return nil, fmt.Errorf("missing op field for %q", k)
    }
    op, ok := opv.(string)
    if !ok {
      return nil, fmt.Errorf("op for field %q must be string, but is %T", k, op)
    }
    value := vals["value"]
    wv := reportmain.WhereValue{Op: op, Value: value}
    r[k] = wv
  }
  return r, nil
}
