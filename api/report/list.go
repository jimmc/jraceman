package report

import (
  "net/http"

  apihttp "github.com/jimmc/jraceman/api/http"
  "github.com/jimmc/jraceman/dbrepo"
  "github.com/jimmc/jraceman/report"

  authlib "github.com/jimmc/auth/auth"
)

// list gets the list of permitted reports base on the current user's permissions.
func (h *handler) list(w http.ResponseWriter, r *http.Request) {
  dbrepos, ok := h.config.DomainRepos.(*dbrepo.Repos)
  if !ok {
    http.Error(w, "Bad database repo", http.StatusInternalServerError)
    return
  }
  currentPerms := authlib.CurrentPermissions(r)
  repAttrs, err := report.ClientPermittedReports(dbrepos, h.config.ReportRoots, currentPerms)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  apihttp.MarshalAndReply(w, repAttrs)
}
