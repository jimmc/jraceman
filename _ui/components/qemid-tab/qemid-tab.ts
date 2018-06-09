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
  }
}
