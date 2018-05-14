class User {
    constructor() {
        this.id = UUID.zero();
        this.token = UUID.zero();
    }

    login(username, password) {
        // TODO: contact server for real userId and auth token
        this.id = new UUID();
        this.token = new UUID();
    }

    logout() {
        this.id = UUID.zero();
        this.token = UUID.zero();
    }
}
