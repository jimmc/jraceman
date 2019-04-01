package crud

// This file contains some wrappers around some external packages
// to make it easier to replace them if desired.

import (
  "encoding/json"
  "fmt"

  "github.com/jimmc/jracemango/domain"

  "github.com/golang/glog"
  "gopkg.in/d4l3k/messagediff.v1"
  jsonpatch "gopkg.in/evanphx/json-patch.v3"
)

type diffReader struct {
  messagediff *messagediff.Diff
}

// Modified returns a map of the fields that differ between the original
// old and new structs.
func (d *diffReader) Modified() map[string]interface{} {
  mods := make(map[string]interface{})
  for k, v := range d.messagediff.Modified {
    // TODO - do we need to handle nested structs?
    mods[k.String()] = v
  }
  return mods
}

// DeepDiff compares two structs and generates a diff.
func deepDiff(oldEntity, newEntity interface{}) (diffs domain.Diffs, equal bool) {
  diff, equal := messagediff.DeepDiff(oldEntity, newEntity)
  return &diffReader{
    messagediff: diff,
  }, equal
}

// PatchToDiffs takes a struct (oldEntity), and empty newEntity, and a datastructure representing a patch,
// applies the patch to the oldEntity, then takes a diff and returns that.
func patchToDiffs(oldEntity, newEntity interface{}, patch interface{}) (domain.Diffs, bool, error) {
  patchBytes, err := json.Marshal(patch)
  if err != nil {
    return nil, false, fmt.Errorf("error marshaling patch data: %v", err)
  }
  glog.V(1).Infof("patchBytes: %v", string(patchBytes))
  oldEntityBytes, err := json.Marshal(oldEntity)
  if err != nil {
    return nil, false, fmt.Errorf("error marshaling oldEntity: %v", err)
  }
  p, err := jsonpatch.DecodePatch(patchBytes)
  if err != nil {
    return nil, false, fmt.Errorf("error decoding patchBytes: %v", err)
  }
  newEntityBytes, err := p.Apply(oldEntityBytes)
  if err != nil {
    return nil, false, fmt.Errorf("error applying patch: %v", err)
  }
  err = json.Unmarshal(newEntityBytes, newEntity)
  if err != nil {
    return nil, false, fmt.Errorf("error unmarshaling newEntity: %v", err)
  }
  diffs, equal := deepDiff(oldEntity, newEntity)
  return diffs, equal, nil
}
