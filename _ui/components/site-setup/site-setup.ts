@Polymer.decorators.customElement('site-setup')
class SiteSetup extends Polymer.Element {

  @Polymer.decorators.property({type: Number})
  selectedTab: number = 0;

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "site",
    Columns: []
    // Columns get based on an API call.
  };

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;

  ready() {
    super.ready();
    this.loadColumns();
  }

  async loadColumns() {
    const result = await ApiManager.xhrJson('/api/query/site/')
    const cols = result.Columns;
    for (let c=0; c<cols.length; c++) {
      const name = cols[c].Name;
      if (name == 'id') {
        cols[c].Label = name.toUpperCase();
      } else {
        cols[c].Label = name[0].toUpperCase() + name.substr(1);
      }
    }
    this.set('tableDesc.Columns', cols);
  }
}
