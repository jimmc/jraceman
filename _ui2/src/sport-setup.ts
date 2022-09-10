import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import './jraceman-tabs.js'
import './table-queryedit.js'

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
            <section slot="panel"><table-queryedit tableName="competition"></table-queryedit></section>
            <span slot="tab">Levels</span>
            <section slot="panel"><table-queryedit tableName="level"></table-queryedit></section>
            <span slot="tab">Genders</span>
            <section slot="panel"><table-queryedit tableName="gender"></table-queryedit></section>
            <span slot="tab">Progressions</span>
            <section slot="panel"><table-queryedit tableName="progression"></table-queryedit></section>
            <span slot="tab">Scoring Systems</span>
            <section slot="panel"><table-queryedit tableName="scoringsystem"></table-queryedit></section>
            <span slot="tab">Scoring Rules</span>
            <section slot="panel"><table-queryedit tableName="scoringrule"></table-queryedit></section>
            <span slot="tab">Stages</span>
            <section slot="panel"><table-queryedit tableName="stage"></table-queryedit></section>
            <span slot="tab">Exceptions</span>
            <section slot="panel"><table-queryedit tableName="exception"></table-queryedit></section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'sport-setup': SportSetup;
  }
}
