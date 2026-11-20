---
created: 2026-11-20
---
# A `baseline update` fixture needs a Setup Manifest the binary wrote, not a hand-written profile and modules object

- Context: the QA gate of Spec 0415 previewed `roundfix baseline update` in a temporary repository with a repository profile without `frontend`, to see `solid` named as a Missing Required Skill.
- Tried: the plan built the fixture with a `docs/agents/setup-context.json` holding only `profile` and `modules`, the shape `roundfix skills check` and restore accept.
- Went wrong: the preview exited 2 with `Category: manifest` and `Skills: not run`, so the QA row failed before skill inspection.
- Cause: the update path parses the complete versioned Setup Manifest (schema, generator, digests, decisions, managed artifacts) before it inspects skills, in `parseManagedSetupManifest` of `internal/baseline/classification.go`; `skills check` reads only the lighter object.
- Solved by: adopting the Baseline in the fixture through `roundfix baseline plan` and `roundfix baseline apply`, with the repository profile's unused stack capabilities and their decisions removed and fixture `SKILL.md` files for the two universal installed-skill capabilities, so the binary writes the manifest.
- Applies when: a test or QA row runs `roundfix baseline update` or another command that reads the managed Setup Manifest in a fixture repository.
- Does not apply when: the row runs only `roundfix skills check`, `baseline profile validate` or `baseline skills restore`, which accept the two-field object.
- Links: Spec 0415 task_02 requirement 5; QA report 2026-11-20 finding F-01; Correction C1.
