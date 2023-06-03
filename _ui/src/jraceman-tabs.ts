import { LitElement, html, css } from 'lit'
import { customElement } from 'lit/decorators.js'

// jraceman-tabs provides a tab bar and slots for tabs and content.
// This ia a typescript version of
//   https://medium.com/@blueyorange/make-a-tabs-web-component-in-litelement-using-slots-and-css-293c55dbd155
// Thank you Russell Johnson.
@customElement('jraceman-tabs')
export class JracemanTabs extends LitElement {
  tabs: Element[] = []
  panels: Element[] = []

  static styles = css`
    nav {
      display: flex;
    }
    nav > ::slotted([slot="tab"]) {
      padding: 0.1rem 0.2rem;
      margin: 0.1rem 0.1rem;
      flex: 1 1 auto;
      color: var(--color-darkGrey);
      border-bottom: 2px solid lightgrey;
      text-align: center;
      font-size: large;
      font-weight: bold;
    }
    nav > ::slotted([slot="tab"][selected]) {
      border-color: black;
    }
    ::slotted([slot="panel"]) {
      display: none;
    }

    ::slotted([slot="panel"][selected]) {
      display: block;
    }
  `;

  constructor() {
    super()
    this.tabs = Array.from(this.querySelectorAll(":scope > [slot=tab]"))
    this.panels = Array.from(this.querySelectorAll(":scope > [slot=panel]"))
    this.selectTab(0)
  }

  // selectTabById can be called from another element to select a tab by
  // the id of that tab slot.
  selectTabById(id: string) {
    const tab = this.querySelector("#"+id)
    if (tab == null) {
      console.error("can't find tab", id)
      return
    }
    const index = this.tabs.indexOf(tab)
    this.selectTab(index)
  }

  selectTab(tabIndex: number) {
    this.tabs.forEach((tab:Element) => tab.removeAttribute("selected"))
    this.tabs[tabIndex].setAttribute("selected", "")
    this.panels.forEach((panel:Element) => panel.removeAttribute("selected"))
    this.panels[tabIndex].setAttribute("selected", "")
    // Dispatch an event so others can take action when a tab changes.
    const event = new Event('jraceman-tab-selected-event', {
      bubbles: true,
      composed: true
    })
    this.dispatchEvent(event)
  }

  onSelect(e:Event) {
    const index = this.tabs.indexOf(e.target as Element)
    if (index<0) {
      // Not our tab, assume it's some nested element in the tab that we can ignore.
      return
    }
    this.selectTab(index)
  }

  onRequestDisplay(e:Event) {
    const target = e.target
    const panelIndex = this.panels.findIndex(p => p.children[0]==target)
    if (panelIndex < 0) {
      console.error("can't find panel", target)
      return
    }
    this.selectTab(panelIndex)  // Display the tab at this level.
    // Let the event propagate up to containing tabs as well.
  }

  render() {
    return html`
      <nav>
        <slot name="tab" @click=${(e:Event) => this.onSelect(e)}></slot>
      </nav>
      <div @jraceman-request-display-event=${(e:Event) => this.onRequestDisplay(e)}>
        <slot name="panel"></slot>
      </div>
    `;
  }
}
