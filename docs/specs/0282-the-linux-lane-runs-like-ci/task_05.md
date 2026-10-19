---
task: task_05
spec: 0282-the-linux-lane-runs-like-ci
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec delivers four things:

- `testfixture.RequireTool` and the External Test Tool rule, with the rust-cli
  journey and the acpx integration test moved onto it;
- `--init` and `delivery.linux_lane_image` for the Linux Lane;
- the Lane Container Failure and the Lane Change Failure, with the
  `linux-lane-container` blocker in the `environment` Park Class;
- `roundfix deliver lane`, the on-demand lane.

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
   - Run the config and lane tests task_02 names, and record that `--init`
     precedes the image, that the fetched default's image was used while the
     item branch's was ignored, and each refused and accepted value.
   - On a temporary copy of the audited head made with `git archive`, remove
     `--init` from the lane's arguments and record that
     `TestLinuxLaneRunsTheChangedPackagesInAContainer` fails (Invariant 1);
     on another, make the gate read the image from the loaded checkout
     configuration and record that `TestLinuxLaneRunsTheConfiguredImage`
     fails (Invariant 2).
3. MUST verify Success Metric 2:
   - Run the rust-cli journey with a `PATH` that holds only `go`, `git` and
     `sh`, and record its `--- SKIP` line and `requires rustfmt, which is not
     on PATH`.
   - On a temporary copy of the audited head, add `exec.LookPath("cargo")` to
     one test file and record that
     `TestEveryExternalTestToolGoesThroughRequireTool` fails naming that file
     and line.
   - Record the class, not only the case: build a temporary directory of
     links to every executable on the gate's `PATH` except `rustfmt`, `cargo`,
     `rustc`, `node`, `npm`, `npx`, `bun` and `docker`, run `go test
     -count=1 ./...` with `PATH` set to it, and record that no failing test's
     output names a missing executable (Invariant 3). Write that evidence
     only after the suite ends.
4. MUST verify Success Metric 3:
   - Run the lane, classification, gate and park tests task_03 names, and
     record each container cause, each change reason, each blocker, each
     Park Class and the lane log's last line.
   - On a temporary copy of the audited head, classify exit `127` as a Lane
     Change Failure and record that
     `TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure` fails
     (Invariant 4); on another, make `gateCandidate` ignore `result.Blocker`
     and record that `TestGateParksWithTheBlockerItNames` fails.
   - Record from the delivery tests that a repository Verification failure
     still parks the bare `gate-failed` (Invariant 5).
5. MUST verify Success Metric 4 through the built product: in a temporary
   repository whose `go.mod` declares `go 1.26.1` and whose `.roundfixrc.yml`
   sets `delivery.linux_lane: false`, with `HOME` set to a temporary
   directory and compiled Docker and Go stand-ins first on `PATH`, run
   `bin/roundfix deliver lane --all` for each of Surface Transcripts 1, 2, 3
   and 4 of `_techspec.md`, and record that standard output, standard error
   and the exit code match each transcript. Record that the stand-in received
   `--init` and `./...`.
6. MUST record the real lane run when the Run sandbox can make it, and
   otherwise record the row as declared unreachable:
   - When `docker info` exits `0` and `docker image inspect golang:1.26`
     finds the image, run `bin/roundfix deliver lane` in the audited head's
     worktree with `HOME` set to a temporary directory and `GOMODCACHE` set
     to the value `go env GOMODCACHE` printed before `HOME` changed. Record
     its outcome line, its exit code, and from its lane log the result of
     `TestCLITestHelpersEndWithTheirTestBinary` and of
     `TestGuidanceCompositionJourney/rust-cli`. Record environment variable
     names only, never their values.
   - Otherwise record `blocked (environment: Docker unreachable from the Run
     sandbox)` and count it in `rows_blocked_declared` against the PRD's
     Unreachable Acceptance.
7. MUST execute the tests task_01, task_02, task_03 and task_04 name in their
   Verification.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read `.github/workflows/ci-verify.yml` at the audited head and record its
     `runs-on` value and that no step installs a Rust toolchain;
   - because this Spec's authorization forbids network access, record each
     web source as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`.
9. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` defines **External Test Tool**, **Lane
     Container Failure** and **Lane Change Failure** and carries the revised
     **Linux Lane**;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every Governed Path changed is bounded in `_authorization.md`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
      under `internal/baseline/assets/`, and `docs/agents/setup-context.json`;
    - the Roundfix skill version was raised and recorded;
    - no `### QA settlement` section changed;
    - no guide, agent guide or skill copy keeps `Linux Lane failed: docker run
      exited` or `Linux Lane exceeded 20m0s`.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0282-the-linux-lane-runs-like-ci --plan --offline`
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
- [ ] Reproduce the four transcripts through the built binary, and run the whole suite without the External Test Tools.
- [ ] Run the real lane when Docker is reachable, record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every lane container runs with `--init` on the fetched default's image, and each removed guard makes its named test fail.
- [ ] No test fails for an absent External Test Tool, and a new literal lookup is refused.
- [ ] Each container cause and change reason parks with its own blocker and Park Class, `roundfix deliver lane` matches its transcripts, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0282-the-linux-lane-runs-like-ci/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-5; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Invariants 1-5; Current behavior; Build Order 5
- ADR-0294; ADR-0288; ADR-0104; ADR-0240; ADR-0276
