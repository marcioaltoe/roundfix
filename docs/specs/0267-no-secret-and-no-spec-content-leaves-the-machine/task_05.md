---
task: task_05
spec: 0267-no-secret-and-no-spec-content-leaves-the-machine
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec refuses a QA Report commit and
every Roundfix push that carry a Secret Match, withholds the Judge Key
Variables from every Agent Session and Verification command, and describes
both in the skills, the QA prompt, the user guides and the glossary. Every row
runs against the built binary or the repository's own tests in temporary
repositories, homes and remotes. No row reaches a provider or the network,
starts a real Agent Session, reads a credential, or reads or writes the real
`~/.roundfix`. This gate's own evidence is held to the rule it tests: it
records environment variable names and whether they are set, never values,
records a Secret Match by path, line and kind only, and never writes a
secret-shaped sample to `qa/`, or the Daemon's own Secret Scan refuses this
report.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact,
   with `bin/roundfix` rebuilt from the audited head first.
2. MUST verify Success Metric 1: at the audited head, run
   `TestNoAgentSessionReceivesAJudgeKeyVariable` and
   `TestVerificationReceivesNoJudgeKeyVariable`, and repeat the
   characterization of `_prd.md` through a `go test -overlay` probe that writes
   nothing to the repository: `acpxCommandEnv` over a base holding the three
   Judge Key Variables returns 0 of them, and an `ExecVerifier.Verify` command
   sees none of them, with a control variable present in both. Record the
   counts and the variable names only, against 3 and 1 measured at `9963ffd3`.
3. MUST verify Success Metric 2: run
   `TestQAReportWithASecretMatchSettlesFailedWithoutACommit` and
   `TestQAReportWithoutASecretMatchKeepsItsCommit`, and
   `TestSecretPatternsMatchTheirKindAndNeverEchoTheValue` for every kind of
   `_techspec.md` → Data Models; on a temporary copy of the tree where the QA
   step's scan call is removed, the first QA step test fails. Record the
   refusal reason and event by path, line and kind.
4. MUST verify Success Metric 3: run `TestPushRefusesAHistoryThatAddsASecret`,
   `TestEnginePushParksSecretFound` and
   `TestFinalPushRefusesAHistoryThatAddsASecret`; on a temporary copy of the
   tree where `GitHubCLI.PushBranch` skips the history read, the first fails.
   Record the `secret-found` blocker, its Park Class `finding` and its next
   action as API Contract 4 states them.
5. MUST verify Success Metric 4: through a `go test -overlay` probe in
   `internal/secretscan` that writes nothing to the repository, run
   `ScanContent` over every file `git ls-files` lists at the audited head and
   record 0 Secret Matches, or each Match by path, line and kind.
6. MUST verify the QA contract text: the `qa-gate` skill's
   `### Evidence carries no secret` heading and the `qaGateContract` line of
   `internal/agent/spec_prompt.go` carry the rule of task_04 Requirements 2
   and 3, and every `### QA settlement` section is byte-identical to its bytes
   at the delivery base, the merge base of the audited head with `main`.
7. MUST verify that `bin/roundfix archive --plan`, run in a temporary
   repository holding a fixture Spec with `HOME` set to a temporary directory
   and each Judge Key Variable removed with `env -u`, as a Run's commands now
   receive their environment, reports each candidate with
   `ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set` and sends no request. The
   Daemon running this gate may predate task_02, so the gate MUST run every
   command that could call the judge, this one included, with the Judge Key
   Variables removed; record the candidate lines, never a value.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence through the passages `_prd.md` quotes: the
   GitHub push protection passage and the two GNU Make passages. For each it
   records whether the passage still supports the behavior, and it opens no
   network connection to do so.
9. MUST perform the glossary check of `docs/agents/domain.md` and record that
   `CONTEXT.md` defines **Secret Pattern**, **Secret Match**, **Secret Scan**
   and **Judge Key Variables**, and whether any other term the Spec
   introduced needs one.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path changed is bounded in `_authorization.md`, that
    `internal/cli/review.go`, the Makefile, the CI workflows, `go.mod`,
    `go.sum` and `.roundfixrc.yml` did not change against the delivery base,
    and that no added line in the Spec's commits matches a Secret Pattern.
    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, supported by the pre-PR equivalent of each control: the
    maintainer's recorded delivery authority, the Daemon's repository
    Verification at the audited head, no unresolved threads because no Pull
    Request exists, no Merge-Ready acceptance yet, and the audited head named
    as the claimed candidate.
12. MUST NOT accept a row whose only evidence is that a file was read, MUST NOT
    persist the output of `env`, `printenv`, `set`, `export -p`,
    `declare -x` or `make -p`, and MUST NOT write outside temporary
    directories and the Spec's `qa/` folder.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Measure Success Metrics 1-4 with the sabotaged copies.
- [ ] Check the QA contract text, the archive plan row and the outside evidence.
- [ ] Record the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] No Agent Session or Verification environment holds a Judge Key
      Variable, and a control variable survives.
- [ ] A QA step and both pushes refuse a sample, and each named test fails on
      its sabotaged copy.
- [ ] The repository's own tracked files hold no Secret Match.
- [ ] The scope audit finds no undeclared Governed Path, and the report holds
      no environment value.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0267-no-secret-and-no-spec-content-leaves-the-machine/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- [_prd.md](_prd.md) — Goals; User Stories 1-5; Core Features 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; Invariants 1-5; Build Order 5
- ADR-0273; ADR-0015; ADR-0258
