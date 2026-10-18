---
spec: 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
prd: _prd.md
created: 2026-10-18
---

# QA evidence keeps its images, and a secret is refused where it is born — Technical Spec

## Executive Summary

Three slices extend existing seams, and no package is new.

1. QA evidence images (ADR-0291).
   - `internal/speccheck` gains one exported classifier, `ReadableEvidence`,
     which the mechanical stage uses instead of its NUL-byte test.
   - The Daemon's QA Report commit uses the same classifier on the paths it
     would stage under the report's `evidence/` directory and refuses a
     non-image binary there.
2. The Task commit Secret Scan (ADR-0292). The Daemon scans the lines a Task
   commit would add before the Task settles, through a new
   `secretscan.ScanAddedLines` that reads a `git diff` patch.
3. The Secret Rewrite (ADR-0292).
   - The Delivery Retry of an archived `secret-found` item asks a new
     `SecretRewriteProver`, implemented in `internal/cli`, whether a
     non-descending head is a Secret Rewrite.
   - `ClassifyPark` prints the rewrite that works.

The trade-offs accepted:

- A gate that saves a video or a PDF under `qa/evidence/` now fails at its own
  report commit, instead of at the next gate.
- A Task with a secret-shaped fixture fails without a repair turn.
- A Secret Rewrite is reviewed again from round 1.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  record schema is minted; the QA Task and the Task settlement each gain one
  refusal reason, and the Delivery Retry one acceptance rule. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every new read is a local file or a Git
  object, and every test uses temporary repositories and fakes. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — each Build Order step follows recorded
  decisions. ADR-0291 governs Build Order 1: "QA evidence holds text and
  images, and the QA Report commit refuses any other binary". ADR-0292
  governs Build Orders 2 and 3: "A Task commit is scanned for secrets, and a
  retry accepts a Secret Rewrite". ADR-0273, amended by both: "No secret and
  no Spec content leaves the machine from a Run". ADR-0283: "A Run Marker
  keeps the judge silent inside a Run". ADR-0237: "A Delivery Retry records a
  merge made outside the queue". ADR-0272: "A review record survives a fix,
  and a correction inside the archived Spec's declared paths returns to
  review". ADR-0197: "A pre-PR reviewer lineage spans at most two rounds".
  ADR-0117: "A defect is checked by the stage that can produce it". ADR-0189:
  "A version names one content." ADR-0280: "A guide Task records its own
  fingerprint." ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0259: "A test runs in parallel unless it names why it
  cannot". ADR-0285: "A repository contract compiles every test for
  Windows." ADR-0125: "A spawned fixture is a compiled binary, never a written
  script". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0276:
  "A row the authorization keeps offline qualifies, and a byte invariant
  compares with the Delivery Base". ADR-0265: "Each claim carries its proof as
  a section." ADR-0193: "A Spec names its prerequisites, and the owner waits
  for them". Set no obligation here, because this Spec changes none of their rules: ADR-0092 (triage intent lives in a typed backlog); ADR-0120 (retired documentation lives under one history root); ADR-0165 (a blocking review after archive parks publication for a corrective Spec); ADR-0167 (the pre-PR Pull Request row never decides a qualifying partial); ADR-0169 (the pre-PR review diffs the candidate from its merge base); ADR-0182 (a Task settles on the facts its gate will check); ADR-0196 (a pre-PR review finding parks only after validation); ADR-0229 (an operator archive resumes any park from the Run start); ADR-0233 (a skill version raise is regenerated at merge, and a review-only correction returns to review); ADR-0240 (one QA partial policy, and rows a Run sandbox cannot reach); ADR-0247 (an archive leaves an Archive Record, and the Spec folder stays in Git); ADR-0248 (existing history is sanitized in batches after a history-full tag); ADR-0250 (a module version is chosen when it is recorded, and the coverage record lists every platform); ADR-0251 (a Legacy Archive Folder is read leniently, and a failed QA keeps its verdict); ADR-0254 (baseline update sanitizes pending history in the change it plans); ADR-0258 (a Spec is authored against the QA rerun classes, and a retirement writes reduced history); ADR-0260 (an archive relinks live citations, and a Baseline update lists the citations it breaks); ADR-0261 (an Archive Record names its pull request, and readers fetch the pull request head); ADR-0262 (a cancelled Spec leaves an Archive Record and returns its sources); ADR-0269 (archive and cancel refuse ignored Spec files and finish a failed removal forward); ADR-0271 (a test bounds by real time only what it asserts); ADR-0277 (a Legacy Archive Folder conversion keeps the summary, links and dates its folder held); ADR-0278 (a retirement deletes, and the History Root holds only Archive Records); ADR-0279 (status reports a Queue Owner proven gone, and a dismissal survives a fix that leaves its anchor); ADR-0281 (a test replaces an executable it rebuilds instead of writing into it); ADR-0282 (the history cleanup names a superseded ADR in its successor, and a large candidate is reviewed in segments); ADR-0286 (a fixed round-2 finding gets one Correction Review of its fix). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  grant of 2026-10-16 ("Governed Paths declarados") and the 2026-09-30 skills
  grant ("considere autorizado a ajustar todas as skills se necessário").
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born/_authorization.md`;
  bounded files: `.agents/skills/implement-task/SKILL.md`,
  `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `skills/implement-task/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-18 at `2ef4f080` (v0.89.0).

- Mechanical evidence check (`internal/speccheck/mechanical.go`).
  - `RunMechanicalStage` (`:539`) calls `detectMechanicalEvidenceScratchState`
    (`:585`, defined at `:2501`) for the newest QA Report.
  - It lists tracked files under the report's `evidence/` directory with
    `git ls-files --stage` (`mechanicalTrackedEvidence`, `:2541`). A gitlink
    raises `QA-EVIDENCE-PATH` "gitlink", and a regular blob that
    `mechanicalBlobIsBinary` (`:2574`) classifies as binary raises it as
    "built binary" (`:2529-2535`).
  - `mechanicalBlobIsBinary` reads the first 8000 bytes through `git cat-file
    blob` and returns whether they hold a NUL byte (`:2599`). Every PNG and
    every JFIF JPEG holds one in its header, so every screenshot is refused.
  - The finding's fix text is "Remove the gate's scratch state from the Spec
    evidence directory and keep only reader-openable artifacts." (`:2608`).
  - The stage runs before the QA Agent's prompt
    (`internal/daemon/task_engine.go:2786`, prompt at `:2880`), so files the
    gate writes during its turn are invisible to it.
  - `TestEvidenceRefusesScratchStateBinary`
    (`internal/speccheck/mechanical_test.go:2025`) commits `\x7fELF\x00…` and
    asserts one finding whose detail contains `binary`.
- QA Report commit (`internal/daemon/task_engine.go`).
  - After the verdict, `prepareQAReportCommit` (`:3006`, defined at `:3715`)
    computes the staged paths.
  - `scanQAReportPaths` (`:3012`, defined at `:3058`) runs the Secret Scan on
    them, before and after the QA format command (`:3016-3023`).
  - A Secret Match sets the reason `QA evidence carries a secret: <match>;
    ...` (`:3031`), settles the QA Task `failed` (`:3033`), publishes one
    `daemon.qa` event with phase `secret_scan` (`:3036-3049`) and returns
    before `commitPreparedQAReport` (`:3050`).
  - `secretscan.ScanContent` (`internal/secretscan/scan.go:41`) returns no
    match for content with a NUL byte, so no step at this commit looks at a
    binary.
- Task settlement (`internal/daemon/task_engine.go`).
  - `executeTask` (`:1145-1227`) runs Verification, then
    `prepareTaskCommit` (`:1185`, defined at `:2125`). Lost stage paths and
    the authorization checks can still set a failure (`:1201-1208`).
  - It then settles (`:1213`) and calls `commitTask` (`:2211`). No Secret
    Scan runs.
  - The Secret Scan reads commits only at the Daemon's Final Push
    (`internal/daemon/daemon.go:343-350`) and the Delivery Queue's push
    (`internal/delivery/github.go:179-188`). Both use
    `secretscan.PushedHistoryArgs` and `ScanPushedHistory` (`scan.go:56-117`),
    which report only added lines and only after a `commit <id>` line.
- Secret Patterns (`internal/secretscan/patterns.go`).
  - `credential-assignment` (`:23`) matches an optional upper-case prefix
    ending in `API_KEY`, `TOKEN`, `SECRET` or `PASSWORD`, an assignment
    operator, an optional quote, and a value of at least 20 characters from
    `[A-Za-z0-9_+/.=-]`.
  - `matchesKind` (`:26-46`) reports it only when the value holds a letter and
    a digit. A probe over every tracked file outside `docs/history/` at
    `2ef4f080` found no Secret Match.
- Delivery Retry of a `secret-found` item (`internal/delivery`).
  - The push parks `secret-found: <match>, <match>` (`engine.go:1114-1123`),
    each match in `secretscan.Match.String` form (`<path>:<line>: <kind>
    (commit <id>)`).
  - `ClassifyPark` prints "remove the secret from every commit the blocker
    names on the item branch in <worktree>, then run roundfix deliver retry
    <slug>" (`park_class.go:65-67`). It is pinned by
    `internal/delivery/park_class_test.go:88` and
    `internal/delivery/push_secret_scan_test.go:114`.
  - `Retry` treats an archived item in its `state.Archived` branch
    (`engine.go:547-582`). A moved head is accepted only when it descends
    from the newest candidate, or from the Run start when no candidate is
    recorded (`:554-567`). Otherwise it refuses with "retry Delivery Queue
    item %q: archived item head %q differs from candidate head %q"
    (`:569-574`).
  - The `corrective-spec-required` branch shows how a proof is wired:
    - `Dependencies.Corrections` (`engine.go:253`) and its
      `ReviewCorrectionProver` interface (`:161-163`);
    - the implementation on `commandDeliveryWorkflow`
      (`internal/cli/deliver_review_correction.go:16-18`), wired at
      `internal/cli/deliver_workflow.go:67`;
    - the refusal suffix ` (<reason>)` (`engine.go:503-526`).
  - `ItemHistory.RunStart` (`internal/cli/deliver_workflow.go:1030`) returns
    the Implement Run's start head.
  - A review of a head that does not descend from the reviewed head starts a
    new round 1 (`decideReviewLineage`, `internal/cli/review_lineage.go:36-45`).
- Documentation.
  - The `implement` guide states the QA Report commit Secret Scan
    (`docs/user-guide/commands/implement.md:157-161`), as does the Roundfix
    `implement` reference (`.agents/skills/roundfix/references/implement.md:167-171`).
  - The `deliver` guide's Park Class table carries the `secret-found` next
    action (`docs/user-guide/commands/deliver.md:408`), as does the Roundfix
    `deliver` reference (`.agents/skills/roundfix/references/deliver.md:323`).
    The guide's retry table is at `deliver.md:521-533`.
  - The `qa-gate` skill says "Store screenshots and focused artifacts under
    `qa/evidence/<YYYY-MM-DD>-<run-slug>/`" (`.agents/skills/qa-gate/SKILL.md:468`)
    and states the Secret Scan refusal under "Evidence carries no secret"
    (`:44-52`).
  - The glossary's **Mechanical Refusal Code** (`CONTEXT.md:381`) describes
    `QA-EVIDENCE-PATH` only as a named path the repository does not carry, and
    **Secret Scan** (`CONTEXT.md:916`) names the QA Report commit and pushes
    only.
  - The test `qaSecretRefusalDocumentation`
    (`internal/docscontract/command_documentation_test.go:153-200`) requires
    exactly one Daemon literal starting with `QA evidence carries a secret`.

## System Architecture

- `internal/speccheck`: a new `evidence_image.go` holds `ReadableEvidence`
  and `EvidenceProbeBytes`; `mechanical.go` uses them.
- `internal/daemon`: `task_engine.go` gains the QA evidence binary refusal
  and the Task commit Secret Scan.
- `internal/secretscan`: `scan.go` gains `ScanAddedLines`, sharing
  `ScanPushedHistory`'s hunk parser.
- `internal/delivery`: `engine.go` gains `SecretRewrite`,
  `SecretRewriteProver`, `Dependencies.SecretRewrites` and the retry branch,
  and `park_class.go` the new next action.
- `internal/cli`: a new `deliver_secret_rewrite.go` implements the prover on
  `commandDeliveryWorkflow`, wired in `deliver_workflow.go`.

## Implementation Design

### Interfaces

```go
// internal/speccheck
const EvidenceProbeBytes = 8000
func ReadableEvidence(probe []byte) bool

// internal/secretscan
func ScanAddedLines(patch []byte) []Match

// internal/delivery
type SecretRewrite struct {
	Accepted bool
	Reason   string
}
type SecretRewriteProver interface {
	ProveSecretRewrite(ctx context.Context, workDir, base, candidate, head string, named []string) (SecretRewrite, error)
}
// Dependencies gains: SecretRewrites SecretRewriteProver
```

### Data Models

No new file format. A Delivery Queue item keeps its fields; an accepted
Secret Rewrite appends the head to `CandidateCommits`, as every accepted
moved head does.

### API Contracts

1. API Contract: `ReadableEvidence`.
   - It reads at most `EvidenceProbeBytes` bytes of a file's beginning and
     returns true when any of these holds:
     - the probe holds no NUL byte (text, SVG included);
     - the probe begins with `89 50 4E 47 0D 0A 1A 0A` (PNG);
     - the probe begins with `FF D8 FF` (JPEG).
   - Otherwise it returns false. It never reads a name or an extension.
2. API Contract: the mechanical evidence check.
   - `detectMechanicalEvidenceScratchState` keeps the gitlink finding and
     raises `QA-EVIDENCE-PATH` for a regular blob only when `ReadableEvidence`
     of its first `EvidenceProbeBytes` bytes is false.
   - The detail becomes `tracked evidence path <path> is a non-image binary`.
     The fix becomes `Remove the gate's scratch state from the Spec evidence
     directory and keep only text, PNG, JPEG or SVG artifacts.`
3. API Contract: the QA Report commit refusal.
   - When the Secret Scan of the prepared paths finds no match, the Daemon
     reads each prepared path under `path.Dir(<report path>) + "/evidence/"`
     that exists. Any path whose first `EvidenceProbeBytes` bytes are not
     `ReadableEvidence` is a non-image binary.
   - The QA Task settles `failed` with the reason `QA evidence carries a
     non-image binary: <path>; <path>`, in path order.
   - The Daemon publishes one `daemon.qa` event with the summary `QA Report
     refused: <n> non-image binary file(s)` and the payload `{"phase":
     "evidence_binary", "paths": [...]}`, prints the summary to progress, and
     creates no QA Report commit; the files stay uncommitted in the Run
     Worktree.
   - A read error is an infrastructure failure, as in `scanQAReportPaths`.
4. API Contract: `ScanAddedLines`.
   - It reads one `git diff` patch, with no `commit` lines, and returns a
     Match for each added line that matches a Secret Pattern. Each Match has
     its path and its line in the new version, and `Commit` empty.
   - `ScanPushedHistory` keeps its behavior and output.
5. API Contract: the Task commit Secret Scan.
   - In `executeTask`, after `prepareTaskCommit` and after the lost-path and
     authorization checks leave no failure, the Daemon scans the prepared
     stageable paths:
     - a path that does not exist in the Run Worktree contributes nothing;
     - a path `HEAD` does not hold is read whole with `ScanContent`;
     - every other path is read from one `git diff --no-color --no-ext-diff
       --no-renames -U0 HEAD -- <paths>` with `ScanAddedLines`.
   - A Secret Match sets the failure `Task commit carries a secret: <match>;
     <match>`, using `Match.String`. The Daemon publishes one `daemon.task`
     event with the summary `Task <id> commit refused: <n> Secret Match(es)`
     and the payload `{"phase": "secret_scan", "matches": [{"path", "line",
     "kind"}]}`.
   - The Task then settles `failed` through the existing path, so no Task
     commit is made and its changes stay in the Run Worktree.
   - A Git or read error is an infrastructure failure of the Run.
6. API Contract: the Secret Rewrite retry.
   - In `Retry`'s archived branch, when a moved head is not accepted by
     descent, `blockerMatches(item.Blocker, BlockerSecretFound)` holds, the
     candidate is known, the head is not empty and `SecretRewrites` is set,
     the engine resolves the base with `history.RunStart(item.RunID)`.
   - It reads the named paths from the blocker. Each entry after
     `secret-found: ` is separated by `, ` and has the form `<path>:<line>:
     <kind>` with an optional ` (commit <id>)`; duplicates collapse.
   - It calls `ProveSecretRewrite(workDir, base, candidate, head, named)`.
   - An accepted proof appends the head to `CandidateCommits` and sets the
     stage `reviewing`.
   - Otherwise the refusal is today's reason followed by ` (<reason>)`, where
     the reason is the proof's reason, the prover's error, or `the item has no
     Run start to bound the rewrite` when the base cannot be resolved.
   - Every other blocker keeps today's text with no suffix.
7. API Contract: `ProveSecretRewrite` (`internal/cli`). It refuses with the
   first failing rule, in this order:
   - `item head does not descend from the Run start <base>` when the base is
     not an ancestor of the head;
   - `the rewrite changes <path>, which the blocker does not name` for a path
     in `git diff --raw --no-renames candidate head` that is not named;
   - `the rewrite adds, removes or changes the mode of <path>` for a named
     path whose diff status is not `M` or whose mode changed;
   - `the rewrite changes the line count of <path>` when the two versions,
     split on `\n`, differ in length;
   - `the rewrite changes line <n> of <path>, which holds no Secret Match at
     the parked candidate` for a differing line whose candidate line
     `ScanContent` does not match;
   - `the rewrite keeps a Secret Match at line <n> of <path>` for any line of
     a named path's head version that `ScanContent` matches.

   Otherwise it accepts. It reads Git objects only and writes nothing.
8. API Contract: the `secret-found` next action. `ClassifyPark` prints `in
   <worktree>, rewrite the item branch from the earliest commit the blocker
   names so that each named line holds a value no Secret Pattern matches and
   nothing else in the final tree changes, then run roundfix deliver retry
   <slug>`.
9. API Contract: everything else stays. These keep their behavior: every
   other Mechanical Refusal Code, the QA Report Secret Scan and its text, the
   QA format command, the Final Push and Delivery Queue push scans, every
   Secret Pattern, every other blocker's retry, the Verification retry
   policy, and every existing test name.

### Surface Transcripts

None. No command gains an output line of its own: the refusals extend
existing reasons, events and the park table; API Contracts 3, 5, 6, 7 and 8
fix each text and the named tests assert it.

### Invariants

```text
1. A file the QA Report commit stages under qa/evidence/ is a file the next mechanical stage accepts.
2. A Task commit that settles completed adds no line the Final Push scan would match.
3. A Delivery Retry accepts a non-descending head of a secret-found item only when the trees differ solely at lines that matched at the parked candidate and match nowhere in the named files at the head.
4. No Secret Pattern, no Secret Match format and no other retry changes.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Goal 2, User Stories 1-2, Core Features 1-2, Success Metrics 1-2 → `ReadableEvidence`, the mechanical check, the QA Report commit refusal (API Contracts 1-3)
- Goal 3, User Story 3, Core Feature 3, Success Metric 3 → `ScanAddedLines` and the Task commit scan (API Contracts 4-5)
- Goal 4, User Story 4, Core Feature 4, Success Metric 4 → `SecretRewriteProver`, the retry branch, the next action (API Contracts 6-8)
- Core Feature 5 → the `implement` and `deliver` guides, the Roundfix references, the `qa-gate` and `implement-task` skills, `CONTEXT.md`
- Success Metrics 1-4, API Contract 9 and Invariants 1-4 → the named tests and the scope audit

## Absence Evidence

1. The readable-evidence classifier — command: `grep -rn 'func ReadableEvidence' internal/speccheck`; result: no match at `2ef4f080`.
2. The QA Report binary refusal — command: `grep -rn 'non-image binary' internal/daemon`; result: no match at `2ef4f080`.
3. The Task commit Secret Scan — command: `grep -rn 'Task commit carries a secret' internal`; result: no match at `2ef4f080`.
4. The added-lines scanner — command: `grep -rn 'func ScanAddedLines' internal/secretscan`; result: no match at `2ef4f080`.
5. The Secret Rewrite — command: `grep -rn 'SecretRewrite' internal`; result: no match at `2ef4f080`.

## ADR Alternatives

- ADR-0291: this Spec's decision; it rejects text-only evidence, detection by extension, accepting every viewable binary, dropping the binary silently and an SVG XML-prefix test.
- ADR-0292: this Spec's decision; it rejects whole-file and line-number comparisons, a value-span comparison, a reset-only next action, exempting tests or fixture-looking values, a whole-file Task scan and a repair turn.
- ADR-0273: amended by ADR-0291 and ADR-0292 with a dated note; it takes, for the Task commit, the rejected "Scanning only `qa/evidence/**`" in reverse, since its own reason, that "any implementation commit can carry a value too", is what the Task commit scan answers.
- ADR-0283: no rejected alternative is taken; no new process reads a judge key.
- ADR-0237: no rejected alternative is taken; a Delivery Retry still records a merge made outside the queue.
- ADR-0272: no rejected alternative is taken; the `corrective-spec-required` correction proof is unchanged.
- ADR-0197: no rejected alternative is taken; a rewritten head starts a new lineage as any non-descending head does.
- ADR-0117: no rejected alternative is taken; each defect is now checked at the commit that produces it.
- ADR-0189: no rejected alternative is taken; each skill edit raises its recorded version.
- ADR-0280: no rejected alternative is taken; each guide Task records its own fingerprint.
- ADR-0256: no rejected alternative is taken; each changed surface's skill text changes with it.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0285: no rejected alternative is taken; every new test compiles for Windows.
- ADR-0125: no rejected alternative is taken; no test writes an executable, and binary fixtures are byte strings built at runtime.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on repository history and published sources.
- ADR-0276: no rejected alternative is taken; network-free rows qualify and byte invariants compare with the Delivery Base.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.
- ADR-0193: no rejected alternative is taken; this Spec names Spec 0279 as its prerequisite.
- ADR-0182: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0240: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0250: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0262: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0271: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0279: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0282: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0269: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0278: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0281: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0092: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0120: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0247: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0248: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0254: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0258: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0260: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0261: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0277: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0165: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0169: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0196: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

None outside the repository. The Secret Scan, the classifier and the proof
read local files and Git objects only.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, builds temporary
repositories with the fixtures already used beside it, and compiles for
Windows. Secret samples are built at runtime from parts, as ADR-0273
requires, and binary samples are byte strings built in the test.

- `internal/speccheck/evidence_image_test.go`:
  - `TestReadableEvidenceDecidesByContent`: text, SVG, PNG and JPEG probes are
    readable; an ELF header, a SQLite header, a PDF with a NUL byte and an
    ELF header named `.png` are not.
  - `TestEvidenceAcceptsTrackedImages`: tracked PNG, JPEG, SVG and an
    extensionless PNG raise no `QA-EVIDENCE-PATH`.
  - `TestEvidenceRefusesNonImageBinaries`: an ELF header, a SQLite header, a
    PDF and an ELF header named `.png` each raise one finding whose detail
    names the path and `non-image binary`.
- `internal/daemon/qa_evidence_binary_test.go`:
  - `TestQAReportCommitsImageEvidence`: a PNG under the evidence directory
    rides in the QA Report commit.
  - `TestQAReportWithANonImageBinarySettlesFailedWithoutACommit`: the QA Task
    settles `failed` with the exact reason, the event names the path, and no
    QA Report commit exists.
- `internal/secretscan/added_lines_test.go`:
  `TestScanAddedLinesReportsOnlyAddedLines`.
- `internal/daemon/task_secret_scan_test.go`:
  - `TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit`;
  - `TestTaskCommitScansOnlyTheLinesItAdds`: an edit beside an already
    committed match commits;
  - `TestTaskCommitAcceptsALettersOnlyFixture`.
- `internal/cli/deliver_secret_rewrite_test.go`:
  - `TestSecretRewriteAcceptsReplacedMatchedLines`;
  - `TestSecretRewriteRefusesAnyOtherChange`, with one subtest per rule of
    API Contract 7.
- `internal/delivery/secret_rewrite_retry_test.go`:
  - `TestRetryAcceptsASecretRewriteOfAnArchivedItem`;
  - `TestRetryRefusesARejectedSecretRewriteWithItsRule`.
- The pinned next action changes in `internal/delivery/park_class_test.go` and
  `internal/delivery/push_secret_scan_test.go`, to API Contract 8's text.

## Build Order

1. QA evidence images: `ReadableEvidence`, the mechanical check, the QA Report
   commit refusal, their tests, the `implement` guide and reference, the
   `qa-gate` skill, the glossary, the skill versions and the Behavior Surface
   Record (no dependency).
2. The Task commit Secret Scan: `ScanAddedLines`, the settlement scan, their
   tests, the `implement` guide and reference, the fixture guidance in the
   guide, the reference and the `implement-task` skill, the glossary, the
   skill versions and the record (depends on: 1, which changes the same
   `task_engine.go`, guide, reference, skill version and record).
3. The Secret Rewrite: the prover, the retry branch, the next action, their
   tests and the two pinned tests, the `deliver` guide and reference, the
   glossary, the skill version and the record (depends on: 2, which raises
   the same Roundfix version and record).
4. QA gate (depends on: 3).

Spec 0279 is a prerequisite and is named in the Task Graph's `requires`. It
changes the `deliver` guide, the Roundfix `deliver` reference, the Roundfix
skill version, the version record and the Behavior Surface Record that steps 1
to 3 also change. This Spec is rebased on the default branch after 0279
merges and checked again before its Run.

## Risks & Considerations

- A browser tool that writes a video or a trace under `qa/evidence/` now fails
  the gate at once. The `qa-gate` skill says where such files go.
- A Secret Rewrite that also reformats a line is refused, and the reason names
  the line. The operator replaces only the value.
- A blocker entry whose path contains `, ` is split wrongly and its path is
  never named, so such a rewrite is refused as touching an unnamed path.
- Steps 1 to 3 each raise the Roundfix skill version and re-record shared
  files, so they run in series.

## Glossary

- adds: **Secret Rewrite**
- changes: **Secret Scan**
- changes: **Mechanical Refusal Code**

## Decisions

- Content decides readability: no NUL byte, the PNG signature or the JPEG
  Start Of Image marker followed by a marker. See ADR-0291.
- No image size bound: the code has no natural limit for QA evidence. See
  ADR-0291.
- The Task scan reads added lines with the push scan's parser. See ADR-0292.
- A Secret Rewrite is compared line by line at the lines matched at the
  candidate and returns to `reviewing`. See ADR-0292.
- The prover is a separate optional dependency, like the review correction
  prover, so the existing `ItemHistory` fakes stay unchanged.
