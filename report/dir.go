package report

import (
  "fmt"
  "log"

  "github.com/jimmc/gtrepgen/gen"
)

type ReportAttributes struct {
  Name string
  Display string
  OrderBy map[string]string
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
    orderbymap, err := extractUserOrderByMap(tplAttrs)
    if err != nil {
      // If we get an error, don't add this report to the list, but still show other reports.
      log.Printf("Error decoding orderby in template %q: %v", name, err)
      continue
    }
    report := &ReportAttributes{
      Name: name,
      Display: display,
      OrderBy: orderbymap,
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

// extractUserOrderByMap looks at the orderby attribute in the given template attributes
// and extacts from that the user-visible fields.
// If there is no orderby attribute, it returns nil for both the value and the error.
func extractUserOrderByMap(tplAttrs map[string]interface{}) (map[string]string, error) {
    orderbyval, ok := tplAttrs["orderby"]
    if !ok {
      return nil, nil   // No orderby attribute, that's OK.
    }
    orderby, ok := orderbyval.(map[string]interface{})
    if !ok {
      return nil, fmt.Errorf("orderby attribute is not map[string]interface{}, it is %T", orderbyval)
    }
    // We return a map from the orderby key to the display value.
    r := map[string]string{}
    for k, v := range orderby {
      orderitem, ok := v.(map[string]interface{})
      if !ok {
        return nil, fmt.Errorf("value for orderby item %s is not map[string]interface{}, it is %T", k, v)
      }
      displayV := orderitem["display"]
      display, ok := displayV.(string)
      if !ok {
        return nil, fmt.Errorf("value for display in orderby item %s is not string, it is %T", k, displayV)
      }
      r[k] = display
    }
    return r, nil
}
