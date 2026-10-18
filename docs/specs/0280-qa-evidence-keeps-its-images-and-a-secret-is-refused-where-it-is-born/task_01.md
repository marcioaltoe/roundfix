---
task: task_01
spec: 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
status: completed
type: backend
complexity: high
---

# Task 01: QA evidence keeps its images, and the QA Report commit refuses any other binary

## Overview

On 2026-10-18 the Conexus QA gate committed PNG screenshots, as the `qa-gate`
skill asks, and the next gate's mechanical stage refused all 67 as built
binaries (`_prd.md` → Sources 1). The mechanical stage calls any blob with a
NUL byte binary, and the QA Report commit never looks at binaries.

This Task adds one content classifier, `ReadableEvidence`, that accepts text,
PNG and JPEG. The mechanical stage and the QA Report commit both use it, so
what the commit lands the next gate accepts, and any other binary fails the
gate that would commit it. It is verifiable on its own through mechanical and
Daemon tests on temporary repositories.

## Requirements

1. MUST answer the Finding of 2026-10-18,
   `references/2026-10-18-the-qa-gate-commits-screenshots-the-mechanical-audit-then-refuses.md`
   (`_prd.md` → Sources 1), by implementing API Contracts 1 to 3 of
   `_techspec.md` (ADR-0291).
2. MUST add `EvidenceProbeBytes` (8000) and `ReadableEvidence` to the new
   `internal/speccheck/evidence_image.go`, exactly as API Contract 1 states:
   - true for a probe with no NUL byte;
   - true for a probe that begins with `89 50 4E 47 0D 0A 1A 0A`;
   - true for a probe that begins with `FF D8 FF`;
   - false otherwise.

   It decides from content only and never reads a name.
3. MUST make `detectMechanicalEvidenceScratchState` raise `QA-EVIDENCE-PATH`
   for a regular blob only when `ReadableEvidence` of its first
   `EvidenceProbeBytes` bytes is false, keep the gitlink finding, and use the
   detail `tracked evidence path <path> is a non-image binary` and the fix
   text of API Contract 2. `mechanicalBlobIsBinary` MUST stop being the
   decision, and the probe size MUST come from `EvidenceProbeBytes`.
4. MUST refuse a non-image binary at the QA Report commit, as API Contract 3
   states.
   - When the Secret Scan of the prepared paths finds no match, the Daemon
     reads each prepared path that exists under the report's `evidence/`
     directory.
   - It settles the QA Task `failed` with the reason `QA evidence carries a
     non-image binary: <path>; <path>`.
   - It publishes one `daemon.qa` event with phase `evidence_binary` and the
     paths, and prints its summary.
   - It creates no QA Report commit.

   The data source is the Run Worktree's files the commit would stage. The
   event and the reason are read by the operator and the Run Database, and
   they carry paths only, never file content. The new reason MUST NOT start
   with `QA evidence carries a secret`, so the single-literal pin of
   `qaSecretRefusalDocumentation` holds.
5. MUST add these tests, each calling `t.Parallel()` first and compiling for
   Windows. Binary samples are byte strings built in the test; no test writes
   an executable file mode or runs a sample.
   - In the new `internal/speccheck/evidence_image_test.go`:
     - `TestReadableEvidenceDecidesByContent`;
     - `TestEvidenceAcceptsTrackedImages`: a tracked PNG, JPEG, SVG and an
       extensionless PNG raise no `QA-EVIDENCE-PATH`;
     - `TestEvidenceRefusesNonImageBinaries`: an ELF header, a SQLite header
       (`SQLite format 3` and a NUL), a PDF holding a NUL and an ELF header
       named `.png` each raise one finding naming the path and `non-image
       binary`.

     These reuse the mechanical fixtures of `mechanical_test.go`, which stays
     unchanged.
   - In the new `internal/daemon/qa_evidence_binary_test.go`:
     - `TestQAReportCommitsImageEvidence`: a PNG under the evidence directory
       rides in the QA Report commit;
     - `TestQAReportWithANonImageBinarySettlesFailedWithoutACommit`: the QA
       Task settles `failed` with the exact reason, the event names the path,
       and no QA Report commit exists.

     These reuse the QA fixtures of `qa_secret_scan_test.go`.

   `TestEvidenceRefusesScratchStateBinary` and every other existing test MUST
   keep its name and pass.
6. MUST document the change, each file naming `non-image binary`:
   - `docs/user-guide/commands/implement.md`, beside the QA Report commit
     Secret Scan paragraph;
   - `.agents/skills/roundfix/references/implement.md`, beside its QA Report
     commit Secret Scan paragraph;
   - `.agents/skills/qa-gate/SKILL.md`, under a new heading `### Evidence files
     are text or images`. It says that screenshots are PNG, JPEG or SVG, and
     that a video, a trace archive, a PDF, a database or a build output is
     recorded as text or kept outside `qa/evidence/`, because the Daemon
     refuses to commit any other binary there.

   MUST search the guides, the agent guides and the skill copies for text
   saying a binary in `qa/evidence/` is always refused or that screenshots
   are the only visual evidence, and qualify each. No text inside any `###
   QA settlement` section may change.
7. MUST revise **Mechanical Refusal Code** in `CONTEXT.md` through
   domain-modeling, with ADR-0291's link, so `QA-EVIDENCE-PATH` also names a
   tracked evidence file that is a gitlink or a non-image binary.
8. MUST run these commands, because this Task edits a user guide page, the
   Roundfix reference and the `qa-gate` skill. No record or mirror is
   hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to raise and record the Roundfix and `qa-gate` skill versions;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Add the classifier and use it in the mechanical stage.
- [ ] Add the QA Report commit refusal.
- [ ] Add the five tests.
- [ ] Document the change in the guide, the reference, the `qa-gate` skill and the glossary.
- [ ] Sync the mirrors, record the versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] Tracked PNG, JPEG and SVG evidence raises no `QA-EVIDENCE-PATH`, and every other binary raises one, whatever its name.
- [ ] A QA run commits its PNG evidence, and a QA run with a non-image binary settles `failed` without a QA Report commit.
- [ ] The skill, the guide, the reference and the glossary describe the rule.

## Consumers

- `speccheck.ReadableEvidence` and `speccheck.EvidenceProbeBytes` — consumers: `./internal/speccheck`, `./internal/daemon`

## Context

- instruction: `docs/adr/0291-qa-evidence-holds-text-and-images-and-the-qa-report-commit-refuses-any-other-binary.md`
- instruction: `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- instruction: `internal/speccheck/mechanical_test.go`
- instruction: `internal/daemon/qa_secret_scan_test.go`
- instruction: `internal/docscontract/command_documentation_test.go`
- creates: `internal/speccheck/evidence_image.go`
- creates: `internal/speccheck/evidence_image_test.go`
- creates: `internal/daemon/qa_evidence_binary_test.go`
- interface: `internal/speccheck/mechanical.go`
- interface: `internal/daemon/task_engine.go`
- interface: `docs/user-guide/commands/implement.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestReadableEvidenceDecidesByContent|TestEvidenceAcceptsTrackedImages|TestEvidenceRefusesNonImageBinaries)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReadableEvidenceDecidesByContent TestEvidenceAcceptsTrackedImages TestEvidenceRefusesNonImageBinaries; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestQAReportCommitsImageEvidence|TestQAReportWithANonImageBinarySettlesFailedWithoutACommit)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQAReportCommitsImageEvidence TestQAReportWithANonImageBinarySettlesFailedWithoutACommit; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/implement.md .agents/skills/roundfix/references/implement.md .agents/skills/qa-gate/SKILL.md CONTEXT.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'non-image binary' || { printf 'missing non-image binary in %s\n' "$f" >&2; exit 1; }; done; grep -qF -- '### Evidence files are text or images' .agents/skills/qa-gate/SKILL.md || { printf 'missing qa-gate heading\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Mechanical Refusal Code**:' || { printf 'missing glossary term Mechanical Refusal Code\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/speccheck || exit 1; go test -count=1 -run 'QA|Evidence|Secret' ./internal/daemon || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the five tests do not exist, so the command fails at its first pass check. Afterwards these checks pass: the mechanical and QA tests, the documentation, the glossary, the skill versions and the Behavior Surface Record.

## References

- `_prd.md` → Sources; Goals; User Story 1; User Story 2; Core Feature 1; Core Feature 2; Core Feature 5; Success Metric 1; Success Metric 2; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 9; Invariant 1; Current behavior; Testing Approach; Build Order 1
- ADR-0291; ADR-0273; ADR-0117; ADR-0125; ADR-0280; ADR-0189; ADR-0256; ADR-0259; ADR-0285


## Result

Implementation handoff for Task 01. Task status and declared Verification remain
Daemon-owned; no commit, push or pull request was made.

### Acceptance evidence

1. **Tracked text and image evidence.** Added `EvidenceProbeBytes = 8000` and
   the content-only `ReadableEvidence` classifier. The mechanical check now
   uses it and emits the exact non-image binary detail and fix from API
   Contract 2; the gitlink check remains. `TestReadableEvidenceDecidesByContent`,
   `TestEvidenceAcceptsTrackedImages` and `TestEvidenceRefusesNonImageBinaries`
   passed in the focused check below. Cases cover PNG, JPEG, SVG, an
   extensionless PNG, ELF, SQLite, a PDF containing NUL, an ELF named `.png`,
   malformed/offset image signatures, and NUL at and beyond the probe boundary.
   The existing `TestEvidenceRefusesScratchStateBinary`, gitlink, ordinary
   artifact and absent-directory tests also passed without changing
   `mechanical_test.go`.
2. **QA Report commit boundary.** After both existing Secret Scans find no
   match, the Daemon probes the existing prepared paths under the report's
   `evidence/` directory in the Run Worktree. A non-image binary sets the exact
   failed reason, publishes one `daemon.qa` event with phase `evidence_binary`
   and sorted paths, prints its summary, and returns without a report commit.
   `TestQAReportCommitsImageEvidence` passed with byte-for-byte PNG, JPEG, SVG
   and extensionless PNG contents in the actual temporary repository commit.
   `TestQAReportWithANonImageBinarySettlesFailedWithoutACommit` passed for one
   and two binaries, checking exact reason, event payload, summary, preserved
   report/evidence, no binary content in output, and no report commit or commit
   event. `TestQAReportBinaryCheckReadsOnlyPreparedEvidence` passed for nested
   images, deleted paths, unprepared and outside paths, sorted results, probe
   limits, read errors and cancellation. Existing QA Secret Scan and formatter
   tests passed without changing `qa_secret_scan_test.go`.
3. **Guidance and glossary.** The implement guide and Roundfix reference name
   the new refusal beside their Secret Scan paragraphs. The qa-gate skill adds
   `### Evidence files are text or images`, names PNG/JPEG/SVG screenshots and
   directs video, trace archives, PDFs, databases and build output to text or
   outside `qa/evidence/`; its focused-artifact instruction now specifies text
   or images. Through domain-modeling, `Mechanical Refusal Code` now includes
   tracked gitlinks and non-image binaries with ADR-0291's link. Searched
   `docs/user-guide`, `docs/agents`, `.agents/skills` and `skills` for old binary
   refusal and screenshot-only wording; no additional contradictory guidance
   was found. A read-only Python audit confirmed canonical/mirror equality and
   byte-identical `### QA settlement` sections for qa-gate, Roundfix and
   archive-spec against HEAD. The existing single-literal Secret Scan
   documentation checks passed.

### Focused checks

- Before production edits,
  `go test ./internal/speccheck -run '^TestEvidence(AcceptsTrackedImages|RefusesNonImageBinaries)$' -count=1`
  exited 1: tracked PNG/JPEG were refused as built binaries and the refusal
  wording was still the old text.
- Before production edits,
  `go test ./internal/daemon -run '^TestQAReport(CommitsImageEvidence|WithANonImageBinarySettlesFailedWithoutACommit)$' -count=1`
  exited 1: a QA run with non-image evidence was accepted.
- After all Go edits,
  `go test ./internal/speccheck ./internal/daemon -run '^Test(ReadableEvidence|Evidence(Accepts|Refuses)|QA(Report(CommitsImageEvidence|WithANonImageBinarySettlesFailedWithoutACommit|BinaryCheck)|ReportWithASecretMatch|ReportWithoutASecretMatch|SecretScan|Formatter))' -count=1`
  exited 0 for both packages.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheQASecretRefusal(IsDocumentedAsTheDaemonPrintsIt|ContractRejectsTheParenthesizedForm)$' -count=1`
  exited 0.
- `GOOS=windows GOARCH=amd64 go test -c -o /private/tmp/roundfix-0280-26a9aa2dcaf8e01d-speccheck.test.exe ./internal/speccheck`
  and `GOOS=windows GOARCH=amd64 go test -c -o /private/tmp/roundfix-0280-26a9aa2dcaf8e01d-daemon.test.exe ./internal/daemon`
  each exited 0. Every new top-level test calls `t.Parallel()` first; binary
  samples are runtime byte strings written as ordinary files and never run.

### Generated artifacts and scope

- `make skills-sync` exited 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0 and generated qa-gate `0.0.17` and Roundfix `0.1.99` in both copies
  and the owned-skill version record.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  exited 0 and regenerated the Behavior Surface Record.
- Repository-required `make baseline-digests` exited 0 with `changed:false`;
  no derived Baseline artifact remains changed.
- The initial read-only scope audit found the 16 declared Task 01 paths,
  including this Result file. The Verification Feedback repair below adds one
  ordinary fixture path; the Task Graph and other Task files remain unchanged. All
  mirrors and records were written only by their owning commands.
- The configured sandbox cannot start because its writable root contains a
  symlink. Commands ran through the approved escalation path. Checks used
  temporary repositories and made no network or provider calls.

The Task's declared Verification, selected repository Verification and terminal
settlement were not run by this Agent. No follow-up implementation was added.


### Verification Feedback repair — attempt 1

The Daemon diagnostic artifact for attempt 1 identified regressions in
`TestQAGateCarriesARowFromAnUnintegratedFailedPass` and
`TestQAFormatFormatsAnImportedPassBeforeThePrecondition`. Both reused
`recordPriorCommit`, whose evidence fixture was a non-image binary. The new
commit refusal correctly left that evidence outside HEAD, so their later
imported-file commit assertions could not read it.

Repaired the shared fixture in `internal/daemon/qa_prior_pass_test.go` to use
PNG content, retaining the `.bin` name and NUL, CRLF and high-byte payload.
Its existing-file preservation check reads the same fixture constant. This
keeps the original byte-preservation, file-identity, mode, import, carry and
format assertions, all existing test names, and the new non-image refusal.
This additional ordinary test path is required by Task 01's rule that existing
tests continue to pass under the changed evidence contract; no other Task's
implementation or protected tooling path changed.

Focused evidence:

- Before the fixture repair,
  `go test ./internal/daemon -run '^Test(QAGateCarriesARowFromAnUnintegratedFailedPass|QAFormatFormatsAnImportedPassBeforeThePrecondition)$' -count=1`
  exited 1, reproducing both reported regressions.
- After the repair,
  `go test ./internal/daemon -run '^Test(QAGateCarriesARowFromAnUnintegratedFailedPass|QAFormat|PriorQAPass|APriorPass|QAPromptNamesTheImportedPassHead|QAReport(CommitsImageEvidence|WithANonImageBinarySettlesFailedWithoutACommit|BinaryCheck))' -count=1`
  exited 0, covering both reported tests, the shared fixture's other prior-pass
  consumers, QA formatting and the new evidence acceptance/refusal cases.
- `GOOS=windows GOARCH=amd64 go test -c -o /private/tmp/roundfix-0280-26a9aa2dcaf8e01d-daemon.test.exe ./internal/daemon`
  exited 0 after the fixture repair.

Task status remains Daemon-owned and unchanged. The Agent did not rerun the
Task's declared Verification, commit, push or open a pull request.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/daemon/qa_prior_pass_test.go`
