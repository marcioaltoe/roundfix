# Loop-order evidence

The settled action order is identical in the three rule-owned carriers:

`implement the graph including its authored gate, archive, open the Pull
Request, watch until Clean, and merge`.

Fresh reads confirmed it in:

- `docs/agents/autonomous-work.md`;
- the shipped generated clause under the formatter fixture;
- `internal/baseline/assets/modules/autonomous-work.json`.

`rtk go test ... -run '^TestCheckLoopOrder'` exited 0 and covers an independent
seeded divergence in each of those sources plus the corrected repository.
The public repository `spec check` also exited 0.

## F-002 — the shipped rationale misstates the Spec 0078 evidence count

Impact: Trust-Damage.

Actor and step: a maintainer reads the canonical order and its recorded proof
that ADR-0080 lets Pull Request-observing acceptance pass before a Pull Request
exists.

Expected: the quantitative claim in the guide, shipped clause, Baseline module,
TechSpec, and task_01 Result matches the cited archived QA report.

Actual: all five Spec 0065 carriers say Spec 0078 passed with **nine of
eighteen** rows blocked on no open Pull Request. The archived source report's
frontmatter says `rows_blocked_environment: 11`; its Coverage says 7 pass and
11 environment-blocked; its Final verdict again says all 11 blocked rows had
equivalent evidence.

Reproduction:

1. Search the Spec 0065 carriers for `nine of eighteen`.
2. Read
   `docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/qa-report-2026-08-05.md`.
3. Compare frontmatter line 7, Coverage, and Final verdict: all state 11.

The broader conclusion remains supported — Spec 0078 did pass with typed
equivalent evidence — but the shipped numeric evidence is false. The gate
cannot call documentation honest while its proof count contradicts its source.

Affected row: R10.
