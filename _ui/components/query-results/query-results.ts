@Polymer.decorators.customElement('query-results')
class QueryResults extends Polymer.Element {

  @Polymer.decorators.property({type: Object})
  queryResults: object = {};
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
