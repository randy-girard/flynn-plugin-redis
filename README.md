# Flynn Redis plugin

[![coverage](.github/badges/coverage.svg)](https://github.com/randy-girard/flynn-plugin-redis/actions/workflows/ci.yml)

Redis resource provider for Flynn. Layout follows
[`flynn-plugin-template`](https://github.com/randy-girard/flynn-plugin-template).

Redis comes from the Ubuntu 24.04 package set and runs as a **single process**.
Treat data as ephemeral (cache, development, test). There are no HA guarantees.

## Install

On a cluster host. A local checkout is optional: the `redis` alias pulls
`https://github.com/randy-girard/flynn-plugin-redis` when no sibling dir exists.

```text
sudo flynn-host plugin:install redis --ref vX
sudo flynn-host plugin:install https://github.com/randy-girard/flynn-plugin-redis.git --ref vX
sudo flynn-host plugin:install /path/to/flynn-plugin-redis
sudo flynn-host plugin:update redis --ref vX
sudo flynn-host plugin:uninstall redis
sudo flynn-host plugin:uninstall redis --force
```

`--ref` is a published GitHub Release tag from **Build and Release**. Override
the org with `--github-org`, `FLYNN_PLUGIN_GITHUB_ORG`, or `/etc/flynn/plugins.json`:

```json
{
  "github_org": "randy-girard",
  "redis": {
    "url": "https://github.com/randy-girard/flynn-plugin-redis.git",
    "ref": "vX"
  }
}
```

Private or draft releases need `flynn-host plugin:credentials-set github` (or
`FLYNN_PLUGIN_GITHUB_TOKEN`). The user `flynn` CLI does not install plugins.
After install, that cluster's CLI catalog lists `redis` (`redis-cli`, `dump`,
`restore`) from the plugin manifest. The `flynn` binary does not compile those
commands in; it fetches usage from the cluster and runs them as jobs.
`plugin:update` deploys a new release without re-asking setup prompts.
Uninstall with `flynn-host plugin:uninstall redis`. Resource-provider uninstall
refuses while other apps still use provisioned resources unless `--force`.

## Usage

After the plugin is installed, provision a database for an app:

```text
flynn resource:add redis
```

That starts a Redis server as a Flynn app and configures the application to
connect to it. The app release gets `FLYNN_REDIS` (the Redis app name),
`REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, and `REDIS_URL` (some libraries
use the URL).

### Commands

After install the laptop fetches these from the cluster catalog. Colon form is
canonical; space form also works.

| Command | Purpose |
| --- | --- |
| `flynn redis:cli` / `flynn redis redis-cli` | Open `redis-cli` on the cluster (any extra args after `--`) |
| `flynn redis:dump [-q] [-f <file>]` | RDB dump to a file or stdout |
| `flynn redis:restore [-q] [-f <file>]` | Restore a dump from a file or stdin |

```text
flynn redis:cli
flynn redis:cli -- PING
flynn redis:dump -f db.dump
flynn redis:restore -f db.dump
```

The console uses `REDIS_HOST` (`leader.<redis-app>.discoverd`), the same host as
`REDIS_URL`. No local Redis client or firewall change is required.

To reach Redis from outside the cluster, add a TCP route to the leader:

```text
flynn -a $(flynn env:get FLYNN_REDIS) route:add tcp --service $(flynn env:get FLYNN_REDIS) --leader
```

Firewall that port. Use it only over the local network, a VPN, or an SSH tunnel.

## Image artifacts

`./script/plugin-build` (and the manual **Build and Release** workflow) writes:

- `dist/image.json` — Flynn Artifact (`type: flynn`)
- `dist/<manifest-id>.json` — ImageManifest at `artifact.uri`
- `dist/layers/<layer-id>.squashfs` — Flynn ubuntu-noble (local overlay) plus the Redis delta
- `dist/flynn-plugin.json` — manifest with `artifacts.image` set to the GitHub
  Release URL for `image.json`

Release assets are `image.json`, hook scripts, and the **Redis delta** squashfs.
Flynn ubuntu-noble is not re-uploaded; hosts fetch that layer from the Flynn
GitHub Release named in `flynn.plugin.base` (same layer id).

```text
https://github.com/<owner>/flynn-plugin-redis/releases/download/<tag>/image.json
https://github.com/<owner>/flynn-plugin-redis/releases/download/<tag>/{delta-id}.squashfs
```

The image stacks Flynn's **ubuntu-noble** layer (from a Flynn GitHub Release;
`images.json.gz` may list it only as postgres/gitreceive layer 0, not as a
named `ubuntu-noble` image) with
`redis-server` + `flynn-redis` / `flynn-redis-api`. Entrypoint is
`/bin/start-flynn-redis` (the API process adds `api`). Child Redis apps use the
same image via `REDIS_IMAGE_ID=self`.

## Develop

```bash
./script/run-unit-tests
./script/plugin-build
```

Unit tests write HTML coverage under `coverage/` (gitignored): open `coverage/index.html`. Set `PLUGIN_SKIP_COVERAGE=1` to skip the report.


On Linux those run natively. On macOS they use Docker Desktop (linux/amd64,
privileged), same idea as Flynn `script/run-unit-tests`, so process tests see
`redis-server` and `plugin-build` can overlay Flynn's ubuntu-noble + mksquashfs.

Pin the Flynn OS with `-flynn-version vYYYYMMDD.N` or `build.base.version` in
`flynn-plugin.json` (default is the latest published Flynn GitHub release).
`PLUGIN_BUILD_DOCKER=0` forces native. GitHub **Build and Release** is still the
publish path (manual, like Flynn).

Flynn APIs come from `go.mod` (`require github.com/randy-girard/flynn`). Do not vendor Flynn.

## GitHub Actions

CI runs `gofmt`, release-note checks, and `go test` on push/PR. Image builds are
**manual**: run **Build and Release**, pass a version like `v20260914.0.0` (`vYYYYMMDD.N.P`; Flynn stays `vYYYYMMDD.N`).. The default is a published GitHub Release (not draft, not
prerelease). Notes group conventional commits the same way Flynn does, with a
Full Changelog compare link and install commands. Turn on **draft** or
**prerelease** only if you want those GitHub flags.
