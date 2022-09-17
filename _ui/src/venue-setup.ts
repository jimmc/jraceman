import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import './jraceman-tabs.js'
import './table-queryedit.js'

/**
 * venue-setup is the tab content that contains other tabs for venue setup.
 *
 * @slot - Slots for contains tabs and tab content
 */
@customElement('venue-setup')
export class VenueSetup extends LitElement {
  static styles = css`
  `;

  render() {
    return html`
        <jraceman-tabs>
            <span slot="tab">Sites</span>
            <section slot="panel"><table-queryedit tableName="site"></table-queryedit></section>
            <span slot="tab">Areas</span>
            <section slot="panel"><table-queryedit tableName="area"></table-queryedit></section>
            <span slot="tab">Lane Order</span>
            <section slot="panel"><table-queryedit tableName="laneorder"></table-queryedit></section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'venue-setup': VenueSetup;
  }
}
