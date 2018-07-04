@Polymer.decorators.customElement('registration-table')
class RegistrationTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "registration",
    Columns: []         // Columns get set from an API call.
  };
}
