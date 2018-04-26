package domain

// The Diffs interface provides a way to retrieve the set of
// differences between two structs.
type Diffs interface {
  Modified() map[string]interface{}
}
