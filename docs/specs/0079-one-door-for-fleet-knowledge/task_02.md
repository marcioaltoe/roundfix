---
task: task_02
spec: 0079-one-door-for-fleet-knowledge
status: completed
type: chore
complexity: high
---

# Task 02: Carry the contracts and the permission to the guides

## Overview

Tooling task one of two, executed under the express maintainer authorization
recorded at `docs/workflow/authorizations/2026-08-06-fleet-knowledge-door.md`.
The findings contract gains its three extensions (rollup kind, archival home,
license pointer) and the secondbrain contract gains the inbox-entry shape plus
the **permissive** carve-out — sessions may create under the brain's
`inbox/**` and nothing else. Everything is authored in the two Baseline
modules, regenerated through the measured Spec 0075 choreography, and adopted
into both guides as postimages. Permission only: the obligating clauses are
task_06's, after the pilot proves the door.

## Requirements

1. MUST read the current content of both modules and both guides' managed
   blocks before editing; the change is authored against the real surface,
   not from memory.
2. MUST extend the findings contract in
   `internal/baseline/assets/modules/context-workflow.json` with, using these
   literal tokens: `kind: rollup` findings declaring resolvable `members:`;
   the `docs/findings/_archived/` home; the `absorbed_by:` license required on
   every archived finding, resolvable to a rollup or a Spec; and the triage
   contract — a pending Inbox Entry resolves into exactly one Finding, one
   Backlog Entry (ADR-0092 boundary: evidence never becomes intent without a
   human choice), or one recorded discard, with the entry's provenance cited
   in the minted artifact.
3. MUST write the inbox-entry contract into
   `internal/baseline/assets/modules/secondbrain.json` with these literal
   frontmatter tokens: `origin:`, `destination:`, `type-hint:`, `created_at:`,
   `capture:` (values `manual` and `auto`), and the triage-time additions
   `resolved_to:` or `discarded_reason:`; positional status — pending at the
   namespace root, resolved under `_triaged/`.
4. MUST replace the secondbrain module's blanket write prohibition with the
   permissive carve-out: sessions MAY create files under the brain's
   `inbox/**`; every other brain path stays read-only, and the guide's
   secret-handling prohibitions are preserved verbatim.
5. MUST NOT author any obligating clause (inbox-first skill steps, mandatory
   session-end capture, mandatory research capture) — permission and contract
   shapes only.
6. MUST run the module chain per the Spec 0075 choreography: bootstrap the
   Source Baseline manifest rows for new clauses (the regenerator maintains
   rows, never creates them), then run `make baseline-digests` twice — the
   maintained fixture is the chain's first step and converges only on the
   second pass.
7. MUST adopt the regenerated managed-block postimages into
   `docs/agents/docs-layout.md` and `docs/agents/secondbrain.md`.
8. MUST change only the authorization's bounded files, their sanctioned
   deterministic digest fallout (ADR-0081), and this task file.

## Subtasks

- [ ] Author the findings-contract extensions in the context-workflow module.
- [ ] Author the inbox-entry contract and permissive carve-out in the
      secondbrain module.
- [ ] Bootstrap manifest rows and run the two-pass digest regeneration.
- [ ] Adopt both guide postimages.

## Acceptance Criteria

- [ ] Both guides carry the new contracts with the literal tokens named in
      Requirements 2–3.
- [ ] The secondbrain guide permits creation under `inbox/**` only and keeps
      every secret-handling prohibition.
- [ ] No obligating clause exists anywhere in the diff.
- [ ] The digest chain is converged: a further regeneration changes nothing.
- [ ] The diff touches only bounded files, sanctioned digest fallout, and
      this task file.

## Context

- instruction: docs/workflow/authorizations/2026-08-06-fleet-knowledge-door.md
- interface: internal/baseline/assets/modules/context-workflow.json
- interface: internal/baseline/assets/modules/secondbrain.json
- interface: docs/agents/docs-layout.md
- interface: docs/agents/secondbrain.md

## Verification

- `grep -q "kind: rollup" docs/agents/docs-layout.md && grep -q "absorbed_by" docs/agents/docs-layout.md && grep -q "_archived" docs/agents/docs-layout.md`
  — expected: exit 0; the findings-contract extensions are adopted in the
  layout guide.
- `grep -q "inbox/" docs/agents/secondbrain.md && grep -q "type-hint" docs/agents/secondbrain.md && grep -q "resolved_to" docs/agents/secondbrain.md && grep -q "_triaged" docs/agents/secondbrain.md`
  — expected: exit 0; the entry contract and carve-out are adopted in the
  secondbrain guide.
- `grep -q "kind: rollup" internal/baseline/assets/modules/context-workflow.json && grep -q "inbox/" internal/baseline/assets/modules/secondbrain.json`
  — expected: exit 0; the contracts live in the inheritable modules, not only
  in this repository's guides.
- `output="$(grep -c "Do not write to the Secondbrain" docs/agents/secondbrain.md)"; [ "$output" = "0" ]`
  — expected: exit 0; the blanket prohibition was replaced by the carve-out,
  not left contradicting it.
- `s1="$(git status --porcelain | sort)"; make baseline-digests >/dev/null 2>&1; s2="$(git status --porcelain | sort)"; [ "$s1" = "$s2" ]`
  — expected: exit 0; the digest chain is converged — one more regeneration
  is a no-op on the worktree.
- `output="$(git status --porcelain | awk '{print $NF}' | grep -vE '^(internal/baseline/assets/|internal/baseline/testdata/|docs/agents/docs-layout\.md$|docs/agents/secondbrain\.md$|docs/agents/setup-context\.json$|docs/specs/0079-one-door-for-fleet-knowledge/task_02\.md$)')"; [ -z "$output" ]`
  — expected: exit 0; nothing outside the bounded files and sanctioned
  fallout changed.

## References

- `_prd.md` → Core Features 1–5; User Stories 1–3; Project Constraints
  (Tooling authority).
- `_techspec.md` → Implementation Design (Interfaces); Testing Approach
  (module choreography); Build Order 1.
- ADR-0095, ADR-0092, ADR-0081.

## Result

### Implementation

- Versioned the context-workflow module, its docs-layout guide, and its layout
  rule. Added separate clauses for `kind: rollup` plus resolvable `members:`,
  the `docs/findings/_archived/` home plus required `absorbed_by:` license, and
  the one-entry/one-outcome Triage contract with the ADR-0092 human-choice
  boundary and minted-artifact provenance.
- Versioned the secondbrain module, its guide, and the edited read-only rule.
  Added the positional Inbox Entry contract with `origin:`, `destination:`,
  `type-hint:`, `created_at:`, `capture: manual # manual | auto`, and exactly
  one triage-time `resolved_to:` or `discarded_reason:`. Pending entries stay
  at `inbox/<destination>/`; resolved entries move under `_triaged/`.
- Replaced the blanket write prohibition with permission for sessions to
  create only under `inbox/**`. Every other Secondbrain path remains read-only;
  the module and guide secret-handling prohibitions remain byte-identical to
  their `HEAD` preimages.
- Bootstrapped five marker-delimited Source Baseline entries and manifest rows:
  `clause.context.findings-08-rollup`,
  `clause.context.findings-09-archive`,
  `clause.context.inbox-01-triage`,
  `clause.secondbrain.inbox-write-permission`, and
  `clause.secondbrain.inbox-entry-contract`. Regeneration replaced every
  temporary span and digest with values calculated from the source bytes.
- Ran the required digest choreography after the final module correction. The
  first `rtk make baseline-digests` pass exited 0 with `changed:true`; the
  second exited 0 with `changed:false`. Adopted both regenerated managed blocks
  into the local guides without changing their repository-authored content
  except the now-permitted Secondbrain access statements.

### Focused checks

- Focused `jq -e` contract assertions over both modules exited 0. They checked
  module/guide/rule versions, every required literal token, resolvable-member
  and archival targets, positional status, the `manual | auto` capture enum,
  the single triage result, the ADR-0092 boundary, provenance, and the exclusive
  `inbox/**` permission.
- Exact `cmp -s` checks showed that each adopted managed guide block is
  byte-identical to its regenerated formatter golden. Five additional
  `cmp -s` checks showed that every new Source Baseline marker span is
  byte-identical to its canonical module guidance.
- Exact `cmp -s` against `HEAD` showed the secondbrain module's
  `rule.secondbrain.secret-safety` guidance is unchanged. A focused absence
  check found no remaining `Do not write to the Secondbrain` blanket clause.
- The forbidden-obligation scan found no inbox-first step, mandatory
  session-end capture, mandatory research capture, or `must capture` clause in
  the edited modules. No authorial skill or hook path changed.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test
  ./internal/baseline -count=1 -run
  "TestCatalogCompatibility|TestFormatterComposition"` passed both focused
  tests. The first attempt with the environment's macOS Go cache did not reach
  compilation because the sandbox denied that cache path; the task-scoped
  cache retry is the behavioral evidence.
- `rtk git diff --check` exited 0 after the guide adoption.

### Acceptance-criterion evidence

1. Both guide managed blocks match their regenerated goldens exactly. The
   docs-layout block contains `kind: rollup`, `members:`, `_archived/`,
   `absorbed_by:`, and the one-result Triage contract; the secondbrain block
   contains every Inbox Entry field and `_triaged/` positional status.
2. The secondbrain module, generated managed block, and repository-owned access
   statements all make `inbox/**` the only writable namespace. Exact preimage
   comparison proves the secret-safety clause is unchanged.
3. The focused forbidden-obligation assertion exited 0, and diff inspection
   found only permission and contract shapes. Inbox-first skill steps,
   session-end capture obligations, and research-capture obligations remain
   outside this diff.
4. The final two-pass regeneration sequence exited 0 and ended with
   `{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":false}`.
5. The 20 changed paths are the two authorized modules, two authorized guides,
   this Task file, two formatter goldens, one profile digest pin, five
   Source Baseline corpus/manifest/index artifacts, three catalog
   digest/normalization/diagnostic artifacts, and four plan-characterization
   goldens. No path outside the authorization boundary or ADR-0081 fallout is
   changed.

### Follow-up

- The focused non-update
  `TestReadoptionCompatibilityMaintainedFixture` check reached the assertion
  and failed because the regenerated Source Baseline now correctly has 111
  entries while `internal/baseline/preservation_test.go` still pins 106. That
  Go file is outside this tooling Task's exact authorization. The Go-owned
  consistency-check slice must update the named expectation before the full
  repository gate; this Task does not widen its diff to do so.
- The Daemon-owned commands under `## Verification` were not run in this Agent
  turn. Task status remains Daemon-owned.
