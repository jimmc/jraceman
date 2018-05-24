interface ColumnDesc {
  Name: string;
  Label: string;
  Type: string;
}

interface TableDesc {
  Table?: string;
  Columns: ColumnDesc[];
}

@Polymer.decorators.customElement('table-query')
class TableQuery extends Polymer.Element {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc;

  clearForm() {
    console.log("in TableQuery.clear()");
  }

  search() {
    console.log("in TableQuery.search()");
    for (let col of this.tableDesc['Columns']) {
      const name = col.Name;
      const formVal = this.$.main.querySelector("#fff"+name).value;
      console.log(name, formVal)
      // TODO - get the operator also, send an API request
    }
  }
}
