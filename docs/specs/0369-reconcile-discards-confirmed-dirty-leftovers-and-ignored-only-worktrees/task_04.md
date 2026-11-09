---
task: task_04
spec: 0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec adds the Dirty
Discard, `roundfix reconcile --discard-dirty <digest>`, which records a
Branch Disposition naming an ended Run's uncommitted leftovers before it
removes the Run Worktree and Run Branch; it makes a recorded Run Worktree
that Git does not register and that holds only ignored files `safe`; it makes
a full-scan `--apply` remove empty repository directories under the Worktree
Root (ADR-0379); and it documents all three in ADR-0053, ADR-0289, ADR-0290,
`CONTEXT.md`, the `reconcile` user guide and the Roundfix skill's
`reconcile` reference.

Every row runs in this repository or in temporary repositories it creates
and removes. No row reaches the network, reads another repository, starts an
Agent session or an ACP adapter, reads or writes the real `~/.roundfix`,
sends a request to a provider, the judge or a reviewer, records an
environment value, or sends a signal to any process. No row creates files in
`/tmp`, `/private/tmp` or `$HOME` outside a directory it made with
`mktemp -d` and removes. Every `go test` carries an explicit `-timeout` and
keeps the session's own `HOME` and `TMPDIR`. Every row that runs
`bin/roundfix` or a binary built from the Delivery Base sets `HOME` and
`TMPDIR` to fresh temporary directories, with `GOMODCACHE` and `GOCACHE`
kept at their `go env` values, and removes them afterwards. The QA never
writes into the checkout while a test suite runs. The gate keeps its
evidence under names the repository does not ignore, `.txt` for command
output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 20m -v`, every test tasks 01
   and 02 name in `./internal/worktree`, `./internal/daemon` and
   `./internal/cli`, then the whole `./internal/worktree` and `./skills`
   packages, `-run 'Branch|Disposition|Superseded' ./internal/daemon`,
   `-run 'Reconcile|DirtyDiscard|ReleaseMergedRuns|AutomaticRelease|TestRunCommandHelp' ./internal/cli`
   and `-tags docscontract ./internal/docscontract`, recording each pass and
   that none of the nineteen new tests reports `SKIP`.
3. MUST build a second binary from the Delivery Base (the merge base of the
   audited head with `main`), from a `git archive` of that commit unpacked
   in a temporary directory, with `go build -buildvcs=false`, and in a
   temporary Git repository with one commit and fresh `HOME` and `TMPDIR`,
   run `reconcile` and `reconcile --format json` with both binaries: the
   outputs are byte-identical, so a report without a Dirty Discard
   candidate or an empty repository directory is unchanged.
4. MUST record Success Metrics 1 and 2 from the `-v` output of the task_01
   tests, naming for each metric the test that proves it:
   `TestReconcileDryRunListsTheDirtyDiscardPathsAndDigest`,
   `TestReconcileDiscardDirtyRecordsThenRemovesTheConfirmedRuns`,
   `TestReconcileDiscardDirtyRefusesAChangedDigestWithoutChangingAnything`,
   `TestReconcileDiscardDirtyRefusesWhenNoRunIsOffered`,
   `TestADirtyDiscardIsNotOfferedForARunWithALiveProcessTree`,
   `TestADirtyDiscardIsNeverOfferedOutsideItsBounds` and
   `TestDiscardDirtyTerminalRunRefusesAChangedPathSet`; and Success Metrics
   3 and 4 from the task_02 tests
   `TestAnUnregisteredRunWorktreeHoldingOnlyIgnoredFilesIsSafe`,
   `TestAnUnregisteredRunWorktreeKeepsTodaysUnknown`,
   `TestReconcileApplyReleasesAnIgnoredOnlyRunWorktreeAndEmptyRepositoryDirectories`,
   `TestReconcileWithARunIDRemovesNoEmptyRepositoryDirectory` and
   `TestEmptyRepositoryDirectoriesListOnlyEmptyRepositoryShapedNames`.
5. MUST reproduce each TechSpec Surface Transcript through `bin/roundfix`
   with fresh `HOME` and `TMPDIR`: Surface Transcript 1 (a Dirty Discard
   needs its digest), Surface Transcript 2 (a malformed digest is refused
   before any inspection) and Surface Transcript 3 (the four mutation
   switches exclude each other), comparing stdout, stderr and the exit code.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails:
   - `DiscardPaths` is set whatever the commit proof says:
     `TestADirtyDiscardIsNeverOfferedOutsideItsBounds` fails on the
     unrepresented Task commit;
   - the digest covers the Run IDs and heads but not the paths:
     `TestReconcileDiscardDirtyRefusesAChangedDigestWithoutChangingAnything`
     fails, because the added file no longer changes the digest;
   - `--discard-dirty` skips `WriteDirtyDiscardRecord`:
     `TestReconcileDiscardDirtyRecordsThenRemovesTheConfirmedRuns` fails;
   - the live-process filter is removed:
     `TestADirtyDiscardIsNotOfferedForARunWithALiveProcessTree` fails;
   - `git check-ignore` gains `--no-index`:
     `TestAnUnregisteredRunWorktreeKeepsTodaysUnknown` fails on the tracked
     path;
   - the apply path does not re-list the files:
     `TestApplyRefusesAnIgnoredOnlyWorktreeWhoseFilesChanged` fails;
   - the name shape check is removed from `EmptyRepositoryDirectories`:
     `TestEmptyRepositoryDirectoriesListOnlyEmptyRepositoryShapedNames` fails
     on `notes`;
   - the sweep also runs with a Run ID:
     `TestReconcileWithARunIDRemovesNoEmptyRepositoryDirectory` fails.
7. MUST execute the tests every non-QA Task names in its Verification.
8. MUST record, as the outside-evidence rows, sources this Spec did not
   write, with their dates and not blocked:
   - the Fiscus report quoted in `_prd.md` → Acceptance evidence ("summary
     dirty=11 unknown=1", "O diretório só tinha .storage/postgres/18/docker"),
     compared offline with row 4, where the same shapes, an outside leftover
     of an archived Spec's ended Run and an unregistered directory holding
     only `.storage/postgres/18/docker/base/1`, leave through the supported
     switches;
   - the `git check-ignore` manual passages that section quotes ("By
     default, tracked files are not shown at all since they are not subject
     to exclude rules" and "None of the provided paths are ignored."),
     compared with a probe in a temporary repository under the QA's own
     temporary directory: `git check-ignore -z --stdin` prints a path below an
     ignored `.storage/` directory, does not print a tracked path that a
     pattern matches, and exits 1 when no path is ignored.
9. MUST record that the PRD's `## Unreachable Acceptance` section is
   `None.`, so the report declares no blocked-declared row.
10. MUST record: the delivered
    `bin/roundfix spec check 0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees --strict`
    on the audited head reports no finding.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `- adds:` for **Dirty Discard** and
    `- changes:` for **Branch Disposition**, **Reconcile Command** and **Run
    Worktree Reconciliation**, that the Task commits change exactly those
    four entries of `CONTEXT.md`, and that
    `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryStatesEachRuleOnce|TestGlossaryClosedListsEqualTheirConstants|TestGlossaryAvoidLinesArePlainLists)$'`
    passes.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - ADR-0053, ADR-0289 and ADR-0290 carry the `Amended by [ADR-0379]`
      paragraph and `updated_at: 2026-11-09T15:00:00Z`;
    - the Task commit range changes none of these: the Makefile, any
      workflow, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/**`, any file
      under `internal/baseline/`, `internal/spec/`, `internal/speccheck/` or
      `internal/suiteguardcontract/`, `internal/cli/deliver_disk_full*.go`,
      `internal/cli/doctor_disk_space.go`, any guide under `docs/agents/`,
      any skill other than `roundfix`, any ADR other than ADR-0053, ADR-0289
      and ADR-0290, or text inside the `### QA settlement` section of any
      skill.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
14. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read another repository, or run `git fetch`, `git push` or `gh`;
    - read or write `~/.roundfix`, or run `roundfix reconcile` with
      `--apply` or `--discard-dirty` outside a temporary home and repository;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - send a signal to any process;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build both binaries and run the repository Verification and the named suites.
- [ ] Compare the reports without candidates and reproduce the three Surface Transcripts.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the outside sources, the glossary check, the self-check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A Dirty Discard lists, confirms, records and removes only an ended Run of an archived Spec whose commit proof holds, and every refusal changes nothing.
- [ ] An ignored-only unregistered directory and empty repository directories leave on a full-scan `--apply`, and every other case keeps today's report.
- [ ] Each sabotage makes its named test fail, this Spec passes its own strict check, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-4; Surface Transcripts 1-3; Testing Approach; Build Order 4
