---
task: task_01
spec: 0340-qa-evidence-links-resolve-in-the-committed-tree
status: completed
type: backend
complexity: medium
---

# Task 01: The shared evidence link check refuses unresolved links in qa-report accept and settle

## Overview

`internal/spec` gains the one check every QA settlement path reads: the
Markdown link extraction of `MissingRelativeLinks` becomes the exported
`RelativeLinks`, shared with `SC-HISTORY-RETAINED`, and
`UnresolvedEvidenceLinks` classifies each local link of a QA Report as
resolved, `ignored` or `missing` by asking Git what the commit would hold.
`roundfix qa-report accept` and the Settle Command for a `qa` Task refuse an
Unresolved Evidence Link with one refusal text. The guides, the Roundfix
settle reference and `CONTEXT.md` state the rule. This answers the adopted
Argus report of 2026-11-01, "QA evidence links resolve in the committed tree"
([references](references/2026-11-01-qa-evidence-links-resolve-in-the-committed-tree.md)).

## Requirements

1. MUST add `RelativeLink` and `RelativeLinks(path string, content []byte)
   []RelativeLink` to `internal/spec/record_links.go` and rewrite
   `MissingRelativeLinks` as a filter over it, keeping its signature and its
   results byte for byte (`_techspec.md` → API Contract 1).
2. MUST create `internal/spec/qa_evidence_links.go` with
   `EvidenceLinkIgnored`, `EvidenceLinkMissing`, `UnresolvedEvidenceLink`,
   `UnresolvedEvidenceLinks`, `EvidenceLinksError` and
   `RequireResolvedEvidenceLinks` as `_techspec.md` → Interfaces and API
   Contracts 2 and 3 state: one `git rev-parse --show-toplevel` and one
   `git ls-files -z --cached --others --exclude-standard` call with literal
   pathspecs through the package's existing Git reader, a disk-only rule
   outside Git, document order, and the refusal text
   `QA Report evidence links do not resolve in the committed tree: <destination> (line <n>, <reason>)`
   joined by `; `. A read or Git failure is returned as an error, never as a
   link.
3. MUST create `internal/spec/qa_evidence_links_test.go` with
   `TestUnresolvedEvidenceLinksNameTheReason`,
   `TestUnresolvedEvidenceLinksOutsideGitReadTheDisk`,
   `TestEvidenceLinksErrorNamesEachLink` and
   `TestMissingRelativeLinksReadsTheSharedExtraction`, each calling
   `t.Parallel()` first, with the cases `_techspec.md` → Testing Approach
   lists, the forced-add `evidence/c.log` and the code-span link included.
4. MUST make `runQAReportCommand` (`internal/cli/qa_report.go`) call
   `spec.RequireResolvedEvidenceLinks` after `QAReportEligibility` accepts
   and print a refusal as `roundfix: QA Report eligibility refused: <error>`
   with exit 1 (API Contract 4). Leave `qaReportUsage` unchanged.
5. MUST make `requireSettleQAReportEligibility` (`internal/cli/settle.go`)
   call it on the same report after `QAReportEligibility` accepts and wrap a
   refusal as `<error> (report <path>)` (API Contract 5). Change no other
   function of that file.
6. MUST add `TestQAReportAcceptRefusesAnUnresolvedEvidenceLink` and
   `TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink` to
   `internal/cli/qa_report_test.go`, each calling `t.Parallel()` first, with
   the exact stderr lines `_prd.md` → Success Metrics 1 and 3 give and the
   accepted rerun after the evidence is saved as `.txt`.
7. MUST document the rule in `docs/user-guide/commands/qa-report.md` and in
   the `## QA Report acceptance` section of the canonical
   `.agents/skills/roundfix/references/settle.md`: what a local link is, the
   two reasons, the refusal text, and the sentence "Every local link of the
   report must resolve in the committed tree" verbatim. In
   `docs/user-guide/commands/settle.md`, after the paragraph on QA eligibility,
   state that the Settle Command applies the same rule to a `qa` Task with its
   stderr line. Change no `### QA settlement` section.
8. MUST add **Unresolved Evidence Link** to `CONTEXT.md` through
   `domain-modeling`, beginning "A local link of a QA Report", with both
   reasons and an `_Avoid_` line, and append to **QA Report** the sentence
   "A QA Report settles only when every local link resolves in the committed
   tree." citing ADR-0351.
9. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`, recording
   `skills/testdata/owned-skill-versions.json`. Never hand-edit a mirror or a
   version field.
10. MUST re-record the Behavior Surface Record as the LAST step, after the
    last edit of any guide, with `go test -count=1 -tags docscontract
    ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'
    -record-skill-coverage`, and commit `docs/references/behavior-surfaces.json`.
11. MUST NOT touch `spec.QAReportEligibility`, `internal/spec/archive.go`,
    `internal/daemon/`, the `qa-gate` skill, `docs/references/skill-coverage.json`,
    any Baseline module, `docs/agents/`, the Makefile, CI workflows, `go.mod`
    or `go.sum`.

## Subtasks

- [ ] Characterize: an eligible report linking an ignored `.log` and an absent `.txt` is accepted by `qa-report accept` today.
- [ ] Extract `RelativeLinks`, rewrite `MissingRelativeLinks` over it and add the classification with its tests.
- [ ] Wire `qa-report accept` and the Settle Command, with their tests.
- [ ] Update the guides, the reference and `CONTEXT.md`, sync the mirror, record the skill version, re-record the fingerprints last.

## Acceptance Criteria

- [ ] `qa-report accept` and the Settle Command refuse an ignored or missing link with the shared text and accept the same report once the evidence is kept as `.txt`.
- [ ] A tracked target resolves whatever the ignore patterns say, and `SC-HISTORY-RETAINED` reads the same links through the shared extraction.
- [ ] The guides, both reference copies and `CONTEXT.md` carry the rule, the mirrors equal their canonical copies, and the owned-skill version and Behavior Surface records pass.

## Context

- instruction: `docs/adr/0351-a-qa-report-settles-only-when-its-local-links-resolve-in-the-committed-tree.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- creates: `internal/spec/qa_evidence_links.go`
- creates: `internal/spec/qa_evidence_links_test.go`
- interface: `internal/spec/record_links.go`
- interface: `internal/cli/qa_report.go`
- interface: `internal/cli/settle.go`
- interface: `internal/cli/qa_report_test.go`
- interface: `docs/user-guide/commands/qa-report.md`
- interface: `docs/user-guide/commands/settle.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- `spec.RelativeLinks` and `spec.MissingRelativeLinks` — consumers: `./internal/spec`, `./internal/speccheck`
- `spec.UnresolvedEvidenceLinks`, `spec.RequireResolvedEvidenceLinks` and `spec.EvidenceLinksError` — consumers: `./internal/cli`
- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestUnresolvedEvidenceLinksNameTheReason|TestUnresolvedEvidenceLinksOutsideGitReadTheDisk|TestEvidenceLinksErrorNamesEachLink|TestMissingRelativeLinksReadsTheSharedExtraction)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestUnresolvedEvidenceLinksNameTheReason TestUnresolvedEvidenceLinksOutsideGitReadTheDisk TestEvidenceLinksErrorNamesEachLink TestMissingRelativeLinksReadsTheSharedExtraction; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^(TestQAReportAcceptRefusesAnUnresolvedEvidenceLink|TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQAReportAcceptRefusesAnUnresolvedEvidenceLink TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/spec ./internal/speccheck && go test -count=1 -timeout 20m -run 'QAReport|Settle|QAPartial|PartialShape' ./internal/cli` — expected: exit 0. Before this Task none of the six tests exists, so the first pass line is missing; afterwards they pass with the Spec, Spec Check, QA Report and Settle suites.
- `grep -qF 'RequireResolvedEvidenceLinks(' internal/cli/qa_report.go || { echo 'qa-report accept does not read the shared check' >&2; exit 1; }; grep -qF 'RequireResolvedEvidenceLinks(' internal/cli/settle.go || { echo 'settle does not read the shared check' >&2; exit 1; }; for f in docs/user-guide/commands/qa-report.md .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Every local link of the report must resolve in the committed tree' || { printf 'missing rule in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/settle.md | grep -qF -- 'QA Report evidence links do not resolve in the committed tree' || { echo 'missing settle refusal in the settle guide' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Unresolved Evidence Link**:' || { echo 'missing term in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'A QA Report settles only when every local link resolves in the committed tree.' || { echo 'missing QA Report sentence in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**QA Report**:' || { echo 'missing QA Report term in CONTEXT.md' >&2; exit 1; }; cmp .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither CLI path reads the check, so the first test fails. Afterwards both read it, the guides, both reference copies and the glossary carry the rule, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 1, 2 and 4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 6; Success Metric 7; Glossary
- `_techspec.md` → Current behavior 3-5; API Contracts 1-5 and 8; Invariants 1-3 and 5; Declared breaks; Testing Approach; Build Order 1


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, no-workarounds, systematic-debugging, evidence-gate.

Implemented the shared relative-link extraction and committed-tree check. Both
CLI paths apply it after the existing QA eligibility decision. Git reads use
the existing reader, one root lookup and one batched literal-pathspec listing
when there are existing repository targets. Disk, Git and cancellation errors
remain infrastructure errors. The qa-report usage and all other Settle
functions retain their bytes.

Acceptance evidence:

- Refusal and repair: the two new CLI tests initially failed because both
  commands returned success with ignored evidence; qa-report also accepted a
  missing target, and Settle omitted the ignored evidence from its fixture
  commit. After implementation, both assert the exact stderr lines. Saving
  evidence as `.txt` makes qa-report accept and Settle succeed; the Settle
  fixture's index includes the saved evidence. The refusal leaves the fixture
  Task file unchanged.
- Shared extraction and resolution: the four new spec tests exercise document
  order and line numbers, a committed forced-add `evidence/c.log`, an addable
  `.txt`, mixed and ignored-only directories, an existing outside-repository
  target, self links, code spans, URLs, anchors, absolute destinations,
  reference definitions, literal glob characters, outside-Git disk behavior,
  read failure, cancellation and corrupt Git metadata/index failures. The
  history filter also retains its pre-existing broken/self-link behavior.
- Documentation and owned skills: qa-report and settle guides, the canonical
  QA Report acceptance reference and its generated mirror state the rule,
  reasons and refusal text. CONTEXT defines Unresolved Evidence Link with an
  Avoid line and extends QA Report with ADR-0351. `make skills-sync` and
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`
  exited 0; the recorder raised both canonical and mirrored fields from
  `0.1.167` to `0.1.168` and recorded that digest. Exact byte comparisons of
  both mirror pairs passed. The Behavior Surface Record is regenerated by the
  final command below, after all guide edits; its outcome is in the handoff.

Focused checks, all exit 0 on their final run:

- `go test ./internal/spec ./internal/cli -run 'Test(UnresolvedEvidenceLinks|EvidenceLinksError|MissingRelativeLinksReads|QAReportAcceptRefusesAnUnresolved|SettleRefusesAQAReportWithAnIgnored)' -count=1`
- `go test ./internal/spec ./internal/speccheck ./internal/cli -run 'Test(MissingRelativeLinks|PlanRecordLinks|HistoryRetained|RunQAReport|QAReportAccept|SettleAppliesEligibility|SettleQATaskRefuses|SettleRefusesAQAReport)' -count=1`
  passed spec and speccheck, and reached CLI PASS, but suiteguard rejected the
  CLI process because skill sync overlapped it. After the writer ended,
  `go test ./internal/cli -run 'Test(RunQAReport|QAReportAccept|SettleAppliesEligibility|SettleQATaskRefuses|SettleRefusesAQAReport)' -count=1`
  exited 0 with the tree held stable. This was an execution-order error,
  not a relaxed guard or assertion.
- `go test ./skills -run 'Test(RoundfixSkillIndex|InstallWritesEveryRoundfixReference|SettlementGuidanceIsOneTable|EveryOwnedSkillVersionIsRecorded)' -count=1`
- Python byte comparisons and changed-path allowlist inspection; `git diff --check`.

Final sanctioned regeneration, after the last guide edit:

`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`

The initial worktree change was the Daemon-owned `pending` to `in_progress`
status write in this Task. That field is preserved. Declared Verification and
repository Verification remain Daemon-owned and were not run. No other Task,
Task Graph, Daemon implementation, archive or eligibility policy was changed.
No commit, push or Pull Request was made. No follow-up outside this slice was
implemented.
