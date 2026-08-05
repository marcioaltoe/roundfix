# Owned-Skill regeneration journey

Build: `9d7f834e`

Scratch clone: `/private/tmp/roundfix-qa-0067.JWZaVp/repo`

The probe appended one harmless line to the canonical
`.agents/skills/qa-gate/SKILL.md`. Running `rtk make baseline-digests` before
syncing the distributed mirror exited 2 at `TestAuthorialSkillSync`; this was
the planned out-of-order-input probe.

Following the documented owned-Skill path:

1. `rtk make skills-sync` exited 0.
2. `rtk make baseline-digests` exited 0 and reported `changed: true`.
3. The command reported these derived changes:
   - `internal/baseline/testdata/catalog.diagnostics.golden.json`
   - `internal/baseline/testdata/catalog.digest`
   - `internal/baseline/testdata/catalog.normalized.json`
   - `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
   - `internal/baseline/testdata/parity-corpus/v1/manifest.json`
4. `rtk make verify` exited 2: 3,340 tests passed, 7 failed, and 3 skipped.
   `TestBaselinePlanCharacterization` failed because four goldens retained the
   old catalog digest. `TestDeclaredStepRegenerationAndFrozenBoundaries` also
   failed in the dedicated plan-characterization subtest.
5. A second `rtk make baseline-digests` exited 0 with `changed: false` and
   `derived artifacts already match their canonical sources`; it did not repair
   the red Verification gate.

The scratch status after the two sanctioned runs contained the canonical and
distributed Skill files, three setup snapshots, the three catalog artifacts,
and two files below the parity corpus. It contained no plan-characterization
golden because the sanctioned command did not regenerate that dedicated
corpus.

