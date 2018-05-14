class Recorder {
    constructor(rate) {
        // TODO: multiple channels and depths
        // NOTE: google cloud only accepts FLAC and LCM16, so for 16 bits is it
        this.channels = 1;          // Mono = 1, Stereo = 2, etc.
        this.rate = rate;           // sample rate in hertz
        this.depth = 16;            // bits per sample
        this.listeners = new Set();
        this.context = null;
        this._offset = 0;
    }

    // http://soundfile.sapp.org/doc/WaveFormat/
    get header() {
        let header = new Uint8Array(44);
        let i = 0;
        [
            // "RIFF" chunk descriptor
            new Uint8Array([0x52, 0x49, 0x46, 0x46]),                       // ChunkID: "RIFF"
            Endian.little(36, 4),                                           // ChunkSize: assumes 0 for data size
            new Uint8Array([0x57, 0x41, 0x56, 0x45]),                       // Format: "WAVE"
            // "fmt " sub-chunk
            new Uint8Array([0x66, 0x6d, 0x74, 0x20]),                       // Subchunk1ID: "fmt "
            Endian.little(16, 4),                                           // Subchunk1Size: PCM = 16
            Endian.little(1, 2),                                            // AudioFormat: PCM = 1
            Endian.little(this.channels, 2),                                // NumChannels: Mono = 1, Stereo = 2, etc.
            Endian.little(this.rate, 4),                                    // SampleRate
            Endian.little(this.rate * this.channels * this.depth / 8, 4),   // ByteRate
            Endian.little(this.channels * this.depth / 8, 2),               // BlockAlign
            Endian.little(this.depth, 2),                                   // BitsPerSample
            // "data" sub-chunk
            new Uint8Array([0x64, 0x61, 0x74, 0x61]),                       // Subchunk2ID: "data"
            Endian.little(0, 4)                                             // most WAVE decoders know that 0 means to read until EOF
        ].forEach(outer => outer.forEach(inner => header[i++] = inner))
        return header;
    }

    addListener(listener) {
        return this.listeners.add(listener)
    }

    removeListener(listener) {
        return this.listeners.delete(listener);
    }

    process(evt) {
        let max = (1 << (this.depth - 1)) - 1,
            step = evt.inputBuffer.sampleRate / this.rate,
            input = evt.inputBuffer.getChannelData(0),
            output = new Int16Array(Math.floor((input.length - this._offset) / step));

        // resample input data
        let i = 0;
        while (this._offset < input.length - 1) {
            // linear interpolation
            let high = input[Math.ceil(this._offset)],
                low = input[Math.floor(this._offset)],
                mid = low + (high - low) * (this._offset % 1);

            // convert floats in [-1.0, 1.0] to ints in [-max, max]
            output[i] = Math.round(max * mid);

            // increment input and output counters
            this._offset += step;
            i += 1;
        }
        this._offset %= (input.length - 1);

        // send processed data to callback listeners
        this.listeners.forEach(listener => listener(output));
    }

    start() {
        if (this.isRecording()) return;

        let self = this;
        this.context = new AudioContext();

        // create stream
        navigator.mediaDevices.getUserMedia({
            audio: true
        }).then(stream => {
            console.log('Start recording...');
            let source = self.context.createMediaStreamSource(stream);

            let processor = self.context.createScriptProcessor(4096, 1, 1);
            processor.onaudioprocess = self.process.bind(self);

            source.connect(processor).connect(self.context.destination);
        });
    }

    stop() {
        if (this.context) {
            console.log('Stop recording...');
            this.context.close();
            this.context = null;
        }
    }

    isRecording() {
        return this.context !== null && this.context.state === 'running';
    }
}
