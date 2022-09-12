import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import { PostMessage } from './message-log.js'

import "./sql-query.js"

// A pane for debugging.
@customElement('debug-pane')
export class DebugPane extends LitElement {
  static styles = css`
  `;

  render() {
    return html`
      <jraceman-tabs>
          <span slot="tab">Messages</span>
          <section slot="panel">
            <button @click=${this.onClick.bind(this,"info","Sample info message")}>Post Info message</button>
            <button @click=${this.onClick.bind(this,"warning","A warning message")}>Post Warning message</button>
            <button @click=${this.onClick.bind(this,"error","Error message")}>Post Error message</button>
          </section>
          <span slot="tab">SQL Query</span>
          <section slot="panel"><sql-query></sql-query></section>
      </jraceman-tabs>
    `;
  }

  onClick(level: string, text: string) {
    PostMessage("debug", level, text)
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'debug-pane': DebugPane;
  }
}
