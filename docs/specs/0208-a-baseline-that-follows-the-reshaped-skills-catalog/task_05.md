---
task: task_05
spec: 0208-a-baseline-that-follows-the-reshaped-skills-catalog
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. This Spec renames three setup snapshots and retires one, drops two skills from the catalog, gives the external-triage module eight clauses with force and Source Baseline rows, requires two skills in `core`, and updates this repository's skills. Every behavior row is exercised through the built binary against a temporary repository or a disposable copy the gate creates, or by executing the named tests against the built tree. No command reaches a provider or a network remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify API Contract 1 through the built binary: `roundfix baseline assets sync --check --source-dir <clone>/setups` against a clone of the local `~/dev/skills` at `b3c45a45f1bccd3b33aaecaaa22947d942f2fc02`, origin set to `https://github.com/marcioaltoe/skills.git`, exits `0`. When the local checkout lacks the commit, the row is recorded blocked with that reason; the gate MUST NOT clone from the network.
3. MUST verify API Contract 2 by executing task_01's tests and reading the embedded setups: exactly `go`, `rust` and `typescript` exist beside any composed snapshot another Spec added, each keeps the Roundfix-owned entry and the `minimumVersion` values the Delivery Base held for its predecessor, and no module, trigger, bundle or snapshot names a removed skill.
4. MUST verify API Contract 3 by executing task_02's tests: the golden renders the eight clauses with force, a Source Baseline adopter with external triage enabled reports `rule.external-triage` as `replaced`, and the undeclared variant is refused. It MUST also verify through the built binary, in a temporary repository adopted with `--profile standard-typescript-monorepo` and `triage.external` true, that the rendered `docs/agents/external-triage.md` holds the eight clauses.
5. MUST verify API Contract 4 by reading `docs/agents/skill-dispatch.md` and the Standard TypeScript Monorepo golden: `trigger.core.typesafe-ai` and `trigger.core.crafting-effective-readmes` are present; `trigger.core.review`, `trigger.typescript.triage` and `trigger.typescript.crafting-effective-readmes` are absent.
6. MUST verify API Contract 5 and Success Metric 5 in a disposable copy of this repository: `roundfix doctor` prints `skills: ok (43 required: 14 Roundfix-owned, 29 external)`, task_03's and task_04's tests pass, no lock entry or directory names `review` or `triage`, and `roundfix baseline update --repo <copy> --format text` reports `current`.
7. MUST verify through the built binary, in a disposable copy of this repository, that `roundfix baseline update --repo <copy> --no-skills --format text` reports `File changes: 0`, and that `internal/baseline/assets/retention/`, `internal/baseline/assets/lock-hash-compatibility-v1.json` and every parity fixture except `asset-sync.json` are byte-identical to the Delivery Base.
8. MUST verify from the Delivery Base diff that no clause other than the eight new ones changed its identifier, guidance or enforcement level, that `rule.external-triage`'s Source Baseline entry is kept, that no Roundfix-owned skill changed, and that every changed `.agents/skills/` file is one the grant lists.
9. MUST record, as evidence this Spec did not author:
    - the upstream skills repository at `b3c45a4` (row 2), whose `setups/` holds `go.txt`, `rust.txt` and `typescript.txt` and no `go-cli`, `go-tui`, `rust-cli` or `typescript-bun` list;
    - the six Secondbrain inbox entries `~/dev/secondbrain/inbox/roundfix/_triaged/2026-10-01-skills-*.md`, read only;
    - the published disclosure policies <https://github.com/jdx/hk/pull/703> and <https://github.com/openedx/.github/blob/master/AI_POLICY.md>, and the Secondbrain mirror `~/dev/secondbrain/projects/onioncry/mirror/docs/agents/triage-labels.md`.
    When a source is unavailable or no longer says so, it MUST record that fact as the row's evidence or as a blocked row with its reason. The gate MUST NOT write to the Secondbrain.
10. MUST verify that this Spec's own artifacts satisfy the promise rule, and check whether it introduced, changed or retired a glossary term that `CONTEXT.md` does not carry.
11. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, that every Governed Path is bounded in `_authorization.md`, and that every regenerated file is an output of a sanctioned command.
12. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance, and support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review precedes archive and publication;
    - review-artifact ancestry: the audited head named as the claimed candidate.
13. MUST NOT accept a row whose only evidence is that a file was read, except the outside-evidence rows, whose subject is what an outside source holds.
14. MUST NOT write to the live Run Database under `~/.roundfix`, to `~/dev/skills`, or to any adopter repository.

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

- `newest="$(find 'docs/specs/0208-a-baseline-that-follows-the-reshaped-skills-catalog/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; User Stories 1-4; Core Features 1-5; Success Metrics
1-5; Declared breaks; Acceptance evidence; `_techspec.md` → Testing Approach
1-5; API Contracts 1-5; ADR-0058; ADR-0059; ADR-0060; ADR-0061; ADR-0067;
ADR-0072; ADR-0073; ADR-0080; ADR-0081; ADR-0088; ADR-0091; ADR-0093;
ADR-0094; ADR-0096; ADR-0099; ADR-0103; ADR-0104; ADR-0117; ADR-0130;
ADR-0149; ADR-0155; ADR-0156; ADR-0166; ADR-0167; ADR-0178; ADR-0179;
ADR-0182; ADR-0186; ADR-0191; ADR-0192; ADR-0194; ADR-0195; ADR-0202;
ADR-0204; ADR-0206.
