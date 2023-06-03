import { LitElement, html, css } from 'lit'
import { customElement } from 'lit/decorators.js'

import './jraceman-dropdown.js'

import { ReportResults } from './report-results.js'

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

  reportResultsElement?: ReportResults

  setReportResults(rr: ReportResults) {
    this.reportResultsElement = rr
  }

  onOpenInNewTab() {
    const newWindow = window.open()!;
    newWindow.document.title = "JRaceman report"
    newWindow.document.body.innerHTML = this.reportResultsElement!.getReportResultsHTML()
  }

  onOpenSourceInNewTab() {
    const newWindow = window.open()!;
    newWindow.document.title = "JRaceman report source"
    newWindow.document.body.innerText = this.reportResultsElement!.getReportResultsHTML()
  }

  onPrint() {
    const newWindow = window.open('','_blank')!;
    newWindow.document.title = "JRaceman report"
    newWindow.document.body.innerHTML = this.reportResultsElement!.getReportResultsHTML()
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
