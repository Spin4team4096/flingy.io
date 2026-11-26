// required vars: settings, settingsData, settingsPage, ui.settings, loadSettings(), playSound(), randomFloat(), renderPanel(), ui.menu.status

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