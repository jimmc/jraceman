import { LitElement, html, css } from 'lit'
import { PropertyValues } from 'lit-element'
import { customElement, property } from 'lit/decorators.js'
import { when } from 'lit/directives/when.js'
import { repeat } from 'lit/directives/repeat.js'

import './sheet-editor.js'

import { ApiHelper } from './api-helper.js'
import { ApiManager, XhrOptions } from './api-manager.js'
import { EventRaces } from './event-races.js'
import { PostError } from './message-log.js'
import { ProgressionLanes } from './progression-lanes.js'
import { SheetEditor } from './sheet-editor.js'
import { QueryResultsData, TableDesc, TableDescSupport } from './table-desc.js'

interface FromRoundInfo {
  RoundNumber: number;
  StageName: string;
}

/**
 * entries-progress is the tab content that contains other tabs for venue setup.
 *
 * @slot - Slots for contains tabs and tab content
 */
@customElement('entries-progress')
export class EntriesProgress extends LitElement {
  static styles = css`
  `;

  // eventId is for the event as selected in our parent and passed in to us.
  @property()
  eventId = ""

  // eventRaces gets loaded when eventId changes.
  eventRaces?: EventRaces;

  // fromRounds is the list of rounds for this race, less one.
  // Round 0 is the "draw" round, the rest are for progressing from that round.
  @property()
  fromRounds: FromRoundInfo[] = []

  // entryTableDesc is the table descriptor for the entry table.
  entryTableDesc: TableDesc = {Table:'', Columns:[]}

  // entries is the set of entry records for the selected event.
  entries: QueryResultsData = TableDescSupport.emptyQueryResults()

  // selectedRoundNumber is the "from" round number as selected by the user.
  @property()
  selectedRoundNumber = 0

  // sheetEditor is our editing sheet
  sheetEditor?: SheetEditor

  // sheetTableDesc is the tableDesc we send to the sheetEditor.
  @property()
  sheetTableDesc: TableDesc = {
    Table: '(unset in entries-progress)',
    Columns: [],
  }

  // sheetQueryResults is the row data we send to the sheetEditor.
  @property()
  sheetQueryResults: QueryResultsData = TableDescSupport.emptyQueryResults()

  firstUpdated(changedProperties:PropertyValues<this>) {
    super.firstUpdated(changedProperties);
    this.sheetEditor = this.shadowRoot!.querySelector("sheet-editor")! as SheetEditor
  }

  async update(changedProperties: Map<string, unknown>) {
    if (changedProperties.has("eventId") || changedProperties.has("selectedRoundNumber")) {
      // When the eventId changes, update our event info.
      // Loading the races doesn't depend on the entries, so we could do these
      // two in parallel. We need the races to load the lanes, and we need the
      // entries and the lanes to collect the sheetQueryResults.
      await this.loadEventEntries()
      await this.loadEventRaces()
      await this.loadEventLanes()
      this.sheetTableDesc = ProgressionLanes.lanesFromRoundTableDesc()
      this.sheetQueryResults = ProgressionLanes.collectLanesFromRound(
          this.entryTableDesc, this.entries, this.eventRaces!, this.selectedRoundNumber)
    }
    super.update(changedProperties)
  }

  async loadEventRaces() {
    if (this.eventId == "") {
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
      PostError("entries-progress", errstr)
      return
    }
    console.log("entries-progress eventRaces", eventRaces)
    this.eventRaces = eventRaces
    this.fromRounds = this.eventRacesToRoundSelectedInfo(eventRaces)
  }

  async loadEventEntries() {
    if (this.eventId == "") {
      return
    }
    const path = '/api/query/entry/'
    const params = [
      { name: "eventId", op: "eq", value: this.eventId },
    ]
    const options: XhrOptions = {
      method: "POST",
      params: params,
    }
    try {
      this.entryTableDesc = await ApiHelper.loadTableDesc('entry')
      const entries = await ApiManager.xhrJson(path, options) as QueryResultsData
      console.log("entries-progress entries", entries)
      this.entries = entries
    } catch (e) {
      console.error(e)
      const errstr = "Error getting entries: " + e
      PostError("entries-progress", errstr)
      return
    }
  }

  // Load the lanes for all of the races that have been loaded into this.eventRaces.
  async loadEventLanes() {
    if (!this.eventRaces || !this.eventRaces.Races) {
      console.error("Attempted to load event lanes but there are no races")
    }
    const races = this.eventRaces!.Races
    if (races.length==0) {
      console.info("No races in event, so no lanes to load")
    }
    for (let race of races) {
      if (race.RaceID=="") {
        console.log("No raceID for this race, not fetching lanes")
      } else {
        // TODO - we only need races for fromRound and fromRound+1
        const lanes = await this.loadRaceLanes(race.RaceID)
        console.log("lanes for RaceID", race.RaceID, lanes)
        race.Lanes = lanes
      }
    }
  }

  // Load all of the lanes for the specified race.
  async loadRaceLanes(raceId: string): Promise<QueryResultsData> {
    const path = '/api/query/lane/'
    const params = [
      { name: "raceId", op: "eq", value: raceId },
    ]
    const options: XhrOptions = {
      method: "POST",
      params: params,
    }
    try {
      const raceLanes = await ApiManager.xhrJson(path, options) as QueryResultsData
      if (!raceLanes.Table) {
        raceLanes.Table = 'lane'
      }
      console.log("entries-progress raceLanes", raceLanes)
      return raceLanes
    } catch (e) {
      console.error(e)
      const errstr = "Error getting entries: " + e
      PostError("entries-progress", errstr)
      return TableDescSupport.emptyQueryResults()
    }
  }

  eventRacesToRoundSelectedInfo(eventRaces: EventRaces): FromRoundInfo[] {
    let rounds: FromRoundInfo[] = []
    let roundInfo: FromRoundInfo = {
      RoundNumber: 0,
      StageName: "Draw",
    }
    rounds.push(roundInfo)
    for (let roundCount of eventRaces.RoundCounts) {
      let roundInfo: FromRoundInfo = {
        RoundNumber: roundCount.Round,
        StageName: roundCount.StageName,
      }
      rounds.push(roundInfo)
    }
    return rounds
  }

  onRoundChange() {
    const newSelectedRoundNumber = (this.shadowRoot!.querySelector('#round_list') as HTMLSelectElement)!.value
    console.log("onRoundChange, selected round number is now", newSelectedRoundNumber)
    this.selectedRoundNumber = parseInt(newSelectedRoundNumber)
    this.requestUpdate()
  }

  render() {
    return html`
      ${when(this.eventId,
        ()=>html`
          From Round:
          <select id="round_list" @change="${this.onRoundChange}">
            ${repeat(this.fromRounds, (round)=>html`
              <option value="${round.RoundNumber}" ?selected=${this.selectedRoundNumber==round.RoundNumber}>
                ${round.RoundNumber}: ${round.StageName}
              </option>
            `)}
          </select>
          Ready to draw First Round
          <button>Add Entry</button>
          <button>Delete Entry</button>
          <button>Draw</button>
          <button>Undraw</button>
          <br/>
          <sheet-editor .tableDesc=${this.sheetTableDesc}
            .queryResults=${this.sheetQueryResults}>
          </sheet-editor>
        `,
        ()=>html`Select an Event`
      )}
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'entries-progress': EntriesProgress;
  }
}
