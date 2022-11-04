import { LitElement, html, css } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { when } from 'lit/directives/when.js'

import './jraceman-tabs.js'
import './table-queryedit.js'

import { ApiHelper, KeySummary } from './api-helper.js'
import { PostError } from './message-log.js'

/**
 * by-event is the tab content for the By Event tab that allows operations on events.
 */
@customElement('by-event')
export class ByEvent extends LitElement {
  static styles = css`
  `;

  @property()
  meetItems: KeySummary[] = []

  @property()
  byChoice = ""

  @property()
  task = ""

  connectedCallback() {
    super.connectedCallback()
    this.loadMeetChoices()      // No need to await here
  }

  async loadMeetChoices() {
    try {
      this.meetItems = await ApiHelper.loadKeySummaries("meet")
      this.onMeetChange()
    } catch(e) {
      console.error("Error getting meet table summary: ", e)
      const evt = e as XMLHttpRequest
      PostError("reports", evt.responseText)
    }
  }

  onMeetChange() {
    const meetId = (this.shadowRoot!.querySelector("#val_meet") as HTMLSelectElement)!.value
    console.log("Meet changed to", meetId)
    // TODO: update the event/race choice lists
    this.onByChoiceChange()     // Init the ByChoice list
    this.onTaskChange()         // Init the task pane
  }

  onByChoiceChange() {
    this.byChoice = (this.shadowRoot!.querySelector("#by_choice") as HTMLSelectElement)!.value
    console.log("By-choice changed to", this.byChoice) // TODO
  }

  onTaskChange() {
    this.task = (this.shadowRoot!.querySelector("#task") as HTMLSelectElement)!.value
    console.log("Task-choice changed to", this.task) // TODO
  }

  render() {
    return html`
        Meet:
        <select id="val_meet" @change="${this.onMeetChange}">
          ${repeat(this.meetItems, (keyitem)=>html`
            <option value="${keyitem.ID}">${keyitem.Summary}</option>
          `)}
        </select>
        <br/>
        <select id="by_choice" @change="${this.onByChoiceChange}">
          <option value="by_event_number">By Event #</option>
          <option value="by_race_number">By Race #</option>
          <option value="by_event_name">By Event Name</option>
        </select>
        ${when(this.byChoice=="by_event_number",()=>html`[by event number list]`)}
        ${when(this.byChoice=="by_race_number",()=>html`[by race number list]`)}
        ${when(this.byChoice=="by_event_name",()=>html`[by event name list]`)}
        <select id="task" @change="${this.onTaskChange}">
          <option value="create_races">Create Races</option>
          <option value="entries_progress">Entries/Progress</option>
          <option value="results">Results</option>
          <option value="reports">Reports</option>
        </select>
        <br/>
        ${when(this.task=="create_races",()=>html`[create races pane]`)}
        ${when(this.task=="entries_progress",()=>html`[entries/progress pane]`)}
        ${when(this.task=="results",()=>html`[results pane]`)}
        ${when(this.task=="reports",()=>html`[reports pane]`)}
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'by-event': ByEvent;
  }
}
