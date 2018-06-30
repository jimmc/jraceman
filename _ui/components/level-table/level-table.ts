@Polymer.decorators.customElement('level-table')
class LevelTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "level",
    Columns: []         // Columns get set from an API call.
  };
}
