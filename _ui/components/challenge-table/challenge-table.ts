@Polymer.decorators.customElement('challenge-table')
class ChallengeTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "challenge",
    Columns: []         // Columns get set from an API call.
  };
}
