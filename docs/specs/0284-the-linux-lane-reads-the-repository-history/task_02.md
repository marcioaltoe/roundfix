---
task: task_02
spec: 0284-the-linux-lane-reads-the-repository-history
status: completed
type: qa
complexity: high
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec delivers one change: the
Linux Lane mounts the repository's object store read-only at `/repo-objects`
and registers it as the container repository's alternate, in the Delivery
Queue gate and in `roundfix deliver lane`, and skips a store it cannot mount
whole with its reason.

Every row runs on temporary repositories and homes with compiled Docker and Go
stand-ins, except the one real lane row, which runs only when the Run sandbox
reaches a Docker daemon. No row reads or writes the real `~/.roundfix`,
reaches the network, pulls an image, sends a request to a provider, the judge
or a reviewer, reads a credential, or reads another repository. The gate
never runs `roundfix spec judge`, never runs `roundfix review` against a real
provider, and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1:
   - Run `TestLinuxLaneMountsTheRepositoryHistoryReadOnly` and
     `TestLinuxLaneRunsTheChangedPackagesInAContainer`, and record, for the
     gate from a linked worktree, `deliver lane --all` and `deliver lane` on
     changed packages, the two `-v` values and that the store is the main
     repository's (Invariants 1, 2 and 5).
   - On a temporary copy of the audited head made with `git archive`, remove
     the object store's `-v` pair from the lane's arguments and record that
     both tests fail; on another, remove `echo /repo-objects >
     .git/objects/info/alternates` from the script and record that both fail
     (Invariant 3); on a third, change the mount's `:ro` to `:rw` and record
     that both fail.
3. MUST verify Success Metric 2:
   - Run `TestLinuxLaneSkipsAnObjectStoreItCannotMount` and
     `TestDeliverLaneObjectStoreSkipTranscript`, and record each skip reason
     and that no Docker or Go start was recorded (Invariant 4).
   - On a temporary copy of the audited head, make the store check ignore
     `info/alternates` and record that
     `TestLinuxLaneSkipsAnObjectStoreItCannotMount` fails.
   - Through the built product, reproduce Surface Transcript 1 of
     `_techspec.md`: in a temporary repository whose `go.mod` declares `go
     1.26.1`, whose `.roundfixrc.yml` sets `delivery.linux_lane: false` and
     whose `.git/objects/info/alternates` names another temporary
     repository's object store, with `HOME` set to a temporary directory and
     compiled Docker and Go stand-ins first on `PATH`, run `bin/roundfix
     deliver lane --all` and record that standard output, standard error and
     the exit code match the transcript and that the stand-ins recorded no
     start. Remove the alternates file, run it again, and record that the
     Docker stand-in received `-v <objects>:/repo-objects:ro`.
4. MUST verify Success Metric 3 when the Run sandbox can, and otherwise record
   the row as declared unreachable:
   - When `docker info` exits `0` and `docker image inspect golang:1.26`
     finds the image, run `bin/roundfix deliver lane --all` in the audited
     head's worktree with `HOME` set to a temporary directory and
     `GOMODCACHE` set to the value `go env GOMODCACHE` printed before `HOME`
     changed. Record its outcome line, its exit code, and from its lane log
     the result lines of `roundfix/internal/baseline` and
     `roundfix/internal/spec` and the `-v` values of the logged command.
     Record environment variable names only, never their values. Write that
     evidence only after the lane ends.
   - Otherwise record `blocked (environment: Docker unreachable from the Run
     sandbox)` and count it in `rows_blocked_declared` against the PRD's
     Unreachable Acceptance.
5. MUST execute the tests task_01 names in its Verification.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read `.github/workflows/ci-verify.yml` at the audited head and record
     its `fetch-depth` value and the comment above it;
   - record from the adopted Backlog Entry that the operator's run on `main`
     failed only the two history tests;
   - because this Spec's authorization forbids network access, record each
     web source as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`.
7. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` carries the revised **Linux Lane**, linking
     ADR-0297;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
8. MUST verify from the repository history:
   - task_01's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, and `docs/agents/setup-context.json`;
   - the Roundfix skill version was raised and recorded;
   - no `### QA settlement` section changed;
   - no guide, agent guide or skill copy still says the lane mounts only the
     module cache.

   This row reads Task commits, so it declares a `commit_range` input.
9. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance. It supports the row with the pre-PR equivalent of each
   control:
   - approval from the recorded delivery authority;
   - checks from the Daemon's repository Verification at the audited head;
   - no unresolved threads, because no Pull Request exists;
   - no Merge-Ready acceptance yet;
   - the audited head named as the claimed candidate.
10. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0284-the-linux-lane-reads-the-repository-history --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It MUST NOT
    do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - start a workflow or pull an image.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Reproduce Surface Transcript 1 and the mount through the built binary.
- [ ] Run the real lane when Docker is reachable, record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every lane container mounts the object store read-only and borrows from it, and each removed guard makes its named tests fail.
- [ ] Each store the lane cannot mount whole skips with its reason, and `roundfix deliver lane` matches Surface Transcript 1.
- [ ] The real lane run passes or is declared unreachable, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0284-the-linux-lane-reads-the-repository-history/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; Invariants 1-5; Current behavior; Build Order 2
- ADR-0297; ADR-0294; ADR-0288; ADR-0104; ADR-0240; ADR-0276
