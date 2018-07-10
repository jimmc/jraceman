@Polymer.decorators.customElement('option-table')
class OptionTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "option",
    Columns: []         // Columns get set from an API call.
  };
}
