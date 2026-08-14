# Fluxus greenfield command evidence

Repository copy:
`/private/tmp/roundfix-qa-0047.nh68hi/fluxus-greenfield`

The source copy started clean at
`1aeed7e8370c3d14137c42b0c789dcbe3bd1ba3b`.

## Plan and approval

Planning with the completed Decision Document exited `0` and wrote
`plan.json`.

- Plan Digest:
  `sha256:1b9f7ae235df14ad9c35cf5bb75cf15b2e831b4f3a9fbb127fab87909c574374`
- Retention rows: `0`
- Planned postimages: `15`
- No `docs/agents/specific-repository.md`,
  `docs/agents/repository.md`, or `docs/agents/repository-rules.md` postimage
  exists.

Applying with an all-zero confirmation digest exited `3`, named the exact
approved digest, and left `git status --short` empty.

Applying with the exact digest exited `0`. `apply-result.json` reports:

```json
{
  "state": "verified",
  "verifiedPostimages": 15,
  "warnings": 13
}
```

The warnings identify existing nested carriers as unchanged and requiring
conflict review.

## Generated output

The generated root presents the hierarchy in this order:

1. Universal instructions
2. Context and documentation
3. Spec workflow
4. Autonomous work
5. Stack guidance
6. Surface guidance
7. Optional knowledge sources

No generic or repository-specific carrier exists after apply. The generated
documentation guide contains the five ADR lifecycle states, accepted-only
active semantics, legacy compatibility, nullable deprecation/superseding
fields, all four Findings states, and the complete evidence/routing/addendum
template.

## Formatter, Verification, and persistence

- `rtk bun run format`: exit `1`, because Fluxus has no `format` script.
- `rtk bun run fmt`: exit `0`, 589 files checked.
- `rtk make verify`: exit `0`; formatting, lint, OnionCry, typecheck, and tests
  passed.

The fresh Plan after formatting and Verification is not empty:

```json
{
  "fileChanges": 1,
  "changed": [
    {
      "path": "AGENTS.fc4892a007880b753cc1fde672f4bafc5a866877b2e464a068e71d6abe70d443.md",
      "action": "create"
    }
  ]
}
```

The persistence criterion requires zero file changes or exact idempotent
reapply, so the greenfield journey fails at its final step.
