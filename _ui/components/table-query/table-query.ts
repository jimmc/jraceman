interface ColumnDesc {
  Name: string;
  Label: string;
  Type: string;
  FKTable: string;
  FKItems: FKItem[];
}

interface TableDesc {
  Table?: string;
  Columns: ColumnDesc[];
}

interface FKItem {
  ID: string;
  Summary: string;
}

@Polymer.decorators.customElement('table-query')
class TableQuery extends Polymer.Element {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc;

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;

  @Polymer.decorators.property({type: String})
  selectedOp: string;

  clear() {
    for (let col of this.tableDesc['Columns']) {
      const name = col.Name;
      this.$.main.querySelector("#val_"+name).value = '';
      this.$.main.querySelector("#op_"+name).value = 'eq';
    }
  }

  async search() {
    console.log("in TableQuery.search()");
    let params = [];
    for (let col of this.tableDesc['Columns']) {
      const name = col.Name;
      const colVal = this.$.main.querySelector("#val_"+name).value;
      const colOp = this.$.main.querySelector("#op_"+name).value;
      console.log(name, colOp, colVal)
      if (colVal && colOp) {
        const colParams = {
          name: name,
          op: colOp,
          value: colVal,
        };
        params.push(colParams);
      }
    }
    const options: XhrOptions = {
      method: "POST",
      params: params,
    }
    const queryPath = '/api/query/' + this.tableDesc.Table + '/';
    try {
      const result = await ApiManager.xhrJson(queryPath, options);
      if (result && !result.Table) {
        result.Table = this.tableDesc.Table;
      }
      console.log(result);
      this.queryResults = result;
    } catch(e) {
      this.queryResults = {
        Error: e.responseText
      }
    }
  }

  isStringColumn(colType: string) {
    return colType == "string";
  }

  static tableDescToCols(tableDesc: TableDesc): ColumnDesc[] {
    const cols = tableDesc.Columns;
    for (let c=0; c<cols.length; c++) {
      const name = cols[c].Name;
      if (name == 'id') {
        cols[c].Label = name.toUpperCase();
      } else {
        cols[c].Label = name[0].toUpperCase() + name.substr(1);
      }
    }
    return cols;
  }
}
