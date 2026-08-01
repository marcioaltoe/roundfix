# Documentation probes

Build: `e45dd37d2f2ced6dcaa3533fcea939a867b3ea6c`

`rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-07-31-qa-01/docs_scope_probe.rb` exited 0. Fresh assertions derived the launcher plus all five platform coordinates from `dist/npm/platforms.json` and confirmed:

- the ordered trusted-publisher setup binds every coordinate to owner `marcioaltoe`, repository `roundfix`, and workflow `release.yml`;
- npm validates each binding only when `npm publish` attempts OIDC;
- the publish-free rehearsal, fallback purpose, external switch, empty-record exit evidence, repository variable/secret removals, and fallback-branch removal remain documented;
- the closing procedure disallows token publication per package for all six coordinates only after the empty fallback record and independently confirms the registry settings;
- `CONTEXT.md` defines Release Set and Publication Preflight with the current workflow meaning;
- the preflight has no cooldown retry and the Release Plan Command remains unchanged.

`rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-08-01-qa-02/failure_vocabulary_probe.rb` exited 0. The workflow and runbook both enumerate exactly `identity`, `publish`, `registry`, `runtime`, and `undetermined`; every runbook row includes its emitting stage, meaning, and first recovery action.
