import { LitElement, html, css } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { when } from 'lit/directives/when.js'
import { repeat } from 'lit/directives/repeat.js'

import { ApiManager, XhrOptions } from './api-manager.js'
import { PostError } from './message-log.js'
import { EventInfo } from './event-info.js'

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

  @property()
  eventId = ""

  @property()
  fromRounds: FromRoundInfo[] = []

  entries: any  // TODO

  selectedRoundNumber = 0

  update(changedProperties: Map<string, unknown>) {
    if (changedProperties.has("eventId")) {
      // When the eventId changes, update our event info.
      this.loadEventInfo()  // No need to await here, just kick it off.
      this.loadEventEntries()
    }
    super.update(changedProperties)
  }

  async loadEventInfo() {
    if (this.eventId == "") {
      return
    }
    const path = '/api/app/event/' + this.eventId + '/info'
    let eventInfo : EventInfo = {
      Summary: "",
      EntryCount: 0,
      GroupCount: 0,
      GroupSize: 0,
      RoundCounts: [],
      Races: [],
    }
    try {
      eventInfo = await ApiManager.xhrJson(path)
    } catch (e) {
      console.error(e);
      const errstr = "Error getting event info: " + e/*.responseText*/
      PostError("entries-progress", errstr)
      return
    }
    console.log("entries-progress eventInfo", eventInfo)
    this.fromRounds = this.eventInfoToRoundSelectedInfo(eventInfo)
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
      const entries = await ApiManager.xhrJson(path, options)
      console.log("entries-progress entries", entries)
      this.entries = entries
    } catch (e) {
      console.error(e)
      const errstr = "Error getting entries: " + e
      PostError("entries-progress", errstr)
      return
    }
  }

  eventInfoToRoundSelectedInfo(eventInfo: EventInfo): FromRoundInfo[] {
    let rounds: FromRoundInfo[] = []
    let roundInfo: FromRoundInfo = {
      RoundNumber: 0,
      StageName: "Draw",
    }
    rounds.push(roundInfo)
    for (let roundCount of eventInfo.RoundCounts) {
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
