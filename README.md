# Flynn Redis plugin

Redis resource provider for Flynn. Layout follows
[`flynn-plugin-template`](https://github.com/randy-girard/flynn-plugin-template).

## Install

On a cluster host:

```text
flynn-host plugin install redis
flynn-host plugin install /path/to/flynn-plugin-redis
flynn-host plugin install https://github.com/randy-girard/flynn-plugin-redis.git --ref vX
```

The user `flynn` CLI does not install plugins. After install, that cluster's
CLI catalog lists `redis` (`redis-cli`, `dump`, `restore`). Provision a database
for an app with:

```text
flynn resource add redis
```

## Image artifacts

`./script/plugin-build` (and the manual **Build and Release** workflow) writes:

- `dist/image.json` — Flynn Artifact (`type: flynn`)
- `dist/<manifest-id>.json` — ImageManifest at `artifact.uri`
- `dist/layers/<layer-id>.squashfs` — zstd squashfs rootfs
- `dist/flynn-plugin.json` — manifest with `artifacts.image` set to the GitHub
  Release URL for `image.json`

Release asset URLs look like:

```text
https://github.com/<owner>/flynn-plugin-redis/releases/download/<tag>/image.json
https://github.com/<owner>/flynn-plugin-redis/releases/download/<tag>/{id}.squashfs
```

The image stacks Flynn's **ubuntu-noble** layer (from a Flynn GitHub Release) with
`redis-server` + `flynn-redis` / `flynn-redis-api`. Entrypoint is
`/bin/start-flynn-redis` (the API process adds `api`). Child Redis apps use the
same image via `REDIS_IMAGE_ID=self`.

## Develop

```bash
./script/run-unit-tests
./script/plugin-build
```

On Linux those run natively. On macOS they use Docker Desktop (linux/amd64,
privileged), same idea as Flynn `script/run-unit-tests`, so process tests see
`redis-server` and `plugin-build` can overlay Flynn's ubuntu-noble + mksquashfs.

Pin the Flynn OS with `-flynn-version vYYYYMMDD.N` or `build.base.version` in
`flynn-plugin.json` (default is the latest published Flynn GitHub release).
`PLUGIN_BUILD_DOCKER=0` forces native. GitHub **Build and Release** is still the
publish path (manual, like Flynn).

Pin Flynn APIs with the `replace` in `go.mod` (`github.com/randy-girard/flynn`).

## GitHub Actions

CI runs `gofmt` and `go test` on push/PR. Image builds are **manual**: run
**Build and Release**, pass a version like `v20260914.0` (same scheme as Flynn).
The default is a draft/prerelease.

## Safety

No HA guarantees. Treat data as ephemeral (cache / dev / test). See
[docs/redis.md](docs/redis.md).
