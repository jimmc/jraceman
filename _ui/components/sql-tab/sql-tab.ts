@Polymer.decorators.customElement('sql-tab')
class SqlTab extends Polymer.Element {

  checkEnter(e: any) {
    if (e.key == 'Enter' && e.shiftKey) {
      e.stopPropagation();
      this.execute();
    }
  }

  // Clears the SQL text area.
  clear() {
    this.$.sqlText.value = "";
  }

  // Executes the SQL text.
  execute() {
    const sql = this.$.sqlText.value;
    console.log("Execute: " + sql);     // TODO
  }
}
