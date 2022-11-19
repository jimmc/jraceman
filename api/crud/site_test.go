package crud_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jraceman/api/test"

  goldenbase "github.com/jimmc/golden/base"
  goldenhttp "github.com/jimmc/golden/http"
)

// TODO: Update, Delete, Patch

func TestList(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunCrudTest("site-list", func() (*http.Request, error) {
    r, err := http.NewRequest("GET", "/api/crud/site/", nil)
    if err != nil {
      return nil, err
    }
    r = apitest.AddTestUser(r, "view_venue")
    return r, nil
  }), "RunCrudTest")
}

// The same test as TestList, but easier to understand.
func TestListAlt(t *testing.T) {
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-list") // Calls t.Fatal on error.
  defer cleanup()
  req, err := http.NewRequest("GET", "/api/crud/site/", nil)
  if err != nil {
    t.Fatalf("http.NewRequest failed: %v", err)
  }
  req = apitest.AddTestUser(req, "view_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  apitest.RequireHttpSuccess(t, req, rr)     // Call t.Fatal if no http success.
  apitest.RequireBodyMatchesGolden(t, rr, "site-list")
}

func TestListLimit(t *testing.T) {
  r := apitest.NewCrudTester()
  goldenbase.FatalIfError(t, r.Init(), "Init")

  goldenbase.FatalIfError(t, goldenhttp.RunTestWith(r, "site-list-limit-1", func() (*http.Request, error) {
    r, err := http.NewRequest("GET", "/api/crud/site/?limit=1&offset=1", nil)
    if err != nil {
      return nil, err
    }
    r = apitest.AddTestUser(r, "view_venue")
    return r, nil
  }), "RunTestWith")

  goldenbase.FatalIfError(t, goldenhttp.RunTestWith(r, "site-list-limit-2", func() (*http.Request, error) {
    r, err := http.NewRequest("GET", "/api/crud/site/?limit=2&offset=2", nil)
    if err != nil {
      return nil, err
    }
    r = apitest.AddTestUser(r, "view_venue")
    return r, nil
  }), "RunTestWith")

  r.Close()
}

func TestGet(t *testing.T) {
  goldenbase.FatalIfError(t, apitest.RunCrudTest("site-get", func() (*http.Request, error) {
    r, err := http.NewRequest("GET", "/api/crud/site/S2", nil)
    if err != nil {
      return nil, err
    }
    r = apitest.AddTestUser(r, "view_venue")
    return r, nil
  }), "RunCrudTest")
}

func TestCreateWithID(t *testing.T) {
  r := apitest.NewCrudTester()
  goldenbase.FatalIfError(t, r.Init(), "Init")

  payloadfile, err := os.Open("testdata/site-create-id.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  goldenbase.FatalIfError(t, goldenhttp.RunTestWith(r, "site-create-id", func() (*http.Request, error) {
    r, err := http.NewRequest("POST", "/api/crud/site/", payloadfile)
    if err != nil {
      return nil, err
    }
    r = apitest.AddTestUser(r, "edit_venue")
    return r, nil
  }), "RunTestWith")

  goldenbase.FatalIfError(t, goldenhttp.RunTestWith(r, "site-create-id-get", func() (*http.Request, error) {
    r, err := http.NewRequest("GET", "/api/crud/site/S10", nil)
    if err != nil {
      return nil, err
    }
    r = apitest.AddTestUser(r, "view_venue")
    return r, nil
  }), "RunTestWith")

  r.Close()
}
