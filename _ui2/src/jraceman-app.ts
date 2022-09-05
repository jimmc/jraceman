import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import '@vaadin/vaadin-split-layout/vaadin-split-layout.js';
import './jraceman-tabs.js'
import './query-results.js'
import './sport-setup.js'

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
    }
    .title-bar {
      width: "100%";
      background-color: lightgray;
      color: black;
    }
  `;

  render() {
    return html`
      <div class="title-bar">JRaceman</div>
      <vaadin-split-layout id="main" orientation="vertical" vertical>
          <div id="top" class="tab-container">
            <jraceman-tabs>
                <h3 slot="tab">Sport Setup</h2>
                <section slot="panel"><sport-setup></sport-setup></section>
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
            <jraceman-tabs>
                <h3 slot="tab">Messages</h2>
                <section slot="panel">Messages content</section>
                <h3 slot="tab">Query Results</h2>
                <section slot="panel"><query-results></query-results></section>
                <h3 slot="tab">Report Results</h2>
                <section slot="panel">Report Results content</section>
                <h3 slot="tab">Help</h2>
                <section slot="panel">Help is not yet implemented</section>
            </jraceman-tabs>
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
