import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {when} from 'lit/directives/when.js';
import '@vaadin/vaadin-split-layout/vaadin-split-layout.js';

import './database-pane.js'
import './debug-pane.js'
import './meet-setup.js'
import './message-log.js'
import './plan-setup.js'
import './query-results.js'
import './reports-pane.js'
import './report-results.js'
import './sport-setup.js'
import './team-setup.js'
import './venue-setup.js'
import { JracemanTabs} from './jraceman-tabs.js'
import { ReportResultsEvent } from './reports-pane.js'
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

  @property()
  unviewedMessageCount = 0

  connectedCallback() {
    super.connectedCallback()
    // We add a listener for query results so that we can make
    // the query results tab visible.
    document.addEventListener("jraceman-query-results-event", this.onQueryResultsEvent.bind(this))
    document.addEventListener("jraceman-report-results-event", this.onReportResultsEvent.bind(this))
    document.addEventListener("jraceman-post-message-event", this.onPostMessage.bind(this))
  }

  // Get a pointer to the message-log-pane
  getMessageLogPanel() {
    const shadowRoot = this.shadowRoot
    if (!shadowRoot) {
      console.error("no shadow root")
      return null
    }
    const panel = shadowRoot.querySelector("#message-log-pane")
    if (!panel) {
      console.error("Can't find message-log-pane")
      return null
    }
    return panel
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
    bottomTabs.selectTabById("query-results-tab")
  }

  // When we get a report-results event, make the report-results tab visible.
  onReportResultsEvent(e:Event) {
    const evt = e as CustomEvent<ReportResultsEvent>
    console.log("JracemanApp got updated results", evt.detail.results)
    const bottomTabs = this.shadowRoot!.querySelector("#bottom-tabs")! as JracemanTabs
    bottomTabs.selectTabById("report-results-tab")
  }

  // WHen we get a post-message event, count it.
  onPostMessage(e:Event) {
    console.log("JracemanApp.onPostMessage", e)
    const messageLogPanel = this.getMessageLogPanel()
    if (!messageLogPanel) {
      return    // console.error already called
    }
    if (messageLogPanel.hasAttribute("selected")) {
      console.log("message panel is selected")
      return    // Messages are being viewed
    }
    // The message panel is not currently selected
    this.unviewedMessageCount++
  }

  // When a tab is selected, see if we need to clear the message count.
  onTabSelected(e:Event) {
    console.log("JracemanApp.onTabSelected", e)
    const messageLogPanel = this.getMessageLogPanel()
    if (!messageLogPanel) {
      return    // console.error already called
    }
    if (!messageLogPanel.hasAttribute("selected")) {
      console.log("message panel is not selected")
      return
    }
    // The message panel is selected, clear the unviewed count
    this.unviewedMessageCount = 0
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
                <section slot="panel"><reports-pane></reports-pane></section>
                <span slot="tab">Database</span>
                <section slot="panel"><database-pane></database-pane></section>
                <span slot="tab">Debug</span>
                <section slot="panel"><debug-pane></debug-pane></section>
            </jraceman-tabs>
          </div>
          <div id="bottom" class="tab-container">
            <jraceman-tabs @jraceman-tab-selected-event=${this.onTabSelected} id="bottom-tabs">
                <span slot="tab">Messages${when(this.unviewedMessageCount>0,
                  ()=>html` [+${this.unviewedMessageCount}]`
                )}</span>
                <section slot="panel" id="message-log-pane"><message-log></message-log></section>
                <span slot="tab" id="query-results-tab">Query Results</span>
                <section slot="panel"><query-results></query-results></section>
                <span slot="tab" id="report-results-tab">Report Results</span>
                <section slot="panel"><report-results></report-results></section>
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
