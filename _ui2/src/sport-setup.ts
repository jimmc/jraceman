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
            <section slot="panel">Content for Genders</section>
            <span slot="tab">Progressions</span>
            <section slot="panel">Content for Progressions</section>
            <span slot="tab">Scoring System</span>
            <section slot="panel">Content for Scoring System</section>
            <span slot="tab">Scoring Rules</span>
            <section slot="panel">Content for Scoring Rules</section>
            <span slot="tab">Stages</span>
            <section slot="panel">Content for Stages</section>
            <span slot="tab">Exceptions</span>
            <section slot="panel">Content for Exceptions</section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'sport-setup': SportSetup;
  }
}
