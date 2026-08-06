---
task: task_02
spec: 0079-one-door-for-fleet-knowledge
status: pending
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
