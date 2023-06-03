export interface ColumnDesc {
  Name: string;
  Label: string;
  Type: string;
  FKTable: string;
  FKItems: FKItem[];
  ReadOnly?: boolean;
  Hidden?: boolean;
}

export interface TableDesc {
  Table: string;
  Columns: ColumnDesc[];
}

export interface FKItem {
  ID: string;
  Summary: string;
}

export interface QueryResultsColumnDesc {
  Name: string;
  Type: string;
}

export interface QueryResultsData {
  Error?: string;
  Table: string;
  Columns: QueryResultsColumnDesc[];
  Rows: any[][];
}

export interface QueryResultsEvent {
  message: string;
  results: QueryResultsData;
}

export interface RequestEditEvent {
  Table: string;
  ID: string;
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

  // Convert a value from a string representation to its actual value.
  static convertToType(val: string, typ: string): any {
    switch (typ) {
    case 'bool':
      return val.toLowerCase()=='true' || val=='1';  // TODO - explicitly check for false values?
    case 'int':
      return parseInt(val);     // TODO - catch and handle parsing errors
    case 'float':
    case 'float32':
      return parseFloat(val);   // TODO - catch and handle parsing errors
    default:
      return val;       // no conversion for strings or unknown types
    }
  }

  // emptyQueryResuts returns an empty QueryResultsData for use in initializers or
  // other places where we don't have any results data.
  static emptyQueryResults(): QueryResultsData {
    const empty: QueryResultsData = {
      Table: '(unset in table-desc)',
      Columns: [],
      Rows: [],
    }
    return empty
  }
}
