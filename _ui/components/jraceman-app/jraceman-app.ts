/* JRaceman app */

@Polymer.decorators.customElement('jraceman-app')
class JRacemanApp extends Polymer.Element {

  @Polymer.decorators.property({type: Number})
  selectedTopTab: number = 0;

  @Polymer.decorators.property({type: Number})
  selectedBottomTab: number = 0;

  @Polymer.decorators.property({type: Object})
  queryResults: object;

  ready() {
    super.ready();
    ClientMessages.append("top", "JRaceman client started");
  }

  // When the query results get changed, display them.
  @Polymer.decorators.observe('queryResults')
  queryResultsChanged() {
    this.selectedBottomTab = 1;
  }
}
