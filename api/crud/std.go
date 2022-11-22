package crud

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strconv"
  "strings"

  "github.com/jimmc/jraceman/domain"
  apihttp "github.com/jimmc/jraceman/api/http"

  "github.com/jimmc/auth/auth"
  "github.com/jimmc/auth/permissions"

  "github.com/golang/glog"
)

// Std defines the type-specific methods that are needed by the generic
// CRUD code.
type std interface {
  EntityTypeName() string   // such as "site"
  EntityGroupName() string      // The table group for auth puposes, such as "venue".
  NewEntity() interface{}       // such as a new Site
  List(offset, limit int) ([]interface{}, error)
  Save(entity interface{}) (string, error)        // function must cast entity to its type
  FindByID(ID string) (interface{}, error)       // returns same type as NewEntity
  DeleteByID(ID string) error
  UpdateByID(ID string, newEntity, oldEntity interface{}, diffs domain.Diffs) error
}

// StdCrud is an http handler that takes care of most of the work for
// CRUD requests on our domain types. The API CRUD handler for a specific
// type turns around and calls this handler with a type-specific std that
// defines the type-specific methods needed by this handler.
func (h *handler) stdcrud(w http.ResponseWriter, r *http.Request, st std) {
  // We require read or write permission for the table group, depending on http mthod.
  permPrefix := "edit_"         // Assume we will need write privileges.
  if r.Method == http.MethodGet {
    permPrefix = "view_"       // GET requires read privilege.
  }
  permissionName := permPrefix + st.EntityGroupName()
  permission := permissions.Permission(permissionName)
  if !auth.CurrentUserHasPermission(r, permission) {
    currentUser := auth.CurrentUser(r)
    username := "(no current user)"
    if currentUser != nil {
      username = currentUser.Id()
    }
    glog.Infof("Not authorized: user %q does not have permission %q", username, permissionName)
    http.Error(w, "Not authorized", http.StatusUnauthorized)
    return
  }
  entityType := st.EntityTypeName()
  entityID := strings.TrimPrefix(r.URL.Path, h.crudPrefix(entityType))
  glog.Infof("%s %s ID='%s'", r.Method, entityType, entityID);
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
      glog.V(3).Infof("Method is Patch")
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
    msg := fmt.Sprintf("Error decoding JSON for create: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  id, err := st.Save(entity)
  if err != nil {
    msg := fmt.Sprintf("Error saving data: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  result := struct {
    Status string
    ID string
  } {
    Status: "ok",
    ID: id,
  }
  apihttp.MarshalAndReply(w, result)
}

func (h *handler) stdList(w http.ResponseWriter, r *http.Request, st std) {
  limit, err := intQueryParam(r, "limit", 0)
  if err != nil || limit < 0 {
    msg := fmt.Sprintf("Bad format for limit")
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  offset, err := intQueryParam(r, "offset", 0)
  if err != nil || offset < 0 {
    msg := fmt.Sprintf("Bad format for offset")
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  if offset > 0 && limit == 0 {
    msg := fmt.Sprintf("Must specify limit when specifying offset")
    http.Error(w, msg, http.StatusBadRequest)
    return
  }

  result, err := st.List(offset, limit)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  // TODO - add offset and limit info to the result?
  apihttp.MarshalAndReply(w, result)
}

func (h *handler) stdGet(w http.ResponseWriter, r *http.Request, st std, entityID string) {
  result, err := st.FindByID(entityID)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  apihttp.MarshalAndReply(w, result)
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
    msg := fmt.Sprintf("Error decoding JSON for update: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }

  diffs, equal := deepDiff(oldEntity, newEntity)
  if equal {
    msg := fmt.Sprintf("No change specified")   // Maybe this should not be an error?
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  glog.V(1).Infof("update entity diffs: %v", diffs.Modified())

  if err := st.UpdateByID(entityID, oldEntity, newEntity, diffs); err != nil {
    msg := fmt.Sprintf("Error updating data: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  apihttp.OkResponse(w)
}

func (h *handler) stdPatch(w http.ResponseWriter, r *http.Request, st std, entityID string) {
  glog.V(3).Infof("Begin stdPatch")
  oldEntity, err := st.FindByID(entityID)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  decoder := json.NewDecoder(r.Body)
  var patch interface{}
  if err := decoder.Decode(&patch); err != nil {
    msg := fmt.Sprintf("Error decoding JSON for patch: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  glog.V(2).Infof("decoded patch: %+v", patch)

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
  glog.V(1).Infof("patch entity diffs: %v", diffs.Modified())

  if err := st.UpdateByID(entityID, oldEntity, newEntity, diffs); err != nil {
    msg := fmt.Sprintf("Error updating data: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  apihttp.OkResponse(w)
}

func (h *handler) stdDelete(w http.ResponseWriter, r *http.Request, st std, entityID string) {
  if err := st.DeleteByID(entityID); err != nil {
    msg := fmt.Sprintf("Error deleting data: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  apihttp.OkResponse(w)
}

func intQueryParam(r *http.Request, paramName string, dflt int) (int, error) {
  s := r.URL.Query().Get(paramName)
  if s == "" {
    return dflt, nil
  }
  n, err := strconv.Atoi(s)
  if err != nil {
    return dflt, err
  }
  return n, nil
}
