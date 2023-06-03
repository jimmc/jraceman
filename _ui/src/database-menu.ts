import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';

import './jraceman-dropdown.js'

import { ApiManager, XhrOptions} from './api-manager.js'
import { PostError, PostInfo } from './message-log.js'

// A drop-down menu for database operations.
@customElement('database-menu')
export class DatabaseMenu extends LitElement {
  static styles = css`
    jraceman-dropdown {
      display: inline-block;    /* Make our menu on same line as the tab label */
    }

    .menu {
      cursor: context-menu;
    }
  `;

  onExportClick() {
    PostError("database-menu", "Export Jraceman is not yet implemented")
    console.error("Export JRaceman is not yet implemented");
  }

  onImportClick() {
    PostError("database-menu", "Import JRaceman is not yet implemented");
    console.error("Import JRaceman is not yet implemented");
  }

  onLoadSqlClick() {
    PostError("database-menu", "Load SQL is not yet implemented");
    console.error("Load SQL is not yet implemented");
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
      PostError("database", "Error getting list of sections: " + e/*.responseText*/);
      return;
    }
    const op = (dryrun ? "Upgrade checking " : "Upgrading ");
    PostInfo("database", op + sections.length + " sections")
    let numSectionsToUpgrade = 0;
    for (let i = 0; i < sections.length; i++) {
      const section = sections[i];
      const sectionPath = '/api/database/upgrade/' + section + dryrunFlag;
      const options : XhrOptions = {
        method: 'POST',
      };
      // PostInfo("database", "Updating " + section);
      let sectionResult;
      try {
        sectionResult = await ApiManager.xhrJson(sectionPath, options);
      } catch (e) {
        console.error("section error ", e);
        PostError("database", "Error in section " + section + ": " + e/*.responseText*/);
        if (!dryrun) {
          // If we are not doing a dryrun, stop after an error.
          return;
        } else {
          continue;
        }
      }
      if (sectionResult.Nop) {
        // PostInfo("database", "No change to " + section);
      } else {
        PostInfo("database", resultPrefix + sectionResult.Message);
        numSectionsToUpgrade++;
      }
    }
    if (numSectionsToUpgrade == 0) {
      PostInfo("database", "Upgrade done, no changes");
    } else {
      PostInfo("database", "Upgrade done");
    }
  }

  render() {
    return html`
      <jraceman-dropdown>
        <span slot="control" class="menu">☰</span>
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
