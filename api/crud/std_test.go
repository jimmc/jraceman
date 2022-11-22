package crud_test

import (
  "net/http"
  "os"
  "testing"

  apitest "github.com/jimmc/jraceman/api/test"

  //goldenbase "github.com/jimmc/golden/base"
  //goldenhttp "github.com/jimmc/golden/http"
)

func TestStdBadMethod(t *testing.T) {
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-list") // Calls t.Fatal on error.
  defer cleanup()
  req := apitest.RequireNewRequest(t, "FOO", "/api/crud/site/", nil) // FOO is not a valid method.
  req = apitest.AddTestUser(req, "edit_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  if got, want := rr.Code, http.StatusMethodNotAllowed; got != want {
    t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }
}

func TestStdListOK(t *testing.T) {
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-list") // Calls t.Fatal on error.
  defer cleanup()
  req := apitest.RequireNewRequest(t, "GET", "/api/crud/site/", nil) // Calls t.Fatal on error.
  req = apitest.AddTestUser(req, "view_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  apitest.RequireHttpSuccess(t, req, rr)     // Calls t.Fatal if not http success.
  apitest.RequireBodyMatchesGolden(t, rr, "site-list")
}

func TestStdListNoUser(t *testing.T) {
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-list") // Calls t.Fatal on error.
  defer cleanup()
  req := apitest.RequireNewRequest(t, "GET", "/api/crud/site/", nil) // Calls t.Fatal on error.
  rr := apitest.ServeCrudRequest(req, repos)
  if got, want := rr.Code, http.StatusUnauthorized; got != want {
    t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }
}

func TestStdListNoPermission(t *testing.T) {
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-list") // Calls t.Fatal on error.
  defer cleanup()
  req := apitest.RequireNewRequest(t, "GET", "/api/crud/site/", nil) // Calls t.Fatal on error.
  req = apitest.AddTestUser(req, "none")        // Add a user, but without the right permission.
  rr := apitest.ServeCrudRequest(req, repos)
  if got, want := rr.Code, http.StatusUnauthorized; got != want {
    t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }
}

func TestStdPostWithID(t *testing.T) {
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-list")
  defer cleanup()

  // Create a site record.
  payloadfile, err := os.Open("testdata/site-create-id.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  req := apitest.RequireNewRequest(t, "POST", "/api/crud/site/S999", payloadfile)
  req = apitest.AddTestUser(req, "edit_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  if got, want := rr.Code, http.StatusBadRequest; got != want {
    t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }
}

func TestStdCreateBadJson(t *testing.T) {
  repos, cleanup := apitest.RequireEmptyDatabase(t)
  defer cleanup()

  // Create a site record.
  payloadfile, err := os.Open("testdata/site-create-bad-json.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  req := apitest.RequireNewRequest(t, "POST", "/api/crud/site/", payloadfile)
  req = apitest.AddTestUser(req, "edit_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  if got, want := rr.Code, http.StatusBadRequest; got != want {
    t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }
}

func TestStdCreateNoSave(t *testing.T) {
  repos, cleanup := apitest.RequireEmptyDatabase(t)
  defer cleanup()

  // Create a site record.
  payloadfile, err := os.Open("testdata/site-create-id.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer payloadfile.Close()

  req := apitest.RequireNewRequest(t, "POST", "/api/crud/site/", payloadfile)
  req = apitest.AddTestUser(req, "edit_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  if got, want := rr.Code, http.StatusBadRequest; got != want {
    t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }
}

func TestStdPutWithoutID(t *testing.T){
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-list")
  defer cleanup()

  updatepayloadfile, err := os.Open("testdata/site-update-id.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer updatepayloadfile.Close()
  req := apitest.RequireNewRequest(t, "PUT", "/api/crud/site/", updatepayloadfile)
  req = apitest.AddTestUser(req, "edit_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  if got, want := rr.Code, http.StatusBadRequest; got != want {
    t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }
}

func TestStdPatchWithoutID(t *testing.T){
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-list")
  defer cleanup()

  patchpayloadfile, err := os.Open("testdata/site-patch.payload")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer patchpayloadfile.Close()
  req := apitest.RequireNewRequest(t, "PATCH", "/api/crud/site/", patchpayloadfile)
  req = apitest.AddTestUser(req, "edit_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  if got, want := rr.Code, http.StatusBadRequest; got != want {
    t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }
}

func TestDeleteWithoutID(t *testing.T) {
  repos, cleanup := apitest.RequireDatabaseWithSqlFile(t, "site-create-id")
  defer cleanup()

  // Make sure the record we want is there as expected.
  siteId := "S1"
  req := apitest.RequireNewRequest(t, "GET", "/api/crud/site/" + siteId, nil)
  req = apitest.AddTestUser(req, "view_venue")
  rr := apitest.ServeCrudRequest(req, repos)
  apitest.RequireHttpSuccess(t, req, rr)
  apitest.RequireBodyMatchesGolden(t, rr, "site-delete-id-before")

  // Issue invalid delete request.
  req = apitest.RequireNewRequest(t, "DELETE", "/api/crud/site/", nil)
  req = apitest.AddTestUser(req, "edit_venue")
  rr = apitest.ServeCrudRequest(req, repos)
  if got, want := rr.Code, http.StatusBadRequest; got != want {
    t.Fatalf("HTTP response status for request %v: got %d, want %d\nBody: %v",
        req.URL, got, want, rr.Body.String())
  }

  // Make sure our record is still there.
  req = apitest.RequireNewRequest(t, "GET", "/api/crud/site/" + siteId, nil)
  req = apitest.AddTestUser(req, "view_venue")
  rr = apitest.ServeCrudRequest(req, repos)
  apitest.RequireHttpSuccess(t, req, rr)
  apitest.RequireBodyMatchesGolden(t, rr, "site-delete-id-before")
}
