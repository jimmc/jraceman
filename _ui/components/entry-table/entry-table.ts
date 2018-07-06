@Polymer.decorators.customElement('entry-table')
class EntryTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "entry",
    Columns: []         // Columns get set from an API call.
  };
}
