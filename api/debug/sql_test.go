package debug_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"

  goldenhttp "github.com/jimmc/golden/http"
)

func TestGet(t *testing.T) {
  request := func() (*http.Request, error) {
    urlstr := "/api/debug/sql/?name=site-report&q=select+name,bar+from+test+where+foo='x' order by id"
    return http.NewRequest("GET", urlstr, nil)
  }
  if err := apitest.StartDebugToGolden("simple-sql", request); err != nil {
    t.Error(err.Error())
  }
}

func TestPost(t *testing.T) {
  repos, handler, err := apitest.StartDebugToSetup()
  if err != nil{
    t.Fatal(err.Error())
  }
  defer repos.Close()

  payloadfile, err := os.Open("testdata/simple-sql.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  request := func() (*http.Request, error) {
    req, err := http.NewRequest("POST", "/api/debug/sql/", payloadfile)
    if err != nil {
      return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    return req, nil
  }
  if err := goldenhttp.SetupToGolden(repos.DB(), handler, "simple-sql", request);
       err != nil {
    t.Error(err.Error())
  }
}
