# NAS AI-Brain

Self-hosted AI-brain platform: Obsidian LiveSync through CouchDB on a Linux NAS,
reached over Tailscale and fronted by Caddy. The filesystem is the vault;
livesync-bridge is the sync transport; CouchDB is the wire format phones speak;
nothing on the NAS runs Obsidian.


## Folder layout

Copied from the Aug 3 contract, with names as they exist on disk (`docker-compose.yml`, stack folder `ai-brain/` not `brain/`).

```
docker-stack/
├── README.md
├── .gitignore                  # **/.env, *-data/, certs/, etc-pihole/,
│                               #   bridge-dat/* except config.json.example
├── scripts/
│   ├── renew.certs.sh          # tailscale cert + caddy reload (systemd timer)
│   └── backup-vault.sh         # Phase E — git commit/push of the vault
│
├── systemd/
│   ├── renew-certs.service
│   └── renew-certs.timer
│
├── edge/            Caddy — owns `web`, holds certs, routes everything
│   ├── docker-compose.yml
│   ├── Caddyfile
│   ├── .env.example
│   └── certs/                  # directory mount (survives renewal), gitignored
│
├── ai-brain/        CouchDB 3 + livesync-bridge (GUI retired)
│   ├── docker-compose.yml      # couchdb on `web`; bridge on brain_internal
│   ├── .env.example            # COUCH_USER, COUCH_PASSWORD, VAULT_PATH
│   ├── couchdb-data/           # gitignored
│   ├── bridge-src/             # headless replicator image
│   └── bridge-dat/             # runtime + live config.json (creds) · gitignored
│       └── config.json.example # committed: couchdb peer + storage peer
│
├── brain-api/       one Go module, three entrypoints (in progress)
│   ├── docker-compose.yml
│   ├── dockerfile
│   ├── justfile
│   ├── cmd/{api,sorter,adminctl}/
│   └── internal/{vault,sort,config}/
│
├── media/           gluetun + qbittorrent + jellyseerr/jellyseerr + *arr
│   ├── docker-compose.yml
│   ├── .env.example
│   └── failback/               # separate Go module, Phase 8
│
├── pihole/          network-wide DNS + ad-blocking
│   ├── docker-compose.yml      # host :53 — the one sanctioned host-port exception
│   ├── .env.example
│   └── etc-pihole/             # gitignored
│
└── shared/          public upload frontend — Funnel + forward-auth (later)
```

Storage, not stacks — outside this repo:

```
/mnt/nas/
├── vault/           written by livesync-bridge, sorter, brain-api
├── inbox/           sorter dropzone; quarantine/ lives here, never under vault/
├── models/
├── media-library/
└── shared-uploads/
```

On this machine the live vault mount is currently `/home/humbertx/vault` (`VAULT_PATH` in `ai-brain/.env`). The contract path remains `/mnt/nas/vault`.

## Networks

| Stack | Contents | Networks | Exposure |
| --- | --- | --- | --- |
| `edge` | Caddy + certs | `web` (owner) | Tailnet `:80`/`:443`/`:5984`/`:8089` |
| `ai-brain` | CouchDB 3 + livesync-bridge | CouchDB: `web` · bridge: `brain_internal` (+ `web` today) | Caddy `:5984` only |
| `brain-api` | api · sorter · adminctl | `web` only (planned) | Caddy; sorter has none |
| `media` | gluetun + qbit + seerr + *arr | host / compose ports (to fold behind Caddy) | private |
| `pihole` | DNS + DHCP | `web` (UI) + host `:53` (and `:67` if DHCP) | LAN DNS; UI via Caddy `:8089` |
| `shared` | Immich / uploads | `shared_edge` only — never joins `web` | public, last |

Rules:

1. `edge` owns `web` with `name: web`. Every other stack joins it as `external: true`. Bring `edge` up first.
2. Private data paths get an `internal: true` network. CouchDB↔bridge traffic should not need the proxy.
3. One `docker-compose.yml` and one `.env` per directory. Commit `.env.example`, never `.env`.
4. `certs/` is a directory mount, never individual files (inode trap on renewal).
5. `bridge-dat/` is gitignored; `config.json.example` is committed. Live `config.json` holds CouchDB creds and the obfuscation passphrase.

## Pipeline

```
phone/iPad Obsidian --(Caddy :5984, HTTPS)--> CouchDB <--(brain_internal)-- livesync-bridge
                                                                                    |
                                                                             VAULT_PATH
                                                                          (plaintext .md)
                                                                                    |
                                                   brain-api · sorter · git backup · Filesystem MCP
```

Phones still hit `https://<tailnet-host>:5984`. The bridge has no inbound HTTP surface.

On the CouchDB peer, `baseDir` only filters. The storage peer folder is the real write target.

## Phase C — livesync-bridge (Aug 22, 2026)

GUI container is gone. Bridge is up (`livesync-bridge`, 200 MB cap).

| Check | Result |
| --- | --- |
| End-to-end encryption | Off |
| Obfuscate Properties | **On and locked** after remote init. There is no separate obfuscation field in the plugin. Put the vault **Passphrase** into `bridge-dat/config.json` as `obfuscatePassphrase`. Leave `passphrase` empty while E2EE is off. |
| Forward (Obsidian → disk) | Passed. After `Replicate now`, notes land under `VAULT_PATH`. Bridge user is UID 1993 (`deno`); the vault dir needs `u:1993:rwX` (ACL) or equivalent. Existing CouchDB docs are not replayed after a restart (`Watch starting from now`). |
| Reverse (disk → phone) | Not recorded yet |
| Conflict (phone + disk same note) | Not recorded yet |

Bring-up:

```bash
cd ~/docker-stack/edge && docker compose up -d
cd ~/docker-stack/ai-brain && docker compose up -d
```

Do not start the bridge until the first Obsidian client has finished initializing a new remote.

## brain-api / `internal/vault`

Filesystem is authoritative. `internal/vault` is a plain-file client (read/write markdown + frontmatter). No CouchDB chunk format, no Obsidian REST plugin, no `OBSIDIAN_*` env. Env is `VAULT_PATH` + `INBOX_PATH`. `COUCHDB_URL` is deferred until something actually needs sync-state inspection.

Writes into the vault must be atomic (write-temp-then-rename) so the bridge never syncs a half-written note. `quarantine/` stays under `inbox/`, never under `vault/`.
