# Todo
- rendering
- uhhhhhhhhhh everything

### Rendering
- drawRect() rounded corner support

### UI
- right click to open color dropdown
- gamemode selector in the menu (FFA / 2TDM), also used to filter `pickServer`
- hide the HUD when scene == 0 (leave button etc.)
- tips are meant to be randomly shuffled every cycle like in old flingy.io

### Debug
- debug menu: draw the squares for the borderScale/mouseScale clamp ranges

### Misc
- live ping: either nothing or a real ping/pong packet (join-handshake RTT isn't live)
- expire `invalidServers` marks after ~30-60s using timestamps instead of blacklisting for the whole session
- interpolation/predict for entity movement (the smoothing setting)
- reuse one global `TextDecoder` instead of creating one per username in packets 4/5
- // in the future, make it wait for a server to be available, since that's the point after all
- the the game is in portrait (taller than width), put an overlay with text "rotate your device to landscape mode". reuse `icon-play`. if possible, add an `i don't care` button.