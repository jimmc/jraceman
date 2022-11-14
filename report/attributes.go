package report

import (
  "fmt"

  "github.com/jimmc/gtrepgen/gen"
)

// ReportAttributes contains the attributes loaded from a report template.
type ReportAttributes struct {
  Name string
  Display string
  Description string
  Permission string
  Where []string
  OrderBy []AttributesOrderByItem
}

/* ReadAllTemplateAttrs loads the attributes from all the template files
 * in all of the given directories.
 */
func ReadAllTemplateAttrs(reportRoots []string) ([]*ReportAttributes, error) {
  allAttrs := make([]*ReportAttributes, 0)
  for _, root := range reportRoots {
    dirAttrs, err := ReadTemplateAttrs(root)
    if err != nil {
      return nil, err
    }
    allAttrs = append(allAttrs, dirAttrs...)
  }
  return allAttrs, nil
}

/* ReadTemplateAttrs loads the attributes from all the template files in
 * the given directory.
 */
func ReadTemplateAttrs(templateDir string) ([]*ReportAttributes, error) {
  newDestPointer := func() interface{} {
    return &ReportAttributes{}
  }
  fileAttrs, err := gen.ReadDirFilesAttributesAs(templateDir, newDestPointer)
  if fileAttrs == nil {
    return nil, err
  }
  reportAttrs := []*ReportAttributes{}
  for _, fattrs := range fileAttrs {
    if fattrs.Err != nil {
      return nil, fmt.Errorf("for template %q received error %v", fattrs.Name, fattrs.Err)
    }
    attrs, ok := fattrs.Attributes.(*ReportAttributes)
    if !ok {
      return nil, fmt.Errorf("invalid data type for template %q", fattrs.Name)
    }
    attrs.Name = fattrs.Name
    reportAttrs = append(reportAttrs, attrs)
  }
  return reportAttrs, err
}

// GetAttributes loads our attributes from the template in one of the given report roots.
func GetAttributes(templateName string, reportRoots []string) (*ReportAttributes, error) {
  attrs := &ReportAttributes{}
  if err := gen.FindAndReadAttributesInto(templateName, reportRoots, attrs); err != nil {
    return nil, fmt.Errorf("reading template attributes: %v", err)
  }
  return attrs, nil
}
