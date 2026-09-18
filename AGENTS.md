# Agent notes

This is the Flynn **Redis** resource-provider plugin (`kind: resource-provider`). Layout follows [flynn-plugin-template](https://github.com/randy-girard/flynn-plugin-template). Go 1.24, `GOFLAGS=-mod=mod`. Images are Linux/amd64 squashfs layered on Flynn’s published ubuntu-noble (GitHub `images.json.gz`).

## Tests are required

Do not land behavior without tests in the **same change**.

- Process/API/client changes: `process_test.go`, `process_persist_test.go`, `cmd/flynn-redis/*_test.go`, `cmd/flynn-redis-api/*_test.go`.
- `cmd/plugin-build` (manifest, Flynn base selection, layer verify): `cmd/plugin-build/*_test.go`.
- Run `./script/run-unit-tests` (native on Linux; Docker on macOS, image includes `redis-server`). `gofmt -s` must be clean. Unit tests write HTML coverage under `coverage/` (gitignored).
- Process tests that exec `redis-server` must skip cleanly when the binary is missing, and must run in the Docker wrapper.

Skip tests only when the change cannot regress (typo in comments, LICENSE). Say so in the commit body.

## Semantic git commits

Use [Conventional Commits](https://www.conventionalcommits.org/):

```text
<type>(optional-scope): <imperative summary>
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `perf`, `ci`, `build`, `chore`.

- One concern per commit. Do not mix a feature and an unrelated cleanup.
- Subject is why it matters, not a file list.
- Put tests in the same commit as the behavior they cover (`feat`/`fix` with tests), not a later `test:` dump unless the commit is tests-only.
- Do not commit `dist/`, `.plugin-build-cache/`, or `script/docker/dev/.image-built`.

Examples:

```text
feat: overlay redis-server on Flynn ubuntu-noble
fix: skip process tests when redis-server is not installed
test: cover Flynn base layer selection from images.json
docs: document FLYNN_VERSION pinning for plugin-build
```

Commit when asked. Push only when asked.

## Plugin contract

- Do not rebuild Ubuntu from a cloud image; `plugin-build` must pull Flynn’s ubuntu-noble layer.
- Pin Flynn with `build.base.version` / `-flynn-version` for published releases; `latest` is for local builds.
- Import Flynn APIs as `github.com/randy-girard/flynn/...`. `go.mod` must `require github.com/randy-girard/flynn`. Do not vendor Flynn and do not `replace` it with a sibling `../flynn`.
- Keep `redis-server` in this plugin’s `img/packages.sh`, not in Flynn’s shared OS layer.
- `flynn-host plugin install` is operator-only; the user `flynn` CLI does not install plugins.
