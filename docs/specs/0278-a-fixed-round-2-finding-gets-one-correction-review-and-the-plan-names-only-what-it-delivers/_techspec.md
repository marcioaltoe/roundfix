---
spec: 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers
prd: _prd.md
created: 2026-10-18
---

# A fixed round-2 finding gets one Correction Review, and the plan names only what it delivers — Technical Spec

## Executive Summary

Three slices extend existing seams, and no package is new.

1. The pre-PR review's lineage gains the **Correction Review** (ADR-0286).
   - When round 2 ended with findings that are all disposed as fixed,
     `roundfix review` reviews the delta from the round-2 head once instead of
     closing at the ceiling.
   - The Delivery Queue's correction proof stops applying ADR-0272's
     declared-path bound to that case.
2. The Spec Consistency Check gains `SC-VERIFY-PATH-UNDECLARED`, an error from
   the Tasks stage (ADR-0287).
3. A test appended to `internal/docscontract/publicdocs_test.go` builds the
   release check's Managed Refresh plan, so every Task settlement reports a
   stale Setup Manifest.

The trade-offs accepted:

- One more reviewer call over a ceiling that never reviews a fix.
- A path rule that misses paths it cannot recognize, chosen over one that
  raises false positives.

The QA-delivery question that the brief proposed was measured on Jev before
the Run and rejected (ADR-0287), so the judge is unchanged.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; the review record's lineage gains the optional
  boolean `correction` and the Spec Consistency Check gains one code. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the reviewer call of a Correction Review
  uses the configured provider exactly as round 2 does, and every test uses a
  fake reviewer. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — each Build Order step follows recorded
  decisions. ADR-0286 governs Build Order 1: "A fixed round-2 finding gets one
  Correction Review." ADR-0287 governs Build Orders 2 and 3: "The Spec
  Consistency Check reports a Verification path the plan does not provide."
  ADR-0197, amended: "A pre-PR reviewer lineage spans at most two rounds".
  ADR-0272, amended: "A review record survives a fix, and a correction inside
  the archived Spec's declared paths returns to review". ADR-0103, amended: "An
  applied refresh republishes the Setup Manifest so the update converges".
  ADR-0165: "A blocking review after archive parks publication for a
  corrective Spec". ADR-0196: "A pre-PR review finding parks only after
  validation". ADR-0169: "The pre-PR review diffs the candidate from its merge
  base". ADR-0282: "a large candidate is reviewed in segments". ADR-0233: "a
  review-only correction returns to review". ADR-0279: "a dismissal survives a
  fix that leaves its anchor". ADR-0237: "A Delivery Retry records a merge
  made outside the queue". ADR-0117: "A defect is checked by the stage that can produce
  it". ADR-0252: "The selective gate runs the repository contracts a change
  makes relevant". ADR-0253: "Every Repository Contract Test runs on main and
  before a release". ADR-0271: "A test bounds by real time only what it
  asserts". ADR-0189: "A version names one content." ADR-0280: "A guide Task
  records its own fingerprint." ADR-0256: "A release waits for the skills that
  describe a changed surface". ADR-0259: "A test runs in parallel unless it
  names why it cannot". ADR-0285: "A repository contract compiles every test
  for Windows." ADR-0104: "A Spec accepts on evidence it did not author".
  ADR-0276: "A row the authorization keeps offline qualifies, and a byte
  invariant compares with the Delivery Base". ADR-0265: "Each claim carries
  its proof as a section." Set no obligation here, because this Spec changes none of their rules: ADR-0173 (History Citations of a Baseline Plan); ADR-0182 (settling on the gate's facts); ADR-0209 (the source-grouping question); ADR-0231 (the Jev Ceiling in User Config); ADR-0239 (each OpenRouter stage's key); ADR-0240 (the QA partial policy); ADR-0250 (module versions and the coverage record's platforms); ADR-0273 (no secret leaves a Run, which ADR-0283 amended); ADR-0278 (retirement deletes); ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0177 (the Relocation Citation scan reads through a repository root); ADR-0208 (work items that share a context share a Spec, and an open item is extended); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0235 (OpenAI and Anthropic models run only through the codex and claude subscriptions); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0238 (a light tier runs low-complexity Tasks on an open model); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0257 (inside a Run the Daemon is the only full gate, and the Baseline states each rule once); ADR-0200, ADR-0201 and ADR-0283 (the advisory judge, its data boundary and its Run Marker, which this Spec leaves unchanged after measuring and rejecting a new question); ADR-0281 (rebuilt executables in tests). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH grant of 2026-10-16 ("Governed Paths declarados") and the 2026-09-30
  skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `internal/docscontract/publicdocs_test.go`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`.

## Current behavior

Characterized on 2026-10-18 at `9471f698` (v0.88.0).

- Review lineage. `decideReviewLineage`
  (`internal/cli/review_lineage.go:35`) has three outcomes:
  - It starts round 1 for a new checkout, provider, merge base or
    non-descending head.
  - It repeats a blocked round and reuses a record at the same head.
  - Otherwise it plans round 2 over the delta from the round-1 head. After
    round 2, or after `ceiling-closed`, it sets `Ceiling` (`:61-63`).

  The ceiling branch of `runReviewCommand` (`internal/cli/review.go:746`) calls
  `closeAtCeiling` (`review_lineage.go:141`). It closes `ceiling-closed` with
  exit 0 when every standing round-2 finding has one valid disposition, `fixed`
  or `dismissed`, and no reviewer reads the change after round 2. Otherwise it
  prints a blocked result and keeps the round-2 record. Further rules:
  - Round 2 diffs the delta with `reviewScopedDiff` and the prompt of
    `buildRoundTwoPrompt` (`review_lineage.go:72`).
  - It is never segmented, a delta above `reviewDiffBound` (`review.go:98`,
    917,504 bytes) blocks, and its session always ends (`review.go:886`).
  - `reviewLineage` (`review_lineage.go:13`) has no field that marks a review
    after round 2.
- Delivery retry.
  - The Delivery Queue parks findings that name an archived Spec as
    `corrective-spec-required: <slugs>`, in either round
    (`internal/delivery/engine.go:980-990`).
  - A retry with a moved head calls `ProveReviewCorrection`
    (`internal/cli/deliver_review_correction.go:18`). It requires:
    - descent from the parked candidate;
    - a Reviewed Head Record with outcome `findings`;
    - exactly one valid disposition per standing finding, `fixed` or
      `dismissed`;
    - every changed path inside the archived Spec or among its Tasks'
      declared paths (`:94-115`). Otherwise it refuses with "changed path %s
      is outside the named archived Specs and the paths their Tasks declared".
  - In Spec 0276's case the correction changed
    `internal/cli/deliver_item_binary_fork_linux_test.go`, which no Task of
    0276 declared, so the retry would have refused it.
  - `roundfix review --base <commit>` reviews a fresh round 1 from that base.
    It overwrites the checkout's current record, so the next default review
    starts a new lineage (`review_lineage.go:37`).
- Spec Consistency Check.
  - The Verification detectors in `internal/speccheck/verification.go` judge a
    command's shape: work independence, inverted exit, wrap fragility,
    hermeticity and vacuity. None reads the paths a command names.
  - The 0276 code `SC-SKILL-VERSION-UNTASKED` shows the wiring every new code
    needs:
    - a detector file;
    - a `stagedDetectors` entry (`internal/speccheck/coherence.go:59`);
    - a call in `Check` (`internal/speccheck/constraints.go:202`);
    - a `corpusFindingCodes` entry (`internal/docscontract/corpus_test.go:294`);
    - a count in `internal/docscontract/testdata/corpus-golden.json`, pinned by
      `TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095`
      (`internal/spec/archive_layout_characterization_test.go`);
    - a row in the `spec` user guide and a bullet in the Roundfix `spec`
      reference.
  - `spec.RecordedTaskPaths` reads `## Recorded paths`.
- Setup Manifest.
  - For the decisions `verification.gate` and `verification.incremental`,
    `validateMakeTarget` (`internal/baseline/profile_alignment.go:1631`)
    returns `contentIdentity` (`:1688`) of the whole Makefile as the
    `declarationDigest`. Any Makefile byte moves both digests in
    `docs/agents/setup-context.json`.
  - `collectReleasePlanBaselineCheck`
    (`internal/cli/releaseplan_checks.go:55`) builds a Managed Refresh plan and
    reports `plan_ready` for any change, but it is read-only and never blocks.
  - No test compares the manifest with a fresh plan.
    `TestRegenerationIsDeclared` only asserts that the paths a regeneration
    changes are declared. It passed on the stale tree in 46 s.
  - At `a5d5a7e1` the plan proposes one change, and at `9471f698` none.
  - The documentation contract package already imports `internal/baseline`.
    It has nine always-relevant files, pinned by
    `TestRepositoryContractTestsDeclareTheirRelevance`.

## System Architecture

- Review (`internal/cli`).
  - `reviewLineage` gains `Correction`.
  - The ceiling branch first asks whether a Correction Review applies, and
    then runs it through the round-2 path with its own prompt.
  - `ProveReviewCorrection` gains one exemption from the path bound.
- Spec check (`internal/speccheck`). A new `verification_path.go` holds
  `SC-VERIFY-PATH-UNDECLARED` and is wired like `SC-SKILL-VERSION-UNTASKED`.
- Documentation contracts (`internal/docscontract`). One test is appended to
  `publicdocs_test.go`.

## Implementation Design

### Interfaces

```go
// internal/cli
// reviewLineage gains: Correction bool `json:"correction,omitempty"`

// internal/speccheck
const CodeVerifyPathUndeclared = "SC-VERIFY-PATH-UNDECLARED"
```

### Data Models

No new file format.

- A Pre-PR Review Record's `lineage` object may carry `"correction": true`.
  Readers that ignore it see a round-2 record.

### API Contracts

1. API Contract: when a Correction Review applies.
   - At the ceiling branch, before `closeAtCeiling`, a Correction Review
     applies when all of these hold:
     - the prior record has outcome `findings`;
     - its lineage is round 2 without `correction`;
     - the current head differs from the prior head and descends from it;
     - every standing prior finding, that is every finding not dismissed by
       validation, has exactly one valid disposition at the prior head, and
       that disposition is `fixed`, with the prior head an ancestor of its
       `fixedBy` and `fixedBy` an ancestor of the current head.
   - Otherwise the ceiling behaves as today.
2. API Contract: what a Correction Review records and sends.
   - It plans `reviewLineage{Round: 2, Correction: true, PreviousHead: <prior
     head>, PreviousFindings: <prior findings>, PreviousDispositions: <the
     fixed dispositions>}`.
   - It runs the round-2 path: the delta `<prior head>..<head>` from
     `reviewScopedDiff`, never segmented, blocked above `reviewDiffBound`, in
     a new Agent Session.
   - Its prompt, from `buildCorrectionReviewPrompt`, opens with the round-2
     header. It continues: `This is the Correction Review of this candidate:
     round 2 reviewed <prior head> and raised the findings below, each fixed
     by the commit named. The diff below is the change since <prior head>; it
     is reviewed once and no further review of this candidate follows.`
   - It then lists each finding as `- <id> | fixed by <commit> | <text>`,
     followed by `--- BEGIN CORRECTION DELTA DIFF ---` and `--- END CORRECTION
     DELTA DIFF ---`.
   - Validation, outcomes, the Reviewed Head Record and every other record
     field follow round 2.
3. API Contract: after a Correction Review.
   - A clean or findings-dismissed Correction Review lets the Delivery Queue
     advance exactly as a clean round 2 does.
   - A findings Correction Review parks exactly as round-2 findings park:
     `corrective-spec-required: <slugs>` for archived Specs, `review-findings`
     otherwise.
   - A blocked Correction Review repeats as a Correction Review.
   - The next review of a lineage whose prior record has `correction: true`
     reaches the ceiling and `closeAtCeiling` exactly as after round 2. It
     never plans a second Correction Review.
4. API Contract: the retry of a round-2 park.
   - `ProveReviewCorrection` skips the declared-path bound when all of these
     hold:
     - the parked record's lineage is round 2 without `correction`;
     - every standing finding's one valid disposition is `fixed`.
   - Every other check is unchanged: descent, the record, one valid
     disposition per finding.
   - A park by round 1, by a Correction Review, or by a round 2 with any
     `dismissed` disposition keeps the bound and today's refusal text.
5. API Contract: `SC-VERIFY-PATH-UNDECLARED`.
   - An error from the Tasks stage, with no horizon. It skips repositories
     without a Task Graph.
   - For each pending Task that is not QA, and each Verification command, the
     command is split into shell words, and the operators `;`, `&&`, `||`,
     `|`, `(`, `)` and redirections begin a new command. A word is a candidate
     path when:
     - after an optional leading `./`, it matches path segments of letters,
       digits, `.`, `_`, `-`, `@` and `+`, joined by `/` (at least two
       segments, an optional trailing `/`);
     - its first segment is an existing top-level directory of the
       repository;
     - it holds no `$`, `*`, `?`, `[`, `<` or `...`.
   - Never a candidate:
     - the pattern operand of `grep` or `rg` (the first non-option operand
       when no `-e` or `-f` is given, else the operand of each `-e` or `-f`);
     - every operand of `printf` and `echo`;
     - the operand of a `test` or `[` command after `!`;
     - a path under the Spec Root.
   - A candidate is reported when all of these hold:
     - it is not a declared path of the Task (`creates:`, `interface:`,
       `deletes:`, `instruction:`, or `## Recorded paths`);
     - it is not a directory prefix of one;
     - it is not a `creates:` path, or a directory prefix of one, of a Task in
       the Task's dependency closure;
     - it does not exist in the checked tree (`os.Lstat` under the repository
       root).
   - The finding:
     - Summary: `Task <id> Verification names <path>, which the Task does not
       declare, no Task it depends on creates, and the checked tree does not
       hold`.
     - Location: the Task file at the Verification command's line.
     - Fix: `Declare <path> in this Task, create it in a Task this Task needs,
       or correct the command, following the write-tasks rule "Verification
       paths".`
6. API Contract: the Setup Manifest contract.
   - `TestTheSetupManifestIsCurrent`, appended to
     `internal/docscontract/publicdocs_test.go`, builds the plan from the
     repository root through these calls:
     - `baseline.LoadEmbeddedCatalog`;
     - `baseline.ResolveManifestInput`;
     - `baseline.BuildPlan` with `PreservationModeManagedRefresh` and the
       directories of `PATH`.
   - It fails when any call fails, when the input needs new decisions, when
     no plan is returned (naming the state and message), or when the plan
     holds any file change or history move.
   - The failure message names each changed path and the command
     `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes`.
   - It writes nothing.
7. API Contract: everything else stays. These keep their behavior: round 1,
   round 2, Review Segments, `ceiling-closed` with any dismissal, the
   corrective-Spec park, every existing refusal text, every judge question
   and threshold, the release plan, the Managed Refresh, and every
   existing test name.

### Surface Transcripts

None. No command gains a new output line of its own: the Correction Review
writes a record field and `SC-VERIFY-PATH-UNDECLARED` uses the check's
existing finding format; API Contracts 2 and 5 fix each text and the named
tests assert it.

### Invariants

```text
1. A Reviewer Lineage calls the reviewer at most once after round 2, and only when every standing round-2 finding was fixed by a commit the head contains.
2. A Correction Review reviews only the delta from the round-2 head.
3. A pending non-QA Task never reaches a Run with a Verification path it neither declares, receives from a Task it needs, nor finds in the tree.
4. Every Task settlement fails while the Setup Manifest differs from a fresh Managed Refresh.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, User Stories 1, Core Features 1-2, Success Metrics 1-2 → `reviewLineage.Correction`, the ceiling branch, `buildCorrectionReviewPrompt`, `ProveReviewCorrection` (API Contracts 1-4)
- Goal 2, User Story 2, Core Feature 3, Success Metric 3 → `verification_path.go` and its wiring (API Contract 5)
- Goal 3, User Story 3, Core Feature 4, Success Metric 4 → `TestTheSetupManifestIsCurrent` (API Contract 6)
- Core Feature 5 → the `review`, `deliver` and `spec` guides, the Roundfix references, the `write-tasks` skill, `CONTEXT.md`
- Success Metrics 1-4, API Contract 7 and Invariants 1-4 → the named tests and the scope audit

## Absence Evidence

1. Correction Review — command: `grep -rn 'json:"correction' internal/cli`; result: no match at `9471f698`.
2. The declared-path exemption — command: `grep -n 'Lineage' internal/cli/deliver_review_correction.go`; result: no match at `9471f698`.
3. `SC-VERIFY-PATH-UNDECLARED` — command: `grep -rn 'SC-VERIFY-PATH-UNDECLARED' internal --include='*.go'`; result: no match at `9471f698`.
4. The Setup Manifest contract — command: `grep -rln 'BuildPlan' internal/docscontract`; result: no match at `9471f698`.

## ADR Alternatives

- ADR-0286: this Spec's decision; it rejects a third round for any disposition, a corrective Spec for every round-2 fix, keeping the declared-path bound for a round-2 fix, and the name "round 3".
- ADR-0287: this Spec's decision; it rejects treating every path-shaped word as a path, requiring existing paths to be declared, the measured QA-delivery question, a Full-Contract-only manifest check and a new Repository Contract Test file.
- ADR-0197: amended by ADR-0286 with a dated note; the ceiling stays for every lineage except the one Correction Review, and no rejected alternative of ADR-0197 is taken.
- ADR-0272: amended by ADR-0286 with a dated note; it takes, for a round-2 park whose findings are all fixed only, the rejected "Accepting any code correction", because the Correction Review reviews that correction, which ADR-0272's rejection assumed nobody would.
- ADR-0103: amended by ADR-0287 with a dated note; no rejected alternative is taken; the manifest still converges through the Managed Refresh, now checked at each settlement.
- ADR-0165: no rejected alternative is taken; findings that name an archived Spec still park for a corrective Spec.
- ADR-0196: no rejected alternative is taken; a Correction Review's findings are validated like any round's.
- ADR-0169: no rejected alternative is taken; the lineage still requires one merge base.
- ADR-0282: no rejected alternative is taken; a Correction Review, like round 2, is never segmented.
- ADR-0233: no rejected alternative is taken; a review-only correction still returns to review.
- ADR-0279: no rejected alternative is taken; dismissals keep their anchors and never open a Correction Review.
- ADR-0237: no rejected alternative is taken; the Delivery Retry still records a merge made outside the queue.
- ADR-0200: no rejected alternative is taken; the judge gains no question.
- ADR-0201: no rejected alternative is taken; the judge's data boundary is unchanged.
- ADR-0283: no rejected alternative is taken; the Run Marker is unchanged, and the source-3 measurement ran outside any Run.
- ADR-0117: no rejected alternative is taken; the path defect is checked at the Tasks stage that produces it.
- ADR-0252: no rejected alternative is taken; the new contract test lives in an always-relevant file.
- ADR-0253: no rejected alternative is taken; the Full Contract Run discovers the new test with no list.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.
- ADR-0189: no rejected alternative is taken; each skill edit raises its recorded version.
- ADR-0280: no rejected alternative is taken; each guide Task records its own fingerprint.
- ADR-0256: no rejected alternative is taken; each changed surface's skill reference changes with it.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0285: no rejected alternative is taken; every new test compiles for Windows.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on repository history and published sources.
- ADR-0276: no rejected alternative is taken; network-free rows qualify and byte invariants compare with the Delivery Base.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.

- ADR-0173: no rejected alternative is taken; a Baseline Plan's History Citations are unchanged.
- ADR-0182: no rejected alternative is taken; a Task still settles on the facts its gate checks.
- ADR-0209: no rejected alternative is taken; the source-grouping question is unchanged.
- ADR-0231: no rejected alternative is taken; the Jev Ceiling is still read from User Config.
- ADR-0239: no rejected alternative is taken; each OpenRouter stage still reads its own key first.
- ADR-0240: no rejected alternative is taken; the QA partial policy is unchanged.
- ADR-0250: no rejected alternative is taken; no test changes its platform set, so the coverage record is not re-recorded.
- ADR-0273: no rejected alternative is taken; no secret and no Spec content leaves a Run.
- ADR-0278: no rejected alternative is taken; this Spec retires nothing.
- ADR-0281: no rejected alternative is taken; no test rebuilds an executable.
- ADR-0092: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0120: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0177: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0208: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0235: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0247: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0248: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0254: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0258: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0260: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0261: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0277: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0238: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0262: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0269: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0257: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

None outside the repository. The Correction Review calls the configured review
provider exactly as round 2 does.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, builds temporary
repositories with the fixtures already used beside it, and compiles for
Windows.

- `internal/cli/review_correction_test.go`:
  - `TestCorrectionReviewReviewsTheFixOnce`: round 2 leaves two findings, both
    disposed `--fixed-by`. The next review calls the fake reviewer exactly once
    with the round-2-head-to-head delta and the Correction Review prompt, and
    records `correction: true`. A further review never calls the reviewer.
  - `TestCorrectionReviewFindingsReachTheCeiling`: a findings answer exits as
    round-2 findings do, and the next review closes at the ceiling.
  - `TestADismissedRoundTwoFindingKeepsTheCeiling`: one `dismissed`
    disposition closes `ceiling-closed` with no reviewer call, as today.
- `internal/cli/deliver_review_correction_test.go`:
  - `TestRetryAdmitsAFixedRoundTwoParkOutsideDeclaredPaths`: a round-2 record
    whose findings are all fixed is accepted with a changed undeclared path.
  - `TestRetryKeepsTheBoundForADismissedRoundTwoFinding`: the same fix with one
    dismissal is refused with today's text.
- `internal/speccheck/verification_path_test.go`:
  - `TestVerificationPathUndeclaredRaisesTheMissingFile`: the 0274 shape.
  - `TestVerificationPathUndeclaredIgnoresWhatThePlanProvides`: declared,
    recorded, created by a needed Task, directory prefix, existing, grep
    pattern, printf operand, negated test and Spec Root paths.
- `internal/docscontract/publicdocs_test.go`: `TestTheSetupManifestIsCurrent`.

## Build Order

1. The Correction Review and the retry exemption, with their tests, the
   `review` and `deliver` guides and references, the glossary, the skill
   version and the Behavior Surface Record (no dependency).
2. `SC-VERIFY-PATH-UNDECLARED` with its tests, wiring, corpus golden and pin,
   the `spec` guide row, the Roundfix `spec` reference, the `write-tasks` rule,
   the skill versions and the record (depends on: 1, which raises the same
   Roundfix version and record).
3. The Setup Manifest contract (no dependency).
4. QA gate (depends on: 2, 3).

No other Spec is a prerequisite.

## Risks & Considerations

- A path the rule cannot recognize, such as one built from a variable, stays
  unchecked. The rule prefers silence to a false positive.
- A Task that edits the Makefile must regenerate the Setup Manifest in the
  same Task. ADR-0287 records this, and the test's message names the command.
- Steps 1 and 2 each raise the Roundfix skill version and re-record shared
  files, so they run in series. Step 3 shares no file with them.

## Glossary

- adds: **Correction Review**
- changes: **Reviewer Lineage**

## Decisions

- A Correction Review is a round-2 lineage with `correction: true`, never a
  round 3. See ADR-0286.
- The retry exemption applies only to a round-2 park whose findings are all
  fixed. See ADR-0286.
- The path rule is conservative and uses the dependency closure for created
  paths. See ADR-0287.
- The QA-delivery question is not adopted: measured on 2026-10-18 through
  `typesafe/jev-1.13-20260917`, it found neither real case and raised 14 of 72
  negatives. See ADR-0287.
- The manifest check is appended to `publicdocs_test.go`. See ADR-0287.
