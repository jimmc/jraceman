import { LitElement, html, css } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';

import { JracemanApp } from './jraceman-app.js'

/**
 * jraceman-dialog displays an interactive dialog and collects a button press.
 */
@customElement('jraceman-dialog')
export class JracemanDialog extends LitElement {
  static styles = css`
    :host {
      display: block;
      height: 100vh;
      width: 100vw;
      margin: 0;
      padding: 0;
      display: flex;
      flex-direction: column;
      position: absolute;
      opacity: 1;
      pointer-events: none;
    }
    .wrapper {
      display: block;
      height: 100vh;
      width: 100vw;
      margin: 0;
      padding: 0;
      display: flex;
      flex-direction: column;
      position: absolute;
    }
    .wrapper.closed {
      visibility: hidden;
      display: none;
      opacity: 0;
    }
    .wrapper.open {
      align-items: center;
      display: flex;
      justify-content: center;
      height: 100vh;
      position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
      opacity: 1;
      visibility: visible;
      pointer-events: auto;
    }
    .overlay {
      background: rgba(0, 0, 0, 0.8);
      height: 100%;
      position: fixed;
        top: 0;
        right: 0;
        bottom: 0;
        left: 0;
      width: 100%;
    }
    .dialog {
      background: #ffffff;
      max-width: 600px;
      padding: 1rem;
      position: fixed;
    }
    `;

  @property()
  title: string = "TODO: Title here"

  @property()
  message: string = "TOTOD: message here"

  @property()
  buttons: string[] = [ "Cancel" ]

  // True when our dialog is open and being displayed.
  @property()
  open = false

  // Call this method to open the singleton dialog.
  static openDialog(title:string, message:string, buttons:string[]) {
    const app = document.querySelector('jraceman-app')! as JracemanApp
    const dialog = app.shadowRoot!.querySelector('jraceman-dialog')! as JracemanDialog
    console.log("dialog is", dialog)

    dialog.onOpen(title, message, buttons)
  }

  onOpen(title: string, message: string, buttons: string[]) {
    this.title = title
    this.message = message
    this.buttons = buttons
    this.open = true
  }

  // Here when the user clicks a button or outside the dialog.
  // If a button, the index is the array index into this.buttons.
  // If outside the dialog, the index is -1.
  onClose(buttonIndex:number) {
    console.log("JracemanDialog.close", buttonIndex)
    this.open = false
  }

  render() {
    return html`
    <div class="wrapper ${this.open? `open` : `closed`}">
      <div class="overlay" @click="${this.onClose.bind(this, -1)}"></div>
      <div class="dialog" role="dialog">
        <h1 id="title">${this.title}<slot name="heading"></slot></h1>
        <div id="content" class="content">
          ${this.message}
          <slot></slot>
        </div>
        ${repeat(this.buttons, (button:string, buttonIndex)=>html`
          <button type=button class="close" @click=${this.onClose.bind(this, buttonIndex)}>${button}</button>
        `)}
      </div>
    </div>`;

  }
}

declare global {
  interface HTMLElementTagNameMap {
    'jraceman-dialog': JracemanDialog;
  }
}
