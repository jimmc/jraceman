import { LitElement, html, css } from 'lit'
import { PropertyValues } from 'lit-element'
import { customElement, property } from 'lit/decorators.js'

import './query-fields.js'

import { ApiManager, XhrOptions } from './api-manager.js'
import { PostError } from './message-log.js'
import { QueryFields } from './query-fields.js'
import { TableDesc, QueryResultsData, QueryResultsEvent } from './table-desc.js'

/**
 * table-query provides a form to do a query on a table.
 */
@customElement('table-query')
export class TableQuery extends LitElement {
  static styles = css`
  `;

  @property({type: Object})
  tableDesc: TableDesc = {
    Table: "(unset-in-table-query)",
    Columns:[],
  };

  @property({type: Object /*, notify: true*/})
  queryResults: object = {};

  @property({type: String})
  selectedOp: string = '';

  queryFields?: QueryFields

  firstUpdated(changedProperties:PropertyValues<this>) {
    super.firstUpdated(changedProperties);
    this.queryFields = this.shadowRoot!.querySelector("query-fields")! as QueryFields
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

  clear() {
    this.queryFields!.clear()
  }

  async search() {
    console.log("TableQuery.search begin");
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
      this.queryResults = result;
    } catch(e) {
      const evt = e as XMLHttpRequest
      PostError("query", evt.responseText)
      console.log("Error in table query", e)
      return    // Don't attempt to update the QueryResults tab.
    }
    console.log("TableQuery.search results", this.queryResults);
    // Now tell results tab to display this data.
    const event = new CustomEvent<QueryResultsEvent>('jraceman-query-results-event', {
      detail: {
        message: 'Query results for table '+this.tableDesc.Table,
        results: this.queryResults
      } as QueryResultsEvent
    });
    // Dispatch the event to the document so any element can listen for it.
    console.log("TableQuery dispatching event", event)
    document.dispatchEvent(event);
  }

  async editInSheet() {
    console.log("TableQuery.editInSheet begin");
    const params = this.queryFields!.fieldsAsParams()
    const options: XhrOptions = {
      method: "POST",
      params: params,
    }
    const queryPath = '/api/query/' + this.tableDesc.Table + '/'
    let result: QueryResultsData
    try {
      result = await ApiManager.xhrJson(queryPath, options) as QueryResultsData
      if (result && !result.Table) {
        result.Table = this.tableDesc.Table
      }
    } catch(e) {
      const evt = e as XMLHttpRequest
      PostError("query", evt.responseText)
      console.log("Error in table query", e)
      return    // Don't attempt to update the QueryResults tab.
    }
    console.log("TableQuery.editInSheet results", result)
    // Now tell results tab to display this data.
    const event = new CustomEvent<QueryResultsEvent>('jraceman-query-to-sheet-event', {
      detail: {
        message: 'Query results for table '+this.tableDesc.Table,
        results: result
      } as QueryResultsEvent
    })
    console.log("TableQuery dispatching event", event)
    this.dispatchEvent(event)
  }

  isStringColumn(colType: string) {
    return colType == "string"
  }

  render() {
    return html`
        <form>
          <button type=button @click="${this.search}">Search</button>
          <button type=button @click="${this.clear}">Clear</button>
          <button type=button @click="${this.editInSheet}">Edit in Sheet</button>
          <query-fields .tableDesc=${this.tableDesc}></query-fields>
        </form>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'table-query': TableQuery;
  }
}
