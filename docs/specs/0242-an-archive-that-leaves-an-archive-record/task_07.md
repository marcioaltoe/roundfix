---
task: task_07
spec: 0242-an-archive-that-leaves-an-archive-record
status: pending
type: backend
complexity: medium
---

# Task 07: Absorbed history resolves without the Spec folders, and the archive guidance describes only the record

## Overview

The first QA gate of this Spec (`qa/qa-report-2026-10-06.md` in Run
`run_20261006T232220Z_7c2e4a372c1dd72e`) failed rows 09 and 14.

Finding F-04 (Blocks-Completion): in a `git clone --no-local` of the audited
head with all 223 folders under `docs/history/specs` removed in a commit,
`make verify` exits 0 but `make verify-docs` exits 2.
`TestCheckCorpusGolden` and `TestCheckActiveCorpusHasNoErrors` report 97
`SC-ARCHIVE-LICENSE` errors. Each names a Finding or Rollup under
`docs/history/findings` whose `absorbed_by` names an archived Spec that then
exists neither as a folder nor as an Archive Record. The record-aware
resolver in `internal/speccheck/citations.go` accepts a `<slug>.md` record,
but `detectFindingsConsistency` has no answer for a Spec whose folder left
the tree and lives only in Git. Core Feature 4 of `_prd.md` lets the Spec
check's absorption check read the pre-archive tree from Git, and Success
Metric 4 requires both gates to pass in that clone.

Finding F-03 (Trust-Damage): the qualifying declared `partial` row of the
three `### QA settlement` tables still says the archive keeps "The Spec, its
QA report and evidence, and the declarations' `satisfied-by` record",
instead of the Exact text. The archive-spec skill's description, purpose,
Steps and Unarchive sections, and the Roundfix skill's
`references/archive.md`, still describe stamping `_prd.md`, moving the folder
and rewriting links, ahead of an appended Archive Record section that
contradicts them. `docs/user-guide/commands/archive.md` lacks the Vocabulary
Contract pattern `archive plan for`, which `spec check` reports as
`SC-VOCABULARY-UNDOCUMENTED`.

This Task makes the absorption check resolve through the record or through
Git, rewrites the stale guidance, and proves the ablation in a disposable
clone. It rewrites no file under `docs/history`.

## Requirements

1. MUST make `SC-ARCHIVE-LICENSE` accept an `absorbed_by` that names an
   archived Spec present in any of these forms:
   - an active Spec folder;
   - a legacy archived folder;
   - an Archive Record `<slug>.md` under the archive root;
   - a folder under the archive root that the repository's Git history at
     `HEAD` records and that was later removed.

   Read Git only when a license is otherwise unresolved, with at most one
   `git` process per check, for example
   `git log --format= --name-only --diff-filter=D HEAD -- <archive-root>`.
   A slug that none of the four forms names still reports
   `SC-ARCHIVE-LICENSE`. Outside a Git repository, or when Git fails, the
   check falls back to the first three forms and never panics. MUST NOT
   rewrite or remove any `absorbed_by` value or any file under
   `docs/history` to silence the check.
2. MUST add `internal/speccheck/archive_license_git_test.go` with
   `TestArchiveLicenseResolvesThroughTheRecordOrGit`. In a temporary Git
   repository it commits `docs/history/specs/0003-legacy/_prd.md`, then
   removes that folder in a later commit, and writes an Archive Record for
   `0002-record-only` with no folder. It writes three archived Findings
   whose `absorbed_by` names `0002-record-only`, `0003-legacy` and
   `0004-never-archived`. `speccheck.Check` reports exactly one
   `SC-ARCHIVE-LICENSE`, and it names `0004-never-archived`.
3. MUST rewrite the Archives cell of the qualifying declared `partial` row
   in all three `### QA settlement` tables (archive-spec, qa-gate and the
   Roundfix skill) to the `_techspec.md` → Exact texts wording: "The Archive
   Record, which names the QA Report and verdict and carries the
   declarations' `satisfied-by` record as `unproven`; the Spec, its QA
   report and evidence stay in Git at the record's `source_revision`." The
   three sections stay byte-identical, as
   `TestSettlementGuidanceIsOneTable` requires. No other cell, outcome name
   or Settles text changes.
4. MUST rewrite the archive-spec skill outside `### QA settlement` so it
   describes only the Archive Record flow:
   - the `description` front matter and the purpose paragraph say the
     archive writes `<slug>.md` under the resolved archive root and removes
     the Spec folder, which stays in Git at the record's `source_revision`.
     They no longer say "then stamp the archive metadata and move" or "with
     the completion stamped in its frontmatter";
   - `## Steps` says the normal archive "verifies the preconditions, writes
     the Archive Record and removes the Spec folder", replacing "stamps the
     archive metadata, and moves the folder";
   - `## Unarchive` says to "Restore the Spec folder from Git" at the
     record's `source_revision` and delete the record, replacing the
     `git mv` instruction.
5. MUST rewrite `## Archive Command` in the Roundfix skill's
   `references/archive.md` so it no longer says the archive stamps
   `_prd.md`, "then moves" the folder, or rewrites links "Before the move".
   It states the record, the `unproven` declarations of a qualifying
   `partial`, and the override fields the record carries. Its confirmation
   example is Surface Transcript 1's line, ending
   `kept in Git at <12-hex>`, and it names the `; promoted <n> file(s) to
   docs/references/` suffix. Existing legacy folders keep their semantics;
   say so in one sentence if needed.
6. MUST add to `docs/user-guide/commands/archive.md`:
   - a plan example whose first line is
     `archive plan for <slug>: removes <n> file(s) (<b> bytes) and writes docs/history/specs/<slug>.md`;
   - the keyless candidate form `no advice (<KEY_VARIABLE> is not set)`;
   - the promotion suffix `; promoted <n> file(s) to docs/references/`.

   This satisfies the Vocabulary Contract pattern `archive plan for`.
7. MUST edit only the canonical copies under `.agents/skills/`, run
   `make skills-sync`, and then run
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   That record command picks the next free version for archive-spec, qa-gate
   and roundfix (since Spec 0228). It writes both version fields of the
   canonical and mirror `SKILL.md`, and appends the entries to
   `skills/testdata/owned-skill-versions.json`. Never pick or hand-edit a
   version, and never replace a recorded digest.
8. MUST keep every other sentence of the edited skills and guides, and the
   two Baseline clauses, byte-identical. Run `make baseline-digests` and
   confirm it leaves no diff. If it changes a derived file, declare that file
   in `## Recorded paths`.
9. MUST prove the ablation: the Verification clones the working tree,
   removes every folder under `docs/history/specs` in a commit, and runs
   `TestCheckCorpusGolden` and `TestCheckActiveCorpusHasNoErrors` there.
   Both must pass. The clone lives in a temporary directory and is removed.

## Subtasks

- [ ] Resolve absorbed Specs through the record or Git, with the fixture test.
- [ ] Rewrite the partial settlement cell identically in the three skills.
- [ ] Rewrite the archive-spec workflow and the Roundfix archive reference.
- [ ] Document the plan, keyless and promotion lines in the archive guide.
- [ ] Sync the mirrors and record the raised versions.
- [ ] Run the ablation and the repository gates.

## Acceptance Criteria

- [ ] With every Spec folder under `docs/history/specs` removed, both corpus
      tests pass.
- [ ] A license naming a slug that never existed still reports
      `SC-ARCHIVE-LICENSE`.
- [ ] No shipped guidance promises stamping, moving or link rewriting for a
      new archive.
- [ ] The three settlement tables carry the Exact texts and stay
      byte-identical.

## Context

- creates: `internal/speccheck/archive_license_git_test.go`
- interface: `internal/speccheck/citations.go`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/archive.md`
- instruction: `internal/speccheck/archive_record_test.go`
- instruction: `internal/docscontract/corpus_test.go`
- instruction: `skills/settlement_guidance_repocontract_test.go`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- instruction: `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchiveLicenseResolvesThroughTheRecordOrGit|TestArchivedSlugsIncludeArchiveRecords)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveLicenseResolvesThroughTheRecordOrGit TestArchivedSlugsIncludeArchiveRecords; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the new test does not exist, so the command fails; with the test and without the resolver change, it reports `0003-legacy` as unresolved.
- `d="$(mktemp -d)"; trap 'rm -rf "$d"' EXIT; git clone -q --no-local . "$d/r" && git diff --binary HEAD | git -C "$d/r" apply --allow-empty && git ls-files -z -o --exclude-standard | tar --null -T - -cf - | tar -xf - -C "$d/r" && cd "$d/r" && git add -A && git -c user.name=ablation -c user.email=ablation@example.invalid commit -qm work --allow-empty && git rm -rq -- docs/history/specs/*/ && git -c user.name=ablation -c user.email=ablation@example.invalid commit -qm ablate && ! ls -d docs/history/specs/*/ >/dev/null 2>&1 && out="$(go test -count=1 -v -tags docscontract ./internal/docscontract -run '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$' 2>&1)" || { printf '%s\n' "$out" | tail -40; exit 1; }; for name in TestCheckCorpusGolden TestCheckActiveCorpusHasNoErrors; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the ablated clone reports unresolved `absorbed_by` licenses and the undocumented `archive plan for`, so the command fails.
- `for f in .agents/skills/archive-spec/SKILL.md .agents/skills/qa-gate/SKILL.md .agents/skills/roundfix/SKILL.md skills/archive-spec/SKILL.md skills/qa-gate/SKILL.md skills/roundfix/SKILL.md; do n="$(tr -s '[:space:]' ' ' < "$f")"; printf '%s' "$n" | grep -qF -- "The Archive Record, which names the QA Report and verdict and carries the declarations'" || { printf 'missing partial settlement wording in %s\n' "$f" >&2; exit 1; }; if printf '%s' "$n" | grep -qF -- "The Spec, its QA report and evidence, and the declarations'"; then printf 'old partial settlement wording in %s\n' "$f" >&2; exit 1; fi; done; for f in .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md; do n="$(tr -s '[:space:]' ' ' < "$f")"; for old in 'then stamp the archive metadata and move' 'with the completion stamped in its frontmatter' 'stamps the archive metadata, and moves the folder' 'git mv'; do if printf '%s' "$n" | grep -qF -- "$old"; then printf 'stale archive guidance in %s: %s\n' "$f" "$old" >&2; exit 1; fi; done; for new in 'verifies the preconditions, writes the Archive Record and removes the Spec folder' 'Restore the Spec folder from Git'; do printf '%s' "$n" | grep -qF -- "$new" || { printf 'missing archive guidance in %s: %s\n' "$f" "$new" >&2; exit 1; }; done; done; for f in .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md; do n="$(tr -s '[:space:]' ' ' < "$f")"; if printf '%s' "$n" | grep -qE -- 'archive stamps .?_prd[.]md|It then moves|Before the move, the archive rewrites'; then printf 'stale archive reference in %s\n' "$f" >&2; exit 1; fi; printf '%s' "$n" | grep -qF -- 'kept in Git at <12-hex>' || { printf 'missing record confirmation in %s\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md && cmp .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md && go test -count=1 ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' && go test -count=1 -tags repocontract ./skills -run '^TestSettlementGuidanceIsOneTable$'` — expected: exit 0; before this Task all six skill files keep the old partial cell and the archive-spec skill still says it stamps and moves, so the command fails.
- `n="$(tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md)"; for p in 'archive plan for <slug>: removes <n> file(s)' 'no advice (<KEY_VARIABLE> is not set)' '; promoted <n> file(s) to docs/references/'; do printf '%s' "$n" | grep -qF -- "$p" || { printf 'missing in archive guide: %s\n' "$p" >&2; exit 1; }; done` — expected: exit 0; before this Task the guide documents none of the three lines, so the command fails.

## References

- `_prd.md` → Core Feature 4; Core Feature 6; Success Metric 4
- `_techspec.md` → Readers; Exact texts; Vocabulary Contract; Surface Transcripts 1, 2, 3 and 4
- QA Report `qa-report-2026-10-06.md` → F-03; F-04; rows 09 and 14
- ADR-0247; ADR-0189; ADR-0233
