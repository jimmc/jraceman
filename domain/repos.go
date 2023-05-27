package domain

// The Repos interface ties together the individual repo interfaces for the
// different types so that methods within one of the type-specific repos can
// access the repos for the other types.
type Repos interface {
  // Table types
  Area() AreaRepo
  Challenge() ChallengeRepo
  Competition() CompetitionRepo
  Complan() ComplanRepo
  ComplanRule() ComplanRuleRepo
  ComplanStage() ComplanStageRepo
  ContextOption() ContextOptionRepo
  Entry() EntryRepo
  Event() EventRepo
  Exception() ExceptionRepo
  Gender() GenderRepo
  Lane() LaneRepo
  LaneOrder() LaneOrderRepo
  Level() LevelRepo
  Meet() MeetRepo
  Option() OptionRepo
  Permission() PermissionRepo
  Person() PersonRepo
  Progression() ProgressionRepo
  Race() RaceRepo
  Registration() RegistrationRepo
  RegistrationFee() RegistrationFeeRepo
  Role() RoleRepo
  RolePermission() RolePermissionRepo
  RoleRole() RoleRoleRepo
  ScoringRule() ScoringRuleRepo
  ScoringSystem() ScoringSystemRepo
  SeedingList() SeedingListRepo
  SeedingPlan() SeedingPlanRepo
  Simplan() SimplanRepo
  SimplanRule() SimplanRuleRepo
  SimplanStage() SimplanStageRepo
  Site() SiteRepo
  Stage() StageRepo
  Team() TeamRepo
  User() UserRepo
  UserRole() UserRoleRepo

  // Composite types
  EventRaces() EventRacesRepo
  SimplanSys() SimplanSysRepo
}
