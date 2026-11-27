function randomInt(num1, num2) {return Math.floor(Math.random()*(Math.max(num1, num2)-Math.min(num1, num2)+1))+Math.min(num1, num2)}
function randomFloat(num1, num2) {return Math.random()*(Math.max(num1, num2)-Math.min(num1, num2))+Math.min(num1, num2)}
function clamp(value, min, max) {return Math.max(min, Math.min(max, value))}
function degToRad(d) {return d*(Math.PI/180)}
function radToDeg(r) {return r*(180/Math.PI)}
function lerp(a, b, t) {return a+(b-a)*t}