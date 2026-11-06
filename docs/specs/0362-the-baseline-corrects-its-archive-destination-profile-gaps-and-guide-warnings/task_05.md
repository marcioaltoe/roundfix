---
task: task_05
spec: 0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec gives the
Baseline's archive clause its current destination, refuses a repository
profile that lacks a decision a template it renders requires in `baseline
profile validate`, `baseline plan` and `baseline update`, stops the
nested-carrier warning for a guide the repository wrote under
`docs/agents/`, and documents both rules in
the baseline command guide and two owned skills.

Every row runs on this repository's files, its embedded catalog, temporary
repositories and temporary copies of the audited head or the Delivery Base.
No row starts an Agent session or an ACP adapter, reads or writes the real
`~/.roundfix`, reads another repository, reaches the network, sends a request
to a provider, the judge or a reviewer, or reads a credential. Every
`go test` carries an explicit `-timeout` and keeps the session's own `HOME`
and `TMPDIR`. Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to
fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept at their
`go env` values, and removes them afterwards. The QA never writes into the
checkout while a test suite runs. The gate keeps its evidence under names
the repository does not ignore, `.txt` for command output, and links only
files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 20m -v`, the tests tasks 01 to
   03 name over `./internal/baseline` and `./internal/cli`, then the whole
   `./internal/baseline` package, `./skills` and
   `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST reproduce Surface Transcripts 1 and 2 of `_techspec.md` through
   `bin/roundfix` in a temporary Git repository built as the TechSpec states
   (a Makefile with `verify` and `verify-incremental` targets, then
   `bin/roundfix baseline profile init --id sqllib --from go-cli-tui` and
   `verification.incremental` removed from the profile's `decisions`),
   comparing stdout, stderr and the exit code under the transcript matching
   rules, and record that `git status --porcelain` in that repository shows
   only the profile file before and after each command. Then add
   `verification.incremental` back to `decisions` and record that
   `bin/roundfix baseline profile validate sqllib` exits 0 (Success Metric 4).
4. MUST verify Success Metric 2 through the built binary on two temporary
   clones made with `git clone --shared` of this checkout into fresh
   temporary directories, removed afterwards: at the Delivery Base (the merge
   base of the audited head with `main`),
   `bin/roundfix baseline update --repo <clone> --no-skills --no-history --format text`
   lists `- update docs/agents/docs-layout.md`; at the audited head the same
   command reports `File changes: 0`. Only those lines are asserted.
5. MUST verify against the Delivery Base that in
   `internal/baseline/assets/modules/spec-workflow.json` every rule and
   clause other than `clause.spec.keep-artifacts-in-spec-folder` is equal,
   as decoded JSON, to its value at the Delivery Base, except the top-level
   `version` and the `version` of `rule.spec.docs-layout`, which each rise by
   exactly one, and that the clause's guidance equals its Delivery Base
   guidance with only the sentence of API Contract 1 replaced.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the named case:
   - in the module clause, the new ending replaced by the old one ending
     "an agent guide or the Secondbrain inbox before or at archive.":
     `TestArchiveRecordClausesAreAppended` fails;
   - `ProfileDecisionGaps` made to return nil:
     `TestAProfileWithoutTheIncrementalDecisionHasOneGap` and
     `TestBaselinePlanRefusesAProfileDecisionGap` fail;
   - the reachable set limited to the profile's own `decisions`, without the
     decisions their effects require: `TestNoBuiltInProfileHasADecisionGap`
     fails on the built-in profiles whose `autonomous.enabled` requires
     `runtime.backend` and `runtime.design`;
   - the gap check removed from `runBaselineProfileValidateCommand`, so a
     gap profile prints `sqllib: valid (repository)`:
     `TestBaselineProfileValidateNamesADecisionGap` fails;
   - `hasMarkerLine` made to return false:
     `TestARepositoryAuthoredGuideRaisesNoCarrierConflict` fails on the guide
     with a begin marker line and no end marker;
   - the managed-path condition removed from the repository-authored guide
     branch: the same test fails on the unmarked `docs/agents/docs-layout.md`;
   - the new kind removed from `warningsForCarrierClassifications`: the same
     test fails on the plain-prose `docs/agents/sql.md`.
7. MUST execute the tests every non-QA Task names in its Verification.
8. MUST record, as an outside-evidence row, a source this Spec did not
   write, with its date and not blocked: the CommonMark 0.31.2 passage
   `_prd.md` → Acceptance evidence quotes ("The block begins with a line
   that meets a start condition (after up to three optional spaces of
   indentation)."), compared offline with the `TestAMarkerLineIsACommentThatStartsALine`
   cases for three and four spaces of indentation.
9. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `None.` under `## Glossary`, that the Task
    commit range does not change `CONTEXT.md`, and that
    `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^TestGlossaryPhrasesAreDefinedTerms$'`
    passes.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `CONTEXT.md`, any file under
      `docs/adr/`, anything under
      `internal/baseline/assets/source-baselines/`,
      `internal/baseline/assets/retention/`, `internal/baseline/assets/templates/`
      or `internal/baseline/assets/decisions.json`, any module other than
      `spec-workflow.json`, or any guide under `docs/agents/` other than
      `docs-layout.md` and `setup-context.json`.

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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings --plan --offline; rm -rf "$h" "$t"`
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
- [ ] Run the named tests and reproduce both Surface Transcripts through the binary.
- [ ] Run the two update previews and the byte invariant against the Delivery Base.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the outside source, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The archive clause names the ai-memory destination in every rendered copy and is retained from the Source Baseline.
- [ ] The three profile commands refuse a gap profile with the named text and write nothing, and a guide the repository wrote stops warning while every other carrier keeps its warning.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-6; Surface Transcript 1; Surface Transcript 2; Testing Approach; Build Order 5
