---
task: task_05
spec: 0268-a-turn-that-did-not-finish-is-not-accepted
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the ACPX Runner reject a
turn the provider failed, lets the Fallback Chain take it under the Lost
Rollout rule, parks an unrecovered one as runtime infrastructure, refuses
Claude background tasks in Run sessions, and settles a QA turn that leaves one
with a pending report as an Unfinished Turn. Every row runs against the built
binary, the repository's tests, fake runners, a fake acpx or fake ACP adapter,
a temporary home and temporary directories. No row starts a real Agent
Session, sends a request to a model provider, reads `~/.acpx` or the real
`~/.roundfix`, or persists an environment value.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact,
   with `bin/roundfix` rebuilt from the audited head first.
2. MUST verify Success Metric 1 and API Contract 1 by executing task_02's six
   runner tests and its review classifier test, and MUST confirm from the diff
   that the exit-code mapping, ADR-0020's parsed-result anomaly and the
   read-only permission refusal paths are unchanged (Invariant 3).
3. MUST re-measure the stream premise against the installed acpx when acpx and
   Node are present: write a fake ACP adapter in a temporary directory that
   answers one `session/prompt` with `threadStatus` `active`, then
   `systemError`, then the capacity text and `end_turn`, configure it through a
   temporary `HOME`, run one `acpx` prompt with `--format json`, record
   stdout, stderr and exit, and feed the recorded stdout to the runner through
   task_02's test seam to show the Provider System Error it reports. When acpx
   or Node is absent, record the row as blocked with that reason and rest it
   on `_prd.md` → Characterization and Acceptance evidence.
4. MUST verify Success Metric 2, Success Metric 3, API Contract 2 and API
   Contract 3 by executing task_03's five daemon tests and its delivery
   workflow test, and MUST confirm from the diff that a fallback happens only
   when the Task is before its First Handoff or the QA report is pending, that
   no same-selection retry follows a Provider System Error, and that the Lost
   Rollout budget is not touched.
5. MUST verify Success Metric 4, API Contract 4 and API Contract 5 by
   executing task_04's three agent tests and three daemon tests, and MUST
   record from the built binary's prompt (through the test seam) that the QA
   prompt carries the foreground rule once.
6. MUST verify that every existing test in
   `internal/daemon/lost_rollout_recovery_test.go`,
   `internal/agent/lost_rollout_test.go` and
   `internal/cli/deliver_runtime_infrastructure_test.go` passes, and that
   those three files are byte-identical to the delivery base, the merge base
   of the audited head with `main`.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence through the passages it quotes: the ACP
   "v2 Prompt Lifecycle" RFD passages on v1 failures, the ACP "Prompt Turn"
   definition of `end_turn`, and the two Claude Code documentation passages on
   `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS`. For each it records whether the
   passage still supports the design, and it opens no network connection to do
   so.
8. MUST verify that the configuration, usage and `deliver` guides and the
   Roundfix Skill's `profiles` and `deliver` references describe the recovery,
   the two Run Event phases, the `claude` refusal and the park; that each
   skill mirror equals its canonical file; and that the raised owned-skill
   version is recorded. It repeats the old-wording search that
   `docs/agents/autonomous-work.md` requires: every guide, skill copy and
   `CONTEXT.md` searched for "the one exception" beside "Lost Rollout", each
   hit recorded, and none may still call the Lost Rollout the only exception.
9. MUST perform the glossary check of `docs/agents/domain.md`: record that
   `CONTEXT.md` defines **Provider System Error**, **Agent Background Task**
   and **Unfinished Turn**, that **Fallback Chain**, **Fallback Selection** and
   **Agent Work Started** name the Provider System Error, and whether any other
   term the Spec introduced needs one.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path changed is bounded in `_authorization.md`, and that the
    Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`,
    `docs/references/skill-coverage.json` and the `qa-gate` skill did not
    change. This row reads Task commits, so it declares a `commit_range`
    input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, supported by the pre-PR equivalent of each control: the
    maintainer's recorded delivery authority, the Daemon's repository
    Verification at the audited head, no unresolved threads because no Pull
    Request exists, no Merge-Ready acceptance yet, and the audited head named
    as the claimed candidate.
12. MUST NOT accept a row whose only evidence is that a file was read, MUST
    NOT write outside temporary directories and the Spec's `qa/` folder, and
    MUST record environment variable names only, never their values.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests and re-measure the acpx stream premise.
- [ ] Confirm the unchanged paths from the diff and the delivery base.
- [ ] Record the outside evidence, the docs, the glossary check and the scope
      audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] A Provider System Error is never accepted, and the recovery and park
      follow the Lost Rollout rule without a same-selection retry.
- [ ] `claude` commands carry the refusal variable and an Unfinished Turn names
      every task.
- [ ] The Lost Rollout tests are byte-identical to the delivery base and pass.
- [ ] The scope audit finds no undeclared Governed Path.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0268-a-turn-that-did-not-finish-is-not-accepted/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- [_prd.md](_prd.md) — Goals 1-5; User Stories 1-4; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Characterization; Acceptance evidence
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Invariants 1-15; Testing Approach; Build Order 5
- ADR-0274; ADR-0245; ADR-0114; ADR-0050; ADR-0020; ADR-0057
