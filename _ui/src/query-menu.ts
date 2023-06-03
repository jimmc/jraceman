import { LitElement, html, css, render } from 'lit'
import { customElement } from 'lit/decorators.js'

import './jraceman-dropdown.js'

import { QueryResults } from './query-results.js'

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

  queryResultsElement?: QueryResults

  setQueryResults(qr: QueryResults) {
    this.queryResultsElement = qr
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

  // Render the query results the same way as we do in the query-results tab.
  renderQueryResults() {
    return html`
      <style>
        ${QueryResults.styles}
      </style>
      ${this.queryResultsElement!.render()}
    `;
  }

  renderSql() {
    const qr = this.queryResultsElement!.getQueryResults()
    let sql = ""
    sql += "-- Column types:"
    const cols = qr.Columns
    for (var col of cols) {
      sql += " "+col.Name+":"+col.Type
    }
    sql += "\n"
    for (var row of qr.Rows) {
      let sqlLine = "INSERT INTO " + qr.Table + "("
      let colsep = ""
      for (var col of cols) {
        sqlLine += colsep
        colsep = ","
        sqlLine += col.Name
      }
      sqlLine += ") VALUES("
      colsep = ""
      for (let c = 0; c < qr.Columns.length; c++) {
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
