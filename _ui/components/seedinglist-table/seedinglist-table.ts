@Polymer.decorators.customElement('seedinglist-table')
class SeedingListTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "seedinglist",
    Columns: []         // Columns get set from an API call.
  };
}
