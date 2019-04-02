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

// ControlsOrderByItem is the information about OrderBy that is given to the user.
type ControlsOrderByItem struct {
  Name string
  Display string
}

/* ClientVisibleReports returns the list of reports and their attributes
 * that should be visible to a client.
 * Once we have user ids and authorizations, this function should accept
 * a user id and return only data that should be visible to that user.
 */
func ClientVisibleReports(reportRoots []string) ([]*ReportControls, error) {
  allControls := make([]*ReportControls, 0)
  for _, root := range reportRoots {
    controls, err := ClientVisibleReportsOne(root)
    if err != nil {
      return nil, err
    }
    allControls = append(allControls, controls...)
  }
  return allControls, nil
}

/* ClientVisibleReportsOne returns the list of reports and their user-visible attributes
 * from one root directory.
 */
func ClientVisibleReportsOne(templateDir string) ([]*ReportControls, error) {
  attrs, err := ReadTemplateAttrs(templateDir)
  if err != nil {
    return nil, err
  }
  reports := make([]*ReportControls, 0)
  for _, tplAttrs := range attrs {
    userOrderByList, err := extractUserOrderByList(tplAttrs)
    if err != nil {
      // If we get an error, don't add this report to the list, but still show other reports.
      glog.Errorf("Error decoding orderby in template %q: %v", tplAttrs.Name, err)
      continue
    }
    if tplAttrs.Display == "" {
      // Don't include templates with no Display attribute.
      continue
    }
    report := &ReportControls{
      Name: tplAttrs.Name,
      Display: tplAttrs.Display,
      Description: tplAttrs.Description,
      OrderBy: userOrderByList,
    }
    reports = append(reports, report)
  }
  return reports, nil
}

// extractUserOrderByList looks at the orderby attribute in the given template attributes
// and extacts from that the user-visible controls.
func extractUserOrderByList(tplAttrs *ReportAttributes) ([]ControlsOrderByItem, error) {
    r := []ControlsOrderByItem{}
    for _, v := range tplAttrs.OrderBy {
      r = append(r, ControlsOrderByItem{
        Name: v.Name,
        Display: v.Display,
      })
    }
    return r, nil
}
