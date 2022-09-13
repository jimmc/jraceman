import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';

import { ApiManager, XhrOptions} from "./api-manager.js"
import "./jraceman-dropdown.js"
import { PostMessage } from "./message-log.js"

// A drop-down menu for database operations.
@customElement('database-menu')
export class DatabaseMenu extends LitElement {
  static styles = css`
    jraceman-dropdown {
      display: inline-block;    /* Make our menu on same line as the tab label */
    }
  `;

  onExportClick() {
    console.log("Export JRaceman is not yet implemented");
  }

  onImportClick() {
    console.log("Import JRaceman is not yet implemented");
  }

  onLoadSqlClick() {
    console.log("Load SQL is not yet implemented");
  }

  onCheckUpgradeClick() {
    this.doUpgrade(true)
  }

  onUpgradeClick() {
    this.doUpgrade(false)
  }

  async doUpgrade(dryrun: boolean) {
    const dryrunFlag = (dryrun ? '?dryrun=true' : '');
    const resultPrefix = (dryrun ? "Would execute: " : "Executed: ");
    const path = '/api/database/upgrade' + dryrunFlag;
    let sections : string[];
    try {
      sections = await ApiManager.xhrJson(path);
    } catch (e) {
      console.error(e);
      PostMessage("database", "error", "Error getting list of sections: " + e/*.responseText*/);
      return;
    }
    const op = (dryrun ? "Upgrade checking " : "Upgrading ");
    PostMessage("database", "info", op + sections.length + " sections")
    let numSectionsToUpgrade = 0;
    for (let i = 0; i < sections.length; i++) {
      const section = sections[i];
      const sectionPath = '/api/database/upgrade/' + section + dryrunFlag;
      const options : XhrOptions = {
        method: 'POST',
      };
      // PostMessage("database", "info", "Updating " + section);
      let sectionResult;
      try {
        sectionResult = await ApiManager.xhrJson(sectionPath, options);
      } catch (e) {
        console.error("section error ", e);
        PostMessage("database", "error", "Error in section " + section + ": " + e/*.responseText*/);
        if (!dryrun) {
          // If we are not doing a dryrun, stop after an error.
          return;
        } else {
          continue;
        }
      }
      if (sectionResult.Nop) {
        // PostMessage("database", "info", "No change to " + section);
      } else {
        PostMessage("database", "info", resultPrefix + sectionResult.Message);
        numSectionsToUpgrade++;
      }
    }
    if (numSectionsToUpgrade == 0) {
      PostMessage("database", "info", "Upgrade done, no changes");
    } else {
      PostMessage("database", "info", "Upgrade done");
    }
  }

  render() {
    return html`
      <jraceman-dropdown>
        <span slot="control">[M]</span>
        <div slot="content">
          <button @click="${this.onExportClick}">Export JRaceman data file</button>
          <button @click="${this.onImportClick}">Import JRaceman data file</button>
          <button @click="${this.onLoadSqlClick}">Load SQL file</button>
          <button @click="${this.onCheckUpgradeClick}">Check Upgrade</button>
          <button @click="${this.onUpgradeClick}">Upgrade</button>
        </div>
      </jraceman-dropdown>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'database-menu': DatabaseMenu;
  }
}
