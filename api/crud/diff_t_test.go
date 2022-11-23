package crud

import (
  "github.com/jimmc/jraceman/domain"
)

// PatchToDiffsForTest makes our private patchToDiffs function available to
// the crud_test package for unit testing.
func PatchToDiffsForTest(oldEntity, newEntity interface{}, patch interface{}) (domain.Diffs, bool, error) {
  return patchToDiffs(oldEntity, newEntity, patch)
}
