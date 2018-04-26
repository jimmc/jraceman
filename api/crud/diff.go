package crud

import (
  "github.com/jimmc/jracemango/domain"

  "gopkg.in/d4l3k/messagediff.v1"
)

type diffReader struct {
  messagediff *messagediff.Diff
}

func (d *diffReader) Modified() map[string]interface{} {
  mods := make(map[string]interface{})
  for k, v := range d.messagediff.Modified {
    // TODO - do we need to handle nested structs?
    mods[k.String()] = v
  }
  return mods
}

func deepDiff(oldEntity, newEntity interface{}) (diffs domain.Diffs, equal bool) {
  diff, equal := messagediff.DeepDiff(oldEntity, newEntity)
  return &diffReader{
    messagediff: diff,
  }, equal
}
