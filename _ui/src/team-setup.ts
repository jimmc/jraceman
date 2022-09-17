import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import './jraceman-tabs.js'
import './table-queryedit.js'

/**
 * team-setup is the tab content that contains other tabs for team setup.
 *
 * @slot - Slots for contains tabs and tab content
 */
@customElement('team-setup')
export class TeamSetup extends LitElement {
  static styles = css`
  `;

  render() {
    return html`
        <jraceman-tabs>
            <span slot="tab">Challenges</span>
            <section slot="panel"><table-queryedit tableName="challenge"></table-queryedit></section>
            <span slot="tab">Teams</span>
            <section slot="panel"><table-queryedit tableName="team"></table-queryedit></section>
            <span slot="tab">People</span>
            <section slot="panel"><table-queryedit tableName="person"></table-queryedit></section>
            <span slot="tab">Seeding Plans</span>
            <section slot="panel"><table-queryedit tableName="seedingplan"></table-queryedit></section>
            <span slot="tab">Seeding Lists</span>
            <section slot="panel"><table-queryedit tableName="seedinglist"></table-queryedit></section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'team-setup': TeamSetup;
  }
}
