import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';
import {when} from 'lit/directives/when.js';

import { ApiManager, XhrOptions } from './api-manager.js'
import { TableDesc, ColumnDesc, RequestEditEvent } from './table-desc.js'

/**
 * table-edit provides a form to edit one record of a table.
 */
@customElement('table-edit')
export class TableEdit extends LitElement {
  static styles = css`
  `;

  @property({type: Object})
  tableDesc: TableDesc = {
    Table: "(unset-in-table-edit)",
    Columns:[],
  };

  // If we are editing a new record, this value is blank.
  @property({type: String})
  recordId: string = '';

  @property({type: String})
  editIdLabel: string = '[New]';

  // databaseResult holds the result of the database query to load a record.
  databaseResult: any;

  constructor() {
    super()
    // We add a listener for edit requests.
    document.addEventListener("jraceman-request-edit-event", this.onRequestEditEvent.bind(this))
  }

  // This function gets called when someone is asking for a row to be edited.
  async onRequestEditEvent(e:Event) {
    const evt = e as CustomEvent<RequestEditEvent>
    const req:RequestEditEvent = evt.detail
    if (!req || req.Table != this.tableDesc.Table) {
      return;   // Not our table
    }
    if (!req.ID) {
      return;   // No ID specified
    }
    console.log("table-edit edit", req.Table, req.ID);
    // Build a query expression to select that row based on the ID
    const name = this.tableDesc.Columns[0].Name         // Typically "id"
    const colOp = 'eq';
    const colVal = req.ID;
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
      console.error("Error from /api/query:", e/*.responseText*/);
      return
    }
    if (result && !result.Table) {
      result.Table = this.tableDesc.Table;
    }
    console.log(result);
    if (result.Rows.length != 1) {
      throw 'Expected exactly one row';  // TODO - more graceful error handling
    }
    this.databaseResult = result;
    this.reset();
    this.editIdLabel = '[' + (this.recordId || 'New') + ']';
    // Make this tab visible
    this.displayThisTab();
  }

  displayThisTab() {
    const event = new Event('jraceman-request-display-event', {
      bubbles: true,
      composed: true
    })
    // Dispatch the request up our element chain.
    console.log("TableEdit dispatching event", event)
    this.dispatchEvent(event)
  }

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

  // Clear clears a record back to empty.
  clear() {
    for (let col of this.tableDesc['Columns']) {
      const name = col.Name;
      this.setSelectValue("#val_"+name, '')
    }
  }

  // Reset resets an existing record to the values that came from the database.
  reset() {
    // TODO - ask for confirmation if the form data has changed.
    const row = this.databaseResult.Rows[0];
    this.clear();
    // Populate the form
    let c = 0;
    for (let col of this.databaseResult.Columns) {
      const name = col.Name;
      this.setSelectValue("#val_"+name, row[c])
      if (name === this.tableDesc.Columns[0].Name) {
        this.recordId = row[c];
      }
      c++;
    }
  }

  // Sets the form back to blank to allow entering a new record.
  newRecord() {
    // TODO - ask for confirmation if the form data has changed.
    this.clear();
    this.recordId = '';
    this.editIdLabel = '[' + (this.recordId || 'New') + ']';
  }

  async save() {
    console.log("in TableQuery.save()");
    let fields: any = {};
    for (let col of this.tableDesc['Columns']) {
      const name = col.Name;
      const colVal = this.getSelectValue("#val_"+name)
      console.log(name, colVal)
      if (colVal) {
        console.log("Type of field " + name + " is " + col.Type);
        // For non-string fields, convert from the string in the form
        // to the appropriate data type for the field.
        const convertedColVal = this.convertToType(colVal, col.Type)
        // For queries, we specify each field with name and value tags,
        // but when calling the CRUD api we use the name as the field
        // name and the value as the field value for that name.
        fields[name] = convertedColVal;
      }
    }
    const queryPath = '/api/crud/' + this.tableDesc.Table + '/' + this.recordId;
    const method = this.recordId ? "PUT" : "POST";
    const options: XhrOptions = {
      method: method,
      params: fields,
    }
    try {
      const result = await ApiManager.xhrJson(queryPath, options);
      if (result && !result.Table) {
        result.Table = this.tableDesc.Table;
      }
      // Use the returned ID if it was set.
      let returnedId = result['ID']
      if (!returnedId) {
        returnedId = result[this.tableDesc.Columns[0].Name];
      }
      if (returnedId) {
        this.recordId = returnedId;
        this.setSelectValue("#val_id", this.recordId);
        this.editIdLabel = '[' + (this.recordId || 'New') + ']';
      }
      console.log(result);
    } catch(e) {
      console.error("Error:", e/*.responseText*/);
    }
  }

  isStringColumn(colType: string) {
    return colType == "string";
  }

  convertToType(val: string, typ: string): any {
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

  render() {
    return html`
        <form>
          <button type=button @click="${this.save}">Save</button>
          ${when(this.recordId, ()=>html`
            <button type=button @click="${this.reset}">Reset</button>
          `)}
          <button type=button @click="${this.newRecord}">New</button>
          &nbsp;&nbsp; Record: ${this.editIdLabel}

          <table>
            ${repeat(this.tableDesc.Columns, (col /*, colIndex*/) => html`
              <tr>
                <td>${col.Label}</td>
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
        </form>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'table-edit': TableEdit;
  }
}
