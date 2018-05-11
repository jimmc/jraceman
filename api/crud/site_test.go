package crud_test

import (
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "reflect"
  "testing"

  "github.com/jimmc/jracemango/api/crud"
  "github.com/jimmc/jracemango/dbrepo/dbtesting"
  "github.com/jimmc/jracemango/domain"
)

// List, Get, Create, Update, Delete, Patch

func TestList(t *testing.T) {
  repos,err := dbtesting.ReposWithSetupFile("testdata/site.setup")
  if err!= nil {
    t.Fatal(err.Error())
  }
  defer repos.Close()
  config := &crud.Config{
    Prefix: "/api/crud/",
    DomainRepos: repos,
  }
  handler := crud.NewHandler(config)

  req, err := http.NewRequest("GET", "/api/crud/site/", nil)
  if err != nil {
    t.Fatalf("error creating list request: %v", err)
  }

  rr := httptest.NewRecorder()
  handler.ServeHTTP(rr, req)

  if got, want := rr.Code, http.StatusOK; got != want {
    t.Errorf("list call status: got %d, want %d", got, want)
  }

  body := rr.Body.String()
  if body == "" {
    t.Fatalf("list response body should not be empty")
  }
  var dat []domain.Site
  if err = json.Unmarshal(rr.Body.Bytes(), &dat); err != nil {
    t.Fatalf("error unmarshaling list body: %v", err)
  }

  anytown := "Anytown"
  somewhere := "Somewhere"
  phone1 := "800-555-1212"
  phone2 := "888-555-1234"
  expectedData := []domain.Site{
    {
      ID: "S1",
      Name: "Site One",
      City: &anytown,
      Phone: &phone1,
    },
    {
      ID: "S2",
      Name: "Site Two",
      City: &somewhere,
      Phone: &phone2,
    },
  }

  if got, want := dat, expectedData; !reflect.DeepEqual(got, want) {
    t.Fatalf("Wrong data, got %+v, want %+v", got, want)
  }
}
