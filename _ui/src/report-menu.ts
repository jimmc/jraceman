import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';

import "./jraceman-dropdown.js"
// import { PostError } from "./message-log.js"
import { ReportResultsData, ReportResultsEvent } from "./reports-pane.js"

// A drop-down menu for report operations.
@customElement('report-menu')
export class ReportMenu extends LitElement {
  static styles = css`
    jraceman-dropdown {
      display: inline-block;    /* Make our menu on same line as the tab label */
    }

    .menu {
      cursor: context-menu;
    }
  `;

  reportResultsHTML: string = "(No report results)"

  connectedCallback() {
    super.connectedCallback()
    document.addEventListener("jraceman-report-results-event", this.onReportResultsEvent.bind(this))
  }

  onReportResultsEvent(e:Event) {
    const evt = e as CustomEvent<ReportResultsEvent>
    console.log("ReportMenu got updated report results", evt.detail.results)
    // Save the report results so we can write it out on request.
    const reportResults = evt.detail.results as ReportResultsData
    this.reportResultsHTML = reportResults.HTML
  }

  onOpenInNewTab() {
    const newWindow = window.open()!;
    newWindow.document.title = "JRaceman report"
    newWindow.document.body.innerHTML = this.reportResultsHTML
  }

  onOpenSourceInNewTab() {
    const newWindow = window.open()!;
    newWindow.document.title = "JRaceman report source"
    newWindow.document.body.innerText = this.reportResultsHTML
  }

  onPrint() {
    const newWindow = window.open('','_blank')!;
    newWindow.document.title = "JRaceman report"
    newWindow.document.body.innerHTML = this.reportResultsHTML
    newWindow.print()
    newWindow.close()
  }

  render() {
    return html`
      <jraceman-dropdown>
        <span slot="control" class="menu">☰</span>
        <div slot="content">
          <button @click="${this.onOpenInNewTab}">View in new tab</button>
          <button @click="${this.onOpenSourceInNewTab}">View source in new tab</button>
          <button @click="${this.onPrint}">Print</button>
        </div>
      </jraceman-dropdown>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'report-menu': ReportMenu;
  }
}
