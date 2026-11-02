---
task: task_05
spec: 0338-the-release-decision-travels-with-the-spec
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the PRD's Delivery
Type reach the Archive Record and the Delivery Queue's pull request title,
lets an approved authorization record carry a Release Ceiling that answers the
release plan's `approval_required` within every ceiling, makes Backlog Entries
planning evidence, rewords the Baseline's release and authorization clauses,
and teaches the authoring and Roundfix skills both rules.

Every row runs on this repository's own files, temporary repositories and
temporary copies of the audited head. No row starts an Agent session or an
ACP adapter, reads or writes the real `~/.roundfix`, reads another
repository, reaches the network, sends a request to a provider, the judge or
a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml`, and removes them afterwards. The QA never
writes into the checkout while a test suite runs. The only byte comparison is
the "no ceiling" row of requirement 4, between two runs of the same binary on
two repositories that differ only in the record; before writing it, the row
runs the same command twice on one repository and confirms the output equals
itself (the release plan prints no clock, PID or temporary path).

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 6 and Invariants 1 to 6 by running
   `go test -count=1 -timeout 20m -v` with
   `TestQueuePullRequestTitleFollowsTheDeliveryType`,
   `TestQueuePullRequestTitleKeepsTheSlugWhileItFits`,
   `TestQueueMessagesStayReadableBySquashReaders`,
   `TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug` and
   `TestReleasePlanAnswersApprovalFromTheReleaseCeiling` over
   `./internal/cli`, `TestPRDDeliveryTypeReadsTheFrontmatter` and
   `TestArchiveRecordCarriesThePRDDeliveryType` over `./internal/spec`,
   `TestPRDTypeFindingNamesTheValue` over `./internal/speccheck`,
   `TestReleaseCeilingIsReadAndValidated` over `./internal/authorization`,
   and `TestAnswerApprovalFollowsEveryCeiling` and `TestMaintenanceOnly` over
   `./internal/releaseplan`, recording each pass.
3. MUST reproduce the Argus title report through the built binary in a
   temporary clone of the audited head that the row creates and removes:
   commit a fixture Spec `9001-fixture` whose `_prd.md` declares `type: fix`,
   run `bin/roundfix spec check 9001-fixture` with fresh `HOME` and `TMPDIR`
   and record that no `SC-PRD-TYPE` appears; set `type: docs`, rerun, and
   record exactly one `SC-PRD-TYPE` naming `docs`. The fixture's other
   findings are expected and are not part of the row.
4. MUST reproduce the Argus release report through the built binary in a
   temporary repository the row creates and removes: tag `v0.1.0`; a commit
   `docs: plan 0001-x` adding `docs/specs/0001-x/_prd.md` and an approved
   `_authorization.md` with `operations` `implement`, `commit`, `push`,
   `pull_request`, `merge`, `release`, `release_ceiling: minor` and
   `release_paths: [CHANGELOG.md]`, and moving a Backlog Entry out of
   `docs/backlog/`; then `feat: deliver 0001-x (#2)` deleting the Spec
   folder, adding `docs/history/specs/0001-x.md` and changing `internal/x.go`.
   With fresh `HOME` and `TMPDIR`, record for `bin/roundfix release plan` in
   text and `--format json`: exit 0, `Decision: ready`, the proposed minor
   version, `Release authorization: answered` and the grant line. Then, each
   in its own copy, record exit 3 and `approval_required` with the reason for
   ceiling `patch` and for an added direct `feat:` commit; and, with the
   record listing `release` but no ceiling, record that the text and JSON
   equal, byte for byte, those of the same repository whose record omits
   `release`, and carry no `Release authorization` line.
5. MUST run each sabotage `_techspec.md` → Testing Approach lists on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails, with
   the case whose input makes the mutation observable.
6. MUST execute the tests task_01 to task_04 name in their Verification.
7. MUST search `docs/user-guide/`, `docs/agents/`, every copy of the
   Roundfix, `write-prd`, `write-techspec` and `write-tasks` skills, and the
   Standard TypeScript goldens, and record that: the deliver guide and both
   copies of the deliver reference describe `<type>: deliver`; the spec guide
   and references list `SC-PRD-TYPE`; the archive guide and references name
   the record's `type`; the release runbook, `usage.md`, both copies of the
   release reference and of the Roundfix `SKILL.md` state when a Release
   Ceiling answers and that it never answers a breaking plan or a manual
   classification, and still say a generic release request authorizes only a
   conclusive patch plan; the two clauses carry the API Contract 11 wording;
   and the authoring skills carry the Delivery Type and the release question.
8. MUST record, as outside-evidence rows, four sources this Spec did not
   write, each with its date and none blocked:
   - the Argus report quoted in
     `references/2026-11-01-the-delivery-pull-request-type-follows-the-spec.md`
     (2026-10-31): compare offline its shape (a fix Spec squashed as `feat:`,
     a refused patch) with requirement 2's title tests and requirement 3;
   - the Argus request quoted in
     `references/2026-11-01-a-spec-authorization-carries-the-release-decision.md`
     (2026-11-01): compare offline its shape (a release authorized with the
     Spec asked again) with requirement 4;
   - the Conventional Commits passage `_prd.md` → Acceptance evidence quotes
     ("`fix` type commits should be translated to `PATCH` releases."):
     compare it with `ClassifyCommit` of a `fix: deliver` title;
   - the Cocogitto passage `_prd.md` → Acceptance evidence quotes, with this
     repository's `cog.toml`: record that the four Delivery Types are among
     the accepted types.
9. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that the
    PRD and TechSpec declare `adds: **Delivery Type**` and `adds: **Release
    Ceiling**`, that `CONTEXT.md` defines both, that **Release Plan** gains
    its sentence, and that the Task commit range changes no other glossary
    entry.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/`,
      `docs/references/skill-coverage.json`, any existing Archive Record, any
      Baseline module other than `core` and `spec-workflow`, the Source
      Baseline corpus, any skill other than the Roundfix, `write-prd`,
      `write-techspec` and `write-tasks` files the Tasks declare;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0338-the-release-decision-travels-with-the-spec --plan --offline; rm -rf "$h" "$t"`
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
- [ ] Run the named tests and reproduce both reports through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A fix Spec's pull request title starts with `fix:` and an unsupported PRD type is reported, in the tests and through the binary.
- [ ] A sufficient Release Ceiling answers the release plan, every other case still asks with its reason, and a range without a ceiling keeps its bytes, in the tests and through the binary.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0338-the-release-decision-travels-with-the-spec/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-12; Invariants 1-6; Declared breaks; Testing Approach; Build Order 5
