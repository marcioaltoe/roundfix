---
task: task_04
spec: 0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the Setup Manifest's
optional `roundfixVersion`, the Recorded Roundfix Version, written by a
Baseline plan that changes a file; makes configuration loading, `roundfix init
--scope project` and the six Baseline commands refuse a newer recorded release
with a message naming both releases and `roundfix upgrade`; keeps Doctor, the
Update Check, the Run readers and `roundfix stop` working without Project
Config; documents it in four command guides, the Roundfix `baseline` and
`setup` references and the glossary; and records ADR-0392.

Every row runs on this repository's own files, temporary repositories,
temporary homes and temporary copies of the audited head. No row starts an
Agent session or an ACP adapter, calls ai-memory, reads or writes the real
`~/.roundfix` or the real home directory, runs the real `gh`, reaches the
network, reads another repository, sends a request to a provider, the judge or
a reviewer, or reads a credential. Every `go test` it runs carries an explicit
`-timeout`. Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to
fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept at their `go
env` values, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY` and every
Judge Key Variable for that command, puts a fake `gh` script that prints `[]`
first on `PATH`, and removes the temporary directories afterwards. The QA never
writes into the checkout while a test suite runs. The gate keeps its evidence
under names the repository does not ignore, `.txt` for command output, and
links only files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix` and
keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes and repositories, and some, such as
`TestGuidanceCompositionJourney/rust-cli`, call an installed formatter that
needs the session's toolchain configuration under `HOME`. Only rows that run
`bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 6 and Invariants 1 to 4 by running
   `go test -count=1 -timeout 20m -v` over `./internal/config` with the six
   tests task_01 names, over `./internal/baseline` with
   `TestAPlanRecordsItsWriterVersionOnlyWithAFileChange`,
   `TestASetupManifestRoundfixVersionMustBeAReleaseVersion`,
   `TestBaselinePlanCharacterization` and `TestBaselineCompatibilityCorpus`,
   and over `./internal/cli` with
   `TestACommandRefusesARepositoryRecordedByANewerRoundfix`,
   `TestReadOnlyCommandsWarnOfARepositoryRecordedByANewerRoundfix`,
   `TestBaselineUpdateRecordsTheRunningRoundfixVersion`,
   `TestBaselineCommandsRefuseARepositoryRecordedByANewerRoundfix`,
   `TestTheUpdateCheckNamesARepositoryRecordedByANewerRoundfix`,
   `TestUpdateCheckTranscripts` and the three parity tests the TechSpec's
   Declared breaks names (`TestBaselinePlanProfileFile`,
   `TestHumanAutomationPlanParity`, `TestBaselineHumanProfileAdaptation`),
   recording every row and subtest by name and that each parity test changed
   only by its direct request's `WriterVersion: app.Version`; then
   the whole `./internal/config`, `./internal/baseline` and `./internal/cli`
   packages and `-tags docscontract ./internal/docscontract`, recording each
   pass.
3. MUST reproduce Surface Transcript 1 of `_techspec.md` through the built
   binary in a temporary repository built as that section's fixture says (a
   Git repository from `git init`), comparing standard output, standard error
   and exit code with the transcript's conventions. Surface Transcript 2 needs
   a development installed version, which only its named test can set; through
   the binary record instead that `bin/roundfix update check` exits `0`, prints
   the warning line of API Contract 5 on stderr, and prints its release line
   followed by `Baseline: failed: <refusal>`.
   Then, in the same repository with `bin/roundfix`:
   - `spec check x`, `follow-up` and `init --scope project --force` exit `2`
     with the refusal, and `.roundfixrc.yml` keeps its bytes;
   - `baseline update --repo . --format text` prints the refusal, `Category:
     version` and `Next action: run roundfix upgrade, then rerun roundfix
     baseline update`, exits `2`, and the repository tree is unchanged;
   - `doctor` prints the `baseline-version: failed (...)` line of API Contract
     5 first and exits `1`; record only that line and the exit code;
   - `stop 00000000` prints the warning line first and then its usual
     not-found refusal.
4. MUST verify the writer end to end with two builds of the audited head,
   `-ldflags "-X roundfix/internal/app.Version=0.900.0"` and `0.901.0`, in a
   temporary copy of the audited head made with `git archive` and committed in
   a fresh Git repository: `0.900.0` `baseline update --repo . --no-skills
   --no-history --format text` reports `current` with zero file changes and no
   `roundfixVersion` in the Setup Manifest; after replacing `Fix root causes.`
   with `Fix root causes!` in `docs/agents/agent-instructions.md`, `0.900.0`
   with `--yes` writes `"roundfixVersion": "0.900.0"` and the next run
   reports `current`; `0.901.0` then reports `current` and leaves `0.900.0`;
   `bin/roundfix` (the audited release, below `0.900.0`) refuses
   `baseline update` there. Record the Setup Manifest's `roundfixVersion` line
   after each step.
5. MUST run each sabotage of `_techspec.md` → Testing Approach on its own
   temporary copy of the audited head made with `git archive` and committed
   in a fresh Git repository (the suite guard refuses a tree outside Git),
   with `-timeout`
   on every `go test`, and record that the named test fails on the named case:
   sabotages 1 to 16 as listed there.
6. MUST execute the tests tasks 01, 02 and 03 name in their Verification.
7. MUST search and record:
   - the four command guides carry task_03's headings, sentences and `text`
     blocks, and both copies of the `baseline` and `setup` references the
     paragraphs of API Contract 11;
   - `.agents/skills/roundfix` and `skills/roundfix` hold the same `baseline`
     and `setup` bytes;
   - `CONTEXT.md` carries API Contract 12, and ADR-0392 is `accepted`;
   - the refusal text of API Contract 2 appears in `internal/config` exactly
     once outside tests;
   - `docs/agents/setup-context.json` and every file under
     `internal/baseline/testdata/` keep their bytes from the Delivery Base.
8. MUST record, as outside-evidence rows, two sources this Spec did not write,
   each with its date (read 2026-11-11) and none blocked, through the passages
   `_prd.md` → Acceptance evidence quotes:
   - Go, "Go Toolchains": compare "The Go toolchain refuses to load a module or
     workspace that declares a minimum required Go version greater than the
     toolchain's own version." with the refusal of row 3;
   - HashiCorp, "Manage Terraform versions": compare "Terraform will update the
     `terraform_version` whenever you apply a change to your configuration
     using a newer Terraform version." with the writer rule of row 4.
9. MUST record that the PRD declares no Unreachable Acceptance, so
   `rows_blocked_declared` is `0`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that the
    PRD and TechSpec declare `adds:` **Recorded Roundfix Version** and
    `changes:` **Setup Manifest**, that `CONTEXT.md` defines the first with the
    text of API Contract 12 and revises the second, and that
    `TestGlossaryPhrasesAreDefinedTerms` passes.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `docs/agents/setup-context.json`, any
      file under `docs/agents/`, `internal/baseline/assets/` or
      `internal/baseline/testdata/`, `docs/references/coverage-record.json`,
      any Vendored Skill under `.agents/skills`, any skill other than the
      Roundfix skill;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - edit `roundfixVersion` in this repository's Setup Manifest or run a
      Baseline command against this checkout;
    - run the real `gh`, call ai-memory or reach the network;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests, the transcripts, the refusal sweep and the writer rows.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A newer recorded release refuses every configuration-loading and Baseline command with both releases and `roundfix upgrade`, and the recovery commands keep working.
- [ ] The record moves only with a plan that changes a file, never lowers, and no golden or this repository's Setup Manifest moves.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-12; Surface Transcript 1; Surface Transcript 2; Invariants 1-5; Declared breaks; Testing Approach; Build Order 4
