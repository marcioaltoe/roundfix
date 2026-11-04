---
task: task_04
spec: 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the bundle install of
`roundfix baseline update` and `roundfix skills install` remove the `ref` of
each `skills-lock.json` entry it rewrites for an installed Roundfix-owned
skill and report it, and makes `roundfix doctor` warn with
`DR-SKILL-LOCK-STALE` for each such entry that still pins another copy.

Every row runs on this repository's own files, temporary repositories and
temporary copies of the audited head. No row starts an Agent session or an
ACP adapter, runs the skills CLI, reads or writes the real `~/.roundfix`,
reads another repository, reaches the network, sends a request to a provider,
the judge or a reviewer, or reads a credential. Every `go test` it runs
carries an explicit `-timeout`. Every row that runs `bin/roundfix` sets
`HOME` and `TMPDIR` to fresh temporary directories, with `GOMODCACHE` and
`GOCACHE` kept at their `go env` values, because the CLI loads
`$HOME/.roundfix/config.yml`, and removes them afterwards. The QA never
writes into the checkout while a test suite runs. No row compares one
command's output with another's byte for byte. The gate keeps its evidence
under names the repository does not ignore, `.txt` for command output, and
links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 5 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` with
   `TestRecordInstalledOwnedSkillsRemovesTheRefOfAnInstalledOwnedSkill` over
   `./internal/baseline`, with
   `TestSkillsInstallRemovesTheRefOfOwnedSkillLockEntries`,
   `TestBaselineUpdateSkillsStageReportsTheUnpinnedLockEntries`,
   `TestDoctorNamesAStaleOwnedSkillLockEntry` and
   `TestDoctorKeepsAFailedSkillsLineWithAStaleOwnedSkillLockEntry` over
   `./internal/cli`, and `TestTheGuidesStateTheUnpinnedOwnedLockEntry` over
   `./skills`. Then run the whole `./internal/baseline`, `./internal/cli` and
   `./skills` packages and `-tags docscontract ./internal/docscontract`
   (`-timeout 40m` for `./internal/cli`), recording each pass.
3. MUST reproduce Surface Transcript 1 of `_techspec.md` through the built
   binary in a temporary Git repository the row creates and removes. Its
   `skills-lock.json` is written in the skills CLI's two-space format with
   entries `archive-spec`, `roundfix`, `testing-boss` and `write-prd`, each
   with `source: marcioaltoe/skills`,
   the `ref` of the catalog commit the adopted report names, `sourceType: github`, its
   catalog `skillPath` and a zero `computedHash`; only the bundle is
   installed. Run `bin/roundfix skills install` from the repository with the
   temporary `HOME` and `TMPDIR`:
   - record that stdout, stderr and the exit code match the transcript under
     its controlled variations;
   - record that `git diff` of the lock removes exactly the three owned `ref`
     lines and changes exactly their three hashes, and that the
     `testing-boss` entry is unchanged;
   - run `bin/roundfix skills install` again and record that it prints no
     `Updated skills-lock.json` or `Removed skills-lock.json` line and leaves
     the lock unchanged.
   Remove the temporary directories and the repository afterwards.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that the
   named test fails on the named case:
   - the `ref` removal dropped from `RecordInstalledOwnedSkills`:
     `TestRecordInstalledOwnedSkillsRemovesTheRefOfAnInstalledOwnedSkill`
     fails on `ref first`;
   - the first-member span ending at the member's own value end: the same
     test fails on `ref first`;
   - the `ref` removed before the installed check: the same test fails on
     `skills tool format`;
   - the scan ignoring the hash: `TestDoctorNamesAStaleOwnedSkillLockEntry`
     fails on `pinned after a bundle install`;
   - the scan reading every lock entry instead of the given names:
     `TestDoctorNamesAStaleOwnedSkillLockEntry` fails on
     `unpinned with the installed hash`;
   - the status left `ok`: `TestDoctorNamesAStaleOwnedSkillLockEntry` fails on
     `pinned after a bundle install`;
   - the next action appended without its "already there" check:
     `TestDoctorKeepsAFailedSkillsLineWithAStaleOwnedSkillLockEntry` fails;
   - the `DR-SKILL-LOCK-STALE` sentence removed from both copies of the
     Roundfix `setup` reference: `TestTheGuidesStateTheUnpinnedOwnedLockEntry`
     fails.
5. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
6. MUST search `docs/user-guide/commands/doctor.md`,
   `docs/user-guide/commands/baseline.md`, `docs/user-guide/commands/skills.md`,
   `CONTEXT.md` and both copies of the Roundfix `baseline` and `setup`
   references, and record:
   - each states the removed `ref`;
   - the baseline texts name `skills.lockUnpinned` and `- lock ref removed`;
   - the skills and setup texts name `Removed skills-lock.json ref:`;
   - the doctor and setup texts name `DR-SKILL-LOCK-STALE`;
   - no text still says the bundle install keeps every field but
     `computedHash`;
   - `CONTEXT.md` defines **Stale Owned Lock Entry**.
7. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked:
   - the marcioaltoe-skills report quoted in
     `references/2026-11-04-the-baseline-update-keeps-a-stale-ref-in-owned-skill-lock-entries.md`
     (2026-11-04): compare offline its Quaestor lock shape (twelve owned
     entries pinned to one catalog commit, only `computedHash` changed) with
     the fixture and results of requirement 3;
   - the `skills@1.7.6` passages `_prd.md` → Acceptance evidence quotes:
     compare them with the state requirement 3 leaves (owned entries without
     `ref`, so the quoted `cloneRepo` call receives no ref and clones the
     default branch) and with the external entry that keeps its `ref`.
8. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `adds: **Stale Owned Lock Entry**`, that
   `CONTEXT.md` defines it in the `**Term**:` form, and that the Task commit
   range changes no other glossary entry.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, the root
      `skills-lock.json`, `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `docs/references/skill-coverage.json`,
      `docs/user-guide/release-runbook.md`, anything under
      `internal/baseline/assets/` or `docs/agents/`, the `domain-modeling`
      skill, any skill other than the Roundfix skill;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce Surface Transcript 1 through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Each installed owned entry loses its `ref` and gets the installed hash, the external entry keeps its bytes, and a second install changes nothing, in the tests and through the binary.
- [ ] Doctor warns with `DR-SKILL-LOCK-STALE` for each stale owned entry and reads `ok` after one bundle install.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-5; Surface Transcript 1; Invariants 1-5; Declared breaks; Testing Approach; Build Order 4
