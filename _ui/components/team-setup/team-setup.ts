@Polymer.decorators.customElement('team-setup')
class TeamSetup extends MidTab {

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;
}
