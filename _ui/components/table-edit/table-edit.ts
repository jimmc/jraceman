interface ColumnDesc {
  Name: string;
  Label: string;
  Type: string;
}

interface TableDesc {
  Table?: string;
  Columns: ColumnDesc[];
}

@Polymer.decorators.customElement('table-edit')
class TableEdit extends Polymer.Element {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc;

  // If we are editing a new record, this value is blank.
  @Polymer.decorators.property({type: String})
  recordId: string = '';

  clear() {
    for (let col of this.tableDesc['Columns']) {
      const name = col.Name;
      this.$.main.querySelector("#val_"+name).value = '';
    }
  }

  async save() {
    console.log("in TableQuery.save()");
    let params = [];
    for (let col of this.tableDesc['Columns']) {
      const name = col.Name;
      const colVal = this.$.main.querySelector("#val_"+name).value;
      console.log(name, colVal)
      if (colVal) {
        const colParams = {
          name: name,
          value: colVal,
        };
        params.push(colParams);
      }
    }
    const options: XhrOptions = {
      method: "POST",
      params: params,
    }
    const queryPath = '/api/edit/' + this.tableDesc.Table + '/';
    /*
    try {
      const result = await ApiManager.xhrJson(queryPath, options);
      if (result && !result.Table) {
        result.Table = this.tableDesc.Table;
      }
      console.log(result);
    } catch(e) {
      // Error: e.responseText
    }
    */
  }

  isNewRecord() {
    return this.recordId == '';
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
