package report

import (
  "github.com/jimmc/jraceman/dbrepo"

  "github.com/jimmc/auth/permissions"

  "github.com/golang/glog"
)

// ReportControls contains the information about a report template that is given to the user.
type ReportControls struct {
  Name string
  Display string
  Description string
  OrderBy []ControlsOrderByItem
  Where []ControlsWhereItem
  permission string
}

/* ClientPermittedReports returns the list of reports and their attributes
 * for the reports that the client has permission to see.
 */
func ClientPermittedReports(dbrepos *dbrepo.Repos, reportRoots []string, perms *permissions.Permissions) ([]*ReportControls, error) {
  visibleReports, err := ClientVisibleReports(dbrepos, reportRoots)
  if err != nil {
    return nil, err
  }
  permittedReports := make([]*ReportControls, 0)
  for _, r := range visibleReports {
    if r.permission != "" && perms.HasPermission(permissions.Permission(r.permission)) {
      permittedReports = append(permittedReports, r)
    }
  }
  return permittedReports, nil
}

/* ClientVisibleReports returns the list of reports and their attributes
 * that should be visible to a client.
 */
func ClientVisibleReports(dbrepos *dbrepo.Repos, reportRoots []string) ([]*ReportControls, error) {
  attrs, err := ReadAllTemplateAttrs(reportRoots)
  if err != nil {
    return nil, err
  }
  reports := make([]*ReportControls, 0)
  for _, tplAttrs := range attrs {
    controls := attrsToControls(dbrepos, tplAttrs)
    if controls == nil {
      continue
    }
    reports = append(reports, controls)
  }
  return reports, nil
}

// attrsToControls creates the user-visible ReportControls from the ReportAttributes
// read from a report template.
func attrsToControls(dbrepos *dbrepo.Repos, tplAttrs *ReportAttributes) *ReportControls {
  if tplAttrs.Display == "" {
    // Don't include templates with no Display attribute.
    return nil
  }
  orderByItems, err := extractControlsOrderByItems(tplAttrs)
  if err != nil {
    // If we get an error, don't add this report to the list, but still show other reports.
    glog.Errorf("Error decoding orderby in template %q: %v", tplAttrs.Name, err)
    return nil
  }
  whereItems, err := extractControlsWhereItems(dbrepos, tplAttrs)
  if err != nil {
    // If we get an error, don't add this report to the list, but still show other reports.
    glog.Errorf("Error decoding where in template %q: %v", tplAttrs.Name, err)
    return nil
  }
  report := &ReportControls{
    Name: tplAttrs.Name,
    Display: tplAttrs.Display,
    Description: tplAttrs.Description,
    OrderBy: orderByItems,
    Where: whereItems,
    permission: tplAttrs.Permission,
  }
  return report
}
