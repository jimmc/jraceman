import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {when} from 'lit/directives/when.js';

import './auth-setup.js'
import './database-menu.js'
import './database-pane.js'
import './debug-pane.js'
import './jraceman-login.js'
import './jraceman-split.js'
import './meet-setup.js'
import './message-log.js'
import './plan-setup.js'
import './query-results.js'
import './reports-pane.js'
import './report-results.js'
import './sport-setup.js'
import './team-setup.js'
import './venue-setup.js'
import { ApiManager } from './api-manager.js'
import { JracemanLogin, LoginStateEvent } from './jraceman-login.js'
import { JracemanTabs} from './jraceman-tabs.js'
import { Message, PostMessageEvent, PostError } from './message-log.js'
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
      height: 100vh;
      width: 100vw;
      margin: 0;
      padding: 0;
      display: flex;
      flex-direction: column;
    }
    .title-bar {
      height: 16px
      width: "100%";
      background-color: lightgray;
      color: black;
      flex: 0 0 auto;
    }
    .right {
      float: right;
    }
    #main {
      height: 0;        /* Let flex fill the size; this prevents it from resizing when content changes. */
      flex: 1 1 0;
    }
    #login {
      display: block;
    }
    [hidden="true"] {
      display: none !important;
    }
  `;

  @property()
  jracemanVersion = ''

  @property()
  unviewedMessageCount = 0

  @property()
  loggedIn = false

  connectedCallback() {
    super.connectedCallback()
    // We add a listener for query results so that we can make
    // the query results tab visible.
    document.addEventListener("jraceman-query-results-event", this.onQueryResultsEvent.bind(this))
    document.addEventListener("jraceman-report-results-event", this.onReportResultsEvent.bind(this))
    document.addEventListener("jraceman-post-message-event", this.onPostMessage.bind(this))
    document.addEventListener("jraceman-login-state-event", this.onLoginState.bind(this))
    JracemanLogin.AnnounceLoginState()  // See if we are logged in.

    this.loadVersion()
  }

  async loadVersion() {
    const queryPath = "/api0/version"
    const options = {}
    let result
    try {
      result = await ApiManager.xhrJson(queryPath, options);
    } catch(e) {
      PostError("jraceman-app", "Error from /api0/version: " + e/*.responseText*/);
      console.error("Error getting version from /api0/version:", e/*.responseText*/);
      return
    }
    this.jracemanVersion = result as string
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

  // When we get a post-message event, count it.
  onLoginState(e:Event) {
    const evt = e as CustomEvent<LoginStateEvent>
    console.log("JracemanApp.onLoginState", evt)
    this.loggedIn = evt.detail.State
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

  // When we get a post-message event, count it.
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
    const evt = e as CustomEvent<PostMessageEvent>
    const m:Message = evt.detail.message
    if (m.level.toLowerCase() == "error") {
      // If an error message is posted, display the Messages tab.
      const bottomTabs = this.shadowRoot!.querySelector("#bottom-tabs")! as JracemanTabs
      bottomTabs.selectTabById("message-log-tab")
    }
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

  logout() {
    console.log("Logging out")
    const jl = this.shadowRoot!.querySelector("#login") as JracemanLogin
    jl.logout()
  }

  render() {
    return html`
      <div class="title-bar">JRaceman ${this.jracemanVersion}
        ${when(this.loggedIn,()=>html`
          <a href="#" class="right" @click="${this.logout}">Logout</a>
        `)}
      </div>
      <jraceman-login id="login" hidden=${this.loggedIn} logged-in=${this.loggedIn}></jraceman-login>
      ${when(this.loggedIn,()=>html`
        <jraceman-split id="main">
          <div id="top" slot="top" class="tab-container">
            <jraceman-tabs>
              <span slot="tab">Auth Setup</span>
              <section slot="panel"><auth-setup></auth-setup></section>
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
              <span slot="tab"><database-menu></database-menu>Database</span>
              <section slot="panel"><database-pane></database-pane></section>
              <span slot="tab">Debug</span>
              <section slot="panel"><debug-pane></debug-pane></section>
            </jraceman-tabs>
          </div>
          <div id="bottom" slot="bottom" class="tab-container">
            <jraceman-tabs @jraceman-tab-selected-event=${this.onTabSelected} id="bottom-tabs">
              <span slot="tab" id="message-log-tab">Messages${when(this.unviewedMessageCount>0,
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
        </jraceman-split>
      `)}
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'jraceman-app': JracemanApp;
  }
}
