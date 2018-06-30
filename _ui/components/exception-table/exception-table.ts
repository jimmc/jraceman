@Polymer.decorators.customElement('exception-table')
class ExceptionTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "exception",
    Columns: []         // Columns get set from an API call.
  };
}
