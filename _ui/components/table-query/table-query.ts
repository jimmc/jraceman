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
      const colVal = this.$.main.querySelector("#val_"+name).value;
      const opItem = this.$.main.querySelector("#op_"+name).selectedItem;
      const colOp = opItem && opItem.getAttribute('name');
      console.log(name, colOp, colVal)
    }
    // TODO - send an API request
  }

  isStringColumn(colType: string) {
    return colType == "string";
  }
}
