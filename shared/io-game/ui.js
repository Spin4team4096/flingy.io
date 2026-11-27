// required vars: settings, settingsData, settingsPage, ui.settings, loadSettings(), playSound(), randomFloat(), renderPanel(), ui.menu.status, and more idk man. just have everything bro

ui = {
    menu: {
        myself: document.getElementById('menu'),
        logo: document.getElementById('logo'),
        menuControls: {
            colorButton: document.getElementById('color-button'),
            username: document.getElementById('username'),
            playButton: document.getElementById('play-button'),
        },
        message: document.getElementById('message'),
        status: document.getElementById('status'),
        gamemodeList: document.getElementById('gamemode-list'),
        shop: {
            coins: document.getElementById('coins'),
            buyColor: document.getElementById('buy-color'),
        },
        bottomRow: {
            settingsButton: document.getElementById('settings-button'),
            changelogButton: document.getElementById('changelog-button'),
            discordButton: document.getElementById('discord-button'),
        },
        tips: document.getElementById('tips'),
    },
    hud: {
        myself: document.getElementById('hud'),
        leaveButton: document.getElementById('leave-button'),
        notifications: document.getElementById('notifications'),
        leaderboard: document.getElementById('leaderboard'),
        bottomMessage: document.getElementById('bottom-message'),
    },
    settings: document.getElementById('settings-panel'),
    changelog: document.getElementById('changelog-panel'),
}

animations = {
    colorButton: {
        type: 0, // 0 no animation, 1 is free spin, 2 is 90 deg spin
        angle: 0,
        progress: 1,
        from: 0,
        to: 0,

        state: 0, // 0 idle, 1 accelerating, 2 decelleraintgjngg
        speed: 0,
        time: 0,
        acceleration: 0,
        lastCrossing: 0,
    },
    playButton: {
        type: 0, // 0 no animation, 1 is free spin, 2 is 90 deg spin, 3 is wiggle
        state: 0, // 0 idle, 1 accelerating, 2 decelleraintgjngg. this is why it's numbers, cuz i can't speel
        result: 0, // 0 waiting green, 1 success green, 2 failed yellow, 3 waiting red
        angle: 0,
        speed: 0,
        time: 0,
    },
    settingsButton: {
        type: 0, // 0 no animation, 1 is 90 spin clockwise, 2 is 90 spin anticlockwise
        angle: 0,
        progress: 1,
    },
    changelogButton: {
        type: 0, // 0 no animation, 1 is 90 spin clockwise, 2 is 90 spin anticlockwise
        angle: 0,
        progress: 1,
    },
    discordButton: {
        type: 0, // 0 no animation, 1 is 180 spin
        angle: 0,
        progress: 1,
    },
    bgImage: {
        opacity: 1, 
    }
}

ui.menu.menuControls.colorButton.addEventListener('pointerdown', () => {
    if (animations.colorButton.type == 1) return
    playSound('bwooop', 1, randomFloat(0.9,1.1))
    animations.colorButton.type = 2
    animations.colorButton.to += 90
    animations.colorButton.from = animations.colorButton.angle
    animations.colorButton.progress = 0

    playerColor = unlockedColors[(unlockedColors.indexOf(playerColor)+1)%unlockedColors.length]
    ui.menu.menuControls.colorButton.style.setProperty('--c', `var(--${playerColor})`)
    ui.menu.menuControls.colorButton.style.setProperty('--c-outline', `var(--${playerColor}-outline)`)
})

ui.menu.menuControls.playButton.addEventListener('pointerdown', () => {
    if (animations.playButton.state != 0 && animations.playButton.result < 2) return
    playSound('bwooop', 1, randomFloat(0.9,1.1))
    animations.playButton.type = 1
    animations.playButton.state = 1
    animations.playButton.result = 0
    animations.playButton.angle = 0
    animations.playButton.speed = 0
    animations.playButton.time = 0
    playerSlot = -1
    pickServer()
})

settingsClearTimeout = null
function toggleSettings() {
    settingsPage = settingsPage != 0 ? 0 : 1
    ui.settings.classList.toggle('open', settingsPage != 0)
    if (changelogOpen) {
        changelogOpen = !changelogOpen
        ui.changelog.classList.toggle('open', changelogOpen)
        animations.changelogButton.type = changelogOpen == true ? 1 : 2
        animations.changelogButton.angle = 0
        animations.changelogButton.progress = 0
    }
    playSound('bwooop', 1, settingsPage != 0 ? 1.1 : 0.9)
    animations.settingsButton.type = settingsPage != 0 ? 1 : 2
    animations.settingsButton.angle = 0
    animations.settingsButton.progress = 0
    clearTimeout(settingsClearTimeout)
    if (settingsPage != 0) {
        renderPanel(ui.settings, settingsData[settingsPage])
    } else {
        settingsClearTimeout = setTimeout(() => {
            ui.settings.replaceChildren()
        }, 250*settings.animationTime)
    }
}
ui.menu.bottomRow.settingsButton.addEventListener('pointerdown', () => {toggleSettings()})

changelogClearTimeout = null
function toggleChangelog() {
    changelogOpen = !changelogOpen
    ui.changelog.classList.toggle('open', changelogOpen)
    if (settingsPage != 0) {
        settingsPage = settingsPage != 0 ? 0 : 1
        ui.settings.classList.toggle('open', settingsPage != 0)
        animations.settingsButton.type = settingsPage != 0 ? 1 : 2
        animations.settingsButton.angle = 0
        animations.settingsButton.progress = 0
    }
    playSound('bwooop', 1, changelogOpen == true ? 1.1 : 0.9)
    animations.changelogButton.type = changelogOpen == true ? 1 : 2
    animations.changelogButton.angle = 0
    animations.changelogButton.progress = 0
    clearTimeout(changelogClearTimeout)
    if (changelogOpen) {
        renderPanel(ui.changelog, changelogData)
    } else {
        changelogClearTimeout = setTimeout(() => {
            ui.changelog.replaceChildren()
        }, 250*settings.animationTime)
    }
}
ui.menu.bottomRow.changelogButton.addEventListener('pointerdown', () => {toggleChangelog()})

ui.menu.bottomRow.discordButton.addEventListener('click', () => {
    if (animations.discordButton.progress < 1) return
    playSound('bwooop', 1, randomFloat(0.9,1.1))
    animations.discordButton.type = 1
    animations.discordButton.angle = 0
    animations.discordButton.progress = 0
    setTimeout(() => {
        window.open('https://discord.gg/z2uZvEXfPx', '_blank')
    }, 300*settings.animationTime)
})

ui.menu.shop.buyColor.addEventListener('click', () => {
    if (animations.colorButton.type != 0) return
    if (coins < colorPrice || remainingColors.length == 0) {
        playSound('negativeaction')
        return
    }
    coins -= colorPrice
    animations.colorButton.type = 1
    animations.colorButton.state = 1
    animations.colorButton.result = 0
    animations.colorButton.angle = 0
    animations.colorButton.speed = 0
    animations.colorButton.time = 0
    animations.colorButton.acceleration = randomFloat(3.5, 4.0)
    animations.colorButton.lastCrossing = 0
})

ui.hud.leaveButton.addEventListener('click', () => {
    playSound('negativeaction2')
    ws.close()
})

ui.menu.menuControls.username.addEventListener('input', () => {
    username = ui.menu.menuControls.username.value
})

addEventListener('pointerdown', e => {
    if (settingsPage != 0 && !ui.settings.contains(e.target) && !ui.menu.bottomRow.settingsButton.contains(e.target)) {
        toggleSettings()
    }
    if (changelogOpen && !ui.changelog.contains(e.target) && !ui.menu.bottomRow.changelogButton.contains(e.target)) {
        toggleChangelog()
    }
})

function parseMarkdown(text) {
    let fragment = document.createDocumentFragment()
    let lines = text.split(" \\n ")
    for (let i = 0; i < lines.length; i++) {
        fragment.append(...parseMarkdownRange(lines[i]))
        if (i < lines.length-1) {
            fragment.append(document.createElement('br'))
        }
    }
    return fragment
}

function parseMarkdownRange(text) {
    let tokens = [
        ['***', '***'],
        ['**', '**'],
        ['*', '*'],
        ['__', '__'],
        ['~~', '~~'],
        ['`', '`'],
        ['==', '=='],
    ]
    let nodes = []
    
    let i = 0
    outer: while (i < text.length) {
        if (text.startsWith('[', i)) {
            let nameEnd = text.indexOf(']', i+1)
            let linkStart = nameEnd != -1 ? text.indexOf('(', nameEnd+1) : -1
            let linkEnd = linkStart != -1 ? text.indexOf(')', linkStart+1) : -1

            if (linkStart == nameEnd + 1 && linkEnd != -1) {
                let a = document.createElement('a')
                a.href = text.slice(linkStart+1, linkEnd)
                a.append(...parseMarkdownRange(text.slice(i+1, nameEnd)))
                nodes.push(a)
                i = linkEnd+1
                continue outer
            }
        }
        for (let t = 0; t < tokens.length; t++) {
            if (text.startsWith(tokens[t][0], i)) {
                let cloze = text.indexOf(tokens[t][1], i+tokens[t][0].length)
                if (cloze != -1) {
                    let children = parseMarkdownRange(text.slice(i+tokens[t][0].length, cloze))
                    let e = null
                    switch (tokens[t][0]) {
                        case '***': 
                            let i = document.createElement('i')
                            i.append(...children)
                            e = document.createElement('b')
                            e.append(i)
                            break
                        case '**': 
                            e = document.createElement('b')
                            e.append(...children)
                            break
                        case '*': 
                            e = document.createElement('i')
                            e.append(...children)
                            break
                        case '__': 
                            e = document.createElement('u')
                            e.append(...children)
                            break
                        case '~~': 
                            e = document.createElement('s')
                            e.append(...children)
                            break
                        case '`': 
                            e = document.createElement('code')
                            e.append(...children)
                            break
                        case '==':
                            e = document.createElement('span')
                            e.style.color = '#ffffff'
                            e.append(...children)
                            break
                    }
                    nodes.push(e)
                    i = cloze+tokens[t][1].length
                    continue outer
                }
            }
        }
        nodes.push(document.createTextNode(text[i]))
        i++
    }
    return nodes
}

function renderPanel(panel, data) {
    panel.replaceChildren()
    let e = document.createElement(`div`)
    e.className = 'title'
    e.textContent = data['title']
    e.onclick = () => {
        switch (settingsPage) {
            case 0: {
                switch (changelogOpen) {
                    case true: {
                        toggleChangelog()
                    } break
                }
            } break
            case 1: {
                toggleSettings()
            } break
            case 2: {
                settingsPage = 1
                renderPanel(ui.settings, settingsData[settingsPage])
                playSound('bwooop', 1, 0.9)
            } break
            case 3: {
                settingsPage = 1
                renderPanel(ui.settings, settingsData[settingsPage])
                playSound('bwooop', 1, 0.9)
            } break
            case 4: {
                settingsPage = 1
                renderPanel(ui.settings, settingsData[settingsPage])
                playSound('bwooop', 1, 0.9)
            } break
            default: {
                settingsPage = 1
                renderPanel(ui.settings, settingsData[settingsPage])
                playSound('bwooop', 1, 0.9)
            } break
        }
    }
    panel.append(e)

    renderPanelNodes(panel, data['children'])
}

function renderPanelNodes(parent, nodes) {
    for (let node of nodes) {
        switch (node.type) {
            case 'section': {
                let details = document.createElement(`details`)
                let summary = document.createElement(`summary`)
                details.className = 'section'
                summary.textContent = node.label
                summary.style.setProperty('--c', `var(--${node.color})`)
                summary.style.setProperty('--c-outline', `var(--${node.color}-outline)`)
                details.style.setProperty('--c', `var(--${node.color})`)
                details.style.setProperty('--c-outline', `var(--${node.color}-outline)`)
                summary.insertAdjacentHTML('beforeend', `<svg class="section-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M5 5a2 2 0 0 1 3.008-1.728l11.997 6.998a2 2 0 0 1 .003 3.458l-12 7A2 2 0 0 1 5 19z"/></svg>`)
                details.open = node.open == true
                details.append(summary)
                summary.onclick = () => {
                    playSound('bwooop', 1, details.open ? 0.9 : 1.1)
                }

                parent.append(details)
                renderPanelNodes(details, node.children)
                
            } break
            case 'button': {
                let e = document.createElement(`button`)
                e.className = 'button'
                e.textContent = node.label
                e.style.setProperty('--c', `var(--${node.color})`)
                e.style.setProperty('--c-outline', `var(--${node.color}-outline)`)
                if (node.page) {
                    e.onclick = () => {
                        settingsPage = node.page
                        renderPanel(ui.settings, settingsData[settingsPage])
                        playSound('bwooop', 1, randomFloat(0.9,1.1))
                    }
                }
                parent.append(e)
            } break
            case 'toggle': {
                let e = document.createElement(`div`)
                e.className = 'label'
                e.textContent = node.label

                let input = document.createElement('input')
                input.type = 'checkbox'
                input.checked = !!settings[node.setting]
                input.addEventListener('change', () => {
                    settings[node.setting] = input.checked
                    loadSettings()
                })

                e.append(input)
                parent.append(e)
            } break
            case 'slider': {
                let label = document.createElement(`div`)
                label.className = 'label'
                label.textContent = node.label

                let value = document.createElement(`div`)
                value.className = 'value'
                value.textContent = settings[node.setting].toFixed(2)

                let input = document.createElement('input')
                input.type = 'range'
                input.min = node.min
                input.max = node.max
                input.step = node.step
                input.value = settings[node.setting]
                input.addEventListener('input', () => {
                    value.textContent = Number(input.value).toFixed(2)
                })
                input.addEventListener('change', () => {
                    settings[node.setting] = Number(input.value)
                    value.textContent = Number(input.value).toFixed(2)
                    loadSettings()
                })

                label.append(value)
                label.append(input)
                parent.append(label)
            } break
            case 'choice': {
                let label = document.createElement(`div`)
                label.className = 'label'
                label.textContent = node.label

                let options = document.createElement(`div`)
                options.className = 'options'

                for (let b = 0; b < node.options.length; b++) {
                    let choice = document.createElement(`button`)
                    choice.className = 'choice'
                    choice.textContent = node.options[b].label
                    choice.classList.toggle('active', settings[node.setting] == node.options[b].value)

                    choice.onclick = () => {
                        settings[node.setting] = node.options[b].value
                        loadSettings()
                        choice.classList.toggle('active', settings[node.setting] == node.options[b].value)
                        for (let o of options.children) {o.classList.toggle('active', o == choice)}
                    }

                    options.append(choice)
                }

                label.append(options)
                parent.append(label)
            } break

            // changelog stuff
            case 'text': {
                let text = document.createElement('div')
                text.className = 'text'
                text.replaceChildren(parseMarkdown(node.label))
                text.style.opacity = node.opacity || 1
                if (node.size != null) {
                    text.style.fontSize = `calc(${node.size}px*var(--ui))`
                }

                parent.append(text)
            } break
            case 'server': {
                let subtitle = document.createElement('div')
                subtitle.className = 'subtitle'
                subtitle.textContent = 'Server:'

                parent.append(subtitle)
                renderPanelNodes(subtitle, node.children)

            } break
            case 'client': {
                let subtitle = document.createElement('div')
                subtitle.className = 'subtitle'
                subtitle.textContent = 'Client:'

                parent.append(subtitle)
                renderPanelNodes(subtitle, node.children)
            } break
            case 'point': {
                let text = document.createElement('div')
                text.className = 'point'
                text.replaceChildren(parseMarkdown(node.label))

                parent.append(text)
                if (node.children) {renderPanelNodes(text, node.children)}
            } break
        }
    }
}

function renderStatus(statusMessage) {
    if (statusMessage && ui.menu.status) {
        ui.menu.status.replaceChildren(parseMarkdown(statusMessage.text))
        ui.menu.status.style.setProperty('--c', `var(--${statusMessage.color})`)
        ui.menu.status.style.setProperty('--c-outline', `var(--${statusMessage.color}-outline)`)
    } else if (ui.menu.status) {
        ui.menu.status.replaceChildren()
    }
}

function tipsLoop() {
    tipIndex = 0
    ui.menu.tips.replaceChildren(parseMarkdown(tips[tipIndex]))
    ui.menu.tips.style.opacity = 1
    setInterval(() => {
        ui.menu.tips.style.opacity = 0
        setTimeout(() => {
            tipIndex = (tipIndex+1)%tips.length
            ui.menu.tips.replaceChildren(parseMarkdown(tips[tipIndex]))
            if (scene == 0) {ui.menu.tips.style.opacity = 1}
        }, 1000*settings.animationTime)
    }, 10000)
}