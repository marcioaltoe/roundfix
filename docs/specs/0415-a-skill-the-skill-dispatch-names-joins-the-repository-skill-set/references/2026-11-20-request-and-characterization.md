# Request and characterization of Spec 0415 — 2026-11-20

The request is copied from the operator's brief for this Spec. Everything under
"Measured" was run by the planning author on 2026-11-20 on the branch
`fix/dispatch-named-skills-join-the-skill-set` at `851024fa` (v0.152.0). Probes
ran as `go test -overlay` files kept outside the repository and wrote nothing to
it. The prototype ran on a scratch clone outside the repository, which was
removed. No file of another repository is copied here.

## 1. The request

The maintainer's request reached the operator through ai-memory message
`01a1f929` on 2026-11-20, about 10:40 BRT. The delivery of a Spec in an adopter
stopped because the Baseline's skill dispatch demands the skill `solid` for
production code and the skill was never installed there. The maintainer's
words: "se é parte da baseline deveria ter sido verificado no update e fazer
parte do setup no projeto skills."

The operator's reading of the ask: the Baseline update, and the setup of a new
repository, check that every skill the skill dispatch demands is installed,
and install it or fail listing what is missing, before a Run dispatches a Task.

The project that owns the upstream skills catalog answered in message
`01a1f92a` that every built-in profile's skills are inside the Setup Snapshots
and the catalog presets, and suspected the update. The catalog does hold
`solid`; the defect is in Roundfix.

The adopter uses a repository-owned Baseline Profile with the modules `core`,
`context-workflow`, `typescript`, `bun`, `autonomous-work`, `spec-workflow`,
`secondbrain` and `repository-extension`: TypeScript without the `frontend`
module. Another session installed `solid` there; this Spec changes nothing in
that repository.

## 2. Measured cause

| Fact | Where |
| --- | --- |
| `bundle.production-code` holds `coding-guidelines`, `clean-code`, `solid` | `internal/baseline/assets/skill-activations.json:54` |
| `trigger.production-code` has the owner `typescript` | `internal/baseline/assets/skill-activations.json:113-116` |
| The `typescript` module's `requiredSkills` lacks `solid`; it depends on `core` only | `internal/baseline/assets/modules/typescript.json:7-28` |
| `solid` is required only by the `frontend` module, which also dispatches it | `internal/baseline/assets/modules/frontend.json:28`, `:184` |
| A repository profile's required skills are the union of its modules' `requiredSkills` | `internal/baseline/skills_restore.go:641-646` |
| A built-in profile's required skills are the same union, written a second time | `internal/baseline/skills_restore.go:786-823` |
| Doctor, the update preview and `skills check` read the same union a third time, from the Setup Manifest's modules | `internal/cli/doctor.go:570-632` |
| The guide renders every bundle of an activation whose owner is selected, whatever the modules require | `internal/baseline/plan.go:2186-2218` |
| The catalog checks that a profile's dispatched and bundled skills are in its Setup Snapshot, never that a module requires them | `internal/baseline/catalog_validate.go:448-466` |
| A module requires every skill its own `skillDispatch` names, held by a test | `internal/baseline/upstream_skill_names_test.go:75-90`, `:179` |
| No rule holds an activation's bundle to its owner | `internal/baseline/catalog_validate.go:1405-1433` |

Every built-in profile that selects `typescript` also selects `frontend`, so
the gap never showed in a built-in profile.

## 3. Reproduction

### In the package

A repository profile made by `InitCustomProfile` from
`standard-typescript-monorepo` and cut to the eight modules of section 1
resolves without error. On `851024fa`:

- `loadSkillSnapshotProfile` returns no `solid` in `RequiredSkills` and no
  contract for it in `Skills`;
- `profileNamedSkills` for the same modules names exactly one skill outside
  the required set: `solid`;
- `renderSkillDispatch` prints
  `` - `bundle.production-code`: `coding-guidelines`, `clean-code`, `solid` ``;
- `selectRestoreSkills` with `solid` fails with "Selected skills are not
  external members of this profile: solid."

### Through the released binary

In a new temporary Git repository: `roundfix baseline profile init --id
ts-backend --from standard-typescript-monorepo`; the profile's `modules` cut to
the eight of section 1; `docs/agents/setup-context.json` holding
`{"profile": "ts-backend", "modules": [the eight]}`; `roundfix skills install
--target project`; and a `skills-lock.json` holding `{"version": 1, "skills":
{}}`.

| Command | Roundfix 0.152.0 | Scratch build with the fix of section 5 |
| --- | --- | --- |
| `roundfix baseline profile validate ts-backend` | `valid (repository)`, exit 0 | the same |
| `roundfix baseline skills restore --profile ts-backend --skill solid --source-dir <absent directory>` | `restore.skill-invalid: Selected skills are not external members of this profile: solid.`, exit 2 | `restore.source-dir-invalid: Offline Git object store is not a directory: <path>.`, exit 2: the selection is accepted |
| `roundfix skills check` | exit 1; 34 lines "required by the Baseline and not installed"; none names `solid` | exit 1; 35 such lines, one of them `.agents/skills/solid: required by the Baseline and not installed; restore with roundfix baseline skills restore --profile ts-backend --skill solid` |

## 4. The sweep for other gaps

A probe walked all 11 activations and all 17 modules of the embedded catalog.

- For each activation it joined the `requiredSkills` of the owner and of every
  module the owner reaches through `dependsOn`, and looked for each skill of
  the bundle there. One gap: `trigger.production-code`, owner `typescript`,
  skill `solid`.
- For each module it looked for each skill of its `skillDispatch` in its own
  `requiredSkills`. No gap.

No skill is in the `requiredSkills` of two modules today, and no rule forbids
it: `catalog.module.requiredSkills.duplicate` looks inside one module only.

## 5. Prototype of the fix

On the scratch clone, `"solid"` was added to the `typescript` module's
`requiredSkills`, after `"onioncry"`.

| Step | Files written |
| --- | --- |
| `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1` | `internal/baseline/assets/modules/typescript.json` (version 10 to the next free one), `internal/baseline/module-versions.json` |
| `make baseline-digests` | `internal/baseline/testdata/catalog.digest`, `internal/baseline/testdata/catalog.normalized.json`, and under `internal/baseline/testdata/plan-characterization/`: `advisory-only-divergences.golden.json`, `clean-adoption.golden.json`, `idempotent-replan-after-verified-apply.golden.json`, `same-baseline-changed-profile-and-catalog-digests.golden.json` |
| `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text` | `docs/agents/setup-context.json` (`catalogDigest` and `roundfixVersion`) |

Without the last step, `TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent`
in `internal/cli` (build tag `docscontract`) fails.

Results on the scratch clone:

- the reproduction of section 3 returns `solid` required, with its contract,
  and the restore selection is accepted;
- `go test ./internal/baseline -count=1` passes;
- the required skills of every built-in profile are the same names before and
  after:

| Built-in profile | Required skills | With `solid` |
| --- | --- | --- |
| `go-cli-tui` | 41 | no |
| `go-cli-typescript-monorepo` | 105 | yes |
| `python` | 38 | no |
| `rust-cli` | 34 | no |
| `standard-typescript-monorepo` | 95 | yes |

## 6. Governed paths

`speccheck.GovernedPath`, in a `go test -overlay` probe of `internal/speccheck`,
read every file this Spec's design names and five controls.

- Governed and changed by this Spec: `internal/baseline/assets/modules/typescript.json`
  and `docs/agents/setup-context.json`.
- Governed and not changed: `internal/baseline/assets/skill-activations.json`,
  `internal/baseline/assets/modules/frontend.json`, the Roundfix skill's
  `references/runtime.md`, and the controls `Makefile`,
  `.github/workflows/ci.yml`, `go.mod` and `.roundfixrc.yml`.
- Not governed: `internal/baseline/module-versions.json`, the six files under
  `internal/baseline/testdata/`, the Go files and tests under
  `internal/baseline` and `internal/cli`, `CONTEXT.md`, the ADRs and this
  Spec's planning files.

## 7. Sources consulted

- Secondbrain, through the Indexus search of 2026-11-20:
  `projects/sql-library/mirror/docs/adr/0004-o-baseline-foi-adotado-com-profile-proprio.md`.
  A second adopter made its own profile by `profile init --from` a built-in
  profile and deleting modules. It confirms that a repository profile is
  usually a subset of a built-in one, so a rule that holds only for the
  built-in compositions does not protect it. `wiki/index.md` and the Roundfix
  entity pages held nothing on this defect. `qmd query` failed three times
  with `SQLITE_BUSY: database is locked`, so the qmd consultation is recorded
  as attempted and unavailable.
- Bazel documentation, "Dependencies", https://bazel.build/concepts/dependencies,
  read on 2026-11-20: "For correct builds, the graph of actual dependencies A
  must be a subgraph of the graph of declared dependencies D." The page's
  example is this defect: a target that uses what it never declared keeps
  building while another target brings the dependency in, and breaks when that
  target leaves. It supports declaring the skill at the owner and checking
  the declaration, and it challenges relying on the `frontend` module.
- Bazel blog, "Strict Java Deps and `unused_deps`",
  https://blog.bazel.build/2017/06/28/sjd-unused_deps.html: the check names
  the missing direct dependency to add. It supports a diagnostic that names
  the activation and the skill.
- Go modules reference, https://go.dev/ref/mod: `go mod tidy` "adds any
  missing module requirements necessary to build the current module's
  packages", and most commands report an error when `go.mod` does not reflect
  reality. It is the alternative this Spec rejects for now: computing the
  required list from the references instead of declaring it. Inference, not
  published evidence: Roundfix can check the declaration completely, because
  the catalog enumerates every reference.
- npm documentation, `npm ci`, https://docs.npmjs.com/cli/v11/commands/npm-ci/:
  it exits with an error when the lock and the manifest disagree and never
  rewrites either. It supports refusing at validation instead of repairing in
  silence.
