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

// PostError posts an error message event, which we listen for and display.
export function PostError(source: string, text: string) {
  PostMessageParts(source, "Error", text)
}

// PostWarning posts an warning message event, which we listen for and display.
export function PostWarning(source: string, text: string) {
  PostMessageParts(source, "Warning", text)
}

// PostInfo posts an info message event, which we listen for and display.
export function PostInfo(source: string, text: string) {
  PostMessageParts(source, "Info", text)
}

// PostDebug posts an debug message event, which we listen for and display.
export function PostDebug(source: string, text: string) {
  PostMessageParts(source, "Debug", text)
}

// PostMessageParts accepts separate fields describing a message,
// collects them into a Message struct, and calls PostMessage.
// Most application code should call PostError, PostWarning,
// PostInfo, or PostDebug instead.
export function PostMessageParts(source: string, level: string, text: string) {
  const m:Message = {
    source: source,
    level: level,
    text: text,
    postTime: Date().toString().substr(0, 24)        // TODO - better format
  }
  PostMessage(m)
}

// PostMessage accepts a Message struct and dispatches a CustomEvent with that message.
// Application code typically calls one of
// PostError, PostWarning, PostInfo, or PostDebug instead,
// which is turn call PostMessage.
export function PostMessage(m: Message) {
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
  // NOTE: if you change the styles here, you should also update
  // MessageMenu.renderMessages().
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

  // Render our messages.
  // NOTE: if you change the formatting here, you should also update
  // MessageMenu.renderMessages().
  render() {
    return html`
      ${when(this.messages.length==0, ()=>html`(No messages)`)}
      ${repeat(this.messages, (message) => html`
        <span class=${message.level.toLowerCase()}>${message.postTime}: [${message.level}](${message.source}) ${message.text}<span><br/>
      `)}
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'message-log': MessageLog;
  }
}
