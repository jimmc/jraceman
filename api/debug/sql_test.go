package debug_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jraceman/api/test"

  authlib "github.com/jimmc/auth/auth"
  goldenbase "github.com/jimmc/golden/base"
)

func TestGet(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunDebugTest("simple-sql", func() (*http.Request, error) {
    urlstr := "/api/debug/sql/?name=site-report&q=select+name,bar+from+test+where+foo='x' order by id"
    req, err := http.NewRequest("GET", urlstr, nil)
    if err != nil {
      return nil, err
    }
    req = apitest.AddTestUser(req, "edit_database")
    ac := &authlib.Config{TokenCookieName:"JRACEMAN_TOKEN"}
    err = authlib.AddTokenCookieForTesting(req, ac)
    if err != nil {
      return nil, err
    }
    return req, nil
  }), "RunDebugTest")
}

func TestPost(t *testing.T) {
  payloadfile, err := os.Open("testdata/simple-sql.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  goldenbase.FatalIfError(t, apitest.RunDebugTest("simple-sql", func() (*http.Request, error) {
    req, err := http.NewRequest("POST", "/api/debug/sql/", payloadfile)
    if err != nil {
      return nil, err
    }
    req = apitest.AddTestUser(req, "edit_database")
    ac := &authlib.Config{TokenCookieName:"JRACEMAN_TOKEN"}
    err = authlib.AddTokenCookieForTesting(req, ac)
    if err != nil {
      return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    return req, nil
  }), "RunDebugTest")
}
