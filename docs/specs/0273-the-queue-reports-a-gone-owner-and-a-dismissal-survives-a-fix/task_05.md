---
task: task_05
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes `roundfix deliver status`
name a recorded Queue Owner proven not running, makes the queue commands run
in an item worktree name the checkout that owns the queue, and makes `roundfix
review dispose --dismiss` accept a dismissal after a fix that leaves the
finding's anchor unchanged. No row reads or writes the real `~/.roundfix`,
reaches a provider, reads a credential or reads another repository.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 and Surface Transcript 1 through the built
   `bin/roundfix` with `HOME` set to a temporary home, `NODE_OPTIONS` unset and
   a temporary Git repository. It seeds that home's Run Database through
   `internal/store` (`CreateDeliveryQueue`, `UpdateDeliveryQueueItem`,
   `ClaimDeliveryQueueOwner`) from a program kept in a temporary worktree of
   the audited head, never in the audited tree, or with `sqlite3` against the
   schema the binary created. For one item at `checking` it records the owner
   PID `1073741824` (`1<<30`) and compares `roundfix deliver status` with the
   transcript; records the PID of a `sleep` process the gate started with an
   identity that carries the platform prefix (`<GOOS>:`) but another start
   time and observes the different-identity line; records it with an
   identity without that prefix and observes no `Owner:` line and the
   could-not-be-checked diagnostic on stderr; records it with the identity
   `store.OwnerProcessIdentity` returns and observes no `Owner:` line; and
   with every item `merged` or `parked` observes no line.
   After each run it reads the owner record back unchanged. It ends the
   `sleep` process it started through that process's own handle.
3. MUST verify Success Metric 3 and Surface Transcripts 2 and 3 through the
   built binary: in the same kind of fixture, an item worktree made with
   `git worktree add` and recorded through `RecordDeliveryQueueItemWorktree`;
   `roundfix deliver retry <slug>` and `roundfix deliver status` run in it
   match the transcripts, `deliver resume` and `deliver stop` there name the
   owning checkout, and the queue and owner record are unchanged afterwards.
   In an unrelated temporary repository the four refusals are compared with
   what a binary built from the delivery base, the merge base of the audited
   head with `main`, prints for the same fixture, and are byte-identical apart
   from the temporary paths.
4. MUST verify Success Metric 2 and Surface Transcripts 4 and 5 through the
   built binary: a temporary repository whose `.roundfixrc.yml` sets
   `defaults.artifact_dir` to a temporary directory, a findings Pre-PR Review
   Record written where `roundfix review` reads it for that checkout
   (`<Artifact Directory>/pre-pr-review/<first 16 hex digits of the SHA-256 of
   the resolved checkout path>/pre-pr-review.json`) with two anchored standing
   findings, and a commit that fixes the second. It records that dismissing
   the first is appended once with `headCommit` equal to the reviewed head and
   `currentHead` equal to HEAD; that a later commit changing or moving the
   anchored lines makes a new dismissal refuse with the transcript's reason and
   the ledger's bytes unchanged; and that a dismissal at the reviewed head
   writes a line without `currentHead`, as the delivery base's binary writes
   it for the same fixture.
5. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification, and run `go test -count=1 -run 'Deliver|ReviewDispose|Disposition|Correction|Lineage|Reus' ./internal/cli`
   and `go test -count=1 -run 'Delivery' ./internal/store`, recording that
   they pass and that no existing test file of those packages changed except
   the new files the Tasks declare.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the psutil issue 2400 passage and the
   r-lib/ps passage on the PID and start-time pair, the GitHub Changelog and
   REST reference passages on review comment outdating, and the
   git-merge-base passage on `--is-ancestor` exit statuses, quoting each and
   stating whether it still supports the design; for a source it cannot reach
   it records `blocked (environment: network denied: <host>)`.
7. MUST perform the glossary check of `docs/agents/domain.md`: record that
   `CONTEXT.md` revises **Queue Owner** with the owner report and **Review
   Finding Disposition** with the dismissal after a fix, both citing
   ADR-0279, and whether any other term the Spec introduced needs one.
8. MUST repeat the old-wording sweep: `git grep` for `an unchanged reviewed`,
   `moved-head` followed by `dismissals`, and `to equal reviewed head` outside
   `docs/history/`, the Spec folder and Go test files finds nothing that
   still describes the old dismissal rule.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`; that every Governed Path
   changed is bounded in `_authorization.md`; that the Makefile, the CI
   workflows, `go.mod`, `go.sum`, every Baseline module and
   `internal/cli/review_conventions.go` did not change; that the Roundfix
   skill version was raised and recorded; and that no `### QA settlement`
   section changed. This row reads Task commits, so it declares a
   `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
11. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT read or write `~/.roundfix`, record any environment value, send any
    request to a provider, start a workflow, or signal a process it did not
    start.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Reproduce the five transcripts through the built binary on seeded fixtures.
- [ ] Compare the unchanged refusals and the reviewed-head dismissal with the delivery base.
- [ ] Record the outside evidence, the glossary check, the sweep and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Each transcript is reproduced through the built binary.
- [ ] A reused PID and a live owner are not reported gone, and no command changed the owner record.
- [ ] Refusals outside an item worktree and the reviewed-head dismissal match the delivery base.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcripts 1-5; Invariants 1-7; Current behavior; Build Order 5
- ADR-0279; ADR-0044; ADR-0268; ADR-0272
