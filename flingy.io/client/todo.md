# Todo

### UI
- right click to open color dropdown
- gamemode selector in the menu (FFA / TDM), also used to filter `pickServer` and drive `gamemodeList`
- hide the HUD when scene == 0 (leave button etc.)
- tips are meant to be randomly shuffled every cycle like in old flingy.io

### HUD
- notifications
- leaderboard
- throw path indicator
- death screen

### Debug
- debug menu: draw the squares for the borderScale/mouseScale clamp ranges

### Misc
- in-game sounds (UI sounds exist; gameplay sfx don't)
- route the fling input through `settings.keybinds.fling` (already `['mouse0', ' ']`); gameplay still reads raw `leftMouseDown`, so derive the packet's `down` flag from the bind state
- live ping: either nothing or a real ping/pong packet (join-handshake RTT isn't live)
- expire `invalidServers` marks after ~30-60s using timestamps instead of blacklisting for the whole session
- interpolation/predict for entity movement (the smoothing setting)
- reuse one global `TextDecoder` instead of creating one per username in packets 4/5
- // in the future, make it wait for a server to be available, since that's the point after all
- the the game is in portrait (taller than width), put an overlay with text "rotate your device to landscape mode". reuse `icon-play`. if possible, add an `i don't care` button.
- a slider or other ui that tells you how long until the fling you flung returns to you. basically, predict that by, when your flail is flung, it looks at da speed and size and simulates x server ticks until the reconnect condition is met, and then uses that tick count as a timer. add this when you add prediction. this'll probably be a prediction-only feature, but i might also add this to interpolation too. a bar with 5s at the top and 0s at the bottom, on the right side of the screen, is good.
- before writing to localstorage, compare values so that you don't lose progress if multiple tabs are open (or maybe not, why should ppl multibox)

### Shared engine
- roadmap and shared bugs live in `../../shared/TODO.md` (no duplicated items here)
