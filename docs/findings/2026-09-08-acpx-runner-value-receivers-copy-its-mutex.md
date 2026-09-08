---
status: pending
created_at: 2026-09-08
updated_at: 2026-09-08
kind: finding
---

# Agent state — value receivers copy the runner mutex (2026-09-08)

The originating session recorded 31 `go vet` copylock diagnostics while the
repository gate passed. Fresh verification on 2026-09-08 reproduced 31
diagnostics: current source still combines a mutex field with value receivers.

Source: Secondbrain `inbox/roundfix/_triaged/2026-08-30-go-vet-acha-31-copylocks-que-o-gate-do-repo-nao-ve.md`.
The original capture and its dated evidence remain intact there.

## 1. Observed behavior

- Symptom / evidence: `internal/agent/acpx_runner.go` contains `stateMu sync.Mutex` and value receivers including `cancellationClock`, `Probe`, `probeACPX`, and `runACPXCommand`. The current Makefile verification composition has no `go vet` step. The fresh command below exited 1 with 31 `passes lock by value` diagnostics across `acpx_runner.go`, `codex_spawn.go`, `fallback.go`, `selection_assignment.go`, and `selection_capabilities.go`.
- Root cause: The value receiver copies the struct containing the mutex. Whether a particular copied lock protects a live mutation was not established; the source-level copying and the missing gate step are separate observations.
- Action / suggestion: Route state ownership to provisional P4, runtime readiness and model capabilities. Correct existing diagnostics before an expressly authorized Verification change adds coverage. Do not suppress `copylocks`.

The implementation group is provisional; no Spec has been assigned and no
implementation or terminal verification is claimed by this triage.

## Verification — 2026-09-08

```sh
rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go vet ./...
```

Exit 1, with 31 copylock diagnostics. This establishes a current analyzer
failure; it does not prove the copied mutex caused a particular production
incident. No suppression or receiver change was made during triage.
