colors = {
    "fill": {
        "red": "#ff6666",
        "coral": "#ff9966",
        "orange": "#ffc266",
        "yellow": "#ffe766",
        "lime": "#87e65c",
        "green": "#71bf4d",
        "blue": "#6696ff",
        "purple": "#9966ff",
        "magenta": "#fc66ff",
        "tan": "#ffccb3",
        "navy": "#4d71bf",
        "brown": "#bf734d",
        "maroon": "#bf4d4d",
        "burgundy": "#bf4d76",
        "pink": "#ff66af",
        "cyan": "#5ce6e6",
        "mint": "#5ce693",
        "black": "#333333",
        "gray": "#bfbfbf",
        "white": "#f2f2f2"
    },
    "outline": {
        "red": "#bf4d4d",
        "coral": "#bf734d",
        "orange": "#bf914d",
        "yellow": "#bfad4d",
        "lime": "#62a642",
        "green": "#4b8033",
        "blue": "#4d71bf",
        "purple": "#734dbf",
        "magenta": "#bd4dbf",
        "tan": "#bf9986",
        "navy": "#334b80",
        "brown": "#804d33",
        "maroon": "#803333",
        "burgundy": "#80334f",
        "pink": "#bf4d84",
        "cyan": "#4dbfbf",
        "mint": "#4dbf7a",
        "black": "#262626",
        "gray": "#999999",
        "white": "#bfbfbf"
    }
}

playerColors = ['red', 'coral', 'orange', 'yellow', 'lime', 'green', 'mint', 'cyan', 'blue', 'navy', 'purple', 'pink', 'magenta', 'burgundy', 'maroon', 'brown', 'tan', 'white', 'gray', 'black']
for (let color of playerColors) {
    document.documentElement.style.setProperty(`--${color}`, colors.fill[color])
    document.documentElement.style.setProperty(`--${color}-outline`, colors.outline[color])
}

function toHex(num) {
    return Math.round(clamp(num, 0, 1)*255).toString(16).padStart(2, "0")
}

function interpolateColors(color1, color2, weight) {
    weight = Math.max(0, Math.min(1, weight))
    
    const r1 = parseInt(color1.slice(1,3),16)
    const g1 = parseInt(color1.slice(3,5),16)
    const b1 = parseInt(color1.slice(5,7),16)
    
    const r2 = parseInt(color2.slice(1,3),16)
    const g2 = parseInt(color2.slice(3,5),16)
    const b2 = parseInt(color2.slice(5,7),16)
    
    const r = Math.round(r1+(r2-r1)*weight)
    const g = Math.round(g1+(g2-g1)*weight)
    const b = Math.round(b1+(b2-b1)*weight)
    
    return `#${r.toString(16).padStart(2, '0')}${g.toString(16).padStart(2, '0')}${b.toString(16).padStart(2, '0')}`
}

// depreciated
function colorsToHex(color) {
    switch (color) {
        case 'red': return '#ff6666'
        case 'coral': return '#ff9966'
        case 'orange': return '#ffc266'
        case 'yellow': return '#ffe766'
        case 'lime': return '#87e65c'
        case 'green': return '#71bf4d'
        case 'blue': return '#6696ff'
        case 'purple': return '#9966ff'
        case 'magenta': return '#fc66ff'
        case 'tan': return '#ffccb3'
        case 'navy': return '#4d71bf'
        case 'brown': return '#bf734d'
        case 'maroon': return '#bf4d4d'
        case 'burgundy': return '#bf4d76'
        case 'pink': return '#ff66af'
        case 'cyan': return '#5ce6e6'
        case 'mint': return '#5ce693'
        case 'black': return '#333333'
        case 'gray': return '#bfbfbf'
        case 'white': return '#f2f2f2'
        case 'special': {
            t = Date.now()/640
            start = [255,102,102]
            end = [255,128,166]
            factor = 0.5+0.5*Math.sin(t)
            r = Math.round(start[0]+(end[0]-start[0])*factor)
            g = Math.round(start[1]+(end[1]-start[1])*factor)
            b = Math.round(start[2]+(end[2]-start[2])*factor)
            return `#${r.toString(16).padStart(2,'0')}${g.toString(16).padStart(2,'0')}${b.toString(16).padStart(2,'0')}`
        }
        default: return '#ff00ff'
    }
}

function colorsToOutlineHex(color) {
    switch (color) {
        case 'red': return '#bf4d4d'
        case 'coral': return '#bf734d'
        case 'orange': return '#bf914d'
        case 'yellow': return '#bfad4d'
        case 'lime': return '#62a642'
        case 'green': return '#4b8033'
        case 'blue': return '#4d71bf'
        case 'purple': return '#734dbf'
        case 'magenta': return '#bd4dbf'
        case 'tan': return '#bf9986'
        case 'navy': return '#334b80'
        case 'brown': return '#804d33'
        case 'maroon': return '#803333'
        case 'burgundy': return '#80334f'
        case 'pink': return '#bf4d84'
        case 'cyan': return '#4dbfbf'
        case 'mint': return '#4dbf7a'
        case 'black': return '#262626'
        case 'gray': return '#999999'
        case 'white': return '#bfbfbf'
        case 'special': {
            t = Date.now()/640
            start = [191,77,77]
            end = [191,96,124]
            factor = 0.5+0.5*Math.sin(t)
            r = Math.round(start[0]+(end[0]-start[0])*factor)
            g = Math.round(start[1]+(end[1]-start[1])*factor)
            b = Math.round(start[2]+(end[2]-start[2])*factor)
            return `#${r.toString(16).padStart(2,'0')}${g.toString(16).padStart(2,'0')}${b.toString(16).padStart(2,'0')}`
        }
        default: return '#ff0000'
    }
}