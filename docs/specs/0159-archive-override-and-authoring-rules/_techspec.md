---
spec: 0159-archive-override-and-authoring-rules
status: active
created: 2026-09-24
surfaces: [backend, cli, docs]
---

# Archive override and authoring rules

## Executive Summary

Give `roundfix archive` the authorized QA Archive Override the guidance already
describes, make Spec check refuse a claimed ADR ordinal, and make the skills
state settlement once and describe every Task declaration.

## Project Constraints

- Identifier strategy: applicable — a Task's `creates:` path under `docs/adr/`
  claims its number. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0154, ADR-0155 and ADR-0156 hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/spec/archive.go`, `internal/spec/archive_test.go`,
  `internal/cli/cli_test.go`, `internal/speccheck/coherence.go`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `skills/write-tasks/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`,
  `skills/qa-gate/SKILL.md`, `.agents/skills/archive-spec/SKILL.md`,
  `skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## The override

`ArchiveRequest` gains an optional override carrying approval source, reason
and revision. With it, `Archive` still requires every non-QA Task completed,
reads the newest QA Report if one exists, and refuses when that report already
qualifies under the declared-acceptance policy. Otherwise it stamps
`qa_override: true`, `qa_override_approval`, `qa_override_reason`,
`qa_override_qa_outcome` (the observed verdict, `missing`, or the read error)
and `qa_override_revision`, and moves the Spec. The QA Task file and every QA
Report move byte-identically. `internal/cli/archive.go` adds `--qa-override`,
`--approval` and `--reason`, requires the latter two with the first and refuses
them without it, resolves the revision from `HEAD`, and reports the archive as
an override.

## Claimed ordinals

A new detector in `internal/speccheck/ordinal.go`, registered in
`coherence.go` at the Tasks stage, collects every `creates:` path matching
`docs/adr/NNNN-*.md` across the active Specs. It reports `SC-ORDINAL-CLAIMED`
when the same number is held on the tree by a file with a different name, or is
created by Tasks of two different active Specs. A path that already exists under
the exact name it declares is a fulfilled claim, not a collision.

## Guidance

The QA gate, archive and Roundfix skills gain one `### QA settlement` section,
identical in all three, with a row for pass, qualifying declared partial,
environment-blocked, failed, missing and override. The task-writing skill and
its template describe the declarations listed in Core Feature 4. A repository
contract test in `skills/` fails when the three sections differ or when a row
is missing. The canonical files are under `.agents/skills/`; the distributed
mirror is regenerated with `make skills-sync`.

## API Contracts

1. `roundfix archive <slug> --qa-override --approval <source> --reason <text>`
   archives with a stamped override; exit 2 when refused.
2. Spec check reports `SC-ORDINAL-CLAIMED` as an error at the Tasks stage.
3. The `### QA settlement` section is identical across the three skills.

## Coverage Map

- Goal 1 → The override; API Contract 1.
- Goal 2 → Claimed ordinals; API Contract 2.
- Goal 3 → Guidance; API Contract 3.
- Core Feature 1 → The override.
- Core Feature 2 → Claimed ordinals.
- Core Feature 3 → Guidance; API Contract 3.
- Core Feature 4 → Guidance.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- API Contracts 1-3 → The override, Claimed ordinals, Guidance.

## Integration Points

- **Spec 0156.** The delivery loop archives on the branch; it never passes an
  override, so a failed QA still parks the item.
- **Spec 0158.** Its two declarations are the ones the guidance now describes.

## Testing Approach

1. **Override.** An authorized override archives a Spec with failed QA and one
   with no report, stamping every field, with the QA Task and report unchanged;
   it is refused with a non-QA Task pending, with QA already qualifying, and
   without approval or reason.
2. **Ordinals.** A Task creating a number held on the tree under another name,
   and two active Specs creating the same number, each fail with
   `SC-ORDINAL-CLAIMED`; distinct numbers and a fulfilled claim pass.
3. **Guidance.** The contract test finds the section identical in the three
   skills; the task-writing skill names every declaration.
4. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The override (depends on: none).
2. Claimed ordinals (depends on: 1).
3. Guidance and decision record (depends on: 1, 2).
4. Terminal QA (depends on: 1, 2, 3, 5, 6).
5. The override waives what normal archive would refuse (depends on: 3).
6. Blame the latecomer, survive a broken neighbour (depends on: 5).

## Risks & Considerations

- **An override read as success.** Every surface reports it as an override and
  the QA evidence moves unchanged; ADR-0154 already records that it grants
  nothing beyond the move.
- **A false collision.** A claim fulfilled under its own name is not a
  collision, so a completed Task never trips its own check.
