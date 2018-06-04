interface QueryResultsData {
  Table?: string;
  Columns?: object[];
  Rows?: object[][];
}

interface SelectedResult {
  Table: string;
  ID: string;
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

  @Polymer.decorators.property({type: Object, notify: true})
  activeItem: object;

  @Polymer.decorators.property({type: Object, notify: true})
  selectedResult: object;

  @Polymer.decorators.observe('queryResults')
  queryResultsChanged() {
    const table = (this.queryResults && this.queryResults.Table) || '';
    const count = (this.queryResults && this.queryResults.Rows && this.queryResults.Rows.length) || 0;
    this.queryResultsMoreLabel = (table || count>0) ? (': ' + table + '[' + count + ']') : '';
  }

  @Polymer.decorators.observe('activeItem')
  activeItemChanged(item: object) {
    // Clicking on a row in vaadin-grid sets it as activeItem.
    // This code make that become the selected item, which highlights it.
    // But this works independently of double-click, which looks weird,
    // so for now we will comment it out, so the row doesn't highlight.
    // this.$.grid.selectedItems = item ? [item] : [];
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

  dblclick(e: any) {
    if (!this.queryResults.Table) {
      // No table, can't edit
      return;
    }
    const row = e.model.item;
    const id = row[0];  // Assume the ID is the first field of the results
    // Let obsevers of selectedResult changes handle the action.
    this.selectedResult = {
      Table: this.queryResults.Table,
      ID: id
    }
  }
}
