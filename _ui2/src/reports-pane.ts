import {LitElement, html, css} from 'lit';
import {customElement, property} from 'lit/decorators.js';
import {repeat} from 'lit/directives/repeat.js';
import {when} from 'lit/directives/when.js';

import { ApiManager } from './api-manager.js'

// ReportAttributes is what we get from the API for each report.
interface ReportAttributes {
  Name: string;
  Display: string;
  OrderBy: {[key:string]:string}[];
  Where: ReportWhereControl[];
}

// ReportWhereControl is we get from the API for each where item.
interface ReportWhereControl {
  Name: string;
  Display: string;
  Ops: string[];
  KeyTable: string;
}

// OrderByControl is what we prepare for our UI for each OrderBy.
interface OrderByControl {
  Name: string;
  Display: string;
}

// WhereControl is what we prepare for our UI for each Where item.
interface WhereControl {
  Name: string;
  Display: string;
  Ops: WhereControlOp[];
  KeyTable: string;
  KeyItems: KeyTableItem[];
}

// WhereControlOp is what we prepare for our UI for each Op on a Where item.
interface WhereControlOp {
  Name: string;
  Display: string;
}

// WhereOption is what we send back to the API for each where field used.
interface WhereOption {
  Name: string;
  Op: string;
  Value: string;
}

// KeyTableItem is what we prepare for the UI for a choice list on a key field.
interface KeyTableItem {
  ID: string;
  Summary: string;
}

const opDisplayMap: {[key:string]:string} = {
  "eq": "=",
  "ne": "!=",
  "gt": ">",
  "ge": ">=",
  "lt": "<",
  "le": "<=",
  "like": "LIKE",
}

// A pane for selecting and generating reports.
@customElement('reports-pane')
export class ReportsPane extends LitElement {
  static styles = css`
    :host {
      height: 100%;
    }
    #main {
      display: flex;
      flex-direction: row;
      height: 100%;
    }
    #left {
      flex: 0 0 30%
    }
    #right {
      flex: 1
    }
    #buttons {
      flex: none;
    }
  `;

  @property({type: Object})
  orderByList: OrderByControl[] = []

  @property({type: Object})
  whereList: WhereControl[] = []

  @property({type: Object})
  reportList: ReportAttributes[] = []

  @property({type: Object})
  reportResults: object = {}

  async connectedCallback() {
    super.connectedCallback()
    await this.loadReportList()
    this.reportNameChanged()
  }

  async loadReportList() {
    console.log("In loadReportList")
    const path = '/api/report/'
    const options = {}
    try {
      const result = await ApiManager.xhrJson(path, options)
      this.reportList = result
    } catch(e) {
      console.log("Error: ", e)         // TODO
    }
  }

  reportNameChanged() {
    const reportName = (this.shadowRoot!.querySelector("#reportName") as HTMLSelectElement)!.value;
    this.updateOrderByList(reportName)
    this.updateWhereList(reportName)
    this.updateKeyChoices()
  }

  updateOrderByList(reportName: string) {
    var obl: OrderByControl[] = []
    const reportAttrs = this.findReport(reportName)
    const reportOrderBys: {[key:string]:string}[] = reportAttrs.OrderBy
    if (reportOrderBys) {
      for (let item of reportOrderBys) {
        const orderby: OrderByControl = {
          Name: item["Name"],
          Display: item["Display"],
        }
        obl.push(orderby)
      }
    }
    this.orderByList = obl
  }

  updateWhereList(reportName: string) {
    const reportAttrs = this.findReport(reportName)
    const reportWhereItems: ReportWhereControl[] = reportAttrs.Where
    if (!reportWhereItems) {
      this.whereList = []
      return
    }
    var wcl: WhereControl[] = []
    for (let item of reportWhereItems) {
      var ops: WhereControlOp[] = []
      for (const opName of item.Ops) {
        const opItem: WhereControlOp = {
          Name: opName,
          Display: opDisplayMap[opName],
        }
        ops.push(opItem);
      }
      const whereItem: WhereControl = {
        Name: item.Name,
        Display: item.Display,
        Ops: ops,
        KeyTable: item.KeyTable,
        KeyItems: [{ID:"", Summary:"(no items available)"}],
      }
      wcl.push(whereItem)
    }
    this.whereList = wcl
  }

  // updateKeyChoices goes through the updated whereList looking for fields that have a KeyTable set.
  // For each of those, it issues a summary request for that table and stores the results in TODO.
  updateKeyChoices() {
    for (var i = 0; i<this.whereList.length; i++) {
      const item = this.whereList[i]
      if (item.KeyTable) {
        this.loadKeyChoices(i, item.KeyTable)
      }
    }
  }

  async loadKeyChoices(i: number, table: string) {
    console.log("In loadKeyChoices for", table)
    const path = '/api/query/' + table + "/summary/"
    const options = {}
    try {
      const result = await ApiManager.xhrJson(path, options)
      const newKeyItems: KeyTableItem[] = [];
      newKeyItems.push({ID: "", Summary: ""});
      for (const row of result.Rows) {
        newKeyItems.push({ID: row[0], Summary: row[1]});
      }
      this.whereList[i].KeyItems = newKeyItems
      this.requestUpdate()
    } catch(e) {
      console.log("Error: ", e)         // TODO
    }
  }

  findReport(reportName: string): ReportAttributes {
    for (const report of this.reportList) {
      if (report.Name == reportName) {
        return report
      }
    }
    throw("Report not found: " + reportName)
  }

  // Generate the report.
  async generate() {
    console.log("Generate")
    const reportName = (this.shadowRoot!.querySelector("#reportName") as HTMLSelectElement)!.value
    const orderBy = (this.shadowRoot!.querySelector("#orderBy") as HTMLSelectElement)!.value
    const path = '/api/report/generate/'
    const formData = {
      name: reportName,
      orderby: orderBy,
      where: this.whereOptions()
    }
    const options = {
      method: 'POST',
      params: formData
    }
    try {
      const result = await ApiManager.xhrJson(path, options)
      this.reportResults = result
    } catch(e) {
      this.reportResults = {
        Error: e//.responseText         // TODO - get repsonseText only
      }
    }
  }

  whereOptions() {
    const whereList: WhereOption[] = [];
    for (const item of this.whereList) {
      const name  = item.Name;
      const opFieldTag = "#op_" + name;
      const valFieldTag = "#val_" + name;
      const op = (this.shadowRoot!.querySelector(opFieldTag) as HTMLSelectElement)!.value
      const val = (this.shadowRoot!.querySelector(valFieldTag) as HTMLSelectElement)!.value
      if (val != '') {
        const itemOption: WhereOption = {
          Name: name,
          Op: op,
          Value: val
        };
        whereList.push(itemOption)
      }
    }
    return whereList
  }

  render() {
    return html`
      <div id="main">
        <div id="left">
          <div id="buttons">
            <button @click="${this.generate}">Generate</button>
          </div>
          <div>
            Available reports:
            <select id="reportName" on-change="reportNameChanged">
              ${repeat(this.reportList, (report)=>html`
                <option value="${report.Name}">${report.Display}</option>
              `)}
            </select>
            <br>
            Order By:
            <select id="orderBy">
              ${repeat(this.orderByList, (item)=>html`
                <option value="${item.Name}">${item.Display}</option>
              `)}
            </select>
          </div>
        </div>  <!-- left -->
        <div id="right">
          Select:
          <br>
          <table>
            ${repeat(this.whereList, (item)=>html`
              <tr>
                <td>${item.Display}</td>
                <td>
                  <select id="op_${item.Name}">
                    ${repeat(item.Ops, (op)=>html`
                      <option value="${op.Name}">${op.Display}</option>
                    `)}
                  </select>
                </td>
                <td>
                  ${when(item.KeyTable, ()=>html`
                    <select id="val_${item.Name}">
                      ${repeat(item.KeyItems, (keyitem)=>html`
                        <option value="${keyitem.ID}">${keyitem.Summary}</option>
                      `)}
                    </select>
                  `, ()=>html`
                    <input type=text id="val_${item.Name}" name=${item.Name} label=${item.Name}></input>
                  `)}
                </td>
              </tr>
            `)}
          </table>
        </div>  <!-- right -->
      </div>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'reports-pane': ReportsPane;
  }
}
