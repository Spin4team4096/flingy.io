# Cobbled.io Todo

### UI / leftovers
- `p`/`m` keybinds exist, but there's no minimap element (leaderboard exists but has no data yet)
- `gamemode-list` element exists but nothing renders it; group from `servers[].config.gamemode`
- `initLobby()` hook for when `connect()` moves to shared `net.js`

### Deploy / archive
- Vercel project: root/output `cobbled.io/client`, build copies shared (`cp -r shared/.` into output)
- `old.cobbled.io` archive project: index + v1.5.2 / v2.7.1 / v3.0.3 single-file builds; state that v2/v3 are archive-only (need their old servers)
- `oldCobbled-2026-10-01.json`: keep private (90 players' names/UUIDs); don't restore the 79 stale servers; status text can be static
- changelog: build v1/v2/v3/v4 sections from the old Discord announcements

### Later
- actual coins earning + game logic (UI/shop shell already copied)
- tailor `settings.json` / `changelog.json` for v4 (right now they're flingy's content)

### Shared engine
- roadmap and shared bugs live in `../shared/TODO.md` (no duplicated items here)
