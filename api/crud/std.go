package crud

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strings"
)

type std interface {
  EntityTypeName() string   // such as "site"
  NewEntity() interface{}       // such as a new Site
  Save(entity interface{}) error        // function must cast entity to its type
  FindById(ID string) (interface{}, error)       // returns same type as NewEntity
}

func (h *handler) stdcrud(w http.ResponseWriter, r *http.Request, st std) {
  // TODO - check authorization
  entityType := st.EntityTypeName()
  entityID := strings.TrimPrefix(r.URL.Path, h.crudPrefix(entityType))
  log.Printf("%s entityID: %s", entityType, entityID);
  switch r.Method {
    case http.MethodGet:
      if entityID == "" {
        h.stdList(w, r, st)
      } else {
        h.stdGet(w, r, st, entityID)
      }
    case http.MethodPost:
      if entityID != "" {
        http.Error(w, "Entity ID may not be specified on a POST", http.StatusBadRequest)
      } else {
        h.stdCreate(w, r, st)
      }
    case http.MethodPut:
      if entityID == "" {
        http.Error(w, "Entity ID must be specified on a PUT", http.StatusBadRequest)
      } else {
        h.stdUpdate(w, r, st ,entityID)
      }
    case http.MethodDelete:
      if entityID == "" {
        http.Error(w, "Entity ID must be specified on a DELETE", http.StatusBadRequest)
      } else {
        h.stdDelete(w, r, st, entityID)
      }
    default:
      http.Error(w, "Method must be GET, POST, PUT, or DELETE", http.StatusMethodNotAllowed)
  }
}

func (h *handler) stdCreate(w http.ResponseWriter, r *http.Request, st std) {
  decoder := json.NewDecoder(r.Body)
  entity := st.NewEntity()
  err := decoder.Decode(entity)
  if err != nil {
    msg := fmt.Sprintf("Error decoding JSON: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  defer r.Body.Close()
  err = st.Save(entity)
  if err != nil {
    msg := fmt.Sprintf("Error saving data: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  res := `{"status": "ok"}`
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(res))
}

func (h *handler) stdList(w http.ResponseWriter, r *http.Request, st std) {
  msg := fmt.Sprintf("List of %s is not implemented", st.EntityTypeName())
  http.Error(w, msg, http.StatusNotImplemented)
}

func (h *handler) stdGet(w http.ResponseWriter, r *http.Request, st std, entityID string) {
  result, err := st.FindById(entityID)
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

func (h *handler) stdUpdate(w http.ResponseWriter, r *http.Request, st std, entityID string) {
  msg := fmt.Sprintf("Update %s is not implemented", st.EntityTypeName())
  http.Error(w, msg, http.StatusNotImplemented)
}

func (h *handler) stdDelete(w http.ResponseWriter, r *http.Request, st std, entityID string) {
  msg := fmt.Sprintf("Delete %s is not implemented", st.EntityTypeName())
  http.Error(w, msg, http.StatusNotImplemented)
}
