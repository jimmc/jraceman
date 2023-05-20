import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';
import {when} from 'lit/directives/when.js';

import { TableDesc } from './table-desc.js'

/**
 * query-fields provides a set of input fields to collect query parameters.
 */
@customElement('query-fields')
export class QueryFields extends LitElement {
  static styles = css`
    .inline {
      display: inline;
    }
    .inline tr {
      display: inline;
    }
  `;

  @property({type: Object})
  tableDesc: TableDesc = {
    Table: "(unset-in-query-fields)",
    Columns: [],
  };

  // tableClass is an optional property that allows selecting a different
  // set of CSS properties to display the table differently.
  @property({type: String})
  tableClass: string = ''

  // getSelectElement gets an HTMLSelectElement by selector.
  getSelectElement(selector: string) {
    const shadowRoot = this.shadowRoot
    if (shadowRoot == null) {
      console.error("shadowRoot is null")
      return null
    }
    return shadowRoot.querySelector(selector) as HTMLSelectElement
  }

  // getSelectValue gets the value of a <select> element.
  getSelectValue(selector: string) {
    const sel = this.getSelectElement(selector)
    if (sel == null) {
      console.error("select element is null in getSelectValue")
      return null
    }
    return sel.value
  }

  // setSelectValue sets the value of a <select> element.
  setSelectValue(selector: string, val: string) {
    const sel = this.getSelectElement(selector)
    if (sel == null) {
      console.error("select element is null in setSelectValue")
      return
    }
    sel.value = val
  }

  clear() {
    for (let col of this.tableDesc.Columns) {
      const name = col.Name;
      this.setSelectValue("#val_"+name, '');
      this.setSelectValue("#op_"+name, 'eq');
    }
  }

  fieldsAsParams() {
    let params = [];
    for (let col of this.tableDesc.Columns) {
      const name = col.Name;
      const colVal = this.getSelectValue("#val_"+name)
      const colOp = this.getSelectValue("#op_"+name)
      console.log("QueryFields.search",name, colOp, colVal)
      if (colVal && colOp) {
        const colParams = {
          name: name,
          op: colOp,
          value: colVal,
        };
        params.push(colParams);
      }
    }
    return params
  }

  isStringColumn(colType: string) {
    return colType == "string";
  }

  render() {
    return html`
      <table class=${this.tableClass}>
        ${repeat(this.tableDesc.Columns, (col /*, colIndex*/) => html`
          <tr>
            <td>${col.Label}</td>
            <td>
              <select id="op_${col.Name}">
                <option value="eq">=</option>
                <option value="ne">!=</option>
                ${when(!col.FKTable, ()=>html`
                  <option value="gt">&gt;</option>
                  <option value="ge">&gt;=</option>
                  <option value="lt">&lt;</option>
                  <option value="le">&lt;=</option>
                  ${when(this.isStringColumn(col.Type), ()=>html`
                    <option value="like">LIKE</option>
                  `)}
                `)}
              </select>
            </td>
            <td>
              ${when(col.FKTable, ()=>html`
                <select id="val_${col.Name}">
                  ${repeat(col.FKItems, (keyitem) => html`
                    <option value="${keyitem.ID}">${keyitem.Summary}</option>
                  `)}
                </select>
              `, ()=>html`
                <input type=text id="val_${col.Name}" name=${col.Name} label=${col.Name}></input>
              `)}
            </td>
          </tr>
          `)}
      </table>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'query-fields': QueryFields;
  }
}
