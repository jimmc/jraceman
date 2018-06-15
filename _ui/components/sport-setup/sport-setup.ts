@Polymer.decorators.customElement('sport-setup')
class SportSetup extends MidTab {

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;
}
