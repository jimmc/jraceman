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
class TableEdit extends LeafTab {

  @Polymer.decorators.property({type: Object})
  tableDesc: TableDesc;

  @Polymer.decorators.property({type: Object})
  selectedResult: SelectedResult;

  // If we are editing a new record, this value is blank.
  @Polymer.decorators.property({type: String})
  recordId: string = '';

  @Polymer.decorators.observe('selectedResult')
  async selectedResultChanged() {
    if (!this.selectedResult || this.selectedResult.Table != this.tableDesc.Table) {
      return;   // Not our record
    }
    if (!this.selectedResult.ID) {
      return;   // No ID specified
    }
    console.log("table-edit edit", this.selectedResult.Table, this.selectedResult.ID);
    // Build a query expression to select that row based on the ID
    const name = "id";
    const colOp = 'eq';
    const colVal = this.selectedResult.ID;
    const colParams = {
      name: name,
      op: colOp,
      value: colVal,
    };
    const params = [colParams];
    const options: XhrOptions = {
      method: "POST",
      params: params,
    }
    const queryPath = '/api/query/' + this.tableDesc.Table + '/';
    let result
    try {
      result = await ApiManager.xhrJson(queryPath, options);
    } catch(e) {
      console.error("Error from /api/query:", e.responseText);
      return
    }
    if (result && !result.Table) {
      result.Table = this.tableDesc.Table;
    }
    console.log(result);
    if (result.Rows.length != 1) {
      throw 'Expected exactly one row';  // TODO - more graceful error handling
    }
    const row = result.Rows[0];
    this.clear();
    // Populate the form
    let c = 0;
    for (let col of result.Columns) {
      const name = col.Name;
      this.$.main.querySelector("#val_"+name).value = row[c];
      c++;
    }
    // Make this tab visible
    this.selectThisTab();
  }

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
