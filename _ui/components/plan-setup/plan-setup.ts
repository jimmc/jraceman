@Polymer.decorators.customElement('plan-setup')
class PlanSetup extends MidTab {

  @Polymer.decorators.property({type: Object, notify: true})
  queryResults: object;
}
