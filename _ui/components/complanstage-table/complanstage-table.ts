@Polymer.decorators.customElement('complanstage-table')
class ComplanStageTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "complanstage",
    Columns: []         // Columns get set from an API call.
  };
}
