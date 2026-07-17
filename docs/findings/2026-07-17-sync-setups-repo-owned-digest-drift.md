---
status: pending
created_at: 2026-07-17
updated_at: 2026-07-17
---

# `sync-setups` replaces Roundfix-owned skill digests

The `sync-setups` command can write an incorrect digest for a skill whose
`source.type` is `repo`. This makes compliant repositories fail audit with
`skills.required.drift` and blocks Roundfix's `make verify`.

## Expected and current behavior

Expected: when Roundfix owns the skill, `sync-setups` preserves the digest from
Roundfix's own authoritative source.

Actual: when another local checkout contains the same skill path, `sync-setups`
calculates the digest from that external checkout and writes it into the bundled
snapshots. The next audit compares that digest with the skill installed by
Roundfix and reports drift.

## Reproduction

1. Keep a skill with `source.type: repo` in a setup snapshot.
2. Point `--source-dir` at an external checkout containing the same skill path
   with different content.
3. Run:

   ```sh
   PYTHONDONTWRITEBYTECODE=1 python3 \
     .agents/skills/setup-context-driven/scripts/context_setup.py \
     sync-setups \
     --source-dir /Users/marcio/dev/skills/setups \
     --format json
   ```

4. Run `make verify` or audit a compliant repository fixture.

`sync-setups` changes the Roundfix-owned skill's `contentDigest`. Audit ends
with `skills.required.drift`; in the observed case, the
`setup-context-driven` Python suite reported 30 failures.

## Code evidence

`normalize_source_skill` in
`.agents/skills/setup-context-driven/scripts/context_setup.py` resolves
`resolve_canonical_skill_file(source_dir, raw_path)` first when the source data
does not provide a digest. The function does not treat `source.type: repo`
differently. A matching path in the checkout named by `source_dir` therefore
takes precedence over the digest already recorded for the Roundfix-owned skill.

A minimal regression test confirms the behavior:

1. clone the setup assets into a temporary directory;
2. select a skill with `source.type: repo`;
3. create a conflicting external `SKILL.md` at the same path;
4. run `sync_setup_snapshots`;
5. observe that the persisted digest becomes the external file's digest.

## Impact

- `make verify` can fail after an apparently successful synchronization.
- Snapshots no longer represent the source declared by their `source` field.
- The result depends on content outside the Roundfix repository.
- `make skills-sync` propagates the incorrect digest into the embedded
  `skills/` copy.

## Suggestion for the spec

- Define source precedence explicitly from `source.type`.
- For `source.type: repo`, use only Roundfix's authoritative source or preserve
  the validated current digest; do not consult the external setup checkout.
- Keep the current external-skill behavior only when the declared source
  requires that checkout.
- Add a regression test with conflicting external content.
- Validate both `.agents/skills/setup-context-driven` and the embedded
  `skills/setup-context-driven` copy.
- Require `setup-context-check` and the full `make verify` as completion gates.

## Changes discarded after diagnosis

The investigation generated snapshots and a local regression test. Those
changes are not part of this finding and were discarded at the user's request.
This document preserves the diagnosis so the correction can be incorporated
into the existing Roundfix spec.

## Planned resolution

Spec [0036 Doctor skill readiness](../specs/0036-doctor-skill-readiness/_prd.md)
owns this correction. Its first Task requires source-type-aware digest
precedence, the conflicting external-content regression, validation of both
Roundfix-owned skill copies, and the full verification gate. This finding
remains pending until that Task and the Spec QA pass.
