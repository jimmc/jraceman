import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';

import { ColumnDesc, QueryResultsData, QueryResultsEvent } from './table-desc.js'

/*
interface SelectedResult {
  Table: string;
  ID: string;
}
*/

/**
 * query-results provides a form to do a query on a table.
 */
@customElement('query-results')
export class QueryResults extends LitElement {
  static styles = css`
  `;

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
        Type: "string"
      },
    ],
    Rows: [
      [ "aaa", 123 ],
      [ "bbb", 456 ],
      [ "ccc", 789 ]
    ],
    */
  };

  constructor() {
    super()
    document.addEventListener("jracemanQueryResultsEvent", this.handleQueryResultsEvent.bind(this))
  }

  handleQueryResultsEvent(e:Event) {
    const evt = e as CustomEvent<QueryResultsEvent>
    console.log("QueryResults got updated results", evt.detail.results)
    this.queryResults = evt.detail.results
  }

  onClick(e: PointerEvent) {
    console.log("QueryResult.onClick",e)
    const td = eval("e.path[0]")
    const rowIndexStr = td.getAttribute('rowIndex')
    if (!rowIndexStr) {
      console.log("no rowIndex in event")
      return
    }
    const rowIndex = parseInt(rowIndexStr)
    this.selectRowByIndex(rowIndex)
  }

  selectRowByIndex(rowIndex: number) {
    console.log("QueryResults.selectRowByIndex", rowIndex)      // TODO - implement this
  }

  render() {
    return html`
      Table: ${this.queryResults.Table}<br/>
      ${this.queryResults.Error}
      <table @click="${this.onClick}">
        <tr>
          ${/*@ts-ignore*/
            repeat(this.queryResults.Columns, (col:ColumnDesc/*, colIndex*/) => html`
            <th>
              ${col.Name}
            </th>
          `)}
        </tr>
        ${/*@ts-ignore*/
          repeat(this.queryResults.Rows, (row:any[], rowIndex) => html`
          <tr>
          ${/*@ts-ignore*/
            repeat(this.queryResults.Columns, (col:ColumnDesc, colIndex) => html`
            <td rowIndex=${rowIndex}>
              ${row[colIndex]}
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
    'query-results': QueryResults;
  }
}
