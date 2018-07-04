@Polymer.decorators.customElement('registrationfee-table')
class RegistrationFeeTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "registrationfee",
    Columns: []         // Columns get set from an API call.
  };
}
