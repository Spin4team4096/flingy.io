# Shared Engine Todo

Classic scripts, no build step — globals are the API. Keep each file's contract comment accurate.
Only extract once a second game needs it; contracts can change as cobbled grows.
One todo, one home: game-specific stuff stays in that game's todo.

### Critical
- game keydown listeners have no typing guard: while typing in the username `<input>`, f/i/o/p/m/escape/tab still fire actions. Add `isTypingTarget(e)` to core and `if (isTypingTarget(e)) return` at the top of each game's listener

### Structure
- input leftovers: the wheel handler still checks `#settings-panel, #changelog-panel` — add `data-scrollable` to both panels and check `closest('[data-scrollable]')`; then move `input.js` back into the core block (tags: core → io-game → game)
- add `isBindDown(binds)` to core (token map: `mouse0/1/2` → mouse state, else `isKeyPressed`); use it for the fling `down` flag instead of raw `leftMouseDown`
- keybind settings UI (rebinding writes `settings.keybinds`; it already persists via the save loop)
- update the `ui.js` contract comment — it now provides `ui`, `renderPanel`, `animations`, `tipsLoop` instead of requiring them
- delete the now-unused `colorsToHex`/`colorsToOutlineHex` from `core/colors.js` (zero call sites; `special` returns later as a theme function value if ever needed)
- per-frame button animation update loops → `io-game/ui.js` (settings/changelog/discord now; color/play once cobbled has game logic)
- FPS/MSPF counters + menu background drift helpers
- mobile + gamepad controls: implement once, shared
- `io-game/net.js`: `pickServer(predicate)`, `connect(id, url, {initLobby, onPacket})`, handshake framing; packet version as a parameter/constant
- `io-game/save.js`: the localStorage loop parameterized by `{prefix, version, defaults, migrate}` (only once cobbled has real saves)
- per-game `config.js` manifest: storagePrefix, mainColor, firebaseBase, soundFiles, menuBg, tipsFile

### core/draw.js
- drawRect() rounded-corner support (currently just warns)

### Data-driven
- `tips.json` per game
- debug panel from data (+ settings toggles to show/hide each section)
- gamemode list derived from `servers[].config.gamemode`; `pickServer` filters by the selected gamemode

### Go
- decide: one module at repo root vs `go.work`
- shared Go packages: mathx, websocket/broadcast, lobby framework

### Tools
- root `dev.sh` at `spin4games/` (one static server for all games); the port split (server+tunnel 4095, static 4096) is already set in flingy's dev.sh
- dev static server should send `Cache-Control: no-cache` for html/js/css (`python -m http.server` lets browsers heuristically cache shared scripts, so edits need a hard refresh today)
- move changelog updater to `spin4games/tools/` (root, tracked)
- multi-game support: per-game dirHandle + webhook, game selector
- import `shared/core/colors.js` for the palette
- recover the old cobbled announcements into changelog data (v1/v2/v3)

# Contracts (reference — not todos)
- `canvas`/`ctx` stay per game; draw/input read the global `canvas`/`ctx` by design (core input assumes `canvas` is CSS-sized to the full viewport for mouse mapping)
- `ui.js` owns the menu DOM refs — games must keep the same element ids
- `ui.js` listeners also touch game globals at event time (`playerColor`, `unlockedColors`, `pickServer`)