@Polymer.decorators.customElement('complanrule-table')
class ComplanRuleTable extends QueryEditMidTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc = {
    Table: "complanrule",
    Columns: []         // Columns get set from an API call.
  };
}
