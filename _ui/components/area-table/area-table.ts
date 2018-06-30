@Polymer.decorators.customElement('area-table')
class AreaTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "area",
    Columns: []         // Columns get set from an API call.
  };
}
