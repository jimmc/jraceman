@Polymer.decorators.customElement('contextoption-table')
class ContextOptionTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "contextoption",
    Columns: []         // Columns get set from an API call.
  };
}
