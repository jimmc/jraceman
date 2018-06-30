@Polymer.decorators.customElement('complan-table')
class ComplanTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "complan",
    Columns: []         // Columns get set from an API call.
  };
}
