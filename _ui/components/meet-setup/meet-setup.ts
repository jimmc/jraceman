@Polymer.decorators.customElement('meet-setup')
class MeetSetup extends MidTab {

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;
}
