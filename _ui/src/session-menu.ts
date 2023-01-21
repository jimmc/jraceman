import { LitElement, html, css } from 'lit'
import { customElement } from 'lit/decorators.js'

import './jraceman-dropdown.js'

import { JracemanApp } from './jraceman-app.js'

// A drop-down menu for session operations.
@customElement('session-menu')
export class SessionMenu extends LitElement {
  static styles = css`
    jraceman-dropdown {
      display: inline-block;    /* Make our menu on same line as the tab label */
    }

    .menu {
      cursor: context-menu;
    }

    .alignright {
      right: 0;
    }
  `;

  jracemanApp?: JracemanApp

  setJracemanApp(app: JracemanApp) {
    this.jracemanApp = app
  }

  logout() {
    this.jracemanApp!.logout()
  }

  refreshLogin() {
    this.jracemanApp!.refreshLogin()
  }

  // Render our menu.
  render() {
    return html`
      <jraceman-dropdown>
        <span slot="control" class="menu">☰</span>
        <div slot="content" class="alignright">
          <button @click="${this.refreshLogin}">Refresh login</button>
          <button @click="${this.logout}">Logout</button>
        </div>
      </jraceman-dropdown>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'session-menu': SessionMenu;
  }
}
