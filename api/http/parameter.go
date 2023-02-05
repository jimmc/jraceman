package http

import (
  "encoding/json"
  "fmt"
  nethttp "net/http"
  "strconv"

  "github.com/golang/glog"
)

// For a POST, decode the body as JSON into a map.
// This is an easy way to get the passed-in parameters if you don't need
// the results to be in a struct.
func GetRequestParameters(r *nethttp.Request) (map[string]interface{}, error) {
  jsonBody := make(map[string]interface{}, 0)
  err := GetRequestParametersInto(r, &jsonBody)
  return jsonBody, err
}

// Decode the body of the request as JSON into the specified destination.
// The caller typically passes &x as destptr to return data into x,
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

// Get the value of a JSON parameter as a string.
// This should be called on the return value from GetRequestParameters.
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

// Get the value of a JSON parameter as an int, or a default value if not defined.
// This should be called on the return value from GetRequestParameters.
func GetJsonIntParameter(jsonBody map[string]interface{}, name string, dflt int) int {
  val, ok := jsonBody[name]
  if !ok {
    return dflt
  }
  s, ok := val.(string)
  if !ok {
    return dflt
  }
  n, err := strconv.Atoi(s)
  if err != nil {
    return dflt
  }
  return n
}

// Get the value of a JSON parameter as a boolean.
// This should be called on the return value from GetRequestParameters.
func GetJsonBoolParameter(jsonBody map[string]interface{}, name string, dflt bool) bool {
  val, ok := jsonBody[name]
  if !ok {
    return dflt
  }
  b, ok := val.(bool)
  if !ok {
    return dflt
  }
  return b
}
