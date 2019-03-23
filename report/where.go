package report

import (
  "fmt"
  "strings"
)

// whereData is what gets returned to the template to build the sql string.
// whereclause and andclause start with a leading space and are intended to
// allow the template to just use that field and have it work whether the
// options contains where values or not.
type whereData struct {
  expr string               // Just the expression after WHERE.
  whereclause string        // WHERE and the expression.
  andclause string          // AND and the expression.
}

// whereDetails specifies how to create a where expression for one field.
type whereDetails struct {
  display string
  table string
  column string
  field string  // This can be used to override table.column such as if
                // the sql expression includes an "as" that names the table
                // something different from the standard name.
}

// whereGroups povides a standard set of expansions into where fields.
var whereGroups = map[string][]string {
  "event": { "event_id", "event_name", "event_number" },
  "person": { "person_id" },
  "team": { "team_id", "team_shortname", "team_name" },
}

// stdWheres provides the expansion of the standard 'where' field names to
// the details needed to bulid the sql expression for that field.
var stdWheres = map[string]whereDetails {
  "event_id": whereDetails{
    display: "Event",
    table: "event",
    column: "id",
  },
  "event_name": whereDetails{
    display: "Event Name",
    table: "event",
    column: "name",
  },
  "event_number": whereDetails{
    display: "Event Number",
    table: "event",
    column: "number",
  },
  "person_id": whereDetails{
    display: "Person",
    table: "person",
    column: "id",
  },
  "team_id": whereDetails{
    display: "Team",
    table: "team",
    column: "id",
  },
  "team_name": whereDetails{
    display: "Team Name",
    table: "team",
    column: "name",
  },
  "team_shortname": whereDetails{
    display: "Team Short Name",
    table: "team",
    column: "shortname",
  },
}

// emptyWhere is what we return to the template when there are no where
// fields specified in the options.
var emptyWhere = whereData{
  expr: "",
  whereclause: "",
  andclause: "",
}

// where generates the data that we return to the template.
func where(attrsMap map[string]interface{}, options *ReportOptions) (*whereData, error) {
  whereList, err := attrsToWhereList(attrsMap)
  if err != nil {
    return nil, err
  }
  whereList, err = expandWhereList(whereList)
  if err != nil {
    return nil, err
  }
  whereMap, err := whereListToMap(whereList)
  if err != nil {
    return nil, err
  }
  return whereMapToData(whereMap, options.WhereValues)
}

// attrsToWhereList extracts the list of standard where field names from
// the template attributes.
func attrsToWhereList(attrsMap map[string]interface{}) ([]string, error) {
  whereAttr := attrsMap["where"]
  if whereAttr == nil {
    return nil, nil
  }
  whereList, ok := whereAttr.([]interface{})
  if !ok {
    return nil, fmt.Errorf("'where' attribute is not array (it is %T)", whereAttr)
  }
  ss := make([]string, len(whereList))
  for i, v := range whereList {
    s, ok := v.(string)
    if !ok {
      return nil, fmt.Errorf("'where' attribute item %v at index %d is not string (it is %T)", v, i, v)
    }
    ss[i] = s
  }
  return ss, nil
}

// expandWhereList expands entries in whereGroups.
func expandWhereList(whereList []string) ([]string, error) {
  return expandStringList([]string{}, []string{}, whereList, whereGroups)
}

// expandStringList expands each element in 'in' and appends to 'out'.
// 'out' may be modified.
// chain may also be modified.
func expandStringList(chain, out, in []string, expansions map[string][]string) ([]string, error) {
  for _, s := range in {
    a, ok := expansions[s]
    if ok {
      // We keep the chain in an array because we want to print it out in order.
      // We assume nesting depths will typically be very short, so we are not
      // concerned about performance here.
      cycle := false
      for _, c := range chain {
        if c == s { cycle = true; break }
      }
      chain = append(chain, s)
      if cycle {
        return nil, fmt.Errorf("cycle detected: %v", chain)
      }
      var err error
      out, err = expandStringList(chain, out, a, expansions)
      if err != nil {
        return nil, err
      }
      chain = chain[:len(chain)-1]      // Remove the string we appended.
    } else {
      out = append(out, s)
    }
  }
  return out, nil
}

// whereListToMap expands a list of standard 'where' names to the
// equivalent customWhere structure that specifies display, table, column,
// and field values.
func whereListToMap(names []string) (map[string]whereDetails, error) {
  r := make(map[string]whereDetails)
  for _, name := range names {
    w, ok := stdWheres[name]
    if !ok {
      return nil, fmt.Errorf("No such standard where field %q", name)
    }
    r[name] = w
  }
  return r, nil
}

// whereMapToData generates the WHERE expression based on the given
// map of fields and the corresponding values.
// It also validates that each where value in the options
// matches a where field in the attributes.
// whereMap is the list of available fields for this report.
// whereValues is the set of values supplied in the options.
// TODO - validate that the type in the options matches the DB type.
func whereMapToData(whereMap map[string]whereDetails, whereValues map[string]WhereValue) (*whereData, error) {
  exprs := make([]string, len(whereValues))
  for whereName, whereValue := range whereValues {
    fieldDetails, ok := whereMap[whereName]
    if !ok {
      return nil, fmt.Errorf("where field %q is not valid for this template", whereName)
    }
    fieldExpr, err := whereString(fieldDetails, whereValue)
    if err != nil {
      return nil, err
    }
    exprs = append(exprs, fieldExpr)
  }
  expr := strings.Join(exprs, " && ")
  return &whereData{
    expr: expr,
    whereclause: " where " + expr,
    andclause: " && " + expr,
  }, nil
}

// whereString generates a where string for the specified field using the
// operator and value that came from the options.
func whereString(fieldDetails whereDetails, whereValue WhereValue) (string, error) {
  // TODO - check for valid op, based on field database type
  // TODO - check for value type, format accordingly.
  field := fieldDetails.field
  if field == "" {
    field = fmt.Sprintf("%s.%s", fieldDetails.table, fieldDetails.column)
  }
  switch v := whereValue.Value.(type) {
  case string:
    return fmt.Sprintf("%s %s %s", field, whereValue.Op, sqlQuotedString(v)), nil
  default:
    return fmt.Sprintf("%s %s %v", field, whereValue.Op, v), nil
  }
}

// sqlQuoteString puts single quotes around the string and escapes single-quotes within
// the string by doubling them.
// TODO - handle more databases? Quoting is database-specific.
func sqlQuotedString(s string) string {
  return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
