import { TableData } from './table-desc.js'

export interface RoundCount {
  Count: number,        // int
  Round: number,        // int
  StageName: string,
}

export interface RaceInfo {
  RaceID: string,
  StageName: string,
  StageNumber: number,  // int
  IsFinal: boolean,
  Round: number,        // int
  Section: number,      // int
  AreaName: string,
  RaceNumber: number,   // float
  LaneCount: number,    // int
  Lanes: TableData, // Not provided by EventRaces call; filled in by separately getting lanes
}

export interface EventRaces {
  Summary: string,
  EntryCount: number,   // int
  GroupCount: number,   // int
  GroupSize: number,    // int
  RoundCounts: RoundCount[],
  Races: RaceInfo[],
}
