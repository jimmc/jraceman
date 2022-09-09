import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import '@vaadin/vaadin-split-layout/vaadin-split-layout.js';

import './query-results.js'
import './sport-setup.js'
import { JracemanTabs} from './jraceman-tabs.js'
import { QueryResultsEvent } from './table-desc.js'

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

  constructor() {
    super()
    // We add a listener for query results so that we can make
    // the query results tab visible.
    document.addEventListener("jraceman-query-results-event", this.onQueryResultsEvent.bind(this))
  }

  // When we get a query-results event, make the query-results tab visible.
  onQueryResultsEvent(e:Event) {
    const evt = e as CustomEvent<QueryResultsEvent>
    console.log("JracemanApp got updated results", evt.detail.results)
    const shadowRoot = this.shadowRoot
    if (shadowRoot == null) {
      console.error("can't find shadowRoot")
      return
    }
    const bottomTabs = shadowRoot.querySelector("#bottom-tabs") as JracemanTabs
    if (bottomTabs == null) {
      console.error("Can't find bottom-tabs")
      return
    }
    bottomTabs.selectTabById("query-results")
  }

  render() {
    return html`
      <div class="title-bar">JRaceman</div>
      <vaadin-split-layout id="main" orientation="vertical" vertical>
          <div id="top" class="tab-container">
            <jraceman-tabs>
                <span slot="tab">Sport Setup</span>
                <section slot="panel"><sport-setup></sport-setup></section>
                <span slot="tab">Tab 2</span>
                <section slot="panel">Content for tab 2</section>
                <span slot="tab">Tab 3</span>
                <section slot="panel">
                    <jraceman-tabs>
                        <span slot="tab">Tab 3.1</span>
                        <section slot="panel">Content for tab 3.1</section>
                        <span slot="tab">Tab 3.2</span>
                        <section slot="panel">Content for tab 3.2</section>
                    </jraceman-tabs>
                </section>
            </jraceman-tabs>
          </div>
          <div id="bottom" class="tab-container">
            <jraceman-tabs id="bottom-tabs">
                <span slot="tab">Messages</span>
                <section slot="panel">Messages content</section>
                <span slot="tab" id="query-results">Query Results</span>
                <section slot="panel"><query-results></query-results></section>
                <span slot="tab">Report Results</span>
                <section slot="panel">Report Results content</section>
                <span slot="tab">Help</span>
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
