---
task: task_04
spec: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec imports the upstream
`python` preset at `65d5f36`, the commit Spec 0317 gave every other upstream
Setup Snapshot, through the asset sync, adds the `python` module with its uv, ruff and pytest guide
and its nine dispatched skills, the `python` profile, and the user guide's
description of it.

Every row runs on temporary copies of the audited head and temporary homes
and repositories under the Run's temporary directory. No row reads or writes
the real `~/.roundfix`, reads `~/dev/skills` or another repository, reads or
writes the Secondbrain, runs a Python tool, reaches the network, sends a
request to a provider, the judge or a reviewer, or reads a credential. The
gate never runs `roundfix spec judge` and never runs an archive plan without
`--offline`. Every command that runs `bin/roundfix` sets `HOME` and `TMPDIR`
to fresh empty directories created for that row under the Run's temporary
directory (the CLI loads User Config from `$HOME/.roundfix/config.yml`); the
binary itself is built once with `make build` under the Run's environment, so
no row moves the Go module cache into a fresh home.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2 and 5 (the import):
   - Run `TestThePythonSetupIsThePublishedUpstreamPreset`,
     `TestEveryUpstreamSetupFollowsThePythonPresetRevision`,
     `TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill`,
     `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot`,
     `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract`,
     `TestRestoreSkillsAcceptsARepositoryProfile` and
     `TestSetupSnapshotDifferencesCompareARepositoryProfile` in
     `./internal/baseline`, and record that they pass.
   - Record from the embedded catalog the `python` setup's source, its 52
     skill names in order (the six `ai-memory-*` skills included, `handoff`
     absent), its 13 Roundfix-owned entries and its bundles, and
     the `source.ref` of `go`, `rust`, `typescript` and `python`.
   - On temporary copies of the audited head, apply one sabotage per copy and
     record that the named case fails:
     - replace `"name": "ai-memory-handoff"` with `"name": "handoff"` in
       the expected list of `internal/baseline/python_setup_test.go` only:
       `TestThePythonSetupIsThePublishedUpstreamPreset` fails, the embedded
       snapshot holding `ai-memory-handoff` at that position;
     - delete the first skill of `internal/baseline/assets/setups/python.json`:
       `TestThePythonSetupIsThePublishedUpstreamPreset` fails.
   - Record from the history of task_01 that `python.json` was created as the
     seed and then rewritten by the sync, that `go.json`, `rust.json`,
     `typescript.json`, `go-cli-typescript-bun.json`, `skills-lock.json` and
     every `.agents/skills/` file are byte-identical to the Delivery Base, and
     that no derived file changed outside the commands of `_techspec.md` → API Contract 1 (the
     Task's `## Result` and its commit's file list).
   - In a temporary copy of the audited head made with `git archive`, after
     `git init` and one commit, run `bin/roundfix baseline update --repo .
     --no-skills --no-history --yes --format text` and record that it reports
     no file change.
3. MUST verify Success Metrics 3 and 4 (the profile):
   - Run `TestThePythonModuleDispatchesItsSkills`,
     `TestThePythonProfileRendersThePythonGuide`,
     `TestEveryModuleRequiresEverySkillItDispatches`,
     `TestBaselineClauseForceIsCharacterized`, `TestEmbeddedCatalog`,
     `TestCatalogCompatibility` and `TestBaselinePlanCharacterization` in
     `./internal/baseline`, and `TestGuidanceCompositionJourney` and
     `TestBaselineHumanProfileAdaptation` in `./internal/cli`, and record that
     they pass.
   - In a temporary repository with a Makefile declaring `verify` and
     `verify-incremental` and one commit, run `bin/roundfix baseline plan
     --profile python` with the decisions the release journey passes for
     `rust-cli` and `--format json`; record the state, the guides it would
     write, the twelve `clause.python.*` identifiers in the rendered
     `docs/agents/python.md`, the nine Python skills and their triggers in the
     rendered `docs/agents/skill-dispatch.md`, and that no
     `docs/agents/go.md`, `docs/agents/rust.md`, `docs/agents/tui.md` or
     `docs/agents/cli.md` is planned.
   - Record the `python` profile's modules and required rules from the
     embedded catalog, and that it selects no `cli-surface` module and no
     `rule.cli.*` rule (maintainer decision of 2026-10-24).
   - On temporary copies, one sabotage per copy, record that the named case
     fails:
     - remove "Only when the repository's command-line interface already
       uses Typer: " from the `typer` trigger in
       `internal/baseline/assets/modules/python.json`:
       `TestThePythonModuleDispatchesItsSkills` fails;
     - prefix the `ty` trigger with "Only when the repository has chosen ty: ":
       `TestThePythonModuleDispatchesItsSkills` fails;
     - remove `--offline` from `clause.python.check-the-lock-offline`:
       `TestThePythonProfileRendersThePythonGuide` fails;
     - remove `--error-on-warning` from `clause.python.type-check-with-ty`:
       `TestThePythonProfileRendersThePythonGuide` fails;
     - reword "selects a new test by its pytest node id" in the same file:
       `TestThePythonProfileRendersThePythonGuide` fails;
     - add `cli-surface` after `python` in the `modules` of
       `internal/baseline/assets/profiles/python.json`:
       `TestThePythonProfileRendersThePythonGuide` fails on the planned
       `docs/agents/cli.md`;
     - remove `ty` from the module's `requiredSkills`:
       `TestEveryModuleRequiresEverySkillItDispatches` fails.
4. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
5. MUST read the rendered Python guide and dispatch guide against
   `_techspec.md` → API Contracts 3 and 4 and record, for each of the
   request's expected items 1, 2, 3 and 5 in the adopted Backlog Entry, the
   clause or dispatch line that delivers it, and that each of the four
   optional skills is dispatched only by a trigger beginning "Only when the
   repository". Record every command the user guide paragraph and the guide
   name, and that `bin/roundfix --help` and `bin/roundfix spec check --help`
   list each Roundfix command among them. A named Roundfix command or flag
   the binary lacks is a finding.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - compare the audited `python` snapshot's source and 52 names with the
     PRD's description of `setups/python.txt` at `65d5f36`, and, because this
     Spec's authorization forbids reading `~/dev/skills` from the gate and the
     network, record the upstream file itself as `blocked (environment:
     network prohibited by Spec authorization)`, quoting the PRD's
     description;
   - record Quaestor's report as `blocked (environment: another repository
     prohibited by Spec authorization)`, quoting the PRD;
   - record the uv and pytest documentation as `blocked (environment: network
     prohibited by Spec authorization)`, quoting the passages the PRD cites.
     Such a row never decides a qualifying `partial`.
7. MUST record the real-adopter criterion as exactly one row, `blocked
   (declared unreachable)`, against the PRD's single Unreachable Acceptance
   entry, and count it once in `rows_blocked_declared`.
8. MUST perform the glossary and documentation checks:
   - record that the PRD and TechSpec declare no glossary term, that the
     glossary contract tests pass, and whether any term the Spec introduced
     needs an entry (`docs/agents/domain.md`); record that `CONTEXT.md` has no
     "Common Setup Snapshot" entry;
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
     other than `core.json` and the new `python.json`, the existing profiles,
     everything under `internal/baseline/assets/source-baselines/`,
     `internal/baseline/assets/retention/` and
     `internal/baseline/assets/formatter-fixtures/`, every file under
     `docs/history/`, every owned skill (`skills/**` and the `.agents/skills/`
     copies of owned skills), every `.agents/skills/` file, `skills-lock.json`,
     every Setup Snapshot other than `python.json`, every user guide page other than
     `context-driven-development.md`, and every Go file that is not a
     `_test.go` file;
   - every parity fixture other than `v1/fixtures/asset-sync.json` and
     `v1/manifest.json` is byte-identical to the Delivery Base (ADR-0072);
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
    `bin/roundfix archive 0309-a-python-baseline-profile-for-uv-pytest-and-ruff --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset and with `HOME`
    and `TMPDIR` each set to a fresh empty temporary directory created for this
    row, and write the report's `## Outcome` as one paragraph that does not
    end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the plan and the refresh in temporary copies, and each sabotage on temporary copies of the audited head.
- [ ] Read the guides against the request, and record the outside sources and the declared row.
- [ ] Record the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The `python` snapshot is the upstream preset at `65d5f36`, every upstream snapshot follows that commit, and each import sabotage makes its named case fail.
- [ ] A `python` plan is ready, renders the Python guide and the nine dispatched skills, four of them only when chosen, and the offline lock check and ty in its gate, and no Go, Rust, TUI or CLI guide.
- [ ] The scope audit finds no undeclared Governed Path change and no change to the files row 9 names.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0309-a-python-baseline-profile-for-uv-pytest-and-ruff/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary; Skills
- `_techspec.md` → Current behavior; API Contracts 1-5; Invariants 1-5; Testing Approach; Build Order 4
- ADR-0322
- ADR-0241
- ADR-0191
- ADR-0072
