@Polymer.decorators.customElement('gender-table')
class GenderTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "gender",
    Columns: []         // Columns get set from an API call.
  };
}
