# Live publish-free dispatch

GitHub Actions run: `30703974453`  
Job: `91379753174`  
Head: `171f6a378c9e640a8a10c9382e28b501b21ff5a0`  
Event: `workflow_dispatch`  
Input: `v0.0.2`

QA invoked the workflow's declared publish-free rehearsal on the exact PR head.
The run completed with the expected overall `failure` because all six
coordinates already contain version `0.0.2`:

| Coordinate | State |
| --- | --- |
| `@roundfix/cli-darwin-arm64` | used |
| `@roundfix/cli-darwin-x64` | used |
| `@roundfix/cli-linux-arm64` | used |
| `@roundfix/cli-linux-x64` | used |
| `@roundfix/cli-win32-x64` | used |
| `roundfix` | used |

`Guard Trusted Publishing runtime`, `Validate tag`, and `Verify gate` each
passed. `Publication preflight` accumulated and named all six blocked
coordinates, then exited 1. GitHub recorded `Cross-compile and stage`, `Publish
to npm`, and `GitHub Release` as skipped. Six independent `rtk npm view
<coordinate>@0.0.2 version` calls each returned `0.0.2`.

The runner warned that `actions/checkout@v4`, `actions/setup-go@v5`, and
`actions/setup-node@v4` target the deprecated Node 20 action runtime and were
forced onto Node 24. The workflow's own Node/npm guard passed, so this warning
did not change the exercised path.

The log also showed setup-node's masked `NODE_AUTH_TOKEN` placeholder in the
preflight environment. This is not the retained repository token. Current npm
Trusted Publishing documentation uses `setup-node` with `registry-url` and a
token-free `npm publish`, and states that npm detects OIDC first before
traditional token fallback: <https://docs.npmjs.com/trusted-publishers/>.

No package, asset, tag, GitHub Release, branch, commit, Pull Request, or review
thread was created or changed by the dispatch.
