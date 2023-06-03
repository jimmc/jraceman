import { LitElement, html, css } from 'lit'
import { customElement } from 'lit/decorators.js'

import './jraceman-tabs.js'
import './table-manager.js'

/**
 * meet-setup is the tab content that contains other tabs for meet setup.
 *
 * @slot - Slots for contains tabs and tab content
 */
@customElement('meet-setup')
export class MeetSetup extends LitElement {
  static styles = css`
  `;

  render() {
    return html`
        <jraceman-tabs>
            <span slot="tab">Meets</span>
            <section slot="panel"><table-manager tableName="meet"></table-manager></section>
            <span slot="tab">Registration Fees</span>
            <section slot="panel"><table-manager tableName="registrationfee"></table-manager></section>
            <span slot="tab">Registrations</span>
            <section slot="panel"><table-manager tableName="registration"></table-manager></section>
            <span slot="tab">Events</span>
            <section slot="panel"><table-manager tableName="event"></table-manager></section>
            <span slot="tab">Entries</span>
            <section slot="panel"><table-manager tableName="entry"></table-manager></section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'meet-setup': MeetSetup;
  }
}
