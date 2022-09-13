import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';

// jraceman-dropdown provides a control and content element
// for drop-down menus.
@customElement('jraceman-dropdown')
export class JracemanDropdown extends LitElement {
  content: Element

  static styles = css`
    #dropdown {
      position: relative;
      display: inline-block;
    }
    ::slotted([slot="content"]) {
      display: block;
      position: absolute;
      background-color: #f9f9f9;
      min-width: 160px;
      box-shadow: 0px 8px 16px 0px rgba(0,0,0,0.2);
      padding: 12px 16px;
      z-index: 1;
    }
    ::slotted([slot="content"][hidden]) {
      display: none;
    }
  `;

  constructor() {
    super()
    this.content = this.querySelector("[slot=content]")!
    this.content!.setAttribute("hidden","")
  }

  onControlClick() {
    if (this.content.hasAttribute("hidden")) {
      this.content!.removeAttribute("hidden")
    } else {
      this.content!.setAttribute("hidden","")
    }
  }

  onContentClick() {
    this.content!.setAttribute("hidden","")
  }

  render() {
    return html`
      <div id="dropdown">
        <slot name="control" @click=${this.onControlClick}></slot>
        <div @click=${this.onContentClick}>
          <slot name="content"></slot>
        </div>
      </div>
    `;
  }
}
