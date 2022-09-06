import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';
import '@vaadin/grid/vaadin-grid.js';
import '@vaadin/grid/vaadin-grid-column.js';

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

  // vaadin-grid needs a path within each row item passed to it, so we
  // convert each row to a map using the column names as field names.
  queryResultsAsMaps(queryResults: QueryResultsData) {
    let res = []
    for (let row of queryResults.Rows) {
      let rr = {}
      let cx = 0
      for (let col of queryResults.Columns) {
        const cc = col as ColumnDesc
        // With javascript, we could use rr[cc.Name]=row[cx]
        // and access that using a field op such as rr.F where F=cc.Name,
        // but typescript won't let us do that, so we resort to using eval instead.
        const ev = "rr." + cc.Name + "=JSON.parse('" + JSON.stringify(row[cx]) + "')"
        eval(ev)        // rr.F = row[cx]   where F=cc.Name
        cx++
      }
      res.push(rr)
    }
    return res
  }

  render() {
    return html`
      Table: ${this.queryResults.Table}<br/>
      ${this.queryResults.Error}
        <vaadin-grid id="grid" items="${JSON.stringify(this.queryResultsAsMaps(this.queryResults))}">
          ${/*@ts-ignore*/
            repeat(this.queryResults.Columns, (col:ColumnDesc, colIndex) => html`
            <vaadin-grid-column header="${col.Name}" path="${col.Name}">
            </vaadin-grid-column>
          `)}
        </vaadin-grid>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'query-results': QueryResults;
  }
}
