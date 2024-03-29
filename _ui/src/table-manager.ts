import { LitElement, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'

import './jraceman-tabs.js'
import './table-edit.js'
import './table-query.js'
import './table-sheet.js'

import { ApiHelper } from './api-helper.js'
import { JracemanTabs } from './jraceman-tabs.js'
import { QueryResultsEvent, TableDesc } from './table-desc.js'
import { TableSheet } from './table-sheet.js'

/**
 * table-manager provides tabs with query, edit, and sheet panels for a table.
 */
@customElement('table-manager')
export class TableManager extends LitElement {

  @property({type: String})
  tableName = "(unset-in-table-manager-name)"

  //@property({type: Object})
  tableDesc: TableDesc = {
    Table: "(unset-in-table-manager-desc)",
    Columns: [],
  };

  // Once everything has had a chance to get set up, we kick off loadColumns.
  firstUpdated(changedProperties: Map<string,any>) {
    super.firstUpdated(changedProperties)
    setTimeout(this.loadColumns.bind(this), 0)  // Start loadColumns "in the background"
  }

  // Loads our table descriptor from the API and updates our component.
  // into this.tableDesc.
  async loadColumns() {
    this.tableDesc = await ApiHelper.loadTableDesc(this.tableName)
    this.requestUpdate()        // Our columns have been changed, update the screen
  }

  async onQueryToSheet(e:Event) {
    const evt = e as CustomEvent<QueryResultsEvent>
    const tableSheet = this.shadowRoot!.querySelector("table-sheet")! as TableSheet
    const results = evt.detail.results
    tableSheet.setSheetData(results)
    const tabs = this.shadowRoot!.querySelector("jraceman-tabs")! as JracemanTabs
    tabs.selectTabById("sheet-tab")
  }

  render() {
    return html`
        <jraceman-tabs>
            <span slot="tab">Query</span>
            <section slot="panel">
              Table: ${this.tableName}
              <table-query .tableDesc=${this.tableDesc}
                  @jraceman-query-to-sheet-event=${this.onQueryToSheet}
                  >
              </table-query>
            </section>
            <span slot="tab">Edit</span>
            <section slot="panel">
              <table-edit .tableDesc=${this.tableDesc}></table-edit>
            </section>
            <span slot="tab" id="sheet-tab">Sheet</span>
            <section slot="panel">
              <table-sheet .tableDesc=${this.tableDesc}></table-edit>
            </section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'table-manager': TableManager;
  }
}
