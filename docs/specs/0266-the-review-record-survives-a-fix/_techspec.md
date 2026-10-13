---
spec: 0266-the-review-record-survives-a-fix
prd: _prd.md
created: 2026-10-13
---

# The review record survives a fix — Technical Spec

## Executive Summary

Three changes close the two defects the PRD names and the retry rule the
maintainer relaxed. First,
`reviewCheckoutDir` names a checkout's review directory from
`config.CanonicalizeResolvedPath` of the checkout path, and one helper,
`sameReviewCheckout`, replaces every lexical comparison of a record's or a
disposition's repository, so the correction proof, which receives the item
worktree in the Delivery Queue's link spelling, meets the record that
`roundfix review` wrote under the resolved path. Second, `persistReviewRecord`
also writes a Reviewed Head Record, `heads/<head commit>.json`, for every
non-reused verdict, and the proof reads the record of the parked candidate
from there, so review round 2 no longer erases what the proof compares
against. Third, the proof's path rule, ADR-0233's third check, also admits a
changed path that the archived Spec's Task files declared under `creates:` or
`interface:` or recorded under `## Recorded paths`, read from Git at the
parked candidate or at the Archive Record's `source_revision`, through a new
`spec.DeclaredTaskPaths`. No command output format, flag, field or exit code
changes; the path refusal's reason names both sets.

The trade-off accepted is a second file per reviewed head instead of carrying
round 1 inside the round-2 record. The lineage already carries
`previousFindings`, but a moved merge base, another provider or the ceiling
replaces the lineage, so the carried copy is lost exactly when an operator
reviews again in a new way. A second trade-off is that the declared paths
are read from the recorded revision, never from the item's head, so a
correction cannot widen its own bound, at the cost of a Git read per archived
Spec at each retry (ADR-0272).

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or record
  format changes; a Reviewed Head Record is a Pre-PR Review Record at another
  path, and recorded repositories are read, never rewritten. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call; tests use temporary homes and repositories, a fake reviewer runner and
  the local `git`. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0272 governs every decision below.
  ADR-0174: "in a directory named from the checkout path". ADR-0233: "every
  path changed from the candidate to the head lies under an archived Spec the
  blocker names", which ADR-0272 refines. ADR-0261: "A reader that opens the
  folder at `source_revision` and finds the commit missing performs" Source
  Revision Recovery, which the proof reuses. ADR-0165: "the correction is a new corrective Spec with its
  own `_authorization.md` and QA gate". ADR-0197: "It keeps the round-2 record
  in place". ADR-0268: "Path comparison resolves the longest existing prefix
  of a path". ADR-0259: "Every top-level test in a Parallel Test Package calls
  `t.Parallel()` as its first statement". ADR-0271: "A test bounds by real
  time only what it asserts". ADR-0256: "A release waits for the skills that
  describe a changed surface". ADR-0189: an owned skill's version names its
  content, raised with the record command. ADR-0244: "A Spec now carries a
  Glossary Declaration". ADR-0257: "Inside a Run the Daemon is the only full
  gate". ADR-0258: "A Spec is authored against the QA rerun classes".
  ADR-0265: "A Spec proves what it claims". ADR-0169: the record's
  `baseCommit` stays the merge base. ADR-0196: validation-dismissed findings
  stay outside the proof's count. ADR-0250 does not apply: no Baseline module
  version or coverage record changes. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. No lint, formatter, test-runner, Makefile or
  CI configuration changes. The Governed Paths task_01 changes, the owned
  `roundfix` skill's review and deliver references and its `SKILL.md` pair,
  were measured with `speccheck.GovernedPath` through a `go test -overlay`
  probe that wrote nothing, and are granted by the cycle AG approval of
  2026-10-13 ("Pode iniciar", with "Concedo" for each Spec's declared Governed
  Paths and "considere autorizado a ajustar todas as skills se necessário").
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0266-the-review-record-survives-a-fix/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Measured on 2026-10-13 on `main` at v0.76.0, on the maintainer's arm64 Mac
with Git 2.56.0, by a characterization test in `internal/cli` removed after
use. The test made a linked worktree at `<tmp>/link/repo-slug/item`, where
`<tmp>/link` points to `<tmp>/real`, ran `roundfix review` there with a fake
reviewer that raised one finding, committed a docs-only fix under
`docs/history/specs/0225-example/`, ran
`roundfix review dispose F1 --fixed-by <fix>`, and called
`ProveReviewCorrection`, once with an explicit and once with the default
Artifact Directory.

**Git registers the resolved path.** `git worktree list --porcelain` names the
worktree `<tmp>/real/repo-slug/item`, and the review record's `repository` is
that path.

**The proof misses the record under the link spelling.** With the link
spelling, the proof returned the error
`read review correction record: open <artifacts>/pre-pr-review/<key>/pre-pr-review.json: no such file or directory`,
in both Artifact Directory modes. With the resolved spelling it accepted. The
Artifact Directory itself was the same in both spellings: the default one is
named from the main checkout that the worktree's Git common directory names.

**Round 2 erases the candidate's record.** After a second `roundfix review`
at the fix, `outcome: reviewed`, `lineage.round: 2`, the proof with the
resolved spelling refused with
`review record must name the item worktree and parked candidate with outcome findings`.

**The sweep.** Every reader, writer and repository comparison of the record:

| Site | What it does | With two spellings |
| --- | --- | --- |
| `reviewCheckoutDir` | digest of `filepath.Clean(checkoutRoot)` | two directories; the root cause |
| `runReviewCommand` prior read, non-agent prior read, `reusableReviewRecord`, `runReviewDisposeCommand`, `persistReviewRecord`, `persistReviewAnswer`, `removeReviewAnswer` | key by `gitState.Root` or `record.Repository`, both from Git | resolved already; unchanged result |
| `ProveReviewCorrection` | keys by and compares `record.Repository != workDir`; `reviewDispositionsAtHead(ledger, workDir, …)` | link spelling: no record, no disposition |
| `decideReviewLineage`, the two prior comparisons in `runReviewCommand`, `reusableReviewRecord`, `runReviewDisposeCommand` | `filepath.Clean` equality | lexical |
| `reviewDispositionsAtHead`, `matchingReviewFindingDispositions`, `reviewFindingAlreadyDisposed`, `reviewFindingEvidenceDismissed`, `validateDismissedReviewRecord`, `validateCeilingClosedReviewRecord` | `==` or `!=` on `Repository` | lexical |
| `endOpenReviewSession` | uses `prior.Repository` as a work directory | no comparison |
| Delivery engine retry (`internal/delivery`) | passes the recorded item worktree to the proof | the link spelling |
| `doctor`, `gc`, `history`, archive | none reads `pre-pr-review/` | not affected |

The record has one writer, `persistReviewRecord`, reached through
`finishReviewCommand`; the disposition ledger has one writer, `review dispose`.

**A disposed code fix is refused whatever the archived Spec declared.** A
second characterization test, also removed after use, archived
`docs/history/specs/0225-example/` with a `task_01.md` declaring
`interface: review.txt`, ran `roundfix review` (one finding on `review.txt`),
committed a fix to `review.txt`, disposed it with `--fixed-by`, and called the
proof with the resolved spelling: it refused with
`changed path review.txt is outside the named archived Specs`. The review
command itself printed its corrective-Spec line, which stays (ADR-0233's
Consequences: only the retry knows what the correction changed).

## System Architecture

- **Resolved checkout key** (`internal/cli/review.go`): `reviewCheckoutDir`
  hashes `roundconfig.CanonicalizeResolvedPath(checkoutRoot)`.
- **Resolved comparisons** (`review.go`, `review_lineage.go`,
  `review_validation.go`, `deliver_review_correction.go`): a new unexported
  `sameReviewCheckout` replaces every repository comparison in the sweep.
- **Reviewed Head Record** (`review.go`, `deliver_review_correction.go`):
  `persistReviewRecord` writes it; the proof reads it through
  `readReviewRecordAtHead`.
- **Declared-path correction** (`internal/spec/declared_paths.go`,
  `deliver_review_correction.go`): `spec.DeclaredTaskPaths` returns a Task
  file's `creates:` and `interface:` paths and its recorded paths; the proof
  collects them for each named archived Spec at the recorded revision and
  admits a changed path found there.
- **Docs**: the review and deliver user guides, the Roundfix skill's review and
  deliver references and `CONTEXT.md`.

No change in `internal/delivery`, `internal/config` or the Delivery Queue
store; the prover interface keeps its signature.

## Implementation Design

### Interfaces

```go
// internal/cli/review.go
func reviewCheckoutDir(artifactDir, checkoutRoot string) string // digest of the resolved path
func sameReviewCheckout(left, right string) bool                // resolved paths equal
// reviewHeadRecordPath returns <checkout dir>/heads/<head>.json, or false when
// head is not 40 or 64 lowercase hexadecimal characters.
func reviewHeadRecordPath(artifactDir, checkoutRoot, head string) (string, bool)
// readReviewRecordAtHead returns the Reviewed Head Record of head, else the
// current record when it names head, else an error matching os.ErrNotExist.
func readReviewRecordAtHead(artifactDir, checkoutRoot, head string) (reviewRecord, error)

// internal/spec/declared_paths.go
// DeclaredTaskPaths returns the creates: and interface: paths of one Task
// file and the paths under its ## Recorded paths, clean and deduplicated.
func DeclaredTaskPaths(content []byte, taskFile string) ([]string, error)

// internal/cli/deliver_review_correction.go
// archivedSpecDeclaredPaths reads every task_NN.md of slug at the recorded
// revision: the candidate for a folder, source_revision for a record.
func (workflow *commandDeliveryWorkflow) archivedSpecDeclaredPaths(ctx context.Context, workDir, candidate, archiveRoot, slug string) ([]string, error)
```

### Data Models

A Reviewed Head Record is a file
`<Artifact Directory>/pre-pr-review/<checkout key>/heads/<head commit>.json`
with the Pre-PR Review Record schema, written atomically (temporary file and
rename, mode `0644`, directory `0755`) like the current record. No field is
added or removed.

### API Contracts

1. API Contract: resolved checkout key. `reviewCheckoutDir` keeps its digest
   (first 16 hexadecimal characters of SHA-256) and its layout, over
   `CanonicalizeResolvedPath(checkoutRoot)` instead of
   `filepath.Clean(checkoutRoot)`. For a path without symbolic links the
   directory is unchanged, so an existing record keeps its directory.
2. API Contract: resolved comparisons. `sameReviewCheckout(left, right)` is
   `CanonicalizeResolvedPath(left) == CanonicalizeResolvedPath(right)`. Every
   site of the sweep's comparison rows uses it; no `==`, `!=` or
   `filepath.Clean` comparison of a `Repository` field remains in
   `review.go`, `review_lineage.go`, `review_validation.go` or
   `deliver_review_correction.go`. A record of another checkout is still
   refused or ignored as today.
3. API Contract: Reviewed Head Record. `persistReviewRecord` writes the
   current record as today and then, when the record is not `reused` and its
   outcome is `reviewed`, `findings`, `findings-dismissed` or
   `ceiling-closed`, writes the same bytes to its head record path. A
   `blocked` or `omitted` record, a reused record and a head that is not a
   commit name write no head record and leave an existing one byte-identical.
   A later round at another head never touches an earlier head's record. The
   next round, reuse and `review dispose` keep reading `pre-pr-review.json`.
4. API Contract: correction proof. `ProveReviewCorrection` reads the record
   with `readReviewRecordAtHead(artifactDir, workDir, candidate)`. A missing
   record is a refusal with the existing reason
   `review record must name the item worktree and parked candidate with outcome findings`,
   no longer an error naming a file. The record must name the item worktree
   (`sameReviewCheckout`), the candidate and `findings`; dispositions come
   from the ledger through `reviewDispositionsAtHead`, which compares resolved
   repositories. ADR-0233's descent and disposition checks and their reasons
   are unchanged; its changed-path check is widened by API Contract 6.
5. API Contract: documentation. The review user guide and the Roundfix
   skill's review reference state that the checkout key is named from the
   resolved checkout path and that each non-reused verdict also keeps a
   Reviewed Head Record under `heads/<head commit>.json`, which later rounds
   never replace. The deliver user guide and the skill's deliver reference
   state that the moved-head `corrective-spec-required` retry reads the
   Reviewed Head Record of the parked candidate in the item worktree's
   resolved checkout, so a symbolic-linked `worktree.location` and a later
   review round leave it readable; that a changed path is admitted under an
   archived Spec the blocker names or among the paths its Tasks declared
   (`creates:`, `interface:` and recorded paths, read at the recorded
   revision); that review round 2 and the repository checks still follow; and
   that any other path still requires a corrective Spec. The table row of the
   moved-head retry says the same.
6. API Contract: declared-path correction. For each archived Spec the blocker
   names, the proof reads its form at the candidate with the existing archive
   reader. A Legacy Archive Folder's `task_*.md` files are read at the
   candidate; an Archive Record's are read under its `source` at its
   `source_revision`, recovered by `spec.EnsureSourceRevision` when missing
   locally. Each file's `spec.DeclaredTaskPaths` joins the declared set. A
   changed path is admitted when it lies under the Spec's archive path, as
   today, or equals a declared path exactly. The first changed path admitted
   by neither, in Git's `--name-only` order, is refused with the reason
   `changed path <path> is outside the named archived Specs and the paths their Tasks declared`,
   which keeps today's reason as its prefix. Task files that cannot be read
   or parsed refuse with
   `declared paths of archived Spec <slug> are unavailable: <cause>`. The
   accepted item still returns to `reviewing` for round 2, the archive and
   repository gate and the required checks, unchanged.

### Surface Transcripts

None. No command's output lines, flags or exit codes change; the retry
already prints the proof's reason inside its existing refusal, and the tests
assert the reasons.

### Invariants

1. For a checkout path without symbolic links, `reviewCheckoutDir` returns
   the directory it returns today.
2. Two spellings of one checkout name one review directory; two checkouts
   never share one.
3. A Reviewed Head Record, once written for a head, changes only when a later
   non-reused verdict for the same head replaces it.
4. The correction proof accepts no changed path outside the named archived
   Specs and their declared paths, and keeps every other ADR-0233 check.
5. The declared set is read at the recorded revision, never from the item's
   head, so a correction cannot add to it.

## Coverage Map

Deviations from the PRD: None.

- Goal 1, User Story 1, Core Feature 1, Core Feature 2, Success Metric 1 →
  resolved checkout key and resolved comparisons (task_02).
- Goal 3, Success Metric 3 → resolved comparisons and the sweep check
  (task_02).
- Goal 2, User Story 2, Core Feature 3, Success Metric 2 → Reviewed Head
  Record and the proof's reader (task_03).
- Goal 4, User Story 3, Core Feature 4, Success Metric 4 → declared-path correction
  (task_04), with the undeclared-file refusal under a linked root (task_02).
- Goal 5, User Story 4, Core Feature 5 → guides and skill references
  (task_01).
- Core Feature 6 → `CONTEXT.md` **Reviewed Head Record** and
  **Pre-PR Review Record** (task_01).

## Absence Evidence

1. Resolved checkout comparison helper — command: `grep -rn 'CanonicalizeResolvedPath\|EvalSymlinks' internal/cli/review.go internal/cli/review_lineage.go internal/cli/review_validation.go internal/cli/deliver_review_correction.go`; result: no match on `main` at v0.76.0; every repository comparison in those files is lexical, and the resolving helper exists only in `internal/config/config.go` (`CanonicalizeResolvedPath`), which this design reuses.
2. Reviewed Head Record — command: `grep -rn 'heads\|previousFindings\|PreviousFindings' internal/cli/review.go internal/cli/review_lineage.go internal/cli/deliver_review_correction.go`; result: four matches and no per-head record: `refs/heads/main` and `refs/heads/master` in the base lookup, and the lineage field, its assignment and its use in the round-2 prompt; the only kept history is `PreviousFindings` and `PreviousDispositions` in the round-2 lineage, which `decideReviewLineage` resets on a new lineage and `closeAtCeiling` replaces, and the proof never reads it.
3. Declared-path correction — command: `grep -rn 'RecordedTaskPaths\|ContextKindInterface\|ContextKindCreates' internal/cli/deliver_review_correction.go internal/cli/review.go; grep -rln 'func DeclaredTaskPaths' internal/spec`; result: no match; the proof admits only paths under the archive prefix, `spec.RecordedTaskPaths` exists for the scope audit, and the Context parser `parseTaskContextRefs` is unexported, so `DeclaredTaskPaths` wraps both.
4. Guide and skill wording — command: `grep -rn 'heads/\|resolved checkout' docs/user-guide/commands/review.md docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/review.md .agents/skills/roundfix/references/deliver.md CONTEXT.md`; result: no match; the guides say only that the record lives under `pre-pr-review/<checkout key>/`.

## ADR Alternatives

- ADR-0272: this Spec's own decision; it rejects lineage-carried history, round-numbered records, a Delivery Queue path migration and a `review dispose --head` flag.
- ADR-0174: no rejected alternative is taken; each checkout keeps its own record, now named from its resolved path, and the disposition ledger stays shared.
- ADR-0233: refined, not superseded, on the maintainer's answer "Relaxar, com limites"; its rejected alternatives (refreshing items from the default branch, declaring every `SKILL.md` derived) stay rejected, and its third retry rule now also admits the archived Spec's declared paths.
- ADR-0261: no rejected alternative is taken; the proof recovers a missing `source_revision` the way every reader does and pushes no ref.
- ADR-0165: no rejected alternative is taken; an archived Spec is still never edited in place, and a correction outside its records and declared paths still parks for a corrective Spec.
- ADR-0197: no rejected alternative is taken; the lineage still spans two rounds and the ceiling still keeps the round-2 record in place.
- ADR-0268: no rejected alternative is taken; the prefix-resolving comparison is reused, and no recorded path is rewritten, as it decided.
- ADR-0259: no rejected alternative is taken; every new test calls `t.Parallel()` first and no Sequential Test is added.
- ADR-0271: no rejected alternative is taken; the new tests wait on no real time, and no test file gains a Real-Time Budget.
- ADR-0256: its rejected model-judging alternative stays rejected; task_01 edits the covering skill files of the changed surfaces.
- ADR-0189: no rejected alternative is taken; the version raise is written and recorded by the record command, never by hand.
- ADR-0244: no rejected alternative is taken; both terms follow the declaration contract.
- ADR-0257: no rejected alternative is taken; Task agents run focused tests and the Daemon runs the declared Verification.
- ADR-0258: no rejected alternative is taken; Verification names the consuming suites, and the QA rows reach no network or provider.
- ADR-0169: no rejected alternative is taken; the review still diffs from the merge base, and a Reviewed Head Record keeps the `baseCommit` it was written with.
- ADR-0196: no rejected alternative is taken; validation still decides which findings park, and the proof still skips a finding dismissed by validation.
- ADR-0250: does not apply; no Baseline module version or coverage record changes, so neither of its rejected alternatives is in play.
- ADR-0265: no rejected alternative is taken; this TechSpec carries the Absence Evidence, ADR Alternatives and Coverage Cross-Read it requires.

## Integration Points

None outside the repository. Git is the local `git` binary; its behavior
under a symbolic link was measured above.

## Testing Approach

New tests, each calling `t.Parallel()` first, all on temporary homes and
repositories, a fake reviewer runner through `withAgentRunner`, and a
worktree root `<tmp>/link` that is a symbolic link to `<tmp>/real`, with the
item worktree added by `git worktree add <tmp>/link/repo-slug/item`:

- `internal/cli/review_resolved_checkout_test.go` (task_02):
  `TestReviewCheckoutDirIsOneForEverySpellingOfACheckout` (the link and the
  resolved spelling name one directory, which is the digest of the resolved
  path; another worktree names another);
  `TestTheCorrectionProofFindsTheRecordUnderASymlinkedWorktreeRoot` with
  subtests `explicit-artifact-directory` and `default-artifact-directory`
  (review with a finding, a docs-only fix under the archived Spec, dispose
  `--fixed-by`, proof with the link spelling accepted);
  `TestTheCorrectionProofRefusesACodeFixUnderASymlinkedWorktreeRoot` (a fix
  to `review.txt` is refused with `changed path review.txt is outside the named archived Specs`);
  `TestTheCorrectionProofRefusesARecordOfAnotherCheckout` (a record naming a
  sibling worktree is refused with the record reason).
- `internal/cli/review_head_record_test.go` (task_03):
  `TestASecondReviewRoundKeepsTheFirstRoundsRecord` (after round 2 the
  candidate's head record still holds the round-1 `findingItems`, and the
  fix's head record holds round 2);
  `TestTheCorrectionProofAcceptsAfterASecondRoundUnderASymlinkedWorktreeRoot`;
  `TestAReusedBlockedOrOmittedReviewWritesNoReviewedHeadRecord` with subtests
  `reused`, `blocked` and `omitted`;
  `TestTheCorrectionProofReadsTheCurrentRecordWithoutAHeadRecord` (an item
  parked before this change);
  `TestTheCorrectionProofRefusesWithoutARecordOfTheCandidate`;
  `TestReviewHeadRecordPathRefusesANonCommitName`.
- `internal/spec/declared_paths_test.go` (task_04):
  `TestDeclaredTaskPathsReadsCreatesInterfaceAndRecordedPaths`.
- `internal/cli/review_declared_correction_test.go` (task_04):
  `TestACodeCorrectionInsideTheDeclaredPathsReturnsToReview` with subtests
  `archive-record` and `legacy-folder`;
  `TestARecordedPathCountsAsDeclared`;
  `TestACodeCorrectionOutsideTheDeclaredPathsNamesTheFirstPath`;
  `TestDeclaredPathsAreReadAtTheRecordedRevision` (a Task file the
  correction itself edits to declare a path does not admit it);
  `TestUnreadableDeclaredPathsRefuseTheCorrection`.
- Existing tests stay unchanged and must pass, notably
  `deliver_review_correction_test.go`, `review_record_checkout_test.go`,
  `review_lineage_test.go` and `review_disposition_test.go`; their fixtures
  already use resolved paths or keep both sides in one spelling.
- Each Task runs its new tests once with `-race` (ADR-0259).

## Build Order

1. Guides, skill references with the version raise, `CONTEXT.md` and the
   Behavior Surface Record, stating API Contract 5 before the code that
   implements it (task_01).
2. Resolved checkout key and comparisons, with their tests (task_02)
   (depends on: 1).
3. Reviewed Head Record and the proof's reader, with their tests (task_03)
   (depends on: 2).
4. Declared-path correction in the proof and `spec.DeclaredTaskPaths`, with
   their tests (task_04) (depends on: 3).
5. QA gate (depends on: 4).

task_01 comes first so that each code Task changes a CLI surface whose guide
a Task it depends on already names. task_02, task_03 and task_04 all edit
`internal/cli/deliver_review_correction.go`, and task_02 and task_03 also
`internal/cli/review.go`, so they run in series. All four ship in one Pull
Request, so no release carries the guides without the code.

## Risks & Considerations

- **A test that compared spellings by accident.** A test that wrote a record
  under an unresolved macOS temporary path and read it under the same path
  still meets it, because both sides now resolve; one that relied on two
  spellings being two checkouts would break, and the sweep found none.
- **Head records accumulate.** One small file per reviewed head per checkout;
  item worktrees are short-lived. Reclaiming them is out of scope.
- **A declared path is broad by design.** A Spec that declared a shared
  file, such as `internal/cli/cli_test.go`, admits a fix anywhere in it. The
  maintainer bounded the rule by declaration, not by hunk; round 2 reviews
  the delta and the repository gate runs, so the fix is still reviewed and
  tested.
- **The current record still moves.** `review dispose` reads the current
  record, so a finding of round 1 can no longer be disposed after round 2 has
  run; the operator disposes first, as the observed flow did. The proof then
  still finds the round-1 findings.

## Glossary

- adds: **Reviewed Head Record**
- changes: **Pre-PR Review Record**
- not a term: **Resolved checkout key** — the component that applies the rule
- not a term: **Resolved comparisons** — the behavior of one helper
- not a term: **Git registers the resolved path** — a heading of the measurements
- not a term: **The proof misses the record under the link spelling** — a heading of the measurements
- not a term: **Round 2 erases the candidate's record** — a heading of the measurements
- not a term: **The sweep** — a heading of the measurements
- not a term: **A test that compared spellings by accident** — a risk heading
- not a term: **Head records accumulate** — a risk heading
- not a term: **The current record still moves** — a risk heading
- not a term: **A declared path is broad by design** — a risk heading
- not a term: **A disposed code fix is refused whatever the archived Spec declared** — a heading of the measurements
- not a term: **Declared-path correction** — the component that applies the widened rule
- not a term: **Current Behavior** — a section heading

## Decisions

- The record directory and every repository comparison use the resolved
  checkout path; see ADR-0272.
- Each non-reused verdict keeps a Reviewed Head Record, and the proof reads it
  for the parked candidate; see ADR-0272.
- A missing record is a refusal with the existing reason, not an error.
- The retry also admits the archived Spec's declared paths, read at the
  recorded revision; see ADR-0272, which refines ADR-0233.
