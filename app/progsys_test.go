package app

import (
  "fmt"
  "strings"
  "testing"

  dbtest "github.com/jimmc/jraceman/dbrepo/test"
)

func TestProgSysForEvent(t *testing.T) {
  tests := []struct{
    testName string
    setupName string
    eventId string
    expectError bool
    errorString string  // If expectError is true, this can have a piece of the expected error string.
  } {
    { "success", "progsys", "EV1", false, "" },
    { "no event", "progsys", "XXX", true, "can't find event" },
    { "no progression in event", "progsys", "EV2", true, "no progression" },
    { "progression not found", "progsys", "EV3", true, "can't find progression" },
    { "no class in progression", "progsys", "EV4", true, "no class" },
    { "complan not implemented", "progsys", "EV5", true, "Complan NYI" },
    { "unknown proression", "progsys", "EV6", true, "unknown progression class" },
    { "error loading simplan", "progsys", "EV7", true, "error loading simplan" },
  }
  for _, tt := range tests {
    t.Run(tt.testName, func(t *testing.T) {

      // Load the database.
      setupfilename := "testdata/" + tt.setupName + ".setup"
      dbRepos, cleanup, err := dbtest.ReposAndLoadFile(setupfilename)
      if err != nil {
        t.Fatalf(err.Error())
      }
      defer cleanup()
 
      // Call the function under test.
      ps, err := ProgSysForEvent(dbRepos, tt.eventId)

      // Check the results.
      if tt.expectError {
        if err == nil {
          t.Fatal("Expected error but did not get one")
        }
        if tt.errorString!="" && !strings.Contains(err.Error(),tt.errorString) {
          t.Fatal(fmt.Errorf("Expected error string to contain %q: %w", tt.errorString, err))
        }
      } else {
        if err != nil {
          t.Fatal(err)
        }
        if ps==nil {
          t.Fatal("Expected progSys but did not get one")
        }
      }
    })
  }
}
