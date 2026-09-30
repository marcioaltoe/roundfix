---
task: task_05
spec: 0200-a-skill-snapshot-that-matches-its-upstream
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. This Spec changes what the asset sync accepts, refreshes four setup snapshots, renames skills in three modules, adds a catalog check, moves a capability's probe, and updates this repository's upstream skills. Every behavior row is exercised by running the built binary against a temporary repository or source the gate creates, or by executing the named tests against the built tree. No command reaches a provider or a live remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary, with a temporary target and a temporary source checkout the gate creates, API Contract 1: a sync whose target names a skill only the renamed source provides refreshes (exit `0`) and its `--check` reports drift (exit `1`); a sync whose result names a skill no source provides exits `2` with "Generated setup snapshots are incompatible with the Baseline catalog"; and neither refused run changes the target's assets.
3. MUST verify, by executing task_01's tests and in a disposable copy of this repository, that the parity fixture's digests follow the synthetic source and that `make baseline-digests` reports `"changed":false`.
4. MUST verify, by executing task_02's tests and reading the four snapshots, that every snapshot pins `a4e18e4fa223196b51d0fd8224e5a33b84f97717`, that no module, bundle or setup names `context7`, `feature-systems-pattern` or `rust`, that `go-cli-tui` takes `go-tui` with `bubbletea` and `tui-design`, that every module requires what it dispatches, and that every setup keeps the `minimumVersion` values and the Roundfix-owned entry it had on the Delivery Base.
5. MUST verify API Contract 2 by executing task_03's dispatch tests, and API Contract 4 through the built binary: `roundfix baseline plan --profile go-cli-tui` in temporary repositories holding only `context7-cli`, only `context7`, and neither, where the third reports `capability.context7` blocking with a next action naming `--skill context7-cli`.
6. MUST verify API Contract 3 by reading `docs/agents/skill-dispatch.md` and the Standard TypeScript Monorepo formatter golden: the renamed triggers, `trigger.core.exa-web-search`, and no old name.
7. MUST verify API Contract 5 and Success Metric 5 in a disposable copy of this repository: `roundfix doctor` prints `skills: ok (42 required: 14 Roundfix-owned, 28 external)`, task_04's tests pass, no lock entry or directory names `context7`, and `roundfix baseline update --repo <copy> --format text` reports `current`.
8. MUST verify through the built binary, in a disposable copy of this repository, that `roundfix baseline update --repo <copy> --no-skills --format text` reports `File changes: 0`, and that the Source Baseline corpus, `internal/baseline/assets/retention/` and `internal/baseline/assets/lock-hash-compatibility-v1.json` are byte-identical to the Delivery Base.
9. MUST verify from the Delivery Base diff that no clause changed its identifier, guidance or enforcement level, that no Roundfix-owned skill changed, and that every changed `.agents/skills/` file is one the grant lists.
10. MUST record, as evidence this Spec did not author:
    - the upstream skill lists in the Secondbrain mirror, read only: `~/dev/secondbrain/projects/skills/mirror/setups/go-tui.txt` (lists `bubbletea` and `tui-design`), `go-cli.txt` (does not), and the renamed paths under `projects/skills/mirror/skills/`;
    - the Secondbrain inbox entry `~/dev/secondbrain/inbox/skills/2026-09-30-skills-vendorizadas-contradizem-o-baseline-ou-estao-quebradas.md`, which lists the three renames;
    - the `skills` tool's update page (<https://vercel-labs-skills.mintlify.app/commands/update>), which removes a skill gone from its source from the lock file.
    When a mirror or the page is unavailable, or no longer says so, it MUST record that fact as the row's evidence or as a blocked row with its reason. The gate MUST NOT write to the Secondbrain.
11. MUST verify that this Spec's own artifacts satisfy the promise rule, and check whether it introduced, changed or retired a glossary term that `CONTEXT.md` does not carry.
12. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, that every Governed Path is bounded in `_authorization.md`, and that every regenerated file is an output of a sanctioned command.
13. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance, and support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review precedes archive and publication;
    - review-artifact ancestry: the audited head named as the claimed candidate.
14. MUST NOT accept a row whose only evidence is that a file was read, except the outside-evidence rows, whose subject is what an upstream file holds.
15. MUST NOT write to the live Run Database under `~/.roundfix`, to `~/dev/skills`, or to any adopter repository.

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

- `newest="$(find 'docs/specs/0200-a-skill-snapshot-that-matches-its-upstream/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; User Stories 1-4; Core Features 1-5; Success Metrics
1-6; Acceptance evidence; `_techspec.md` → Testing Approach 1-5; API
Contracts 1-5; ADR-0058; ADR-0060; ADR-0067; ADR-0072; ADR-0073; ADR-0080;
ADR-0081; ADR-0088; ADR-0091; ADR-0093; ADR-0094; ADR-0096; ADR-0099;
ADR-0103; ADR-0104; ADR-0117; ADR-0130; ADR-0149; ADR-0155; ADR-0156;
ADR-0166; ADR-0178; ADR-0179; ADR-0191.

## Result
