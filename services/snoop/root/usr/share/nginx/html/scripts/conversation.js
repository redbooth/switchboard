class Conversation {
    constructor(user) {
        this.id = UUID.zero()
        this.user = user;
    }

    create() {
        // TODO: contact server to create conversation
        this.id = new UUID();
    }
}