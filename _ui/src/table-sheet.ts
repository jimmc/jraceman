import { LitElement, html, css } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { PropertyValues } from 'lit-element'
import { when } from 'lit/directives/when.js'

import './sheet-editor.js'

import { ApiManager, XhrOptions } from './api-manager.js'
import { JracemanDialog } from './jraceman-dialog.js'
import { PostError } from './message-log.js'
import { QueryFields } from './query-fields.js'
import { SheetEditor } from './sheet-editor.js'
import { TableCustom } from './table-custom.js'
import { TableDesc, TableDescSupport, ColumnDesc, QueryResultsData } from './table-desc.js'

/**
 * table-sheet provides a panel to edit fields in multiple rows and columns.
 */
@customElement('table-sheet')
export class TableSheet extends LitElement {
  static styles = css`
  `;

  @property({type: Object})
  tableDesc: TableDesc = {
    Table: "(unset-in-table-sheet)",
    Columns: [],
  };

  @property({type: Object /*, notify: true*/})
  queryResults: QueryResultsData = {
    Table: "(unset-in-query-results)",
    Columns: [],
    Rows: [],
  };

  @property({type: Number})
  selectedRowIndex: number = -1;

  @property({type: Boolean})
  haveResults = false

  sheetEditor?: SheetEditor

  queryFields?: QueryFields

  firstUpdated(changedProperties:PropertyValues<this>) {
    super.firstUpdated(changedProperties);
    this.sheetEditor = this.shadowRoot!.querySelector("sheet-editor")! as SheetEditor
    this.queryFields = this.shadowRoot!.querySelector("query-fields")! as QueryFields
    this.sheetEditor!.setFieldUpdatedCallback(this.onFieldUpdated.bind(this))
  }

  // getSelectElement gets an HTMLSelectElement by selector.
  getSelectElement(selector: string) {
    const shadowRoot = this.shadowRoot
    if (shadowRoot == null) {
      console.error("shadowRoot is null")
      return null
    }
    return shadowRoot.querySelector(selector) as HTMLSelectElement
  }

  // getSelectValue gets the value of a <select> element.
  getSelectValue(selector: string) {
    const sel = this.getSelectElement(selector)
    if (sel == null) {
      console.error("select element is null in getSelectValue")
      return null
    }
    return sel.value
  }

  // setSelectValue sets the value of a <select> element.
  setSelectValue(selector: string, val: string) {
    const sel = this.getSelectElement(selector)
    if (sel == null) {
      console.error("select element is null in setSelectValue")
      return
    }
    sel.value = val
  }

  clear() {
    for (let col of this.tableDesc['Columns']) {
      const name = col.Name;
      this.setSelectValue("#val_"+name, '');
      this.setSelectValue("#op_"+name, 'eq');
    }
  }

  // Set the given data into our sheet.
  async setSheetData(results: QueryResultsData) {
    console.log("TableSheet.setSheetData", results)
    const maxSafeRows = 500
    if (results.Rows.length > maxSafeRows) {
      const msg = "This query is returning " + results.Rows.length + " rows. " +
          "This may cause your browser to hang or crash. " +
          "You can use the Query tab and Edit-in-Sheet to select a smaller set of rows."
      const buttons =
          [ "Cancel", "Display all data anyway", "Display first " + maxSafeRows + " rows"]
      const b = await JracemanDialog.messageDialog("Data size warning", msg, buttons)
      if (b <= 0) {
        // Canceled, don't show the rows
        return
      }
      if (b == 2) {
        results.Rows = results.Rows.slice(0, maxSafeRows)
      }
      // if b==1, just continue and display all rows.
    }
    this.haveResults = true;
    this.selectedRowIndex = -1
    this.queryResults = results;
    this.haveResults = true;
  }

  async search() {
    console.log("TableSheet.search begin");
    const params = this.queryFields!.fieldsAsParams()
    const options: XhrOptions = {
      method: "POST",
      params: params,
    }
    const queryPath = '/api/query/' + this.tableDesc.Table + '/';
    try {
      const result = await ApiManager.xhrJson(queryPath, options);
      if (result && !result.Table) {
        result.Table = this.tableDesc.Table;
      }
      this.setSheetData(result)
    } catch(e) {
      const evt = e as XMLHttpRequest
      PostError("query", evt.responseText)
      console.log("Error in table query", e)
      return
    }
    console.log("TableSheet.search results", this.queryResults);
  }

  async rowSelected(e:CustomEvent) {
    console.log("Row selected", e)
    this.selectedRowIndex = e.detail as number;
  }

  async rowDeleted(e:CustomEvent) {
    console.log("Row deleted", e)
    const rowIndex = e.detail as number;
    this.queryResults.Rows.splice(rowIndex, 1)
    this.requestUpdate()
    // TODO - send a delete request to the API
  }

  // onFieldUpdated is called from our sheet-editor when the user changes a field.
  // We use a callback here rather than an event so that we can return a value
  // indicating an error. On success, the return value is null.
  // This function must confirm to SheetEditor.FieldUpdatedCallback.
  async onFieldUpdated(tableDesc: TableDesc, _tableData: QueryResultsData,
      _rowIndex: number, colIndex: number, id: string, colVal: string) {
    console.log("in TableSheet.onFieldUpdated")
    const col = tableDesc.Columns[colIndex]
    const name = col.Label      // TODO - Is this the right way to get the field name?
    console.log("Type of field " + name + " is " + col.Type)
    // For non-string fields, convert from the string in the form
    // to the appropriate data type for the field.
    const convertedColVal = TableDescSupport.convertToType(colVal, col.Type)
    const updatePath = '/api/crud/' + tableDesc.Table + '/' + id
    const method = "PATCH"
    const options: XhrOptions = {
      method: method,
      params: [
        { "op": "replace", "path": "/" + name, "value": convertedColVal },
      ]
    }
    try {
      const result = await ApiManager.xhrJson(updatePath, options)
      if (result && !result.Table) {
        result.Table = tableDesc.Table
      }
      console.log(result)
      return null       // No error
    } catch(e) {
      console.error("Error:", e/*.responseText*/)
      return e
    }
  }

  async edit() {
    this.sheetEditor!.editSelectedRow();
  }

  async delete() {
    this.sheetEditor!.deleteSelectedRow();
  }

  isStringColumn(colType: string) {
    return colType == "string";
  }

  // filterFieldTableDesc generates a TableDesc with a ColDesc for each
  // column for which we want to provide a filter choice.
  // Typically there is either zero or one filter field.
  filterFieldTableDesc() {
    const filterColumnName = TableCustom.sheetFilterFieldName(this.tableDesc.Table)
    let filterColumns: ColumnDesc[] = []
    if (filterColumnName) {
      for (let col of this.tableDesc.Columns) {
        if (col.Name == filterColumnName) {
          filterColumns.push(col)
        }
      }
    }
    const ffTableDesc: TableDesc = {
      Table: this.tableDesc.Table,
      Columns: filterColumns,
    }
    return ffTableDesc
  }

  render() {
    return html`
        <form>
          ${when(this.haveResults, ()=>html`[${this.queryResults.Rows.length}]`)}
          <query-fields .tableDesc=${this.filterFieldTableDesc()} tableClass=inline>
          </query-fields>
          <button type=button @click="${this.search}">Search</button>
          <button type=button @click="${this.edit}" ?disabled="${this.selectedRowIndex<0}">Edit</button>
          <button type=button @click="${this.delete}" ?disabled="${this.selectedRowIndex<0}">Delete</button>
        </form>
        <sheet-editor .tableDesc=${this.tableDesc}
            .queryResults=${this.queryResults}
            @row-selected="${this.rowSelected}" @row-deleted="${this.rowDeleted}">
        </sheet-editor>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'table-sheet': TableSheet;
  }
}
