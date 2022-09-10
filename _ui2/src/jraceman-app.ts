import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import '@vaadin/vaadin-split-layout/vaadin-split-layout.js';

import './meet-setup.js'
import './plan-setup.js'
import './query-results.js'
import './sport-setup.js'
import './team-setup.js'
import './venue-setup.js'
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
                <span slot="tab">Plan Setup</span>
                <section slot="panel"><plan-setup></sport-setup></section>
                <span slot="tab">Venue Setup</span>
                <section slot="panel"><venue-setup></sport-setup></section>
                <span slot="tab">Team Setup</span>
                <section slot="panel"><team-setup></sport-setup></section>
                <span slot="tab">Meet Setup</span>
                <section slot="panel"><meet-setup></sport-setup></section>
                <span slot="tab">By Event</span>
                <section slot="panel">By Event is not yet implemented</section>
                <span slot="tab">Reports</span>
                <section slot="panel">Reports is not yet implemented</section>
                <span slot="tab">Database</span>
                <section slot="panel">Database is not yet implemented</section>
                <span slot="tab">Debug</span>
                <section slot="panel">Debug is not yet implemented</section>
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
