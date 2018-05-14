class UUID {
    constructor(size = 16) {
        this.size = size
        this.bytes = new Uint8Array(size).map(byte => Math.floor(Math.random() * 256));
    }

    toBuffer() {
        return this.bytes.buffer;
    }

    toArray() {
        return Array.prototype.slice.call(this.bytes);
    }

    toString() {
        return this.toArray().map(byte => byte.toString(16)).join('');
    }

    static zero(size) {
        let id = new UUID(size);
        id.bytes.fill(0);
        return id;
    }
}
