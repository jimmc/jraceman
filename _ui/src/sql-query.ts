import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';

import { ApiManager } from "./api-manager.js"
import { PostError } from "./message-log.js"
import { QueryResultsEvent } from "./table-desc.js"

// Allows the user to enter and execute an SQL query.
@customElement('sql-query')
export class SqlQuery extends LitElement {
  static styles = css`
    :host {
      height: 100%;
    }
    #main {
      display: flex;
      flex-direction: column;
      height: 100%;
    }
    #sqlText {
      flex: auto;
    }
    #buttons {
      flex: none;
    }
  `;

  checkEnter(e: any) {
    if (e.key == 'Enter' && e.shiftKey) {
      e.stopPropagation();
      this.execute();
    }
  }

  // Clears the SQL text area.
  clear() {
    (this.shadowRoot!.querySelector("#sqlText")! as HTMLSelectElement).value = ""
  }

  // Executes the SQL text.
  async execute() {
    const sql = (this.shadowRoot!.querySelector("#sqlText")! as HTMLSelectElement).value
    console.log("Execute:", sql);     // TODO
    const path = '/api/debug/sql/';
    const formData = {
      q: sql
    };
    const options = {
      method: 'POST',
      params: formData
    };
    try {
      const result = await ApiManager.xhrJson(path, options)
      // Now tell results tab to display this data.
      const event = new CustomEvent<QueryResultsEvent>('jraceman-query-results-event', {
        detail: {
          message: 'Query results for SQL query '+ sql,
          results: result
        } as QueryResultsEvent
      });
      // Dispatch the event to the document so any element can listen for it.
      console.log("SqlQuery dispatching event", event)
      document.dispatchEvent(event);
    } catch(e) {
      const evt = e as XMLHttpRequest
      PostError("sql", evt.responseText)
      console.error("Error",e)
    }
  }

  render() {
    return html`
      <div id="main">
        <textarea id="sqlText" cols="60" autofocus="true"
            placeholder="Enter SQL text" @keypress="${this.checkEnter}"></textarea>
        <div id="buttons">
          <button @click="${this.execute}">Execute</button>
          <button @click="${this.clear}">Clear</button>
        </div>
      </div>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'sql-query': SqlQuery;
  }
}
