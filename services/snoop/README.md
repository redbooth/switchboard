Snoop quietly listens to conversations.

### Recording Rates

https://en.wikipedia.org/wiki/Voice_frequency

Human voice frequency ranges from about 85 HZ to 2550Hz. Per the 
Nyquist–Shannon sampling theorem, the sampling frequency must be 
at least twice the highest component of the voice frequency via 
appropriate filtering prior to sampling at discrete times for 
effective reconstruction of the voice signal. Therefore 8000Hz 
should be sufficient to encode the useful information from human
speech.

8-bit [LPCM](https://en.wikipedia.org/wiki/Linear_pulse-code_modulation) allows sufficient granularity to make speech recognizable
(to my ear, anyway). If the algorithms have trouble distinguishing
speakers at this resolution then we can double it to 16-bit.

There is no particular advantage to recording business meetings in
stereo.

8kHz 8-bit mono recordings required 64kbits/s, or 28.8 MB/hour. This
is large but manageable. If the servers or network become overloaded
we can investigate incorporating client-side compression (either lossless
or lossy) into the pipeline.

### Variable Length Streams

Here's an interesting problem: because the WAVE header also specifies the 
total byte length of the audio data in the file but there's no way that 
we can know this ahead of time. Therefore the WAVE header will contain a 
byte-length of 0 initially, which most WAVE decoders will know means to 
just read until EOF. (So says [this NodeJS WAVE library](https://www.npmjs.com/package/wav).)