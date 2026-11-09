---
task: task_04
spec: 0366-a-qa-report-relies-only-on-evidence-its-own-run-produced
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes every QA settlement
path, the Daemon's QA settlement, `roundfix qa-report accept` and the Settle
Command for a `qa` Task, refuse a QA Report that relies on Borrowed Evidence:
a linked evidence file the audited head or the QA Report Import already held
byte for byte, a `user_flow_executables` entry outside the report's
worktree, or a path inside another Run's Run Worktree. The Daemon first gives
the gate one repair turn and names the Run Worktree, the audited head and the
evidence directory in the QA prompt; the `qa-gate` skill writes each pass's
evidence under its own directory and lists its executables.

Every row runs on this repository's own files, temporary repositories and
worktrees, and temporary copies of the audited head. No row starts an Agent
session or an ACP adapter, reads or writes the real `~/.roundfix`, reads
another repository, reaches the network, sends a request to a provider, the
judge or a reviewer, or reads a credential. Every `go test` it runs carries
an explicit `-timeout` and keeps the session's own `HOME` and `TMPDIR`. Every
row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary
directories, with `GOMODCACHE` and `GOCACHE` kept at their `go env` values,
and removes them afterwards. The QA never writes into the checkout while a
test suite runs. No row compares one command's output with another's byte for
byte. The gate keeps its own evidence under `qa/evidence/<report name>/`,
under names the repository does not ignore, `.txt` for command output, and
links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 8 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` with every test task_01 and task_02
   name in their Verification over `./internal/spec`, `./internal/cli` and
   `./internal/daemon`, then the whole `./internal/spec`, `./internal/daemon`
   and `./skills` packages and `-tags docscontract ./internal/docscontract`,
   recording each pass.
3. MUST reproduce the Quaestor case through the built binary in a temporary
   Git repository the row creates and removes: an empty first commit; a
   worktree added with `git worktree add -b roundfix/run-old <dir>/run-old`;
   a commit holding
   `docs/specs/9001-fixture/qa/evidence/2026-11-06-gate/head.execution.txt`
   (whose text names `<dir>/run-old/.venv/bin/tool`) and
   `docs/specs/9001-fixture/qa/evidence/2026-11-06-gate/replay.py.txt`
   (`HEAD = 4bc81f4`); a later commit; and an uncommitted
   `docs/specs/9001-fixture/qa/qa-report-2026-11-09.md` with `verdict: pass`
   and a Results row linking both. Run
   `bin/roundfix qa-report accept docs/specs/9001-fixture/qa/qa-report-2026-11-09.md`
   from the repository with fresh `HOME` and `TMPDIR` and record exit 1 and
   a stderr line naming both files with `reused`. Then write
   `evidence/qa-report-2026-11-09/head.execution.txt` still naming
   `<dir>/run-old/.venv/bin/tool`, relink, rerun, and record exit 1 with
   `names another Run Worktree`. Then rewrite it to name
   `.venv/bin/tool` inside the repository, add
   `user_flow_executables: [.venv/bin/tool]` to the frontmatter, rerun, and
   record exit 0 with empty output; then remove the temporary homes and the
   repository.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the case whose input makes the mutation
   observable:
   - the `reused` comparison removed from `FindBorrowedEvidence`:
     `TestBorrowedEvidenceRefusesAnEarlierRunsEvidenceDirectory` fails on
     the Quaestor layout;
   - the audited head of a clean committed report taken as `HEAD`:
     `TestBorrowedEvidenceReadsACommittedReportAgainstItsOwnCommit` fails on
     the report committed with its evidence;
   - the Carried Row exemption removed:
     `TestBorrowedEvidenceExemptsACarriedRow` fails on the `carried (...)`
     line;
   - executable entries resolved through symbolic links before comparing:
     `TestBorrowedEvidenceNamesAnExecutableOutsideTheWorktree` fails on
     `.venv/bin/python`;
   - the own Task branch exemption removed:
     `TestBorrowedEvidenceNamesAnotherRunWorktree` fails on the
     `roundfix/run-own-task_01` worktree;
   - the check removed from `requireSettleQAReportEligibility`:
     `TestSettleRefusesAQAReportWithBorrowedEvidence` fails;
   - the Borrowed Evidence repair turn removed from `runQAGate`:
     `TestQABorrowedEvidenceRepairTurnRegeneratesTheEvidence` fails;
   - the commit-boundary check removed:
     `TestQABorrowedEvidenceAfterRepairSettlesFailed` fails;
   - `ImportCommit` and `ImportedPaths` left empty in `borrowedQAEvidence`:
     `TestQABorrowedEvidenceNamesARerunRowCitingImportedEvidence` fails.
5. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
6. MUST search `docs/user-guide/`, every copy of the Roundfix skill and
   every copy of the `qa-gate` skill and record that the qa-report guide and
   both copies of the Roundfix `settle` reference begin the rule with
   "Every evidence file the report relies on must belong to this Run." and
   name the four reasons; that the settle guide carries the settle refusal
   line; that the implement guide and both copies of the Roundfix
   `implement` reference state "On Borrowed Evidence the Daemon returns the
   QA turn once" and the refusal without a QA Report commit; that both
   copies of `qa-gate` carry `### Evidence belongs to this Run`, the
   template key `user_flow_executables: []` and no `run-slug`.
7. MUST record, as the outside-evidence row, the two sources
   `_prd.md` → Acceptance evidence names, each with its date and not
   blocked: the Quaestor report adopted in
   `references/2026-11-07-a-later-runs-qa-gate-reused-an-earlier-runs-evidence.md`
   (measured 2026-11-06), comparing offline its layout (the
   `2026-11-06-gate` directory left on the item branch, the
   `run_20261106T154834Z_fa21ea4f4f57f94a/.venv/bin/quaestor` executable,
   `HEAD = 4bc81f4`) with the row of requirement 3; and the SLSA Build
   Provenance v1.2 passage the PRD quotes ("a build that takes a git
   repository URI as a parameter might record the specific git commit that
   the URI resolved to as a dependency"), comparing it with the audited head
   the `reused` reason reads.
8. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that
   the PRD and TechSpec declare `adds: **Borrowed Evidence**` and
   `changes: **QA Report**`, that `CONTEXT.md` defines the first and carries
   the new sentence in the second, that the Task commit range changes no
   other glossary entry, and that
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^TestGlossaryPhrasesAreDefinedTerms$'`
   passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `internal/spec/qa.go`,
      `internal/spec/qa_evidence_links.go`, `internal/spec/archive.go`,
      anything under `internal/speccheck/`, `internal/agent/`,
      `internal/baseline/` or `docs/agents/`, any file under `docs/adr/`,
      or any skill other than the Roundfix and `qa-gate` files the Tasks
      declare;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0366-a-qa-report-relies-only-on-evidence-its-own-run-produced --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, cite evidence an earlier Run left, or
      keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the Quaestor case through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Borrowed Evidence is refused with its line and reason by `qa-report accept`, the Settle Command and the Daemon, in the tests and through the binary, and an honest report is accepted.
- [ ] The Daemon runs exactly one repair turn, names the Run Worktree, the audited head and the evidence directory, and never commits a refused report.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0366-a-qa-report-relies-only-on-evidence-its-own-run-produced/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-9; Invariants 1-5; Declared breaks; Testing Approach; Build Order 4
