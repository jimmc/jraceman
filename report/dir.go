package report

import (
  "fmt"
  "log"

  "github.com/jimmc/gtrepgen/gen"
)

type OrderByItem struct {
  Name string
  Display string
}

type ReportAttributes struct {
  Name string
  Display string
  OrderBy []OrderByItem
}

/* ClientVisibleReports returns the list of reports and their attributes
 * that should be visible to a client.
 * Once we have user ids and authorizations, this function should accept
 * a user id and return only data that should be visible to that user.
 */
func ClientVisibleReports(reportRoots []string) ([]*ReportAttributes, error) {
  allAttrs := make([]*ReportAttributes, 0)
  for _, root := range reportRoots {
    attrs, err := ClientVisibleReportsOne(root)
    if err != nil {
      return nil, err
    }
    allAttrs = append(allAttrs, attrs...)
  }
  return allAttrs, nil
}

/* ClientVisibleReportsOne returns the list of reports and their attributes
 * from one root directory.
 */
func ClientVisibleReportsOne(templateDir string) ([]*ReportAttributes, error) {
  attrs, err := ReadTemplateAttrs(templateDir)
  if err != nil {
    return nil, err
  }
  reports := make([]*ReportAttributes, 0)
  for _, tplAttrs := range attrs {
    name, ok := tplAttrs["name"].(string)
    if !ok {
      continue  // Name is not a string value, ignore this entry
    }
    display, ok := tplAttrs["display"].(string)
    if !ok {
      continue  // Display is not a string value, ignore this entry
    }
    userOrderByList, err := extractUserOrderByList(tplAttrs)
    if err != nil {
      // If we get an error, don't add this report to the list, but still show other reports.
      log.Printf("Error decoding orderby in template %q: %v", name, err)
      continue
    }
    report := &ReportAttributes{
      Name: name,
      Display: display,
      OrderBy: userOrderByList,
    }
    reports = append(reports, report)
  }
  return reports, nil
}

/* ReadTemplateAttrs loads the attributes from all the template files in
 * the given directory.
 */
func ReadTemplateAttrs(templateDir string) ([]map[string]interface{}, error) {
  attrs, err := gen.ReadDirFilesAttributes(templateDir)
  if attrs == nil {
    return nil, err
  }
  attrMaps := []map[string]interface{}{}
  for _, fattrs := range attrs {
    if fattrs.Err != nil {
      return nil, fmt.Errorf("for template %q received error %v", fattrs.Name, fattrs.Err)
    }
    fmap, ok := fattrs.Attributes.(map[string]interface{})
    if !ok {
      return nil, fmt.Errorf("invalid data type for template %q", fattrs.Name)
    }
    fmap["name"] = fattrs.Name
    attrMaps = append(attrMaps, fmap)
  }
  return attrMaps, err
}

// extractUserOrderByList looks at the orderby attribute in the given template attributes
// and extacts from that the user-visible fields.
// If there is no orderby attribute, it returns nil for both the value and the error.
func extractUserOrderByList(tplAttrs map[string]interface{}) ([]OrderByItem, error) {
    tplName := tplAttrs["name"]
    orderbyval, ok := tplAttrs["orderby"]
    if !ok {
      return nil, nil   // No orderby attribute, that's OK.
    }
    orderbyList, ok := orderbyval.([]interface{})
    if !ok {
      return nil, fmt.Errorf("orderby attribute for template %q is not []interface{}, it is %T", tplName, orderbyval)
    }
    r := []OrderByItem{}
    for _, v := range orderbyList {
      orderitem, ok := v.(map[string]interface{})
      if !ok {
        return nil, fmt.Errorf("value for orderby item %v is not map[string]interface{}, it is %T", v, v)
      }
      nameV := orderitem["name"]
      name, ok := nameV.(string)
      if !ok {
        return nil, fmt.Errorf("value for name in orderby item %s is not string, it is %T", v, nameV)
      }
      displayV := orderitem["display"]
      display, ok := displayV.(string)
      if !ok {
        return nil, fmt.Errorf("value for display in orderby item %s is not string, it is %T", name, displayV)
      }
      r = append(r, OrderByItem{
        Name: name,
        Display: display,
      })
    }
    return r, nil
}
