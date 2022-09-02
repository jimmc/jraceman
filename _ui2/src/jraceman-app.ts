import {LitElement, html, customElement, css} from 'lit-element';
import '@vaadin/vaadin-split-layout/vaadin-split-layout.js';
import './jraceman-tabs.js'

/**
 * jraceman-app is the top-level component that contains the entire JRaceman application.
 *
 * @slot - This element has a slot
 * @csspart button - The button
 */
@customElement('jraceman-app')
export class JracemanApp extends LitElement {
  static styles = css`
    :host {
      display: block;
      border: solid 1px black;
      padding: 16px;
      max-width: 800px;
    }
  `;

  render() {
    return html`
      <h1>JRaceman</h1>
      <vaadin-split-layout id="main" orientation="vertical" vertical>
          <div id="top" class="tab-container">
            <jraceman-tabs>
                <h3 slot="tab">Tab 1</h2>
                <section slot="panel">Content for tab 1</section>
                <h3 slot="tab">Tab 2</h2>
                <section slot="panel">Content for tab 2</section>
                <h3 slot="tab">Tab 3</h2>
                <section slot="panel">
                    <jraceman-tabs>
                        <h3 slot="tab">Tab 3.1</h2>
                        <section slot="panel">Content for tab 3.1</section>
                        <h3 slot="tab">Tab 3.2</h2>
                        <section slot="panel">Content for tab 3.2</section>
                    </jraceman-tabs>
                </section>
            </jraceman-tabs>
          </div>
          <div id="bottom" class="tab-container">
            BOTTOM
          </div>
      </vaadin-split-layout>
      <slot></slot>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'jraceman-app': JracemanApp;
  }
}
