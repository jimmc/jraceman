@Polymer.decorators.customElement('gender-setup')
class GenderSetup extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "gender",
    Columns: []         // Columns get set from an API call.
  };
}
