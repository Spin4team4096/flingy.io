# TODO
- everything

### Bugs
- there is no duplicate connection check. a client can just: connect multiple times lol.
  - in fact, there's no checks at all. a player can just join midgame rn.

### Packets
- send statetick so the client can calc this for himself

### Commands
- `status` that lists the servers AND their status (lagging/joinable/full), or `status <serverID>` for only that server's status
- `stats` that measures and gives avg tps and mspt etc, and `stats <serverID>` for that
- give command finish feedback on those currently silent commands, like `unregister all` would log `unregistered 7 lobbies` for example.
- `list <parameter>`, defaults to assuming `local`, but you can also list everything in `firebase`

### Random Ideas
- Have a server health tracking system, like a command `health` and a struct that contains health info;
  - it also pings @spin4team4096 on discord when all lobbies are ongoing (will need to track this somehow) and there are no empty lobbies on this instance at all (terrible news cuz that means nobody can play due to full lobbies). maybe name it `health.allBusy` ?
- Track player counts, busy lobby counts, and other stats over time. Push to firebase periodically so this can be viewed via stats.spin4team.com later
  - in fact, this means we'd need to distinguish lobbies from servers in firebase. instead of only `flingy/servers`, we'd rename that to `flingy/lobbies`, so that `flingy/servers` will have server instances here. or better nameset: keep `flingy/servers` but add `flingy/instances`. or if i feel like refactoring, do that, but also rename`flingy/servers` into `flingy/lobbies.` man did i overexplain that 