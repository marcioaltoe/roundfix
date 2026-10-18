---
task: task_01
spec: 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
status: pending
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
