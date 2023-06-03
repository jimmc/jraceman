import {LitElement, html, css} from 'lit'
import {customElement} from 'lit/decorators.js'

import './jraceman-tabs.js'
import './table-manager.js'

// A pane for database stuff
@customElement('database-pane')
export class DatabasePane extends LitElement {
  static styles = css`
  `;

  render() {
    return html`
      <jraceman-tabs>
        <span slot="tab">Checks</span>
        <section slot="panel"><reports-pane apiName='check'></reports-pane></section>
        <span slot="tab">Context Options</span>
        <section slot="panel"><table-manager tableName="contextoption"></table-manager></section>
      </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'database-pane': DatabasePane;
  }
}
