package crud

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strings"
)

func (h *handler) site(w http.ResponseWriter, r *http.Request) {
  // TODO - check authorization
  path := strings.TrimPrefix(r.URL.Path, h.crudPrefix("site"))
  log.Printf("site path: %s", path);
  switch r.Method {
    case http.MethodGet:
      if path == "" {
        h.siteList(w, r)
      } else {
        h.siteGet(w, r, path)
      }
    case http.MethodPost:
      if path != "" {
        http.Error(w, "Path may not be specified on a POST", http.StatusBadRequest)
      } else {
        h.siteCreate(w, r)
      }
    case http.MethodPut:
      if path == "" {
        http.Error(w, "Path must be specified on a PUT", http.StatusBadRequest)
      } else {
        h.siteUpdate(w, r, path)
      }
    case http.MethodDelete:
      if path == "" {
        http.Error(w, "Path must be specified on a DELETE", http.StatusBadRequest)
      } else {
        h.siteDelete(w, r, path)
      }
    default:
      http.Error(w, "Method must be GET, POST, PUT, or DELETE", http.StatusMethodNotAllowed)
  }
}

func (h *handler) siteCreate(w http.ResponseWriter, r *http.Request) {
  http.Error(w, "Create site is not implemented", http.StatusNotImplemented)
}

func (h *handler) siteList(w http.ResponseWriter, r *http.Request) {
  http.Error(w, "Listing all sites is not implemented", http.StatusNotImplemented)
}

func (h *handler) siteGet(w http.ResponseWriter, r *http.Request, path string) {
  result, err := h.config.DomainRepos.Site.FindById(path)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  b, err := json.MarshalIndent(result, "", "  ")
  if err != nil {
    http.Error(w, fmt.Sprintf("Failed to marshall json results: %v", err), http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write(b)
}

func (h *handler) siteUpdate(w http.ResponseWriter, r *http.Request, path string) {
  http.Error(w, "Update site is not implemented", http.StatusNotImplemented)
}

func (h *handler) siteDelete(w http.ResponseWriter, r *http.Request, path string) {
  http.Error(w, "Delete site is not implemented", http.StatusNotImplemented)
}
