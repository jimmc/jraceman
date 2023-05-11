import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
//import {repeat} from 'lit/directives/repeat.js';
//import {when} from 'lit/directives/when.js';

import { ApiManager, XhrOptions } from './api-manager.js'
import { PostError } from './message-log.js'
import { TableDesc /*, QueryResultsEvent*/ } from './table-desc.js'

import './sheet-editor.js'

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
    Columns:[],
  };

  @property({type: Object /*, notify: true*/})
  queryResults: object = {
    Table: "(unset-in-query-results)",
    Columns: [],
    Rows: [],
  };

  @property({type: String})
  selectedOp: string = '';

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

  async search() {
    console.log("TableSheet.search begin");
    let params:any[] = [];
    // TODO - if we have a selection field, add it here (see table-query).
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
      this.queryResults = result;
    } catch(e) {
      const evt = e as XMLHttpRequest
      PostError("query", evt.responseText)
      console.log("Error in table query", e)
      return
    }
    console.log("TableSheet.search results", this.queryResults);

    // TODO - display the results in our editing table.
    /*
    const event = new CustomEvent<QueryResultsEvent>('jraceman-query-results-event', {
      detail: {
        message: 'Query results for table '+this.tableDesc.Table,
        results: this.queryResults
      } as QueryResultsEvent
    });
    // Dispatch the event to the document so any element can listen for it.
    console.log("TableSheet dispatching event", event)
    document.dispatchEvent(event);
    */
  }

  async add() {
    console.log("TableSheet.add NYI");
  }

  async edit() {
    console.log("TableSheet.edit NYI");
  }

  async delete() {
    console.log("TableSheet.delete NYI");
  }

  isStringColumn(colType: string) {
    return colType == "string";
  }

  render() {
    return html`
        <form>
          <button type=button @click="${this.search}">Search</button>
          <button type=button @click="${this.add}">Add</button>
          <button type=button @click="${this.edit}">Edit</button>
          <button type=button @click="${this.delete}">Delete</button>
        </form>
        <sheet-editor tableDesc=${JSON.stringify(this.tableDesc)}
            queryResults=${JSON.stringify(this.queryResults)}>
        </sheet-editor>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'table-sheet': TableSheet;
  }
}
