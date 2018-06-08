/* JRaceman app */

@Polymer.decorators.customElement('jraceman-app')
class JRacemanApp extends Polymer.Element {

  @Polymer.decorators.property({type: Number})
  selectedTopTab: number = 0;

  @Polymer.decorators.property({type: Number})
  selectedBottomTab: number = 0;

  // requestSelect is when a subtab of ours wants to be selected
  @Polymer.decorators.property({type: Object})
  requestSelect: (tabName: string|null)=>void;

  @Polymer.decorators.property({type: Object})
  queryResults: object;

  ready() {
    super.ready();
    ClientMessages.append("top", "JRaceman client started");
    this.requestSelect = this.selectTopTab.bind(this);
  }

  // When the query results get changed, display them.
  @Polymer.decorators.observe('queryResults')
  queryResultsChanged() {
    this.selectedBottomTab = 1;
  }

  selectTopTab(tab: string) {
    if (!tab) {
      console.error("Blank tab requested");
      return;
    }
    const tabElem = this.$.toptabs.querySelector('[name="' + tab + '"]');
    if (!tabElem) {
      console.error("No tab named " + tab);
      return;
    }
    tabElem.click();    // select this tab
  }
}
