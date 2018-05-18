package http

import (
  "encoding/json"
  "fmt"
  "net/http"
)

func OkResponse(w http.ResponseWriter) {
  res := `{"status": "ok"}`
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(res))
}

func MarshalAndReply(w http.ResponseWriter, result interface{}) {
  b, err := json.MarshalIndent(result, "", "  ")
  if err != nil {
    http.Error(w, fmt.Sprintf("Failed to marshall json results: %v", err), http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Write(b)
}
