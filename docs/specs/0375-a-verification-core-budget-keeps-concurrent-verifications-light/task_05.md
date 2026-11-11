---
task: task_05
spec: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the Verification Core
Budget: a new `internal/corebudget` package that computes half of the logical
CPUs divided by `runs.max_verifying` (or `runs.max_active` while it is `0`),
unless a positive whole-number `GOMAXPROCS` overrides it; `GOMAXPROCS` and
`ROUNDFIX_VERIFICATION_CORES` on every Verification and Agent Session of an
Implement Run, on the Delivery Queue's repository gate and in the Linux Lane
container; the `verification-cores:` line of `roundfix doctor`; the guides,
the Roundfix references, the glossary and ADR-0385.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head. No row starts an Agent
Session or an ACP adapter, calls ai-memory, reads or writes the real
`~/.roundfix` or the real home directory, runs the real `gh` or Docker,
reaches the network, reads another repository, sends a request to a
provider, the judge or a reviewer, signals a process it did not start, or
reads a credential. Every `go test` it runs carries an explicit `-timeout`.
Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh
temporary directories, with `GOMODCACHE` and `GOCACHE` kept at their
`go env` values, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY` and
every Judge Key Variable for that command, and removes the temporary
directories afterwards. The QA never writes into the checkout while a test
suite runs. The gate keeps its evidence under names the repository does not
ignore, `.txt` for command output, and links only files it kept. Evidence
records variable names and the budget numbers the tests print, never another
environment value.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter that needs the session's toolchain configuration
under `HOME`. Only rows that run `bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 4 and Invariants 1 to 4 by running
   `go test -count=1 -timeout 20m -v` with every test tasks 01 to 03 name
   (`_techspec.md` → Coverage Map) over `./internal/corebudget`,
   `./internal/daemon`, `./internal/agent` and `./internal/cli`, recording
   every row by name; then the whole `./internal/corebudget`,
   `./internal/daemon`, `./internal/agent` and `./internal/delivery`
   packages, `go test -count=1 -timeout 40m -run 'Implement|Doctor|Deliver|Gate|Lane' ./internal/cli`,
   `GOOS=windows go build -buildvcs=false ./internal/...` and `-tags
   docscontract ./internal/docscontract`, recording each pass.
3. MUST record through the built binary, in a fresh temporary `HOME`, that
   `bin/roundfix doctor --help` exits 0, and record that the
   `verification-cores:` line is verified only through
   `TestDoctorReportsTheVerificationCoreBudget` in requirement 2: the real
   `bin/roundfix doctor` proves Agent Selection Profiles by starting ACP
   adapters, which this gate never does.
4. MUST run each sabotage of `_techspec.md` → Testing Approach, 1 to 14, on
   its own temporary copy of the audited head made with `git archive`, turned
   into a Git repository with one commit before the mutation (the spawning
   packages install the suite guard, which refuses a tree outside Git), with
   `-timeout` on every `go test`, running only the named test, and record that
   the named test passes on the unmutated copy and fails on the named row
   after the mutation.
5. MUST execute the tests tasks 01 to 04 name in their Verification.
6. MUST search and record:
   - `docs/user-guide/configuration.md` carries `## Verification Core Budget`
     and the four sentences of task_02 requirement 9;
   - `docs/user-guide/commands/doctor.md` carries the `verification-cores:`
     bullet and `docs/user-guide/commands/deliver.md` names both variables;
   - the canonical `implement`, `deliver` and `setup` references and their
     mirrors under `skills/roundfix/references/` hold the same bytes and name
     the budget;
   - `CONTEXT.md` carries API Contract 7 and ADR-0385 is `accepted`;
   - no non-test Go file under `internal/` outside `internal/corebudget`
     writes the literal `GOMAXPROCS` or `ROUNDFIX_VERIFICATION_CORES`, and
     `corebudget.GoVariable` and `corebudget.Variable` are used only in
     `internal/daemon/daemon.go` (through `corebudget.Environment`),
     `internal/agent/acpx_runner.go` and `internal/cli/deliver_linux_lane.go`.
7. MUST record, as outside-evidence rows, sources this Spec did not write,
   each with its date and none blocked, through the passages `_prd.md` →
   Acceptance evidence quotes:
   - the vortex measurement of 2026-11-10 in Spec 0374's adopted copy, read
     from the committed file the PRD's Sources link names: compare "O `go
     test` usa todos os núcleos por padrão" with requirement 2's verifier and
     gate rows, where the command sees `GOMAXPROCS` equal to the budget;
   - the Go command's own help, read offline with `go help build`: compare
     "The default is GOMAXPROCS, normally the number of CPUs available." with
     the same rows, so `go test -p` follows the budget;
   - the Go runtime documentation (https://pkg.go.dev/runtime, read
     2026-11-10): compare "If the GOMAXPROCS environment variable is set to a
     positive whole number, GOMAXPROCS defaults to that value." with the rows
     of `TestAnExplicitGOMAXPROCSWinsTheVerificationCoreBudget`;
   - Jest's worker default and Vitest's `maxWorkers` page (read 2026-11-10):
     compare "In watch mode, Jest should be unobtrusive and not use all
     available CPUs." with the `default on twelve` row, where the budget is
     half of the logical CPUs.
8. MUST record the PRD's one Unreachable Acceptance entry, the doctor line
   and a Verification on the maintainer's machine after the release, as one
   row counted in `rows_blocked_declared`, with its `satisfied-by`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `adds:` for **Verification Core Budget**, that
   `CONTEXT.md` defines it with the text of API Contract 7, and that
   `TestGlossaryPhrasesAreDefinedTerms` passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`, any file under
      `internal/baseline/` or `docs/agents/`, any Vendored Skill under
      `.agents/skills`, any file under `.claude/`, `.codex/` or `.opencode/`,
      any skill other than the Roundfix skill;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0375-a-verification-core-budget-keeps-concurrent-verifications-light --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh` or Docker, call ai-memory or reach the network;
    - record any environment value other than a budget number a test prints;
    - send any request to a provider, the judge or a reviewer;
    - signal a process the row did not start, or end a process by name;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests, the packages and the Windows build.
- [ ] Run each sabotage on a temporary Git copy.
- [ ] Record the wording search, the outside sources, the unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The budget is half of the logical CPUs divided by the bound, an explicit `GOMAXPROCS` wins, and every Verification and Agent Session of an Implement Run, the repository gate and the Linux Lane receive it.
- [ ] `roundfix doctor`'s `verification-cores:` line matches each prescribed form.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files and no Baseline change.

## Verification

- `newest="$(find 'docs/specs/0375-a-verification-core-budget-keeps-concurrent-verifications-light/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-5; Core Features 1-7; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-7; Invariants 1-4; Declared breaks; Testing Approach; Build Order 5
