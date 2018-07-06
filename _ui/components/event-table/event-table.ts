@Polymer.decorators.customElement('event-table')
class EventTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "event",
    Columns: []         // Columns get set from an API call.
  };
}
