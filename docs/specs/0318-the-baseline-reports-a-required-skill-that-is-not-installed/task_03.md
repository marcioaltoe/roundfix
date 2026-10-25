---
task: task_03
spec: 0318-the-baseline-reports-a-required-skill-that-is-not-installed
status: pending
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the `roundfix baseline
update` preview list each Missing Required Skill under `skills.missing` with
its restore command and report `plan_ready`, exit `3`, instead of `current`;
makes `roundfix skills check` fail on the same skills inside an adopted
repository; pins `roundfix doctor`'s existing failure for both removal shapes;
and documents the result in the guides, the Roundfix skill references and
`CONTEXT.md`.

Every row runs on this repository's own files, temporary copies of the
audited head and temporary repositories. No row reads or writes the real
`~/.roundfix`, starts an Agent session, reaches the network, sends a request
to a provider, the judge or a reviewer, reads a credential, or reads another
repository. Every `go test` it runs carries an explicit `-timeout`, and every
row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary
directories, because the CLI loads `$HOME/.roundfix/config.yml`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2, 4 and 5 and Invariants 1-4 by running
   `TestBaselineUpdatePreviewNamesAMissingRequiredSkill`,
   `TestDoctorNamesAMissingRequiredSkill`,
   `TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState`,
   `TestBaselineUpdatePreviewReportsAnOlderOwnedSkill` and
   `TestBaselineUpdatePreviewStaysCurrentWhenOwnedSkillsMatch` in
   `./internal/cli` with `-v`, recording each named subtest's pass.
3. MUST verify Success Metric 3 by running
   `TestSkillsCheckNamesAMissingRequiredSkill`, `TestRunSkillsCheck` and
   `TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness` with `-v`,
   recording each named subtest's pass.
4. MUST exercise the built binary end to end, with `HOME` and `TMPDIR` at
   fresh temporary directories and the Git identity set only in that `HOME`:
   - in a fresh temporary Git repository holding a `Makefile` with `verify`
     and `verify-incremental` targets and stub `.agents/skills/context7` and
     `.agents/skills/exa-web-search` folders, adopt `go-cli-tui` with
     `bin/roundfix baseline plan` (the decisions of
     `baselineApplyTestPlan` in `internal/cli/baseline_apply_test.go`) and
     `bin/roundfix baseline apply --confirm-plan <digest>`, run
     `bin/roundfix skills install`, and give every external skill the
     adopted modules require a stub folder and a `skills-lock.json` entry
     whose `computedHash` matches it, computed by a throwaway `main` package
     placed in an underscore-prefixed directory of a `git archive` copy of the
     audited head (never in the checkout) that calls
     `skills.SkillFolderHash`; read the required set from the adopted
     modules' `requiredSkills` minus the Roundfix-owned names;
   - record that `bin/roundfix baseline update --repo . --no-history --format
     json` exits `0` with `state` `current`, `bin/roundfix doctor` prints
     `skills: ok`, and `bin/roundfix skills check` exits `0`;
   - remove one required skill's folder and lock entry, as `bunx skills
     remove` does, and record the JSON preview's exit `3`, `state`,
     `skills.status` and `skills.missing`, the text preview's `Skills
     missing` lines, `bin/roundfix doctor`'s exit `1` and `skills:` line, and
     `bin/roundfix skills check`'s exit `1` and stderr;
   - restore the folder only, leaving the lock entry removed, then remove the
     folder only with its lock entry back (the stale-lock shape of the
     `skills` CLI issue #577), and record the same three outputs.
5. MUST verify Invariant 3 with the binary: in a fresh temporary directory
   that is not a Git repository, and in the repository of row 4 with every
   skill installed, `bin/roundfix skills check` prints its unchanged
   `Roundfix skill check passed: ...` line and exits `0`; and
   `bin/roundfix baseline update --repo . --no-history --no-skills --format
   json` on the damaged repository exits `0` with `state` `current`.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named subtest fails, naming the
   case whose input makes the mutation observable:
   - drop the names a lock-entry error carries (the inspection's
     `RepositoryReadinessError.MissingExternal`) from the helper:
     `TestBaselineUpdatePreviewNamesAMissingRequiredSkill/folder_and_lock_entry_removed`
     fails, because its skill has no lock entry;
   - drop `readiness.MissingExternal` from the helper:
     `.../folder_removed_with_a_stale_lock_entry` fails, because its lock
     entry remains and only the folder is gone;
   - drop `readiness.OutdatedExternal` from the helper:
     `.../folder_no_longer_matches_its_lock_entry` fails, because its folder
     exists with a different hash;
   - make the preview ignore a non-empty list when choosing its state:
     `.../folder_and_lock_entry_removed` fails on `plan_ready`;
   - fill the list even under `--no-skills`:
     `.../folder_and_lock_entry_removed` fails at its `--no-skills` check;
   - stop appending the repository diagnostics in `skills check`:
     `TestSkillsCheckNamesAMissingRequiredSkill/adopted_repository_with_a_missing_required_skill`
     fails;
   - make the repository section report a diagnostic when the Setup
     Manifest is absent:
     `TestSkillsCheckNamesAMissingRequiredSkill/repository_without_a_Setup_Manifest`
     fails.
7. MUST execute the tests task_01 and task_02 name in their Verification,
   and run `go test -timeout 15m -count=1 ./skills` and
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract`.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - because this Spec's authorization forbids network access, record the
     `skills` CLI issue #577 as `blocked (environment: network prohibited by
     Spec authorization)`, quoting the passage the PRD cites. Such a row
     never decides a qualifying `partial`;
   - read the adopted Backlog Entry under `references/` in the audited
     checkout and record that it quotes the Fluxus report (`ok:true`,
     `state:current`, `skills.differing:[]` at catalog
     `sha256:a5c687aa60ad35f3d967cb4bdc3094f404c8ff103f39447c4c5e83c8d4fa22ba`)
     and the request "verificar também que a validação detecte
     obrigatoriedade sem instalação", and that row 4's characterization
     reproduces that preview on the Delivery Base binary when one is built
     from the merge base in a temporary copy.
9. MUST record the PRD's single Unreachable Acceptance entry as exactly one
   row, `blocked (declared unreachable)`, counted once in
   `rows_blocked_declared`: no row may read Fluxus or another repository.
10. MUST perform the glossary and documentation checks:
    - record that `CONTEXT.md` holds `**Missing Required Skill**` with the
      definition of `_techspec.md` → Interfaces (`docs/agents/domain.md`);
    - search `docs/user-guide`, `.agents/skills/roundfix` and `CONTEXT.md`
      for wording that says the `baseline update` preview reports `current`
      while a required skill is missing, or that `skills check` checks only
      the shipped bundle without the adopted-repository exception, and
      record what it finds;
    - record that the Roundfix skill's two version fields are equal in the
      canonical copy and the mirror and recorded in
      `skills/testdata/owned-skill-versions.json`.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `internal/cli/doctor.go`,
      every file under `internal/baseline/` and `docs/agents/`, and every
      skill other than `roundfix`;
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
    `HOME="$(mktemp -d)" bin/roundfix archive 0318-the-baseline-reports-a-required-skill-that-is-not-installed --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - start an Agent session;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the end-to-end binary rows and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A removed or lock-mismatched required skill makes the preview exit `3` with the skill under `skills.missing`, and `skills check` and `doctor` exit `1` naming it.
- [ ] A complete set, `--no-skills` and a directory outside an adopted repository keep their current results.
- [ ] Each sabotage makes its named subtest fail, and the scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0318-the-baseline-reports-a-required-skill-that-is-not-installed/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; Interfaces; API Contracts 1-6; Invariants 1-4; Testing Approach; Build Order 3
