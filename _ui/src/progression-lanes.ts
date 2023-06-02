import { EventRaces } from './event-races.js'
import { QueryResultsColumnDesc, QueryResultsData, TableDesc } from './table-desc.js'

// ProgressionLanes provides support functions for dealing with the rows
// of the sheet in the ByEvent tab Entries/Progress task.
export class ProgressionLanes {

  // These are the fields that we collect for a progression lane. Some are hidden, most are read-only.
  static progressionLaneColumnDescriptors = [
    { Name: 'entryid', Label: 'EntryId', Type: 'string', FKTable: 'entry', FKItems: [], Hidden: true },
    { Name: 'fromlaneid', Label: 'FromLaneId', Type: 'string', FKTable: 'lane', FKItems: [], Hidden: true },
    { Name: 'tolaneid', Label: 'ToLaneId', Type: 'string', FKTable: 'lane', FKItems: [], Hidden: true },
    { Name: 'personid', Label: 'PersonId', Type: 'string', FKTable: 'person', FKItems: [] },
    { Name: 'teamid', Label: 'TeamId', Type: 'string', FKTable: 'team', FKItems: [], ReadOnly: true },
    { Name: 'nonscoring', Label: 'NonScoring', Type: 'boolean', FKTable: '', FKItems: [], ReadOnly: true },
    { Name: 'seed', Label: 'Seed', Type: 'int', FKTable: '', FKItems: [], ReadOnly: true },
    { Name: 'groupname', Label: 'Group', Type: 'string', FKTable: '', FKItems: [] },
    { Name: 'alternate', Label: 'Alternate', Type: 'boolean', FKTable: '', FKItems: [] },
    { Name: 'scratched', Label: 'Scratched', Type: 'boolean', FKTable: '', FKItems: [] },
    { Name: 'fromround', Label: 'Round', Type: 'int', FKTable: '', FKItems: [], ReadOnly: true },
    { Name: 'fromsection', Label: 'Section', Type: 'int', FKTable: '', FKItems: [], ReadOnly: true },
    { Name: 'fromlane', Label: 'Lane', Type: 'int', FKTable: '', FKItems: [], ReadOnly: true },
    { Name: 'fromresult', Label: 'Result', Type: 'float', FKTable: '', FKItems: [], ReadOnly: true },
    { Name: 'fromplace', Label: 'Place', Type: 'int', FKTable: '', FKItems: [], ReadOnly: true },
    { Name: 'fromscoreplace', Label: 'ScorePlace', Type: 'int', FKTable: '', FKItems: [], ReadOnly: true },
    { Name: 'toround', Label: 'ToRound', Type: 'int', FKTable: '', FKItems: [] },
    { Name: 'tosection', Label: 'ToSection', Type: 'int', FKTable: '', FKItems: [] },
    { Name: 'tolane', Label: 'ToLane', Type: 'int', FKTable: '', FKItems: [] },
    { Name: 'toresult', Label: 'ToResult', Type: 'float', FKTable: '', FKItems: [], ReadOnly: true },
  ]

  public static lanesFromRoundTableDesc(): TableDesc {  // Return value is suitable for EntriesProgress.sheetTableDesc
    const tableDesc: TableDesc = {
      Table: 'progressionLane',         // Not a real table
      Columns: ProgressionLanes.progressionLaneColumnDescriptors,
    }
    return tableDesc
  }

  public static collectLanesFromRound(
        entries: QueryResultsData,      // from EntriesProgress.entries
        eventRaces: EventRaces,         // from EntriesProgress.eventRaces
        fromRound: number               // from EntriesProgress.selectedRoundNumber
  ): QueryResultsData {                 // Return value is suitable for use as EntriesProgress.sheetQueryResults
    const toRound = fromRound + 1
    const entryColMap = ProgressionLanes.makeColumnIndexMap(entries.Columns)
    const entryRowMap: {[entryId:string]:any[]} = {}
    const progLaneColMap = ProgressionLanes.makeColumnIndexMap(ProgressionLanes.progressionLaneColumnDescriptors)
    const progLaneFieldCount = Object.keys(progLaneColMap).length
    const resultRows = []
    for (let entryRow of entries.Rows) {
      const entryId = entryRow[entryColMap["id"]]
      const progressionLaneData = new Array(progLaneFieldCount)
      resultRows.push(progressionLaneData)
      entryRowMap[entryId] = progressionLaneData
      // Copy fields out of entryRow into progressionLaneData
      progressionLaneData[progLaneColMap["entryid"]] = entryId
      progressionLaneData[progLaneColMap["personid"]] = entryRow[entryColMap["personid"]]
      progressionLaneData[progLaneColMap["groupname"]] = entryRow[entryColMap["groupname"]]
      progressionLaneData[progLaneColMap["alternate"]] = entryRow[entryColMap["alternate"]]
      progressionLaneData[progLaneColMap["scratched"]] = entryRow[entryColMap["scratched"]]
      // TODO - Nonscoring field comes from person record
      // progressionLaneData[progLaneColMap["nonscoring"]] = entryRow[entryColMap["nonscoring"]]
      // TODO - Seed field comes from person record
      // progressionLaneData[progLaneColMap["seed"]] = entryRow[entryColMap["seed"]]
      // TODO - TeamId field comes from person record
    }
    for (let race of eventRaces.Races) {
      if (race.Round == fromRound) {
        // TODO - look up entryId, copy fields from race into that row
      }
      if (race.Round == toRound) {
        // TODO - look up entryId, copy fields from race into that row
      }
    }
    const result:QueryResultsData = {
      Table: 'progressionLanes',        // Not a real table name
      Columns: ProgressionLanes.progressionLaneColumnDescriptors,        // TODO - strip down to name and type?
      Rows: resultRows,
    }
    return result
  }

  // makeColumnIndexMap takes an array of column descriptors (it just uses the Name field)
  // and returns a map from column name to array index, suitable for providing direct
  // by name to columns in the rows of QueryResultsData.
  static makeColumnIndexMap(columnDescriptors: QueryResultsColumnDesc[]): {[columnName:string]:number} {
    const columnMap: {[columnName:string]:number} = {}
    let colIndex = 0
    for (let colDesc of columnDescriptors) {
      columnMap[colDesc.Name] = colIndex
      colIndex++
    }
    return columnMap
  }
}
