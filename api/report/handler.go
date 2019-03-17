package report

import (
  "fmt"
  "net/http"
  "strings"

  "github.com/jimmc/jracemango/domain"
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
  mux.HandleFunc(h.apiPrefix("generate"), h.generate)
  mux.HandleFunc(h.config.Prefix, h.root)
  return mux
}

func (h *handler) root(w http.ResponseWriter, r *http.Request) {
  p := strings.TrimPrefix(r.URL.Path, h.config.Prefix)
  if p=="" {
    // If no path components after our root, return a list of reports.
    h.list(w, r)
    return
  }
  http.Error(w, "Try /api/report/\n", http.StatusForbidden)
}

func (h *handler) apiPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
