let audio = new AudioContext()
masterGain = audio.createGain()
masterGain.connect(audio.destination)
let soundBuffers = {}

addEventListener('pointerdown', () => audio.resume(), {once: true})

async function loadSounds(soundFiles) {
    for (let file of soundFiles) {
        try {
            let name = file.replace(/\.[^.]+$/, '')
            let response = await fetch(`sounds/${file}`, {cache: 'no-store'})
            soundBuffers[name] = await audio.decodeAudioData(await response.arrayBuffer())
        } catch (e) {
            console.warn(`couldn't load ${file}`, e)
        }
    }
    console.log(`loaded ${Object.keys(soundBuffers).length} sounds`)
}

function playSound(name, volume = 1, pitch = 1) {
    let buffer = soundBuffers[name]
    if (!buffer) {console.warn(`missing sound ${name}`); return}
    let source = audio.createBufferSource()
    source.buffer = buffer
    let gain = audio.createGain()
    gain.gain.value = volume
    source.playbackRate.value = pitch
    source.connect(gain).connect(masterGain)
    source.start()
}