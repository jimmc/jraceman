import {LitElement, html, css, render} from 'lit';
import {customElement} from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';

import "./jraceman-dropdown.js"
// import { PostError } from "./message-log.js"
import { QueryResultsData, QueryResultsEvent } from "./table-desc.js"

// A drop-down menu for query operations.
@customElement('query-menu')
export class QueryMenu extends LitElement {
  static styles = css`
    jraceman-dropdown {
      display: inline-block;    /* Make our menu on same line as the tab label */
    }

    .menu {
      cursor: context-menu;
    }
  `;

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

  connectedCallback() {
    super.connectedCallback()
    document.addEventListener("jraceman-query-results-event", this.onQueryResultsEvent.bind(this))
  }

  onQueryResultsEvent(e:Event) {
    const evt = e as CustomEvent<QueryResultsEvent>
    console.log("QueryMenu got updated query results", evt.detail.results)
    // Save the query results so we can write it out on request.
    this.queryResults = evt.detail.results as QueryResultsData
  }

  onViewInNewTab() {
    const newWindow = window.open()!;
    newWindow.document.title = "JRaceman query results"
    render(this.renderQueryResults(), newWindow.document.body)
  }

  onDumpSqlInNewTab() {
    const newWindow = window.open()!;
    newWindow.document.title = "JRaceman query results as SQL"
    newWindow.document.body.innerText = this.renderSql()
  }

  renderSql() {
    let sql = ""
    sql += "-- Column types:"
    const cols = this.queryResults.Columns
    for (var col of cols) {
      sql += " "+col.Name+":"+col.Type
    }
    sql += "\n"
    for (var row of this.queryResults.Rows) {
      let sqlLine = "INSERT INTO " + this.queryResults.Table + "("
      let colsep = ""
      for (var col of cols) {
        sqlLine += colsep
        colsep = ","
        sqlLine += col.Name
      }
      sqlLine += ") VALUES("
      colsep = ""
      for (let c = 0; c < this.queryResults.Columns.length; c++) {
        sqlLine += colsep
        colsep = ","
        if (cols[c].Type=="string") {
          sqlLine += this.sqlQuoteAndEscape(row[c])     // Quote and escape the string.
        } else {
          // No quoting needed
          sqlLine += row[c]
        }
      }
      sqlLine += ");\n"
      sql += sqlLine
    }
    return sql
  }

  // This method should be identical to QueryResults.render(),
  // except without any of the selection or interactive parts.
  renderQueryResults() {
    return html`
      Table: ${this.queryResults.Table}<br/>
      ${this.queryResults.Error}
      <table>
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

  sqlQuoteAndEscape(s: string) {
    return "'" + this.sqlEscape(s) + "'"
  }

  sqlEscape(s: string) {
    s = s.replace(/[\0\n\r\b\t\\'"\x1a]/g, function (ch) {
      switch (ch) {
        case "\0":
          return "\\0";
        case "\n":
          return "\\n";
        case "\r":
          return "\\r";
        case "\b":
          return "\\b";
        case "\t":
          return "\\t";
        case "\x1a":
          return "\\Z";
        case "'":
          return "''";
        case '"':
          return '""';
        default:
          return "\\" + ch;
      }
    })
    return s
  }

  render() {
    return html`
      <jraceman-dropdown>
        <span slot="control" class="menu">☰</span>
        <div slot="content">
          <button @click="${this.onViewInNewTab}">View in new tab</button>
          <button @click="${this.onDumpSqlInNewTab}">Dump SQL in new tab</button>
        </div>
      </jraceman-dropdown>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'query-menu': QueryMenu;
  }
}
