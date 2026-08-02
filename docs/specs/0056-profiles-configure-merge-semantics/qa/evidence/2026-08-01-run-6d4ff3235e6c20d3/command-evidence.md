# Command evidence — QA rerun 2026-08-01

Build under test: `fa43a7d93ca49b60f6d2fc0aa5c141dff0db472a`.

## PC-01 — Project Constraints, task completion, and changed paths

- `rtk proxy rg -n '^status: completed$' docs/specs/0056-profiles-configure-merge-semantics/task_*.md` — exit 0; all seven Task files report `status: completed`.
- `rtk proxy rg -n 'ADR-0037|ADR-0039|ADR-0049|ADR-0055|ADR-0086|Identifier strategy|Authentication and HTTP|Tooling authority' <PRD> <TechSpec>` — exit 0; both active artifacts account for identifier strategy, authentication/HTTP, the applicable selection/atomicity/removal ADRs, and bounded tooling authority through operative `docs/agents/` sources.
- `rtk git diff-tree --no-commit-id --name-status -r <task-commit>` over `b70fab7`, `487dacc`, `540e763`, `a6a6cd0`, `a556edc`, `e1c3b1e`, `c7e7c33`, `f26be50`, and `fa43a7d` — every implementation Task commit is limited to its Task file and declared implementation/docs roots. No protected tooling path was changed, so the PRD/TechSpec Skill authorization and sanctioned digest fallout remained unused. The remediation planning commit precedes task 07, and task 07 contains only `internal/config/` plus its own Task file.
- `rtk git merge-base --is-ancestor f26be50 HEAD` — exit 0; task 07 is a chronological descendant of the remediation artifact commit.
- `rtk git status --short` before QA writes — clean; after opening this required resumable gate, only this QA Report and its evidence are untracked.

Outcome: pass.

## Focused current-build evidence

- `rtk go test -v ./internal/config -run '^(TestProfilesConfigWriterCharacterization|TestEffectiveChangeSet|TestProfilesConfigureMergePreservesOtherCategories|TestProfilesConfigureRemovalPreservesSpacing)$' -count=1` — exit 0; 32 tests passed.
- Two consecutive `rtk go test ./internal/config -run '^TestProfilesConfigWriterCharacterization$' -count=1` invocations — exit 0; 12 tests passed on each run. SHA-256 for all 11 current goldens was identical before and after both normal runs.
- `rtk go test -v ./internal/cli -run '^(TestProfilesDocumentationContractMatchesPublicGuidance|TestProfilesConfigureDryRunAndFailedConfigurationLeaveBytesUnchanged|TestProfilesConfigureExitCodes|TestProfilesConfigureProofRunsBeforeConfirmationAndWrite|TestProfilesConfigureChangeSummary|TestProfilesConfigureProofScope|TestProfilesShowJSONRendersProfileAndRecommendations|TestProfilesShowTextAndJSONAreByteStableAndConsistent|TestProfilesShowDoesNotMutateConfigOrRunState|TestProfilesValidateDeduplicatesProofsAndReportsEveryReference)$' -count=1` — exit 0; 23 tests passed.
- Task 01's completed Result supplies its named deliberate-drift reproduction: a mutated `non_ascii_values` golden failed with the fixture name and a readable hunk, then only the explicit `-update-profiles-config-goldens` flag restored it. Fresh current Verification and repeated normal corpus runs credit this Task-only destructive probe without mutating a QA build.

These checks cover all 47 Task acceptance criteria when combined with the Task
Results' changed-path evidence and the live public journeys below.

## Live CLI environment

- Production-like entry point: `bin/roundfix` built by SG-01.
- Main scratch Git root: `/private/tmp/roundfix-qa-0056-rerun.CFglvw`, branch `ma/qa-0056-rerun`.
- Alias scratch Git root: `/private/tmp/roundfix-qa-0056-alias-rerun.S531Xy`, branch `ma/qa-0056-alias-rerun`.
- Inputs are preserved under `live-fixtures/`; only scratch Project Configs were written.

## F-RECHECK and T07-01 — prior findings and removal trivia

- The PRD and TechSpec now both name ADR-0055 and ADR-0086, and PRD Core Feature 6 now matches the written-category proof scope exposed by the TechSpec, Task 04, CLI help, and user guide.
- Public `profiles configure --scope project --remove frontend --yes --json` on the five-category Project Config exited 0 with `changed:true`. `git diff --no-index` showed exactly the three-line `frontend` mapping deletion; the blank separator before `watch:` remained.
- Fresh `profiles show --category frontend --json` returned `source:"user"`, independently proving the Project Config no longer supplied `frontend`; the persisted Project Config reparsed.
- Public removal of middle `backend` from the adjacent-trivia fixture exited 0. The exact diff removed its leading comment, mapping, and trailing blank line only; the preceding `general` blank line and following `qa` comment remained. Fresh `profiles show --category qa --json` returned `source:"project"` with the original tuple.
- The focused removal-spacing and corpus tests passed without modifying any earlier golden; task 07 added two new input/golden pairs and changed no pre-existing corpus file in commit `fa43a7d`.

Outcome: prior F-001, F-002, and F-003 are resolved; T07-01 passes.

## US-02 and US-04 — summary and explicit removal

- Text dry-run with `--remove frontend --remove data --dry-run --yes` exited 0 and printed exactly `removed: frontend` then `removed: data`; SHA-256 remained `21f843e302d7a3d2ab4e0c19532dbaf2d8b02f8b614b8884698049cdd943477f`.
- JSON dry-run exited 0 with the same ordered `changes` array and unchanged SHA-256.
- Fragment plus `--remove frontend` exited 2, named `profiles.frontend cannot be both configured and removed`, returned `changed:false`/`refused:false`, and preserved the original hash.
- Absent `--remove data --yes --json` exited 0 with `changed:false`; present `frontend` removal exited 0 with `changed:true`; repeating it exited 0 with `changed:false` and preserved the post-removal hash.
- Focused Effective Change Set and summary tests equivalently cover added/replaced classification and repeatable deterministic derivation while the live proof environment blocks fragment writes.

Outcome: US-04 passes. US-02 is environment-blocked only for live add/replace, with equivalent current-build evidence.

## US-03 — refusal and exit codes

- Non-interactive EOF `profiles configure --scope project --remove qa --json` exited 1 and preserved the original hash.
- Captured stdout contained only `roundfix/profiles-configure/v1` JSON with `changed:false` and `refused:true`; stderr contained the preview, prompt, and `Profile configuration unchanged: confirmation declined` diagnostic.
- The live fragment/removal validation failure exited 2. Live applied removal, absent/repeated no-op, and dry-run each exited 0. The focused exit-code matrix independently covered all cases and the additive false marker on non-refusal results.

Outcome: pass.

## T03-02 — dangling alias

- Public removal of `backend`, whose model owned `&backend_model` while `defaults.artifact_dir` retained `*backend_model`, exited 2 with `yaml: unknown anchor 'backend_model' referenced`.
- The candidate reported `changed:false`; SHA-256 stayed `4986f4d9036599c3876c8e2a97101f6753148d7e5d7e4d668ce967153ce5e9b5`.
- Fresh `profiles show --category backend --json` returned the original Project Config profile.

Outcome: pass.

## Environment block — proof-gated add/replace journeys

Public replacement command:

`bin/roundfix profiles configure --scope project --file backend-fragment.yml --yes --json`

- Exit 2 before confirmation or persistence.
- JSON classified the requested category as `replaced`, returned `changed:false`/`refused:false`, and named `selection_rejected` because `/Users/marcio/.local/bin/codex` has an invalid or missing code signature.
- Target SHA-256 remained `21f843e302d7a3d2ab4e0c19532dbaf2d8b02f8b614b8884698049cdd943477f`.
- The product's unblocking action is to reinstall the officially signed Codex CLI into `~/.local/bin`, set `CODEX_PATH` to it, then rerun `roundfix profiles validate` and these journeys.
- Equivalent evidence: the 32 passing config tests cover four-category byte identity, single-value diff, atomic replacement, add, empty/sectionless files, removal, aliases, and corpus parseability. The 23 passing CLI tests cover add/replace/remove summaries, proof before confirmation/write, written-category proof scope, dry-run/no-write, exit codes, show/validate canaries, and docs synchronization.

Affected rows: US-01, US-02, US-05, and CF-06. Each has equivalent fresh
current-build evidence, so this environment cause does not cap the verdict.

## T06-01 and NG-01 — docs and Non-Goals

- `bin/roundfix profiles configure --help` and `docs/user-guide/commands.md` agree on category merge, repeatable explicit removal, written-category proof, dry-run, and text/JSON behavior; the focused documentation synchronization test passed.
- `rtk git diff --name-only cdbed1c..HEAD -- internal/cli/profiles_show.go internal/cli/profiles_validate.go internal/config/profile_resolver.go` produced no path. Public `profiles show` continued to read scratch Project Configs. Public `profiles validate` reached its existing proof boundary and reported the same signed-Codex environment cause without mutation.
- Atomic replacement, unchanged resolution/fallback behavior, and no general YAML pass are covered by the focused exact-output tests and the two live removal diffs.

Outcome: T06-01 and NG-01 pass.

## PR-01 — Pull Request environment

The QA prompt states that no Pull Request is open and Pull Request journeys are
environment-blocked. No lookup was attempted against the never-pushed Run
Worktree branch. Equivalent supervised evidence is the ancestry/scope audit,
full Verification, public built-CLI flows, focused current-build tests, and
resolved prior-finding reproductions. Unblock by opening a Pull Request on
`ma/profiles-configure-merge-semantics`.

## SG-01 — Repository Verification

- `rtk make verify` — exit 0. It reported 3,002 passing Go tests across 24
  packages, 4 passing focused Skill tests, a passing Repository Skill Set
  check, and a successful production binary build at `bin/roundfix` with
  `BuildCommit=fa43a7d-dirty` (the expected QA-artifact delta).

Outcome: pass.
