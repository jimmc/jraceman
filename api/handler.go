package api

import (
  "fmt"
  "net/http"

  "github.com/jimmc/jraceman/api/crud"
  apidb "github.com/jimmc/jraceman/api/database"
  apidebug "github.com/jimmc/jraceman/api/debug"
  "github.com/jimmc/jraceman/api/query"
  apireport "github.com/jimmc/jraceman/api/report"
  "github.com/jimmc/jraceman/domain"
)

type handler struct {
  config *Config
}

// Config is used to configure the api handler.
// The Prefix field must be set to the part of the URL path that was
// used to route the request to this handler.
type Config struct {
  Prefix string
  DomainRepos domain.Repos
  CheckRoots []string
  ReportRoots []string
}

// NewHandler creates the http handler that is used to route api requests.
func NewHandler(c *Config) http.Handler {
  h := handler{config: c}
  mux := http.NewServeMux()

  crudPrefix := h.apiPrefix("crud")
  crudConfig := &crud.Config{
    Prefix: crudPrefix,
    DomainRepos: c.DomainRepos,
  }
  crudHandler := crud.NewHandler(crudConfig)
  mux.Handle(crudPrefix, crudHandler)

  dbPrefix := h.apiPrefix("database")
  dbConfig := &apidb.Config{
    Prefix: dbPrefix,
    DomainRepos: c.DomainRepos,
  }
  dbHandler := apidb.NewHandler(dbConfig)
  mux.Handle(dbPrefix, dbHandler)

  debugPrefix := h.apiPrefix("debug")
  debugConfig := &apidebug.Config{
    Prefix: debugPrefix,
    DomainRepos: c.DomainRepos,
  }
  debugHandler := apidebug.NewHandler(debugConfig)
  mux.Handle(debugPrefix, debugHandler)

  checkPrefix := h.apiPrefix("check")
  checkConfig := &apireport.Config{
    Prefix: checkPrefix,
    DomainRepos: c.DomainRepos,
    ReportRoots: c.CheckRoots,
  }
  checkHandler := apireport.NewHandler(checkConfig)
  mux.Handle(checkPrefix, checkHandler)

  reportPrefix := h.apiPrefix("report")
  reportConfig := &apireport.Config{
    Prefix: reportPrefix,
    DomainRepos: c.DomainRepos,
    ReportRoots: c.ReportRoots,
  }
  reportHandler := apireport.NewHandler(reportConfig)
  mux.Handle(reportPrefix, reportHandler)

  queryPrefix := h.apiPrefix("query")
  queryConfig := &query.Config{
    Prefix: queryPrefix,
    DomainRepos: c.DomainRepos,
  }
  queryHandler := query.NewHandler(queryConfig)
  mux.Handle(queryPrefix, queryHandler)

  mux.HandleFunc(h.config.Prefix, h.blank)
  return mux
}

func (h *handler) blank(w http.ResponseWriter, r *http.Request) {
  http.Error(w, "Try one of /api/crud, /api/database, /api/debug, /api/check, /api/report, /api/query",
      http.StatusForbidden)
}

// ApiPrefix composes our prefix with the next path component so that we can
// provide the right prefix to the handler that handles that next
// path component.
func (h *handler) apiPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
