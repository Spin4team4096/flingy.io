# TODO
- everything

### Bugs
- there is no duplicate connection check. a client can just: connect multiple times lol.
  - in fact, there's no checks at all. a player can just join midgame rn.
- `go http.ListenAndServe(...)` error is ignored: a second server instance logs "Started websocket server" but can't bind, and silently registers duplicate lobbies. Log the error (and maybe exit).

### Team Mode
- rank/placement per team instead of per player
- send Team to the client (packet 3 flags bits 3-7 are all colourIndex, so needs an encoding decision)

### Bots
- implement bots
- use bots to fill empty team slots in TDM so odd counts aren't 2v1/3v2; until then uneven teams are allowed

### Packets
- send statetick so the client can calc this for himself

### Commands
- range support: `unregister 0001-0004` and `kill 0001-0004`, same as the single/all versions
  - also would be good to specify `ffa` or `tdm` icl. `kill ffa_0001-0004` is nice. 
  - maybe even super-shorthand support: `kill ffa_1-4` is awesome. so is `kill f1-4` that is beautiful. maybe that's the meta.
- `status` that lists the servers AND their status (lagging/joinable/full), or `status <serverID>` for only that server's status
- `stats` that measures and gives avg tps and mspt etc, and `stats <serverID>` for that
- give command finish feedback on those currently silent commands, like `unregister all` would log `unregistered 7 lobbies` for example.
- `list <parameter>`, defaults to assuming `local`, but you can also list everything in `firebase`

### Auto-Update
- the server laptop should automatically scan for github updates, and if there is one (that actually changes server code, not just a todo list, or /client /tools etc.):
  - pulls files from git
  - launches a new terminal, runs `go run ./server`
  - on the existing server terminal, `close all`, wait for that to happen, then kill that terminal.
- also, automatically starting the cloudflare tunnel upon system startup would be good.

### Misc
- the player can spin an extra revolution sometimes, tho sometimes the player just stacked a lot of spins but they changed their mine cuz they're not spinning anymore. in that case, the spin accumullator should -360 early. trigger should probably be: the mouse angle goes the other way (not: the mouse stops moving. cuz client can lag.). 
  - what i mean is: if the accumulated spin is positive, and there's a lot of accumulated spin, maybe if > `spinSpeed` amount is in it, but the mouse angle diff between now and last tick is negative, then the direction changed: -360 the accumulator to canel it early.

### Random Ideas
- Have a server health tracking system, like a command `health` and a struct that contains health info;
  - it also pings @spin4team4096 on discord when all lobbies are ongoing (will need to track this somehow) and there are no empty lobbies on this instance at all (terrible news cuz that means nobody can play due to full lobbies). maybe name it `health.allBusy` ?
- Track player counts, busy lobby counts, and other stats over time. Push to firebase periodically so this can be viewed via stats.spin4team.com later
  - in fact, this means we'd need to distinguish lobbies from servers in firebase. instead of only `flingy/servers`, we'd rename that to `flingy/lobbies`, so that `flingy/servers` will have server instances here. or better nameset: keep `flingy/servers` but add `flingy/instances`. or if i feel like refactoring, do that, but also rename`flingy/servers` into `flingy/lobbies.` man did i overexplain that 