---
task: task_05
spec: 0176-baseline-follow-ups-and-the-incremental-tier
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and
settles the Spec on evidence. This Spec allows these production changes:

- the reconcile and restore plan builders in `internal/baseline/`;
- the Verification projection in `internal/baseline/profile_alignment.go`;
- the command help in `internal/cli/cli.go`;
- the Baseline catalog assets named in `_authorization.md`.

Every row is judged against that scope.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary, in a disposable Git repository with an
   offline source that lacks a Profile-required skill, that
   `roundfix baseline skills reconcile --format json` does all of the
   following: exits `3`, prints `plannedChanges` as `[]`, `planDigest` as
   `null` and `finding.code` `reconcile.required-removed`, and leaves
   `skills-lock.json` byte-identical. It MUST also verify the mirror operation:
   a source that only drops an unrequired skill previews with a Plan Digest and
   exit `3`, and the confirmed apply removes the entry.
3. MUST verify through the built binary that a `skills-lock.json` rewritten
   during source acquisition keeps its rewritten bytes after a confirmed apply
   of a digest previewed on the old bytes. The confirmed apply must refuse, for
   both `baseline skills reconcile` and `baseline skills restore`. Use a `git`
   script first on `PATH` that rewrites the lock on `init --bare`. MUST also
   execute the Task 02 builder tests against the built tree to show that
   `lock.changed-during-plan` refuses a changed lock and accepts the planned
   one.
4. MUST verify that the built binary's `roundfix baseline update --repo .
   --no-skills --format json` reports this repository `current`. MUST also
   verify that `docs/agents/docs-layout.md`, and the `docs/agents/docs-layout.md`
   of a disposable repository adopted with the built binary, name
   `roundfix archive <slug> --qa-override --approval <source> --reason <text>`,
   contain `Never hand-edit the override stamp`, and no longer contain
   `preserve any supplied reason`.
5. MUST verify, through the built binary in a disposable adopted repository
   whose Setup Manifest lacks `verification.incremental`, each step of the
   migration:
   - `roundfix baseline update` exits `3` naming the decision and writes
     nothing;
   - `--adopt-suggested` without a `verify-incremental` declaration exits `3`
     and writes nothing;
   - with the declaration, `--yes --adopt-suggested` applies, the rendered
     `docs/agents/agent-instructions.md` names the command, and the manifest
     records the decision and a role-`incremental` projection;
   - a second update reports `current`.

   MUST also verify that a greenfield `roundfix baseline plan` without the
   decision exits `3` naming it, and that with
   `--decision verification.incremental=<declared command>` it plans both tiers.
6. MUST replay, in a disposable Git repository and with the built binary, the
   Setup Manifest and managed guides of the Fiscus mirror at
   `/Users/marcio/dev/secondbrain/projects/fiscus/mirror/`. That manifest
   records only `verification.gate`. The replay shows that `roundfix baseline
   update` names `verification.incremental`, and that once the decision is
   answered with a command the disposable repository declares, the rendered
   agent instructions name it, so the Fiscus waiver has nothing left to waive.
   MUST record the mirror as evidence this Spec did not author. MUST record the
   row as blocked with its reason when the mirror is absent, or when its
   manifest fails to load under the current catalog for a reason this Spec does
   not own. MUST NOT change any Fiscus file.
7. MUST verify that the Roundfix and setup-context-driven skills, their mirrors,
   `docs/user-guide/commands.md`, `docs/user-guide/context-driven-development.md`
   and `CONTEXT.md` describe four things: the required-removed exit, the
   `lock.changed-during-plan` refusal, the override command, and the
   `verification.incremental` decision and its migration.
8. MUST verify that the PRD records Spec 0121's supersession by this Spec as an
   operator step at archive, and that no Task changed a Spec 0121 file.
9. MUST verify, from Git evidence, that every changed governed path is listed in
   `_authorization.md`. The parity corpus, the Source Baselines,
   `internal/baseline/assets/setups/` and `docs/references/coverage-record.json`
   must be unchanged.
10. MUST verify that this Spec's own artifacts satisfy the promise rule.
11. MUST NOT accept a row whose only evidence is that a file was read.

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

- `newest="$(find 'docs/specs/0176-baseline-follow-ups-and-the-incremental-tier/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; Core Features 1-4; Success Metrics 1-4; Acceptance
evidence; `_techspec.md` → Testing Approach 1-5; API Contracts 1-5; ADR-0058;
ADR-0066; ADR-0071; ADR-0072; ADR-0073; ADR-0080; ADR-0081; ADR-0085;
ADR-0091; ADR-0093; ADR-0096; ADR-0097; ADR-0103; ADR-0104; ADR-0117;
ADR-0118; ADR-0130; ADR-0144; ADR-0149; ADR-0154; ADR-0155; ADR-0156.
