@Polymer.decorators.customElement('scoringsystem-setup')
class ScoringSystemSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "scoringsystem",
    Columns: []         // Columns get set from an API call.
  };
}
