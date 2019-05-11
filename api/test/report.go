package test

import (
  "net/http"

  "github.com/jimmc/jracemango/api/report"
)

// CreateReportHandler create an http handler for our report calls.
func CreateReportHandler(r *Tester, reportRoots []string) http.Handler {
  config := &report.Config{
    Prefix: "/api/report/",
    DomainRepos: r.repos,
    ReportRoots: reportRoots,
  }
  return report.NewHandler(config)
}

// NewReportTester returns a new Tester for testing our report calls.
func NewReportTester(reportRoots []string) *Tester {
  return NewTester(func (r *Tester) http.Handler {
    return CreateReportHandler(r, reportRoots)
  })
}

// RunReportTest creates a new Tester and runs a test for a report call.
func RunReportTest(basename string, reportRoots []string, callback func() (*http.Request, error)) error {
  r := NewReportTester(reportRoots)
  return r.Run(basename, callback)
}
