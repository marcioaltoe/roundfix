# Static gate evidence

Build: `75161e9c3a5f7554cd1e0b9290bce6c61820b5c7`.

Command:

```text
rtk make verify
```

Result: exit `2`, environment-blocked.

```text
Go test: 2848 passed, 5 failed, 2 skipped in 24 packages
cli: TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion
store: TestOwnerProcessControllerMatchingOwnerIdentityProceeds
store: TestOwnerProcessControllerProveOwnerLeavesProvenOwnerRunning
store: TestOwnerProcessIdentityIsStableForOneProcess
store: TestOwnerProcessIdentityIgnoresCallerTimezone
fork/exec /bin/ps: operation not permitted
```

All five failures share the macOS process-identity seam and fail because this
ACP sandbox denies `/bin/ps`. The gate reached 2,848 passing tests through the
repository-local cache; it did not fail on formatting, a product assertion, or
Go-cache access. A full-access session must rerun the exact command.

The raw bare-build probe without Make reached the same host-cache sandbox
boundary:

```text
rtk go build -buildvcs=false ./cmd/roundfix
open /Users/marcio/Library/Caches/go-build/...: operation not permitted
```

Repeating that build with the repository-local cache succeeded twice and left
Git status unchanged apart from QA artifacts. This parity deviation affects
only the sandboxed raw-Go probe; the selected repository Verification owns its
local cache default.
