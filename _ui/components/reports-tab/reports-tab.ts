// ReportAttributes is what we get from the API for each report.
interface ReportAttributes {
  Name: string;
  Display: string;
  OrderBy: {[key:string]:string}[];
}

// OrderBy is what we prepare for our UI for each OrderBy.
interface OrderBy {
  Name: string;
  Display: string;
}

@Polymer.decorators.customElement('reports-tab')
class ReportsTab extends Polymer.Element {

  @Polymer.decorators.property({type: Object, notify: true})
  orderByList: OrderBy[];

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
  }

  updateOrderByList(reportName: string) {
    var obl: OrderBy[] = []
    const reportAttrs = this.findReport(reportName)
    const reportOrderBys: {[key:string]:string}[] = reportAttrs.OrderBy
    if (reportOrderBys) {
      for (let item of reportOrderBys) {
        const orderby: OrderBy = {
          Name: item["Name"],
          Display: item["Display"],
        }
        obl.push(orderby)
      }
    }
    this.orderByList = obl
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
      orderby: orderBy
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
}
