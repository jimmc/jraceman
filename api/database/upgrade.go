package database

import (
  "net/http"
  "strings"

  apihttp "github.com/jimmc/jracemango/api/http"
  "github.com/jimmc/jracemango/dbrepo"
)

func (h *handler) upgrade(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case http.MethodGet:
      h.upgradeList(w, r);
    case http.MethodPost:
      // When using a POST, we expect the section as the next path component.
      sectionName := strings.TrimPrefix(r.URL.Path, h.apiPrefix("upgrade"))
      if sectionName == "" {
        http.Error(w, "Section name must be specified on a POST", http.StatusBadRequest)
        return
      }
      h.upgradeSection(w, r, sectionName)
    default:
      http.Error(w, "Method must be GET or POST", http.StatusMethodNotAllowed)
  }
}

func (h *handler) upgradeList(w http.ResponseWriter, r *http.Request) {
  dbrepos, ok := h.config.DomainRepos.(*dbrepo.Repos)
  if !ok {
    http.Error(w, "Bad database repo", http.StatusInternalServerError)
    return
  }
  sectionNames := dbrepos.SectionNames()

  apihttp.MarshalAndReply(w, sectionNames)
}

func (h *handler) upgradeSection(w http.ResponseWriter, r *http.Request, sectionName string) {
  /*
  dryrunStr := r.URL.Query().Get("dryrun")
  dryrun := (dryRunStr == "true")
  dbrepos, ok := h.config.DomainRepos.(*dbrepo.Repos)
  if !ok {
    http.Error(w, "Bad database repo", http.StatusInternalServerError)
    return
  }
  */

  // TODO - call a (new) method in dbrepos to upgrade the specified section.
  /*
  db := dbrepos.DB()
  result, err := strsql.QueryStarAndCollect(db, sqlStr)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error executing sql: %v", err), http.StatusBadRequest)
    return
  }
  */
  result := "NYI"

  apihttp.MarshalAndReply(w, result)
}
