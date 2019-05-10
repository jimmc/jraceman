package debug_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jracemango/api/test"
)

func TestGet(t *testing.T) {
  apitest.RunDebugTest(t, "simple-sql", func() (*http.Request, error) {
    urlstr := "/api/debug/sql/?name=site-report&q=select+name,bar+from+test+where+foo='x' order by id"
    return http.NewRequest("GET", urlstr, nil)
  })
}

func TestPost(t *testing.T) {
  payloadfile, err := os.Open("testdata/simple-sql.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  apitest.RunDebugTest(t, "simple-sql", func() (*http.Request, error) {
    req, err := http.NewRequest("POST", "/api/debug/sql/", payloadfile)
    if err != nil {
      return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    return req, nil
  })
}
