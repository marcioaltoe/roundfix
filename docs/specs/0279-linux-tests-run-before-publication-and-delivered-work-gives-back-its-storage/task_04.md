---
task: task_04
spec: 0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage
status: completed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec delivers three things:

- the Linux Lane of the Delivery Queue gate and this repository's opt-in;
- the directory-scoped leftover rule and the named dirty reason of Run
  Worktree Reconciliation;
- Item Binary reclamation at the merged release and in the GC Command, with
  the retention policy.

Every row runs on temporary repositories, homes and copies of the audited head.
No row reads or writes the real `~/.roundfix`, reaches the network, pulls an
image, sends a request to a provider, the judge or a reviewer, reads a
credential, or reads another repository. The gate never runs
`roundfix spec judge`, never runs `roundfix review` against a real provider,
never runs `roundfix gc` or `roundfix reconcile --apply` against the real
Roundfix Home, and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1:
   - Run the tests task_01 names, and record from their output the
     `docker run` arguments the fake `docker` observed, each skip reason, the
     failing gate's reason, and that the gate without the lane started no
     `docker` process (Invariants 1 and 2).
   - Make a temporary copy of the audited head with `git archive`. In it,
     remove `--pull never` from the lane's arguments, and record that
     `TestLinuxLaneRunsTheChangedPackagesInAContainer` fails.
   - Record that `.roundfixrc.yml` sets `delivery.linux_lane: true` and that
     `bin/roundfix` loads it without a config error from a temporary home.
3. MUST record the real-container row. Read the image name from the audited
   head's `go.mod` `go` directive. When `docker image inspect golang:<that
   version>` succeeds on the gate's machine, unpack `git archive` of the
   audited head into a temporary directory and run, with `--pull never`,
   `--network none`, that directory and the host module cache mounted
   read-only, and a writable copy made inside the container,
   `go test -count=1 -run '^TestDeliveryItemBinaryCopyStartsWhileTheOwnerForks$' ./internal/cli`;
   record its outcome. When the image is absent, or `docker` is unavailable,
   record the row as `blocked (environment: image golang:<version> not present locally)`
   or `blocked (environment: docker unavailable)`, with the exact output of
   the failed check. The gate never pulls the image.
4. MUST verify Success Metric 2:
   - Run the tests task_02 names, and record that a leftover beneath a
     declared directory is `superseded` and removed by apply, that the
     prefix-sharing sibling stays `dirty`, and the exact named reasons
     (Invariants 3 and 4).
   - On a temporary copy of the audited head, restore exact-equality matching
     in the leftover scope, and record that
     `TestMergedSpecRunWithLeftoversBeneathADeclaredDirectoryIsSuperseded`
     fails.
5. MUST verify Success Metric 3:
   - Run the tests task_03 names, and record the dry-run lines, that a parked
     item keeps its Item Binaries, and that the release keeps both logs
     (Invariant 5).
   - On a temporary copy of the audited head, make the GC sweep treat every
     stage as reclaimable, and record that
     `TestGCRemovesOnlyItemBinariesOfMergedOrUnqueuedItems` fails.
   - Record that the existing GC tests pass unchanged, which holds the output
     of a repository without Item Binaries byte-identical.
6. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence.
   - Read the repository-history rows offline with `git show`: the first line
     of `internal/cli/deliver_item_binary_fork_linux_test.go` as commit
     `23ab0904` added it, the two `deletes:` declarations of Spec 0272's
     task_04 at its Archive Record's `source_revision`, and the last
     consequence of ADR-0282 at the Delivery Base.
   - Because this Spec's authorization forbids network access, record each web
     source as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`.
8. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` defines **Linux Lane** and **Item Binary** and
     carries the revised **Run Worktree Reconciliation** and **GC Command**;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its
     `## Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - in `.roundfixrc.yml` only the `delivery.linux_lane: true` line and its
     comment changed;
   - these did not change: the Makefile, the CI workflows, `go.mod`, `go.sum`,
     every Baseline module, `docs/agents/setup-context.json`,
     `docs/references/coverage-record.json`, `internal/delivery/engine.go`
     and every file under `internal/judge/`;
   - the Roundfix skill version was raised and recorded;
   - no `### QA settlement` section changed;
   - no guide or skill copy still says Item Binary copies are not reclaimed.

   This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It MUST NOT
    do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Run or block the real-container row without pulling an image.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The lane runs one bounded container over the changed packages, skips with named reasons, and each removed guard makes its named test fail.
- [ ] A leftover beneath a declared directory is released, and a kept dirty Run names why.
- [ ] Item Binaries go only for merged or unqueued items, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- `_prd.md` → Goals; User Stories 1-5; Core Features 1-7; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; API Contract 9; API Contract 10; Invariants 1-5; Current behavior; Build Order 4
- ADR-0288; ADR-0289; ADR-0290; ADR-0212; ADR-0282; ADR-0104; ADR-0276
