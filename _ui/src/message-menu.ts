import {LitElement, html, css, render} from 'lit';
import {customElement} from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';
import {when} from 'lit/directives/when.js';

import "./jraceman-dropdown.js"
import { Message, PostMessageEvent } from './message-log.js'

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

  messages: Message[] = []

  connectedCallback() {
    super.connectedCallback()
    document.addEventListener("jraceman-post-message-event", this.onPostMessage.bind(this))
  }

  onPostMessage(e:Event) {
    const evt = e as CustomEvent<PostMessageEvent>
    const m:Message = evt.detail.message
    this.messages.push(m)    // TODO - limit size to a max size
  }

  onViewInNewTab() {
    const newWindow = window.open()!;
    newWindow.document.title = "JRaceman messages"
    render(this.renderMessages(), newWindow.document.body)
  }

  // Render our messages the same way as we do in the message log.
  // This method should be the same as MessageLog.render(),
  // except without any highlighting or listeners,
  // and the style portion of this function should match the
  // static styles variable at the top of the MessageLog class.
  renderMessages() {
    return html`
      <style>
        .error {
          color: darkred;
        }
        .warning {
          color: darkorange;
        }
      </style>
      ${when(this.messages.length==0, ()=>html`(No messages)`)}
      ${repeat(this.messages, (message) => html`
        <span class=${message.level.toLowerCase()}>${message.postTime}: [${message.level}](${message.source}) ${message.text}<span><br/>
      `)}
    `;
  }

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
