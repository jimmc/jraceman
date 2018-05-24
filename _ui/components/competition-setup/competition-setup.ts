@Polymer.decorators.customElement('competition-setup')
class CompetitionSetup extends Polymer.Element {

  @Polymer.decorators.property({type: Number})
  selectedTab: number = 0;

  @Polymer.decorators.property({type: Object})
  tableDesc: object = {
    Table: "competition",
    Columns: [
      {
        Name: "id",
        Label: "ID",
        Type: "string"
      },
      {
        Name: "name",
        Label: "Name",
        Type: "string"
      }
    ]
  };
  // TODO - get the above description from the server
}
