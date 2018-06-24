@Polymer.decorators.customElement('database-menu')
class DatabaseMenu extends Polymer.Element {

  exportClicked() {
    ClientMessages.append("database", "Export JRaceman is not yet implemented");
  }

  importClicked() {
    ClientMessages.append("database", "Import JRaceman is not yet implemented");
  }

  loadSqlClicked() {
    ClientMessages.append("database", "Load SQL is not yet implemented");
  }

  checkUpgradeClicked() {
    this.doUpgrade(true)
  }

  upgradeClicked() {
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
      ClientMessages.append("database", "Error getting list of sections: " + e.responseText);
      return;
    }
    const op = (dryrun ? "Upgrade checking " : "Upgrading ");
    ClientMessages.append("database", op + sections.length + " sections")
    let numSectionsToUpgrade = 0;
    for (let i = 0; i < sections.length; i++) {
      const section = sections[i];
      const sectionPath = '/api/database/upgrade/' + section + dryrunFlag;
      const options : XhrOptions = {
        method: 'POST',
      };
      // ClientMessages.append("database", "Updating " + section);
      let sectionResult;
      try {
        sectionResult = await ApiManager.xhrJson(sectionPath, options);
      } catch (e) {
        console.error("section error ", e);
        ClientMessages.append("database", "Error in section " + section + ": " + e.responseText);
        if (!dryrun) {
          // If we are not doing a dryrun, stop after an error.
          return;
        } else {
          continue;
        }
      }
      if (sectionResult.Nop) {
        // ClientMessages.append("database", "No change to " + section);
      } else {
        ClientMessages.append("database", resultPrefix + sectionResult.Message);
        numSectionsToUpgrade++;
      }
    }
    if (numSectionsToUpgrade == 0) {
      ClientMessages.append("database", "Upgrade done, no changes");
    } else {
      ClientMessages.append("database", "Upgrade done");
    }
  }
}
