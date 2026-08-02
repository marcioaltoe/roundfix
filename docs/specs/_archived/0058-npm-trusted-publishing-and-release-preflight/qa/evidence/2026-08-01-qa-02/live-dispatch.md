# Live dispatch evidence

GitHub Actions run: <https://github.com/marcioaltoe/roundfix/actions/runs/30699898430>

- Event: `workflow_dispatch`
- Target branch: `ma/npm-trusted-publishing-and-release-preflight`
- Head: `04d5bbb99f63f3343b29a8ed9e7d6b93923cd93d`
- Input: `version=v0.0.2`
- Job: `release` (`91368950597`)
- Started: `2026-08-01T12:32:26Z`
- Completed: `2026-08-01T12:37:36Z`
- Workflow conclusion: `failure` — expected because the requested version is
  already used on all six coordinates.

## Step results

| Step | Result |
| --- | --- |
| Guard Trusted Publishing runtime | success |
| Validate tag | success (`preflighting 0.0.2`) |
| Verify gate | success |
| Publication preflight | failure after evaluating all six coordinates |
| Cross-compile and stage | skipped |
| Publish to npm | skipped |
| GitHub Release | skipped |

## Registry observations

| Coordinate | Version | State | Diagnostic |
| --- | --- | --- | --- |
| `@roundfix/cli-darwin-arm64` | `0.0.2` | `used` | `registry: ... is already used` |
| `@roundfix/cli-darwin-x64` | `0.0.2` | `used` | `registry: ... is already used` |
| `@roundfix/cli-linux-arm64` | `0.0.2` | `used` | `registry: ... is already used` |
| `@roundfix/cli-linux-x64` | `0.0.2` | `used` | `registry: ... is already used` |
| `@roundfix/cli-win32-x64` | `0.0.2` | `used` | `registry: ... is already used` |
| `roundfix` | `0.0.2` | `used` | `registry: ... is already used` |

The run proves the real published dispatch path reads the live npm registry,
accumulates every blocked coordinate, stops after preflight, and executes no
mutating release stage.
