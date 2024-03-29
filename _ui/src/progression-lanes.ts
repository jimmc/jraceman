import { EventRaces } from './event-races.js'
import { ColumnDesc, TableDataColumnDesc, TableData, TableDesc } from './table-desc.js'

// ProgressionLanes provides support functions for dealing with the rows
// of the sheet in the ByEvent tab Entries/Progress task.
export class ProgressionLanes {

  // These are the fields that we collect for a progression lane. Some are hidden, most are read-only.
  static progressionLaneColumnDescriptors: ColumnDesc[] = [
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
        entryTableDesc: TableDesc,      // from entriesProgress.entryTableDesc
        entries: TableData,      // from EntriesProgress.entries
        eventRaces: EventRaces,         // from EntriesProgress.eventRaces
        fromRound: number               // from EntriesProgress.selectedRoundNumber
  ): TableData {                 // Return value is suitable for use as EntriesProgress.sheetQueryResults
    const toRound = fromRound + 1
    const entryColMap = ProgressionLanes.makeColumnIndexMap(entries.Columns)
    const entryRowMap: {[entryId:string]:any[]} = {}
    const progLaneColMap = ProgressionLanes.makeColumnIndexMap(ProgressionLanes.progressionLaneColumnDescriptors)
    const progLaneFieldCount = Object.keys(progLaneColMap).length
    let laneColMap: {[columnName:string]:number} = {}
    const resultRows = []
    // Plug in the FKItems from entryTableDesc into progressionLaneColumnDescriptors
    ProgressionLanes.progressionLaneColumnDescriptors[progLaneColMap['personid']].FKItems =
        entryTableDesc.Columns[entryColMap['personid']].FKItems
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
      if (race.Lanes && race.Lanes.Rows.length>0) {
        if (Object.keys(laneColMap).length==0) {
          laneColMap = ProgressionLanes.makeColumnIndexMap(race.Lanes.Columns)
        }
        for (let laneRow of race.Lanes.Rows) {
          const laneId = laneRow[laneColMap['id']]
          const laneEntryId = laneRow[laneColMap['entryid']]
          const progressionLaneData = entryRowMap[laneEntryId]
          if (race.Round == fromRound) {
            progressionLaneData[progLaneColMap['fromlaneid']] = laneId
            progressionLaneData[progLaneColMap['fromround']] = laneRow[laneColMap['round']]
            progressionLaneData[progLaneColMap['fromsection']] = laneRow[laneColMap['section']]
            progressionLaneData[progLaneColMap['fromlane']] = laneRow[laneColMap['lane']]
            progressionLaneData[progLaneColMap['fromresult']] = laneRow[laneColMap['result']]
            progressionLaneData[progLaneColMap['fromplace']] = laneRow[laneColMap['place']]
            progressionLaneData[progLaneColMap['fromscoreplace']] = laneRow[laneColMap['scoreplace']]
          }
          if (race.Round == toRound) {
            progressionLaneData[progLaneColMap['tolaneid']] = laneId
            progressionLaneData[progLaneColMap['toround']] = laneRow[laneColMap['round']]
            progressionLaneData[progLaneColMap['tosection']] = laneRow[laneColMap['section']]
            progressionLaneData[progLaneColMap['tolane']] = laneRow[laneColMap['lane']]
            progressionLaneData[progLaneColMap['toresult']] = laneRow[laneColMap['result']]
          }
        }
      }
    }
    const result:TableData = {
      Table: 'progressionLanes',        // Not a real table name
      Columns: ProgressionLanes.progressionLaneColumnDescriptors,        // TODO - strip down to name and type?
      Rows: resultRows,
    }
    return result
  }

  // makeColumnIndexMap takes an array of column descriptors (it just uses the Name field)
  // and returns a map from column name to array index, suitable for providing direct
  // by name to columns in the rows of TableData.
  static makeColumnIndexMap(columnDescriptors: TableDataColumnDesc[]): {[columnName:string]:number} {
    const columnMap: {[columnName:string]:number} = {}
    let colIndex = 0
    for (let colDesc of columnDescriptors) {
      columnMap[colDesc.Name] = colIndex
      colIndex++
    }
    return columnMap
  }
}
