---
type: fix
status: promoted
created: 2026-10-12
spec: 0257-the-baseline-respects-the-repository
---

# A skill newer than the Setup Snapshot is reported as trailing, and the remedy downgrades it

## Symptom

In the vortex, on 2026-10-09, with v0.65.0: `baseline update` applied cleanly.
Then `make skills-update` installed the external skills from
`marcioaltoe/skills` `5290708`, a commit made one hour after the v0.65.0
release. From then on `roundfix doctor` printed
`DR-SKILL-TRAILS-SNAPSHOT: trails the Setup Snapshot: context7-cli, diagnosing-bugs, grilling, handoff, tdd; next: roundfix baseline update`,
and the update preview listed the five skills as drifted, with
`restore with roundfix baseline update --yes`.

`TrailingSetupSkills` (`internal/baseline/skills_trailing.go:32`) flags a skill
when `portableRestoreDigest(files) != contract.TreeDigest`. That is an
equality check, not an ordering. A skill newer than the embedded snapshot
reads as "trails", and the suggested remedy would restore the older copy over
the newer one.

## Where

`roundfix doctor` and the `baseline update` preview: `internal/baseline/skills_trailing.go`.

## Expected

- Tell "behind the snapshot" apart from "different from the snapshot". For
  example, a copy whose digest matches `skills-lock.json` but not the snapshot
  is ahead or unknown, not trailing.
- Never propose restoring the snapshot over a copy that may be newer.

## Evidence

- Secondbrain `inbox/roundfix/_triaged/2026-10-09-dr-skill-trails-snapshot-acusa-skill-mais-nova-que-o-snapshot.md`
  (vortex, 2026-10-09).
