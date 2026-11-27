
const app = firebase.initializeApp({
    databaseURL: "https://spinbase4096-default-rtdb.europe-west1.firebasedatabase.app/"
})
const database = firebase.database()
        
function watchFirebaseServers(path) {
    database.ref(path).on("value", (snapshot) => {
        servers = snapshot.val()
    })
}

function watchFirebaseStatus(path, cache) {
    database.ref(path).on("value", (snapshot) => {
        statusMessage = snapshot.val()
        renderStatus(statusMessage)
        if (statusMessage && cache) {
            setStoredKey(cache, JSON.stringify(statusMessage))
        } else {
            deleteStoredKey(cache)
        }
    })
}