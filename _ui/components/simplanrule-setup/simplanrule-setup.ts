@Polymer.decorators.customElement('simplanrule-setup')
class SimplanRuleSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "simplanrule",
    Columns: []         // Columns get set from an API call.
  };
}
