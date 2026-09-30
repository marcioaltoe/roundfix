---
task: task_05
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec changes production code in
`internal/daemon`, `internal/speccheck`, `internal/config` and `internal/cli`,
and the guides and skills that describe how a Task settles. Every behavior row
is exercised through the built binary or by executing the named tests against
the built tree. Each binary command runs in a disposable repository with a
disposable Roundfix Home. No command reaches GitHub, a provider or a live
remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's tests against the built tree:
   - a Task of a gated graph runs its declared commands and then the
     repository Verification;
   - a failure of that command returns Verification Feedback, and a final
     failure settles the Task `failed` with no Task commit;
   - a graph without a QA gate Task, and a gated graph with
     `verification.repository_at_settlement: false`, run only the declared
     commands;
   - the pre-work probe never runs the appended command;
   - the key defaults to `true` and refuses a non-boolean.
3. MUST verify, by executing task_02's tests against the built tree:
   - a prospective commit outside the grant is refused, and a bounded one and
     a sanctioned regeneration output are accepted;
   - self-approval is refused;
   - the prospective audit and the gate's audit of the created commit report
     the same findings;
   - a check that cannot be evaluated fails;
   - a final failure carries API Contract 3's reason.
4. MUST verify, by executing task_03's tests against the built tree, that a
   refusing Spec Consistency finding the Task introduced fails the check, that
   one present at Task start does not, and that a baseline error fails the
   Task before Agent work.
5. MUST verify through the built binary, in a disposable repository with a
   gated Task Graph and a fake runner, that the Run Event Stream of one non-QA
   Task carries `verification` events for the repository command,
   `settlement check: spec consistency` and `settlement check: authorization`.
6. MUST record, as evidence this Spec did not author:
   - the corrective Task of the archived Spec 0187, task_06 "The grant rule
     sits outside the shared QA settlement section", read from the repository
     history. It records that task_04 settled `completed` and that the gate
     then refused on the repository Verification;
   - the event stream and the verification log `batch-005-attempt-1.log` of
     Run `run_20260930T102010Z_231eaa2e314ed8d8`, read without writing, when
     this machine still holds them;
   - Chromium's presubmit documentation
     (<https://www.chromium.org/developers/how-tos/depottools/presubmit-scripts/>),
     for the two statements the PRD's Acceptance evidence quotes.

   When a source is unavailable, it MUST record that row as blocked with its
   reason.
7. MUST verify that `docs/user-guide/commands.md`, the Roundfix skill and its
   mirror carry `Settlement Checks` and `verification.repository_at_settlement`,
   that `docs/user-guide/configuration.md` carries the key, that the
   `write-tasks` skill and its mirror carry
   `leaves the repository Verification green`, and that
   `make skills-sync-check` exits `0`.
8. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   check whether it introduced, changed or retired a glossary term that
   `CONTEXT.md` does not carry. **Settlement Check** is the term it coined.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, and that every Governed
   Path is bounded in `_authorization.md`.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited
      head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
11. MUST NOT accept a row whose only evidence is that a file was read.
12. MUST NOT write to the live Run Database under `~/.roundfix`.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no
      declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the
      sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0190-a-task-settles-on-the-facts-its-gate-will-check/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; User Stories 1-4; Core Features 1-6; Success Metrics
1-5; Acceptance evidence; `_techspec.md` → Testing Approach 1-6; API Contracts
1-3; ADR-0014; ADR-0038; ADR-0056; ADR-0057; ADR-0080; ADR-0088; ADR-0091;
ADR-0093; ADR-0094; ADR-0096; ADR-0104; ADR-0111; ADR-0117; ADR-0130;
ADR-0135; ADR-0148; ADR-0155; ADR-0156; ADR-0159; ADR-0160; ADR-0166; ADR-0167;
ADR-0176; ADR-0178; ADR-0179; ADR-0182.

## Result
