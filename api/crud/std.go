package crud

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "strings"

  "github.com/jimmc/jracemango/domain"
)

type std interface {
  EntityTypeName() string   // such as "site"
  NewEntity() interface{}       // such as a new Site
  Save(entity interface{}) error        // function must cast entity to its type
  FindByID(ID string) (interface{}, error)       // returns same type as NewEntity
  DeleteByID(ID string) error
  UpdateByID(ID string, newEntity, oldEntity interface{}, diffs domain.Diffs) error
}

func (h *handler) stdcrud(w http.ResponseWriter, r *http.Request, st std) {
  // TODO - check authorization
  entityType := st.EntityTypeName()
  entityID := strings.TrimPrefix(r.URL.Path, h.crudPrefix(entityType))
  log.Printf("%s %s '%s'", r.Method, entityType, entityID);
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
    case http.MethodPatch:
      if entityID == "" {
        http.Error(w, "Entity ID must be specified on a PATCH", http.StatusBadRequest)
      } else {
        h.stdPatch(w, r, st ,entityID)
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
  if err := decoder.Decode(entity); err != nil {
    msg := fmt.Sprintf("Error decoding JSON: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  if err := st.Save(entity); err != nil {
    msg := fmt.Sprintf("Error saving data: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  h.stdOkResponse(w)
}

func (h *handler) stdList(w http.ResponseWriter, r *http.Request, st std) {
  msg := fmt.Sprintf("List of %s is not implemented", st.EntityTypeName())
  http.Error(w, msg, http.StatusNotImplemented)
}

func (h *handler) stdGet(w http.ResponseWriter, r *http.Request, st std, entityID string) {
  result, err := st.FindByID(entityID)
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
  oldEntity, err := st.FindByID(entityID)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  decoder := json.NewDecoder(r.Body)
  newEntity := st.NewEntity()
  if err := decoder.Decode(newEntity); err != nil {
    msg := fmt.Sprintf("Error decoding JSON: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }

  diffs, equal := deepDiff(oldEntity, newEntity)
  if equal {
    msg := fmt.Sprintf("No change specified")   // Maybe this should not be an error?
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  log.Printf("entity diffs: %v", diffs.Modified())

  if err := st.UpdateByID(entityID, oldEntity, newEntity, diffs); err != nil {
    msg := fmt.Sprintf("Error updating data: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  h.stdOkResponse(w)
}

func (h *handler) stdPatch(w http.ResponseWriter, r *http.Request, st std, entityID string) {
  oldEntity, err := st.FindByID(entityID)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  decoder := json.NewDecoder(r.Body)
  var patch interface{}
  if err := decoder.Decode(&patch); err != nil {
    msg := fmt.Sprintf("Error decoding JSON patch: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }

  newEntity := st.NewEntity()
  diffs, equal, err := patchToDiffs(oldEntity, newEntity, patch)
  if err != nil {
    msg := fmt.Sprintf("Error merging JSON patch: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  if equal {
    msg := fmt.Sprintf("No change specified")   // Maybe this should not be an error?
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  log.Printf("entity diffs: %v", diffs.Modified())

  if err := st.UpdateByID(entityID, oldEntity, newEntity, diffs); err != nil {
    msg := fmt.Sprintf("Error updating data: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  h.stdOkResponse(w)
}

func (h *handler) stdDelete(w http.ResponseWriter, r *http.Request, st std, entityID string) {
  if err := st.DeleteByID(entityID); err != nil {
    msg := fmt.Sprintf("Error deleting data: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  h.stdOkResponse(w)
}

func (h *handler) stdOkResponse(w http.ResponseWriter) {
  res := `{"status": "ok"}`
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(res))
}
