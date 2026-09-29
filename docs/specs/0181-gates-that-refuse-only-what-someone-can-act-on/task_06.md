---
task: task_06
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: completed
type: backend
complexity: low
---

# Task 06: Citation checks read only what the Spec's authors wrote

## Overview

`readSpecCitations` in `internal/speccheck/citations.go` walks every Markdown file under a Spec folder, and `SC-ADR-UNLISTED` then requires the PRD to list every ADR it finds. That walk includes an Agent's `## Result`, the Daemon-owned `## Recorded paths` and `## Carry-forward provenance` sections, and the QA reports and evidence under `qa/`. None of that text is written by the Spec's author, and the author is the only one who can list an ADR. On 2026-09-29 this Spec's own QA gate refused before measuring any row. Task 03's Result named ADR-0002 as a fixture number, and the refused report then quoted the finding. This corrective Task makes the walk read only the authored projection (ADR-0176).

## Requirements

1. MUST make `readSpecCitations` skip every file under the Spec folder's `qa/` directory.
2. MUST make `readSpecCitations` ignore, in each `task_*.md` file, the lines of a `## Result`, `## Recorded paths` or `## Carry-forward provenance` section, from its heading up to the next line that starts with `## `. Every other line keeps its line number in the reported location.
3. MUST keep every other file under the Spec folder read exactly as before, including `_prd.md`, `_techspec.md`, `_tasks.md`, `_authorization.md`, `_idea.md` and `references/`. A citation in authored text MUST still report `SC-ADR-UNLISTED` with its unchanged summary and fix text.
4. MUST leave `SC-ADR-RELATED`, `SC-CITATION-UNSUPPORTED`, the horizon from task_03 and the corpus golden unchanged. It MUST rename or remove no top-level test and change no exported function signature.
5. MUST put the new tests in `internal/speccheck/citation_projection_test.go`, over plain temporary directories.
6. MUST NOT name any decision by its `ADR-` identifier in this Task's `## Result`. This Spec's QA gate still runs on the v0.20.0 auditor, which reads that section. Describe fixtures by their role instead.

## Subtasks

- [ ] Skip the `qa/` directory in the Spec citation walk.
- [ ] Skip the Agent- and Daemon-owned sections of Task files.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] An unlisted ADR cited only in a Task's `## Result` opens no `SC-ADR-UNLISTED`.
- [ ] An unlisted ADR cited only in `## Recorded paths` or `## Carry-forward provenance` opens none.
- [ ] An unlisted ADR cited only in a QA report or QA evidence file under `qa/` opens none.
- [ ] An unlisted ADR cited in a Task's Requirements, in `_techspec.md` or in `references/` still opens `SC-ADR-UNLISTED`, at its original line.
- [ ] A Task file whose authored section follows its `## Result` still has that authored citation reported.

## Context

- instruction: `docs/adr/0176-citation-checks-read-only-what-a-specs-authors-wrote.md`
- interface: `internal/speccheck/citations.go`
- creates: `internal/speccheck/citation_projection_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestASpecCitationInATaskResultIsNotAnObligation|TestASpecCitationInADaemonSectionIsNotAnObligation|TestASpecCitationInAQAReportIsNotAnObligation|TestAnAuthoredSpecCitationStillMustBeListed|TestAnAuthoredSectionAfterAResultIsStillRead|TestCheckADRClosureDepthOne|TestTheHorizonLeavesUnlistedCitationsChecked)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestASpecCitationInATaskResultIsNotAnObligation TestASpecCitationInADaemonSectionIsNotAnObligation TestASpecCitationInAQAReportIsNotAnObligation TestAnAuthoredSpecCitationStillMustBeListed TestAnAuthoredSectionAfterAResultIsStillRead TestCheckADRClosureDepthOne TestTheHorizonLeavesUnlistedCitationsChecked; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the five new named tests exists, so the command fails.

## References

- `_prd.md` → Core Feature 4
- `_techspec.md` → Citation checks read the authored projection; Testing Approach 6
- ADR-0176

## Result

Implementation:

- The Spec citation walk now prunes the direct `qa/` subtree before reading files.
- Task Markdown is projected line by line: Agent- and Daemon-owned sections are ignored through the next second-level heading, while authored lines retain their original indexes.
- Five top-level regression tests use plain temporary repositories and separate the Result, each Daemon section, QA report/evidence, authored-file, and post-Result authored-section cases.

Focused checks:

- `rtk env GOCACHE=/tmp/roundfix-0181-task06-gocache go test -count=1 ./internal/speccheck -run 'TestASpecCitationInATaskResultIsNotAnObligation|TestASpecCitationInADaemonSectionIsNotAnObligation|TestASpecCitationInAQAReportIsNotAnObligation|TestAnAuthoredSpecCitationStillMustBeListed|TestAnAuthoredSectionAfterAResultIsStillRead'` — failed before the source change because every excluded fixture opened an unlisted-decision finding and the post-Result fixture reported the excluded line; passed after the source change.
- `rtk git diff --check` — passed after the implementation and tests.
- `rtk make verify-incremental` — the restricted run could not inspect the host process table and also observed one transient store timing failure; the elevated rerun cleared both and passed `internal/speccheck`, but remained non-zero because `skills/TestSettlementGuidanceIsOneTable` found the committed canonical QA-settlement text differs from `.agents/skills/archive-spec/SKILL.md`.

Acceptance evidence:

- `TestASpecCitationInATaskResultIsNotAnObligation` proves a decision mentioned only in Result evidence opens no finding.
- `TestASpecCitationInADaemonSectionIsNotAnObligation` separately covers recorded paths and carry-forward provenance; neither opens a finding.
- `TestASpecCitationInAQAReportIsNotAnObligation` separately covers the report and nested evidence paths under `qa/`; neither opens a finding.
- `TestAnAuthoredSpecCitationStillMustBeListed` covers Task Requirements, the TechSpec, and an adopted reference, asserting the exact existing summary, fix, path, and original line.
- `TestAnAuthoredSectionAfterAResultIsStillRead` proves scanning resumes at the next authored section and reports line 9 rather than the excluded line 5.

Not run:

- The Task's declared Verification command is reserved for the Daemon.

Follow-up:

- Reconcile the archive skill's QA-settlement table with its canonical source in the owning Task. That skill is outside this corrective slice.

## Carry-forward provenance

- Source Run: `run_20260929T185005Z_5bdc71810126cc68`
- Source commit: `d7a1df1204bcffda87b888939d8ef51545a93537`
