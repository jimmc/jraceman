interface ReportResultsData {
  HTML?: string;
}

@Polymer.decorators.customElement('report-results')
class ReportResults extends Polymer.Element {

  @Polymer.decorators.property({type: Object})
  reportResults: ReportResultsData = {};

  @Polymer.decorators.property({type: String, notify: true})
  reportResultsMoreLabel: string;

  @Polymer.decorators.observe('reportResults')
  reportResultsChanged() {
    this.reportResultsMoreLabel = " (no name yet)"
  }

}
