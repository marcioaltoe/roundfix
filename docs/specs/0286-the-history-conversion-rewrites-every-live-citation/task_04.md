---
task: task_04
spec: 0286-the-history-conversion-rewrites-every-live-citation
status: pending
type: backend
complexity: high
---

# Task 04: The check reports a removed history citation, and the guides describe the conversion

## Overview

After the adopter's conversion, `roundfix spec check --strict` exited 0 with
18 citations of removed files still in place (`_prd.md` → Sources 1). This
Task extends `SC-HISTORY-RETAINED` to report each citation of a Removed
History Path that is not a Git Locator, with the locator to use. It also
documents the whole Spec: the four conversion corrections, the check and the
`--no-history` interaction, in the guides, the Roundfix skill and its
references, and the glossary.

## Requirements

1. MUST implement API Contract 6 of `_techspec.md` (ADR-0299) in
   `internal/speccheck/history_retained.go`, reusing `CitationSources`,
   `SourceHistoryCitations`' detector, `RemovedHistoryPaths` and
   `PendingTaskCreates` from Tasks 01 and 03 (`internal/spec/citation_sources.go`
   and `internal/spec/pending_targets.go`). The finding is an error with the
   summary `<path> cites a removed history file without a Git Locator:
   <target>`, the location `<path>:<line>` and the fix `Replace the citation
   with <locator>.`.
2. MUST state the hazard in the code's behavior, as API Contract 6 does:
   - the detector reads only the repository's tracked and untracked,
     unignored files and its Git history, and sends nothing anywhere;
   - files under `docs/history/` are exempt, Archive Records included, as are
     a pending Task's `creates:` target and a line already reported as a
     missing relative link to the same target;
   - text already written as a Git Locator, or followed by `(not reachable
     from any branch or tag)` or `(pull request #`, is never reported;
   - without the retirement rule the detector stays skipped, and outside Git
     the finding is skipped under the existing reason.
3. MUST keep the `SC-HISTORY-RETAINED` count of the characterized corpus at 0
   without changing `internal/docscontract/testdata/corpus-golden.json`, and
   keep every existing `SC-HISTORY-RETAINED` finding and its wording.
4. MUST add in `internal/speccheck/history_retained_test.go`
   `TestHistoryRetainedReportsARemovedHistoryCitationWithoutALocator`, its
   top-level test calling `t.Parallel()` first and compiling for Windows. On a
   temporary repository of the shape an earlier conversion left, with the
   retirement rule in the guide, the five citations of Surface Transcript 1's
   fixture (the Makefile comment, the ADR, the agent guide, the reference page
   and the README, the history files deleted by a commit) are each one error
   with that exact summary, line and fix. A Git Locator, an Archive Record
   under `docs/history/specs`, a skill tree file, a `testdata` file and a
   pending Task's `creates:` target are not reported. Without the retirement
   rule no such finding exists.
5. MUST document the change. These sentences are written verbatim, with
   inline code formatting allowed around command names, flags and paths:
   - in `docs/user-guide/commands/baseline.md`, replacing the paragraph that
     says other citations are listed and left unchanged: "After removals the
     update rewrites every History Citation it lists, in every Citation
     Source: each tracked text file outside the skill trees and testdata
     directories. A prose path that starts with ./ or ../ is resolved against
     the directory of the citing file. Outside Markdown the Git Locator is
     written without backticks.";
   - in `docs/user-guide/commands/baseline.md`,
     `docs/user-guide/commands/history.md`,
     `.agents/skills/roundfix/references/baseline.md` and
     `.agents/skills/roundfix/references/archive.md`, replacing the sentence
     that dates the lineage line with HEAD's committer date: "A retired ADR
     whose file name is not a number followed by a slug is named by its file
     name without the extension. The date is the day the retired ADR was
     superseded: its updated_at in UTC, else the committer date of the newest
     commit that changed it, else the committer date of HEAD.";
   - in the same four files, in the `live-links` paragraph: "The repair
     leaves a link to a file that a Task of an active Spec declares in
     creates: and has not completed. It resolves a missing link by its file
     name when exactly one removed history file has that name.";
   - in `docs/user-guide/commands/history.md`,
     `.agents/skills/roundfix/references/baseline.md` and
     `.agents/skills/roundfix/references/archive.md`, the first sentence of
     the first bullet of this list, with that page's own subject in place of
     "the update" (for example "After removals the plan rewrites every History
     Citation it lists, in every Citation Source");
   - in `docs/user-guide/commands/baseline.md` and
     `.agents/skills/roundfix/references/baseline.md`, beside
     `--no-history`: "While Pending History remains, an update with
     --no-history still writes the retirement rule into the guide, so
     SC-HISTORY-RETAINED reports every retained history file on each active
     Spec, and the path to green is the history conversion: rerun roundfix
     baseline update without --no-history, or apply roundfix history sanitize
     batches until nothing is pending.";
   - in `docs/user-guide/commands/spec.md` and
     `.agents/skills/roundfix/references/spec.md`, in the
     `SC-HISTORY-RETAINED` paragraph: "It also reports a citation that cites a
     removed history file without a Git Locator in any Citation Source outside
     docs/history, and skips a link to a file a pending Task declares in
     creates:.";
   - in `.agents/skills/roundfix/SKILL.md`, replacing "never rewrites citing
     files": the update rewrites every listed citation to a Git Locator.
6. MUST search the guides, the agent guides and the skill copies for the old
   wording and change each occurrence: "Other citations are listed and left
   unchanged", "other citations are listed and left", "is HEAD's committer
   date in UTC" and "never rewrites citing files". No text inside any `###
   QA settlement` section may change.
7. MUST write these terms into `CONTEXT.md` through domain-modeling, each
   linking `docs/adr/0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md`:
   - add **Citation Source**: a tracked regular text file outside the skill
     trees and testdata directories, holding no NUL byte in its first 8,000
     bytes, in which History Citations are found;
   - add **Removed History Path**: a path under the History Root's Findings,
     Backlog or ADR family that HEAD or its history holds and the working tree
     lacks;
   - revise **History Citation**: it is found and rewritten in every Citation
     Source;
   - revise **Git Locator**: outside Markdown it is written as plain text
     without backticks.
8. MUST run these commands, because this Task edits user guide pages and the
   Roundfix skill. No record or mirror is hand-edited. Before them, raise
   both version fields of the Roundfix skill's canonical `SKILL.md`
   (`metadata.version` and the top-level `version`) by one patch over the
   value the tree holds when this Task starts, because the record command
   refuses changed content under an unchanged version. The Agent runs these
   commands itself before handing back: the Daemon only runs the Verification
   and never regenerates a mirror or a record.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Add the removed-citation finding with its exemptions.
- [ ] Add its test on the shape an earlier conversion left.
- [ ] Document the corrections, the check and the `--no-history` note in the guides and references.
- [ ] Write the glossary terms, raise the skill version, sync and record.

## Acceptance Criteria

- [ ] `spec check --strict` fails on the shape an earlier conversion left, with one error per citation naming its file, line and locator.
- [ ] Files under `docs/history/`, Git Locators and pending `creates:` targets are never reported.
- [ ] The guides, references, skill and glossary describe the four corrections, the check and the path to green under `--no-history`.

## Context

- instruction: `docs/adr/0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md`
- instruction: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/speccheck/history_retained.go`
- interface: `internal/speccheck/history_retained_test.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/commands/history.md`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- the `SC-HISTORY-RETAINED` findings `roundfix spec check` reports — consumers: `./internal/speccheck`, `./internal/docscontract`

## Verification

- `grep -qF -- 'cites a removed history file without a Git Locator' internal/speccheck/history_retained.go || { printf 'missing removed-citation finding\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^TestHistoryRetainedReportsARemovedHistoryCitationWithoutALocator$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestHistoryRetainedReportsARemovedHistoryCitationWithoutALocator ' || { printf 'missing pass: TestHistoryRetainedReportsARemovedHistoryCitationWithoutALocator\n' >&2; exit 1; }; for phrase in 'rewrites every History Citation it lists, in every Citation Source' 'is resolved against the directory of the citing file' 'Outside Markdown the Git Locator is written without backticks' 'the path to green is the history conversion'; do tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- "$phrase" || { printf 'missing %s in the baseline guide\n' "$phrase" >&2; exit 1; }; done; for file in docs/user-guide/commands/baseline.md docs/user-guide/commands/history.md .agents/skills/roundfix/references/baseline.md .agents/skills/roundfix/references/archive.md; do for phrase in 'is named by its file name without the extension' 'The date is the day the retired ADR was superseded' 'has not completed. It resolves a missing link by its file name when exactly one removed history file has that name'; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing %s in %s\n' "$phrase" "$file" >&2; exit 1; }; done; ! tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'committer date in UTC (' || { printf 'old lineage date wording in %s\n' "$file" >&2; exit 1; }; done; for file in docs/user-guide/commands/history.md .agents/skills/roundfix/references/baseline.md .agents/skills/roundfix/references/archive.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'every History Citation it lists, in every Citation Source' || { printf 'missing rewrite rule in %s\n' "$file" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- 'the path to green is the history conversion' || { printf 'missing no-history note in the baseline reference\n' >&2; exit 1; }; for file in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'cites a removed history file without a Git Locator in any Citation Source' || { printf 'missing check rule in %s\n' "$file" >&2; exit 1; }; done; for file in docs/user-guide/commands/baseline.md docs/user-guide/commands/history.md .agents/skills/roundfix/references/archive.md .agents/skills/roundfix/references/baseline.md .agents/skills/roundfix/SKILL.md; do ! tr -s '[:space:]' ' ' < "$file" | grep -qiF -- 'other citations are listed and left' || { printf 'old citation scope in %s\n' "$file" >&2; exit 1; }; ! tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'never rewrites citing files' || { printf 'old rewrite wording in %s\n' "$file" >&2; exit 1; }; done; for term in '**Citation Source**:' '**Removed History Path**:' '**History Citation**:' '**Git Locator**:'; do tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "$term" || { printf 'missing glossary term %s\n' "$term" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'it is found and rewritten in every Citation Source' || { printf 'History Citation not revised\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'outside Markdown it is written as plain text without backticks' || { printf 'Git Locator not revised\n' >&2; exit 1; }; grep -qF -- '0299-the-history-conversion-rewrites-every-citation-it-lists' CONTEXT.md || { printf 'glossary does not link ADR-0299\n' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/speccheck || exit 1; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/speccheck || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the finding does not exist, so the command fails at its first check. Afterwards the new check test passes, the guides, references, skill and glossary carry the new rules and none of the old wording, `internal/speccheck` builds for Windows and passes whole, and the skill version, the mirrors, the Behavior Surface Record and the documentation contracts are current.

## References

- `_prd.md` → Sources; Goals; Core Feature 4; Core Feature 5; Success Metric 4; Glossary
- `_techspec.md` → API Contract 6; API Contract 7; Invariant 6; Current behavior; Testing Approach; Build Order 4; Glossary
- ADR-0299; ADR-0278; ADR-0282; ADR-0260; ADR-0189; ADR-0256; ADR-0280; ADR-0259; ADR-0285; ADR-0265
