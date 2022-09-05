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
  /* Sample data looks something like this:*/
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
    ]
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

  render() {
    console.log("QueryResults render")
    return html`
      Query Results here
      ${this.queryResults.Error}
        <vaadin-grid id="grid" items="${JSON.stringify(this.queryResults.Rows)}">
          ${/*@ts-ignore*/
            repeat(this.queryResults.Columns, (col:ColumnDesc, colIndex) => html`
            <vaadin-grid-column header="${col.Name}" path="[colIndex]">
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
