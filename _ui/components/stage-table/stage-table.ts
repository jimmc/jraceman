@Polymer.decorators.customElement('stage-table')
class StageTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "stage",
    Columns: []         // Columns get set from an API call.
  };
}
