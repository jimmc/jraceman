package report

import (
  "fmt"
)

// AttributesOrderByItem contains the OrderBy info loaded from a report template.
type AttributesOrderByItem struct {
  Name string
  Display string
  Sql string
}

// ControlsOrderByItem is the information about OrderBy that is given to the user.
type ControlsOrderByItem struct {
  Name string
  Display string
}

type ComputedOrderBy struct {
  Expr string    // The orderby sql corresponding to the OrderBy in the Options.
  Display string
  Clause string  // Blank if no OrderBy, else "ORDER BY " and the expression.
}

func findOrderByItem(orderByList []AttributesOrderByItem, orderByName string) (*AttributesOrderByItem, error) {
  for _, item := range orderByList {
    if item.Name == orderByName {
      return &item, nil
    }
  }
  return nil, fmt.Errorf("orderby item %q not found", orderByName)
}

// extractControlsOrderByItems looks at the orderby attribute in the given template attributes
// and extacts from that the user-visible controls.
func extractControlsOrderByItems(tplAttrs *ReportAttributes) ([]ControlsOrderByItem, error) {
  r := []ControlsOrderByItem{}
  for _, v := range tplAttrs.OrderBy {
    r = append(r, ControlsOrderByItem{
      Name: v.Name,
      Display: v.Display,
    })
  }
  return r, nil
}

func computeOrderBy(tplName string, options *ReportOptions, attrs *ReportAttributes) (*ComputedOrderBy, error) {
  if options == nil || options.OrderBy == "" {
    return &ComputedOrderBy{}, nil
  }
  if attrs == nil || len(attrs.OrderBy) == 0 {
    return nil, fmt.Errorf("invalid orderby option %q, template %s does not permit orderby",
        options.OrderBy, tplName)
  }
  orderByItem, err := findOrderByItem(attrs.OrderBy, options.OrderBy)
  if err != nil {
    return nil, fmt.Errorf("invalid orderby option %q for template %s",
        options.OrderBy, tplName)
  }
  computedOrderBy := &ComputedOrderBy{
    Expr: orderByItem.Sql,
    Display: orderByItem.Display,
  }
  if computedOrderBy.Expr != "" {
    computedOrderBy.Clause = "ORDER BY " + computedOrderBy.Expr
  }
  return computedOrderBy, nil
}
