@Polymer.decorators.customElement('site-setup')
class SiteSetup extends MidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "site",
    Columns: []         // Columns get set from an API call.
  };

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;

  @Polymer.decorators.property({type: Object})
  selectedResult: SelectedResult;

  ready() {
    super.ready();
    this.loadColumns();
  }

  async loadColumns() {
    const result: TableDesc = await ApiManager.xhrJson('/api/query/site/')
    const cols = TableQuery.tableDescToCols(result);
    this.set('tableDesc.Columns', cols);
  }
}
