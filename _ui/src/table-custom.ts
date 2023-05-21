// TableCustom provides per-table customizations to the generic code
// in the various table-*.ts files for cases where they require
// information that is not in the table and column descriptions.
export class TableCustom {
  // sheetFilterFieldName returns the name of the column to be used
  // as the one filter field on the table-sheet pane for that table,
  // or an empty string if there should not be a filter field.
  static sheetFilterFieldName(tablename: string): string {
    switch (tablename) {
    case 'entry':
      return 'eventid'
    case 'event':
      return 'meetid'
    case 'person':
      return 'teamid'
    case 'race':
      return 'eventid'
    case 'seedinglist':
      return 'seedingplanid'
    default:
      return ''
    }
  }
}
