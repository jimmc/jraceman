@Polymer.decorators.customElement('competition-setup')
class CompetitionSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "competition",
    Columns: []         // Columns get set from an API call.
  };
}
