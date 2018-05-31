interface QueryResultsData {
  Table?: string;
  Columns?: object[];
  Rows?: object[][];
}

@Polymer.decorators.customElement('query-results')
class QueryResults extends Polymer.Element {

  @Polymer.decorators.property({type: Object})
  queryResults: QueryResultsData = {};
  /* Sample data looks something like this: {
    Columns: [
      {
        Name: "col1",
        Type: "string"
      },
      {
        Name: "col2",
        Type: "string"
      },
    ],
    Rows: [
      [ "aaa", 123 ],
      [ "bbb", 456 ],
      [ "ccc", 789 ]
    ]
  };
  */

  @Polymer.decorators.property({type: String, notify: true})
  queryResultsMoreLabel: string;

  @Polymer.decorators.observe('queryResults')
  queryResultsChanged() {
    const table = (this.queryResults && this.queryResults.Table) || '';
    const count = (this.queryResults && this.queryResults.Rows && this.queryResults.Rows.length) || 0;
    this.queryResultsMoreLabel = (table || count>0) ? (': ' + table + '[' + count + ']') : '';
  }

  data(row: object[], col: number) {
    if (row === undefined) {
      return undefined;
    } else {
      return row[col];
    }
  }

  colString(col: number) {
    return col.toString();
  }
}
