@Polymer.decorators.customElement('scoringrule-setup')
class ScoringRuleSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "scoringrule",
    Columns: []         // Columns get set from an API call.
  };
}
