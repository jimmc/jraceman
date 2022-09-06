import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';

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

  selectTab(tabIndex: number) {
    this.tabs.forEach((tab:Element) => tab.removeAttribute("selected"))
    this.tabs[tabIndex].setAttribute("selected", "")
    this.panels.forEach((panel:Element) => panel.removeAttribute("selected"))
    this.panels[tabIndex].setAttribute("selected", "")
  }

  handleSelect(e:Event) {
    const index = this.tabs.indexOf(e.target as Element)
    this.selectTab(index)
  }

  render() {
    return html`
      <nav>
        <slot name="tab" @click=${(e:Event) => this.handleSelect(e)}></slot>
      </nav>
      <slot name="panel"></slot>
    `;
  }
}
