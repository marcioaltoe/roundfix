---
task: task_15
spec: 0119-spec-contained-authorization
status: pending
type: backend
complexity: medium
---

# Task 15: Withdraw the execution boundary from this Spec

## Overview

The maintainer promoted the execution trust boundary in full to its own Spec on
2026-09-10, after the terminal QA gate proved it carries a performance
regression on top of two adversarial bypasses. This slice removes it from this
Spec's delivery so the record half can close. It is a removal, not a rewrite:
Task 06's delivery is superseded by a scope decision, and authored commands must
execute exactly as they did before this Spec touched them.

## Requirements

1. MUST remove the authored-command source decision and every production caller
   of it, so no path consults it before executing an authored command: the
   pre-dispatch probe, the post-Agent verification call, the Settle path, and
   the read-only checker's probing path.
2. MUST restore each of those call sites to the behavior it had before this
   Spec, so an authored command executes exactly as it did at the delivery
   target's revision. Do not leave a disabled flag, a dead parameter, or a
   check that always permits.
3. MUST NOT spawn a subprocess per authored command anywhere in the Task cycle
   after this change, which is the regression being withdrawn.
4. MUST preserve the record half completely: the typed reader in
   `internal/authorization`, the operation vocabulary and every boundary that
   asks it, the constraint reader's refusal, the Governed Path set, the
   changed-path audit's resolved reference, and the suite guard's single
   parser. Removing the execution boundary must not weaken any of them.
5. MUST remove the `SC-SOURCE-UNTRUSTED` token and scope its glossary entry to
   the code that remains, so the glossary does not define a token the CLI
   cannot emit.
6. MUST leave the complete repository Verification passing, including the
   Daemon Task-cycle cases that miss their deadlines today, run as the full
   concurrent suite rather than in isolation.
7. MUST NOT change any Task status, and MUST NOT edit a sibling Task file.
   Task 06 completed; its delivery is withdrawn by this scope decision, and the
   PRD records that.

## Subtasks

- [ ] Remove the source decision and its four production call sites.
- [ ] Restore each call site's pre-Spec execution behavior.
- [ ] Remove the coined execution token and scope the glossary entry.
- [ ] Prove the record half is untouched.
- [ ] Prove the full concurrent suite passes.

## Acceptance Criteria

- [ ] No production file references the authored-command source decision, and
      the file that implemented it is gone rather than orphaned.
- [ ] No production code under `internal/daemon` or `internal/speccheck` spawns
      a process per authored command; a search for the spawn site finds none in
      the execution path.
- [ ] An authored command from a modified or untracked Spec artifact executes,
      as it did before this Spec, because no execution gate remains.
- [ ] Every record-half contract still passes: the typed reader, the operation
      boundaries, the constraint refusal, the governed set, the audit's resolved
      reference, and the suite guard's single parser.
- [ ] `CONTEXT.md` defines only the refusal code the CLI still emits.
- [ ] The complete Verification passes as one concurrent run, including the four
      Daemon Task-cycle cases named in the QA report's F-005.

## Context

- interface: `internal/speccheck/verification_source.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/verification_probe.go`
- interface: `internal/cli/settle.go`
- interface: `internal/cli/spec_check.go`

## Verification

- `matches="$(grep -rln 'AuthorizeAuthoredCommands\|ProbeAuthoredCommands\|SC-SOURCE-UNTRUSTED' internal/ --include='*.go' || true)"; test -z "$matches" || { printf '%s\n' "$matches"; exit 1; }` — no production or test file references the withdrawn boundary; it fails today.
- `test ! -f internal/speccheck/verification_source.go` — the implementation file is removed rather than left orphaned.
- `grep -q 'SC-SOURCE-UNTRUSTED' CONTEXT.md && exit 1; grep -q 'SC-TOOLING-UNAPPROVED' CONTEXT.md` — the glossary defines the code that remains and no longer defines the one the CLI cannot emit.
- `test ! -f internal/speccheck/verification_source.go || exit 1; go test -count=1 ./internal/authorization ./internal/spec ./internal/speccheck ./internal/suiteguardcontract` — every record-half package still passes once the boundary is gone. The guard reads the removal, so this preservation check cannot pass before the work.
- `make verify` — the complete concurrent gate passes, which is the regression this Task withdraws.

## References

- `_prd.md` → Core Features 6; Goals 3; Promoted work and known limitations; Decisions: Declared intentional breaks 3.
- `_techspec.md` → Vocabulary Contract; Coverage Map.
- `qa/qa-report-2026-09-10-02.md` → F-005.
