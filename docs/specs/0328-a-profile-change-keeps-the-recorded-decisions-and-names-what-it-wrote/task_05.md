---
task: task_05
spec: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes `baseline plan`
carry the Setup Manifest's recorded decisions and refuse once, makes
`baseline update` name the plan command for a profile change, names each
Root Backup an apply writes, explains a no-change reconcile, and keeps the
repository's Claude Skill Links on restore.

Every row runs on this repository's own files, temporary copies of the
audited head and temporary repositories. No row reads or writes the real
`~/.roundfix`, reads `~/dev/skills` or another repository, reaches the
network, sends a request to a provider, the judge or a reviewer, or reads a
credential. Every `go test` it runs carries an explicit `-timeout`, and every
row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary
directories, because the CLI loads `$HOME/.roundfix/config.yml`. The QA never
writes into the checkout while a test suite runs. No row compares one
command's output with another's byte for byte; where a row checks a printed
line, the Plan Digest, the backup's content digest and temporary paths are
placeholders.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 and Invariant 2 through the built binary: in
   a temporary Git repository holding a copy of the audited
   `docs/agents/setup-context.json` and a committed `AGENTS.md`, run
   `bin/roundfix baseline update --repo <repo> --decision profile=python
   --no-skills --no-history --format text` and the same with `--format
   json`; record exit `2`, a next action that contains `roundfix baseline
   plan --repo <repo> --profile python`, and `git status --porcelain` empty.
3. MUST verify Success Metric 2 through the built binary in temporary
   repositories with fake `context7`, `context7-cli` and `exa-web-search`
   skill folders: adopt `go-cli-tui` with `baseline plan` (greenfield) and
   `baseline apply`, commit, then run `baseline plan --profile python
   --decision preservation.mode=preservation --format json` and record exit
   `0` and that its `decisions` hold the recorded value of every decision of
   the manifest that `python` selects; run the same plan in text and record
   the `Carried decisions:` line; run it without `preservation.mode` and
   record exit `3` with `instruction-preservation mode is required`; in a
   fresh repository with no manifest, run `baseline plan --profile python`
   with no decision and record exit `3` and one message naming every base
   decision, `required when autonomous.enabled is true: runtime.backend,
   runtime.design` and the preservation mode.
4. MUST verify Success Metric 3 and Invariant 3 through the built binary:
   apply the `python` plan of row 3 and record one `Warning:
   baseline.preservation.backup.written: AGENTS.<digest>.md: exact copy of
   AGENTS.md before this apply; …` line whose path is the one new root file;
   repeat with `--format json` on a fresh copy and record the warning in
   `warnings`; then, in the adopted repository, run `baseline update --yes
   --no-skills --no-history --format text` and record that no
   `baseline.preservation.backup.written` line appears.
5. MUST verify Success Metrics 4 and 5 and Invariants 1, 4, 5 and 6 by
   running `go test -count=1 -timeout 20m -v` over `./internal/baseline` and
   `./internal/cli` with the test names task_02, task_03 and task_04 name in
   their Verification, and recording each pass, including the subtests of
   `TestSkillsRestoreAddsNoLinkWithoutClaudeLinks`.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - make a recorded value win over a flag in the carried-decision
     computation: `TestCarryRecordedDecisionsTakesWhatTheProfileStillSelects`
     fails on its case whose flag differs from the recorded value;
   - stop the plan command from setting `CarryRecordedDecisions`:
     `TestBaselinePlanCarriesTheRecordedDecisionsAcrossAProfileChange` fails
     with a missing-decision refusal;
   - drop the conditional-dependent part of the refusal:
     `TestBaselinePlanNamesEveryMissingAnswerInOneRefusal` fails on
     `runtime.backend`;
   - remove the profile refusal from `DeclareManifestDecisions`:
     `TestBaselineUpdateRefusesAProfileDecisionAndNamesThePlanCommand` fails
     on the next action;
   - drop the backup warnings from `ApplyPlan`:
     `TestBaselineApplyNamesTheRootBackupItWrote` fails;
   - print the `--revision` value in place of the Setup Snapshot refs:
     `TestBaselineSkillsReconcileNoChangesNamesTheSetupSnapshotRevision`
     fails on its case whose revision differs from the snapshot ref;
   - plan a link whenever `.claude/skills` is a directory:
     `TestSkillsRestoreAddsNoLinkWithoutClaudeLinks` fails on its
     copies-only subtest;
   - skip the symbolic-link postimage in the rollback:
     `TestSkillsRestoreRollsBackTheClaudeLinkWithTheTree` fails.
7. MUST execute the tests task_01 to task_04 name in their Verification.
8. MUST record, as the outside-evidence row, the published passage
   `_prd.md` → Acceptance evidence quotes from the skills CLI documentation
   ("Canonical copy: Skills are copied to `.agents/skills/`" and "Symlinks:
   Each agent's skill directory symlinks to the canonical location", Claude
   Code's project path `.claude/skills/`), and the adopter's report quoted in
   the adopted Backlog Entry under `references/`; compare offline the link a
   passing `TestSkillsRestoreLinksARestoredSkillWhereClaudeLinksAreKept`
   asserts (`.claude/skills/<name>` → `../../.agents/skills/<name>`) with
   that layout. The row is never blocked; name vercel-labs/skills and the
   message id `01a17abd` with its date as the source.
9. MUST record Success Metric 6 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    `CONTEXT.md` defines **Carried Decision**, **Root Backup** and **Claude
    Skill Link**, each citing ADR-0341 and no Spec.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      anything under `internal/baseline/assets/`, `docs/agents/` or
      `.agents/skills/` other than the canonical `roundfix` and
      `setup-context-driven` files task_01 declares;
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
    `HOME="$(mktemp -d)" TMPDIR="$(mktemp -d)" bin/roundfix archive 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's ` Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the binary rows in temporary repositories and the named tests.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the outside source, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A profile change plans with only `--profile` and `preservation.mode`, and a missing answer set is named in one refusal.
- [ ] `baseline update --decision profile=python` names the plan command and writes nothing.
- [ ] An apply names its Root Backup; a managed refresh names none.
- [ ] Reconcile names the Setup Snapshot revision; restore keeps the Claude Skill Links and adds none elsewhere.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-6; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-6; Testing Approach; Build Order 5
