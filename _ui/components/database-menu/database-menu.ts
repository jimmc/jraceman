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
    ClientMessages.append("database", "Check Upgrade is not yet implemented");
  }

  upgradeClicked() {
    ClientMessages.append("database", "Upgrade is not yet implemented");
  }
}
