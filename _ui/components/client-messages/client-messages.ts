@Polymer.decorators.customElement('client-messages')
class ClientMessages extends Polymer.Element {

  static singleton: ClientMessages;

  @Polymer.decorators.property({type: String})
  messagetext: string = '';

  // Adds a line to our message area.
  static append(source: string, message: string) {
    if (ClientMessages.singleton) {
      ClientMessages.singleton.appendMessage(source, message)
    } else {
      console.log('No message area for message: ' + source + ': ' + message);
    }
  }

  appendMessage(source: string, message: string) {
    const now = Date().toString().substr(0, 24)        // TODO - better format
    const newMessageLine = now + ' ' + source + ': ' + message
    // TODO - keep a list of messages and keep only N so we don't grow without bound
    this.messagetext = this.messagetext + newMessageLine + '\n';
  }

  ready() {
    super.ready();
    ClientMessages.singleton = this;
  }
}
