@Polymer.decorators.customElement('simplan-setup')
class SimplanSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "simplan",
    Columns: []         // Columns get set from an API call.
  };
}
