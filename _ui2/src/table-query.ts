import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';
import {when} from 'lit/directives/when.js';
import '@vaadin/button/vaadin-button.js';

import { ApiManager, XhrOptions } from './api-manager.js'
import { TableDesc } from './table-desc.js'

/**
 * table-query provides a form to do a query on a table.
 */
@customElement('table-query')
export class TableQuery extends LitElement {
  static styles = css`
  `;

  @property({type: Object})
  tableDesc: TableDesc = {
    Table: "(unset-in-table-query)",
    Columns:[],
  };

  @property({type: Object /*, notify: true*/})
  queryResults: object = {};

  @property({type: String})
  selectedOp: string = '';

  // getSelectElement gets an HTMLSelectElement by selector.
  getSelectElement(selector: string) {
    var shadowRoot = this.shadowRoot
    if (shadowRoot == null) {
      console.error("shadowRoot is null")
      return null
    }
    return shadowRoot.querySelector(selector) as HTMLSelectElement
  }

  // getSelectValue gets the value of a <select> element.
  getSelectValue(selector: string) {
    var sel = this.getSelectElement(selector)
    if (sel == null) {
      console.error("select element is null in getSelectValue")
      return null
    }
    return sel.value
  }

  // setSelectValue sets the value of a <select> element.
  setSelectValue(selector: string, val: string) {
    var sel = this.getSelectElement(selector)
    if (sel == null) {
      console.error("select element is null in setSelectValue")
      return
    }
    sel.value = val
  }

  clear() {
    for (let col of this.tableDesc['Columns']) {
      const name = col.Name;
      this.setSelectValue("#val_"+name, '');
      this.setSelectValue("#op_"+name, 'eq');
    }
  }

  async search() {
    console.log("in TableQuery.search()");
    let params = [];
    for (let col of this.tableDesc['Columns']) {
      const name = col.Name;
      const colVal = this.getSelectValue("#val_"+name)
      const colOp = this.getSelectValue("#op_"+name)
      console.log(name, colOp, colVal)
      if (colVal && colOp) {
        const colParams = {
          name: name,
          op: colOp,
          value: colVal,
        };
        params.push(colParams);
      }
    }
    const options: XhrOptions = {
      method: "POST",
      params: params,
    }
    const queryPath = '/api/query/' + this.tableDesc.Table + '/';
    try {
      const result = await ApiManager.xhrJson(queryPath, options);
      if (result && !result.Table) {
        result.Table = this.tableDesc.Table;
      }
      this.queryResults = result;
    } catch(e) {
      this.queryResults = {
        //Error: e.responseText
        Error: e        // TODO: figure out type so we can just return responseText
      }
    }
    console.log("queryResults", this.queryResults);
    // TODO - display queryResults in results tab
  }

  isStringColumn(colType: string) {
    return colType == "string";
  }

  render() {
    console.log("in TableQuery.render tableDesc is", this.tableDesc)
    console.log("this", this)
    return html`
        <form>
          <vaadin-button @click="${this.search.bind(this)}">Search</vaadin-button>
          <vaadin-button @click="${this.clear.bind(this)}">Clear</vaadin-button>
          <table>
            ${repeat(this.tableDesc.Columns, (col /*, colIndex*/) => html`
              <tr>
                <td>${col.Label}</td>
                <td>
                  <select id="op_${col.Name}">
                    <option value="eq">=</option>
                    <option value="ne">!=</option>
                    <template is="dom-if" if="${!col.FKTable}">
                      <option value="gt">&gt;</option>
                      <option value="ge">&gt;=</option>
                      <option value="lt">&lt;</option>
                      <option value="le">&lt;=</option>
                      <option value="like" hidden="${!this.isStringColumn(col.Type)}">LIKE</option>
                    </template>
                  </select>
                </td>
                <td>
                  <!-- $when seems to not quite work: it inserts "false" before rendering
                       its subexpression, and it doesn't render the false leg of a
                       when expresion with both true and false functions. -->
                  ${when(!col.FKTable, ()=>html`
                    <input type=text id="val_${col.Name}" name=${col.Name} label=${col.Name}></input>
                  `)}
                  ${when(col.FKTable, ()=>html`
                    <select id="val_${col.Name}">
                      ${repeat(col.FKItems, (keyitem) => html`
                        <option value="${keyitem.ID}">${keyitem.Summary}</option>
                      `)}
                    </select>
                  `)}
                </td>
              </tr>
              `)}
          </table>
        </form>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'table-query': TableQuery;
  }
}
