Measured at: 513b22ebab2d62a23b4222073c937c2862907ec2
History files: 5271
History bytes: 42686990
QA evidence files: 2089
QA evidence bytes: 12725738
Binary files: 18
Binary bytes: 12179652

This measurement is read-only. Nothing under docs/history was deleted, moved or rewritten by this measurement.

## Inventory

`git ls-tree -r -l 513b22ebab2d62a23b4222073c937c2862907ec2 -- docs/history`
reported the totals above. It also reported 326 QA Reports (5,549,584 bytes),
720 core Spec artifacts (`_prd.md`, `_techspec.md`, `_tasks.md` and
`_authorization.md`; 5,244,976 bytes), and 672 files in other history
families (2,426,993 bytes). Binaries are paths ending in `.png`, `.jpg`,
`.jpeg`, `.gif`, `.pdf`, `.db` or `.zip`. The measured commit is an ancestor
of `HEAD`; `git diff --name-status 513b22ebab2d62a23b4222073c937c2862907ec2
-- docs/history` has no deleted or renamed path.

## Readers

The static scan was `git grep -l -I -e docs/history -- ':!docs/history/**'`
plus `git grep -l -I -e ArchiveDir -e ArchiveSpecRoot HEAD --
':!docs/history/**'`. It found 83 tracked files naming the history root and
25 building the archive root. The naming readers are the archive, brainstorming,
Roundfix, write-idea, write-prd, write-tasks and write-techspec skills and
references; `.coderabbit.yaml`; `CHANGELOG.md`; ADRs 0120 and 0138 through
0149, 0163 and 0173; `docs/agents/docs-layout.md`,
`docs/agents/specific-repository.md`; the dated backlog entry; the coverage
record; active Specs 0217 and 0218; the archive and context-driven user
guides; baseline assets and tests; archive, baseline, CLI, docscontract,
judge, spec, speccheck, suiteguardcontract and worktree Go sources and tests;
and their mirrored `skills/` files. The exact 83 paths are the output of the
command above, retained in the measurement session; the 25 archive-root
builders are `internal/baseline/history_layout.go`,
`internal/cli/{archive_test,deliver,deliver_workflow,review,spec_check,supersede}.go`,
`internal/docscontract/publicdocs_test.go`, `internal/spec/{archive.go,archive_layout_characterization_test.go,archive_test.go}`,
`internal/specaudit/audit.go`, `internal/speccheck/{backlog,citations}.go`,
`internal/speccheck/{citations_test,constraints_characterization_test}.go`,
`internal/worktree/{merged_head,merged_head_test,merged_spec_absent_target_test,merged_spec_leftovers_test,worktree,worktree_test}.go`,
`skills/owned_skill_edit_repocontract_test.go`, plus the corresponding
`.agents/` skill references and `docs/references/coverage-record.json`.

Direct archived QA-evidence readers outside history are
`docs/references/coverage-record.json` and `internal/cli/review_scope_test.go`.
They respectively retain coverage package inputs and prove that archived QA
evidence is omitted from review scope. No outside file names an archived
binary as an input.

## Ablation

The clone was created outside the repository and Roundfix Home with
`git clone --no-local <repository> /private/tmp/roundfix-0214-measure.4pOmZ0/clone`.
QA evidence directories and binary extensions were removed only there, and
the clone was removed afterwards. The timed command was
`/usr/bin/time -p git -C <clone> grep -n 'docs/history' -- . >/dev/null`,
five times per phase. Pre-removal `real` samples were 0.17, 0.13, 0.15,
0.19, 0.17 seconds (median 0.15s); post-removal samples were 0.14, 0.11,
0.19, 0.28, 0.29 seconds (median 0.19s). All ten grep commands, clone and
removal commands exited 0.

`make verify` exited 2 after removal. Its failing test was
`roundfix/internal/spec` `TestCoverage`, reporting coverage regressions for
archived-evidence packages under Specs 0141, 0157 and 0162.
`make verify-docs` exited 2. Its failing test was
`TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical`, because
`docs/history/specs/0001-implement-command/qa/evidence/live-cli-transcript-2026-07-04.txt`
was absent. These failures show that removal does not leave either gate at
exit 0.

## Agent reads

The Run Database was opened only as
`file:/Users/marcio/.roundfix/roundfix.db?mode=ro&immutable=1`, selecting only
`agent.tool_started` rows. Query used:

```sql
SELECT run_id, cursor, kind, tool_id, tool_state, summary, payload
FROM run_events
WHERE kind = 'agent.tool_started' AND payload LIKE '%docs/history/%';
```

There were 876 matching events across 127 distinct Runs:

| path kind | events | Runs |
| --- | ---: | ---: |
| QA evidence | 40 | 19 |
| QA Report | 343 | 50 |
| core Spec artifact | 125 | 59 |
| other | 368 | 104 |

By `json_extract(payload, '$.params.update.kind')`, tool kinds were
`execute` 476 / 110 Runs, `edit` 399 / 73, and `think` 1 / 1.

## Secondbrain mirror

`/Users/marcio/dev/secondbrain/wiki/index.md` was read first. The required
`qmd query "Roundfix archived history evidence mirror reads" --all --files
--min-score 0.3` was attempted, but its embedding backend failed to
initialize and its cache was read-only; no result was used and this limitation
is recorded here. The read-only mirror at
`/Users/marcio/dev/secondbrain/projects/roundfix/mirror/docs/history` contains
5,271 files and 42,686,990 bytes, equal to the Git inventory.

## Proposal

Candidates: none

Decision rule 3 requires no repository reader outside history and an ablation
that leaves both `make verify` and `make verify-docs` at exit 0. QA evidence
has direct readers in the coverage record and review-scope tests, and the
combined ablation fails both gates. Binaries have no direct static reader,
but the same ablation fails both gates. Agent reads and the mirror are costs,
not vetoes. No files or bytes are proposed for removal, no Agent reads would
be lost, and no prerequisite test can be declared repaired. No Backlog Entry
was written. Any later removal waits for the maintainer's explicit approval.

## What removal would not reclaim

History keeps every removed version. Removing a file from the current tree
does not remove its reachable versions from Git history, and a clone can
still download them. This cites the Git book page named by the PRD's
Acceptance evidence: <https://git-scm.com/book/en/v2/Git-Internals-Maintenance-and-Data-Recovery>.
