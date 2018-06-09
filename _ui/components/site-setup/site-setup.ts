@Polymer.decorators.customElement('site-setup')
class SiteSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "site",
    Columns: []         // Columns get set from an API call.
  };
}
