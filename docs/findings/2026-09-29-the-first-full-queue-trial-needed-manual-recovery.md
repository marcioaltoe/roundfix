---
status: partial
created_at: 2026-09-29
updated_at: 2026-09-29
---

# Delivery — the first full queue trial still needed manual recovery (2026-09-29)

Wave 4 (Specs 0181 and 0182) was the first cycle driven through one `roundfix deliver start`, as the 2026-09-29 handoff asked, using the v0.20.0 plan, limits, revalidation and renewed budget. The queue started at 14:01 with two retries per item. Neither item reached merge on its own: 0182 merged at 17:03 and 0181 around 18:40, both finished by the Supervisor. This report records each manual intervention and what it exposed. Specs 0181, 0182 and 0185 cover most of them.

## 1. Every Spec parked at least twice before its QA gate passed

- Symptom / evidence:
  - 0181, first gate: `SC-ADR-UNLISTED` from a fixture number in task_03's Agent-written `## Result`. The refused report then quoted it again.
  - 0181, second gate: `make verify-changed` failed on `TestSettlementGuidanceIsOneTable`, because the archive-spec skill was outside the grant.
  - 0181, third gate: `QA-AUTH-PATHS`, because the widened grant was not in the merge base of task_07's commit.
  - 0182, first gate: failed only its scope row, over two invalidated tests task_03 had not declared.
  - Evidence: the QA reports on each Run branch, and `roundfix events` for Runs `run_20260929T170240Z_b8558d4fb9103028`, `run_20260929T185005Z_5bdc71810126cc68`, `run_20260929T195318Z_43a2234f5bbc8e85`, `run_20260929T200238Z_ca4a5877c9dc9d5c` and `run_20260929T174719Z_24238e6aebd23928`.
- Root cause: the four frictions Wave 4 set out to fix hit the Wave 4 Specs themselves, because they ran on the v0.20.0 binary. The mechanical authorization audit reads a grant at the merge base of each Task commit (Spec 0178), so a grant widened mid-delivery covers only Task commits rebased onto it.
- Action / suggestion: covered by Spec 0181 (recorded paths, the pre-PR row and the ADR horizon, plus task_06's authored-text citation walk) and Spec 0182. Still unaddressed: the per-commit grant base forces a rebase after any mid-delivery grant widening. That needs a Backlog Entry if it recurs.

## 2. A corrective amendment before carry-forward makes the retry refuse the whole set

- Symptom / evidence: the first `deliver retry` of 0181 refused, because the Supervisor's amendment of `_prd.md`, `_tasks.md` and `_techspec.md` moved every Task's declared inputs. The refusal only pointed to `reconcile --carry-forward`.
- Root cause: carry-forward proves each Task against unmoved declared inputs, and a Spec amendment moves them.
- Action / suggestion: operating rule for now: carry forward first, then commit the amendment. Spec 0182 filters Tasks already completed on the target, which removes the retry side of this. The retry message could name the ordering; this is a Backlog candidate.

## 3. The pre-PR review verdict was unclassifiable twice

- Symptom / evidence: the reviews of `3e77ed70` (Wave 5 authoring) and of 0182's `764e0d5e` ended `blocked`. In both, the answer held a real verdict ("Findings:" or "No findings.") glued to the Agent's progress text.
- Root cause: `internal/agent/acpx_runner.go` joins every `agent_message_chunk` without a boundary.
- Action / suggestion: covered by Spec 0185.

## 4. The queue's retry limit ended both items' delivery

- Symptom / evidence: `--max-retries 2` was exhausted by 0182 at `review-blocked` and by 0181 at its third QA failure. The Supervisor then finished review, archive, gate, Pull Request and merge by hand with `roundfix implement`, `review`, `archive` and `gh`.
- Root cause: the limit counts every park, including parks caused by the defects in sections 1 and 3.
- Action / suggestion: the Wave 5 queue started with five retries. Revisit the default once 0181 and 0185 have shipped.

## What worked — keep

- Revalidation on each item's starting main, `deliver plan` and one detached owner surviving across both items.
- `reconcile --carry-forward` in the item worktree, followed by `deliver retry`, recovered proved Tasks every time it was used in the right order.
- `roundfix reopen --spec` reopened a settled QA gate cleanly after a late corrective Task.
