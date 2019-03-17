package report

import (
  "net/http"

  apihttp "github.com/jimmc/jracemango/api/http"
  "github.com/jimmc/jracemango/report"
)

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
  repAttrs, err := report.ClientVisibleReports(h.config.ReportRoots)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  apihttp.MarshalAndReply(w, repAttrs)
}
