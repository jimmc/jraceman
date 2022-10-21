package report

import (
  "net/http"

  apihttp "github.com/jimmc/jraceman/api/http"
  "github.com/jimmc/jraceman/dbrepo"
  "github.com/jimmc/jraceman/report"
)

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
  dbrepos, ok := h.config.DomainRepos.(*dbrepo.Repos)
  if !ok {
    http.Error(w, "Bad database repo", http.StatusInternalServerError)
    return
  }
  repAttrs, err := report.ClientVisibleReports(dbrepos, h.config.ReportRoots)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  apihttp.MarshalAndReply(w, repAttrs)
}
