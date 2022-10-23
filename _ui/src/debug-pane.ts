import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import { PostMessageParts } from './message-log.js'

import "./sql-query.js"

// A pane for debugging.
@customElement('debug-pane')
export class DebugPane extends LitElement {
  static styles = css`
  `;

  render() {
    return html`
      <jraceman-tabs>
        <span slot="tab">SQL Query</span>
        <section slot="panel"><sql-query></sql-query></section>
        <span slot="tab">Edit ID</span>
        <section slot="panel">Edit ID is not yet implemented</section>
        <span slot="tab">Options</span>
        <section slot="panel"><table-queryedit tableName="option"></table-queryedit></section>
        <span slot="tab">Races</span>
        <section slot="panel"><table-queryedit tableName="race"></table-queryedit></section>
        <span slot="tab">Lanes</span>
        <section slot="panel"><table-queryedit tableName="lane"></table-queryedit></section>
        <span slot="tab">Messages</span>
        <section slot="panel">
          <button @click=${this.onClick.bind(this,"info","Sample info message")}>Post Info message</button>
          <button @click=${this.onClick.bind(this,"warning","A warning message")}>Post Warning message</button>
          <button @click=${this.onClick.bind(this,"error","Error message")}>Post Error message</button>
        </section>
      </jraceman-tabs>
    `;
  }

  onClick(level: string, text: string) {
    PostMessageParts("debug", level, text)      // "debug" is our source name.
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'debug-pane': DebugPane;
  }
}
