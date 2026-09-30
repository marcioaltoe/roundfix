---
task: task_05
spec: 0192-owned-skills-that-describe-the-product-as-it-is
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. This Spec changes no production code: it corrects twelve owned skills and five guide files and adds two repository test files. Every row is therefore a comparison of written text with what the built binary, the Daemon prompt or a Baseline clause does, or an execution of the new tests against the built tree. Each binary command runs in a disposable repository with a disposable Roundfix Home. No command reaches GitHub, a provider or a live remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's three tests against the built tree, that every command path the root help lists is named in the Roundfix skill. It MUST also run the built binary's own help for each newly documented form (`window --help`, `qa-report --help`, `baseline capabilities check --help`, `baseline profile init --help`, `init --help`) and confirm that the skill's synopsis, flags and exit codes for that form match the output.
3. MUST verify, by executing task_02's four tests against the built tree, that every command is named in the user guide, that no guide or README link is dead, and that no refused `--qa` flag is named. It MUST also run the built binary with `implement --spec <slug> --qa` in a disposable repository and record the unknown-flag refusal, and compare the guide's `spec audit` and `storage report` entries with each command's `--help`.
4. MUST prove that both contracts can fail. In a disposable copy of the candidate tree, it removes every `roundfix window clear` from the canonical Roundfix skill and records that `TestEveryCommandIsNamedInTheRoundfixSkill` fails and names `window clear`. It then adds one link to a missing path in `docs/user-guide/commands.md` and records that `TestUserGuideLinksResolve` fails and names that link. The copy is discarded.
5. MUST verify that an authorization record written as the authoring skills now teach is one the product accepts, by executing `TestDeliverPlanApprovesAnExplicitEmptyPathsGrant` and `TestDeliverStartAcceptsAnExplicitEmptyPathsGrant` against the built tree, and by confirming that the `operations` and `paths: []` wording in `write-prd`, `write-techspec` and `write-tasks` matches `docs/user-guide/commands.md`.
6. MUST verify each corrected statement in `qa-gate`, `archive-spec`, `setup-context-driven`, `implement-spec`, `council` and `write-idea` against its source of truth, and record the source beside the row:
   - the static-gate rule against the Daemon's QA prompt text in `internal/daemon/task_engine.go`;
   - the report header against the Results table a recent archived QA Report carries;
   - the archive steps against `roundfix archive --help` and the commit subject in `internal/cli/deliver_workflow.go`;
   - the Finding statuses against `docs/agents/docs-layout.md`;
   - the `implement-spec` argument hint against `roundfix implement --help`.
7. MUST verify that `write-prd`, `write-idea`, `brainstorming`, `business-analyst` and `council` each show one question with two or three options, the recommended one first and labelled `(Recommended)`, and no `Other` option, as the structured-question clause in `docs/agents/agent-instructions.md` requires, and that `council` asks its two closing questions one at a time.
8. MUST verify that the existing skill contract tests pass without an edit to any of them, that `git diff` from the Delivery Base shows no change to `skills/baseline_skill_contract_test.go`, `skills/settlement_guidance_repocontract_test.go`, `internal/docscontract/publicdocs_test.go` or any file under `internal/baseline/assets/` and `docs/agents/`, and that `make skills-sync-check` exits `0`.
9. MUST verify that each of the eleven changed skills other than the Roundfix skill declares, in both version fields, a version one patch step above the one on the Delivery Base; that the Roundfix skill keeps `0.0.2` in both fields, as the PRD's Core Feature 7 states; and that no unchanged owned skill moved its version.
10. MUST record, as evidence this Spec did not author:
    - the output of the built binary's root help, which lists the command paths the contract reads;
    - Git's own check of the same class, `t/t0450-txt-doc-vs-help.sh` (<https://github.com/git/git/blob/e9019fca/t/t0450-txt-doc-vs-help.sh>), which asserts that each builtin's documented synopsis agrees with its `-h` output.

    When the page is unavailable, it MUST record that part of the row as blocked with its reason.
11. MUST verify that this Spec's own artifacts satisfy the promise rule, and check whether it introduced, changed or retired a glossary term that `CONTEXT.md` does not carry.
12. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, and that every Governed Path is bounded in `_authorization.md`.
13. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance, and support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows the gate;
    - review-artifact ancestry: the audited head named as the claimed candidate.
14. MUST NOT accept a row whose only evidence is that a file was read.
15. MUST NOT write to the live Run Database under `~/.roundfix`.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0192-owned-skills-that-describe-the-product-as-it-is/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; Core Features 1-7; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-5; ADR-0057; ADR-0080; ADR-0088;
ADR-0091; ADR-0096; ADR-0104; ADR-0116; ADR-0155; ADR-0156; ADR-0167; ADR-0179.
