import {LitElement, html} from 'lit';
import {customElement, property} from 'lit/decorators.js';

import { ApiManager } from "./api-manager.js"
import { PostError } from "./message-log.js"
import { TableDesc, TableDescSupport, FKItem } from "./table-desc.js"
import "./table-edit.js"
import "./table-query.js"

/**
 * table-queryedit provides tabs with query and edit forms for a table.
 */
@customElement('table-queryedit')
export class TableQueryedit extends LitElement {

  @property({type: String})
  tableName = "(unset-in-table-queryedit-name)"

  //@property({type: Object})
  tableDesc: TableDesc = {
    Table: "(unset-in-table-queryedit-desc)",
    Columns: [],
  };

  // Once everything has had a chance to get set up, we kick off loadColumns.
  firstUpdated(changedProperties: Map<string,any>) {
    super.firstUpdated(changedProperties)
    setTimeout(this.loadColumns.bind(this), 0)  // Start loadColumns "in the background"
  }

  // Loads our column info from the API, builds a new TableDesc, and sets it
  // into this.tableDesc.
  async loadColumns() {
    const td = {
      Table: this.tableName,
      Columns: []       // We will fill this in later
    } as TableDesc
    const result: TableDesc = await ApiManager.xhrJson('/api/query/' + this.tableName + '/')
    const cols = TableDescSupport.tableDescToCols(result);
    td.Columns = cols;
    for (let i = 0; i<cols.length; i++) {
      const col = cols[i]
      if (col.FKTable) {
        td.Columns[i].FKItems = [{ID: "", Summary: ""} as FKItem];
        await this.loadFKChoices(td, i, col.FKTable)
      }
    }
    this.tableDesc = td
    this.requestUpdate()        // Our columns have been changed, update the screen
  }

  async loadFKChoices(td: TableDesc, i: number, table: string) {
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
      PostError("queryedit", "Error loading foreign key choices: " + e)
      console.error("Error: ", e)         // TODO
    }
  }

  render() {
    return html`
        <jraceman-tabs>
            <span slot="tab">Query</span>
            <section slot="panel">
              Table: ${this.tableName}
              <table-query tableDesc=${JSON.stringify(this.tableDesc)}></table-query>
            </section>
            <span slot="tab">Edit</span>
            <section slot="panel">
              <table-edit tableDesc=${JSON.stringify(this.tableDesc)}></table-edit>
            </section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'table-queryedit': TableQueryedit;
  }
}
