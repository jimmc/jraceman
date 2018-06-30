@Polymer.decorators.customElement('laneorder-table')
class LaneOrderTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "laneorder",
    Columns: []         // Columns get set from an API call.
  };
}
