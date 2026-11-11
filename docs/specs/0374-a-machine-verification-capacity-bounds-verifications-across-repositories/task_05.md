---
task: task_05
spec: 0374-a-machine-verification-capacity-bounds-verifications-across-repositories
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the Machine
Verification Capacity: the User Config key `runs.max_verifying`, `1` by
default; entry files under `verification-capacity/` in the Roundfix Home,
served in ticket order and freed only by proof that their owner is gone; a
machine unit that every Daemon Verification takes after its Run's own and
that the Delivery Queue's repository gate takes too; a Run Budget that pauses
while a Verification waits; the capacity block of `runs list` and `deliver
status` and the `verification-capacity:` line of `doctor`; the guides, the
Roundfix references, the glossary, ADR-0384 and ADR-0056's amendment notice.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head. No row starts an Agent
session or an ACP adapter, calls ai-memory, reads or writes the real
`~/.roundfix` or the real home directory, runs the real `gh` or Docker,
reaches the network, reads another repository, sends a request to a
provider, the judge or a reviewer, signals a process it did not start, or
reads a credential. Every `go test` it runs carries an explicit `-timeout`.
Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh
temporary directories, with `GOMODCACHE` and `GOCACHE` kept at their
`go env` values, puts a fake `gh` script that prints `[]` first on `PATH`,
unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY` and every Judge Key
Variable for that command, and removes the temporary directories afterwards.
The QA never writes into the checkout while a test suite runs. The gate keeps
its evidence under names the repository does not ignore, `.txt` for command
output, and links only files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter that needs the session's toolchain configuration
under `HOME`. Only rows that run `bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 6 and Invariants 1 to 6 by running
   `go test -count=1 -timeout 20m -v` with every test tasks 01 to 03 name
   (`_techspec.md` → Coverage Map) over `./internal/store`,
   `./internal/config`, `./internal/daemon` and `./internal/cli`, recording
   every row by name; then the whole `./internal/store`, `./internal/config`,
   `./internal/daemon` and `./internal/delivery` packages, `go test -count=1
   -timeout 40m -run 'Implement|Doctor|RunsList|Deliver|Gate' ./internal/cli`,
   `GOOS=windows go vet ./internal/store` and `-tags docscontract
   ./internal/docscontract`, recording each pass.
3. MUST reproduce Surface Transcript 1 of `_techspec.md` through the built
   binary: in a temporary `HOME` with User Config `runs:\n  max_verifying: 1\n`
   and no Run Database, write the two entry files of the transcript in the
   format of `_techspec.md` → Data Models into
   `$HOME/.roundfix/verification-capacity/`, each with `owner_pid` the PID of
   a `sleep 600` process the row starts and an empty `owner_identity`, run
   `bin/roundfix runs list --all --format plain` from a temporary Git
   repository, compare standard output, standard error and exit code with the
   transcript, then end the `sleep` process by the PID the row recorded.
   Repeat with `owner_pid` the PID of a `sleep` process the row started and
   already waited for, and record that the block disappears (the owner is
   proven gone) and that the two entry files are still on disk (the reader
   removes nothing).
4. MUST record through the built binary, each in its own temporary `HOME`:
   that `bin/roundfix init --scope user` writes `max_verifying: 1` under
   `runs:`; that a User Config with `runs:\n  max_verifying: -1\n` makes
   `bin/roundfix runs list` fail naming
   `runs.max_verifying must be greater than or equal to 0`; that a Project
   Config with `runs:\n  max_verifying: 3\n` makes it print
   `config: runs.max_verifying in Project Config is ignored; set runs.max_verifying in User Config`
   on standard error. The `doctor` line is verified only through
   `TestDoctorReportsTheMachineVerificationCapacity` in requirement 2: the
   real `bin/roundfix doctor` proves Agent Selection Profiles by starting
   ACP adapters, which this gate never does.
5. MUST run each sabotage of `_techspec.md` → Testing Approach, 1 to 11, on
   its own temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, running only the named test, and record
   that the named test fails on the named case.
6. MUST execute the tests tasks 01 to 04 name in their Verification.
7. MUST search and record:
   - `docs/user-guide/configuration.md` carries the sentences of task_02
     requirement 9 and the section `## Machine Verification Capacity`;
   - the `runs`, `deliver` and `doctor` command guides carry the block, the
     queue log line and the `verification-capacity:` line of API Contracts 7
     and 8;
   - the four canonical Roundfix references and their mirrors under
     `skills/roundfix/references/` hold the same bytes and name the Machine
     Verification Capacity;
   - `CONTEXT.md` carries API Contract 10, ADR-0056 carries the amendment
     notice of task_04 requirement 3, and ADR-0384 is `accepted`;
   - no file under `internal/` outside `internal/store/machine_verification*.go`
     removes a capacity entry, and no code path compares an entry's age with
     a limit (search `verification-capacity` and `machine_verification` in
     `internal/`, excluding test files).
8. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked, through the passages
   `_prd.md` → Acceptance evidence quotes:
   - the vortex measurement of 2026-11-10 in this Spec's `references/`
     (read from the committed file): compare "Com a máquina sem outra Run, a
     mesma árvore do Vortex passou inteira: 771 arquivos, 7.140 testes, 285
     s" with requirement 2's ticket-order and gate rows, which keep a second
     repository's Verification from running beside it;
   - Michał Górny, "One jobserver to rule them all" (read 2026-11-10):
     compare "The availability of job tokens therefore limits the total
     number of processes started." with requirement 3's block, where a fixed
     number of units bounds the Verifications of every repository.
9. MUST record the PRD's one Unreachable Acceptance entry, two live Runs of
   two real repositories after the release, as one row counted in
   `rows_blocked_declared`, with its `satisfied-by`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `adds:` for **Machine Verification Capacity**
    and `changes:` for **Verification Capacity**, **Waiting for
    Verification** and **Run Budget**, that `CONTEXT.md` defines them with
    the text of API Contract 10, and that `TestGlossaryPhrasesAreDefinedTerms`
    passes.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`,
      `docs/user-guide/run-database-lifecycle.md`, any Vendored Skill under
      `.agents/skills`, any file under `internal/baseline/`, any file under
      `.claude/`, `.codex/` or `.opencode/`, any skill other than the
      Roundfix skill;
    - the line `const schemaVersion = 23` of `internal/store/store.go` is the
      same at the audited head as at the Delivery Base (the merge base of the
      audited head with the default branch);
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0374-a-machine-verification-capacity-bounds-verifications-across-repositories --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh` or Docker, call ai-memory or reach the network;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - signal a process the row did not start, or end a process by name;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the transcript, the configuration rows and the doctor line through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Waiters are served in ticket order within the capacity, a holder is freed only by proof that its owner is gone, and the capacity at `0` writes nothing.
- [ ] Every Daemon Verification and the repository gate hold a machine unit while their commands run, a Stop Request ends a wait cleanly, and the Run Budget pauses during it.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files and no Run Database schema change.

## Verification

- `newest="$(find 'docs/specs/0374-a-machine-verification-capacity-bounds-verifications-across-repositories/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-5; Core Features 1-9; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-10; Surface Transcript 1; Invariants 1-6; Declared breaks; Testing Approach; Build Order 5
