window.onload = (evt) => {
    let button = document.getElementById('record'),
        canvas = document.getElementById('visualizer'),
        recorder = new Recorder(8000),
        switchboard;

    let handlers = {
        onButtonClick: (evt) => {
            handlers.toggleStyle();
            handlers.toggleRecording();
        },

        toggleStyle: () => {
            button.classList.toggle('recording');
        },

        toggleRecording: () => {
            if (recorder.isRecording()) {
                recorder.stop();
                switchboard.close();
            } else {
                let user = new User();
                user.login();
                let convo = new Conversation(user);
                convo.create();
                switchboard = new Switchboard('ws://localhost:10000', convo, recorder);
                recorder.start();
            }
        },

        visualizeRecording: (analyser, stream) => {
            let WIDTH = canvas.width,
                HEIGHT = canvas.height,
                bufferLength = analyser.frequencyBinCount,
                data = new Uint8Array(bufferLength);

            context = canvas.getContext("2d");
            context.clearRect(0, 0, WIDTH, HEIGHT);

            let draw = function() {
                requestAnimationFrame(draw);
                analyser.getByteFrequencyData(data);
                context.fillStyle = 'rgb(0, 0, 0)';
                context.fillRect(0, 0, WIDTH, HEIGHT);

                let barWidth = (WIDTH / bufferLength) * 2.5;
                    barHeight = 0;
                for(let i = 0, x = 0; i < bufferLength; i++, x += barWidth + 1) {
                    barHeight = data[i];
                    context.fillStyle = `rgb(${ barHeight + 100 },50,50)`;
                    context.fillRect(x, HEIGHT - barHeight / 2, barWidth, barHeight / 2);
                }
            }
            draw();
        }
    };

    button.onclick = handlers.onButtonClick;
}
