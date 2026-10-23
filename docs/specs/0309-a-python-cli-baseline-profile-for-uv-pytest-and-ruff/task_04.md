---
task: task_04
spec: 0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds a Common Setup Snapshot
kind to the Baseline, the `python` common setup of `go` and `rust`, the
`python` module with its uv, ruff and pytest guide, the `python-cli` profile,
the user guide's description of it and the glossary term **Common Setup
Snapshot**.

Every row runs on temporary copies of the audited head and temporary homes
and repositories under the Run's temporary directory. No row reads or writes
the real `~/.roundfix`, reads the upstream skills checkout, reads or writes the
Secondbrain, reads another repository, runs a Python tool, reaches the
network, sends a request to a provider, the judge or a reviewer, or reads a
credential. The gate never runs `roundfix spec judge` and never runs an
archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 3 (the common kind):
   - Run `TestACommonSetupKeepsOnlyTheSkillsEveryComponentCarries`,
     `TestACommonSetupThatDiffersFromItsComponentsIsReported` and
     `TestTheAssetSyncWritesACommonSetupAfterItsComponents` in
     `./internal/baseline`, with the composed setup tests of
     `setup_composition_test.go`, and record that they pass.
   - On temporary copies of the audited head, apply one sabotage per copy and
     record that the named case fails:
     - make the common derivation keep a skill whose name every component
       carries without comparing the entries:
       `TestACommonSetupKeepsOnlyTheSkillsEveryComponentCarries` fails, the
       `differs` skill now kept;
     - make the common drift detail the composed one ("differ from the
       component union"):
       `TestACommonSetupThatDiffersFromItsComponentsIsReported/skill_dropped`
       fails;
     - make the asset sync's upstream refresh loop skip only `composed`
       snapshots: `TestTheAssetSyncWritesACommonSetupAfterItsComponents`
       fails.
3. MUST verify Success Metrics 1, 2 and 4 (the profile):
   - Run `TestThePythonSetupIsTheCommonPartOfGoAndRust`,
     `TestThePythonCliProfileRendersThePythonGuide`,
     `TestBaselineClauseForceIsCharacterized`, `TestEmbeddedCatalog`,
     `TestCatalogCompatibility` and `TestBaselinePlanCharacterization` in
     `./internal/baseline`, and `TestGuidanceCompositionJourney` and
     `TestBaselineHumanProfileAdaptation` in `./internal/cli`, and record that
     they pass.
   - In a temporary repository with a Makefile declaring `verify` and
     `verify-incremental` and one commit, run `bin/roundfix baseline plan
     --profile python-cli` with the decisions the release journey passes for
     `rust-cli` and `--format json`, with `HOME` set to a temporary directory;
     record the state, the guides it would write, the ten
     `clause.python.*` identifiers in the rendered `docs/agents/python.md`,
     and that no `docs/agents/go.md`, `docs/agents/rust.md` or
     `docs/agents/tui.md` is planned.
   - Record the `python` setup's skill names and count from the embedded
     catalog, and that none is a `golang-*`, `rust*`, `clap-rust`,
     `cut-release`, `bubbletea` or `tui-design` skill.
   - On temporary copies, one sabotage per copy, record that the named case
     fails:
     - reword "selects a new test by its pytest node id" in
       `internal/baseline/assets/modules/python.json`:
       `TestThePythonCliProfileRendersThePythonGuide` fails;
     - delete the first skill of `internal/baseline/assets/setups/python.json`:
       catalog validation fails `TestThePythonSetupIsTheCommonPartOfGoAndRust`.
   - In a temporary copy of the audited head made with `git archive`, after
     `git init` and one commit, run
     `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
     with `HOME` set to a temporary directory and record that it reports no
     file change.
4. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
5. MUST read the rendered Python guide against `_techspec.md` → API Contract
   3 and record, for each of the request's expected items 1, 2 and 3 in the
   adopted Backlog Entry, the clause that delivers it. Record every command the
   user guide paragraph and the guide name, and that `bin/roundfix --help` and
   `bin/roundfix spec check --help` list each Roundfix command among them. A
   named Roundfix command or flag the binary lacks is a finding.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - record that `internal/baseline/assets/setups/go.json` and
     `internal/baseline/assets/setups/rust.json` are byte-identical to the
     Delivery Base, count on the Delivery Base the skills they share with
     deeply equal entries, and compare that count and those names with the
     audited `python` setup and with the thirty skills the PRD quotes from
     Quaestor's decision record;
   - because this Spec's authorization forbids network access, record the uv
     and pytest documentation as `blocked (environment: network prohibited by
     Spec authorization)`, quoting the passages the PRD cites. Such a row
     never decides a qualifying `partial`.
7. MUST record the real-adopter criterion as exactly one row, `blocked
   (declared unreachable)`, against the PRD's single Unreachable Acceptance
   entry, and count it once in `rows_blocked_declared`.
8. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` defines **Common Setup Snapshot** with its
     ADR-0322 link, that the glossary contract tests pass, and whether any
     other term the Spec introduced needs an entry (`docs/agents/domain.md`);
   - record that `TestNoPythonBaselineRuntime` passes and that
     `docs/user-guide/context-driven-development.md` and the Roundfix skill
     contain neither `python3` nor `Python fallback`;
   - record that `docs/references/skill-coverage.json` carries the review of
     `_prd.md` → Skills and that `TestTheSkillCoverageMapIsCurrent` passes.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     other than `core.json` and the new `python.json`, the existing setups and
     profiles, everything under `internal/baseline/assets/source-baselines/`,
     `internal/baseline/assets/retention/` and
     `internal/baseline/assets/formatter-fixtures/`, every file under
     `docs/history/`, every skill, every user guide page other than
     `context-driven-development.md`, and every Go file outside
     `internal/baseline/setup_composition.go`,
     `internal/baseline/catalog_validate.go`,
     `internal/baseline/assets_sync.go` that is not a `_test.go` file;
   - in `core.json`, only the module version, the `root.python` hierarchy
     entry and the `guide.python` semantic owner changed;
   - no `### QA settlement` section changed.

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
    `bin/roundfix archive 0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset and with `HOME`
    set to a fresh empty temporary directory created for this row (the CLI
    loads User Config from `$HOME/.roundfix/config.yml`, so the default home
    would read `~/.roundfix`), and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the plan and the refresh in temporary copies, and each sabotage on temporary copies of the audited head.
- [ ] Read the guide against the request, and record the outside sources and the declared row.
- [ ] Record the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The common kind keeps only equal shared entries, is validated and synced, and each sabotage makes its named case fail.
- [ ] A `python-cli` plan is ready, renders the Python guide and no Go, Rust or TUI guide, and its setup is the common part of `go` and `rust`.
- [ ] The scope audit finds no undeclared Governed Path change and no change to the files row 9 names.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary; Skills
- `_techspec.md` → Current behavior; API Contracts 1-4; Invariants 1-5; Testing Approach; Build Order 4
- ADR-0322
- ADR-0204
- ADR-0191
