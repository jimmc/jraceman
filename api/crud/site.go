package crud

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strings"

  "github.com/jimmc/jracemango/domain"
)

func (h *handler) site(w http.ResponseWriter, r *http.Request) {
  // TODO - check authorization
  entityID := strings.TrimPrefix(r.URL.Path, h.crudPrefix("site"))
  log.Printf("site entityID: %s", entityID);
  switch r.Method {
    case http.MethodGet:
      if entityID == "" {
        h.siteList(w, r)
      } else {
        h.siteGet(w, r, entityID)
      }
    case http.MethodPost:
      if entityID != "" {
        http.Error(w, "Entity ID may not be specified on a POST", http.StatusBadRequest)
      } else {
        h.siteCreate(w, r)
      }
    case http.MethodPut:
      if entityID == "" {
        http.Error(w, "Entity ID must be specified on a PUT", http.StatusBadRequest)
      } else {
        h.siteUpdate(w, r, entityID)
      }
    case http.MethodDelete:
      if entityID == "" {
        http.Error(w, "Entity ID must be specified on a DELETE", http.StatusBadRequest)
      } else {
        h.siteDelete(w, r, entityID)
      }
    default:
      http.Error(w, "Method must be GET, POST, PUT, or DELETE", http.StatusMethodNotAllowed)
  }
}

func (h *handler) siteCreate(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var site domain.Site
  err := decoder.Decode(&site)
  if err != nil {
    msg := fmt.Sprintf("Error decoding JSON: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  defer r.Body.Close()
  err = h.config.DomainRepos.Site().Save(&site)
  if err != nil {
    msg := fmt.Sprintf("Error saving data: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  res := `{"status": "ok"}`
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(res))
}

func (h *handler) siteList(w http.ResponseWriter, r *http.Request) {
  http.Error(w, "Listing all sites is not implemented", http.StatusNotImplemented)
}

func (h *handler) siteGet(w http.ResponseWriter, r *http.Request, entityID string) {
  result, err := h.config.DomainRepos.Site().FindById(entityID)
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

func (h *handler) siteUpdate(w http.ResponseWriter, r *http.Request, entityID string) {
  http.Error(w, "Update site is not implemented", http.StatusNotImplemented)
}

func (h *handler) siteDelete(w http.ResponseWriter, r *http.Request, entityID string) {
  http.Error(w, "Delete site is not implemented", http.StatusNotImplemented)
}
