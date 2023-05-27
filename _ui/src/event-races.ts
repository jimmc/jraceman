export interface RoundCount {
  Count: number,
  Round: number,
  StageName: string,
}

export interface RaceInfo {
  RaceID: string,
  StageName: string,
  StageNumber: number,
  IsFinal: boolean,
  Round: number,
  Section: number,
  AreaName: string,
  RaceNumber: number,
  LaneCount: number,
}

export interface EventRaces {
  Summary: string,
  EntryCount: number,
  GroupCount: number,
  GroupSize: number,
  RoundCounts: RoundCount[],
  Races: RaceInfo[],
}
