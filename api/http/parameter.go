package http

import (
  "encoding/json"
  "fmt"
  nethttp "net/http"

  "github.com/golang/glog"
)

func GetRequestParameters(r *nethttp.Request) (map[string]interface{}, error) {
  jsonBody := make(map[string]interface{}, 0)
  err := GetRequestParametersInto(r, &jsonBody)
  return jsonBody, err
}

// Decode the body of the request as JSON into the specified destination.
// The caller typically passed &x as destptr to return data into x,
// where x is an instance of the desired data type for the JSON data.
func GetRequestParametersInto(r *nethttp.Request, destPtr interface{}) error {
  contentType := r.Header.Get("content-type")
  glog.V(1).Infof("content-type: %v\n", contentType)
  if contentType != "application/json" {
    return fmt.Errorf("POST requires content-type=application/json")
  }
  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(destPtr); err != nil {
    return fmt.Errorf("Error decoding JSON body: %v", err)
  }
  return nil
}

func GetJsonStringParameter(jsonBody map[string]interface{}, name string) string {
  val, ok := jsonBody[name]
  if !ok {
    return ""
  }
  s, ok := val.(string)
  if !ok {
    return ""
  }
  return s
}
