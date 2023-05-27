import { LitElement, html, css } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { when } from 'lit/directives/when.js'

import { ApiManager, XhrOptions } from './api-manager.js'
import { PostError, PostInfo } from './message-log.js'
import { EventRaces, RaceInfo } from './event-races.js'

interface CreateRacesResult {
  EventRaces: EventRaces,
  RacesToCreate: RaceInfo[],
  RacesToDelete: RaceInfo[],
  RacesToModFrom: RaceInfo[],
  RacesToModTo: RaceInfo[],
}

/**
 * create-races is the by-event tab that for creating the races for an event.
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

  @state()
  hasLanes = false

  async onCreateRaces() {
    await this.callCreateRaces(false)
  }
  async onDryRun() {
    await this.callCreateRaces(true)
  }
  async callCreateRaces(dryRun: boolean) {
    const numEntries = (this.shadowRoot!.querySelector("#entries") as HTMLInputElement)!.value
    console.log("Create races for", numEntries, "entries")
    if (this.eventId == "") {
      PostError("create-races", "No event selected")
      return
    }
    const allowDeleteLanes = !!(this.shadowRoot!.querySelector("#allowdelete:checked"))
    console.log("allowdelete is", allowDeleteLanes)
    const path = '/api/app/event/' + this.eventId + '/createraces'
    const params = {
      laneCount: numEntries,
        // laneCOunt is number of groups for group events, number of entries for non-group events.
      dryRun: dryRun,
        // If dryRun is true, tell us what would happen without actually doing it.
      allowDeleteLanes:  allowDeleteLanes,
        // If allowDeleteLanes is false and we try to delete a race, that will be an error.
    }
    const options:XhrOptions = {
      method: 'POST',
      params: params,
    }
    var results
    try {
      const response = await ApiManager.xhrJson(path, options)
      results = response as CreateRacesResult
      console.log("Results of Create Races:", results)
    } catch (e) {
      const evt = e as XMLHttpRequest
      console.error(e);
      const errstr = "Error creating races: " + evt.responseText
      PostError("create-races", errstr)
      return;
    }
    const changeCount = results.RacesToCreate.length + results.RacesToDelete.length + results.RacesToModFrom.length
    if (changeCount==0) {
      if (dryRun) {
        PostInfo("create-races", "No changes to races would be made for event " + results.EventRaces.Summary)
      } else {
        PostInfo("create-races", "No changes to races were made for event " + results.EventRaces.Summary)
      }
      return
    }
    if (dryRun) {
      PostInfo("create-races", "The following changes would be made for event " + results.EventRaces.Summary)
    } else {
      PostInfo("create-races", "The following changes were made for event " + results.EventRaces.Summary)
    }
    results.RacesToCreate.forEach( (race) => {
      PostInfo("create-races", "+ Create Race "+this.raceToString(race))
    })
    let deleteWouldFail = false
    results.RacesToDelete.forEach( (race) => {
      if (race.LaneCount > 0 && !allowDeleteLanes) {
        deleteWouldFail = true
      }
      PostInfo("create-races", "+ Delete Race "+this.raceToString(race))
    })
    for (let i=0; i<results.RacesToModFrom.length; i++) {
      const raceFrom = results.RacesToModFrom[i]
      const raceTo = results.RacesToModTo[i]
      PostInfo("create-races", "+ Modify Race "+this.raceToString(raceFrom) + " to " + this.raceToString(raceTo))
    }
    if (deleteWouldFail) {
      PostInfo("create-races", "* NOTE: This operation would fail because some races to be deleted have lane data, and 'Allow deleting races with lane data' is not selected")
    }
  }

  update(changedProperties: Map<string, unknown>) {
    if (changedProperties.has("eventId")) {
      // When the eventId changes, update our event info.
      this.loadEventRaces()  // No need to await here, just kick it off.
    }
    super.update(changedProperties)
  }

  raceToString(race: RaceInfo):string {
    var raceNumber = ""
    if (race.RaceNumber) {
      raceNumber = " #" + race.RaceNumber
    }
    var idinfo = ""
    if (race.RaceID) {
      idinfo = " [" + race.RaceID + "]"
    }
    var laneInfo = ""
    if (race.LaneCount>0) {
      laneInfo = " (NOTE: This race includes lane data)"
    }
    return "Stage="+race.StageName+" Round="+race.Round+" Section="+race.Section + raceNumber + idinfo + laneInfo
  }

  async loadEventRaces() {
    if (this.eventId == "") {
      this.eventSummary = "(Select an event)"
      this.entryUnit = "entries"
      const detail = this.shadowRoot!.querySelector("#eventDetail")
      if (detail) {
        detail.innerHTML = ""
      }
      return
    }
    const path = '/api/app/event/' + this.eventId + '/races'
    let eventRaces : EventRaces = {
      Summary: "",
      EntryCount: 0,
      GroupCount: 0,
      GroupSize: 0,
      RoundCounts: [],
      Races: [],
      }
    try {
      eventRaces = await ApiManager.xhrJson(path)
    } catch (e) {
      console.error(e);
      const errstr = "Error getting event info: " + e/*.responseText*/
      PostError("create-races", errstr)
      return;
    }
    this.eventSummary = "Selected Event: " + eventRaces.Summary
    const inputField = (this.shadowRoot!.querySelector("#entries")! as HTMLInputElement)
    let eventDetailHTML = ""
    if (eventRaces.GroupSize>1) {
      eventDetailHTML += "Number of groups: " + eventRaces.GroupCount + "<br/>"
      this.entryUnit = "groups"
      inputField.value = ""+eventRaces.GroupCount
    } else {
      eventDetailHTML += "Number of entries: " + eventRaces.EntryCount + "<br/>"
      this.entryUnit = "entries"
      inputField.value = ""+eventRaces.EntryCount
    }
    let raceTotal = 0
    let raceSummary = ""
    for (let roundInfo of eventRaces.RoundCounts) {
      raceTotal += roundInfo.Count
      if (raceSummary!="") {
        raceSummary += ", "
      }
      raceSummary += roundInfo.Count + " " + roundInfo.StageName
    }
    if (raceTotal==0) {
      raceSummary = "no races."
    } else {
      raceSummary = "" + raceTotal + " race" + (raceTotal>1?"s":"") + " (" + raceSummary + ")."
    }
    raceSummary = "This event currently has " + raceSummary
    eventDetailHTML += raceSummary + "<br/>"
    this.hasLanes = false
    for (let raceInfo of eventRaces.Races) {
      if (raceInfo.LaneCount>0) {
        this.hasLanes = true
        break
      }
    }
    this.shadowRoot!.querySelector("#eventdetail")!.innerHTML = eventDetailHTML
  }

  render() {
    return html`
        <span id="eventsummary">${this.eventSummary}</span><br/>
        <span id="eventdetail"></span>
        For this event: <button @click="${this.onCreateRaces}">Create Races</button>
            for <input id="entries" size=4></input> ${this.entryUnit}
            <button @click="${this.onDryRun}">Dry Run</button>
        ${when(this.hasLanes,()=>html`<br/>
            <b>NOTE: Some races include lane data</b><br/>
            <input type="checkbox" id="allowdelete" name="allowdelete"></input>
            <label for="allowdelete">Allow deleting races with lane data</label>
        `)}
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'create-races': CreateRaces;
  }
}
