@Polymer.decorators.customElement('simplanstage-setup')
class SimplanStageSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "simplanstage",
    Columns: []         // Columns get set from an API call.
  };
}
