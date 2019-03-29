@Polymer.decorators.customElement('reports-tab')
class ReportsTab extends Polymer.Element {

  @Polymer.decorators.property({type: Object, notify: true})
  reportList: object;

  @Polymer.decorators.property({type: Object, notify: true})
  reportResults: object;

  ready() {
    super.ready()
    this.loadReportList()
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

  // Generate the report.
  async generate() {
    console.log("Generate");
    const reportName = this.$.main.querySelector("#reportName").value;
    const path = '/api/report/generate/';
    const formData = {
      name: reportName
    };
    const options = {
      method: 'POST',
      params: formData
    };
    try {
      const result = await ApiManager.xhrJson(path, options)
      this.reportResults = result;
    } catch(e) {
      this.reportResults = {
        Error: e.responseText
      }
    }
  }
}
