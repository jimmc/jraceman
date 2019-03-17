package report

import (
  "fmt"

  "github.com/jimmc/gtrepgen/gen"
)

type ReportAttributes struct {
  Name string
  Display string
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
    report := &ReportAttributes{
      Name: name,
      Display: display,
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
