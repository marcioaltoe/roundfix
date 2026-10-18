---
task: task_03
spec: 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
status: pending
type: backend
complexity: high
---

# Task 03: A Delivery Retry accepts a Secret Rewrite, and `secret-found` names the rewrite that works

## Overview

On 2026-10-18 the Conexus operator followed the `secret-found` next action,
rewrote the archived item branch, and the retry refused the head because it
did not descend from the candidate (`_prd.md` → Sources 2). The only way out
was a reset to the commit before the archive.

This Task adds the **Secret Rewrite**: a non-descending head whose tree
differs from the parked candidate's only at lines that held a Secret Match
there, and whose history starts at or after the item's Run start. The retry
accepts it and returns the item to `reviewing`. `secret-found` then prints the
rewrite that works. It is verifiable on its own through prover tests on
temporary Git repositories and engine tests with a fake prover.

## Requirements

1. MUST answer the Finding of 2026-10-18,
   `references/2026-10-18-secret-found-asks-for-a-rewrite-the-retry-of-an-archived-item-refuses.md`
   (`_prd.md` → Sources 2, finding 1), by implementing API Contracts 6 to 8
   of `_techspec.md` (ADR-0292).
2. MUST add `SecretRewrite`, `SecretRewriteProver` and
   `Dependencies.SecretRewrites` to `internal/delivery/engine.go`, and in
   `Retry`'s archived branch apply API Contract 6.
   - The rule runs only for a `secret-found` item whose moved head was not
     accepted by descent.
   - The base comes from `history.RunStart(item.RunID)`.
   - The named paths come from the blocker's Secret Match entries.
   - An accepted proof appends the head and sets `reviewing`.
   - A refusal keeps today's archived-head reason followed by ` (<rule>)`.

   Every other blocker's retry, refusal text and stage MUST stay unchanged.
3. MUST implement `ProveSecretRewrite` on `commandDeliveryWorkflow` in the new
   `internal/cli/deliver_secret_rewrite.go` and wire it as `SecretRewrites` in
   `internal/cli/deliver_workflow.go`. It applies API Contract 7's rules in
   order with their exact reasons.
   - The data source is the item worktree's Git objects: the Run start, the
     parked candidate and the head, read with `git merge-base --is-ancestor`,
     `git diff --raw --no-renames` and `git cat-file`.
   - Matched lines come from `secretscan.ScanContent` of each version.
   - It writes nothing. Its reason, read by the operator, names paths and line
     numbers only, never a value.
4. MUST make `ClassifyPark` print API Contract 8's next action for
   `secret-found`. This is a declared break of two pinned assertions:
   - `internal/delivery/park_class_test.go` (line 88);
   - `internal/delivery/push_secret_scan_test.go` (line 114).

   Both move to the new text and keep their names.
5. MUST add these tests, each calling `t.Parallel()` first and compiling for
   Windows. Every secret-shaped sample is built at runtime from parts
   (ADR-0273).
   - In the new `internal/cli/deliver_secret_rewrite_test.go`, using temporary
     Git repositories from the fixtures already in the package:
     - `TestSecretRewriteAcceptsReplacedMatchedLines`: a rewrite from the
       first secret commit that replaces only the matched lines with
       letters-only values is accepted;
     - `TestSecretRewriteRefusesAnyOtherChange`, with one subtest per rule of
       API Contract 7, each asserting its exact reason.
   - In the new `internal/delivery/secret_rewrite_retry_test.go`, with a fake
     prover:
     - `TestRetryAcceptsASecretRewriteOfAnArchivedItem`: the stage is
       `reviewing`, the head is the newest candidate, and the prover received
       the Run start, the candidate, the head and the named paths;
     - `TestRetryRefusesARejectedSecretRewriteWithItsRule`: the refusal is
       the archived-head reason followed by ` (<rule>)`, and the item is
       unchanged.

   Every existing retry and archived-head test MUST keep its name and pass.
6. MUST document the change, each file naming `Secret Rewrite`:
   - `docs/user-guide/commands/deliver.md`, where the Park Class table's
     `secret-found` action becomes API Contract 8's text and the retry table
     gains the row "Archived Spec parked `secret-found` with a Secret Rewrite
     of its parked candidate | `reviewing`";
   - `.agents/skills/roundfix/references/deliver.md`, where the same action
     and rule are stated.

   MUST search the guides, the agent guides and the skill copies for the old
   next action "remove the secret from every commit the blocker names" and
   replace each. `CHANGELOG.md` and the ADRs stay as they are. No text inside
   any `### QA settlement` section may change.
7. MUST add **Secret Rewrite** to `CONTEXT.md` through domain-modeling, with
   ADR-0292's link.
8. MUST run these commands, because this Task edits a user guide page and the
   Roundfix reference. No record or mirror is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to raise and record the Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Add the prover interface, the dependency and the retry rule.
- [ ] Implement and wire the prover.
- [ ] Change the next action and its two pinned assertions.
- [ ] Add the four tests.
- [ ] Document the change in the guide, the reference and the glossary.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] A rewrite that replaces only the matched lines returns the archived item to `reviewing`.
- [ ] Every other rewrite is refused with the archived-head reason and the rule that failed.
- [ ] `secret-found` prints the rewrite that works, and no guide or skill copy keeps the old action.

## Consumers

- `delivery.SecretRewriteProver`, `delivery.Dependencies.SecretRewrites` and the `secret-found` next action — consumers: `./internal/delivery`, `./internal/cli`

## Context

- instruction: `docs/adr/0292-a-task-commit-is-scanned-for-secrets-and-a-retry-accepts-a-secret-rewrite.md`
- instruction: `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- instruction: `internal/cli/deliver_review_correction.go`
- instruction: `internal/secretscan/scan.go`
- creates: `internal/cli/deliver_secret_rewrite.go`
- creates: `internal/cli/deliver_secret_rewrite_test.go`
- creates: `internal/delivery/secret_rewrite_retry_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/park_class.go`
- interface: `internal/delivery/park_class_test.go`
- interface: `internal/delivery/push_secret_scan_test.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestSecretRewriteAcceptsReplacedMatchedLines|TestSecretRewriteRefusesAnyOtherChange)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSecretRewriteAcceptsReplacedMatchedLines TestSecretRewriteRefusesAnyOtherChange; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestRetryAcceptsASecretRewriteOfAnArchivedItem|TestRetryRefusesARejectedSecretRewriteWithItsRule)$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRetryAcceptsASecretRewriteOfAnArchivedItem TestRetryRefusesARejectedSecretRewriteWithItsRule; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Secret Rewrite' || { printf 'missing Secret Rewrite in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'rewrite the item branch from the earliest commit the blocker names' || { printf 'missing next action in %s\n' "$f" >&2; exit 1; }; done; if grep -rqF -- 'remove the secret from every commit the blocker names' docs/user-guide docs/agents .agents/skills skills internal/delivery; then printf 'old secret-found next action remains\n' >&2; exit 1; fi; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Secret Rewrite**:' || { printf 'missing glossary term Secret Rewrite\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/delivery || exit 1; go test -count=1 -run 'Retry|Deliver|Secret|Archived' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the four tests do not exist, so the command fails at its first pass check. Afterwards these checks pass: the prover and retry tests, the documentation, the removal of the old next action, the glossary, the skill version and the Behavior Surface Record.

## References

- `_prd.md` → Sources; Goals; User Story 4; Core Feature 4; Core Feature 5; Success Metric 4; Glossary
- `_techspec.md` → API Contract 6; API Contract 7; API Contract 8; API Contract 9; Invariant 3; Current behavior; Testing Approach; Build Order 3
- ADR-0292; ADR-0273; ADR-0237; ADR-0272; ADR-0197; ADR-0280; ADR-0189; ADR-0256; ADR-0259; ADR-0285
