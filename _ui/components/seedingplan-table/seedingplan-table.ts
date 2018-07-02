@Polymer.decorators.customElement('seedingplan-table')
class SeedingPlanTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "seedingplan",
    Columns: []         // Columns get set from an API call.
  };
}
