@Polymer.decorators.customElement('area-setup')
class AreaSetup extends Polymer.Element {

  @Polymer.decorators.property({type: Number})
  selectedTab: number = 0;

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "area",
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
    const result: TableDesc = await ApiManager.xhrJson('/api/query/area/')
    const cols = TableQuery.tableDescToCols(result);
    this.set('tableDesc.Columns', cols);
  }
}
