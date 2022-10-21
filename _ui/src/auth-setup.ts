import {LitElement, html, css} from 'lit';
import {customElement} from 'lit/decorators.js';
import './jraceman-tabs.js'
import './table-queryedit.js'

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
            <section slot="panel"><table-queryedit tableName="user"></table-queryedit></section>
            <span slot="tab">Roles</span>
            <section slot="panel"><table-queryedit tableName="role"></table-queryedit></section>
            <span slot="tab">Permissions</span>
            <section slot="panel"><table-queryedit tableName="permission"></table-queryedit></section>
            <span slot="tab">User Roles</span>
            <section slot="panel"><table-queryedit tableName="userrole"></table-queryedit></section>
            <span slot="tab">Role Permissions</span>
            <section slot="panel"><table-queryedit tableName="rolepermission"></table-queryedit></section>
            <span slot="tab">Role Roles</span>
            <section slot="panel"><table-queryedit tableName="rolerole"></table-queryedit></section>
        </jraceman-tabs>
    `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'auth-setup': AuthSetup;
  }
}
