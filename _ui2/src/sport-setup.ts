import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import './jraceman-tabs.js'

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
            <h3 slot="tab">Competitions</h2>
            <section slot="panel">Content for Competitions</section>
            <h3 slot="tab">Levels</h2>
            <section slot="panel">Content for Levels</section>
            <h3 slot="tab">Genders</h2>
            <section slot="panel">Content for Genders</section>
            <h3 slot="tab">Progressions</h2>
            <section slot="panel">Content for Progressions</section>
            <h3 slot="tab">Scoring System</h2>
            <section slot="panel">Content for Scoring System</section>
            <h3 slot="tab">Scoring Rules</h2>
            <section slot="panel">Content for Scoring Rules</section>
            <h3 slot="tab">Stages</h2>
            <section slot="panel">Content for Stages</section>
            <h3 slot="tab">Exceptions</h2>
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
