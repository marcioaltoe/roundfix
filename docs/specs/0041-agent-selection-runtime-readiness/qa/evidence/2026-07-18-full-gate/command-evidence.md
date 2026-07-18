# Command evidence — 2026-07-18 full gate

Build: `f8bcc130b1b48a321a82f15ff577fe259b1a0172`

## Static gate

- `rtk make verify` — exit 0; 1,680 Go tests in 20 packages, 79 setup-context-driven tests, asset catalogs, 14 shipped Roundfix Skill contracts, and CLI build passed.

## Built CLI probes

- Scratch repository: `/tmp/roundfix-qa-0041.ngHidk`.
- `roundfix init --scope project` — exit 0; created a Project Config.
- `roundfix profiles show --json` — exit 0; required Codex profiles resolved to Sol/high with GPT-5.5/xhigh fallback; frontend resolved to Claude Fable/medium with Sol/high fallback; Terra and Luna remained advisory entries.
- Partial-selection probes across `resolve`, `watch`, `implement`, and detached `implement` — exit 2; each reported that `--agent`, `--model`, and `--reasoning-effort` must be supplied together or all omitted. The scratch repository gained no Run artifact or worktree.
- `profiles configure --scope project --file duplicate-profile.yml --yes --json` — exit 2; JSON reported `changed:false` and required one additional distinct authorized and proven Agent Selection. Project Config SHA-256 was `ec654c32b3fac6316b7b40ec431f1a4adefc3f79cec6f401afeffe9c72776696` before and after.
- Built help for `fetch`, `resolve`, `watch`, `implement`, `setup`, `doctor`, `profiles configure`, and `profiles validate` — exit 0 and matched the documented profile-led, exact-proof, read-only, and complete-override contracts.

## Focused current-build regressions

- Adapter identity: 12 agent tests and 6 Setup/Doctor tests passed.
- Capability projection: 18 tests passed; the private-runtime-state guard passed.
- Exact proof/application: 23 tests passed; 8 cleanup, cancellation, timeout, and no-prompt tests passed.
- Shared readiness: 5 coordinator tests and 2 operational/fallback tests passed.
- Setup transaction: 3 config tests and 16 CLI tests passed.
- Profile configuration transaction: 14 CLI tests and 3 config tests passed.
- Doctor: 7 profile-readiness tests and 1 independent-check continuation test passed.
- Override grammar and help: 14 invocation tests, 12 command tests, and 12 help/fetch tests passed.
- Documentation and skill contract: 4 documentation tests and 13 skill tests passed; `make skills-sync-check` and `roundfix skills check` passed.

## Race-focused regressions

- Adapter readiness: 18 tests in 2 packages passed with `-race`.
- Capability acquisition: 19 tests passed with `-race`.
- Exact proof/application: 23 tests passed with `-race`.
- Shared readiness/fallback: 6 tests in 2 packages passed with `-race`.
- Setup/config defaults: 27 tests in 2 packages passed with `-race`.
- Profile configuration proof/no-mutation: 9 tests passed with `-race`.
- Doctor profile readiness/continuation: 7 tests passed with `-race`.

## Scope and traceability checks

- `584df1f..HEAD` changed-file review found no Spec 0036 implementation files and no authorial edits outside repo-owned Roundfix and setup-context skills.
- ADR-0055, the Spec validation record, the profile-preflight dogfood finding, and the Spec 0036 PRD all resolve.
