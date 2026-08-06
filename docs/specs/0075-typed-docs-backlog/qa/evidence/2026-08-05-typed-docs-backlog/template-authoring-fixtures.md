# Template-only authoring fixtures

These disposable entries were authored from the freshly generated layout
guide without reading the Spec. An independent probe found exactly five
frontmatter fields in every entry and every type-specific heading.

## feat

```markdown
---
type: feat # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-05
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Show the fastest relevant local check

## Opportunity

Contributors could see the narrowest valid verification command for the files they changed.

## Value

Shorter feedback loops might help contributors find defects before the full gate.

## Shape

Add non-binding command guidance to the existing verification output.
```

## fix

```markdown
---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-05
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Missing decision diagnostic omits the decision source

## Symptom

An operator sees a missing-decision message without the file or flag that can supply it.

## Where

`roundfix baseline plan`.

## Expected

The diagnostic names the accepted decision input paths.

## Evidence

none yet
```

## perf

```markdown
---
type: perf # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-05
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Baseline planning latency

## Slow

Repeated Baseline planning delays an adopting maintainer's feedback loop.

## Measured

The QA scratch-repository plan completed in under one second on 2026-08-05.

## Target

Keep the same scratch-repository plan under one second.
```

## refactor

```markdown
---
type: refactor # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-05
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Separate plan rendering from plan calculation

## Tangled

Plan calculation and JSON rendering are coupled at the CLI boundary.

## Cost

The coupling makes alternate output formats risky to change.

## Shape

Return a typed plan result before selecting its non-binding output representation.
```

## Lifecycle probes

- A promoted `feat` probe set `status: promoted`, set
  `spec: 0099-probe`, and moved from `docs/backlog/` to
  `docs/specs/0099-probe/references/`.
- A declined `fix` probe set `status: declined` with
  `reason: duplicate of an existing tracked intent`.
- A standalone `refact` search over the generated guide returned no match.
