---
task: task_01
spec: 0340-qa-evidence-links-resolve-in-the-committed-tree
status: pending
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
