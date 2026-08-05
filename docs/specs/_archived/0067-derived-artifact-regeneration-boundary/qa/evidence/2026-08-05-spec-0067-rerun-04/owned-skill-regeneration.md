# Owned-Skill regeneration evidence

Build: `e91bf4088b7547ab1f1c4a15c78d1427e769f032`.

QA cloned the build locally to
`/private/tmp/spec0067-qa-rerun04.ut3PMq/repo`, added a scratch-only comment to
the owned `.agents/skills/qa-gate/SKILL.md`, and used the public maintainer
commands below. Nothing from the scratch clone is part of the Run Worktree.

1. `rtk make skills-sync` exited 0. Fresh status showed only the canonical
   owned Skill and its `skills/qa-gate/` mirror changed.
2. The first `rtk make baseline-digests` exited 0 with `changed:true`. Its
   changed-artifact list included the three setup assets, catalog diagnostics,
   catalog digest and normalized catalog, both parity exceptions, and four
   plan-characterization goldens.
3. The first unpiped `rtk make verify` exited 0: 3,382 Go tests passed across
   26 packages; the isolated corpus budget, four Skill tests, Repository Skill
   Set check, build, and Spec consistency checks also passed.
4. The second `rtk make baseline-digests` exited 0 with `no changes` and
   `changed:false`.
5. The second unpiped `rtk make verify` exited 0 with the same gate counts.

The scratch changed-path read after the first regeneration showed the owned
Skill pair and sanctioned derived outputs only. Under the parity corpus, only
`v1/manifest.json` and `v1/fixtures/asset-sync.json` changed; all 15 frozen
paths remained untouched.
