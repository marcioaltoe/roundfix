---
task: task_04
spec: 0301-the-light-tier-yields-to-the-category-profile
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes a light Task whose
light runtime offers no mode for the requested access policy skip the Light
Tier with the reason code `access_mode_unavailable`, makes any other refusal of
a light session before Agent work yield to the category's Preferred Selection,
and adds a `light-tier:` Doctor line that warns about the combination.

Every row runs on temporary homes, repositories and configuration under the
Run's temporary directory, with fake runners and fake or static runtime
factories. No row reads or writes the real `~/.roundfix`, starts a real ACP
runtime or adapter, reaches the network, sends a request to OpenRouter, a
provider, the judge or a reviewer, reads a credential, or reads another
repository. The gate never runs `roundfix spec judge`, never runs `roundfix
review` against a real provider, and never runs an archive plan without
`--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2 and 3:
   - Run `TestServingSelectionsKeepsOnlyRuntimesThatServeTheRequestedAccess`
     and `TestLightTierSkipsWhenTheLightRuntimeCannotServeTheRequestedAccess`
     and record, per subtest, which selections were kept, the skip event's
     reason code, the warning line and the model of the work request
     (Invariants 1 and 2).
   - Run `TestImplementSkipsTheLightTierUnderProjectFullAccess` and record
     the work request's runtime and requested policy, the stderr warning, and
     that no key value reached the output.
3. MUST verify Success Metrics 4 and 5:
   - Run `TestLightTierYieldsToTheCategoryProfileBeforeAgentWork` and
     `TestLightTierDoesNotYieldAfterWorkOrOutsideTheLightTier` and record,
     per subtest, the outcome, the last work request's model, the
     `light_tier_yielded` events and the warning lines (Invariants 3 and 5).
   - Run `TestLightTierDispatchOrder`, `TestLightTierSkipsWithAWarning`,
     `TestLightTierRecordsSpend`, `TestLightTierEscalatesOnce`,
     `TestNoFallbackAfterAgentWorkStarted`,
     `TestAgentSelectionFallbackPublishesBeforeNextSession`, `TestTierFor`,
     `TestPlanEnabled` and `TestImplementBuildsTheLightTierPlan`, and record
     that the three existing skip reasons, the light dispatch order, the
     escalation and the no-fallback-after-work rule are unchanged
     (Invariant 4).
4. MUST verify Success Metric 6: run
   `TestDoctorLightTierResultNamesTheLightRuntimeThatCannotServeFullAccess`,
   `TestDoctorPrintsTheLightTierLineOnlyWithFullAccessAndLightModels`,
   `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn`,
   `TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts`,
   `TestRunDoctorMissingRepositoryRoot` and
   `TestDoctorReportsPrePRReviewPolicy`, and record the exact `warn` and `ok`
   lines, their position, the exit code, and that every other Doctor line is
   unchanged (Invariant 6).
5. MUST, on temporary copies of the audited head made with `git archive`,
   apply one sabotage per copy and record that the named test fails, with
   the case whose input makes the mutation observable:
   - make `applyLightTier` skip the `ServingSelections` check:
     `TestLightTierSkipsWhenTheLightRuntimeCannotServeTheRequestedAccess`
     fails in its `full-access` subtest, whose runner refuses the `opencode`
     session;
   - make `ServingSelections` refuse every selection whenever any is refused:
     `TestServingSelectionsKeepsOnlyRuntimesThatServeTheRequestedAccess`
     fails in its `second-model-serves` subtest, and
     `TestLightTierSkipsWhenTheLightRuntimeCannotServeTheRequestedAccess`
     fails in its `second-model-serves` subtest;
   - make `ServingSelections` refuse a selection whose runtime factory fails:
     `TestServingSelectionsKeepsOnlyRuntimesThatServeTheRequestedAccess`
     fails in its `factory-error` subtest;
   - move the access check after the key check:
     `TestLightTierSkipsWhenTheLightRuntimeCannotServeTheRequestedAccess`
     fails in its `full-access-key-missing` subtest, whose reason code becomes
     `key_missing`;
   - remove the yield from `activate`:
     `TestLightTierYieldsToTheCategoryProfileBeforeAgentWork` fails in its
     `factory-refused` and `setup-refused` subtests;
   - remove the yield from `Run`: the same test fails in its
     `prompt-failed-before-work` subtest;
   - drop the `AGENT_WORK_STARTED` condition of the yield:
     `TestLightTierDoesNotYieldAfterWorkOrOutsideTheLightTier` fails in its
     `after-work-started` subtest;
   - drop the light-candidate condition of the yield: the same test fails in
     its `category-refused` subtest;
   - make the yield continue at the next light model instead of the
     category's Preferred Selection:
     `TestLightTierYieldsToTheCategoryProfileBeforeAgentWork` fails on its
     prepared-session list, which then holds the second light model;
   - make `doctorLightTierResult` report the line with full access off:
     `TestDoctorPrintsTheLightTierLineOnlyWithFullAccessAndLightModels`
     fails in its `built-in` subtest;
   - make it always report `warn`:
     `TestDoctorLightTierResultNamesTheLightRuntimeThatCannotServeFullAccess`
     fails in its `serving runtime` subtest.
6. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read the adopted Backlog Entry under `references/` in the audited
     checkout, which records the Fluxus report of 2026-10-21, and record that
     its refusal text, with whitespace runs collapsed to one space,
     `requested access policy "full-access" failed predicate
     "runtime_has_access_mode" for runtime "opencode"; recovery: disable full
     access or select a runtime that supports it` is byte-identical to the
     error message the full-access case of
     `TestServingSelectionsKeepsOnlyRuntimesThatServeTheRequestedAccess`
     asserts for the refused `opencode` selection;
   - because this Spec's authorization forbids network access, record
     OpenCode's ACP documentation and the Agent Client Protocol's Session
     Modes page as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passages the PRD cites. Such a row never
     decides a qualifying `partial`.
8. MUST record Success Metric 7 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may read another
   repository or reach OpenRouter.
9. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` revises **Light Tier** and links ADR-0314, that
     the glossary contract tests pass, and whether any other term the Spec
     introduced needs an entry (`docs/agents/domain.md`);
   - record the paragraph the configuration guide gained and its four
     reasons, the bullet and row the Doctor guide gained, the sentences both
     copies of the Roundfix `runtime` reference gained, the bullet, row and
     check list both copies of the `setup` reference gained, and that
     `TestTheSkillCoverageMapIsCurrent` passes;
   - search the guides, the agent guides and the skill copies for wording
     that says a light Task is skipped for only three reasons, that a light
     session's refusal fails the Task, or that lists every Doctor line
     without `light-tier:`, and record what it finds.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every Governed Path changed is bounded in `_authorization.md`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
      under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
      every file under `docs/history/`, every file under `internal/agent`,
      `internal/lighttier/tier.go`, `internal/daemon/task_engine.go`,
      `internal/cli/implement.go`, `internal/cli/profiles_validate.go`, every
      `_test.go` file under `internal/lighttier`, `internal/daemon` and
      `internal/cli` that exists at the Delivery Base, and
      `docs/references/skill-coverage.json`;
    - the Roundfix skill's version was raised and recorded;
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
    `bin/roundfix archive 0301-the-light-tier-yields-to-the-category-profile --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, OpenRouter, the judge or a reviewer;
    - start a real ACP runtime or adapter;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Compare the Backlog Entry's refusal text with the access check's text, and record the blocked outside sources.
- [ ] Record the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A light Task under full access runs on its category profile with an `access_mode_unavailable` skip, other pre-work refusals yield, and each removed guard makes its named test fail.
- [ ] The existing Light Tier, selection fallback and Doctor behavior is unchanged.
- [ ] The scope audit finds no undeclared Governed Path change and no change to the files row 10 names.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0301-the-light-tier-yields-to-the-category-profile/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-6; Testing Approach; Build Order 4
- ADR-0314
