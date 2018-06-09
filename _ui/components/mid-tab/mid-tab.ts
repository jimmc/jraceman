// MidTab is a base class for any tab which has child tabs.
class MidTab extends Polymer.Element {

  // selectedTab is the index of our child tab
  @Polymer.decorators.property({type: Number})
  selectedTab: number = 0;

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
    tabElem.click();    // select the requested child tab
    // Make this tab visible
    this.selectThisTab();
  }

  selectThisTab() {
    if (this.selectUp) {
      this.selectUp(this.getAttribute("name"));
    }
  }
}
