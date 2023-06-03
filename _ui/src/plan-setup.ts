import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import './jraceman-tabs.js'
import './table-manager.js'

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
            <section slot="panel"><table-manager tableName="simplan"></table-manager></section>
            <span slot="tab">Simplan Rules</span>
            <section slot="panel"><table-manager tableName="simplanrule"></table-manager></section>
            <span slot="tab">Simplen Stages</span>
            <section slot="panel"><table-manager tableName="simplanstage"></table-manager></section>
            <span slot="tab">Complans</span>
            <section slot="panel"><table-manager tableName="complan"></table-manager></section>
            <span slot="tab">Complan Rules</span>
            <section slot="panel"><table-manager tableName="complanrule"></table-manager></section>
            <span slot="tab">Complen Stages</span>
            <section slot="panel"><table-manager tableName="complanstage"></table-manager></section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'plan-setup': PlanSetup;
  }
}
