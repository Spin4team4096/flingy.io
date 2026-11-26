function setStoredKey(key, value) {
    localStorage.setItem(key, value)
}

function getStoredKey(key) {
    return localStorage.getItem(key) || ''
}

function deleteStoredKey(key) {
    localStorage.removeItem(key)
}