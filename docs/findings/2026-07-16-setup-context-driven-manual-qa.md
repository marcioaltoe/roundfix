---
status: done
created_at: 2026-07-16
updated_at: 2026-07-16
---

# Setup Context-Driven — Manual Variation QA (2026-07-16)

Manual dogfood exercised the real `context_setup.py` process against disposable
repositories across all three profiles, multiple durable-decision combinations,
Secondbrain opt-in/out, adoption, skill drift, invalid inputs, and restoration.
The complete matrix and command evidence live in
`docs/specs/0030-context-driven-agent-instructions/qa/`.

## 1. First-run decision output cannot preview the managed change plan

- **Symptom / evidence**: QA-03 invoked `apply --profile
  typescript-bun-monorepo` without answers. The CLI correctly returned exit `3`
  and all nine `decision.required` findings, but returned
  `"plannedChanges": []`. This leaves the skill unable to fulfill its required
  pre-apply explanation of every managed file or block. See
  `docs/specs/0030-context-driven-agent-instructions/qa/evidence/2026-07-16-manual-variations-03/qa-03-first-apply-questions-and-preview.json`.
- **Root cause**: `plan_apply` returns immediately when `resolve_decisions`
  appends findings. Expected artifacts and the change plan are built only after
  every decision is resolved, while `plannedChanges` is derived only from
  finding codes. The current result schema therefore has no representation of
  a profile-level prospective plan while decisions remain unanswered. Audit
  cannot fill this gap because a missing manifest returns before the profile
  override is resolved.
- **Action / suggestion**: extend the spec and task graph with a non-writing
  preview contract that resolves the profile and all decision-independent
  artifacts before confirmation. The machine-readable response should name
  create/refresh/remove operations and explicitly mark any operations whose
  inclusion still depends on an unanswered decision.

## 2. Durable customization answers are stored but do not drive generated guidance

- **Symptom / evidence**: QA-07 applied the Rust profile with
  `spec.scaffold=false`, `domain.layout=multi-context`,
  `triage.external=true`, `autonomous.enabled=false`, a custom verification
  command, and alternate backend/design runtimes. Apply exited `0`, audit
  returned no findings, and the manifest preserved the answers. Generated
  Markdown did not contain the selected verification command or runtime values.
  Static templates also continued to require local `docs/specs/` and autonomous
  delegation regardless of the false answers. See
  `docs/specs/0030-context-driven-agent-instructions/qa/evidence/2026-07-16-manual-variations-03/qa-07-rust-alternate-decision-behavior.json`.
- **Root cause**: `render_expected_path` receives only static template content.
  The decision map is not passed into template rendering or rule selection.
  `ordered_modules_for_decisions` interprets only `secondbrain.enabled`; every
  other initial decision is used for manifest validation/storage but not for
  generated behavior. Audit compares files to the same static templates, so it
  considers contradictory guidance compliant.
- **Action / suggestion**: add explicit decision-to-rule/module/render mappings
  to the declarative asset contract. Conditional workflow decisions should
  include or exclude the relevant managed rules, and string decisions should be
  safely rendered into their owning guides. Audit must validate semantic
  consistency between decisions and generated content. Add macro cases for
  every non-default value, not only manifest persistence.

## What worked — keep

- All three profiles apply and audit through the real CLI when using the current
  default-oriented templates.
- Stored compatible answers prevent repeated questions, and repeated apply is
  byte-for-byte idempotent.
- Secondbrain opt-in produces the compact pointer and complete read-only guide;
  opt-out removes only managed content.
- Adoption, atomic invalid-input handling, required/extra skill classification,
  portable snapshots, embedded-skill synchronization, stdout/stderr discipline,
  nested-instruction safety, and source restoration all passed.
- The misleading `Reserved` help copy found during the sweep was fixed directly
  in commit `d19ca85` and passed the focused retest.

## Routing addendum — 2026-07-16

Both open findings are routed to
`docs/specs/0031-decision-driven-setup-generation/`. Its minimal PRD preserves
the corrective scope, its technical spec defines the shared Decision Plan and
compatible migration, and ADR-0047 records the declarative decision-effect
boundary.
