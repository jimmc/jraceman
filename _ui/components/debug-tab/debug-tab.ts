@Polymer.decorators.customElement('debug-tab')
class DebugTab extends MidTab {

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;
}
