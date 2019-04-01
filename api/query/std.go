package query

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strings"

  apihttp "github.com/jimmc/jracemango/api/http"
  "github.com/jimmc/jracemango/dbrepo"
  "github.com/jimmc/jracemango/dbrepo/strsql"
  "github.com/jimmc/jracemango/dbrepo/structsql"

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
  return "ERROR"        // This will result in invalid SQL.
}

func (qp *queryParam) CleanValue() interface{} {
  return qp.Value
}

// Std defines the type-specific methods that are needed by the generic
// query code.
type std interface {
  EntityTypeName() string   // such as "site"
  NewEntity() interface{}       // such as a new Site
}

// Stdquery is an http handler that takes care of most of the work for
// query requests on our domain types. The API query handler for a specific
// type turns around and calls this handler with a type-specific std that
// defines the type-specific methods needed by this handler.
func (h *handler) stdquery(w http.ResponseWriter, r *http.Request, st std) {
  // TODO - check authorization
  entityType := st.EntityTypeName()
  pathPrefix := h.queryPrefix(entityType)
  morePath := strings.TrimPrefix(r.URL.Path, pathPrefix)
  if morePath != "" {
    http.Error(w, "Additional path elements may not be specified after "+ pathPrefix, http.StatusBadRequest)
    return
  }
  glog.V(1).Infof("%s %s", r.Method, entityType);
  switch r.Method {
    case http.MethodGet:
      h.stdGetColumns(w, r, st)
    case http.MethodPost:
      h.stdList(w, r, st)
    default:
      http.Error(w, "Method must be POST", http.StatusMethodNotAllowed)
  }
}

func (h *handler) stdGetColumns(w http.ResponseWriter, r *http.Request, st std) {
  entity := st.NewEntity()
  columnInfos := structsql.ColumnInfos(entity)
  result := &GetColumnResults{
    Columns: columnInfos,
  }
  apihttp.MarshalAndReply(w, result)
}

func (h *handler) stdList(w http.ResponseWriter, r *http.Request, st std) {
  decoder := json.NewDecoder(r.Body)
  var queryParams []queryParam
  if err := decoder.Decode(&queryParams); err != nil {
    msg := fmt.Sprintf("Error decoding JSON query parameters: %v", err)
    http.Error(w, msg, http.StatusBadRequest)
    return
  }
  tableName := st.EntityTypeName()
  query := "select * from " + tableName
  whereVals := make([]interface{}, len(queryParams))
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

  glog.V(1).Infof("query: %v", query)

  db := h.config.DomainRepos.(*dbrepo.Repos).DB()
  result, err := strsql.QueryStarAndCollect(db, query, whereVals...)
  if err != nil {
    http.Error(w, fmt.Sprintf("Error executing sql: %v", err), http.StatusBadRequest)
    return
  }

  apihttp.MarshalAndReply(w, result)
}
