import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';
import {when} from 'lit/directives/when.js';

export interface Message {
  source: string;
  level: string;        // one of: Error, Warning, Info, Debug
  text: string;
  postTime: string;
}

export interface PostMessageEvent {
  message: Message;
}

export function PostMessage(source: string, level: string, text: string) {
  const m:Message = {
    source: source,
    level: level,
    text: text,
    postTime: Date().toString().substr(0, 24)        // TODO - better format
  }
  const event = new CustomEvent<PostMessageEvent>('jraceman-post-message-event', {
      detail: {
        message: m,
      } as PostMessageEvent
    });
    // Dispatch the event to the document so any element can listen for it.
    document.dispatchEvent(event);
}

// MessageLog collects and displays log messages.
@customElement('message-log')
export class MessageLog extends LitElement {
  static styles = css`
    .error {
      color: darkred;
    }
    .warning {
      color: darkorange;
    }
  `;

  messages: Message[] = []

  constructor() {
    super()
    document.addEventListener("jraceman-post-message-event", this.onPostMessage.bind(this))
  }

  onPostMessage(e:Event) {
    const evt = e as CustomEvent<PostMessageEvent>
    const m:Message = evt.detail.message
    this.messages.push(m)    // TODO - limit size to a max size
    this.requestUpdate()
  }

  render() {
    return html`
      ${when(this.messages.length==0, ()=>html`(No messages)`)}
      ${repeat(this.messages, (message) => html`
        <span class=${message.level}>${message.postTime}: [${message.level}](${message.source}) ${message.text}<span><br/>
      `)}
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'message-log': MessageLog;
  }
}
