import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import './jraceman-tabs.js'
import './table-manager.js'

/**
 * sport-setup is the tab content that contains other tabs for sport setup.
 *
 * @slot - Slots for contains tabs and tab content
 */
@customElement('sport-setup')
export class SportSetup extends LitElement {
  static styles = css`
  `;

  render() {
    return html`
        <jraceman-tabs>
            <span slot="tab">Competitions</span>
            <section slot="panel"><table-manager tableName="competition"></table-manager></section>
            <span slot="tab">Levels</span>
            <section slot="panel"><table-manager tableName="level"></table-manager></section>
            <span slot="tab">Genders</span>
            <section slot="panel"><table-manager tableName="gender"></table-manager></section>
            <span slot="tab">Progressions</span>
            <section slot="panel"><table-manager tableName="progression"></table-manager></section>
            <span slot="tab">Scoring Systems</span>
            <section slot="panel"><table-manager tableName="scoringsystem"></table-manager></section>
            <span slot="tab">Scoring Rules</span>
            <section slot="panel"><table-manager tableName="scoringrule"></table-manager></section>
            <span slot="tab">Stages</span>
            <section slot="panel"><table-manager tableName="stage"></table-manager></section>
            <span slot="tab">Exceptions</span>
            <section slot="panel"><table-manager tableName="exception"></table-manager></section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'sport-setup': SportSetup;
  }
}
