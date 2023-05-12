import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';
import {when} from 'lit/directives/when.js';

import { TableDesc,ColumnDesc, QueryResultsData, QueryResultsEvent, RequestEditEvent } from './table-desc.js'

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

  constructor() {
    super()
    //document.addEventListener("jraceman-sheet-editor-event", this.onQueryResultsEvent.bind(this))
  }

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
    return (col.Name=='id')
  }

  onClick(e: PointerEvent) {
    console.log("SheetEditor.onClick",e)
    const td = e.target as HTMLElement
    if (!td) {
      console.log("no target")
      return
    }
    const rowIndexStr = td.getAttribute('rowIndex')
    let rowIndex = -1
    if (rowIndexStr) {
      rowIndex = parseInt(rowIndexStr)
    } else {
      console.log("no rowIndex in event")
    }
    this.selectRowByIndex(rowIndex)
  }

  selectRowByIndex(rowIndex: number) {
    console.log("SheetEditor.selectRowByIndex", rowIndex)
    this.selectedRowIndex = rowIndex
    this.requestUpdate()

    // Let our parent component know about the selected row.
    this.dispatchEvent(new CustomEvent("row-selected", {
      bubbles: true,
      detail: rowIndex,
    }));
  }

  editSelectedRow() {
    const rowIndex = this.selectedRowIndex
    if (rowIndex<0) {
      console.log("No row selected")
      return
    }
    // Get the ID for the selected row
    const row = this.queryResults.Rows[rowIndex]
    const idColumnIndex = this.queryResults.Columns.findIndex(col => col.Name == "id")
    if (idColumnIndex < 0) {
      console.warn("SheetEditor.selectRowByIndex: no id field found in row", rowIndex)
      return
    }
    const rowId = row[idColumnIndex]

    // send request-edit event
    const event = new CustomEvent<RequestEditEvent>('jraceman-request-edit-event', {
      detail: {
        Table: this.queryResults.Table,
        ID: rowId
      } as RequestEditEvent
    });
    // Dispatch the event to the document so any element can listen for it.
    console.log("SheetEditor dispatching event", event)
    document.dispatchEvent(event);
  }

  deleteSelectedRow() {
    console.log("deleteSelectedRow NYI")
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
          <tr selected=${this.isRowIndexSelected(rowIndex)}>
          ${/*@ts-ignore*/
            repeat(this.tableDesc.Columns, (col:ColumnDesc, colIndex) => html`
            <td rowIndex=${rowIndex} selected=${this.isRowIndexSelected(rowIndex)}>
              ${when(this.isReadOnly(col),()=>html`
                ${row[colIndex]}
              `, ()=>html`
                ${when(col.FKTable, ()=>html`
                  <select id="val_${col.Name}_${row[0]}">
                    ${repeat(col.FKItems, (keyitem) => html`
                      <option value="${keyitem.ID}" ?selected=${row[colIndex]==keyitem.ID}>${keyitem.Summary}</option>
                    `)}
                  </select>
                `, ()=>html`
                  <input type=text value="${row[colIndex]}"></input>
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
