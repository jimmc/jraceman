import { LitElement, html, css } from 'lit'
import { customElement } from 'lit/decorators.js'

import './jraceman-tabs.js'
import './table-manager.js'

/**
 * auth-setup is the tab content that contains other tabs for auth setup.
 *
 * @slot - Slots for contains tabs and tab content
 */
@customElement('auth-setup')
export class AuthSetup extends LitElement {
  static styles = css`
  `;

  render() {
    return html`
        <jraceman-tabs>
            <span slot="tab">Users</span>
            <section slot="panel"><table-manager tableName="user"></table-manager></section>
            <span slot="tab">Roles</span>
            <section slot="panel"><table-manager tableName="role"></table-manager></section>
            <span slot="tab">Permissions</span>
            <section slot="panel"><table-manager tableName="permission"></table-manager></section>
            <span slot="tab">User Roles</span>
            <section slot="panel"><table-manager tableName="userrole"></table-manager></section>
            <span slot="tab">Role Permissions</span>
            <section slot="panel"><table-manager tableName="rolepermission"></table-manager></section>
            <span slot="tab">Role Roles</span>
            <section slot="panel"><table-manager tableName="rolerole"></table-manager></section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'auth-setup': AuthSetup;
  }
}
