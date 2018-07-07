@Polymer.decorators.customElement('race-table')
class RaceTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "race",
    Columns: []         // Columns get set from an API call.
  };
}
