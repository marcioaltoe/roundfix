# Command evidence — Spec 0056 QA — 2026-08-01

Build: `e1c3b1e`

## PC-01 — Project Constraints and changed-file audit — pass

- `_prd.md` and `_techspec.md` explicitly account for identifier strategy,
  authentication/HTTP, active ADR obligations, and tooling authority with
  operative `docs/agents/` citations.
- Tasks `task_01.md` through `task_06.md` all carry `status: completed`.
- Accepted obligations inspected: ADR-0037, ADR-0039, ADR-0049, and ADR-0086.
- Express authorization commit `6096f1d` predates TechSpec/ADR commit
  `a138577` and all Task commits.
- `git diff-tree --no-commit-id --name-only -r` over `b70fab7`, `487dacc`,
  `540e763`, `a6a6cd0`, `a556edc`, and `e1c3b1e` found only Task files,
  implementation/tests/corpus files, and `docs/user-guide/commands.md`.
  None of the protected Skill or derived digest paths was changed, so the
  bounded authorization was unused and there is no prerequisite/consequent
  protected-tooling choreography to audit.
- Current delta before execution contained only the untracked QA directory.
  Git emitted a non-fatal fsmonitor IPC diagnostic but returned the requested
  path evidence.

## SG-01 — Repository Verification — pass

Command: `rtk make verify`

Exit: `0`

Observed summary:

```text
Go test: 2995 passed in 24 packages
Go test: 4 passed in 1 packages
Roundfix skill check passed: roundfix, write-idea, write-prd, write-techspec,
write-tasks, setup-context-driven, implement-task, implement-spec,
brainstorming, council, business-analyst, archive-spec, qa-gate, evidence-gate
go build -buildvcs=false ... -o bin/roundfix ./cmd/roundfix
```

## T01-01 — Characterization corpus — pass

Commands:

- `go test -v ./internal/config -run '^(TestProfilesConfigWriterCharacterization|TestEffectiveChangeSet|TestProfilesConfigureMergePreservesOtherCategories)$' -count=1`
- Two consecutive `go test ./internal/config -run '^TestProfilesConfigWriterCharacterization$' -count=1`
- SHA-1 of every `*.golden.yml` before and after both runs.

Exit: `0` for every test command. The verbose run named and passed all nine
required shapes: five categories with Fallback Chains, comments, non-default
indentation, non-alphabetical order, multiline scalar, YAML anchor/alias,
non-ASCII values, empty file, and no `profiles` section. The nine hashes were
identical before and after both normal runs.

The current test declares only `-update-profiles-config-goldens` as the write
gate and reports mismatch as `characterization config <fixture> changed` plus
`profilesConfigGoldenDiff`, which emits `--- golden`, `+++ actual`, and an
`@@` changed-region hunk. The named deliberate-drift reproduction in
`task_01.md` is credited after current `make verify` and the current normal
comparison both passed; QA did not mutate repository goldens.

## Live CLI environment

- Binary: the `bin/roundfix` built by SG-01 from `e1c3b1e`.
- Main scratch Git root: `/private/tmp/roundfix-qa-0056-20260801`, branch
  `ma/qa-0056`; five explicitly configured categories, four-space indentation,
  comments, non-alphabetical order, and a multiline scalar.
- Alias scratch Git root: `/private/tmp/roundfix-qa-0056-alias-20260801`, branch
  `ma/qa-0056-alias`.
- No repository or user config was intentionally mutated; all writes targeted
  the scratch Project Configs.

## Environment block — proof-gated add/replace journeys

Command:

`bin/roundfix profiles configure --scope project --file backend-fragment.yml --yes --json`

Exit: `2` before persistence. JSON reported `changed:false`,
`refused:false`, `kind:"replaced"`, and classification
`selection_rejected`. The exact cause was:

```text
codex runtime is not safe for acpx launch:
/Users/marcio/.local/bin/codex has an invalid or missing code signature
```

The target hash remained
`d6be9f4c8c5776a3eb4f7a7cea19ef71194335bb`. Unblocking action from the
product diagnostic: reinstall Codex with the official curl installer into
`~/.local/bin`, set `CODEX_PATH` to that signed binary, then rerun
`roundfix profiles validate` and the proof-gated QA rows. QA did not bypass or
retry-loop this safety refusal.

Equivalent current-build evidence:

- `TestProfilesConfigureMergePreservesOtherCategories` passed replacement,
  atomicity, add, empty/sectionless, and alias cases.
- `TestProfilesConfigureChangeSummary` passed add/replace/remove text/JSON and
  dry-run equivalence.
- `TestProfilesConfigureProofScope` passed proof-before-confirmation and the
  implemented written-category proof scope.
- `TestProfilesConfigureExitCodes` passed applied, no-op, dry-run, refusal,
  validation, and proof-failure paths.

## US-02 — Per-category summary — blocked environment with equivalent evidence

Runnable removal dry-run commands exited `0`:

- Text: exactly `removed: frontend` then `removed: data`; no untouched
  category appeared.
- JSON: `changes` was exactly
  `[{"category":"frontend","kind":"removed"},{"category":"data","kind":"removed"}]`.
- Both reported the same ordered categories, `changed:false`, and left the
  original SHA-1 unchanged.

The focused current-build summary test supplied equivalent add and replace
evidence because those live classifications could not pass the signed-runtime
precondition.

## US-03 — Refusal and exit-code contract — pass

Live non-interactive command:

`bin/roundfix profiles configure --scope project --remove qa --json`

- Captured exit: `1`.
- Stdout contained only the deterministic `roundfix/profiles-configure/v1`
  JSON with `changed:false`, `refused:true`, and one `removed:qa` change.
- Stderr contained the preview, confirmation prompt, and
  `Profile configuration unchanged: confirmation declined` diagnostic.
- The Project Config hash stayed
  `6e2aed3aa798dda8cea220c34cb17f016e018d17`.
- An explicit `n` confirmation also declined without changing bytes.

The focused `TestProfilesConfigureExitCodes` independently passed declined,
EOF, validation code 2, applied/no-op/dry-run code 0, stream separation, and
the additive refusal field.

## US-04 — Explicit removal — fail (F-003)

- Fragment plus `--remove frontend` exited `2`, named
  `profiles.frontend cannot be both configured and removed`, reported
  `changed:false`/`refused:false`, and preserved the target hash.
- Absent `--remove data --yes --json` exited `0`, reported one removed change
  with `changed:false`, and preserved the hash.
- Present `--remove frontend --yes --json` exited `0`, reported
  `changed:true`, and a fresh `profiles show --category frontend --json`
  confirmed Project Config no longer supplied that category.
- Repeating the same removal exited `0`, reported `changed:false`, and kept the
  post-removal hash unchanged.

Independent `git diff --no-index` against the saved original showed that the
writer removed the `frontend` mapping and also deleted the unrelated blank
line between the closing `profiles` entry and `defaults`. That extra diff is
F-003.

## T03-02 — Dangling alias failure — pass

The alias scratch config defined `&backend_model` inside `profiles.backend`
and referenced it from `defaults.artifact_dir`. Public
`profiles configure --scope project --remove backend --yes --json` exited `2`
with `yaml: unknown anchor 'backend_model' referenced`, reported
`changed:false`, and preserved SHA-1
`9ebc14e42f18f6fd589be4c0b69cb5fb485d5c4a`. A fresh `profiles show` still
read the original anchored backend selection.

## Focused current-build CLI/backend evidence

Both commands exited `0`:

- Config: verbose focused run passed the 9 corpus cases, all 8 merge scenarios,
  and all 5 Effective Change Set scenarios.
- CLI: verbose focused run passed documentation synchronization, the 9
  exit-code subcases, 4 change-summary cases, and proof scope.

Additional non-goal canaries passed:

- `go test ./internal/cli -run 'Profiles(Show|Validate)' -count=1`
- Public `profiles show --category backend --json` before and after removal.
- `git diff cdbed1c..HEAD` showed no change in
  `profiles_show.go`, `profiles_validate.go`, or `profile_resolver.go`.

Live `profiles validate --category backend --json` reached its existing exact
proof boundary and returned the same invalid-code-signature environment cause;
the fallback proof remained pending. This is environment evidence, not a
product regression.

## T06-01 — Documentation contract — pass

`bin/roundfix profiles configure --help` and
`docs/user-guide/commands.md` both expose merge-by-category, atomic named
replacement, omission preservation, repeatable `--remove`, written-category
proof, dry-run, and the text/JSON surface. The command reference accurately
lists added/replaced/removed and exit 0/1/2 situations. The focused
`TestProfilesDocumentationContractMatchesPublicGuidance` passed.

## Contract findings from artifact comparison

- F-001: PRD/TechSpec Project Constraints omit accepted ADR-0055, which
  explicitly governs proof before configuration mutation; the PRD also omits
  accepted ADR-0086 even though removal is this Spec's first-class behavior.
- F-002: PRD Core Feature 6 promises proof of every distinct tuple in the
  resulting effective map. TechSpec Decisions, task 04, implementation tests,
  help, and user docs instead intentionally prove only categories written by
  the operation. The contradiction is explicit and current.
- F-003: live removal rewrites one unrelated blank separator, violating the
  minimal-diff/no-general-formatting promises.

## PR-01 — blocked environment with equivalent evidence

Prompt fact: no Pull Request is open; Pull Request journeys are
`blocked (environment: no open Pull Request)`. No `gh` lookup was attempted
against the per-Run branch. Equivalent local supervised evidence consists of
the Task commit ancestry audit, `rtk make verify`, public built-CLI flows,
focused regressions, and current diff/scope inspection. Unblocking requires an
Open Pull Request on `ma/profiles-configure-merge-semantics`.
