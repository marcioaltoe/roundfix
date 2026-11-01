---
type: fix
status: open
created: 2026-11-01
spec: null
reason: null
---

# QA evidence links resolve in the committed tree

## Symptom

A QA Report settled `completed` (qualifying partial, Run Clean) while it linked
evidence files that the QA commit did not contain: the `.log` files were
ignored by the repository's `.gitignore`, so they never reached the commit, and
the Run Worktree that held them was removed. After integration,
`roundfix spec check --strict` reported `SC-HISTORY-RETAINED` for both missing
links.

## Where

QA settlement (`roundfix qa-report accept` and the Daemon's QA settlement) and
the qa-gate skill's evidence rules.

## Expected

Before the QA Task settles, every local evidence link in the report resolves
in the prospective committed tree; an unresolved link fails settlement with the
link and the reason (ignored by Git, missing), so the gate keeps or renames the
evidence instead of losing it silently. The qa-gate skill tells the Agent to
save evidence under names the repository does not ignore. Report truth is
preserved; no check is weakened.

## Evidence

Reported by Argus (`default/gesttione-solutions-argus`, ai-memory message
`01a19895`, 2026-11-01, roundfix 0.127.0): Spec `0005-publicar-downloads-completos`,
Run `run_20261101T191252Z_133b75a4f68959a6`; the report at
`edd9b619:docs/specs/0005-publicar-downloads-completos/qa/qa-report-2026-11-01.md`
linked `evidence/2026-11-01-local/consumer-tests.log` and `journeys.log`;
`git ls-tree -r edd9b619 -- …/qa/evidence` shows only `journeys.ts.txt` and
`scope-audit.txt`. Argus repaired it in `c1f73185` by replaying the evidence
into `.txt` files.
