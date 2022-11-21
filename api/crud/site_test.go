package crud_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jraceman/api/test"

  goldenbase "github.com/jimmc/golden/base"
  goldenhttp "github.com/jimmc/golden/http"
)

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
  apitest.RequireHttpSuccess(t, req, rr)     // Calls t.Fatal if not http success.
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

func TestListError(t *testing.T) {
  // We ask for an empty database so that the SQL query on the site table
  // will return an error.
  repos, cleanup := apitest.RequireEmptyDatabase(t)
  defer cleanup()
  req, err := http.NewRequest("GET", "/api/crud/site/", nil)
  if err != nil {
    t.Fatalf("http.NewRequest failed: %v", err)
  }
  req = apitest.AddTestUser(req, "view_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  if got, want := rr.Code, http.StatusBadRequest; got != want {
  t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
      req.URL, got, want, rr.Body.String())
  }
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

func TestCreate(t *testing.T) {
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-list")
  defer cleanup()

  // Create a site record.
  payloadfile, err := os.Open("testdata/site-create-id.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  req, err := http.NewRequest("POST", "/api/crud/site/", payloadfile)
  if err != nil {
    t.Fatalf("http.NewRequest failed: %v", err)
  }
  req = apitest.AddTestUser(req, "edit_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  apitest.RequireHttpSuccess(t, req, rr)
  apitest.RequireBodyMatchesGolden(t, rr, "site-create-id")
}

func TestUpdate(t *testing.T) {
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-create-id")
  defer cleanup()

  // Make sure the record we want is there as expected.
  siteId := "S1"
  req, err := http.NewRequest("GET", "/api/crud/site/" + siteId, nil)
  if err != nil {
    t.Fatalf("http.NewRequest failed: %v", err)
  }
  req = apitest.AddTestUser(req, "view_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  apitest.RequireHttpSuccess(t, req, rr)
  apitest.RequireBodyMatchesGolden(t, rr, "site-update-id-before")

  // Update the record.
  updatepayloadfile, err := os.Open("testdata/site-update-id.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer updatepayloadfile.Close()
  req, err = http.NewRequest("PUT", "/api/crud/site/" + siteId, updatepayloadfile)
  if err != nil {
    t.Fatalf("http.NewRequest failed: %v", err)
  }
  req = apitest.AddTestUser(req, "edit_venue")
  rr = apitest.ServeCrudRequest(req, repos)
  apitest.RequireHttpSuccess(t, req, rr)
  apitest.RequireBodyMatchesGolden(t, rr, "status-ok")

  // Make sure our record has been updated.
  req, err = http.NewRequest("GET", "/api/crud/site/" + siteId, nil)
  if err != nil {
    t.Fatalf("http.NewRequest failed: %v", err)
  }
  req = apitest.AddTestUser(req, "view_venue")
  rr = apitest.ServeCrudRequest(req, repos)
  apitest.RequireHttpSuccess(t, req, rr)
  apitest.RequireBodyMatchesGolden(t, rr, "site-update-id-after")
}

func TestDelete(t *testing.T) {
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-create-id")
  defer cleanup()

  // Make sure the record we want is there as expected.
  siteId := "S1"
  req, err := http.NewRequest("GET", "/api/crud/site/" + siteId, nil)
  if err != nil {
    t.Fatalf("http.NewRequest failed: %v", err)
  }
  req = apitest.AddTestUser(req, "view_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  apitest.RequireHttpSuccess(t, req, rr)
  apitest.RequireBodyMatchesGolden(t, rr, "site-delete-id-before")

  // Delete the record.
  req, err = http.NewRequest("DELETE", "/api/crud/site/" + siteId, nil)
  if err != nil {
    t.Fatalf("http.NewRequest failed: %v", err)
  }
  req = apitest.AddTestUser(req, "edit_venue")
  rr = apitest.ServeCrudRequest(req, repos)
  apitest.RequireHttpSuccess(t, req, rr)
  apitest.RequireBodyMatchesGolden(t, rr, "status-ok")

  // Make sure our record is no longer there.
  req, err = http.NewRequest("GET", "/api/crud/site/" + siteId, nil)
  if err != nil {
    t.Fatalf("http.NewRequest failed: %v", err)
  }
  req = apitest.AddTestUser(req, "view_venue")
  rr = apitest.ServeCrudRequest(req, repos)
  if got, want := rr.Code, http.StatusBadRequest; got != want {
    t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }
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
