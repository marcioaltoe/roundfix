---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The PRD template and the checker disagree on the tooling row

## Symptom

The `write-prd` template instructs, for a Spec that mutates no tooling:

> With no protected tooling mutation, record: `applicable — no protected tooling
> mutation proposed or authorized`.

Following that instruction, `roundfix spec check --stage prd` refuses:

```text
[error] SC-TOOLING-UNBOUNDED: _prd.md makes Tooling authority applicable without bounded files
  fix: Add an exact bounded files list to the Tooling authority row in _prd.md.
```

The detector requires a bounded-file list for **any** `applicable`, including the
one that declares there is no mutation. Both wordings the template suggests were
tested, with and without a trailing explanation; both refuse. Only
`not applicable — <reason>` passes.

The cost lands on every Spec that touches no tooling, which is most of them, and
it lands at the moment an author is trusting the template.

## Where

The `write-prd` template's Project Constraints guidance and the
`SC-TOOLING-UNBOUNDED` detector.

## Expected

One of the two changes: the template records `not applicable — <reason>` for a
Spec with no tooling mutation, or the detector accepts `applicable` without
bounded files when the row's own reason says no mutation is proposed. The first
is the smaller change and matches what already passes.

Worth settling in the same work: which reading the row is meant to carry.
`applicable` currently means both "this constraint governs this Spec" and "this
Spec mutates tooling", and the detector reads the second while the template
writes the first.

## Evidence

Minted from the Inbox Entry
`inbox/roundfix/2026-08-12-spec-check-e-template-do-write-prd-discordam-na-linha-de-tooling.md`
in the Secondbrain. Reproduced in `roundfix` on 2026-08-12 while authoring the
seventeen PRDs of the current Spec set: the first PRD written from the template
refused for exactly this reason, and all seventeen use `not applicable` as a
result.
