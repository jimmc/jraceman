export interface ColumnDesc {
  Name: string;
  Label: string;
  Type: string;
  FKTable: string;
  FKItems: FKItem[];
}

export interface TableDesc {
  Table: string;
  Columns: ColumnDesc[];
}

export interface FKItem {
  ID: string;
  Summary: string;
}

export class TableDescSupport {
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
