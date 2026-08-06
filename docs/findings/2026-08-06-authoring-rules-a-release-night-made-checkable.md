---
status: pending
created_at: 2026-08-06
updated_at: 2026-08-06
---

# Spec authoring — rules a release night made checkable (2026-08-06)

Eight Specs went from queue to released v0.4.0 in one supervised session. Every
QA gate that failed did so on an authoring defect, not an implementation one,
and each defect is mechanically checkable at authoring time — the same
conclusion `2026-08-06-every-run-that-failed-tonight-failed-on-a-contract.md`
reached for Daemon contracts, now reached for Spec artifacts. Spec 0065 shipped
the `SC-*` enforcement surface the same night; these are its next rules.

## 1. A promise in a Spec artifact that no Task was asked to keep ships nothing

- Symptom / evidence: twice in one night. Spec 0059's TechSpec declared
  `roundfix gc compact [--apply]` in its API Contracts table; no Task
  requirement asked for the route, so `internal/store` shipped `Compact` with
  every guard and `roundfix gc compact` exited `2` with
  `unexpected argument` — gate F-001, `Blocks-Completion`. Spec 0073's Success
  Metric 7 (provenance for three external Go Skills) had no owning Task and was
  discovered undelivered by the gate (F-002).
- Root cause: Coverage checks map Core Features and User Stories
  (`SC-COVERAGE-UNMAPPED`, `SC-COVERAGE-UNTASKED`) but neither API-Contracts
  table rows nor Success Metrics to Task requirements.
- Action / suggestion: two `SC-*` rules — every API Contracts row and every
  Success Metric must resolve to at least one Task requirement, exactly as
  Core Features already must.

## 2. Wave parallelism requires file-disjoint Tasks, and the graph declares enough to check it

- Symptom / evidence: Spec 0059's wave 2 ran task_02, task_03, task_04
  concurrently; all three write `internal/store/journal.go` and its test. Two
  failed at the serialized cherry-pick with
  `integration conflict: internal/store/journal_test.go` and cost a settle
  cycle plus two hand-merged integrations.
- Root cause: independence was judged by concept (compaction vs discovery vs
  policy), not by file surface. The `## Context` sections declared the shared
  interface, so the collision was visible at authoring time.
- Action / suggestion: an `SC-*` rule flagging same-wave Tasks whose declared
  Context paths intersect; the fix at authoring is a dependency edge or a
  disjoint split.

## 3. A Verification written before reading the surface it asserts fails on the truth

- Symptom / evidence: three instances. `^type: (feat|fix|perf|refactor)$`
  rejected the Spec's own template, which documents the enum inline
  (`type: perf # feat | fix | perf | refactor`). A provenance assertion
  targeted `roundfix skills list`, which lists only binary-bundled skills, so
  it could never pass for an external skill — `skills-lock.json` was the
  surface. An absence search over five carriers omitted Task files, then the
  widened version would have failed on the Spec's own prose *quoting* the wrong
  string it corrects.
- Root cause: assertions authored from memory of the surface rather than from
  the surface. The third case adds the boundary: an absence check that cannot
  distinguish a claim from a citation of that claim is unusable in both
  directions.
- Action / suggestion: authoring discipline for `write-tasks` — before a
  Verification lands, run its command against the real surface and record the
  red-then-green pair in the Task; scope negative assertions to shipped
  carriers and cover the rest with positive ones.

## 4. Proving a property by running the whole gate makes the Spec hostage to every flake

- Symptom / evidence: Spec 0073's headline property — an owned-skill edit needs
  no regeneration step — was proven by `TestOwnedSkillEditLeavesMakeVerifyGreen`
  running the complete `make verify` nested inside a test. The gate observed it
  fail on `TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`,
  an unrelated concurrency-sensitive test that passes in isolation.
- Root cause: the instrument was wider than the claim. The claim is "these
  artifacts do not move"; the instrument was "the world is green".
- Action / suggestion: assert artifact stability directly — byte-compare the
  derived artifacts across the edit. Candidate `SC-*` rule: a Task Verification
  may not invoke the repository gate from inside a test.

## 5. Load-sensitive tests now block releases, not just pull requests

- Symptom / evidence: four flake instances in one night, all passing in
  isolation and failing under concurrent Runs:
  `TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree`,
  `TestACPXRunAppliesFullAccessSessionSetup`, `TestProjectDecisionJourney`
  (8.10s in CI vs 3.54s locally), and `TestSealedACPXPromptDiscardsLargeThoughtStreamIncrementally`
  — the last one **inside the release workflow's Verify gate**, stopping the
  v0.4.0 publication until a rerun.
- Root cause: timing- and allocation-sensitive assertions in
  `internal/cli`/`internal/agent` under host load; the machine ran up to four
  concurrent Runs across repositories.
- Action / suggestion: fold into the open test-performance campaign as its own
  backlog entry: identify the load-sensitive set, make their budgets
  load-tolerant or serialize them, because a 1-in-N flake is noise in PR CI
  and a publication blocker in the release gate.

## 6. Standing policies that removed the human from the loop, worth codifying

- Symptom / evidence: the decisions that required maintainer interaction
  tonight were almost all repeats of a policy stated once earlier in the same
  session: spend `coderabbit:review` only on substantive Spec PRs; when review
  capacity blocks, merge on QA-pass + `make verify` + read diff and record the
  open findings on the PR; the 02:00 timeboxed fallback (cut the largest
  remaining Spec) was pre-declared and executed without a wake-up.
- Root cause: session-scoped instructions live in conversation, not in
  `docs/agents/autonomous-work.md`, so every future session re-asks.
- Action / suggestion: promote to the autonomous-work clause set: the review
  economy policy, the merge-with-gates fallback wording (with its mandatory
  PR-comment audit trail), and pre-declared timeboxed fallbacks as a pattern.
  The 2026-08-04 standing-grant shape (purpose-bounded, checkable covered
  list) is the template — it survived six Specs and one `SC-TOOLING-UNAUTHORIZED`
  audit without a single re-ask.

## What worked — keep

- The authored QA gates caught nine real findings across five Specs, including
  three errors in artifacts this session itself authored — one proven against
  the archived source rather than the carrier quoting it.
- Agents refused out-of-scope work three times rather than passing their own
  checks: a bounded Batch declined to edit another PR's Review Issue files, a
  tooling Task declined to widen its own grant to `skills-lock.json`, and a
  skill-sync Task declined to document a command that did not exist.
- Both release-gate refusals (`Validate tag`, nested `Verify gate`) stopped the
  run before anything published — fail-closed release plumbing turned two
  mistakes into reruns instead of half-published coordinates.
