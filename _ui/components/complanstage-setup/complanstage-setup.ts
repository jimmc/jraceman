@Polymer.decorators.customElement('complanstage-setup')
class ComplanStageSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "complanstage",
    Columns: []         // Columns get set from an API call.
  };
}
