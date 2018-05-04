class Switchboard {
    constructor(url, conversation, recorder) {
        let self = this;
        this.socket = new WebSocket(url);
        this.socket.addEventListener('open', function() {
            // send switchboard header data
            self.send(conversation.user.id.bytes);
            self.send(conversation.user.token.bytes);
            self.send(conversation.id.bytes);
            self.send(recorder.header);

            // react to socket events
            self.socket.addEventListener('message', self.receive.bind(self));

            // react to sound events
            recorder.addListener(self.send.bind(self));
        });
    }

    send(data) {
        this.socket.send(data);
    }

    receive(evt) {
        console.log('Message received: ', evt);
    }

    close() {
        this.socket.close();
    }
}
