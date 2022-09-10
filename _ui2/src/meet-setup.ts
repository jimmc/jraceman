import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import './jraceman-tabs.js'
import './table-queryedit.js'

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
            <section slot="panel"><table-queryedit tableName="meet"></table-queryedit></section>
            <span slot="tab">Registration Fees</span>
            <section slot="panel"><table-queryedit tableName="registrationfee"></table-queryedit></section>
            <span slot="tab">Registrations</span>
            <section slot="panel"><table-queryedit tableName="registration"></table-queryedit></section>
            <span slot="tab">Events</span>
            <section slot="panel"><table-queryedit tableName="event"></table-queryedit></section>
            <span slot="tab">Entries</span>
            <section slot="panel"><table-queryedit tableName="entry"></table-queryedit></section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'meet-setup': MeetSetup;
  }
}
