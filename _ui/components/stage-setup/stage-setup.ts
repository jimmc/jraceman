@Polymer.decorators.customElement('stage-setup')
class StageSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "stage",
    Columns: []         // Columns get set from an API call.
  };
}
