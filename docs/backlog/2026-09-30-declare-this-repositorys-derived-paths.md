---
type: chore
status: open
created: 2026-09-30
spec: null
reason: null
---

# Declare this repository's derived paths once the release carries the key

## Opportunity

Spec 0201 adds the Project Config key `delivery.derived_paths`, which lets the queue owner resolve a Pull Request conflict confined to regenerable files. A binary older than that Spec rejects the unknown key, so this repository cannot declare its own derived paths until a release that carries it is installed.

## Value

Without the declaration, a conflict on the Baseline digest pin, catalog snapshots, plan goldens, `docs/agents/setup-context.json` or `skills/testdata/owned-skill-versions.json` still parks as `pull-request-conflict` here. In the v0.22.0 cycle, two items needed a manual merge and regeneration for exactly that reason.

## Shape

After the release that ships Spec 0201, add `delivery.derived_paths` to `.roundfixrc.yml` for those paths with their sanctioned commands (`make baseline-digests`, the managed refresh, `-record-skill-versions`), under a grant that names `.roundfixrc.yml`.
