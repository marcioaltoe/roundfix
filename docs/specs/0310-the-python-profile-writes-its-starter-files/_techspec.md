---
spec: 0310-the-python-profile-writes-its-starter-files
prd: _prd.md
created: 2026-10-23
---

# The Python profile writes its Starter Files — Technical Spec

## Executive Summary

Three changes. First, the Baseline catalog learns Starter File declarations:
a module field `starterFiles` and a template kind `starter-file`, validated
when the catalog loads, and the `python` module declares a Makefile, a Project
Config and a CI workflow whose bytes live under
`internal/baseline/assets/templates/starters/python/`. Second, a Baseline
operation previews and writes a profile's Starter Files: only absent paths,
only with the current Plan Digest, exclusively and with removal of what it
created on failure; profile alignment names it in the next action of an
undeclared Verification command for a profile that has Starter Files. Third,
the CLI exposes it as `roundfix baseline starters`, with its help, text and
JSON output, command guide, Roundfix skill reference, Skill Coverage Map entry
and the glossary term **Starter File**.

The trade-off accepted: a Starter File is never managed after it is written,
so an improved template reaches only repositories that have not written that
file yet.

## Project Constraints

- Identifier strategy: not applicable — no project-owned identifier or record
  schema is minted; catalog identifiers follow the dotted naming and the JSON
  schema name follows `roundfix/baseline-*` payloads. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files only; tests use
  temporary repositories and homes. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Orders 1-3 follow ADR-0323 for
  the Starter Files and their command, ADR-0068 for the confirmation-gated
  shape, ADR-0075 for the divergence whose next action changes, ADR-0219 and
  ADR-0241 for repository-owned drafts that keep the module's Starter Files,
  and ADR-0322 for the profile and module they extend. All six are
  named with their titles in `_prd.md` → Project Constraints. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 ("Preset de baseline Python"; authorization "same
  as cycle AO", carrying the cycle AJ grant "Mexer em Governed Paths
  declarados"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0310-the-python-profile-writes-its-starter-files/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `docs/agents/setup-context.json`,
  `internal/baseline/assets/modules/python.json`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/starters/python/Makefile`,
  `internal/baseline/assets/templates/starters/python/roundfixrc.yml`,
  `internal/baseline/assets/templates/starters/python/verify.yml`,
  `internal/cli/baseline_documentation_contract_test.go`,
  `internal/docscontract/publicdocs_test.go`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-23 at `dde9d5f5`, and on a prototype with
Spec 0309's changes applied. Line numbers move; identifiers do not.

1. The Baseline writes managed guides, root blocks and the Setup Manifest
   only. No module field or template kind carries a repository tooling file.
2. `internal/baseline/catalog_load.go` allows the module fields `schemaVersion`
   … `semanticOwners`; an unknown field is `catalog.module.field.unknown`.
3. `internal/baseline/catalog_validate.go` `validateTemplates` checks a
   template's path, file and tokens, but not its kind; a kind is checked only
   where a root block, guide or repository extension references the template.
   The token pattern is `{{ name }}`, so a template must not contain `{{`.
4. `internal/baseline/profile_alignment.go`: an undeclared selected
   Verification command adds the blocking divergence
   `verification.command.undeclared` with the next action "select a command
   declared by the repository or add its local declaration".
5. `roundfix baseline skills restore` (`internal/cli/baseline_skills_restore.go`)
   is the model of a separate confirmation-gated Baseline command: a non-empty
   preview exits 3 with its Plan Digest, `--confirm-plan` applies the exact
   preview, an empty preview exits 0; `--repo` and `--format <text|json>`.
6. `internal/cli/cli.go`: the root usage and `commandUsage("baseline")` list
   the Baseline subcommands; `internal/docscontract/publicdocs_test.go`
   `TestBaselineDocumentationContract` checks help snippets per subcommand;
   `internal/cli/baseline_documentation_contract_test.go`
   `TestBaselineExamplesParse` parses every `roundfix baseline` example in the
   README, the context-driven guide, the setup skill and the Roundfix skill
   through `parsePublishedBaselineExample`.
7. `docs/references/skill-coverage.json` has one `command baseline <name>`
   entry per subcommand, covered by `.agents/skills/roundfix/references/baseline.md`.
8. `internal/config` `Load` reads a Project Config `.roundfixrc.yml` with
   `defaults.verification`, `worktree.bootstrap`, `worktree.bootstrap_timeout`,
   `worktree.caches` and `verification.tools`.

## System Architecture

- `internal/baseline/catalog_load.go`, `catalog_validate.go`: the field and
  its validation.
- `internal/baseline/starter_files.go` (new): the resolver.
- `internal/baseline/starters.go` (new): the preview and the confirmed write.
- `internal/baseline/profile_alignment.go`: the next action.
- `internal/baseline/assets/modules/python.json`,
  `assets/templates/index.json`, `assets/templates/starters/python/*` (new).
- `internal/cli/baseline_starters.go` (new), `internal/cli/baseline_profile.go`
  (routing), `internal/cli/cli.go` (usage).
- Docs: `docs/user-guide/commands/baseline.md`, the Roundfix skill reference
  and its mirror, `docs/references/skill-coverage.json`,
  `docs/references/behavior-surfaces.json`, `CONTEXT.md`.

No new package, import edge or dependency.

## Implementation Design

### Interfaces

```go
// internal/baseline/starter_files.go
type StarterFile struct {
	ID              string
	Module          string
	Path            string // repository-relative, slash-separated
	Content         []byte
	ContentIdentity string // "sha256:<hex of Content>"
}
func (c *Catalog) StarterFiles(modules []string) ([]StarterFile, error)

// internal/baseline/starters.go
const StartersSchemaVersion = "roundfix/baseline-starters/v1"
type StartersRequest struct {
	Repository  string
	ProfileID   string
	ConfirmPlan string
}
type StarterEntry struct {
	Path            string `json:"path"`
	Action          string `json:"action"` // "create" or "keep"
	ContentIdentity string `json:"contentIdentity"`
}
type StartersPayload struct {
	SchemaVersion string         `json:"schemaVersion"`
	Profile       string         `json:"profile"`
	State         string         `json:"state"` // "preview", "current" or "applied"
	Starters      []StarterEntry `json:"starters"`
	PlanDigest    string         `json:"planDigest"`
}
var ErrStartersStale error // wraps the code starters.plan.stale
func Starters(ctx context.Context, request StartersRequest, catalog *Catalog) (StartersPayload, error)
```

### Data Models

1. A Starter File declaration in a module: `{"id": "starter.<module>.<name>",
   "path": "<repository-relative path>", "template": "<template id>"}` in the
   optional array `starterFiles`.
2. A Starter File template: `{"id": "template.starter.<module>.<name>",
   "version": 1, "kind": "starter-file", "path": "starters/<module>/<file>",
   "tokens": []}`; its bytes are written verbatim.
3. The Plan Digest: `sha256:` followed by the catalog's canonical SHA-256 of
   `{"profile": <id>, "starters": [<entries in order>]}`.

### API Contracts

1. API Contract: validation. A module without `starterFiles` behaves as
   before. With it, the value is an array of objects
   (`catalog.module.starterFiles.invalid`), each with a unique `id` across the
   catalog, a safe relative `path` (`catalog.starter.file.path.invalid`) used
   by no other declaration (`catalog.starter.file.path.duplicate`), and a
   `template` that exists (`catalog.starter.file.template.unknown`) with kind
   `starter-file` (`catalog.starter.file.template.kind`). `StarterFiles`
   returns the given modules' declarations in module order and, within a
   module, in declaration order, with the template bytes and their identity.
2. API Contract: the Python starters. The `python` module declares, in order,
   `starter.python.makefile` → `Makefile`
   (`template.starter.python.makefile`, `starters/python/Makefile`),
   `starter.python.project-config` → `.roundfixrc.yml`
   (`template.starter.python.project-config`, `starters/python/roundfixrc.yml`)
   and `starter.python.verify-workflow` → `.github/workflows/verify.yml`
   (`template.starter.python.verify-workflow`, `starters/python/verify.yml`).
   - The Makefile starts with `SHELL := /bin/bash`, declares `.PHONY: help
     sync fmt fmt-check lint test verify verify-incremental`, a `help` target
     that prints the targets, `sync: uv sync --group dev`, and the recipes
     `uv run --frozen ruff format src tests` (fmt), `uv run --frozen ruff
     format --check src tests` (fmt-check), `uv run --frozen ruff check src
     tests` (lint), `uv run --frozen pytest -q` (test), then the lines
     `verify: fmt-check lint test` and `verify-incremental: verify`.
   - The Project Config holds, with English comments, `defaults.verification:
     make verify`, `worktree.bootstrap: uv sync --group dev --frozen`,
     `worktree.bootstrap_timeout: 5m`, `worktree.caches: [.venv,
     .pytest_cache, .ruff_cache]` and `verification.tools: [uv, make, git]`,
     and nothing else.
   - The workflow is named `Verify`, runs on `push` to `main` and on
     `pull_request`, sets `permissions: contents: read`, and has one job
     `verify` named `make verify` on `ubuntu-latest` with the steps
     `actions/checkout@v7` (`persist-credentials: false`),
     `astral-sh/setup-uv@v10.3.0` (`python-version: "3.12"`), `run: uv sync
     --group dev --locked` and `run: make verify`. No template contains `{{`.
3. API Contract: the preview. `Starters` resolves the profile as the other
   Baseline commands do (built-in or repository-owned), takes its Starter Files
   in profile module order, marks each `create` when nothing exists at the path
   (`os.Lstat` fails) and `keep` otherwise, and computes the Plan Digest. With
   no `create` the state is `current`; otherwise `preview`. A preview writes
   nothing.
4. API Contract: the confirmed write. With `ConfirmPlan` set and at least one
   `create`: a value different from the current digest returns an error
   wrapping `ErrStartersStale` and writes nothing; the current digest creates
   each `create` path in order, its parent directories at mode 0755 and the
   file at 0644 with exclusive creation, and returns state `applied`. If a
   write fails, every file this call created is removed before the error
   (`starters.write.failed: <path>: <cause>`) returns; directories it created
   may remain. A path that appears after the preview changes the digest and is
   refused as stale. With `ConfirmPlan` set and no `create`, the result is
   `current` and nothing is written.
5. API Contract: the next action. When the resolved profile's modules declare
   at least one Starter File, the next action of
   `verification.command.undeclared` is exactly "select a command declared by
   the repository or add its local declaration; `roundfix baseline starters
   --profile <id>` previews this profile's Starter Files", with `<id>` the
   profile ID; otherwise it keeps its bytes.
6. API Contract: the command. `roundfix baseline starters --profile <id>
   [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]`;
   `--profile` is required, `--repo` defaults to `.`, `--format` to `text`.
   JSON writes one `StartersPayload` on stdout. Text writes the lines of the
   Surface Transcripts. Exit codes: 0 for `current` or `applied`; 3 for
   `preview` and for a stale confirmation, which prints nothing on stdout and
   one line on stderr, "roundfix baseline starters: starters.plan.stale: the
   confirmed digest no longer matches; preview again"; 2 for invalid
   arguments, an unknown profile or an unknown format, with the usage error on
   stderr; 1 for a write failure, with the error on stderr. `--help` prints a
   text that starts with `Usage:`, names the synopsis, says the command writes
   only paths that do not exist and never changes an existing file, that a
   non-empty preview exits 3 and that an empty one exits 0. The root usage and
   `commandUsage("baseline")` gain the synopsis line after `roundfix baseline
   skills reconcile`.

### Surface Transcripts

Each runs in a temporary Git repository `<repo>` with one commit and no
Makefile, Project Config or workflow; `<digest>` is `sha256:` and 64
lowercase hex digits, the same value in both transcripts. A later preview
prints `- keep <path>` lines and `State: current` without a `Next:` line and
exits 0.

1. Surface Transcript: the preview writes nothing and asks for confirmation.

   ```transcript
   $ roundfix baseline starters --profile python --repo <repo>
   stdout:
   Baseline starters: python
   - create Makefile
   - create .roundfixrc.yml
   - create .github/workflows/verify.yml
   State: preview
   Plan Digest: <digest>
   Next: roundfix baseline starters --profile python --confirm-plan <digest>
   stderr:
   exit: 3
   ```

2. Surface Transcript: the confirmed preview creates the absent files.

   ```transcript
   $ roundfix baseline starters --profile python --repo <repo> --confirm-plan <digest>
   stdout:
   Baseline starters: python
   - create Makefile
   - create .roundfixrc.yml
   - create .github/workflows/verify.yml
   State: applied
   Plan Digest: <digest>
   stderr:
   exit: 0
   ```

### Invariants

```text
1. No existing file is written, truncated, compared or removed.
2. Nothing is written without the current Plan Digest.
3. A failed confirmed write leaves none of the files it created.
4. A Starter File is never part of the Change Plan, the Setup Manifest or a managed refresh.
5. Profiles without Starter Files and their alignment output keep their bytes.
6. No test runs uv, make on a starter, pytest, ruff or a workflow.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 2, Core Features 1 and 2, Success Metrics 1 and 2, User Story 1 → API Contracts 1 and 2
- Goals 1 and 3, Core Feature 3, Success Metric 3, User Stories 2 and 3, Invariants 1-3 → API Contracts 3 and 4
- Goal 4, Core Feature 4, Success Metric 4, User Story 4, Invariant 5 → API Contract 5
- Core Feature 5, Success Metric 5, Invariants 4 and 6 → API Contract 6, Surface Transcripts

## Absence Evidence

1. Starter File declarations — command: `grep -rn "starterFiles\|starter-file" internal/baseline`; result: no match at `dde9d5f5`.
2. The command — command: `grep -rn "baseline starters" internal docs .agents`; result: no match at `dde9d5f5`.
3. The glossary term — command: `grep -n "Starter File" CONTEXT.md`; result: no match at `dde9d5f5`.

## ADR Alternatives

- ADR-0323: none of its rejected alternatives is taken; the files are not written by `baseline apply`, not managed with markers, and not left as documentation-only examples.
- ADR-0068: its rejected alternative, separate human audit and apply journeys exposing engine phases, is not taken; the command previews and applies in one confirmation-gated workflow with an exact digest.
- ADR-0075: its rejected alternatives are not taken; profile divergence still resolves through a confirmed repository-owned adaptation, and the new next action only adds a pointer to the Starter Files for a profile that has them.
- ADR-0219: no rejected alternative is taken; a repository-owned draft adapted from the `python` profile binds to it and, keeping the `python` module, keeps its Starter Files.
- ADR-0241: no rejected alternative is taken; Starter Files come from the embedded module, as a repository profile's skill contracts come from the embedded snapshots.
- ADR-0322: none of its rejected alternatives is taken; the `python` module and its common setup stay as Spec 0309 delivers them, and this Spec only adds Starter Files to the module.

## Integration Points

The embedded Baseline catalog and its derived files; the CLI; the repository's
Managed Refresh; the owned skill mirrors and version record; the Project Config
loader, used by a test to read a written starter. No test, Verification
command or QA row runs uv, pytest, ruff, a starter Makefile target or a
workflow, reaches the network, reads the real `~/.roundfix`, reads the
Secondbrain or reads another repository.

## Testing Approach

Every new top-level test and subtest calls `t.Parallel()` first, passes no
real-time duration and uses temporary directories; the new files use no
platform API and compile for Windows.

- `internal/baseline/starter_files_test.go` (new, Build Order 1):
  - `TestThePythonModuleDeclaresItsStarterFiles`: for the `python` profile's ordered
    modules the three paths in order with `sha256:` identities, each template
    holding API Contract 2's lines (`verify: fmt-check lint test`,
    `verify-incremental: verify`, the four `uv run --frozen` recipes; the five
    Project Config values; `pull_request:`, `run: uv sync --group dev
    --locked`, `run: make verify`, `python-version: "3.12"`), and no Starter
    File for the other four profiles.
  - `TestAStarterFileDeclarationIsValidated`: on `cloneEmbeddedAssets`, the
    `python` module with an unsafe path, a duplicate path, an unknown template
    and a guide template reports the four codes of API Contract 1.
- `internal/baseline/starters_test.go` (new, Build Order 2), each over a
  repository from `newPlanRepository` (whose Makefile exists):
  - `TestStartersPreviewsAndCreatesOnlyAbsentFiles`: the preview is `preview`
    with `keep Makefile` and two `create` entries and writes nothing; the
    confirmed apply is `applied`, creates the two files with the template
    bytes, leaves the Makefile byte-identical; a second preview is `current`.
  - `TestStartersRefuseADigestThatNoLongerMatches`: after a preview, a
    `.roundfixrc.yml` written by the test makes the confirmation fail with
    `ErrStartersStale`, and no `.github` exists afterwards.
  - `TestStartersRemoveWhatTheyCreatedWhenAWriteFails`: a regular file named
    `.github` makes the workflow write fail; the call returns an error and
    `.roundfixrc.yml`, created earlier in the call, no longer exists.
  - `TestTheUndeclaredGateNamesTheStartersCommand`: through
    `ResolveProfileAlignment` in a repository without a Makefile, the
    `verification.gate` divergence's next action is API Contract 5's text with
    `python`, and for `rust-cli` it is the unchanged text.
- `internal/cli/baseline_starters_test.go` (new, Build Order 3):
  - `TestBaselineStartersPreviewTranscript` and
    `TestBaselineStartersApplyTranscript`: each Surface Transcript asserted
    whole, every stdout line in order, empty stderr and the exit code, with
    the digest replaced by its controlled variation after checking its shape.
  - `TestBaselineStartersWritesAProjectConfigRoundfixLoads`: after the apply,
    `config.Load` with a temporary home and the repository as working
    directory returns `make verify`, the bootstrap command, five minutes, the
    three caches and the tools `uv`, `make`, `git`.
  - `TestBaselineStartersJSONAndRefusals`: `--format json` returns the
    payload with schema `roundfix/baseline-starters/v1`; a stale
    `--confirm-plan` exits 3 with the stale line on stderr and writes nothing;
    a missing `--profile`, an unknown profile and `--format yaml` exit 2.
  - `TestBaselineStartersHelp`: `--help` exits 0 with the synopsis and the
    three sentences of API Contract 6.
- Declared changes: `TestBaselineDocumentationContract`
  (`internal/docscontract/publicdocs_test.go`) gains the synopsis snippet in
  the `baseline` case and a `baseline starters` help case;
  `parsePublishedBaselineExample`
  (`internal/cli/baseline_documentation_contract_test.go`) parses a
  `starters` example with `parseBaselineStartersCommand`. No other existing
  test changes.
- Prototype record (2026-10-23, scratch worktree with Spec 0309's prototype):
  the Build Order 1 tests and the first three Build Order 2 tests (then with
  rollback and staleness in one test) passed after the change and could not
  exist before it, and the `internal/baseline` suite passed with the starters
  declared. These sabotages each failed the named test: the template kind
  check removed → `TestAStarterFileDeclarationIsValidated/guide_template`;
  `verify-incremental: verify` removed from the Makefile template →
  `TestThePythonModuleDeclaresItsStarterFiles`; the digest comparison
  skipped → the staleness test ("stale apply error = <nil>"); the removal on
  failure skipped → the rollback test ("failed apply left a created file").

## Build Order

1. Declarations: the module field, its validation, `StarterFile` and
   `StarterFiles`, the three templates and index entries, the `python`
   module's declarations, the module record command, `make baseline-digests`,
   the Managed Refresh twice, and `starter_files_test.go` (depends on Spec
   0309's `python` module, through `requires`).
2. The operation: `Starters`, its payload and errors, the next action in
   profile alignment, and `starters_test.go` (depends on 1).
3. The command and its documentation: `baseline_starters.go`, routing and
   usage, `baseline_starters_test.go`, the two declared test changes, the
   command guide, the Roundfix reference, mirror, version and record, the
   Skill Coverage Map entry, **Starter File** in the glossary, and the
   Behavior Surface Record re-recorded last (depends on 2).
4. QA gate (depends on: 1, 2, 3).

`requires: [0309-a-python-baseline-profile-for-uv-pytest-and-ruff]`: Build
Order 1 extends the module Spec 0309 creates.

## Risks & Considerations

- The action versions in the workflow (`actions/checkout@v7`,
  `astral-sh/setup-uv@v10.3.0`) age; the repository owns the file once
  written, and a later Spec can raise the template's version.
- `uv run --frozen` in the gate does not notice a stale lockfile; CI's
  `uv sync --locked` does, and the Python guide (Spec 0309) states both.
- A repository with its own Makefile that lacks `verify` keeps it; the preview
  says `keep`, and alignment still reports the undeclared gate.

## Glossary

- adds: **Starter File**

## Decisions

- A separate command writes Starter Files only where absent; see ADR-0323.
- The Makefile and bootstrap use `--frozen`; CI uses `--locked`.
- The command's exit codes follow `baseline skills restore`.
- The profile is `python` and fixes no application surface (maintainer
  decisions of 2026-10-24, quoted in `_prd.md` → Decisions); the Starter
  Files carry nothing specific to a command-line tool, so only the profile
  identifier in the command, the transcripts and the next action follows the
  rename.
