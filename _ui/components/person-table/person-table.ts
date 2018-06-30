@Polymer.decorators.customElement('person-table')
class PersonTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "person",
    Columns: []         // Columns get set from an API call.
  };
}
