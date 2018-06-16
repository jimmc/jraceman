@Polymer.decorators.customElement('level-setup')
class LevelSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "level",
    Columns: []         // Columns get set from an API call.
  };
}
