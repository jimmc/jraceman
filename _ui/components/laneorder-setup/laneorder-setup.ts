@Polymer.decorators.customElement('laneorder-setup')
class LaneOrderSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "laneorder",
    Columns: []         // Columns get set from an API call.
  };
}
