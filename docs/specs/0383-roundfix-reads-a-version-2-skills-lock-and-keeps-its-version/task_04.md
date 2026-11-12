---
task: task_04
spec: 0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version
status: completed
type: qa
complexity: low
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec makes every
`skills-lock.json` reader accept versions 1 and 2, writes and compares each
`computedHash` in the hash algorithm of the lock's own version without
changing that version, refuses any other version by name with `roundfix
upgrade` as Doctor's and the Baseline update's next action, and documents it
in the `doctor`, `skills` and `baseline` guides, the Roundfix skill's
`setup`, `baseline` and `runtime` references and the glossary.

Every row runs in this repository or in temporary directories it creates
and removes. No row creates, removes or changes any file in `/tmp`,
`/private/tmp` or `$HOME` outside a directory made by `mktemp -d`, reaches
the network, runs the skills CLI (`bunx skills`, `npx skills`) or `npm`,
reads another repository, starts an Agent session or an ACP adapter, reads
or writes the real `~/.roundfix`, sends a request to a provider, the judge or
a reviewer, records an environment value, or sends a signal to any process.
Every `go test` carries an explicit `-timeout` and keeps the session's own
`HOME` and `TMPDIR`. Every row that runs a `roundfix` binary sets `HOME` and
`TMPDIR` to fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept
at their `go env` values, and removes them afterwards. The QA never writes
into the checkout while a test suite runs. The gate keeps its evidence under
names the repository does not ignore, `.txt` for command output, and links
only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 20m -v`, the tests
   `TestSumForLockVersionMatchesThePublishedSkillsCLI` and
   `TestParseLockVersionReadsOnlyVersionsOneAndTwo` (`./internal/skillhash`),
   `TestCheckRepositoryReadsASkillsLockOfVersionTwo`,
   `TestCheckRepositoryRefusesAnUnsupportedSkillsLockVersionByName` and
   `TestCheckRepositoryRejectsMalformedLockAndUnsafeRequiredNames`
   (`./skills`), `TestOwnedSkillLockEntriesUseTheHashOfTheLockVersion`,
   `TestARetiredSkillIsProvenUneditedByTheHashOfTheLockVersion`,
   `TestSkillsRestoreWritesTheHashOfTheLockVersion`,
   `TestEveryLockReaderRefusesAnUnsupportedVersionByName` and
   `TestRecordInstalledOwnedSkillsRefusesAnUnsafeLock`
   (`./internal/baseline`), and
   `TestTheSkillsCommandsWorkOnAVersionTwoSkillsLock`,
   `TestDoctorNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion`,
   `TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion`
   and `TestDoctorKeepsReadinessWhenOwnedLockComparisonIsUnavailable`
   (`./internal/cli`), recording the pass line of each test and subtest;
   then `go test -count=1 -timeout 30m ./internal/skillhash ./skills`,
   `-run 'Lock|Restore|Retired|Reconcile|OwnedSkill'` over
   `./internal/baseline` and
   `-run 'Doctor|Skills|Lock|Retired|OwnedSkill|MissingRequired|ReleasePlan'`
   over `./internal/cli` with `-timeout 30m`, and
   `GOOS=windows go vet ./internal/skillhash ./skills ./internal/baseline ./internal/cli`.
3. MUST record Success Metric 1 and the version refusal against the starting
   behavior: unpack a `git archive` of the Delivery Base (the merge base of
   the audited head with `main`) in a `mktemp -d` directory and build its
   binary there with `go build -buildvcs=false -o <dir>/roundfix
   ./cmd/roundfix`, and use the audited head's `bin/roundfix`. For each
   binary, in a fresh `mktemp -d` directory made a Git repository with `git
   init`, write the version 2 lock `{"version": 2, "skills": {"roundfix":
   {"source": "marcioaltoe/skills", "ref": "old", "sourceType": "github",
   "skillPath": "skills/roundfix/SKILL.md", "computedHash": <64 zeros>}}}`
   with two-space indentation and run `roundfix skills install --target
   project` there. Record that the Delivery Base binary exits 2 with
   `skills install failed: skills-lock.json must use version 1 with a skills
   object.` and leaves the `roundfix` entry's `ref` and zero hash, and that
   the audited head's binary exits 0, prints `Updated skills-lock.json
   computedHash: roundfix` and `Removed skills-lock.json ref: roundfix`, and
   leaves `"version": 2` with a 64-character hash and no `ref` in the entry.
   Then, with the audited head's binary only, write `{"version":3,"skills":{}}`
   in a fresh repository and record that the same command exits 2 with
   `skills-lock.json version 3 is not supported; this Roundfix reads versions
   1 and 2.`. Remove every directory afterwards.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive` and turned into a Git repository (`git init`, `git add
   -A`, one commit with a throwaway identity passed with `-c`), with
   `-timeout` on every `go test`, and record that the named test fails:
   - `sumVersion2` (`internal/skillhash/lock.go`) skips the CRLF rule:
     `TestSumForLockVersionMatchesThePublishedSkillsCLI/version_2_edge_folder`;
   - `SumForLockVersion` returns `Sum(files)` for version 2:
     `TestSumForLockVersionMatchesThePublishedSkillsCLI/version_2_one_file`;
   - `ParseLockVersion` accepts 3 instead of 2:
     `TestParseLockVersionReadsOnlyVersionsOneAndTwo`;
   - `skillsLockDocument.digest` always hashes in version 1:
     `TestOwnedSkillLockEntriesUseTheHashOfTheLockVersion/version_2`,
     `TestARetiredSkillIsProvenUneditedByTheHashOfTheLockVersion/version_2`
     and `TestSkillsRestoreWritesTheHashOfTheLockVersion`;
   - `readLocalSkillsLock` (`skills/repository.go`) refuses version 2:
     `TestCheckRepositoryReadsASkillsLockOfVersionTwo`;
   - `loadSkillsLock` refuses version 2:
     `TestTheSkillsCommandsWorkOnAVersionTwoSkillsLock`;
   - `doctorSkillReadinessResult` keeps the skills command as the next action
     for an unsupported version:
     `TestDoctorNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion`;
   - `roundfix baseline update` names `repair skills-lock.json and rerun
     roundfix baseline update` for an unsupported version:
     `TestBaselineUpdateNamesRoundfixUpgradeForAnUnsupportedSkillsLockVersion/version_3`.
5. MUST execute the tests every non-QA Task names in its Verification.
6. MUST record, as the outside-evidence row, sources this Spec did not
   write, with their dates and not blocked: the adopted vortex and pantheon
   messages (`references/`, reported by vortex and by pantheon, measured on
   Roundfix 0.141.1 and skills 1.8.1 in repositories this Spec did not
   build), whose refusal `skills-lock.json must use version 1 with a skills
   object` matches row 3's Delivery Base result; and the skills CLI README
   passages and test vectors quoted in `_prd.md` → Acceptance evidence
   (`e3b0c442…`, `93c36952…`, `b22c4d22…`), compared with the version 2
   subtests of `TestSumForLockVersionMatchesThePublishedSkillsCLI` that pin
   the same three values. The row reads the files in this repository, never
   the network.
7. MUST repeat the wording search of task_03: the three guides and the
   canonical `setup` and `baseline` references and their mirrors state that
   versions 1 and 2 are read, hashed in the version the file declares and
   kept; the `doctor` guide and both `runtime` texts name `roundfix upgrade`
   as the only next action for any other version; both `baseline` texts name
   the update's upgrade next action; and no non-test file under
   `internal/cli`, `internal/baseline` or `skills` contains `must use version
   1 with a skills object` or `want 1"`.
8. MUST perform the glossary check (`docs/agents/domain.md`): record that
   the PRD and TechSpec declare `changes: **Stale Owned Lock Entry**`, that
   `CONTEXT.md` defines it with "other than the hash of the installed copy in
   the algorithm of the lock version" and the ADR-0393 link as the Task
   commit range left it, and that ADR-0393 is accepted.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is listed in `_authorization.md` `paths`;
   - the Task commit range changes none of these: the Makefile, any
     workflow, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
     `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/**`, any ADR, any
     Baseline module, asset or derived Baseline file
     (`internal/baseline/assets/**`), any guide under `docs/agents/`, or the
     `### QA settlement` section of any skill.

   This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read another repository, or run `git fetch`, `git push` or `gh`;
    - read or write `~/.roundfix`;
    - create a file in `/tmp`, `/private/tmp` or `$HOME` outside a
      `mktemp -d` directory it removes;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - send a signal to any process;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification and the named suites.
- [ ] Characterize the Delivery Base binary against the audited head on a version 2 and a version 3 lock, and run each sabotage on a temporary copy.
- [ ] Record the outside sources, the wording search, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] On a version 2 lock the audited head's bundle install exits 0, writes the version 2 hash and keeps the version, where the Delivery Base binary exits 2.
- [ ] Each sabotage makes its named test fail, and a version 3 lock is refused by name with `roundfix upgrade`.
- [ ] The guides, references and glossary describe the versions, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; Testing Approach; Build Order 4
