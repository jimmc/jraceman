import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';
import {when} from 'lit/directives/when.js';

import { ApiManager, XhrOptions } from './api-manager.js'
import { JracemanDialog } from './jraceman-dialog.js'
import { TableDesc, TableDescSupport, ColumnDesc, QueryResultsData, QueryResultsEvent, RequestEditEvent } from './table-desc.js'
import { TableEdit } from './table-edit.js'

/**
 * sheet-editor provides a table-layout for editing.
 */
@customElement('sheet-editor')
export class SheetEditor extends LitElement {
  static styles = css`
    table.sheet-editor {
      border: 2;
    }
    table.sheet-editor th {
      background-color: lightgray;
    }
    table.sheet-editor th.readonly {
      font-weight: normal;
    }
    tr[selected="true"] {
      background-color: lightblue;
    }
    td[editing="true"] input {
      background-color: lightgreen;
    }
    td[error="true"] input {
      background-color: lightpink;
    }
  `;

  @property({type: Object})
  tableDesc: TableDesc = {
    Table: "(unset-in-sheet-editor)",
    Columns:[],
  };

  @property({type: Object})
  queryResults: QueryResultsData = {
    Table: "",
    Columns: [],
    Rows: [],
  /* Sample data looks something like this:
    Table: "fake",
    Columns: [
      {
        Name: "col1",
        Type: "string"
      },
      {
        Name: "col2",
        Type: "int"
      },
    ],
    Rows: [ "aaa", 123 ],
    */
  };

  // The index of the selected row, or -1 if no row is selected.
  @property()
  selectedRowIndex = -1

  onQueryResultsEvent(e:Event) {
    const evt = e as CustomEvent<QueryResultsEvent>
    console.log("SheetEditor got updated results", evt.detail.results)
    this.queryResults = evt.detail.results
    this.selectedRowIndex = -1
  }

  isRowIndexSelected(rowIndex: number) {
    return rowIndex == this.selectedRowIndex
  }

  isReadOnly(col:ColumnDesc) {
    return (col.Name=='id' || col.ReadOnly)
  }

  onClick(e: PointerEvent) {
    console.log("SheetEditor.onClick",e)
    const td = e.target as HTMLElement
    if (!td) {
      console.log("no target")
      return
    }
    const tr = td.closest('tr')!
    const rowIndexStr = tr.getAttribute('rowIndex')
    let rowIndex = -1
    if (rowIndexStr) {
      rowIndex = parseInt(rowIndexStr)
    } else {
      console.log("no rowIndex in event")
    }
    this.selectRowByIndex(rowIndex)
  }

  async onChange(e: InputEvent) {
    console.log("onChange", e)
    const target = e.target as HTMLInputElement
    const value = target.value
    console.log("value is", value)
    const td = target.closest('td')!    // Get our containing table cell
    const colIndexStr = td.getAttribute('colIndex')
    const colIndex = colIndexStr ? parseInt(colIndexStr) : -1
    const tr = td.closest('tr')!        // Get our containing row
    const rowIndexStr = tr.getAttribute('rowIndex')
    const rowIndex = rowIndexStr ? parseInt(rowIndexStr) : -1
    if (rowIndex<0) {
      console.log("No row selected")
      return
    }
    const id = this.idForRowIndex(rowIndex)
    console.log(" ..for rowId", id, " column", this.tableDesc.Columns[colIndex].Name)
    const err = await this.saveChange(id, colIndex, value)
    if (err) {
      td.removeAttribute("editing")
      td.setAttribute("error","true")
    } else {
      td.removeAttribute("editing")
    }
  }

  // onInput gets called each time the user edits the text in one of our text fields.
  onInputText(e: InputEvent) {
    console.log("onInputText", e)
    const target = e.target as HTMLInputElement
    const td = target.closest('td')!    // Get our containing table cell
    td.setAttribute("editing","true")
      // Set the attribute so we can visually notify the user that the field
      // is being edited. It won't be saved until the user presses Enter or
      // exits the field (such as by tabbing out).
  }

  selectRowByIndex(rowIndex: number) {
    console.log("SheetEditor.selectRowByIndex", rowIndex)
    this.selectedRowIndex = rowIndex
    this.requestUpdate()

    // Let our parent component know about the selected row.
    this.dispatchEvent(new CustomEvent("row-selected", {
      bubbles: true,
      detail: rowIndex,
    }))
  }

  idForRowIndex(rowIndex: number) {
    // Get the ID for the selected row
    const row = this.queryResults.Rows[rowIndex]
    const idColumnIndex = this.queryResults.Columns.findIndex(col => col.Name == "id")
    if (idColumnIndex < 0) {
      console.warn("SheetEditor.selectRowByIndex: no id field found in row", rowIndex)
      return ""
    }
    const rowId = row[idColumnIndex]
    if (!rowId) {
      console.log("No rowId found in row")
    }
    return rowId
  }

  editSelectedRow() {
    const rowIndex = this.selectedRowIndex
    if (rowIndex<0) {
      console.log("No row selected")
      return
    }
    const rowId = this.idForRowIndex(rowIndex)
    if (!rowId) {
      return    // Issue already logged to console
    }

    // send request-edit event
    const event = new CustomEvent<RequestEditEvent>('jraceman-request-edit-event', {
      detail: {
        Table: this.queryResults.Table,
        ID: rowId
      } as RequestEditEvent
    })
    // Dispatch the event to the document so any element can listen for it.
    console.log("SheetEditor dispatching event", event)
    document.dispatchEvent(event)
  }

  async deleteSelectedRow() {
    if (this.selectedRowIndex < 0) {
        JracemanDialog.messageDialog("Error", "No row is selected", ["Dismiss"])
        return
    }
    const rowIndex = this.selectedRowIndex
    const id = this.queryResults.Rows[rowIndex][0]
    const ok = await TableEdit.deleteRow(this.tableDesc.Table, id)
    if (!ok) {
      return
    }

    this.queryResults.Rows.splice(rowIndex, 1)  // Update our table of rows.
    this.selectRowByIndex(-1)       // No row is selected now.
    // Let our parent component know that we changed the data.
    this.dispatchEvent(new CustomEvent("row-deleted", {
      bubbles: true,
      detail: rowIndex,
    }))
  }

  async saveChange(id: string, colIndex: number, colVal: string) {
    console.log("in SheetEditor.saveChange")
    const col = this.tableDesc.Columns[colIndex]
    const name = col.Label      // TODO - Is this the right way to get the field name?
    console.log("Type of field " + name + " is " + col.Type)
    // For non-string fields, convert from the string in the form
    // to the appropriate data type for the field.
    const convertedColVal = TableDescSupport.convertToType(colVal, col.Type)
    const updatePath = '/api/crud/' + this.tableDesc.Table + '/' + id
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
        result.Table = this.tableDesc.Table
      }
      console.log(result)
      return null       // No error
    } catch(e) {
      console.error("Error:", e/*.responseText*/)
      return e
    }
  }

  // TODO: We assume below that the ID column is 0 and its value is in row[0].
  render() {
    return html`
      Table: ${this.tableDesc.Table}<br/>
      ${this.queryResults.Error}
      <table class=sheet-editor @click="${this.onClick}">
        <tr>
          ${/*@ts-ignore*/
            repeat(this.tableDesc.Columns, (col:ColumnDesc/*, colIndex*/) => html`
            ${when(this.isReadOnly(col),()=>html`
              <th class=readonly>${col.Label}</th>
            `,()=>html`
              <th>${col.Label}</th>
            `)}
          `)}
        </tr>
        ${/*@ts-ignore*/
          repeat(this.queryResults.Rows, (row:any[], rowIndex) => html`
          <tr rowIndex=${rowIndex} selected=${this.isRowIndexSelected(rowIndex)}>
          ${/*@ts-ignore*/
            repeat(this.tableDesc.Columns, (col:ColumnDesc, colIndex) => html`
            <td colIndex=${colIndex} selected=${this.isRowIndexSelected(rowIndex)}>
              ${when(this.isReadOnly(col),()=>html`
                ${row[colIndex]}
              `, ()=>html`
                ${when(col.FKTable, ()=>html`
                  <select id="val_${col.Name}_${row[0]}" @change="${this.onChange}">
                    ${repeat(col.FKItems, (keyitem) => html`
                      <option value="${keyitem.ID}"
                          ?selected=${row[colIndex]==keyitem.ID}>
                        ${keyitem.Summary}
                      </option>
                    `)}
                  </select>
                `, ()=>html`
                  <input @input="${this.onInputText}" @change="${this.onChange}"
                      type=text value="${row[colIndex]}">
                  </input>
                `)}
              `)}
            </td>
          `)}
          </tr>
        `)}
      </table>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'sheet-editor': SheetEditor;
  }
}
