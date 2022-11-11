import { LitElement, html, css } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';

import { ApiManager, XhrOptions } from './api-manager.js'
import { PostError } from './message-log.js'

interface RoundCount {
  Count: number,
  Round: number,
  StageName: string,
}

interface EventInfo {
  Summary: string,
  EntryCount: number,
  GroupCount: number,
  GroupSize: number,
  RaceCounts: RoundCount[],
}

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
  eventSummary = ""

  @state()
  entryUnit = "entries"

  async onCreateRaces() {
    const numEntries = (this.shadowRoot!.querySelector("#entries") as HTMLInputElement)!.value
    console.log("Create races for", numEntries, "entries")
    // TODO - check for eventId=="" or bad numEntries and abort if so.
    const path = '/api/app/event/' + this.eventId + '/createraces'
    const params = { entryCount: numEntries }
    const options:XhrOptions = {
      method: 'POST',
      params: params,
    }
    try {
      /*const results =*/ await ApiManager.xhrJson(path, options)
      // TODO look at results
    } catch (e) {
      const evt = e as XMLHttpRequest
      console.error(e);
      const errstr = "Error creating races: " + evt.responseText
      PostError("create-races", errstr)
      return;
    }
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
      this.entryUnit = "entries"
      const detail = this.shadowRoot!.querySelector("#eventDetail")
      if (detail) {
        detail.innerHTML = ""
      }
      return
    }
    const path = '/api/app/event/' + this.eventId + '/info'
    let eventInfo : EventInfo = { Summary:"", EntryCount:0, GroupCount:0, GroupSize:0, RaceCounts:[]}
    try {
      eventInfo = await ApiManager.xhrJson(path)
    } catch (e) {
      console.error(e);
      const errstr = "Error getting event info: " + e/*.responseText*/
      PostError("create-races", errstr)
      return;
    }
    this.eventSummary = "Selected Event: " + eventInfo.Summary
    const inputField = (this.shadowRoot!.querySelector("#entries")! as HTMLInputElement)
    let eventDetailHTML = ""
    if (eventInfo.GroupSize>1) {
      eventDetailHTML += "Number of groups: " + eventInfo.GroupCount + "<br/>"
      this.entryUnit = "groups"
      inputField.value = ""+eventInfo.GroupCount
    } else {
      eventDetailHTML += "Number of entries: " + eventInfo.EntryCount + "<br/>"
      this.entryUnit = "entries"
      inputField.value = ""+eventInfo.EntryCount
    }
    let raceTotal = 0
    let raceInfo = ""
    for (let roundInfo of eventInfo.RaceCounts) {
      raceTotal += roundInfo.Count
      if (raceInfo!="") {
        raceInfo += ", "
      }
      raceInfo += roundInfo.Count + " " + roundInfo.StageName
    }
    if (raceTotal==0) {
      raceInfo = "no races."
    } else {
      raceInfo = "" + raceTotal + " race" + (raceTotal>1?"s":"") + " (" + raceInfo + ")."
    }
    raceInfo = "This event currently has " + raceInfo
    eventDetailHTML += raceInfo + "<br/>"
    this.shadowRoot!.querySelector("#eventdetail")!.innerHTML = eventDetailHTML
  }

  render() {
    return html`
        <span id="eventsummary">${this.eventSummary}</span><br/>
        <span id="eventdetail"></span>
        For this event: <button @click="${this.onCreateRaces}">Create Races</button> for <input id="entries" size=4></input> ${this.entryUnit}
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'create-races': CreateRaces;
  }
}
