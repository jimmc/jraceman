@Polymer.decorators.customElement('team-table')
class TeamTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "team",
    Columns: []         // Columns get set from an API call.
  };
}
