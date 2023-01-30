package domain

import (
  "fmt"
  "strings"
)

// ProgressionRepo describes how Progression records are loaded and saved.
type ProgressionRepo interface {
  FindByID(ID string) (*Progression, error)
  List(offset, limit int) ([]*Progression, error)
  Save(*Progression) (string, error)
  UpdateByID(ID string, oldProgression, newProgression *Progression, diffs Diffs) error
  DeleteByID(ID string) error
}

// Progression defines a named progression for use in events.
type Progression struct {
  ID string
  Name string
  Class string
  Parameters *string
}

// ParmsAsMap parses the Parameters field and returns s
// map[string]string with all of the values.
// Each parameter is name=value, and parameters are separated by commas.
// There is no additional whitespace around either the equals or the commas.
func (p *Progression) ParmsAsMap() (map[string]string, error) {
  values := make(map[string]string)
  if p.Parameters == nil || *p.Parameters == "" {
    return values, nil
  }
  pkvs := strings.Split(*p.Parameters, ",")
  for _, kvs := range pkvs {
    kva := strings.Split(kvs, "=")
    if len(kva) != 2 {
      return nil, fmt.Errorf("Invalid syntax for progression parameter %s, should be name=value", kvs)
    } else {
      values[kva[0]] = kva[1]
    }
  }
  return values, nil
}

// ProgressionMeta provides funcions related to the Progression struct.
type ProgressionMeta struct {}

func (m *ProgressionMeta) EntityTypeName() string {
  return "progression"
}

func (m *ProgressionMeta) EntityGroupName() string {
  return "sport"
}

func (m *ProgressionMeta) NewEntity() interface{} {
  return &Progression{}
}
