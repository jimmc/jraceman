@Polymer.decorators.customElement('progression-setup')
class ProgressionSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "progression",
    Columns: []         // Columns get set from an API call.
  };
}
