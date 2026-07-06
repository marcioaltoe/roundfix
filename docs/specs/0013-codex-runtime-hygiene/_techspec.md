---
spec: 0013-codex-runtime-hygiene
prd: _prd.md
created: 2026-07-06
---

# Codex Runtime Hygiene — Technical Spec

## Executive Summary

Add a diagnosis-only Doctor Command and stop Runs from spawning a
quarantine-blocked codex. The primary trade-off is where the machine-health
checks live: the Setup Command already runs Node/acpx/agent-probe checks
through a `report(step, status, detail)` runner, so the checks are extracted
into a shared health-check unit that both Setup (which prepares) and Doctor
(which only diagnoses) consume — no duplicated probes. Codex hygiene is a new
macOS-only check (quarantine attribute + notarization acceptance) that reports
the curl-reinstall fix; on other platforms it reports not-applicable. Spawning
codex through `codex-acp` resolves a verified-clean binary so an agent loop's
repeated execs stop tripping XProtect. Detection never remediates automatically
(ADR-0032).

## System Architecture

- `internal/cli` — extract the Setup Command's per-check logic (node, acpx,
  agent probe) into a shared health-checks unit; the Setup Command keeps its
  preparing/mutating steps (install acpx, write config) and consumes the shared
  read-only checks. A new `doctor` command runs the shared checks plus the codex
  hygiene check and mutates nothing.
- New codex-hygiene inspector (macOS-gated) — resolves the codex on PATH and any
  configured codex path, inspects the `com.apple.quarantine` extended attribute
  and notarization/signing acceptance, and returns a health result with the
  curl-reinstall remediation string.
- `internal/agent` — the codex-acp launch path resolves a verified-clean codex
  (respecting a configured codex path) and spawns that binary.
- No store, journal, or TUI changes.

## Implementation Design

### Interfaces

```go
// Shared health checks consumed by both Setup and Doctor.
type CheckStatus string // "ok" | "skipped" | "failed"

type CheckResult struct {
    Name       string // "node" | "acpx" | "agent" | "codex"
    Status     CheckStatus
    Detail     string // human line
    NextAction string // remediation when Status == failed (e.g. curl reinstall)
}

type HealthChecker interface {
    Node(ctx context.Context) CheckResult
    ACPX(ctx context.Context) CheckResult
    Agent(ctx context.Context, spec agent.RuntimeSpec) CheckResult
    Codex(ctx context.Context) CheckResult // macOS: real; else Status=skipped, "not applicable"
}
```

```go
// Codex hygiene inspection (macOS).
type CodexHygiene struct {
    Path        string // resolved codex path inspected
    Quarantined bool   // com.apple.quarantine present
    Accepted    bool   // notarization/signing accepted
}
// Failure when Quarantined || !Accepted; NextAction names the
// curl-to-~/.local/bin reinstall fix. On non-darwin, the check is not run.
```

### Doctor Command

`roundfix doctor` (support command, non-interactive): run every shared check,
print one line per check (`node`, `acpx`, `agent`, `codex`) with its status and,
on failure, the next action; exit non-zero when any Run-breaking check fails.
Doctor performs no installs, no config writes, no downloads — it is the
read-only sibling of the Setup Command.

### Codex hygiene check (ADR-0032)

macOS-only. Resolve the codex the same way a Run would (configured codex path if
set, else PATH). Inspect:

- **Quarantine** — presence of the `com.apple.quarantine` extended attribute on
  the resolved binary.
- **Acceptance** — Gatekeeper/notarization assessment of the binary (a
  read-only assess; no mutation).

If quarantined or not accepted, the check fails with the curl-to-`~/.local/bin`
reinstall command as the next action. On non-darwin platforms the check returns
`skipped` with "not applicable" and never fails the command.

### Verified-clean codex on spawn (ADR-0032)

The codex-acp launch path resolves the codex binary through the same
configured-path-then-PATH order and, on macOS, prefers a binary that passes the
hygiene inspection. It never silently spawns a known-quarantined binary without
surfacing the risk. Non-macOS spawning is unchanged.

## Coverage Map

- Stories 1-2 → codex hygiene check + curl-reinstall next action (ADR-0032)
- Story 3 → Doctor Command running the shared checks
- Story 4 → verified-clean codex on the codex-acp spawn path (ADR-0032)
- Story 5 → macOS-gated check returning not-applicable elsewhere

## Integration Points

Local process/filesystem only: extended-attribute read and a read-only
Gatekeeper assessment of the codex binary on macOS; codex resolution reuses the
existing PATH/config lookup. No new external systems.

## Testing Approach

- Shared checks: the Setup Command's existing tests keep passing after the
  extraction (behavior byte-stable); the extraction is covered by reusing them
  through both Setup and Doctor entry points.
- Doctor: buffer-captured CLI tests with a fake `HealthChecker` — assert one line
  per check, the failure next-action text, and the non-zero exit when a check
  fails; a passing set exits zero.
- Codex hygiene: inject the extended-attribute and assessment probes as
  interfaces; table tests over {quarantined, not-accepted, clean, non-darwin}
  asserting status, next action, and that non-darwin is skipped. No real
  `xattr`/`spctl` in unit tests.
- Spawn resolution: a test asserts the codex-acp launch resolves the configured
  clean path over a quarantined PATH entry.

## Build Order

1. Extract shared read-only health checks (node, acpx, agent) from the Setup
   Command; Setup consumes them with byte-stable output (no deps)
2. Codex hygiene inspector: quarantine + acceptance probes behind interfaces,
   macOS-gated, with the curl-reinstall next action (ADR-0032) (no deps)
3. Doctor Command running the shared checks plus codex hygiene, diagnosis-only,
   non-zero on Run-breaking failure (depends on: 1, 2)
4. Verified-clean codex resolution on the codex-acp spawn path (ADR-0032)
   (depends on: 2)
5. Docs and skill sync (depends on: 3, 4)

## Risks & Considerations

- The macOS gate must be a build/runtime platform check, not a config toggle —
  the check is meaningless off darwin and must never fail a Linux/Windows
  Doctor run.
- Notarization assessment must be strictly read-only; no `spctl` mutation, no
  de-quarantine — Doctor diagnoses, the developer remedies (ADR-0032).
- Spawn-time resolution must not add latency to every codex exec; inspect once
  per Run's codex resolution, not per exec.
- Extracting the Setup checks risks drifting Setup's output; the extraction task
  must keep Setup's reported lines byte-stable, guarded by its existing tests.

## Decisions

- Codex hygiene is exposed through a new diagnosis-only Doctor Command, sharing
  read-only checks with the Setup Command; glossary gains Doctor Command.
- Quarantine/notarization detection is macOS-only; other platforms report
  not-applicable. See ADR-0032.
- Roundfix prefers a verified-clean codex when spawning codex-acp and never
  auto-remediates. See ADR-0032.
