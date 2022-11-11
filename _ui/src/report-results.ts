import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';

import { ReportResultsData, ReportResultsEvent } from './reports-pane.js'

// Show the results of running a report from report-pane
@customElement('report-results')
export class ReportResults extends LitElement {
  static styles = css`
  `;

  reportResultsHTML: string = ""

  connectedCallback() {
    super.connectedCallback()
    document.addEventListener("jraceman-report-results-event", this.onReportResultsEvent.bind(this))
  }

  getReportResultsHTML() {
    return this.reportResultsHTML
  }

  onReportResultsEvent(e:Event) {
    const evt = e as CustomEvent<ReportResultsEvent>
    console.log("ReportResults got updated results", evt.detail.results)
    const reportResults = evt.detail.results as ReportResultsData
    this.reportResultsHTML = reportResults.HTML
    this.shadowRoot!.querySelector("#results")!.innerHTML = reportResults.HTML
    this.requestUpdate()
  }

  render() {
    return html`
      <span id="results">(Generated reports will appear here)</span>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'report-results': ReportResults;
  }
}
