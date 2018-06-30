@Polymer.decorators.customElement('progression-table')
class ProgressionTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "progression",
    Columns: []         // Columns get set from an API call.
  };
}
