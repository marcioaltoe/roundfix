---
schema: spec-tasks/v1
spec: 0415-a-skill-the-skill-dispatch-names-joins-the-repository-skill-set
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — A skill the skill dispatch names joins the Repository Skill Set

| id      | title                                                                                               | type    | complexity | needs   |
| ------- | --------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The owner of a Skill Activation requires every skill of its bundle, and one function reads the required skills | backend | high       | —       |
| task_02 | Run the final QA gate                                                                               | qa      | high       | task_01 |

Waves: 1 → task_01 · 2 → task_02

The graph follows `_techspec.md` → Build Order: step 1 is task_01 and step 2, the gate, is task_02. Step 1 stays one Task because each of its parts breaks the tree without the next: the activation rule alone makes the embedded catalog fail to load until `typescript` requires `solid`, the module edit fails the digest and version tests until the repository's commands regenerate the records, and the regenerated catalog fails `TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent` until the Managed Refresh rewrites this repository's Setup Manifest. Moving the three readers onto `Catalog.RequiredSkills` has no Verification that fails before it on its own, since the three unions agree on both sides of the move, so it rides in the same Task, whose Verification proves it by the absence of each hand-written union. The glossary sentence belongs to the same context, so it joins task_01 rather than paying for its own cycle.

## What the Tasks stage settled beyond the TechSpec

- **Subtest names.** The TechSpec names the subtests of `TestAnActivationOwnerRequiresEverySkillOfItsBundle` and one of `TestTheRequiredSkillsCoverTheRenderedDispatch`; task_01 also names `every_module_with_its_dependencies` and `every_built_in_profile`, so that the Verification and the gate can check the two halves of invariant 3 by name.
- **The `./internal/cli` suites.** The Doctor, `skills check`, restore and update tests consume `Catalog.RequiredSkills` and the embedded catalog's digest. task_01 runs them by a name filter in a command that first proves Doctor no longer holds its own union, so the command fails before the work; the gate runs the whole package.
- **The update preview through the binary.** The characterization reproduced `skills check` and restore through the binary; the gate adds the update preview without `--yes`, which writes nothing, and records the row as blocked by the environment only if that preview reaches the network.

## Shared files

No other active Spec is on the default branch. A separate Spec about the knowledge-base search, waiting on a decision, is expected to edit another Baseline module; whichever delivers second regenerates the module version record, the digests, the goldens and the Setup Manifest on its own tree, and task_01 writes "the next version", never a number. No Verification of this Spec reads another Spec's artifact, and the graph needs no `requires`.

## Time in tests

No test of this Spec reads the clock or waits on a timer; every new test reads the embedded catalog, a mutated in-memory copy of it, or files in a temporary repository.
