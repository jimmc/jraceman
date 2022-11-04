import { ApiManager, XhrOptions } from './api-manager.js'

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
  public static async loadKeySummaries(table: string, params?: {}) {
    const path = '/api/query/' + table + '/summary/'
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
}
