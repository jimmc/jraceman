// QueryEditMidTab is a base class for middle tabs which have both
// Query and Edit child tabs.
class QueryEditMidTab extends MidTab {

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;

  @Polymer.decorators.property({type: Object})
  selectedResult: SelectedResult;

  ready() {
    super.ready();
    this.loadColumns();
  }

  // Loads our column info from the API and sets the result into this.tableDesc.Columns.
  // Assumes that this.tableDesc exists and has a Table field.
  async loadColumns() {
    const tableName: string = this.get('tableDesc.Table');
    const result: TableDesc = await ApiManager.xhrJson('/api/query/' + tableName + '/')
    const cols = TableQuery.tableDescToCols(result);
    this.set('tableDesc.Columns', cols);
    for (var i = 0; i<cols.length; i++) {
      const col = cols[i]
      if (col.FKTable) {
        this.set('tableDesc.Columns.'+i+'.FKItems', [{ID: "", Summary: ""}]);
        this.loadFKChoices(i, col.FKTable)
      }
    }
  }

  async loadFKChoices(i: number, table: string) {
    console.log("In loadFKChoices for", table)
    const path = '/api/query/' + table + "/summary/"
    const options = {}
    try {
      const result = await ApiManager.xhrJson(path, options)
      const newFKItems: FKItem[] = [];
      newFKItems.push({ID: "", Summary: ""});
      for (const row of result.Rows) {
        newFKItems.push({ID: row[0], Summary: row[1]});
      }
      this.set('tableDesc.Columns.'+i+'.FKItems', newFKItems)
    } catch(e) {
      console.log("Error: ", e)         // TODO
    }
  }
}
