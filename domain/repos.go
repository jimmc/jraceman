package domain

// The Repos interface ties together the individual repo interfaces for the
// different types so that methods within one of the type-specific repos can
// access the repos for the other types.
type Repos interface {
  Area() AreaRepo
  Challenge() ChallengeRepo
  Competition() CompetitionRepo
  Complan() ComplanRepo
  ComplanRule() ComplanRuleRepo
  ComplanStage() ComplanStageRepo
  Exception() ExceptionRepo
  Gender() GenderRepo
  LaneOrder() LaneOrderRepo
  Level() LevelRepo
  Person() PersonRepo
  Progression() ProgressionRepo
  ScoringRule() ScoringRuleRepo
  ScoringSystem() ScoringSystemRepo
  Simplan() SimplanRepo
  SimplanRule() SimplanRuleRepo
  SimplanStage() SimplanStageRepo
  Site() SiteRepo
  Stage() StageRepo
  Team() TeamRepo
}
