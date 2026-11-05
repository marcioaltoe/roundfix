---
task: task_04
spec: 0354-an-external-skill-is-checked-only-for-presence
status: completed
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec checks a required external
skill only for presence, a lock entry and a folder, removes the Setup
Snapshot comparison from Doctor and the managed refresh with the
`skills.differing` field, keeps the owned-skill version rule and this
repository's own lock and snapshot rule (ADR-0364), and states the rule in
the guides, the README, the Roundfix skill, the glossary and the amendment
notices of ADR-0221, ADR-0224, ADR-0241 and ADR-0264.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head. No row starts an Agent
session or an ACP adapter, reads or writes the real `~/.roundfix` or the real
home directory, reads another repository, reaches the network, sends a
request to a provider, the judge or a reviewer, or reads a credential. Every
`go test` it runs carries an explicit `-timeout`. Every row that runs
`bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary directories, with
`GOMODCACHE` and `GOCACHE` kept at their `go env` values, because the CLI
loads `$HOME/.roundfix/config.yml`, and removes them afterwards. The QA never
writes into the checkout while a test suite runs. The gate keeps its evidence
under names the repository does not ignore, `.txt` for command output, and
links only files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter (`rustfmt` through rustup) that needs the
session's toolchain configuration under `HOME`. Only rows that run
`bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 6 and Invariants 1 to 4 by running
   `go test -count=1 -timeout 20m -v` with
   `TestCheckRepositoryChecksAnExternalSkillOnlyForPresence` over `./skills`
   and `TestDoctorChecksAnExternalSkillOnlyForPresence`,
   `TestBaselineUpdateNeverRestoresAPresentExternalSkill`,
   `TestDoctorKeepsAnExternalSkillThatDiffersFromItsLockReady`,
   `TestBaselineUpdateNeverRestoresAnExternalSkillThatDiffersFromItsLock`,
   `TestBaselineUpdatePreviewNamesAMissingRequiredSkill`,
   `TestBaselineSkillsRestoreCommand` and
   `TestDoctorFailsForAnOwnedSkillOlderThanTheBundle` over `./internal/cli`,
   and `TestSkillsRestoreOfflinePreviewApplyAndIdempotence` over
   `./internal/baseline`, recording every row by name; then the whole
   `./skills`, `./internal/cli` and `./internal/baseline` packages with and
   without `-tags docscontract`,
   `go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec`
   and `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST reproduce the fluxus report through the built binary in a temporary
   copy of the audited head made with `git archive` and committed in a fresh
   `git init` repository, with a fresh temporary `HOME` and `TMPDIR`, after
   appending one line to `.agents/skills/testing-boss/SKILL.md`:
   - with the lock entry's `computedHash` left as it was, record the `skills:`
     line of `bin/roundfix doctor` (`skills: ok (...)`, with no `outdated:`
     and no `DR-SKILL-DIFFERS-FROM-SNAPSHOT`) and that `bin/roundfix baseline
     update --no-history --format text` exits `0` with no `Skills differing`
     line and no `testing-boss` under `Skills missing`;
   - with `--format json`, record that the result has no `differing` member;
   - after removing `.agents/skills/testing-boss`, record that Doctor prints
     `skills: failed (missing: testing-boss; next: bunx skills add marcioaltoe/skills@testing-boss)`
     and the preview exits `3` with `1 required skill(s) are missing:
     testing-boss`;
   - with the folder back and `.agents/skills/qa-gate/SKILL.md` declaring a
     version above the binary's in both version fields, record that Doctor's
     `skills:` line is `ok`, and with one below, that it is `failed (below
     minimum: skill "qa-gate" requires <bundle>, found <lower>; next:
     roundfix skills install --target project)`.
   Remove every temporary directory afterwards.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the named case:
   - the folder hash and its comparison with `computedHash` back in the
     external loop: `TestCheckRepositoryChecksAnExternalSkillOnlyForPresence`
     fails on `changed file`, and
     `TestDoctorKeepsAnExternalSkillThatDiffersFromItsLockReady` fails;
   - the lock entry requirement dropped: the first fails on `no lock entry`;
   - the snapshot suffix restored in `repositorySkillsCheck`, appending
     `DR-SKILL-DIFFERS-FROM-SNAPSHOT` for every present external skill:
     `TestDoctorChecksAnExternalSkillOnlyForPresence` fails on `differs from
     its Setup Snapshot`;
   - `baselineUpdateExternalDrift` returning every required external skill:
     `TestBaselineUpdateNeverRestoresAPresentExternalSkill` fails on `present
     copy that differs from its Setup Snapshot`;
   - the lock hash check removed from `checkThisRepositorySkillSet`, with one
     required skill's `computedHash` in the copy's `skills-lock.json`
     replaced by 64 zeros: `TestThisRepositoryHoldsEveryRequiredExternalSkill`
     (`-tags docscontract`) passes, and with the check in place it fails
     naming the skill.
5. MUST execute the tests tasks 01 to 03 name in their Verification.
6. MUST search the Doctor and Baseline guides, the README, both copies of the
   Roundfix `baseline`, `setup` and `runtime` references, `CONTEXT.md` and
   the four amended ADRs, and record that the first six carry the rule
   sentence of API Contract 7, both `runtime` copies say each remaining
   skill is checked only for presence, the README says its content is not
   compared, none of these names `DR-SKILL-DIFFERS-FROM-SNAPSHOT` or
   `skills.differing`, the two glossary entries carry their revised bodies,
   and ADR-0221, ADR-0224, ADR-0241 and ADR-0264 carry the amendment notice
   naming ADR-0364 and an `updated_at` of 2026-11-05.
7. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked:
   - the fluxus report quoted in
     `references/2026-11-05-only-roundfix-owned-skills-are-checked-by-digest-or-version.md`
     (2026-11-05): compare offline its Doctor line and its `skills.differing`
     entries with the results of requirement 3;
   - vercel-labs/skills issue 806 (opened 2026-03-31), through the passage
     `_prd.md` → Acceptance evidence quotes ("the stored hash cannot be
     recomputed from the installed copy"): record that the `changed file`
     and `lock entry without computedHash` rows and requirement 3 no longer
     depend on that hash.
8. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `changes:` for **Missing Required Skill** and
   **Setup Snapshot Difference**, that `CONTEXT.md` carries both revisions,
   and that `TestGlossaryPhrasesAreDefinedTerms` passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`,
      `docs/references/skill-coverage.json`, any Vendored Skill under
      `.agents/skills`, anything under `internal/baseline/assets/` or
      `docs/agents/`, any skill other than the Roundfix skill files task_03
      declares;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0354-an-external-skill-is-checked-only-for-presence --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the report through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A present external skill is ready whatever its content, is never listed under `skills.differing` or restored, and a missing one is still reported.
- [ ] An owned skill above the binary's version passes and one below fails, with no network.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0354-an-external-skill-is-checked-only-for-presence/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-7; Invariants 1-4; Declared breaks; Testing Approach; Build Order 4
