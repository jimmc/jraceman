import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import './jraceman-tabs.js'
import './table-queryedit.js'

/**
 * plan-setup is the tab content that contains other tabs for plan setup.
 *
 * @slot - Slots for contains tabs and tab content
 */
@customElement('plan-setup')
export class PlanSetup extends LitElement {
  static styles = css`
  `;

  render() {
    return html`
        <jraceman-tabs>
            <span slot="tab">Simplans</span>
            <section slot="panel"><table-queryedit tableName="simplan"></table-queryedit></section>
            <span slot="tab">Simplan Rules</span>
            <section slot="panel"><table-queryedit tableName="simplanrule"></table-queryedit></section>
            <span slot="tab">Simplen Stages</span>
            <section slot="panel"><table-queryedit tableName="simplanstage"></table-queryedit></section>
            <span slot="tab">Complans</span>
            <section slot="panel"><table-queryedit tableName="complan"></table-queryedit></section>
            <span slot="tab">Complan Rules</span>
            <section slot="panel"><table-queryedit tableName="complanrule"></table-queryedit></section>
            <span slot="tab">Complen Stages</span>
            <section slot="panel"><table-queryedit tableName="complanstage"></table-queryedit></section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'plan-setup': PlanSetup;
  }
}
