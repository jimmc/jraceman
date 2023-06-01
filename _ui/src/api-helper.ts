import { ApiManager, XhrOptions } from './api-manager.js'
import { PostError } from './message-log.js'
import { TableDesc, TableDescSupport, FKItem } from './table-desc.js'

// KeySummary is a summary line for one row of a table, typically for a foreign key.
export interface KeySummary {
  ID: string;
  Summary: string;
}

// ApiHelper provides common functions that access the api.
export class ApiHelper {
  // loadKeySummaries loads ID and a summary string for each row of
  // the specified table. This is useful for providing a list of choices
  // for a foreign key to that table.
  // If there is an error, this function will throw an XMLHttpRequest event.
  public static async loadKeySummaries(table: string, params?: {}, format?: string) {
    let path = '/api/query/' + table + '/summary/'
    if (format) {
      path = path + format + "/"
    }
    const options: XhrOptions = {}
    if (params) {
      options.method = "POST"
      options.params = params
    }
    const result = await ApiManager.xhrJson(path, options)
    const newKeyItems: KeySummary[] = [];
    newKeyItems.push({ID: "", Summary: ""});
    for (const row of result.Rows) {
      newKeyItems.push({ID: row[0], Summary: row[1]});
    }
    return newKeyItems
  }

  // loadTableDesc loads the table description (column info) for the given table.
  public static async loadTableDesc(tableName: string): Promise<TableDesc> {
    const td = {
      Table: tableName,
      Columns: []       // We will fill this in later
    } as TableDesc
    const result: TableDesc = await ApiManager.xhrJson('/api/query/' + tableName + '/')
    const cols = TableDescSupport.tableDescToCols(result);
    td.Columns = cols;
    for (let i = 0; i<cols.length; i++) {
      const col = cols[i]
      if (col.FKTable) {
        td.Columns[i].FKItems = [{ID: "", Summary: ""} as FKItem];
        await ApiHelper.loadFKChoices(td, i, col.FKTable)
      }
    }
    return td
  }

  // loadFKChoices loads the foreign key choice values from the API
  // for column i in the TableDesc and inserts those values into the
  // FKItems field of that column.
  static async loadFKChoices(td: TableDesc, i: number, table: string) {
    const path = '/api/query/' + table + "/summary/"
    const options = {}
    try {
      const result = await ApiManager.xhrJson(path, options)
      const newFKItems: FKItem[] = [];
      newFKItems.push({ID: "", Summary: ""});
      for (const row of result.Rows) {
        newFKItems.push({ID: row[0], Summary: row[1]});
      }
      td.Columns[i].FKItems = newFKItems
    } catch(e) {
      PostError("api-helper", "Error loading foreign key choices: " + e)
      console.error("Error: ", e)         // TODO
    }
  }
}
