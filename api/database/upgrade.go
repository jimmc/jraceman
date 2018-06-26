package database

import (
  "fmt"
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
      // When using a POST, we expect the table as the next path component.
      tableName := strings.TrimPrefix(r.URL.Path, h.apiPrefix("upgrade"))
      if tableName == "" {
        http.Error(w, "Table name must be specified on a POST", http.StatusBadRequest)
        return
      }
      h.upgradeTable(w, r, tableName)
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
  tableNames := dbrepos.TableNames()

  apihttp.MarshalAndReply(w, tableNames)
}

func (h *handler) upgradeTable(w http.ResponseWriter, r *http.Request, tableName string) {
  dryrunStr := r.URL.Query().Get("dryrun")
  dryrun := (dryrunStr == "true")
  dbrepos, ok := h.config.DomainRepos.(*dbrepo.Repos)
  if !ok {
    http.Error(w, "Bad database repo", http.StatusInternalServerError)
    return
  }

  nop, message, err := dbrepos.UpgradeTable(tableName, dryrun)
  if err != nil {
    http.Error(w, fmt.Sprintf("error upgrading table %s: %v", tableName, err), http.StatusBadRequest)
    return
  }
  type upgradeResult struct {
    Nop bool
    Message string
  }
  result := upgradeResult{
    Nop: nop,
    Message: message,
  }

  apihttp.MarshalAndReply(w, result)
}
