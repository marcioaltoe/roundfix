---
status: done
created_at: 2026-09-30
updated_at: 2026-09-30
absorbed_by: 0205-an-advisory-judge-for-spec-authoring
---

# Jev judgments measured against the Spec archive (2026-09-30)

Five typed judgments were scored against labels derived from `docs/history/specs` (171 archived Specs), `docs/adr` (146 accepted ADRs) and the adopted findings and Backlog Entries. Two pass as advisory checks, one is inconclusive, and two do not pass.

- **Model:** `jev-latest` resolved to `jev-1.13.0` on every response.
- **Calls:** direct `POST https://api.typesafe.ai/v1/systemone`, no gateway. 5,240 requests, 5.19M input tokens, **US$ 0.218**, zero failures and zero retries.
- **Latency:** p50 330–335 ms in every benchmark; p95 between 493 ms and 1.33 s.
- **Data sent:** fragments of PRDs, TechSpecs, Task files, ADRs, findings and Backlog Entries of this repository only. No source code, diffs, QA evidence or Secondbrain content.

## Results

| # | Benchmark | n | Base rate | Baseline | Jev | Verdict |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | Citation support: does the ADR say what the Spec attributes to it? | 600 pairs from 200 claims | 33% supports | Word overlap: AUROC 0.80; 90% accuracy against random swaps | AUROC 0.89 [0.87, 0.92]; 94.5% accuracy against random swaps (AUROC 0.99); 0.80 against nearest-ADR swaps (baseline 0.65) | **Passes** as an advisory check |
| 5 | Goal → mechanism: does the TechSpec section deliver the PRD goal? | 286 pairs, 41 Specs | 50% | Word overlap: AUROC 0.77 [0.72, 0.83] | AUROC 0.86 [0.81, 0.90]; real section beats a random one in 83.9% of pairs (baseline 72.7%) | **Passes** as an advisory check, with low recall |
| 3 | Task Type | 399 tasks | 60% backend | Majority 60.1%; keyword rule 75.7% | 84.2% population-weighted; 89.5% at confidence ≥ 0.9 | Passes as a calibration benchmark; **not worth adopting** |
| 2 | Overlap shortlist: which Spec covers this source? | 99 sources × 20 candidates, two negative sets | 5% | TF-IDF recall@1 0.82 (random set), 0.71 (hard set) | recall@1 0.84 (random), 0.64 (hard); recall@5 0.95 and 0.86 (TF-IDF 0.93 and 0.82) | **Does not pass**: no gain over TF-IDF |
| 4 | Task lint: do three Nouls predict a Verification repair? | 300 tasks | 23% repaired | Chance, AUROC 0.50 | Vague acceptance: AUROC 0.62 [0.55, 0.70]; verification targets the Task: 0.60; file outside context: 0.50 | **Inconclusive**: a weak real signal, too weak to flag a Task |

### 1. Citation support

- **Operating point:** raise a finding when the answer is not `supports` at confidence ≥ 0.8.
- **On real citations:** 3 of 200 were flagged. All three quote the claim that ADR-0083 makes `make verify` the authoritative gate, which ADR-0116 records as a false citation. Jev was right on all three, so the measured false-alarm rate is 0 of 197.
- **On planted swaps:** it caught 139 of 200 random swaps (69.5%) and 21 of 200 nearest-ADR swaps (10.5%).
- **Without the confidence floor:** 14 of 200 real citations are flagged. The extra 11 are context-dependent sentences ("makes this the highest-stakes path") or low-confidence answers.
- **Nearest-ADR swaps are a noisy label.** A related ADR often does support the claim. The top "false accepts" pair a claim about the QA Task node with ADR-0088 instead of ADR-0091, and both records say it.

### 5. Goal → mechanism

- **Operating point:** raise a finding when P(section delivers goal) < 0.3.
- **On real mapped sections:** 3 of 143 were flagged (2.1%).
- **On random sections of the same TechSpec:** 32 of 143 were caught, with 91% precision among the flags.
- The negatives are random sections, which is easier than a real coverage defect. Recall on real defects is unmeasured.

### 3. Task Type

- Calibration is monotone: 33% accuracy in the lowest confidence decile, 100% in the top three.
- `qa` is 50 of 50, `backend` and `test` 87–88%. `chore` is 3 of 24 and `data` 3 of 11, and many of those labels are arguable (a `chore` that only moves documents).
- At confidence ≥ 0.99 it still disagrees with the author on 12 of 210 tasks. The authoring model already chooses the type, so a check here saves nothing.

### 2. Overlap shortlist

- Against 19 random Specs, Jev and TF-IDF are within two points of each other.
- Against the 19 lexically nearest Specs, Jev is worse at rank 1 (0.64 against 0.71) and slightly better at rank 5.
- As a filter on the hard set, keeping score ≥ 0.5 keeps 59% of candidates for 96% recall; score ≥ 1.0 keeps 26% and drops recall to 82%.
- The label favors the baseline: each PRD was written from its adopted source, so the two share wording by construction. Finding a duplicate written in different words was not tested.

### 4. Task lint

- The direction matches the hypothesis. Tasks that needed a Verification repair score lower on "every acceptance criterion is observable" (mean 0.65 against 0.71).
- Backend Tasks only: AUROC 0.67 [0.59, 0.76], n = 198.
- Of the 30 Tasks scored below 0.5, 37% needed a repair, against a 23% base rate. That is too weak to act on per Task.
- "Requirements name a file outside Context" carries no signal.
- The label is noisy: it is whether the Agent-written `## Result` mentions verification feedback, which reflects implementation difficulty as much as Task quality.

## Recommendation

**Adopt as advisory checks** (a finding the authoring model must answer; never a gate):

1. **Citation support**, inside `write-prd` and `write-techspec`, for every sentence that attributes something to an ADR. It complements the repository's word-overlap detector, which it beats against both kinds of swap.
2. **Goal → mechanism**, inside `write-techspec`, for each Coverage Map line. It is cheap, false alarms are rare, and its recall is modest.

**Do not adopt:**

- **Overlap shortlist.** A lexical shortlist already does as well.
- **Task Type.** Accurate but redundant.
- **Task lint.** Keep the acceptance question as a candidate, and re-measure it once Runs record a cleaner label (Verification attempts per Task in the Run Database).

## Limitations

- **One model version, one pass.** Everything ran once on `jev-1.13.0`. Repeatability was not measured. The thresholds belong to this version.
- **English only.** Every artifact is English. Nothing here supports Portuguese repositories.
- **Citation leakage was removed, not proven absent.** ADR numbers and `(Spec NNNN)` were stripped from claims. Recent Specs already passed the repository's word-overlap check, which flatters the lexical baseline.
- **Citation claims come from 80 ADRs**, capped at six claims per ADR, using sentences with exactly one ADR number.
- **Overlap:** 99 sources, not 100, because only indexed findings and Backlog Entries are within the data authorization. One Spec with seven adopted sources counts seven times. Only the first 1,500 characters of each text were sent; slugs and dated filenames were scrubbed.
- **Goal → mechanism:** 41 Specs, not 60. Only 41 of 170 TechSpecs have a Coverage Map line that names a section heading. The section title was part of the state.
- **Task Type** labels are the authors' own choices and some are arguable. The sample is stratified, and the headline number is weighted back to the archive's mix.
- **Thresholds were chosen on the same data they are reported on.** Expect some optimism.

## Clause-consistency scan of the Baseline

A second run sent only Baseline clause text (`internal/baseline/assets/modules/*.json`) and the text of the 14 Roundfix-owned skills.

- **Run facts:** `jev-1.13.0`, 2,758 requests, zero failures, 2.26M input tokens, US$ 0.095, p50 328 ms.
- **Clause pairs** (1,091 shortlisted): one confirmed defect, the duplicated `clause.backend.boundary-contracts` / `rule.backend.boundary-contracts` entry that Spec 0193 removed. Every flagged "contradiction" between 0.50 and 0.70 confidence was a wording tension or a false alarm.
- **Planted probe:** Jev caught 4 of 4 planted product-truth contradictions and 23 of 30 force flips. It labelled 10 of 10 identical texts `duplicate`.
- **Per-clause Nouls:** "ambiguous" saturates (mean 0.79), so it does not discriminate.
- **Skill vs clause** (1,485 pairs): Jev kept nothing at confidence ≥ 0.70. Its 24 lower-confidence `contradicts` answers held 4 real contradictions (17% precision).
- **Conclusion:** the recurring pre-release skills and guides check rests on plain-code checks (duplicate text, force table, citation and command coverage; Specs 0192, 0193, 0195), not on Jev.

## Where the material lives

The scripts, labelled sets and JSONL logs were produced in a session scratchpad and are not committed. The Spec that adopts the two advisory checks must carry the question texts and the thresholds above as its own reviewable source.

## Addendum (2026-09-30): routing

- Recommendations 1 and 2 (citation support and goal → mechanism as advisory checks) → Spec `0205-an-advisory-judge-for-spec-authoring`, which carries the question texts and thresholds of this finding as its reviewable source, pins `jev-1.13.0` and records the maintainer's data, key and spending decisions in ADR-0200 and ADR-0201.
- The overlap shortlist and Task Type are not adopted, as recommended; no Backlog Entry is needed.
- Task lint stays a candidate to re-measure once Runs record Verification attempts per Task → Backlog Entry `docs/backlog/2026-09-30-re-measure-the-task-lint-judgment-on-a-cleaner-label.md`.
- The clause-consistency scan needs no action: its one confirmed defect was removed by Spec 0193.
