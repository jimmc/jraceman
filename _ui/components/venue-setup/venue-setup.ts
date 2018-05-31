@Polymer.decorators.customElement('venue-setup')
class VenueSetup extends Polymer.Element {

  @Polymer.decorators.property({type: Number})
  selectedTab: number = 0;

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;
}
