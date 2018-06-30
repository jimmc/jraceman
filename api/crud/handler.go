package crud

import (
  "fmt"
  "net/http"

  "github.com/jimmc/jracemango/domain"
)

type handler struct {
  config *Config
}

// Config provides configuration of the http handler for CRUD calls.
type Config struct {
  Prefix string
  DomainRepos domain.Repos
}

// NewHandler creates the http handler for CRUD calls.
func NewHandler(c *Config) http.Handler {
  h := handler{config: c}
  mux := http.NewServeMux()
  mux.HandleFunc(h.crudPrefix("area"), h.area)
  mux.HandleFunc(h.crudPrefix("competition"), h.competition)
  mux.HandleFunc(h.crudPrefix("complan"), h.complan)
  mux.HandleFunc(h.crudPrefix("complanrule"), h.complanrule)
  mux.HandleFunc(h.crudPrefix("complanstage"), h.complanstage)
  mux.HandleFunc(h.crudPrefix("exception"), h.exception)
  mux.HandleFunc(h.crudPrefix("gender"), h.gender)
  mux.HandleFunc(h.crudPrefix("laneorder"), h.laneorder)
  mux.HandleFunc(h.crudPrefix("level"), h.level)
  mux.HandleFunc(h.crudPrefix("progression"), h.progression)
  mux.HandleFunc(h.crudPrefix("scoringrule"), h.scoringrule)
  mux.HandleFunc(h.crudPrefix("scoringsystem"), h.scoringsystem)
  mux.HandleFunc(h.crudPrefix("simplan"), h.simplan)
  mux.HandleFunc(h.crudPrefix("simplanrule"), h.simplanrule)
  mux.HandleFunc(h.crudPrefix("simplanstage"), h.simplanstage)
  mux.HandleFunc(h.crudPrefix("site"), h.site)
  mux.HandleFunc(h.crudPrefix("stage"), h.stage)
  return mux
}

func (h *handler) crudPrefix(s string) string {
  return fmt.Sprintf("%s%s/", h.config.Prefix, s)
}
