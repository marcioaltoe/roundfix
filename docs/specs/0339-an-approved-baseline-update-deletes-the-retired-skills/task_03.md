---
task: task_03
spec: 0339-an-approved-baseline-update-deletes-the-retired-skills
status: pending
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes an approved
`roundfix baseline update` delete each Retired Skill whose copies it proves
unedited against `skills-lock.json`, binds the deletions to the Plan Digest,
keeps an edited, unrecorded or unsafe copy with its reason, never deletes a
skill outside the Retired Skills set, and replaces ADR-0246 with ADR-0350 in
the glossary, the guide and the Roundfix skill.

Every row runs on this repository's own files, temporary repositories and
temporary copies of the audited head. No row starts an Agent session or an
ACP adapter, reads or writes the real `~/.roundfix`, reads another
repository, reaches the network, sends a request to a provider, the judge or
a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, with `GOMODCACHE` and `GOCACHE` set
to the values `go env` prints, because the CLI loads
`$HOME/.roundfix/config.yml`, and removes them afterwards. The QA never
writes into the checkout while a test suite runs. No row compares one
command's output with another's byte for byte.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 5 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` with
   `TestRetiredSkillDispositionsProveACopyUnedited`,
   `TestDeleteRetiredSkillsDeletesOnlyThePlannedCopies`,
   `TestInstalledRetiredSkillsListsTreesLinksAndLockEntries` and
   `TestContext7AndExaAreRequiredByEverySnapshotAndNeverRetired` over
   `./internal/baseline`, and
   `TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills`,
   `TestBaselineUpdateRefusesARetiredSkillEditedAfterThePlan`,
   `TestBaselineUpdateAppliesAndListsRetiredSkills`,
   `TestBaselineUpdateWithoutSkillsOmitsRetiredSkills` and
   `TestBaselineUpdateRetiredSkillDeletionTargets` over `./internal/cli`,
   recording each pass.
3. MUST reproduce the Argus case through the command path: run
   `go test -count=1 -timeout 20m -v -run '^TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills$' ./internal/cli`,
   which drives `roundfix baseline update` through `RunContext` in a
   temporary adopted Git repository, and record from its source and its pass
   that the preview reports `plan_ready`, the approved apply deletes the
   `tdd` folder, link and lock entry, and `design-patterns`, `context7-cli`,
   `exa-web-search` and `local-only` keep their folders, links and lock
   entries.
4. MUST run each sabotage `_techspec.md` → Testing Approach lists on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails, with
   the case whose input makes the mutation observable:
   - `DeleteRetiredSkills` returns before deleting:
     `TestDeleteRetiredSkillsDeletesOnlyThePlannedCopies` fails on "deletes
     the unedited retired copy and keeps every other skill", and
     `TestBaselineUpdateApprovedApplyDeletesTheRetiredSkills` fails;
   - the hash comparison never reports a difference (the comparison replaced
     by one against a constant that never matches):
     `TestRetiredSkillDispositionsProveACopyUnedited` fails on the
     `design-patterns` case;
   - `retiredSkillsPlanDigest` returns its input unchanged:
     `TestBaselineUpdateRefusesARetiredSkillEditedAfterThePlan` fails;
   - the Retired Skills name guard of `DeleteRetiredSkills` removed:
     `TestDeleteRetiredSkillsDeletesOnlyThePlannedCopies` fails on "refuses a
     name outside the Retired Skills";
   - any `.claude/skills` link target accepted:
     `TestRetiredSkillDispositionsProveACopyUnedited` fails on the
     `to-prompt` case.
5. MUST execute the tests task_01 and task_02 name in their Verification.
6. MUST search `docs/user-guide/`, `CONTEXT.md`, every copy of the Roundfix
   skill and `docs/adr/` (excluding `docs/history/` and this Spec's folder)
   and record that ADR-0246 no longer exists and no live document links it,
   that ADR-0191 carries its amendment, and that no passage still says
   `baseline update` never deletes a retired copy.
7. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked:
   - the Argus report quoted in
     `references/2026-11-01-an-approved-baseline-update-deletes-the-retired-skills.md`
     (2026-11-01): compare offline its shape (an approved apply that printed
     `Skills retired: 7` and deleted nothing) with the row of requirement 3;
   - the FreeBSD `pkg` FAQ passage `_prd.md` → Acceptance evidence quotes
     ("fails original SHA256 checksum, not removing"): compare it with the
     `edited` keep reason of API Contract 1.
8. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `changes: **Retired Skill**`, that `CONTEXT.md`
   defines it with the deletion and links ADR-0350, and that the Task commit
   range changes no other glossary entry.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `docs/references/skill-coverage.json`,
      `internal/baseline/retired_skills.go`, anything under
      `internal/baseline/assets/` or `docs/agents/`, any skill other than the
      Roundfix files task_02 declares;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0339-an-approved-baseline-update-deletes-the-retired-skills --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the Argus case through the command path.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An approved update deletes the unedited retired copy and keeps every other skill, through the command path.
- [ ] Each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files, and no live document cites ADR-0246.

## Verification

- `newest="$(find 'docs/specs/0339-an-approved-baseline-update-deletes-the-retired-skills/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-7; Invariants 1-5; Declared breaks; Testing Approach; Build Order 3
