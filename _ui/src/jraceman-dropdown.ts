import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';

// jraceman-dropdown provides a control and content element
// for drop-down menus.
@customElement('jraceman-dropdown')
export class JracemanDropdown extends LitElement {
  content: Element
  onDocumentClickListener: (e:MouseEvent)=>void

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
    this.onDocumentClickListener = this.onDocumentClick.bind(this)
  }

  onControlClick(e: Event) {
    if (this.content.hasAttribute("hidden")) {
      e.stopPropagation()       // Consume the click event
      this.showPopup()
    } else {
      this.hidePopup()
    }
  }

  onContentClick() {
    this.content!.setAttribute("hidden","")
  }

  showPopup() {
    this.content!.removeAttribute("hidden")
    // Add an event listener that grabs clicks outside the popup
    // and takes down the popup.
    document.addEventListener("click", this.onDocumentClickListener)
  }

  // This listener gets added to the document when our dropdown
  // becomes visible. It checks for clicks outside the dropdown
  // and hides the dropdown when it sees one.
  // Note that this does not prevent other listeners from taking their
  // action based on the click.
  onDocumentClick(e: MouseEvent) {
    const isClosest = ((e!.target) as HTMLElement)!.closest("#dropdown");

    if (!isClosest && !this.content.hasAttribute("hidden")) {
      this.hidePopup()
    }
  }

  hidePopup() {
    this.content!.setAttribute("hidden","")
    document.removeEventListener("click", this.onDocumentClickListener)
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
