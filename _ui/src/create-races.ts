import { LitElement, html, css } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';

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
  eventInfo = "Select an event"

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
    this.eventInfo = await ("Info for event "+this.eventId)
  }

  render() {
    return html`
        <span id="eventinfo">${this.eventInfo}</span>
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
