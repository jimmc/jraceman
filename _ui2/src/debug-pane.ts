import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import { PostMessage } from './message-log.js'

// A pane for debugging.
@customElement('debug-pane')
export class DebugPane extends LitElement {
  static styles = css`
  `;

  render() {
    return html`
      <button @click=${this.onClick.bind(this,"info","Sample info message")}>Post Info message</button>
      <button @click=${this.onClick.bind(this,"warning","A warning message")}>Post Warning message</button>
      <button @click=${this.onClick.bind(this,"error","Error message")}>Post Error message</button>
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
