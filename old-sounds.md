# Old Flingy.io Sound Reference

> DISCLAIMER: This .md was generated with AI, for my reference. Nothing is vibe-coded.

Extracted from `Very Serious Game Jam/index.html` (direct `playSound()` calls) and
`Very Serious Game Jam/server.html` (events sent to clients as `playSound` packets
via `spawnSound()`). All files now live in `client/sounds/` as OGG Vorbis.

## Used

| Sound | Plays when |
|---|---|
| `bonus` | Buy-colour wheel finishes spinning and unlocks the new colour |
| `bwooop` | 3 cases: switching game mode button, colour button click (cycles colour), buy-colour wheel crossing each 90° |
| `eliminated` | A knocked-out player respawns (alive set back to 1) — only that player hears it |
| `gamestart` | Match starts — everyone, full volume, no distance fade |
| `hit1` | You hit another player — random 50/50 with `hit2` |
| `hit2` | Same, the other half of the random pair |
| `lose` | A player's HP hits 0 and they're out for good — only that player hears it |
| `matchfound` | A server was picked, starting the connection |
| `negativeaction2` | Every failure path: no servers, all servers busy, no available server, peer connect error, server disconnect packet, peer error, lost connection (10s no updates), leaving a game |
| `ping` | A player or bot joins (`initializePlayer`) — everyone, no fade |
| `sling` | A player flings (releases the flail) |
| `success` | While waiting on busy servers, one frees up ("Game found!") |
| `uipopup` | Discord button click |
| `victory` | Last player standing wins — only the winner hears it |

## Unused in old flingy.io

`achievementmc`, `negativeaction` (only the `2` variant was ever called), `pointsadd`,
`throw`, `uiclick`, `uiclick2`, `uipopup2`, `uipopupclose`

## Notes for reimplementation

- **Positional sounds** (`sling`, `hit1`, `hit2`, `victory`, `lose`, `eliminated`) were
  sent with world coordinates and faded by distance from the camera (640-unit range).
- **Private sounds** (`victory`, `lose`, `eliminated`) were sent only to one player.
- **`gamestart`** and **`ping`** were sent with `fade: false`, i.e. always full volume.
- World/menu split: everything except `sling`, `hit1`, `hit2`, `victory`, `lose`,
  `eliminated`, `gamestart`, `ping` is a menu/UI sound.
