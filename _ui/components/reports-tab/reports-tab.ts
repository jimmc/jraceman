// ReportAttributes is what we get from the API for each report.
interface ReportAttributes {
  Name: string;
  Display: string;
  OrderBy: {[key:string]:string}[];
  Where: {[key:string]:string}[];
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
}

// WhereOption is what we send back to the API for each where field used.
interface WhereOption {
  Name: string;
  Op: string;
  Value: string;
}

@Polymer.decorators.customElement('reports-tab')
class ReportsTab extends Polymer.Element {

  @Polymer.decorators.property({type: Object, notify: true})
  orderByList: OrderByControl[];

  @Polymer.decorators.property({type: Object, notify: true})
  whereList: WhereControl[];

  @Polymer.decorators.property({type: Object, notify: true})
  reportList: ReportAttributes[];

  @Polymer.decorators.property({type: Object, notify: true})
  reportResults: object;

  async ready() {
    super.ready()
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
    const reportName = this.$.main.querySelector("#reportName").value;
    this.updateOrderByList(reportName)
    this.updateWhereList(reportName)
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
    var wcl: WhereControl[] = []
    const reportAttrs = this.findReport(reportName)
    const reportWhereItems: {[key:string]:string}[] = reportAttrs.Where
    if (reportWhereItems) {
      for (let item of reportWhereItems) {
        const whereItem: WhereControl = {
          Name: item["Name"],
          Display: item["Display"],
        }
        wcl.push(whereItem)
      }
    }
    this.whereList = wcl
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
    const reportName = this.$.main.querySelector("#reportName").value
    const orderBy = this.$.main.querySelector("#orderBy").value
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
        Error: e.responseText
      }
    }
  }

  whereOptions() {
    const whereList: WhereOption[] = [];
    for (const item of this.whereList) {
      const name  = item.Name;
      const opFieldTag = "#op_" + name;
      const valFieldTag = "#val_" + name;
      const op = this.$.main.querySelector(opFieldTag).value
      const val = this.$.main.querySelector(valFieldTag).value
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
}
