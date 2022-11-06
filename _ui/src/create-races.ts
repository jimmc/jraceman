import { LitElement, html, css } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';

import { ApiManager } from './api-manager.js'
import { PostError } from './message-log.js'

/**
 * create-races is the tab content that contains other tabs for venue setup.
 *
 * @slot - Slots for contains tabs and tab content
 */
@customElement('create-races')
export class CreateRaces extends LitElement {
  static styles = css`
  `;

  @property()
  eventId = ""

  @state()
  eventSummary = "Select an event"
  @state()
  eventEntryCount = 0

  onCreateRaces() {
    const numEntries = (this.shadowRoot!.querySelector("#entries") as HTMLInputElement)!.value
    console.log("Create races for", numEntries, " entries")
  }

  update(changedProperties: Map<string, unknown>) {
    if (changedProperties.has("eventId")) {
      // When the eventId changes, update our event info.
      this.loadEventInfo()  // No need to await here, just kick it off.
    }
    super.update(changedProperties)
  }

  async loadEventInfo() {
    if (this.eventId == "") {
      this.eventSummary = "(Select an event)"
      this.eventEntryCount = 0
      return
    }
    const path = '/api/app/event/' + this.eventId + '/info'
    let eventInfo = {Summary: "", EntryCount:0}
    try {
      eventInfo = await ApiManager.xhrJson(path)
    } catch (e) {
      console.error(e);
      const errstr = "Error getting event info: " + e/*.responseText*/
      PostError("create-races", errstr)
      return;
    }
    this.eventSummary = eventInfo.Summary
    this.eventEntryCount = eventInfo.EntryCount
  }

  render() {
    return html`
        Selected Event: <span id="eventsummary">${this.eventSummary}</span>
        <br/>
        Entry Count: <span id="evententrycount">${this.eventEntryCount}</span>
        <br/>
        For this event: <button @click="${this.onCreateRaces}">Create Races</button> for <input id="entries" size=4></input> entries
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'create-races': CreateRaces;
  }
}
