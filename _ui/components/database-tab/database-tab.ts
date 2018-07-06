@Polymer.decorators.customElement('database-tab')
class DatabaseTab extends MidTab {

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;
}
