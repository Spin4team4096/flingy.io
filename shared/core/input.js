mouseX = 0
mouseY = 0
lastMouseX = 0
lastMouseY = 0
leftMouseDown = false
middleMouseDown = false
rightMouseDown = false
isScrollingUp = false
isScrollingDown = false

inputMethod = 'keyboard'

addEventListener('mousemove', e => {
    let k = canvas.width/innerWidth
    mouseX = e.clientX*k-canvas.width/2
    mouseY = e.clientY*k-canvas.height/2
})
addEventListener('mousedown', e => {
    if (e.button === 0) {
        leftMouseDown = true
    } else if (e.button === 1) {
        middleMouseDown = true
    } else if (e.button === 2) {
        rightMouseDown = true
    }
})
addEventListener('mouseup', e => {
    if (e.button === 0) {
        leftMouseDown = false
    } else if (e.button === 1) {
        middleMouseDown = false
    } else if (e.button === 2) {
        rightMouseDown = false
    }
})
addEventListener('contextmenu', e => e.preventDefault())

let scrollTimeout
addEventListener('wheel', (e) => {
    if (e.target.closest('#settings-panel, #changelog-panel')) return
    e.preventDefault()
    isScrollingUp = (e.deltaY < 0)
    isScrollingDown = (e.deltaY > 0)

    clearTimeout(scrollTimeout)
    scrollTimeout = setTimeout(() => {
        isScrollingUp = false
        isScrollingDown = false
    }, 33.33)
})

let keys = {}
addEventListener('keydown', e => {
    if (e.target.matches?.('input, textarea, select') || e.target.isContentEditable) return
    let key = e.key.toLowerCase()
    inputMethod = 'keyboard'
    keys[key] = true
    if(['Tab'].includes(e.key)) {
        e.preventDefault()
    }
})
addEventListener('keyup', e => {
    let key = e.key.toLowerCase()
    keys[key] = false
})

function clearInputState() {
    keys = {}
    leftMouseDown = false
    middleMouseDown = false
    rightMouseDown = false
}
addEventListener('blur', () => clearInputState())
document.addEventListener('visibilitychange', () => {if (document.hidden) clearInputState()})

function isKeyPressed(key) {
    return keys[key.toLowerCase()] == true
}