---
task: task_12
spec: 0119-spec-contained-authorization
status: pending
type: docs
complexity: medium
---

# Task 12: Restore the preserved template guidance and settle the catalog digest

## Overview

Task 08 replaced the templates' tooling-authority wording instead of adding to
it, which breaks the repository's required gate, and Task 07 left the assembled
Baseline manifest carrying a catalog digest one refresh behind. Both are
regressions against work that already passed, and both live in files this
Spec's grant already bounds.

This is an authorized tooling Task. It may change only
`.agents/skills/write-prd/references/prd-template.md`,
`.agents/skills/write-techspec/references/techspec-template.md`,
`skills/write-prd/references/prd-template.md`,
`skills/write-techspec/references/techspec-template.md`,
`docs/agents/setup-context.json`, the derived pins that the sanctioned
regeneration commands rewrite, and this Task file. Stop before any other
mutation. The bounded set comes from the approved grant in
[_authorization.md](_authorization.md).

## Requirements

1. MUST restore the exact preserved phrases `express maintainer authorization`
   and `bounded files` in both the PRD and TechSpec templates, in both the
   canonical and shipped trees.
2. MUST keep the Spec-contained record placement guidance Task 08 added. The
   defect is replacement, not addition: both contracts hold at once, and the
   template must instruct an author to cite the Spec-contained record *and* to
   record express authorization with exact bounded files.
3. MUST regenerate the shipped copies from their canonical source rather than
   editing them by hand.
4. MUST commit the current catalog digest in the assembled Baseline manifest so
   the first public managed refresh reports the repository already current and
   changes zero files. A first refresh that rewrites a tracked manifest is the
   defect, not the fix.
5. MUST NOT hand-edit a digest value; every pin comes from its regeneration
   command.
6. MUST leave the repository's complete Verification passing, which it does not
   today.

## Subtasks

- [ ] Restore both preserved phrases in both canonical templates.
- [ ] Keep the Spec-contained placement guidance alongside them.
- [ ] Regenerate the shipped bundle from canonical source.
- [ ] Regenerate and commit the catalog digest until the first refresh is a no-op.

## Acceptance Criteria

- [ ] The skills contract tests for both templates pass, where they report four
      missing-guidance diagnostics today.
- [ ] Both templates contain the preserved phrases and the Spec-contained
      record placement guidance at the same time.
- [ ] Each shipped template is byte-identical to its canonical source.
- [ ] A public managed refresh on a clean tree reports the repository current
      and changes no tracked file on its first run.
- [ ] The complete repository Verification passes.

## Context

- interface: `skills/baseline_skill_contract_test.go`
- instruction: `docs/agents/specific-repository.md`

## Verification

- `go test -count=1 ./skills -run '^(TestWritePRDProjectConstraints|TestWriteTechSpecProjectConstraints)$'` — the preserved guidance is back; this fails today with four diagnostics.
- `grep -q 'express maintainer authorization' .agents/skills/write-prd/references/prd-template.md && grep -q 'bounded files' .agents/skills/write-prd/references/prd-template.md && grep -q '_authorization.md' .agents/skills/write-prd/references/prd-template.md` — both contracts hold in the same template rather than one replacing the other.
- `grep -q 'express maintainer authorization' .agents/skills/write-techspec/references/techspec-template.md && grep -q 'bounded files' .agents/skills/write-techspec/references/techspec-template.md && grep -q '_authorization.md' .agents/skills/write-techspec/references/techspec-template.md` — the same for the TechSpec template.
- `grep -q 'express maintainer authorization' skills/write-prd/references/prd-template.md || exit 1; grep -q 'express maintainer authorization' skills/write-techspec/references/techspec-template.md || exit 1; make skills-sync-check` — the shipped templates carry the restored phrase and match their canonical source, so the bundle was regenerated rather than hand-edited. The guard reads the phrase missing today, not one already present.
- `committed="$(grep -o '"catalogDigest": *"[^"]*"' docs/agents/setup-context.json | head -1 | sed 's/.*: *"//; s/"$//')" || exit 1; regenerated="$(head -1 internal/baseline/testdata/catalog.digest | tr -d '[:space:]')" || exit 1; test -n "$committed" || exit 1; test -n "$regenerated" || exit 1; test "$committed" = "$regenerated" || { printf 'committed=%s regenerated=%s\n' "$committed" "$regenerated"; exit 1; }` — the assembled manifest carries the current catalog digest, so a managed refresh has nothing to rewrite. The two values disagree today, which is exactly why the first refresh is not a no-op; comparing them directly avoids a subprocess that behaves differently inside a disposable checkout.

## References

- `_prd.md` → Core Features 3; Goals 1; Decisions: Regression locks.
- `_techspec.md` → System Architecture: Authoring guidance; Build Order 5.
- `qa/qa-report-2026-09-09.md` → F-001, F-004.
