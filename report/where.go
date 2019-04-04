package report

import (
  "fmt"
  "strings"
)

// ControlsWhereItem is the 'where' info that we pass to the user.
type ControlsWhereItem struct {
  Name string
  Display string
}

// OptionsWhereItem contains the value specified in ReportOptions for one where field.
type OptionsWhereItem struct {
  Op string     // The comparison operation to use for this field.
  Value interface{}     // The value to use on the RHS of the comparison.
}

// ComputedWhere is what gets returned to the template to build the sql string.
// WhereClause and AndClause start with a leading space and are intended to
// allow the template to just use that field and have it work whether the
// options contains where values or not.
type ComputedWhere struct {
  Expr string               // Just the expression after WHERE.
  WhereClause string        // WHERE and the expression.
  AndClause string          // AND and the expression.
}

// whereDetails specifies how to create a where expression for one field.
// This information comes from the template attributes.
type whereDetails struct {
  display string
  table string
  column string
  field string  // This can be used to override table.column such as if
                // the sql expression includes an "as" that names the table
                // something different from the standard name.
}

// whereGroups povides a standard set of expansions into where fields.
// The template attributes can include one of the keywords in this map,
// which then get expanded into the corresponding list of where fields.
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
  "site_id": whereDetails{
    display: "Site",
    table: "site",
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
var emptyWhere = ComputedWhere{
  Expr: "",
  WhereClause: "",
  AndClause: "",
}

// extractControlsWhereItems generates the 'where' info that we pass to the user.
func extractControlsWhereItems(tplAttrs *ReportAttributes) ([]ControlsWhereItem, error) {
  whereList, err := expandWhereList(tplAttrs.Where)
  if err != nil {
    return nil, err
  }
  whereMap, err := whereListToMap(whereList)
  if err != nil {
    return nil, err
  }
  r := []ControlsWhereItem{}
  for _, name := range whereList {
    details := whereMap[name]
    item := ControlsWhereItem{
      Name: name,
      Display: details.display,
    }
    r = append(r, item)
  }
  return r, nil
}

// computeWhere computes the 'where' data that we supply to the template during generation.
func computeWhere(attrs *ReportAttributes, options *ReportOptions) (*ComputedWhere, error) {
  whereList, err := expandWhereList(attrs.Where)
  if err != nil {
    return nil, err
  }
  whereMap, err := whereListToMap(whereList)
  if err != nil {
    return nil, err
  }
  if options == nil {
    options = &ReportOptions{}
  }
  whereListInUse := extractWhereListInUse(whereList, options.WhereValues)
  return whereMapToData(whereMap, whereListInUse, options.WhereValues)
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

// extractWhereListInUse takes a list of all defined where fields for a
// template and returns a list that only contains the fields that are
// used in the options.
func extractWhereListInUse(whereList []string, whereValues map[string]OptionsWhereItem) []string {
  r := []string{}
  for _, s := range whereList {
    if _, ok := whereValues[s]; ok {
      r = append(r, s)
    }
  }
  return r
}

// whereMapToData generates the WHERE expression based on the given
// map of fields and the corresponding values.
// It also validates that each where value in the options
// matches a where field in the attributes.
// whereMap is the list of available fields for this report.
// whereListInUse is the list of field names from whereValues,
// which we use to order the where expressions.
// whereValues is the set of values supplied in the options.
// TODO - validate that the type in the options matches the DB type.
func whereMapToData(whereMap map[string]whereDetails, whereListInUse []string, whereValues map[string]OptionsWhereItem) (*ComputedWhere, error) {
  exprs := []string{}
  for whereName, _ := range whereValues {
    _, ok := whereMap[whereName]
    if !ok {
      return nil, fmt.Errorf("where field %q is not valid for this template", whereName)
    }
  }
  for _, whereName := range whereListInUse {
    whereValue := whereValues[whereName]
    fieldDetails := whereMap[whereName]
    fieldExpr, err := whereString(fieldDetails, whereValue)
    if err != nil {
      return nil, err
    }
    exprs = append(exprs, fieldExpr)
  }
  expr := strings.Join(exprs, " AND ")
  if expr == "" {
    return &emptyWhere, nil
  }
  return &ComputedWhere{
    Expr: expr,
    WhereClause: " where " + expr,
    AndClause: " AND " + expr,
  }, nil
}

// whereString generates a where string for the specified field using the
// operator and value that came from the options.
func whereString(fieldDetails whereDetails, whereValue OptionsWhereItem) (string, error) {
  // TODO - check for valid op, based on field database type
  // TODO - check for value type, format accordingly.
  field := fieldDetails.field
  if field == "" {
    field = fmt.Sprintf("%s.%s", fieldDetails.table, fieldDetails.column)
  }
  op, err := whereOpStr(whereValue.Op)
  if err != nil {
    return "", fmt.Errorf("invalid op %q for field %s", whereValue.Op, fieldDetails.display)
  }
  switch v := whereValue.Value.(type) {
  case string:
    return fmt.Sprintf("%s %s %s", field, op, sqlQuotedString(v)), nil
  default:
    return fmt.Sprintf("%s %s %v", field, op, v), nil
  }
}

// whereOpStr converts the Op field from a OptionsWhereItem to the string to be used
// in the sql for that expression.
func whereOpStr(op string) (string, error) {
  switch op {
  case "eq": return "=", nil
  case "ne": return "!=", nil
  case "gt": return ">", nil
  case "ge": return ">=", nil
  case "lt": return "<", nil
  case "le": return "<=", nil
  // TODO: add "in" operator
  default: return "", fmt.Errorf("unknown op %q", op)
  }
}

// sqlQuoteString puts single quotes around the string and escapes single-quotes within
// the string by doubling them.
// TODO - handle more databases? Quoting is database-specific.
func sqlQuotedString(s string) string {
  return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
