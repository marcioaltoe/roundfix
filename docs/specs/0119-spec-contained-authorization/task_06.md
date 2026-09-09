---
task: task_06
spec: 0119-spec-contained-authorization
status: pending
type: backend
complexity: high
---

# Task 06: Execute authored commands only on committed provenance

## Overview

Executing a Spec's authored Verification runs shell the Spec's author wrote.
Bind that execution to committed provenance: commands run when the artifact
carrying them is tracked here and byte-identical to its committed bytes, and
otherwise require an execution approval naming the approved revision. The slice
is verifiable on its own — the same Spec executes before an edit and refuses
after one, while read-only checking keeps working throughout.

## Requirements

1. MUST record today's execution behavior at the authored-command boundary
   before changing it, so the change is measured against the present contract.
2. MUST execute an authored command only when its carrying artifact is tracked
   in this repository and byte-identical to its committed bytes at the resolved
   revision.
3. MUST refuse with `SC-SOURCE-UNTRUSTED` when the Spec Root resolves outside
   the repository's Git tree, when the carrying artifact is untracked or
   modified against its committed bytes, or when the command text differs from
   the committed bytes, and MUST name which of those conditions fired.
4. MUST accept an execution approval only when the approval record itself
   reaches the trusted bar it is granting: it comes from a committed ancestor
   of the delivery target, not from the same untracked or modified source it
   authorizes. An approval a source can append to itself, pointing at a
   revision that source just created, is self-approval and grants nothing.
5. MUST bind an accepted approval to the repository, the approved revision, the
   carrying artifact, and a digest of the exact approved command text, and MUST
   fall back to refusal once any of those changes, without withdrawing the
   historical approval record.
6. MUST apply the same decision at every path that executes an authored
   command, emitting one shared code rather than private diagnostics. The
   pre-dispatch probe is not that boundary on its own: post-Agent verification
   and Settle each call the verifier directly rather than through the prober, so
   guarding only the probe leaves the two paths that actually run Agent-supplied
   commands open. Cover the probe, the post-Agent verification call, and the
   Settle call, and test each through its own path rather than through a shared
   helper.
7. MUST keep `roundfix spec check` without command execution usable against any
   source, trusted or not.
8. MUST report an unreadable Git object or an unavailable revision as an
   unresolved source and refuse to execute, never as trusted and never as a
   manufactured verdict.
9. MUST document both `SC-TOOLING-UNAPPROVED` and `SC-SOURCE-UNTRUSTED` as
   glossary entries, so the coined tokens have a durable owner.
10. MUST NOT grant network access, credential access, a different sandbox, or any
    execution privilege beyond running the already-approved commands.

## Subtasks

- [ ] Record the present execution behavior at the authored-command boundary.
- [ ] Compare the carrying artifact and command text against committed bytes.
- [ ] Refuse with the shared code, naming which condition fired.
- [ ] Accept a revision-bound execution approval and expire it on edit.
- [ ] Apply the decision at probing, Implement dispatch and Settle.
- [ ] Give both coined codes a glossary entry.

## Acceptance Criteria

- [ ] A Spec tracked here and unmodified executes its authored commands, and the
      same Spec with one command's text edited in the working tree refuses with
      `SC-SOURCE-UNTRUSTED` naming the modified-artifact condition.
- [ ] A Spec Root resolving outside the repository's Git tree refuses, naming
      the out-of-tree condition.
- [ ] A Spec carrying an execution approval already committed in the delivery
      target executes; an approval appended to the same untracked or modified
      source it authorizes refuses as self-approval.
- [ ] The approval stops authorizing once the repository, revision, artifact or
      approved command digest changes, and the historical approval record is
      unchanged.
- [ ] Read-only checking without command execution returns its normal report for
      an untrusted source and executes nothing.
- [ ] The refusal is reported identically from the probe, from the post-Agent
      verification call, and from Settle, each exercised through its own path.
- [ ] An unreadable Git object reports an unresolved source and executes
      nothing.
- [ ] Both coined codes appear in the glossary with the meaning the readers
      emit, so the Vocabulary Contract check runs instead of skipping.

## Context

- interface: `internal/speccheck/verification.go`
- interface: `internal/cli/spec_check.go`
- interface: `internal/daemon/verification_probe.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/cli/settle.go`

## Verification

- `grep -q 'SC-SOURCE-UNTRUSTED' internal/speccheck/verification.go && grep -q 'func TestVerificationRefusesUntrustedSource' internal/speccheck/verification_test.go && go test -count=1 ./internal/speccheck -run '^TestVerificationRefusesUntrustedSource$'` — an out-of-tree root, an untracked artifact, a modified artifact and altered command text each refuse with the shared code naming the condition.
- `grep -q 'func TestVerificationExecutesOnCommittedProvenance' internal/speccheck/verification_test.go && go test -count=1 ./internal/speccheck -run '^TestVerificationExecutesOnCommittedProvenance$'` — a tracked, unmodified Spec executes, and a revision-bound execution approval expires when the approved commands change.
- `grep -q 'func TestSettleRefusesUntrustedCommandSource' internal/cli/settle_test.go && go test -count=1 ./internal/cli -run '^TestSettleRefusesUntrustedCommandSource$'` — Settle refuses through its own direct verifier call, not through the prober.
- `grep -q 'func TestPostAgentVerificationRefusesUntrustedCommandSource' internal/daemon/engine_test.go && go test -count=1 ./internal/daemon -run '^TestPostAgentVerificationRefusesUntrustedCommandSource$'` — the post-Agent verification call refuses through its own path.
- `grep -q 'SC-SOURCE-UNTRUSTED' CONTEXT.md && grep -q 'SC-TOOLING-UNAPPROVED' CONTEXT.md` — both coined codes carry a glossary owner.
- `grep -q 'SC-SOURCE-UNTRUSTED' internal/speccheck/verification.go && go build -buildvcs=false ./...` — the shared code exists and the tree compiles with it wired in.

## References

- `_prd.md` → User Stories 2; Core Features 6; Goals 3; Decisions: Declared intentional breaks 3.
- `_techspec.md` → Vocabulary Contract; Implementation Design: Audit and compatibility; API Contracts; Build Order 4.
- ADR-0014, ADR-0096.
