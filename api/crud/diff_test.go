package crud_test

import (
  "encoding/json"
  //"net/http"
  //"os"
  "testing"

  "github.com/jimmc/jraceman/api/crud"
  //apitest "github.com/jimmc/jraceman/api/test"

  //goldenbase "github.com/jimmc/golden/base"
  //goldenhttp "github.com/jimmc/golden/http"

  jsonpatch "gopkg.in/evanphx/json-patch.v3"
)

// Define a sample struct we use for testing.
type sample struct {
  I int
  S string
}

func TestDiffsHappyPath(t *testing.T) {
  oldEntity := &sample{ 123, "abc" }
  newEntity := &sample{}
  patchCmdBytes := []byte(`[{"op":"replace","path":"/S","value":"def"}]`)
  patch := make(jsonpatch.Patch,0)      // Patch is typedef for an array.
  json.Unmarshal(patchCmdBytes,&patch)
  diffs, eq, err := crud.PatchToDiffsForTest(oldEntity, newEntity, patch)
  if err != nil {
    t.Fatalf("Error from patchToDiffs: %v", err)
  }
  if eq {
    t.Fatalf("Expected to be not equal")
  }
  if diffs == nil {
    t.Fatal("Expected to have some diffs")
  }
}

func TestDiffsErrors(t *testing.T) {
  oldEntity := &sample{ 123, "abc" }
  newEntity := &sample{}
  unmarshalable := make(chan int)    // Make a value that can't be marshaled.
  unapplyablepatchCmdBytes := []byte(`[{"op":"replace","path":"/Foo","value":"x"}]`)
  unapplyablepatch := make(jsonpatch.Patch,0)      // Patch is typedef for an array.
  json.Unmarshal(unapplyablepatchCmdBytes,&unapplyablepatch)

  _, _, err := crud.PatchToDiffsForTest(oldEntity, newEntity, unmarshalable)
  if err == nil {
    t.Errorf("Expected error with unmarshalable patch")
  }
  _, _, err = crud.PatchToDiffsForTest(unmarshalable, newEntity, nil)
  if err == nil {
    t.Errorf("Expected error with unmarshalable oldEntity")
  }
  _, _, err = crud.PatchToDiffsForTest(oldEntity, newEntity, oldEntity)
  if err == nil {
    t.Errorf("Expected error with invalid patch")
  }
  _, _, err = crud.PatchToDiffsForTest(oldEntity, newEntity, unapplyablepatch)
  if err == nil {
    t.Errorf("Expected error with unapplyable patch")
  }
  _, _, err = crud.PatchToDiffsForTest(oldEntity, nil, nil)
  if err == nil {
    t.Errorf("Expected error unmarshaling to newEntity")
  }
}
