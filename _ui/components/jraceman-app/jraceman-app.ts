/* JRaceman app */

@Polymer.decorators.customElement('jraceman-app')
class JRacemanApp extends Polymer.Element {

  @Polymer.decorators.property({type: Number})
  selectedTopTab: number = 0;

  @Polymer.decorators.property({type: Number})
  selectedBottomTab: number = 0;
}
