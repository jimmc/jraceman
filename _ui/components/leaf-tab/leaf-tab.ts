class LeafTab extends Polymer.Element {

  @Polymer.decorators.property({type: Object})
  selectUp: (tabName: string|null)=>void;

  selectThisTab() {
    // Make this tab visible
    if (this.selectUp) {
      this.selectUp(this.getAttribute("name"));
    }
  }
}
