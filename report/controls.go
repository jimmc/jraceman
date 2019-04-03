package report

import (
  "github.com/golang/glog"
)

// ReportControls contains the information about a report template that is given to the user.
type ReportControls struct {
  Name string
  Display string
  Description string
  OrderBy []ControlsOrderByItem
}

/* ClientVisibleReports returns the list of reports and their attributes
 * that should be visible to a client.
 * Once we have user ids and authorizations, this function should accept
 * a user id and return only data that should be visible to that user.
 */
func ClientVisibleReports(reportRoots []string) ([]*ReportControls, error) {
  attrs, err := ReadAllTemplateAttrs(reportRoots)
  if err != nil {
    return nil, err
  }
  reports := make([]*ReportControls, 0)
  for _, tplAttrs := range attrs {
    controls := attrsToControls(tplAttrs)
    if controls == nil {
      continue
    }
    reports = append(reports, controls)
  }
  return reports, nil
}

// attrsToControls creates the user-visible ReportControls from the ReportAttributes
// read from a report template.
func attrsToControls(tplAttrs *ReportAttributes) *ReportControls {
  if tplAttrs.Display == "" {
    // Don't include templates with no Display attribute.
    return nil
  }
  userOrderByList, err := extractControlsOrderByList(tplAttrs)
  if err != nil {
    // If we get an error, don't add this report to the list, but still show other reports.
    glog.Errorf("Error decoding orderby in template %q: %v", tplAttrs.Name, err)
    return nil
  }
  report := &ReportControls{
    Name: tplAttrs.Name,
    Display: tplAttrs.Display,
    Description: tplAttrs.Description,
    OrderBy: userOrderByList,
  }
  return report
}
