import { LitElement, html, css } from 'lit'
import { customElement } from 'lit/decorators.js'

import './jraceman-dropdown.js'

import { JracemanDialog } from './jraceman-dialog.js'

// A drop-down menu for debug operations.
@customElement('debug-menu')
export class DebugMenu extends LitElement {
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

  async openDialog() {
    const b = await JracemanDialog.messageDialog("Test Dialog", "Sample message for testing dialog", ["Cancel", "OK"])
    console.log("Button index is", b)
  }

  // Render our menu.
  render() {
    return html`
      <jraceman-dropdown>
        <span slot="control" class="menu">☰</span>
        <div slot="content" class="alignright">
          <button @click="${this.openDialog}">Open Dialog</button>
        </div>
      </jraceman-dropdown>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'debug-menu': DebugMenu;
  }
}
