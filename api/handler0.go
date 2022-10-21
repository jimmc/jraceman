package api

import (
  "fmt"
  "net/http"

  apihttp "github.com/jimmc/jraceman/api/http"
)

// handler0 is the handler for api0 calls, which are api calls
// that can be executed without authentication.
type handler0 struct {
    prefix string
    version string
}

// NewHandler0 creates the http handler that is used to route api0 requests.
func NewHandler0(prefix, version string) http.Handler {
  h := handler0{
    prefix: prefix,
    version: version,
  }
  mux := http.NewServeMux()

  mux.HandleFunc(h.api0Prefix("version"), h.sendVersion)

  return mux
}

func (h *handler0) sendVersion(w http.ResponseWriter, r *http.Request) {
  result := h.version
  apihttp.MarshalAndReply(w, result)
}

// ApiPrefix composes our prefix with the next path component so that we can
// provide the right prefix to the handler that handles that next
// path component.
func (h *handler0) api0Prefix(s string) string {
  return fmt.Sprintf("%s%s/", h.prefix, s)
}
