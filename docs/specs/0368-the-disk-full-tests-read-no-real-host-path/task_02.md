---
task: task_02
spec: 0368-the-disk-full-tests-read-no-real-host-path
status: pending
type: qa
complexity: low
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec makes
`doctorDiskSpaceAncestor` walk past a file that stands where a directory is
expected (`ENOTDIR`), as it walks past a missing path, and moves every
disk-full test path that reaches the filesystem under the test's own
`t.TempDir()`, so a file on the host can no longer change those tests.

Every row runs in this repository or in temporary directories it creates
and removes. No row creates, removes or changes any file in `/tmp`,
`/private/tmp` or `$HOME` outside a directory made by `mktemp -d`, reaches
the network, reads another repository, starts an Agent session or an ACP
adapter, reads or writes the real `~/.roundfix`, sends a request to a
provider, the judge or a reviewer, records an environment value, or sends a
signal to any process. Every `go test` carries an explicit `-timeout` and
keeps the session's own `HOME` and `TMPDIR`. Every row that runs
`bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary directories, with
`GOMODCACHE` and `GOCACHE` kept at their `go env` values, and removes them
afterwards. The QA never writes into the checkout while a test suite runs.
The gate keeps its evidence under names the repository does not ignore,
`.txt` for command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 20m -v`,
   `TestADiskFullPathBelowAFileIsReadAtItsNearestDirectory` and
   `TestTheDiskSpaceAncestorWalksPastAFile` in `./internal/cli`, then
   `-run 'DiskFull|DiskSpace|LowOnDisk' ./internal/cli` and
   `-run 'DiskFull|LowDisk' ./internal/delivery ./internal/store`, recording
   each pass, and whether the subtest
   `TestTheDiskSpaceAncestorWalksPastAFile/an_unreadable_directory_still_fails`
   passed or was skipped (it skips only on Windows or as root).
3. MUST record Success Metric 1: `grep -nE -- 'tm[p]/' internal/cli/deliver_disk_full_test.go`
   prints nothing, and the four tests
   `TestTheItemBuildOnAFullDiskIsADiskFullFailure`,
   `TestTheRepositoryGateOnAFullDiskParksDiskFull`,
   `TestDiskFullFailureUsesOnlyTheLogTailAndPreservesTheError` and
   `TestADiskFullPathMayHoldSpaces` pass; compare their assertions with the
   Delivery Base (the merge base of the audited head with `main`) through
   `git diff <delivery base> -- internal/cli/deliver_disk_full_test.go` and
   record that only paths changed in them.
4. MUST record Success Metrics 2 and 3 against the starting behavior: in a
   `git archive` of the Delivery Base unpacked in a temporary directory,
   copy the audited head's `internal/cli/deliver_disk_full_test.go` and
   `internal/cli/doctor_disk_space_test.go` over the archive's copies and run
   the two new tests there with `-timeout`; record that both fail, the first
   with 0 bytes and no capacity read, the second on the path below a file
   with `not a directory`. On the audited head both pass (row 2).
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails:
   - `doctorDiskSpaceAncestor` loses the `syscall.ENOTDIR` condition:
     `TestADiskFullPathBelowAFileIsReadAtItsNearestDirectory` and
     `TestTheDiskSpaceAncestorWalksPastAFile` fail;
   - `doctorDiskSpaceAncestor` walks to the parent on every stat error:
     `TestTheDiskSpaceAncestorWalksPastAFile/an_unreadable_directory_still_fails`
     fails (record the skip instead when the gate runs on Windows or as
     root).
6. MUST execute the tests every non-QA Task names in its Verification.
7. MUST record, as the outside-evidence row, sources this Spec did not
   write, with their dates and not blocked: the adopted Backlog Entry,
   recorded by the operator on 2026-11-09 before this Spec existed, whose
   symptom ("made the reported free space `0 MiB` instead of the stubbed
   value") matches row 4's starting behavior; and the POSIX `stat()` error
   list quoted in `_prd.md` → Acceptance evidence ("ENOTDIR: A component of
   the path prefix is not a directory."), compared with row 2, where a path
   below a file resolves to its nearest existing directory. The row reads
   the entry in the Spec's `references/`, never the Run diagnostics under
   `~/.roundfix`.
8. MUST perform the glossary check (`docs/agents/domain.md`): record that
   the PRD and TechSpec declare `## Glossary` as `None.` and that the Task
   commit range does not change `CONTEXT.md`.
9. MUST verify from the repository history:
    - task_01's changed files stay within its declarations or its `##
      Recorded paths`;
    - `_authorization.md` lists `paths: []` and no Governed Path changed;
    - the Task commit range changes none of these: the Makefile, any
      workflow, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/**`,
      `internal/cli/deliver_disk_full.go`, `internal/cli/deliver_free_space.go`,
      any guide, skill, ADR, record or `CONTEXT.md`.

    This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0368-the-disk-full-tests-read-no-real-host-path --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read another repository, or run `git fetch`, `git push` or `gh`;
    - read or write `~/.roundfix`;
    - create a file in `/tmp`, `/private/tmp` or `$HOME` outside a
      `mktemp -d` directory it removes;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - send a signal to any process;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification and the named suites.
- [ ] Characterize the Delivery Base with the new tests and run each sabotage on a temporary copy.
- [ ] Record the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A failing path below a file is measured at its nearest existing directory on the audited head and not on the Delivery Base.
- [ ] No disk-full test names a `/tmp/` path, and the existing disk-full tests keep their assertions.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification


- `newest="$(find 'docs/specs/0368-the-disk-full-tests-read-no-real-host-path/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-2; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence; Glossary
- `_techspec.md` → Interfaces; Testing Approach; Build Order 2
