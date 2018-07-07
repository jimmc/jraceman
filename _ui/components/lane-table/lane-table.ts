@Polymer.decorators.customElement('lane-table')
class LaneTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "lane",
    Columns: []         // Columns get set from an API call.
  };
}
