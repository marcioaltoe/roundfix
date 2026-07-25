# QA command evidence

Build under test:

- source revision:
  `288559754c9cc0423083a30e6b63ecc98067c71b`
- `bin/roundfix` SHA-256:
  `264adbd3b8b194b8a3dc56b73764a111a63bc974741efdf1ede0131bbbd57688`
- `bin/roundfix --version`:
  `roundfix 0.0.1 (2885597-dirty, built 2026-07-24 23:19:52 -0300)`

Source repositories were clean before and after QA:

- Fluxus:
  `1aeed7e8370c3d14137c42b0c789dcbe3bd1ba3b`
- Oraculum:
  `ad74f46197500de63dc0d9ff0d3e09f61a6a43ce`
- `git -c core.fsmonitor=false status --short` returned no paths for either
  source checkout before or after the disposable-copy journeys.

## Roundfix static gate

The first `rtk make verify` attempt failed after 2,197 passing tests because
`TestFormatterComposition` inherited
`commit.gpgsign=true` from `/Users/marcio/.gitconfig` and could not sign a
temporary fixture commit.

The isolated reproduction passed:

```text
rtk proxy env GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=commit.gpgsign \
  GIT_CONFIG_VALUE_0=false go test -count=1 ./internal/baseline \
  -run '^TestFormatterComposition$'
ok roundfix/internal/baseline
```

The full unchanged gate then passed with the same process-scoped isolation:

```text
rtk proxy env GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=commit.gpgsign \
  GIT_CONFIG_VALUE_0=false make verify
Go test: 2198 passed in 22 packages
Go test: 4 passed in 1 packages
Roundfix skill check passed
roundfix build passed
```

The focused Task 01–07 acceptance suite passed all four packages, the
Task 08 documentation suite passed 28 tests in two packages, and the Task 09
real-boundary suite passed 34 tests in two packages.

## Public documentation check

The user guide labels its Decision Document a “complete greenfield input.”
Running that exact document against Fluxus exited `3`:

```text
required Baseline decisions are missing: runtime.backend, runtime.design
```

Adding those two required decisions allowed planning to exit `0`. The
documentation example is therefore parser-valid but not complete enough to
produce the Plan it introduces.

The verified Fluxus apply recommended `bun run format`, but the real repository
returned:

```text
error: Script not found "format"
```

The repository's actual formatter command, `rtk bun run fmt`, exited `0`, and
the declared repository gate, `rtk make verify`, also exited `0`.
