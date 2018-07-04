@Polymer.decorators.customElement('meet-table')
class MeetTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "meet",
    Columns: []         // Columns get set from an API call.
  };
}
