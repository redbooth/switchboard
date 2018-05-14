class Endian {
    static big(val, size) {
        let bytes = new Uint8Array(size);
        for (let i = 0; i < size; i++) {
            bytes[i] = 0xff & (val >>> (8 * (size - i - 1)))
        }
        return bytes
    }

    static little(val, size) {
        let bytes = new Uint8Array(size);
        for (let i = 0; i < size; i++) {
            bytes[i] = 0xff & (val >>> (8 * i))
        }
        return bytes
    }
}
