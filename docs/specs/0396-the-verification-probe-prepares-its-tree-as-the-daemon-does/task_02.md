---
task: task_02
spec: 0396-the-verification-probe-prepares-its-tree-as-the-daemon-does
status: pending
type: qa
complexity: medium
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes
`roundfix spec check --run-verification` prepare its disposable checkout with
the `worktree.copy` files and the Worktree Bootstrap, as the Daemon prepares
a Task Worktree, once before the first probed command; makes a failed
bootstrap end the check without probing; documents it in the `spec` guide,
the Roundfix skill's `spec` reference, the **Worktree Bootstrap** entry and
the ADR-0148 and ADR-0034 amendment notices; and records ADR-0406.

Every row runs on this repository's own files, temporary repositories,
temporary homes and temporary copies of the audited head. No row starts an
Agent session or an ACP adapter, calls ai-memory, reads or writes the real
`~/.roundfix` or the real home directory, runs the real `gh`, reaches the
network, reads another repository, sends a request to a provider, the judge
or a reviewer, or reads a credential. Every `go test` and `go build` it runs
carries an explicit `-timeout` where the command takes one and
`GOPROXY=off`. Every row that runs `bin/roundfix` or the delivery-base binary
sets `HOME` and `TMPDIR` to fresh temporary directories, with `GOMODCACHE`
and `GOCACHE` kept at their `go env` values, unsets `GH_TOKEN`,
`GITHUB_TOKEN`, `TYPESAFE_API_KEY` and every Judge Key Variable for that
command, puts a fake `gh` script that prints `[]` first on `PATH`, and
removes the temporary directories afterwards. The QA never writes into the
checkout while a test suite runs. The gate keeps its evidence under names the
repository does not ignore, `.txt` for command output, and links only files
it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes and repositories, and some, such as
`TestGuidanceCompositionJourney/rust-cli`, call an installed formatter that
needs the session's toolchain configuration under `HOME`. Only rows that run
a built binary use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 4 and Invariants 1 to 3 by running
   `go test -count=1 -timeout 30m -v` over `./internal/cli` with
   `TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes`,
   `TestSpecCheckStopsWhenTheWorktreeBootstrapFails`,
   `TestSpecCheckRunsTheWorktreeBootstrapOnceBeforeTheFirstProbedCommand` and
   `TestSpecCheckWithoutWorktreePreparationPrintsWhatItPrinted`, recording
   every test and subtest by name; then the `SpecCheck` tests of
   `./internal/cli` (`-run 'SpecCheck'`), the whole `./internal/worktree`,
   `./skills` and `-tags docscontract ./internal/docscontract` packages, and
   `GOOS=windows go build -buildvcs=false ./internal/cli ./internal/worktree`,
   recording each pass.
3. MUST reproduce Surface Transcript 1 and Surface Transcript 2 of
   `_techspec.md` through `bin/roundfix`, recording each command, both
   streams and the exit code:
   - build each fixture repository in a temporary directory the way
     `TestSpecCheckPreparesTheVerificationTreeAsTheDaemonDoes` and
     `TestSpecCheckStopsWhenTheWorktreeBootstrapFails` build theirs: copy
     `docs/agents`, `docs/adr`, `docs/workflow` and `docs/specs/clean` from
     `internal/speccheck/testdata/repo`, write the one-Task `_tasks.md` and
     `task_01.md` with the transcript's commands, the `.gitignore`, the
     untracked `local.env` and the `.roundfixrc.yml` the transcript names,
     and commit with a local Git identity;
   - run `bin/roundfix spec check clean --run-verification` in each, and
     Surface Transcript 2 also with `--format json`, and compare every line
     with the transcript, the `...` line matching the fixture's `Skipped:`
     block and nothing else;
   - in the Surface Transcript 1 repository, run
     `bin/roundfix spec check clean --run-verification --format json` and
     record `verification.bootstrap`.
4. MUST verify Invariant 1 and the starting defect against the delivery base:
   `git archive` the merge base of the audited head with `main` into a
   temporary directory, build `roundfix-base` there with `go build
   -buildvcs=false` and `GOPROXY=off`, and record that in the Surface
   Transcript 1 repository it reports `test -f .deps-installed` and
   `test -f local.env` as `honest` and exits `0`, and that in a third
   fixture repository without `.roundfixrc.yml`, with the commands
   `test -f output-the-task-writes` and `true`, `roundfix-base` and
   `bin/roundfix` print the same standard output, the same standard error
   and the same exit code, in text and in JSON.
5. MUST run each sabotage of `_techspec.md` → Testing Approach on its own
   temporary copy of the audited head made with `git archive` and committed
   in a fresh Git repository (the suite guard refuses a tree outside Git),
   with `-timeout` on every `go test`, and record that the named test fails
   on the named subtest: sabotages 1 to 8 as listed there.
6. MUST execute the tests task_01 names in its Verification.
7. MUST search and record:
   - the `spec` guide and both copies of the Roundfix `spec` reference carry
     Text A of API Contract 5 right after the paragraph it names, and the two
     reference copies hold the same bytes;
   - the **Worktree Bootstrap** entry of `CONTEXT.md` carries Text B of API
     Contract 6;
   - ADR-0148 carries Text C and ADR-0034 Text D of API Contract 7, each
     stays `accepted`, and ADR-0406 is `accepted`;
   - `internal/cli/spec_check.go` calls `runworktree.Provision` exactly once
     and `ProvisionItem` in `internal/worktree/worktree.go` returns
     `Provision(`.
8. MUST record, as outside-evidence rows, sources this Spec did not write,
   each with its date (read 2026-11-14) and none blocked, through the
   passages `_prd.md` → Acceptance evidence quotes:
   - the adopted indexus message, a measurement on a repository this Spec
     did not build: compare its "honest" from the check against the Daemon's
     vacuous refusal with the `roundfix-base` and `bin/roundfix` verdicts of
     requirements 3 and 4 on the Surface Transcript 1 repository;
   - the git-worktree documentation: compare the quoted passage with a
     measurement of the installed `git`, a tool this Spec did not build: in a
     temporary repository whose main worktree holds an ignored
     `node_modules/marker`, `git worktree add --detach <dir> HEAD` produces a
     checkout without `node_modules/`; record `git --version` and the result.
9. MUST record that `_prd.md` → Unreachable Acceptance declares none, so
   `rows_blocked_declared` is `0`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `- changes: **Worktree Bootstrap**`, that
    task_01 changed that entry and added no other term, and that
    `TestGlossaryPhrasesAreDefinedTerms` passes.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `docs/agents/setup-context.json`, any
      file under `docs/agents/` or `internal/baseline/`,
      `docs/references/coverage-record.json`,
      `docs/references/skill-coverage.json`, any Vendored Skill under
      `.agents/skills`, any skill other than the Roundfix skill, and any
      Roundfix reference other than `spec`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0396-the-verification-probe-prepares-its-tree-as-the-daemon-does --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent session or an ACP adapter;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests, the packages and the Windows build.
- [ ] Reproduce both Surface Transcripts and compare with the delivery-base binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A command that passes only after the Worktree Bootstrap is reported `vacuous` by `bin/roundfix` and `honest` by the delivery-base binary, and a repository without preparation prints the same report with both.
- [ ] A failed bootstrap probes nothing and prints Surface Transcript 2.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0396-the-verification-probe-prepares-its-tree-as-the-daemon-does/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-7; Surface Transcript 1; Surface Transcript 2; Invariants 1-4; Declared breaks; Testing Approach; Build Order 2
- ADR-0406, ADR-0148, ADR-0034
