---
task: task_01
spec: 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
status: pending
type: backend
complexity: medium
---

# Task 01: The profile proof's setup deadline comes from config

## Overview

On 2026-10-19 and 2026-10-20 the Delivery Queue parked five times because a
profile proof's fixed 30-second setup deadline expired twice (`_prd.md` →
Sources 1). This Task makes that deadline the config key
`defaults.profile_proof_setup_timeout` and carries it from every
profile-readiness consumer to the agent's setup context, and documents the
key in the configuration guide and the Roundfix `profiles` reference. It is
verifiable on its own: an injected runner sees the configured value, the fake acpx's codex
hygiene probe sees the deadline on the setup context, and `roundfix profiles
show` refuses a value outside the range.

## Requirements

1. MUST answer the Secondbrain Inbox Entry of 2026-10-20,
   [`references/2026-10-20-a-prova-de-perfil-do-preflight-estoura-o-tempo-em-serie.md`](references/2026-10-20-a-prova-de-perfil-do-preflight-estoura-o-tempo-em-serie.md)
   (`_prd.md` → Sources 1), by implementing API Contract 1 of `_techspec.md`
   with the shapes `_techspec.md` → Interfaces sketches (ADR-0298).
2. MUST add `defaults.profile_proof_setup_timeout` to `internal/config`:
   - a Go duration read from User Config and Project Config, Project Config
     winning, exported as `Defaults.ProfileProofSetupTimeout`;
   - default `DefaultProfileProofSetupTimeout`, `30s`, and the agent's
     `acpxPreflightSetupTimeout` defined from that constant, so the value is
     written once;
   - `Validate` refuses a value not greater than zero or greater than
     `MaxProfileProofSetupTimeout` (`10m`) with exactly
     `defaults.profile_proof_setup_timeout must be greater than 0 and at most
     10m`.
   The `roundfix init` template and this repository's `.roundfixrc.yml` do
   not change.
3. MUST add `SetupTimeout time.Duration` to `agent.ProbeRequest`.
   `ProveExactSelection` bounds its setup with it, using
   `acpxPreflightSetupTimeout` when it is not positive, and `ACPXRunner.Probe`
   forwards it. The cleanup deadline, the Fallback Chain probe in
   `fallback.go` and Run sessions keep their values.
4. MUST set `SetupTimeout` from `config.Defaults.ProfileProofSetupTimeout` in
   every `ProbeRequest` that `proveProfileSelectionsWithOptions` builds, so
   `profiles validate`, the operational Preflight, the pre-PR review, Doctor,
   Setup and `profiles configure` all pass it. State the hazard as API
   Contract 1 does: the value comes from User and Project Config and reaches
   only the setup context of each proof attempt; an expired deadline stops
   the acpx command and closes the disposable Session, as today.
5. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows and waiting on no real time:
   - `internal/config/profile_proof_setup_timeout_test.go`:
     `TestProfileProofSetupTimeoutConfig`, with the subtests of
     `_techspec.md` → Testing Approach (built-in value, User Config `45s`,
     Project Config `90s` over User Config `45s`, `0s`, `-1s` and `11m`
     refused with the exact message, `soon` refused with `invalid duration`);
   - `internal/agent/proof_setup_deadline_test.go`:
     `TestExactSelectionProofSetupUsesTheRequestedDeadline`, reusing the fake
     acpx harness and a codex hygiene probe that records the time left until
     the setup context's deadline: `SetupTimeout` of 2 h leaves more than 1 h
     on every recorded context through `ProveExactSelection` and through
     `Probe`; no `SetupTimeout` leaves more than zero and at most
     `acpxPreflightSetupTimeout`;
   - `internal/cli/profile_proof_setup_deadline_test.go`:
     `TestProfileProofPassesTheConfiguredSetupDeadline` (Project Config
     `90s` reaches every tuple's `ProbeRequest` in `profiles validate`;
     without it, `DefaultProfileProofSetupTimeout`) and
     `TestProfilesShowRefusesAnInvalidSetupDeadlineTranscript`, asserting
     Surface Transcript 1 of `_techspec.md` whole: empty standard output,
     both standard error lines in order and exit `2`.
6. MUST document the key, because a config key and its covering skill
   change together (ADR-0256):
   - `docs/user-guide/configuration.md`: a row for
     `defaults.profile_proof_setup_timeout` in the General settings table
     (default `30s`, a duration greater than zero and at most 10m, User and
     Project Config, Project Config winning) and, under Agent selection
     profiles, a short subsection stating that it bounds the disposable Agent
     Session setup of each proof attempt;
   - `.agents/skills/roundfix/references/profiles.md`: a new section on the
     proof's setup deadline with the key, its default `30s`, its range "at
     most 10m", and that every configured tuple is still proven.
   The attempts and the backoff are documented by task_04. No text inside
   any `### QA settlement` section may change.
7. MUST run these commands, because this Task edits a user guide page, a
   Roundfix reference and the config keys the Behavior Surface Record
   fingerprints. Before them, raise both version fields of the Roundfix
   skill's canonical `SKILL.md` (`metadata.version` and the top-level
   `version`) by one patch over the value the tree holds when this Task
   starts, because the record command refuses changed content under an
   unchanged version. The Agent runs these commands itself before handing
   back: the Daemon only runs the Verification and never regenerates a
   mirror or a record. No record or mirror is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Add the config key, its default and its validation.
- [ ] Add `ProbeRequest.SetupTimeout` and use it in the setup context and `Probe`.
- [ ] Pass the configured value from every profile-readiness consumer.
- [ ] Add the config, agent and CLI tests, including Surface Transcript 1.
- [ ] Document the key in the configuration guide and the Roundfix `profiles` reference.
- [ ] Raise the skill version, sync and record.

## Acceptance Criteria

- [ ] Every proof attempt's setup context carries the configured deadline, or 30 s when none is configured.
- [ ] User and Project Config set the key with the usual precedence, and a value outside the range is refused with the exact message.
- [ ] `roundfix profiles show` reproduces Surface Transcript 1.
- [ ] The cleanup deadline, the Fallback Chain probe and Run sessions keep their values.
- [ ] The configuration guide and the Roundfix `profiles` reference name the key, its default and its range.

## Consumers

- `defaults.profile_proof_setup_timeout` and `agent.ProbeRequest.SetupTimeout` — consumers: `./internal/config`, `./internal/agent`, `./internal/cli`
- the `config keys` Behavior Surface fingerprint and the Roundfix `profiles` reference — consumers: `./internal/docscontract`, `./skills`

## Context

- instruction: `docs/adr/0298-a-timed-out-profile-proof-gets-a-third-attempt-and-its-setup-deadline-is-configurable.md`
- instruction: `internal/agent/fallback.go`
- instruction: `internal/cli/profile_proof_retry_test.go`
- interface: `internal/config/config.go`
- creates: `internal/config/profile_proof_setup_timeout_test.go`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/selection_assignment.go`
- creates: `internal/agent/proof_setup_deadline_test.go`
- instruction: `internal/agent/acpx_runner_test.go`
- interface: `internal/cli/profiles_validate.go`
- creates: `internal/cli/profile_proof_setup_deadline_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/references/profiles.md`
- interface: `skills/roundfix/references/profiles.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^TestProfileProofSetupTimeoutConfig$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestProfileProofSetupTimeoutConfig ' || { printf 'missing pass: TestProfileProofSetupTimeoutConfig\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestExactSelectionProofSetupUsesTheRequestedDeadline|TestACPXProbeSelectionSetupUsesBoundedContext)$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestExactSelectionProofSetupUsesTheRequestedDeadline TestACPXProbeSelectionSetupUsesBoundedContext; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestProfileProofPassesTheConfiguredSetupDeadline|TestProfilesShowRefusesAnInvalidSetupDeadlineTranscript|TestProfileProofRetriesATimedOutSetupOnce|TestProfileProofDoesNotRetryARejectedSelection|TestProfileProofDoesNotRetryACancelledCommand)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestProfileProofPassesTheConfiguredSetupDeadline TestProfilesShowRefusesAnInvalidSetupDeadlineTranscript TestProfileProofRetriesATimedOutSetupOnce; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/config ./internal/agent ./internal/cli || exit 1; GOOS=windows go vet ./internal/config ./internal/agent ./internal/cli || exit 1; go test -count=1 ./internal/config || exit 1; for file in docs/user-guide/configuration.md .agents/skills/roundfix/references/profiles.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'defaults.profile_proof_setup_timeout' || { printf 'missing the key in %s\n' "$file" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'at most 10m' || { printf 'missing the range in %s\n' "$file" >&2; exit 1; }; done; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check. Afterwards the config, agent and CLI tests pass with the existing setup-context and retry tests, the three packages and their tests compile for Windows, the config package passes whole, the configuration guide and the Roundfix `profiles` reference name the key and its range, and the skill mirror, the skill version, the Behavior Surface Record and the documentation contracts are current.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Core Feature 4; Success Metric 3
- `_techspec.md` → API Contract 1; API Contract 4; Surface Transcript 1; Invariants 1-2; Interfaces; Data Models; Current behavior; Testing Approach; Build Order 1
- ADR-0298; ADR-0050; ADR-0140; ADR-0107; ADR-0184; ADR-0259; ADR-0271; ADR-0285; ADR-0189; ADR-0256; ADR-0280
