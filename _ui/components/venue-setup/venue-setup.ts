@Polymer.decorators.customElement('venue-setup')
class VenueSetup extends MidTab {

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;
}
