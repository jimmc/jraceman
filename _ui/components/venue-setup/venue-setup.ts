@Polymer.decorators.customElement('venue-setup')
class VenueSetup extends Polymer.Element {

  @Polymer.decorators.property({type: Number})
  selectedTab: number = 0;

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;

  // requestSelect is when a subtab of ours wants to be selected
  @Polymer.decorators.property({type: Object})
  requestSelect: (tabName: string|null)=>void;

  // selectUp is what we call when we want to be selected
  @Polymer.decorators.property({type: Object})
  selectUp: (tabName: string|null)=>void;

  ready() {
    super.ready();
    this.requestSelect = this.selectChildTab.bind(this);
  }

  selectChildTab(tab: string) {
    if (!tab) {
      console.error("Blank tab requested");
      return;
    }
    const tabElem = this.$.tabs.querySelector('[name="' + tab + '"]');
    if (!tabElem) {
      console.error("No tab named " + tab);
      return;
    }
    tabElem.click();    // select this tab
    // Make this tab visible
    if (this.selectUp) {
      this.selectUp(this.getAttribute("name"));
    }
  }
}
