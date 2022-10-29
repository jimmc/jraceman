import {LitElement, html, css, render} from 'lit';
import {customElement} from 'lit/decorators.js';

import './jraceman-dropdown.js'
import { MessageLog } from './message-log.js'

// A drop-down menu for message operations.
@customElement('message-menu')
export class MessageMenu extends LitElement {
  static styles = css`
    jraceman-dropdown {
      display: inline-block;    /* Make our menu on same line as the tab label */
    }

    .menu {
      cursor: context-menu;
    }
  `;

  messageLog?: MessageLog

  setMessageLog(log: MessageLog) {
    this.messageLog = log
  }

  onViewInNewTab() {
    const newWindow = window.open()!;
    newWindow.document.title = "JRaceman messages"
    render(this.renderMessages(), newWindow.document.body)
  }

  // Render our messages the same way as we do in the message log.
  renderMessages() {
    return html`
      <style>
        ${MessageLog.styles}
      </style>
      ${this.messageLog!.render()}
    `;
  }

  // Render our menu.
  render() {
    return html`
      <jraceman-dropdown>
        <span slot="control" class="menu">☰</span>
        <div slot="content">
          <button @click="${this.onViewInNewTab}">View in new tab</button>
        </div>
      </jraceman-dropdown>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'message-menu': MessageMenu;
  }
}
