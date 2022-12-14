package query

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strings"

  apihttp "github.com/jimmc/jraceman/api/http"
  "github.com/jimmc/jraceman/dbrepo"
  "github.com/jimmc/jraceman/dbrepo/strsql"
  "github.com/jimmc/jraceman/dbrepo/structsql"

  "github.com/jimmc/auth/auth"
  "github.com/jimmc/auth/permissions"

  "github.com/golang/glog"
)

// QueryParam defines one column comparison for an SQL query.
type queryParam struct {
  Name string
  Op string
  Value string
}

// GetColumnResults provides info about the available columns for a query.
type GetColumnResults struct {
  Columns []structsql.ColumnInfo
}

func (qp *queryParam) CleanName() string {
  // TODO - make sure the name is a valid column name
  return qp.Name
}

func (qp *queryParam) CleanOp() string {
  switch qp.Op {
  case "eq": return "="
  case "ne": return "!="
  case "gt": return ">"
  case "ge": return ">="
  case "lt": return "<"
  case "le": return "<="
  case "like": return "LIKE"
  }
  glog.Warning("ERROR: invalid op '"+qp.Op+"'")
  return "ERROR"        // This will result in invalid SQL.
}

func (qp *queryParam) CleanValue() interface{} {
  return qp.Value
}

// Std defines the type-specific methods that are needed by the generic
// query code.
type std interface {
  EntityTypeName() string   // such as "site"
  EntityGroupName() string      // The table group for auth puposes, such as "venue".
  NewEntity() interface{}       // such as a new Site
  SummaryQuery(format string) string    // An SQL query that is suitable for use as a summary for each row.
}

// Stdquery is an http handler that takes care of most of the work for
// query requests on our domain types. The API query handler for a specific
// type turns around and calls this handler with a type-specific std that
// defines the type-specific methods needed by this handler.
func (h *handler) stdquery(w http.ResponseWriter, r *http.Request, st std) {
  // We require view permission for the table group.
  permissionName := "view_" + st.EntityGroupName()
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
  pathPrefix := h.queryPrefix(entityType)
  morePath := strings.TrimPrefix(r.URL.Path, pathPrefix)
  // TODO - we might want to get rid of these defaults and require that the
  // next part of the path be specified.
  getOp := "column"     // Default op for GET request.
  postOp := "row"       // Default op for POST request.
  subOp := ""
  if morePath != "" {
    morePath = strings.TrimSuffix(morePath, "/")
    moreParts := strings.SplitN(morePath, "/", 2)
    getOp = moreParts[0]
    postOp = moreParts[0]
    if len(moreParts) == 2 && getOp=="summary" {
      subOp = moreParts[1]
    } else if len(moreParts) > 1 {
      msg := fmt.Sprintf("Too many additional path elements after %s (%v)", pathPrefix, morePath)
      http.Error(w, msg, http.StatusBadRequest)
      return
    }
  }
  glog.V(1).Infof("%s %s %s|%s", r.Method, entityType, getOp, postOp);
  switch r.Method {
    case http.MethodGet:
      switch getOp {
      case "column":
        h.stdGetColumns(w, r, st)
      case "row", "summary":
        h.stdGetRows(w, r, st, []queryParam{}, getOp, subOp)       // Get all rows.
      default:
        http.Error(w, "Invalid GET operation", http.StatusBadRequest)
        return
      }
    case http.MethodPost:
      switch postOp {
      case "column":
        h.stdGetColumns(w, r, st)
      case "row", "summary":
        var queryParams []queryParam
        if r.Body != nil {
          decoder := json.NewDecoder(r.Body)
          if err := decoder.Decode(&queryParams); err != nil {
            msg := fmt.Sprintf("Error decoding JSON query parameters: %v", err)
            http.Error(w, msg, http.StatusBadRequest)
            return
          }
        }
        h.stdGetRows(w, r, st, queryParams, postOp, subOp)
      default:
        http.Error(w, "Invalid POST operation", http.StatusBadRequest)
        return
      }
    default:
      http.Error(w, "Method must be GET or POST", http.StatusMethodNotAllowed)
  }
}

func (h *handler) stdGetColumns(w http.ResponseWriter, r *http.Request, st std) {
  entity := st.NewEntity()
  columnInfos := structsql.ColumnInfos(entity)

  // Special-case the column info if needed.
  dbrepos := h.config.DomainRepos.(*dbrepo.Repos)
  tableRepo, err := dbrepos.TableRepo(st.EntityTypeName())
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  columnsUpdater, ok := tableRepo.(structsql.ColumnInfosUpdater)
  if ok {
    columnInfos = columnsUpdater.UpdateColumnInfos(columnInfos)
  }

  result := &GetColumnResults{
    Columns: columnInfos,
  }
  apihttp.MarshalAndReply(w, result)
}

func (h *handler) stdGetRows(w http.ResponseWriter, r *http.Request, st std, queryParams []queryParam, op, subOp string) {
  tableName := st.EntityTypeName()
  query := "select * from " + tableName
  orderBy := ""
  if op == "summary" {
    query = st.SummaryQuery(subOp)      // The subOp is used as the summary format.
    orderBy = "summary"
  }
  whereVals := make([]interface{}, len(queryParams))
  if op == "summary" && strings.Contains(query, "ORDER BY") {
    queryParts := strings.Split(query, "ORDER BY")
    orderBy = queryParts[len(queryParts)-1]      // Get the last element of the slice.
    queryParts = queryParts[:len(queryParts)-1]  // Remove the order by value
    query = strings.Join(queryParts, "ORDER BY")
  }
  if len(queryParams) > 0 {
    whereParts := make([]string, len(queryParams))
    for i, qp := range queryParams {
      op := qp.CleanOp()
      whereParts[i] = qp.CleanName() + " " + op + " ?"
      whereVals[i] = qp.CleanValue()
    }
    query = query + " WHERE " + strings.Join(whereParts, " AND ")
  } else {
    // No query params, so no WHERE clause
  }
  if orderBy != "" {
    query = query + " ORDER BY " + orderBy
  }

  glog.V(1).Infof("query: %v", query)

  db := h.config.DomainRepos.(*dbrepo.Repos).DB()
  result, err := strsql.QueryStarAndCollect(db, query, whereVals...)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error executing sql: %v", err), http.StatusBadRequest)
    return
  }

  apihttp.MarshalAndReply(w, result)
}
