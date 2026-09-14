# Lint finding evidence

Task 04 requires the repository linter to be recorded clean. The repository has no `lint` Make target and no `.golangci.yml`, so the gate used Go's pinned native analyzer:

`rtk env GOTOOLCHAIN=go1.26.7 GOCACHE=/tmp/roundfix-go-build-0134-qa go vet ./...`

The command exited 1 with 31 `copylocks` diagnostics. Every diagnostic is in `internal/agent`; representative failures are:

```text
internal/agent/acpx_runner.go:172:14: cancellationClock passes lock by value: roundfix/internal/agent.ACPXRunner contains sync.Mutex
internal/agent/acpx_runner.go:548:14: Probe passes lock by value: roundfix/internal/agent.ACPXRunner contains sync.Mutex
internal/agent/codex_spawn.go:146:14: commandEnv passes lock by value: roundfix/internal/agent.ACPXRunner contains sync.Mutex
internal/agent/fallback.go:26:14: ProbeFallback passes lock by value: roundfix/internal/agent.ACPXRunner contains sync.Mutex
internal/agent/selection_assignment.go:96:14: ProveExactSelection passes lock by value: roundfix/internal/agent.ACPXRunner contains sync.Mutex
internal/agent/selection_capabilities.go:348:14: AcquireSelectionCapabilities passes lock by value: roundfix/internal/agent.ACPXRunner contains sync.Mutex
```

This is pre-existing. A local clone checked out the pre-Task base `e9746a4543097beb4506e296411d8bd3f9c2a797`; `rtk env GOTOOLCHAIN=go1.26.7 GOCACHE=/tmp/roundfix-go-build-0134-history go vet ./internal/agent` exited 1 with the same 31 diagnostics. `git diff --name-only e9746a45..HEAD` contains no `internal/agent` path. The affected file last changed in `8a33f8fdec67a82c74592ac30c1e0e91c8aa678d`, which predates this Spec.

The failure does not invalidate the passing repository Verification or the feature-specific flows, but it makes Task 04's explicit clean-linter acceptance criterion false on this build.
