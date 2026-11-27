function drawRect(x, y, width, height, fillColor, outlineColor, thickness = 0, dir = 0, roundedCorners = null) {
    ctx.save()

    ctx.translate(x, y)
    if (dir != 0) {ctx.rotate(degToRad(dir))}
    
    if (roundedCorners == null) {
        if (fillColor) {
            ctx.fillStyle = fillColor
            
            ctx.fillRect(-width/2, -height/2, width, height)
        }

        if (outlineColor && thickness > 0) {
            ctx.strokeStyle = outlineColor
            ctx.lineWidth = thickness

            ctx.strokeRect(-width/2+thickness/2, -height/2+thickness/2, width-thickness, height-thickness)
        }
    } else {
        // roundedCorners can be [10,10,10,10] or even wacky like [0,0,10,20]
        console.warn(`you dumbo. you need to code rounded corners in if you wanna use em`)
    }
    
    ctx.restore()
}

function drawShape(x, y, radius, n, fillColor, outlineColor, thickness = 0, dir = 0) {
    ctx.save()

    ctx.translate(x, y)
    if (dir != 0) {ctx.rotate(degToRad(dir))}

    if (fillColor) {
        ctx.fillStyle = fillColor

        ctx.beginPath()
        ctx.moveTo(radius, 0)
        for (let p = 1; p < n; p++) {
            ctx.lineTo(Math.cos(degToRad(360*p/n))*radius, Math.sin(degToRad(360*p/n))*radius)
        }
        ctx.closePath()
        ctx.fill()
    }

    if (outlineColor && thickness > 0) {
        let t = Math.min(thickness, radius)
        let r = radius-t/2

        ctx.strokeStyle = outlineColor
        ctx.lineWidth = t

        ctx.beginPath()
        ctx.moveTo(r, 0)
        for (let p = 1; p < n; p++) {
            ctx.lineTo(Math.cos(degToRad(360*p/n))*(r), Math.sin(degToRad(360*p/n))*r)
        }
        ctx.closePath()
        ctx.stroke()
    }

    ctx.restore()
}

function drawCircle(x, y, radius, fillColor, outlineColor, thickness = 0) {
    ctx.save()

    ctx.translate(x, y)

    if (fillColor) {
        ctx.fillStyle = fillColor

        ctx.beginPath()
        ctx.arc(0, 0, radius, 0, 2*Math.PI)
        ctx.closePath()
        ctx.fill()
    }

    if (outlineColor && thickness > 0) {
        let t = Math.min(thickness, radius)
        let r = radius-t/2

        ctx.strokeStyle = outlineColor
        ctx.lineWidth = t

        ctx.beginPath()
        ctx.arc(0, 0, r, 0, 2*Math.PI)
        ctx.closePath()
        ctx.stroke()
    }

    ctx.restore()
}

function drawLine(x, y, points, fillColor, outlineColor, thickness = 0, dir = 0, closed = false) {
// this one doesn't inset the line thickness bc impossible.
    ctx.save()

    ctx.translate(x, y)
    if (dir != 0) {ctx.rotate(degToRad(dir))}

    ctx.beginPath()
    ctx.moveTo(points[0][0], points[0][1])
    for (let n = 1; n < points.length; n++) {
        ctx.lineTo(points[n][0], points[n][1])
    }

    if (fillColor) {
        ctx.fillStyle = fillColor
        ctx.fill()
    }

    if (outlineColor && thickness > 0) {
        ctx.lineWidth = thickness
        ctx.strokeStyle = outlineColor

        if (closed) ctx.closePath()
        ctx.stroke()
    }

    ctx.restore()
}

function drawText(text, x, y, size, fillColor, outlineColor, thickness = 0, dir = 0, align = 'center', stackDirection = 'center') {
    text = text.split('\n')

    for (let i = 0; i < text.length; i++) {
        ctx.save()

        let lineY = 0
        switch (stackDirection) {
            case 'up': lineY = -(text.length-1-i)*size*1.2; break
            case 'down': lineY = i*size*1.2; break
            case 'center': lineY = (i-(text.length-1)/2)*size*1.2; break
        }
        
        ctx.translate(x, y+lineY)
        if (dir != 0) {ctx.rotate(degToRad(dir))}

        let bold = false
        if (text[i].startsWith('**') && text[i].endsWith('**') && text[i].length > 4) {
            bold = true
            text[i] = text[i].slice(2,-2)
        }

        ctx.textAlign = align
        ctx.textBaseline = 'middle'
        ctx.font = `${bold ? 'bold ' : ''}${size}px Ubuntu`

        if (outlineColor && thickness > 0) {
            ctx.strokeStyle = outlineColor
            ctx.lineWidth = thickness

            ctx.strokeText(text[i], 0, 0)
        }

        if (fillColor) {
            ctx.fillStyle = fillColor
            
            ctx.fillText(text[i], 0, 0)
        }

        ctx.restore()
    }
}