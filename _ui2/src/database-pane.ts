import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';

import "./jraceman-tabs.js"
import "./table-queryedit.js"

// A pane for database stuff
@customElement('database-pane')
export class DatabasePane extends LitElement {
  static styles = css`
  `;

  render() {
    return html`
      <jraceman-tabs>
        <span slot="tab">Context Options</span>
        <section slot="panel"><table-queryedit tableName="contextoption"></table-queryedit></section>
      </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'database-pane': DatabasePane;
  }
}
