import { LitElement, html, css } from 'lit'
import { customElement } from 'lit/decorators.js'

// jraceman-split provides a split screen with top and bottom slots.
@customElement('jraceman-split')
export class JracemanSplit extends LitElement {
  static styles = css`
    :host {
      display: flex;
      flex-direction: column;
      margin: 0;
      padding: 0;
    }
    #divider {
      width: 100vw;
      flex: 0 0 5px;
      background-color: black;
      cursor: ns-resize;
    }
    ::slotted([slot="top"]) {
      width: 100vw;
      flex: 1 1 50%;
      overflow: auto;
    }
    ::slotted([slot="bottom"]) {
      width: 100vw;
      flex: 1 1 50%;
      overflow: auto;
    }
  `;

  top?: HTMLElement
  bottom?: HTMLElement

  connectedCallback() {
    super.connectedCallback()
    this.top = this.querySelector("[slot=top]")! as HTMLElement
    this.bottom = this.querySelector("[slot=bottom]")! as HTMLElement
  }

  onMouseDown(md: MouseEvent) {
    md.preventDefault()          // Turn off other event processing, such as selecting text on mouse drag.
    if (!this.top || !this.bottom) {
      console.error("Can't find one of our panes (top or bottom)")
      return
    }
    const topPane = this.top
    const bottomPane = this.bottom

    const sizeProp = "offsetHeight"
    const posProp = "pageY"
    let lastPos = md[posProp]
    let topSize = topPane[sizeProp]
    let bottomSize = bottomPane[sizeProp]
    const totalSize = topSize + bottomSize
    //const totalGrow = Number(topPane.style.flexGrow) + Number(bottomPane.style.flexGrow)
    const totalGrow = 2 // Attempting to read topPane.style.flexGrow returns empty string

    function onMouseMove(mm: MouseEvent) {
      const pos = mm[posProp]
      let delta = pos - lastPos
      lastPos = pos
      if (delta + topSize < 0) {
        delta = -topSize
      } else if (bottomSize - delta < 0) {
        delta = bottomSize
      }
      topSize += delta
      bottomSize -= delta
      topPane.style.flexGrow = (totalGrow*topSize/totalSize).toString()
      bottomPane.style.flexGrow = (totalGrow*bottomSize/totalSize).toString()
      topPane.style.flexBasis = topSize.toString()+'px'
      bottomPane.style.flexBasis = bottomSize.toString()+'px'
    }

    function onMouseUp(/*mu: MouseEvent*/) {
      window.removeEventListener("mousemove", onMouseMove)
      window.removeEventListener("mouseup", onMouseUp)
    }

    window.addEventListener("mousemove", onMouseMove)
    window.addEventListener("mouseup", onMouseUp)
  }

  render() {
    return html`
      <slot name="top"></slot>
      <div @mousedown="${this.onMouseDown}" id="divider"></div>
      <slot name="bottom"></slot>
    `;
  }
}
