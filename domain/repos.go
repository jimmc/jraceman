package domain

// The Repos interface ties together the individual repo interfaces for the
// different types so that methods within one of the type-specific repos can
// access the repos for the other types.
type Repos interface {
  Area() AreaRepo
  Site() SiteRepo
}
