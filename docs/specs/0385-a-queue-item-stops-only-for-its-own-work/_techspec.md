---
spec: 0385-a-queue-item-stops-only-for-its-own-work
prd: _prd.md
created: 2026-11-12
---

# A queue item stops only for its own work — Technical Spec

## Executive Summary

`commandDeliveryEngineDependencies` takes the owner's log writer and gives it
to both the Delivery Engine and the `commandDeliveryWorkflow` it builds;
`newCommandDeliveryEngine` passes `os.Stderr`, which is the detached owner's
console log. Every workflow write goes through one accessor that falls back to
standard error, so no step writes to a nil writer. `runDeliveryOwner` replaces
its deferred unconditional owner release with the idle release the normal exit
already uses, so an owner that fails or panics while an item can advance keeps
its record and ADR-0279's status line names it. `parseReviewFindingAnchor`
accepts `.` or `;` after the line when a space or the end follows. The review
validator copies the reviewed checkout's regular `.ai-memory.toml` into its
private working directory, and the agent's sealed-directory rule admits that
one file. The trade-off is a named, optional coupling to ai-memory's marker in
exchange for sessions recorded under the right project without exposing the
checkout (ADR-0395).

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted or changed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0395 (new), ADR-0279, ADR-0373,
  ADR-0384, ADR-0196 and ADR-0069 apply; ADR-0056, ADR-0385, ADR-0286,
  ADR-0317, ADR-0180 and ADR-0344 do not (see ADR Alternatives). Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BD answers of 2026-11-12 quoted in the authorization record; the
  Governed Paths are the Roundfix skill files task_04 changes. No Makefile,
  workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0385-a-queue-item-stops-only-for-its-own-work/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Measured at `d79b30b4`:

1. `commandDeliveryEngineDependencies(runStore, loaded)` builds
   `commandDeliveryWorkflow{store, loaded, git, gh}` with no `log` and sets the
   engine's `Log` to `os.Stderr`. Only `deliver lane` and tests set `log`.
2. `Gate` writes `workflow.log` three times: `printSharedCacheTrim` (which
   returns on a nil writer, so the production gate never prints a trim), the
   Machine Verification wait callback and the lease-release warning (both
   `fmt.Fprintf`, which panic on nil). The item-binary probe and the two Linux
   Lane writers fall back to `os.Stderr` themselves.
3. A gate waiting behind a held unit panics in `Gate.func1`, the frame in the
   0382 owner's console log and quaestor's trace.
4. `runDeliveryOwner` defers `ReleaseDeliveryQueueOwner` unconditionally. A
   panic runs it while unwinding, and so does `printDeliverFailure` after an
   engine error; the queue records no owner while the item stays at
   `gating`, and `deliveryOwnerLine` returns no line for `OwnerPID <= 0`.
5. `parseReviewFindingAnchor` refuses any rune after the anchor other than a
   space, `:`, `,` or `)`.
6. `validateReviewConventions` creates `os.MkdirTemp("",
   "roundfix-review-validator-")` (mode 0700, empty) and passes it as
   `SealedPromptRequest.WorkDir`; `validateSealedPromptRequest` refuses a
   directory with any entry. `internal/baselineacp` creates its own empty
   directories with `withPrivateDirectory`.
7. `deliver status`, `runs list` and `doctor` already name a waiting
   repository gate in the Machine Verification Capacity block, from its
   `waiting` entry.

## System Architecture

No new package, command or file format. The changes stay in
`internal/cli` (Delivery Queue workflow, owner, anchor parser, validator) and
`internal/agent` (the sealed-directory rule and the marker copy).

## Implementation Design

### Interfaces

```go
// internal/cli/deliver_workflow.go
func commandDeliveryEngineDependencies(runStore *store.Store, loaded roundconfig.Loaded, log io.Writer) delivery.EngineDependencies
func (workflow *commandDeliveryWorkflow) console() io.Writer // log, or os.Stderr when nil

// internal/agent/sealed.go
const SealedIdentityMarker = ".ai-memory.toml"
func CopySealedIdentityMarker(repositoryRoot, workDir string) error
```

### Data Models

None: no record, schema or config key changes. The owner record keeps its
columns; only the moment it is cleared changes.

### API Contracts

1. API Contract 1 — the owner's log. `newCommandDeliveryEngine` calls
   `commandDeliveryEngineDependencies(runStore, loaded, os.Stderr)`; the
   function sets the workflow's `log` and the engine's `Log` to the same
   writer. `console()` returns `workflow.log`, or `os.Stderr` when it is nil.
   Every read of `workflow.log` outside `console()` is replaced by
   `workflow.console()`: in `Gate` (the shared-cache trim, the wait callback,
   the release warning), in the item-binary probe and in both Linux Lane
   writers of `deliver_linux_lane.go`, whose own nil fallbacks go away. The
   wait line keeps its text:
   `roundfix: Delivery Queue item <spec>: waiting for the Machine Verification Capacity; <held> of <capacity> held by <holders>`.
2. API Contract 2 — the owner's record. `runDeliveryOwner`'s deferred call
   becomes `runStore.ReleaseIdleDeliveryQueueOwner(context.WithoutCancel(ctx),
   loaded.GitRoot, pid, identity)` with its result ignored, as today's call's
   is. It clears the record only when no item is outside `merged` and `parked`
   and the record is still this owner's; otherwise the record stays. The
   normal exit path is unchanged. `deliver status`, `resume`, `retry`,
   `start` and `stop` are unchanged.
3. API Contract 3 — the anchor. After the `reviewAnchorPattern` match, a next
   rune `.` or `;` is accepted when it is the last rune of the text or the rune
   after it is a space (`unicode.IsSpace`); any other next rune keeps today's
   rule (space, `:`, `,`, `)`).
4. API Contract 4 — the sealed directory. `validateSealedPromptRequest`
   accepts a directory whose entries are none, or exactly one regular file
   named `SealedIdentityMarker`; otherwise it refuses with
   `sealed Agent working directory must be empty apart from its .ai-memory.toml`.
   `CopySealedIdentityMarker` reads `<repositoryRoot>/.ai-memory.toml` with
   `os.Lstat`: absent, not a regular file, or larger than 65536 bytes copies
   nothing and returns nil; a read or write failure returns an error; a copy is
   written with mode 0600.
5. API Contract 5 — the validator. After `os.MkdirTemp`,
   `validateReviewConventions` calls `agent.CopySealedIdentityMarker(repo.Root,
   workDir)`; an error makes the validator `unavailable` with that reason,
   which leaves every asked finding standing (ADR-0196). The directory is still
   removed when the validator returns.
6. API Contract 6 — the `deliver` guide. In
   `docs/user-guide/commands/deliver.md`, right after the sentence that ends
   "for `deliver resume` to recover.", this paragraph, word for word:
   "A Queue Owner that stops while an item can still advance, because its
   engine returned an error or the process panicked, keeps its owner record,
   so status reports it with the line above. It releases the record when it
   exits with no item left to advance, and `deliver stop` releases it after
   the owner it stopped has exited."
   and, right after the sentence "Each holder in `<holders>` is `<holder
   description>, <step>`, joined by `; `.", the sentence: "The queue log is
   the Queue Owner's console log, the file `deliver start` and `deliver
   resume` print as `Console Log:`."
7. API Contract 7 — the `review` guide. In `docs/user-guide/commands/review.md`,
   right after "For a comma-listed anchor, only the first range counts.":
   "An anchor followed by a period or a semicolon reads as one followed by a
   space when a space or the end of the finding comes next, so an anchor may
   end a sentence; `path:2.5` and `path:2.md` are not anchors."
   and, right after "It uses the review's selected runtime in a sealed session
   and may use no tools.": "The session works in a new private directory
   outside every checkout. When the reviewed checkout's root holds a regular
   `.ai-memory.toml` of at most 64 KiB, Roundfix copies it into that
   directory, which holds nothing else, so a tool that names its project from
   the working directory records the session under the reviewed repository."
8. API Contract 8 — the Roundfix skill. Both copies of
   `references/deliver.md` and `references/review.md` carry the four texts of
   API Contracts 6 and 7, word for word, at the same anchors.
9. API Contract 9 — the glossary. The **Queue Owner** entry of `CONTEXT.md`
   gains, at the end of its definition paragraph: "A Queue Owner that stops
   while an item can still advance keeps that record, so the report names it
   (ADR-0395).", where `ADR-0395` is a Markdown link to
   `docs/adr/0395-a-queue-item-stops-only-for-its-own-work.md`, as the
   entry's other ADR links are written.

### Surface Transcripts

None. The Spec adds no command, flag or output line: the gate's wait line and
the `Owner:` status line keep their texts, and the states that print them, a
held capacity unit and a dead owner process, are built in-process by the
named tests of Testing Approach.

### Invariants

1. A workflow message is never written to a nil writer; the production owner
   writes every one to its console log.
2. The wait line's text, the capacity entries and the order of the gate's
   steps are unchanged; waiting time still does not count against any limit.
3. An owner record is cleared only by an idle exit, `deliver stop`, or a
   proof that the owner is gone; status never reclaims it.
4. Punctuation never shortens a line number or a path: the path stays the
   regular expression's capture, and `.` or `;` followed by anything but a
   space keeps the anchor unreadable.
5. A sealed session's working directory is private, outside every checkout,
   and holds nothing but an optional regular copy of the repository's
   `.ai-memory.toml`; the session's tools stay denied. Baseline analysis
   copies nothing.

### Declared breaks

- `commandDeliveryEngineDependencies` gains a parameter. Its five test
  callers pass `io.Discard`, a change of their call only:
  `internal/cli/deliver_archived_resume_test.go`,
  `internal/cli/deliver_check_timeout_test.go` (two calls),
  `internal/cli/deliver_turn_test.go`,
  `internal/cli/deliver_review_capacity_test.go` and
  `internal/cli/deliver_free_space_test.go`; each gains the `io` import where
  it lacks it. A Task that finds another caller on the default branch when it
  runs changes it the same way and says so.
- The production gate now prints its shared-cache trim lines, which it
  silently dropped.
- An owner that fails or panics with an item left leaves its record, so
  `deliver status` prints the `Owner:` line and `deliver resume` prints
  `roundfix: Delivery Queue owner PID <pid> is not running; reclaimed its owner
  record.` before starting a new owner.
- `TestReviewFindingAnchorParsesEachForm` gains rows; no row changes.
- `conventionValidatorRunner` in
  `internal/cli/review_convention_validator_test.go` accepts the marker and
  records its bytes; every existing test keeps its fixture without a marker.
- The sealed refusal's message changes from `sealed Agent working directory
  must be empty`; no test or document pins the old text.

## Coverage Map

- Goal 1, User Story 1, Core Feature 1, Success Metric 1 → API Contract 1,
  Invariants 1 and 2.
- Goal 2, User Story 2, Core Feature 2, Success Metric 2 → API Contract 2,
  Invariant 3.
- Goal 3, User Story 3, Core Feature 3, Success Metric 3 → API Contract 3,
  Invariant 4.
- Goal 4, User Story 4, Core Feature 4, Success Metric 4 → API Contracts 4 and
  5, Invariant 5.
- Core Feature 5, Success Metric 5 → API Contracts 6 to 9.

Deviations from the PRD: None.

## Absence Evidence

1. The owner log accessor and the dependencies parameter — command: `grep -n 'func (workflow \*commandDeliveryWorkflow) console\|log:    log' internal/cli/deliver_workflow.go`; result: no match at `d79b30b4`, and the dependencies literal sets no `log`.
2. The sealed identity marker — command: `grep -rn 'SealedIdentityMarker\|ai-memory.toml' internal`; result: no match at `d79b30b4`.

## ADR Alternatives

- ADR-0395: the design is the decision; it takes none of the six options the
  record rejects.
- ADR-0279: it rejected reporting a queue with no recorded owner; status keeps
  that rule, and the owner keeps its record instead.
- ADR-0373: it rejected treating an identity mismatch as a live owner; a kept
  record whose PID was reused is still proven gone by it.
- ADR-0384: it rejected freeing a holder on elapsed time; the gate's wait and
  its holder rules are unchanged, only the writer of its line is set.
- ADR-0196: it rejected dropping a dismissed finding and letting a reviewer
  that ignores the grammar pass; an anchor that ends a sentence is still the
  grammar's anchor, and the validator still fails closed.
- ADR-0069: it rejected giving analysis checkout access; Baseline analysis
  still gets an empty directory, and the validator's directory gains only the
  marker copy, never a checkout file.
- ADR-0056: not applicable — it separates a Run's Task and Verification
  Capacity, which this Spec does not touch.
- ADR-0385: not applicable — it exports a core budget to Verification, which
  this Spec does not change.
- ADR-0286: not applicable — it governs the Correction Review of a fixed
  round-2 finding, whose rounds this Spec does not change.
- ADR-0317: not applicable — the Result Status Note dismissal runs before the
  validator and is unchanged.
- ADR-0180: not applicable — it derives built-in selections from the
  Recommended Profile, which this Spec does not touch.
- ADR-0344: not applicable — it governs profile proofs whose close fails,
  which this Spec does not change.

## Integration Points

- The detached owner's standard error is its console log; `deliver retry`
  builds the same dependencies with its terminal's standard error.
- ai-memory reads a `.ai-memory.toml` in a non-Git working directory outside
  the home directory only in that directory; Codex and Claude sessions do not
  load that file as instructions.

## Testing Approach

New tests, each calling `t.Parallel()` first:

- `internal/cli/deliver_owner_log_test.go`:
  `TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate` (on
  `newDeliveryItemFixture` with the lane off, `runs.max_verifying` 1 and
  Verification `exit 0`: builds the dependencies with a buffer whose first
  write releases a unit the test acquired with
  `store.AcquireMachineVerification` as Run
  `run_20261112T101500Z_0123456789abcdef` of `0096-other`, step `Task task_qa
  Verification attempt 1`; asserts first that the workflow's `log` and the
  engine's `Log` are that buffer, then that `Gate` passes, the buffer holds
  exactly the wait line naming that Run, and no capacity entry remains) and
  `TestADeliveryWorkflowWithoutALogWritesToStandardError` (a workflow without
  a log returns `os.Stderr` from `console()`, one with a log returns it) and
  `TestAQueueOwnerThatStopsWithAnItemLeftKeepsItsOwnerRecord` (subtests
  `panic`, `error`, `idle` on `newDeliveryOwnerStatusQueue` with one `gating`
  item: `runDeliveryOwner` with a fake engine that panics on a nil writer,
  returns an error, or merges the item; the first two keep owner PID
  `os.Getpid()` and an identity, the third clears both and exits `0`).
- `internal/cli/review_validation_test.go`: nine rows added to
  `TestReviewFindingAnchorParsesEachForm` (`path.go:2. Failure: broken`,
  `path.go:2; Failure: broken`, `` `path.go:2-4`. broken ``, `path.go:2.`,
  `notes.v2.md:3.` followed by a tab, readable; `path.go:2.5 broken`,
  `path.go:2.md broken`, `path.go:2;x broken`, `path.go:2.. broken`,
  unreadable) and `TestReviewReadsAnAnchorFollowedBySentencePunctuation` (two
  findings `review.txt:2.` and `` `review.txt:2`; `` with `Failure:` clauses:
  exit `exitRunFailed`, outcome `findings`, both anchored at `review.txt` line
  2 and standing).
- `internal/agent/sealed_identity_test.go`:
  `TestASealedWorkingDirectoryHoldsNothingButTheIdentityMarker` (empty and
  marker accepted; another file, the marker with another file, a marker
  directory and a marker symbolic link refused with the message of API
  Contract 4) and `TestCopySealedIdentityMarkerCopiesOnlyARegularMarker`
  (absent, directory, symbolic link and 65537 bytes copy nothing; a regular
  marker is copied with its bytes and mode 0600 and then passes validation).
- `internal/cli/review_convention_validator_test.go`:
  `TestTheReviewValidatorSessionCarriesTheRepositoryMarker` (subtests
  `marker`, where the fixture commits a `.ai-memory.toml` at its root, and `no
  marker`: the fake runner records the working directory's marker bytes, equal
  to the committed file or nil).

Sabotages, each on a temporary copy of the audited head made with `git
archive` and committed in a fresh Git repository, each observed by the named
test on the named case (all run by the author on 2026-11-12 against the
prototype):

1. Drop `log: log` from the workflow literal →
   `TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate` (the wiring
   assertion).
2. Give the engine `os.Stderr` instead of `log` → the same test.
3. Make `console()` return `workflow.log` with no fallback →
   `TestADeliveryWorkflowWithoutALogWritesToStandardError`.
4. Restore the unconditional `ReleaseDeliveryQueueOwner` defer →
   `TestAQueueOwnerThatStopsWithAnItemLeftKeepsItsOwnerRecord/panic` and
   `/error`.
5. Remove the `.`/`;` case → `TestReviewFindingAnchorParsesEachForm/path.go:2._Failure:_broken`
   and `TestReviewReadsAnAnchorFollowedBySentencePunctuation`.
6. Accept `.` or `;` without checking the next rune →
   `TestReviewFindingAnchorParsesEachForm/path.go:2.5_broken`.
7. Remove the `CopySealedIdentityMarker` call →
   `TestTheReviewValidatorSessionCarriesTheRepositoryMarker/marker`.
8. Accept any entry in `validateSealedPromptRequest` →
   `TestASealedWorkingDirectoryHoldsNothingButTheIdentityMarker/other_file`.
9. Accept the marker name whatever its type →
   `TestASealedWorkingDirectoryHoldsNothingButTheIdentityMarker/marker_symbolic_link`.
10. Follow symbolic links in `CopySealedIdentityMarker` (`os.Stat`) →
    `TestCopySealedIdentityMarkerCopiesOnlyARegularMarker/symbolic_link`.

## Build Order

1. The owner's log, the owner's record, their tests, the `deliver` guide text
   and the **Queue Owner** glossary entry (task_01).
2. The anchor punctuation, its tests and its `review` guide sentence (task_02)
   (depends on: 1, for the shared Behavior Surface Record).
3. The sealed marker rule and copy, the validator call, their tests and the
   `review` guide sentence (task_03) (depends on: 2, for the shared `review`
   guide and Behavior Surface Record).
4. The Roundfix `deliver` and `review` references and the skill version
   (task_04) (depends on: 3).
5. The QA gate (task_05) (depends on: 4).

No other Spec is required. Specs 0383 and 0384 deliver first and edit other
sections of the shared files `_prd.md` → Prerequisites names; each Task works on
the tree it runs on, cites functions rather than lines, and task_04 raises the
skill version with the record command.

## Risks & Considerations

- A kept record of a dead owner keeps showing until someone runs `deliver
  resume`; that is the intended signal.
- A repository's `.ai-memory.toml` may name a server profile; the copy is the
  repository's own file, the same one every session in its checkout reads.
- An ai-memory release that changes where it reads markers would leave the
  session in a stray project again, never break the review.
- Baseline analysis and profile proofs still open sessions in bare temporary
  directories (Non-Goals).

## Glossary

- changes: **Queue Owner**

## Decisions

- One log passed by the constructor, read through `console()`. See ADR-0395.
- The owner keeps its record on a failed exit with an item left, instead of a
  status line for a queue with no recorded owner. See ADR-0395 and ADR-0279.
- Punctuation ends an anchor only before a space or the end. See ADR-0395 and
  ADR-0196.
- The validator's working directory carries the repository's own marker, the
  only place ai-memory reads one outside a checkout and the home directory.
  See ADR-0395 and ADR-0069.
