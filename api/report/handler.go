package report

import (
  "fmt"
  "net/http"
  "strings"

  "github.com/jimmc/jraceman/domain"
)

type handler struct {
  config *Config
}

// Config provides configuration of the http handler for our calls.
type Config struct {
  Prefix string
  DomainRepos domain.Repos
  ReportRoots []string
}

// NewHandler creates the http handler for our calls.
func NewHandler(c *Config) http.Handler {
  h := handler{config: c}
  mux := http.NewServeMux()
  // Permissions are checked in h.list and h.generate.
  mux.HandleFunc(h.apiPrefix("generate"), h.generate)
  mux.HandleFunc(h.config.Prefix, h.root)
  return mux
}

// root calls either h.list or h.generateByName, which both check permissions.
func (h *handler) root(w http.ResponseWriter, r *http.Request) {
  p := strings.TrimPrefix(r.URL.Path, h.config.Prefix)
  if p=="" {
    // If no path components after our root, return a list of reports.
    h.list(w, r)
    return
  }
  ss := strings.Split(p, "/")
  // Since we know p!="", ss must have at least one element.
  reportName := ss[0]
  action := "attributes"        // Preset to the default action.
  if len(ss) > 1 && ss[1] != "" {
    action = ss[1]
  }

  switch action {
  // case "attributes": h.attributesByName(w, r, reportName)
  case "generate": h.generateByName(w, r, reportName)
  default: http.Error(w, fmt.Sprintf("Try %s\n", h.config.Prefix), http.StatusForbidden)
  }
}

func (h *handler) apiPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
