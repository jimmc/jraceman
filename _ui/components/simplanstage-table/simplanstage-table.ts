@Polymer.decorators.customElement('simplanstage-table')
class SimplanStageTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "simplanstage",
    Columns: []         // Columns get set from an API call.
  };
}
