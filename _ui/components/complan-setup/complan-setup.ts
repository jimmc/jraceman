@Polymer.decorators.customElement('complan-setup')
class ComplanSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "complan",
    Columns: []         // Columns get set from an API call.
  };
}
