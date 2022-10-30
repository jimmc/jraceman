import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {PropertyValues} from 'lit-element';

import { ApiManager } from "./api-manager.js"
import { PostError } from "./message-log.js"
import { hash as sha256hash } from "./sha256.js"

export interface LoginStateEvent {
  State: boolean;
  Permissions: string;
}

/**
 * jraceman-login is the login screen that shows up when the user is not logged in.
 */
@customElement('jraceman-login')
export class JracemanLogin extends LitElement {
  static styles = css`
      :host {
      }
      #container {
        background-color: white;
        position: absolute;
        width: 400px;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        padding: 15px;
        display: table;
      }
      #container > div {
        display: table-row;
      }
      #container > div > span, input {
        display: table-cell;
      }
      .error {
        color: red;
        margin-top: 10px;
      }
  `;

  @property()
  loginError = ""

  permissions: string[] = [];
  username?: HTMLInputElement
  password?: HTMLInputElement
  focused = ''

  constructor() {
    super()
    ApiManager.AuthErrorCallback = JracemanLogin.AnnounceLoginState
  }

  firstUpdated(changedProperties:PropertyValues<this>) {
    super.firstUpdated(changedProperties);
    this.username = this.shadowRoot!.querySelector("#username") as HTMLInputElement
    this.password = this.shadowRoot!.querySelector("#password") as HTMLInputElement
    this.username!.addEventListener('keydown', this.keydown.bind(this));
    this.password!.addEventListener('keydown', this.keydown.bind(this));
    setTimeout(() => { this.username!.focus(); }, 0);
  }

  // AnnounceLoginState gets our current login state and dispatches a LoginStateEvent.
  public static async AnnounceLoginState() {
    try {
      const statusUrl = "/auth/status/"
      const response = await ApiManager.xhrJson(statusUrl)
      const loggedIn = response.LoggedIn
      const permissions = response.Permissions
      JracemanLogin.SendLoginStateEvent(loggedIn, permissions)
    } catch (e) {
      PostError("login", "Call to /auth/status failed: " + e)
      console.error("auth status call failed", e)
    }
  }

  public static SendLoginStateEvent(loggedIn: boolean, permissions: string) {
    // Dispatch an event so others can take action when the login state changes.
    const event = new CustomEvent<LoginStateEvent>('jraceman-login-state-event', {
      bubbles: true,
      composed: true,
      detail: {
        State: loggedIn,
        Permissions: permissions
      } as LoginStateEvent
    })
    document.dispatchEvent(event)
  }

  async login() {
    const username = this.username!.value;
    const password = this.password!.value;
    const hashword = this.sha256sum(username + "/" + password);
    try {
      const loginUrl = "/auth/login/";
      const formData = new FormData();
      formData.append("username", username);
      formData.append("hashword", hashword);
      const options = {
        method: "POST",
        params: formData,
        encoding: 'direct',
      };
      /* const response = */ await ApiManager.xhrJson(loginUrl, options);
      console.log("Login succeeded")
      this.username!.value = '' // Clear out username and password fields.
      this.password!.value = ''
    } catch (e) {
      this.loginError = "Login failed";
    }
    JracemanLogin.AnnounceLoginState()
  }

  async logout() {
    try {
      const loginUrl = "/auth/logout/";
      /* const response = */ await ApiManager.xhrJson(loginUrl);
      this.permissions = [];
      console.log("Logout succeeded");
    } catch (e) {
      PostError("login", "Logout failed: " + e);
      console.error("Logout failed", e);
    }
    JracemanLogin.AnnounceLoginState()
  }

  hasPermission(perm: string) {
    return (this.permissions.indexOf(perm) >= 0);
  }

  keydown(e: any) {
    this.loginError = "";
    if (e.key == "Enter") {
      if (this.focused == 'username') {
        this.password!.focus();
      } else if (this.focused == 'password') {
        this.login();
      }
    }
  }

  sha256sum(s: string) {
    const s8a = new TextEncoder().encode(s);
    const r8a = sha256hash(s8a);
    const rs = this.toHexString(r8a)
    return rs
  }
  toHexString(bytes:Uint8Array) {
    return Array.prototype.map.call(bytes, (b) => {
      return ('0'+(b & 0xFF).toString(16)).slice(-2)
    }).join('')
  }

  onFocus(fname: string) {
    this.focused = fname
  }

  onBlur() {
    this.focused = ''
  }

  render() {
    return html`
      <div id="container">

        <div>
          <span class=label>Username:</span>
          <input type=text id="username" @focus="${this.onFocus.bind(this,'username')}" @blur="${this.onBlur}"></input>
        </div>
        <div>
          <span class=label>Password:</span>
          <input type=password id="password" @focus="${this.onFocus.bind(this,'password')}" @blur="${this.onBlur}"></input>
        </div>

        <div class="buttons">
          <button type=button raised class="primary" @click="${this.login}">
            Login
          </button>
          <span class="error">${this.loginError}</span>
        </div>

      </div>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'jraceman-login': JracemanLogin;
  }
}
