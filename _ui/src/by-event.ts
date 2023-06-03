import { LitElement, html, css } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { repeat } from 'lit/directives/repeat.js'
import { when } from 'lit/directives/when.js'

import './create-races.js'
import './entries-progress.js'
import './jraceman-tabs.js'
import './table-manager.js'

import { ApiHelper, KeySummary } from './api-helper.js'
import { PostError } from './message-log.js'

/**
 * by-event is the tab content for the By Event tab that allows operations on events.
 */
@customElement('by-event')
export class ByEvent extends LitElement {
  static styles = css`
  `;

  @state()
  meetItems: KeySummary[] = []

  @state()
  meetId = ""

  @state()
  eventItems: KeySummary[] = []

  @state()
  eventId = ""

  @state()
  byChoice = ""

  @state()
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
      PostError("by-event", evt.responseText)
    }
  }

  async loadEventChoices() {
    try {
      let params = []
      const meetWhere = {
        name: "meetId",
        op: "eq",
        value: this.meetId,
      }
      params.push(meetWhere)
      // The contents of the event list depends on the setting of the by-choice.
      if (this.byChoice == "by_event_number") {
        this.eventItems = await ApiHelper.loadKeySummaries("event", params, "bynumber")
      } else if (this.byChoice == "by_event_name") {
        this.eventItems = await ApiHelper.loadKeySummaries("event", params, "byname")
      } else if (this.byChoice == "by_event_id") {
        this.eventItems = await ApiHelper.loadKeySummaries("event", params, "byid")
      } else if (this.byChoice == "by_race_number") {
        this.eventItems = await ApiHelper.loadKeySummaries("event", params, "byracenumber")
      } else {
        PostError("by-event", 'Bad byChoice value "' + this.byChoice + '"')
      }
      // Let the event choice list get updated, then propagate possible change of selected item.
      setTimeout(()=>this.onEventChange(),0)
    } catch(e) {
      console.error("Error getting event table summary: ", e)
      const evt = e as XMLHttpRequest
      PostError("by-event", evt.responseText)
    }
  }

  onMeetChange() {
    this.meetId = (this.shadowRoot!.querySelector("#meet_list") as HTMLSelectElement)!.value
    console.log("Meet changed to", this.meetId)
    this.onByChoiceChange()     // Init the ByChoice list
    this.onTaskChange()         // Init the task pane
  }

  onByChoiceChange() {
    this.byChoice = (this.shadowRoot!.querySelector("#by_choice") as HTMLSelectElement)!.value
    console.log("By-choice changed to", this.byChoice)
    setTimeout(()=>this.loadEventChoices())        // After dependent components are available, update them.
  }

  onEventChange() {
    this.eventId = (this.shadowRoot!.querySelector("#event_list") as HTMLSelectElement)!.value
    console.log("Event changed to", this.eventId) // TODO
  }

  onTaskChange() {
    this.task = (this.shadowRoot!.querySelector("#task") as HTMLSelectElement)!.value
    console.log("Task-choice changed to", this.task) // TODO
  }

  render() {
    return html`
        Meet:
        <select id="meet_list" @change="${this.onMeetChange}">
          ${repeat(this.meetItems, (keyitem)=>html`
            <option value="${keyitem.ID}">${keyitem.Summary}</option>
          `)}
        </select>
        <br/>
        <select id="by_choice" @change="${this.onByChoiceChange}">
          <option value="by_event_number">By Event #</option>
          <option value="by_event_name">By Event Name</option>
          <option value="by_event_id">By Event ID</option>
          <option value="by_race_number">By Race #</option>
        </select>
        <select id="event_list" @change="${this.onEventChange}">
          ${repeat(this.eventItems, (keyitem)=>html`
            <option value="${keyitem.ID}" ?selected=${this.eventId==keyitem.ID}>${keyitem.Summary}</option>
          `)}
        </select>
        <select id="task" @change="${this.onTaskChange}">
          <option value="create_races">Create Races</option>
          <option value="entries_progress">Entries/Progress</option>
          <option value="results">Results</option>
          <option value="reports">Reports</option>
        </select>
        <br/>
        ${when(this.task=="create_races",()=>html`<create-races .eventId="${this.eventId}"></create-races>`)}
        ${when(this.task=="entries_progress",()=>html`<entries-progress .eventId="${this.eventId}"></entries-progress>`)}
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
