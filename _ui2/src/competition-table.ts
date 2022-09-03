import {css} from 'lit';
import {customElement} from 'lit/decorators.js';

import { TableQueryedit } from './table-queryedit.js'

/**
 * competition-table is the tab content that contains the query and edit tabs
 * for the competition table.
 */
@customElement('competition-table')
export class CompetitionTable extends TableQueryedit {
  static styles = css`
  `;

  constructor() {
    super()
    this.tableName = "competition"
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'competition-table': CompetitionTable;
  }
}
