# Flingy.io

> DISCLAIMER: This README is temporarily written with AI, for lack of time lol.

> I did not vibe code this. Not only is that boring, but I already played with fire by vibecoding once, and I never want to do that ever again.

> I did use Claude Code to spot bugs, do minor things (not coding), find shit, and teach me Go/html/css cuz all I knew beforehand was js.

Go game server + JS/canvas client.

## Layout

- `client/` — HTML/JS client (served separately in production)
- `server/` — Go game server
- `lobbies.json` — lobby configs created at startup
- `.env` — secrets and URLs (copy `.env.example`, never commit)

## Requirements

- Go 1.27+
- `cloudflared` (production only)

## Local development

```bash
cp .env.example .env    # fill in the Firebase values
python3 -m http.server 8000 -d client
go run ./server
```

Open `http://localhost:8000`.

`go run ./server` compiles the package in `server/` and runs it immediately, leaving no binary
behind. Run it from the repo root — it reads `lobbies.json`, `.env` and `serviceAccountKey.json`
relative to the working directory.

The server listens on `127.0.0.1:4096` and only handles `/ws`. The client is served separately.

To test against a server running on your own machine, set `WS_URL=ws://localhost:4096/ws` in
`.env` — that URL is what the server writes into its lobby entries, and it's where the client
will connect.

## Production (single machine)

```bash
go build -o flingy ./server
./flingy
```

Keep it running under `tmux`/`screen`/systemd. Then start the tunnel:

```bash
cloudflared tunnel run --url http://127.0.0.1:4096 flingy
```

No ports are forwarded; the tunnel dials out to Cloudflare.

## Domain layout

| Host | Points at |
|---|---|
| `flingy.io` | Vercel (client) |
| `ws.flingy.io` | tunnel → machine 1 :4096 |
| `ws2.flingy.io` | tunnel → machine 2 :4096 |

Each game server sets its own `WS_URL` in `.env`, and clients connect to whatever URL a lobby
advertises. Players only ever visit `flingy.io`.

## Adding another game server machine

Run this on the new machine. Each machine needs its own tunnel — never run the same tunnel on
two machines, Cloudflare will balance connections between them randomly.

```bash
# install cloudflared (Debian/Ubuntu)
sudo mkdir -p --mode=0755 /usr/share/keyrings
curl -fsSL https://pkg.cloudflare.com/cloudflare-main.gpg | sudo tee /usr/share/keyrings/cloudflare-main.gpg >/dev/null
echo "deb [signed-by=/usr/share/keyrings/cloudflare-main.gpg] https://pkg.cloudflare.com/cloudflared any main" | sudo tee /etc/apt/sources.list.d/cloudflared.list
sudo apt update
sudo apt install cloudflared

# authenticate and create its own tunnel
cloudflared tunnel login
cloudflared tunnel create flingy2

# point a hostname at it
cloudflared tunnel route dns flingy2 ws2.flingy.io

# run it
cloudflared tunnel run --url http://127.0.0.1:4096 flingy2
```

If the new machine has no browser, do `tunnel login` / `tunnel create` / `route dns` on a machine
that does, then copy `~/.cloudflared/cert.pem` and `~/.cloudflared/<UUID>.json` over and run
`cloudflared tunnel run flingy2`.

Then on that machine:

```bash
git clone <repo>
cd Flingy.io
cp .env.example .env    # WS_URL=wss://ws2.flingy.io/ws
go build -o flingy ./server
./flingy
```

## Conventions

- `OriginPatterns` in the server must list every site that hosts the client
  (`flingy.io`, `beta.flingy.io`, `localhost:*` for dev). It matches the page's `Origin`
  header, not the WebSocket hostname.
