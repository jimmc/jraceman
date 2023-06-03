import { LitElement, html, css } from 'lit'
import { customElement } from 'lit/decorators.js'

import './jraceman-tabs.js'
import './table-manager.js'

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
            <section slot="panel"><table-manager tableName="challenge"></table-manager></section>
            <span slot="tab">Teams</span>
            <section slot="panel"><table-manager tableName="team"></table-manager></section>
            <span slot="tab">People</span>
            <section slot="panel"><table-manager tableName="person"></table-manager></section>
            <span slot="tab">Seeding Plans</span>
            <section slot="panel"><table-manager tableName="seedingplan"></table-manager></section>
            <span slot="tab">Seeding Lists</span>
            <section slot="panel"><table-manager tableName="seedinglist"></table-manager></section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'team-setup': TeamSetup;
  }
}
