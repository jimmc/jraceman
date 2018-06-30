@Polymer.decorators.customElement('site-table')
class SiteTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "site",
    Columns: []         // Columns get set from an API call.
  };
}
