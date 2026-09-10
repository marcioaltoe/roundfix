---
task: task_06
spec: 0119-spec-contained-authorization
status: completed
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
2. MUST compare an authored projection of the carrying artifact, not its whole
   file bytes. The projection covers the authored contract — the Verification
   command text and the requirements that bound it — and excludes the fields the
   Daemon owns, its `status` frontmatter and its appended `## Result`. The
   Daemon sets a Task to `in_progress` before the Agent turn and therefore
   before Verification, so a whole-file comparison differs from the committed
   bytes on every normal Task and would refuse the entire implementation loop.
   The TechSpec already draws this line: Daemon-owned status and Result updates
   are execution records, not permission.
3. MUST execute an authored command only when its carrying artifact is tracked
   in this repository and its authored projection matches the projection of its
   committed bytes at the resolved revision.
4. MUST record today's outcome at each authored-command entry point as a
   characterization before changing it, so an unintended change of execution
   behavior fails against the recorded present contract rather than against a
   test written afterwards to agree with the new code.
5. MUST refuse with `SC-SOURCE-UNTRUSTED` when the Spec Root resolves outside
   the repository's Git tree, when the carrying artifact is untracked, or when
   its authored projection differs from the committed projection, and MUST name
   which of those conditions fired.
6. MUST accept an execution approval only when the approval record itself
   reaches the trusted bar it is granting: it lives in this repository, in the
   consuming Spec's authorization record, and is already committed in the
   delivery target's ancestry. It is never read from the untrusted source it
   authorizes, which is what makes an out-of-tree or modified source
   approvable at all: the approval names that source, the source does not
   supply it. An approval a source appends to itself, pointing at a revision
   that source just created, is self-approval and grants nothing.
7. MUST bind an accepted approval to the repository, the approved revision, the
   carrying artifact, and a digest of the exact approved command text, and MUST
   fall back to refusal once any of those changes, without withdrawing the
   historical approval record.
8. MUST apply the same decision at every path that executes an authored
   command, emitting one shared code rather than private diagnostics. The
   pre-dispatch probe is not that boundary on its own: post-Agent verification
   and Settle each call the verifier directly rather than through the prober, so
   guarding only the probe leaves the two paths that actually run Agent-supplied
   commands open. Cover the probe, the post-Agent verification call, and the
   Settle call, and test each through its own path rather than through a shared
   helper.
9. MUST keep `roundfix spec check` without command execution usable against any
   source, trusted or not.
10. MUST report an unreadable Git object or an unavailable revision as an
   unresolved source and refuse to execute, never as trusted and never as a
   manufactured verdict.
11. MUST emit the `SC-SOURCE-UNTRUSTED` token exactly as the glossary's Grant
    Refusal Code entry defines it, so the Vocabulary Contract runs clean once
    the token exists. The glossary entry is already in place; this Task must
    match it rather than restate it.
12. MUST NOT grant network access, credential access, a different sandbox, or any
    execution privilege beyond running the already-approved commands.

## Subtasks

- [ ] Record the present execution behavior at each authored-command entry point.
- [ ] Compare the carrying artifact and command text against committed bytes.
- [ ] Refuse with the shared code, naming which condition fired.
- [ ] Accept a revision-bound execution approval and expire it on edit.
- [ ] Apply the decision at probing, Implement dispatch and Settle.
- [ ] Give both coined codes a glossary entry.

## Acceptance Criteria

- [ ] A Spec tracked here and unmodified executes its authored commands, and the
      same Spec with one command's text edited in the working tree refuses with
      `SC-SOURCE-UNTRUSTED` naming the modified-artifact condition.
- [ ] A Task the Daemon has set to `in_progress`, and one carrying an appended
      `## Result`, still execute: the projection ignores both, so the normal
      implementation loop is not refused by its own trust rule.
- [ ] A characterization case records today's outcome at the probe, the
      post-Agent verification call and Settle, and fails if any of those
      outcomes changes for a source the Spec did not intend to affect.
- [ ] A Spec Root resolving outside the repository's Git tree refuses, naming
      the out-of-tree condition.
- [ ] An out-of-tree Spec Root and a modified in-tree artifact each execute
      when this repository's committed authorization record approves that exact
      source, proving the exception path works rather than only that the happy
      path does.
- [ ] An approval appended to the same untracked or modified source it
      authorizes refuses as self-approval, even when it names a revision that
      exists.
- [ ] The approval stops authorizing once the repository, revision, artifact or
      approved command digest changes, and the historical approval record is
      unchanged.
- [ ] Read-only checking without command execution returns its normal report for
      an untrusted source and executes nothing.
- [ ] The refusal is reported identically from the probe, from the post-Agent
      verification call, and from Settle, each exercised through its own path.
- [ ] An unreadable Git object reports an unresolved source and executes
      nothing.
- [ ] The emitted `SC-SOURCE-UNTRUSTED` token matches the glossary's Grant
      Refusal Code entry, and the whole-Spec-Root check reports no vocabulary
      finding.

## Context

- interface: `internal/speccheck/verification.go`
- interface: `internal/cli/spec_check.go`
- interface: `internal/daemon/verification_probe.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/cli/settle.go`

## Verification

- `grep -q 'SC-SOURCE-UNTRUSTED' internal/speccheck/verification.go && grep -q 'func TestVerificationRefusesUntrustedSource' internal/speccheck/verification_test.go && go test -count=1 ./internal/speccheck -run '^TestVerificationRefusesUntrustedSource$'` — an out-of-tree root, an untracked artifact, a modified artifact and altered command text each refuse with the shared code naming the condition.
- `grep -q 'func TestVerificationExecutesOnCommittedProvenance' internal/speccheck/verification_test.go && go test -count=1 ./internal/speccheck -run '^TestVerificationExecutesOnCommittedProvenance$'` — a tracked, unmodified Spec executes, and a revision-bound execution approval expires when the approved commands change.
- `grep -q 'func TestAuthoredProjectionIgnoresDaemonOwnedFields' internal/speccheck/verification_test.go && go test -count=1 ./internal/speccheck -run '^TestAuthoredProjectionIgnoresDaemonOwnedFields$'` — a Task set to `in_progress` and one carrying an appended Result still execute, so the trust rule does not refuse the normal loop.
- `grep -q 'func TestExecutionEntryPointCharacterization' internal/speccheck/verification_test.go && go test -count=1 ./internal/speccheck -run '^TestExecutionEntryPointCharacterization$'` — today's outcome at each entry point is recorded and asserted as the present contract.
- `grep -q 'func TestProbeRefusesUntrustedCommandSource' internal/daemon/verification_probe_test.go && go test -count=1 ./internal/daemon -run '^TestProbeRefusesUntrustedCommandSource$'` — the pre-dispatch probe refuses through its own path.
- `grep -q 'func TestSettleRefusesUntrustedCommandSource' internal/cli/settle_test.go && go test -count=1 ./internal/cli -run '^TestSettleRefusesUntrustedCommandSource$'` — Settle refuses through its own direct verifier call, not through the prober.
- `grep -q 'func TestPostAgentVerificationRefusesUntrustedCommandSource' internal/daemon/engine_test.go && go test -count=1 ./internal/daemon -run '^TestPostAgentVerificationRefusesUntrustedCommandSource$'` — the post-Agent verification call refuses through its own path.
- `grep -q 'SC-SOURCE-UNTRUSTED' internal/speccheck/verification.go || exit 1; go run -buildvcs=false ./cmd/roundfix spec check` — the token exists and every active Spec still checks clean, so the Vocabulary Contract runs rather than skips and the emitted token matches its glossary entry.
- `grep -q 'SC-SOURCE-UNTRUSTED' internal/speccheck/verification.go && go build -buildvcs=false ./...` — the shared code exists and the tree compiles with it wired in.

## References

- `_prd.md` → User Stories 2; Core Features 6; Goals 3; Decisions: Declared intentional breaks 3.
- `_techspec.md` → Vocabulary Contract; Implementation Design: Audit and compatibility; API Contracts; Build Order 4.
- ADR-0014, ADR-0096.

## Result

Implemented one shared authored-command source decision before the probe,
post-Agent Verification, and Settle shell boundaries. A Task executes from
committed provenance when its Spec Root is inside the delivery repository, its
artifact is tracked at the resolved revision, its requested commands match the
carrying file, and its authored projection matches the committed projection.
The projection excludes only `status` frontmatter, an appended `## Result`, and
the line-break separator before that Result.

An untrusted source can execute only through `execution_approvals` in the
committed consuming `_authorization.md`. Each approval binds a Git repository
identity, the carrying artifact's last committed revision, its
repository-relative path, and `sha256:` of the exact command text. The reader
loads the record from delivery-target ancestry, so an uncommitted approval or
an approval carried by an external source cannot authorize itself. The check
performs Git and filesystem reads only; it changes no sandbox, environment,
network, credential, or command-runner privilege.

Acceptance evidence:

1. `TestVerificationExecutesOnCommittedProvenance/tracked_unmodified_artifact_is_authorized`
   accepts a committed Task. `TestVerificationRefusesUntrustedSource/altered_command_text`
   and `/requested_command_differs_from_carrying_artifact` return
   `SC-SOURCE-UNTRUSTED` with `modified-artifact` before execution.
2. `TestAuthoredProjectionIgnoresDaemonOwnedFields` changes `status` to
   `in_progress`, then appends `## Result`; both projections remain authorized.
3. Before implementation, the existing probe, Task-cycle, `--run-verification`,
   and Settle characterization selections passed unchanged. After the change,
   `TestExecutionEntryPointCharacterization` preserves the committed-source
   outcome for the probe, post-Agent Verification, and Settle identities.
4. `TestVerificationRefusesUntrustedSource/Spec_Root_outside_repository_Git_tree`
   returns the shared code with `out-of-tree-spec-root`.
5. `TestVerificationExecutesOnCommittedProvenance` accepts both a modified
   in-tree artifact and an out-of-tree artifact when the committed delivery
   record carries their exact approval. Existing external Implement fixtures
   now carry the same committed approval and pass through the public flow.
6. `TestVerificationRefusesUntrustedSource/modified_source_cannot_append_its_own_approval`
   and `/out-of-tree_source_cannot_carry_its_own_approval` prove that source-local
   self-approval grants nothing even when it names an existing revision.
7. `TestVerificationExecutesOnCommittedProvenance/approval_expires_with_every_bound_identity_field`
   varies repository, revision, artifact, and command digest independently.
   Every case refuses and re-reads the historical authorization record
   byte-identically.
8. `TestSpecCheckRunVerification/keeps_read-only_checking_available_for_an_edited_source`
   returns the normal report without `--run-verification`, creates no marker,
   then returns the shared refusal when execution is requested.
9. `TestProbeRefusesUntrustedCommandSource`,
   `TestPostAgentVerificationRefusesUntrustedCommandSource`, and
   `TestSettleRefusesUntrustedCommandSource` exercise their own entry points.
   Each reports `SC-SOURCE-UNTRUSTED` plus `modified-artifact`, and each asserts
   that the verifier or shell was not reached.
10. `TestVerificationRefusesUntrustedSource/unreadable_Git_object`,
    `/unavailable_revision`, and `/unavailable_approved_source_revision` report
    `unresolved-source` and execute nothing.
11. Fresh source inspection finds `SC-SOURCE-UNTRUSTED` in
    `internal/speccheck/verification.go` and both declared Grant Refusal Codes
    in `CONTEXT.md`'s `Grant Refusal Code` entry. The whole-Spec-Root command is
    reserved for Daemon Verification and was not run in this Agent turn.

Focused checks:

- Pre-change, `rtk env GOCACHE=/tmp/roundfix-task06-go-build go test ./internal/daemon -run '^(TestProbeCommands|TestTaskCycleExecutesAgentVerifySettleCommitContract)$' -count=1`
  passed, recording that the probe and post-Agent path executed their current
  committed source. The matching CLI selection for
  `TestSpecCheckRunVerification` and
  `TestSettleVerificationRunsSurfaceCommandsVerbatim` also passed. Initial
  attempts using the default Go cache were sandbox-blocked before compilation
  and are not behavior evidence.
- After the final source edit, the cross-entry-point selection
  `rtk env GOCACHE=/tmp/roundfix-task06-go-build go test ./internal/speccheck ./internal/daemon ./internal/cli -run '^(TestVerificationRefusesUntrustedSource|TestVerificationExecutesOnCommittedProvenance|TestAuthoredProjectionIgnoresDaemonOwnedFields|TestExecutionEntryPointCharacterization|TestProbeRefusesUntrustedCommandSource|TestPostAgentVerificationRefusesUntrustedCommandSource|TestSettleRefusesUntrustedCommandSource|TestSettleVerificationRunsSurfaceCommandsVerbatim|TestSpecCheckRunVerification|TestRunImplementUsesConfiguredExternalSpecRootEndToEnd)$' -count=1`
  passed in all three packages.
- A broader affected-package run reported `internal/speccheck` and
  `internal/daemon` passing. Its CLI process hit one intermittent concurrent
  Git-worktree fixture error while reading another Task Worktree's `commondir`;
  the isolated failing test passed unchanged, and a fresh full
  `rtk env GOCACHE=/tmp/roundfix-task06-go-build go test ./internal/cli -count=1`
  rerun passed in 80.258 seconds. No Task-06 code change was made for that
  non-reproducible lifecycle observation.
- `rtk env GOCACHE=/tmp/roundfix-task06-go-build go vet ./internal/speccheck ./internal/daemon ./internal/cli`
  passed.

The Task's declared `## Verification` commands were not run; Daemon
Verification remains pending.
