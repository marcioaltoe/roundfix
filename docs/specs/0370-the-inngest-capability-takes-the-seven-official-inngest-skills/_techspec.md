---
spec: 0370-the-inngest-capability-takes-the-seven-official-inngest-skills
prd: _prd.md
created: 2026-11-09
---

# The Inngest capability takes the seven official Inngest skills — Technical Spec

## Executive Summary

The `backend` module's `requiredSkills` and `skillDispatch` replace
`inngest` with the seven `inngest-*` skills, one trigger each, and the
module record command raises its version. `inngest` joins `retiredSkills` in
`internal/baseline/retired_skills.go`, so the existing ADR-0350 machinery
plans and applies the deletion of an adopter's unedited copy without new
code. The Catalog Skill Bump then moves the five upstream Setup Snapshots,
this repository's vendored catalog skills and their lock entries to catalog
commit `e45466d2`; the asset-sync parity row follows through the sanctioned
`jq` transform, and `make baseline-digests` and a managed refresh built after
it regenerate every derived file. Two new tests pin the result, and the ADR,
glossary and guides say why `inngest` was followed to its replacements
(ADR-0380). The trade-off this design accepts is that the bump moves every
catalog skill that changed since `c8b2e85`, not only the Inngest ones: the
snapshots and the lock name one Catalog Revision by rule (ADR-0360), so a
narrower move is not available.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the trigger ids follow the existing `trigger.<module>.<skill>`
  form. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files only; no Task, test
  or QA row reaches the network. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0380 (new) records the rule,
  ADR-0380: "A skill the catalog replaces is followed to its replacements,
  and the replaced skill becomes a Retired Skill." It amends ADR-0206,
  ADR-0206: "It is never followed to a skill with a similar purpose under
  another name"; task_02 writes the notice. ADR-0350 deletes the unedited
  copy, ADR-0350: "It is `delete` when every copy it finds is provably the
  one the repository installed". ADR-0360 makes the vendored catalog
  skills and the snapshots follow the catalog through the bump, ADR-0360: "The Catalog Skill Bump is the only writer." ADR-0191, ADR-0191: "A setup snapshot follows its upstream
  by name". ADR-0364, ADR-0364: "A required external skill is checked only
  for presence". ADR-0204, ADR-0204: "A composed profile takes a setup
  composed from upstream setups by name". ADR-0322, ADR-0322: "A Python
  profile imports the published upstream Python preset". ADR-0345, ADR-0345:
  "A derived-record writer names its command to the suite guard". Kept and
  not changed: ADR-0219, ADR-0221, ADR-0224, ADR-0241 and ADR-0264, each
  cited because it cites a listed ADR. Source:
  `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BA batch 2 answers of 2026-11-09 quoted in the authorization record;
  no Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0370-the-inngest-capability-takes-the-seven-official-inngest-skills/_authorization.md`;
  bounded files: `.agents/skills/context7-cli/references/setup.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `docs/agents/setup-context.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/modules/backend.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/setups/go-cli-typescript-bun.json`,
  `internal/baseline/assets/setups/go.json`,
  `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/setups/rust.json`,
  `internal/baseline/assets/setups/typescript.json`,
  `skills/roundfix/SKILL.md`.

## System Architecture

- `internal/baseline/assets/modules/backend.json` (task_01): `requiredSkills`
  lists `inngest-cli`, `inngest-durable-functions`, `inngest-events`,
  `inngest-flow-control`, `inngest-middleware`, `inngest-setup` and
  `inngest-steps` in lexical order where it listed `inngest`; `skillDispatch`
  holds one entry per skill, in the same place and order, each with the one
  trigger API Contract 1 names. Nothing else in the module changes except
  the version line the record command writes.
- `internal/baseline/retired_skills.go` (task_01): `"inngest"` joins
  `retiredSkills`. `InstalledRetiredSkills`, `DeleteRetiredSkills`, the asset
  sync's drop of a Retired Skill and the `baseline update` report already
  read that set, so no other production code changes.
- The Setup Snapshots, `skills-lock.json`, the vendored `context7-cli` file,
  the parity row and the derived files (task_01) change only through the
  commands of API Contract 3.
- Documentation (task_02): the ADR-0206 notice, the **Retired Skill**
  sentence, the `baseline` command guide and the Roundfix skill's `baseline`
  reference.

## Implementation Design

### Interfaces

No exported or unexported Go signature changes. The one production change is
an entry in an existing map, and it is what removes an adopter's old copy:
with `inngest` in `retiredSkills`, `InstalledRetiredSkills` inspects
`.agents/skills/inngest`, `.claude/skills/inngest` and the `inngest` lock
entry and plans `delete` when the copy hashes to the lock's `computedHash`,
or `keep` with the reason `edited` or `unrecorded`; `DeleteRetiredSkills`
removes the planned copy on an approved `roundfix baseline update`; and the
asset sync drops `inngest` from every snapshot it writes.

```go
var retiredSkills = map[string]struct{}{
	// existing entries unchanged
	"inngest": {},
}
```

### Data Models

No schema changes. The module keeps `setup-context-driven/module-v2`, the
snapshots `setup-context-driven/setup-snapshot/0.0.1`, and the lock its
version 1 shape.

### API Contracts

1. API Contract: the `backend` module's seven dispatch triggers, each the
   only trigger of its skill. Each line below is the trigger id, a colon and
   a space, then its `when` text, exactly:

   ```text
   trigger.backend.inngest-cli: Installing or running the Inngest CLI and its local Dev Server.
   trigger.backend.inngest-durable-functions: Building or changing Inngest functions, their triggers, retries, and schedules.
   trigger.backend.inngest-events: Designing Inngest events, fan-out, and idempotent event handling.
   trigger.backend.inngest-flow-control: Configuring Inngest concurrency, throttling, rate limits, debounce, or batching.
   trigger.backend.inngest-middleware: Adding Inngest middleware for logging, tracing, error tracking, or dependency injection.
   trigger.backend.inngest-setup: Setting up Inngest in a project: the SDK, client configuration, environment variables, and the serve endpoint.
   trigger.backend.inngest-steps: Building or changing the steps of Inngest functions: step runs, sleeps, waits for events, and invocations.
   ```
2. API Contract: `roundfix baseline update` reports `inngest` with the
   existing Retired Skill lines of ADR-0350, unchanged in shape:
   `- retired inngest: no longer required by the Baseline; the approved update deletes .agents/skills/inngest, .claude/skills/inngest and its skills-lock.json entry`
   for a copy proven unedited, and
   `- retired inngest: kept, its copy differs from its skills-lock.json computedHash; delete .agents/skills/inngest, .claude/skills/inngest and its skills-lock.json entry by hand if the repository no longer uses it`
   for an edited one.
3. API Contract: the regeneration sequence, in this order, every Roundfix run
   with a fresh temporary `HOME` and the module and build caches kept at
   their `go env` values:
   - step 1: `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`;
   - step 2: `GIT_ALLOW_PROTOCOL=file make skills-third-party-bump SKILLS_CATALOG=/Users/marcio/.roundfix-operator/scratch/spec-0370/skills SKILLS_REF=e45466d2b9e15480eb286a4c19b513494a867bce`,
      exit 0, ending with `skills-third-party-bump: 33 catalog skills and every upstream Setup Snapshot pinned to e45466d2b9e15480eb286a4c19b513494a867bce`;
   - step 3: the asset-sync parity transform of Spec 0208 with Spec 0317's
      composed-snapshot filter: write the sorted snapshots with
      `jq -s 'sort_by(.id)[]' internal/baseline/assets/setups/*.json` to a
      temporary file, run
      `jq --indent 2 --slurpfile s <that file> -f <program> internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
      into a second temporary file, and copy it over the fixture, with this
      program:

      ```jq
      .input.sourceRevision as $rev
      | def row($s): {activationBundles: $s.activationBundles, digest: "", id: $s.id,
          schemaVersion: $s.schemaVersion, version: $s.version,
          source: {path: $s.source.path, ref: $rev, repository: "example/skills", type: "github"},
          skills: [$s.skills[] | if .source.type == "github"
            then {name, path, source: {path: .source.path, ref: $rev, repository: "example/skills", type: "github"}, treeDigest: ""}
            else {minimumVersion, name, path, source: {name: .source.name, type: "repo"}} end]}
          | if .activationBundles == null then del(.activationBundles) else . end;
        .manifest.setups = [$s[] | select(.source.type != "composed") | row(.)]
      | .managedEntryLedger = [.manifest.setups[] | {digest: "", id}]
      | .normalizedOutput.result.findings = [.manifest.setups[] | {
          action: "Review the snapshot diff and run asset validation.",
          code: "skills.setup-snapshot.updated", managedId: ("setup." + .id),
          message: "Setup snapshot was synchronized from the canonical source.",
          path: ("assets/setups/" + .id + ".json"), severity: "info"}]
      | .normalizedOutput.result.summary.info = (.manifest.setups | length)
      ```
   - step 4: `make baseline-digests`, which fills the row's digests and rewrites
      the corpus manifest;
   - step 5: `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
      twice, the second reporting `File changes: 0`. The bump refreshes the
      Setup Manifest with a binary it built before `make baseline-digests`
      rewrote the profile's golden digest, so without this step the recorded
      catalog digest is stale and `TestTheSetupManifestIsCurrent` fails.

### Documentation

task_02 writes the ADR-0206 amendment notice, appends one sentence to the
**Retired Skill** definition in `CONTEXT.md`, adds one sentence and the
`inngest` paths to the Retired skills section of
`docs/user-guide/commands/baseline.md`, adds the same sentence to
`.agents/skills/roundfix/references/baseline.md`, mirrors it with
`make skills-sync`, raises the Roundfix skill's version by the record command
and re-records the Behavior Surface Record after its last guide edit.

### Surface Transcripts

None. `roundfix baseline update` keeps its report layout; the `inngest`
lines are existing Retired Skill lines, asserted by the named CLI test.

## Coverage Map

- Goal 1, Story 1, Core Feature 1, Success Metric 1 → System Architecture
  (the `backend` module's required skills and dispatch); API Contract 1;
  Testing Approach (`TestTheBackendModuleTakesTheSevenOfficialInngestSkills`,
  `TestTheBackendModuleRequiresTheSevenInngestSkills`); Build Order 1.
- Goal 2, Story 3, Core Feature 3, Success Metric 2 → System Architecture
  (the snapshots, lock and derived files that change only through commands);
  API Contract 3 (the bump at `e45466d2`, the parity transform, the digests
  and the managed refresh); Build Order 1.
- Goal 3, Story 2, Core Feature 2, Success Metric 3 → System Architecture
  (`inngest` in `retiredSkills`, read by `InstalledRetiredSkills` and
  `DeleteRetiredSkills`); Interfaces (the map entry); API Contract 2 (the
  deletion and kept lines); Testing Approach
  (`TestBaselineUpdateApprovedApplyDeletesTheReplacedInngestSkill`); Build
  Order 1.
- Goal 4, Core Feature 4 → Documentation (the ADR-0206 notice, the glossary
  sentence, the guide and the skill reference); Build Order 2.
- Success Metric 4 → Testing Approach (the new tests fail on the unchanged
  tree); Build Order 3.

Deviations from the PRD: None.

## Absence Evidence

1. A Baseline module that requires an `inngest-*` skill — command: `grep -rn 'inngest-' internal/baseline/assets/modules internal/baseline/assets/setups`; result: no match (exit 1).
2. `inngest` as a Retired Skill — command: `grep -n 'inngest' internal/baseline/retired_skills.go`; result: no match (exit 1).
3. A test of the Inngest skills or of an `inngest` migration — command: `grep -rln 'inngest' internal --include='*_test.go'`; result: no match (exit 1).

## ADR Alternatives

- ADR-0380: the decision this Spec records; it rejects dropping `inngest`
  without a replacement, a partial set, a hand-deleted removed skill and a
  conditional requirement, and the design takes none of them.
- ADR-0206: amended for a skill the catalog replaces; a skill removed without
  a named replacement is still dropped, and the design takes neither of the
  alternatives it rejected.
- ADR-0350: applied unchanged; it rejected deleting an edited copy, and the
  design keeps one.
- ADR-0360: applied; it rejected bumping only the lock and re-pinning test
  constants, and the new tests name no Catalog Revision.
- ADR-0191: applied; the asset sync stays the snapshots' only writer.
- ADR-0364: kept; the seven are checked for presence like every required
  external skill, and the design adds no snapshot comparison it rejected.
- ADR-0204: kept; the composed snapshot is recomposed by the sync, not by
  hand.
- ADR-0322: kept; the Python preset moves to the same Catalog Revision and
  its list is unchanged.
- ADR-0345: applied; every regeneration command is declared with its
  outputs.
- ADR-0219: kept; profile drafts are unchanged, and the design takes none of
  its rejected options.
- ADR-0221: kept; a required skill is still held to its embedded snapshot,
  which now pins the seven.
- ADR-0224: kept; this repository's vendored skills and lock stay in step
  with the `go` snapshot through the bump.
- ADR-0241: kept; a repository profile still takes its contracts from the
  embedded snapshots.
- ADR-0264: kept; the bump's restore replaces only a vendored skill it proves
  older, as it rejected otherwise.

## Integration Points

None outside the repository at delivery. After the release, an adopter's
`roundfix baseline update` reads the new module and snapshots and the
skills stage installs the seven from GitHub, the Unreachable Acceptance
entry of the PRD.

## Testing Approach

- `internal/baseline/inngest_skills_test.go` (new, task_01):
  `TestTheBackendModuleTakesTheSevenOfficialInngestSkills` reads the embedded
  catalog: the `backend` module requires each of the seven and dispatches
  each once with the trigger of API Contract 1, and names `inngest` in
  neither list; `loadRestoreProfile` for `standard-typescript-monorepo` and
  `go-cli-typescript-monorepo` requires and pins each of the seven at
  `skills/05-implementation-loop/<skill>` with a tree digest and neither
  requires nor pins `inngest`; `RetiredSkills()` holds `inngest`. It names
  no Catalog Revision.
- `internal/cli/baseline_update_inngest_migration_test.go` (new, task_01):
  `TestTheBackendModuleRequiresTheSevenInngestSkills` resolves the required
  external skills of a Setup Manifest naming the `typescript` and `backend`
  modules; `TestBaselineUpdateApprovedApplyDeletesTheReplacedInngestSkill`
  plants `.agents/skills/inngest`, its `.claude/skills` link and a lock entry
  with its computed hash in an adopted repository and asserts the preview
  line, the planned deletion, the approved apply and an unchanged rest of
  the lock; its `edited` subtest asserts the kept line and unchanged bytes.
- `internal/baseline/retired_skills_test.go`:
  `TestNoCatalogEntryNamesARetiredSkill` pins the Retired Skill list, so its
  expected list gains `inngest`.
- A prototype on `cfc08683` passed these tests,
  `go test ./internal/baseline ./skills`,
  `-tags docscontract ./internal/docscontract`,
  `-run 'Baseline|Skill|Doctor' ./internal/cli` and `make repo-test`; on the
  unchanged tree the new tests fail, and each sabotage the QA Task lists
  failed its named test.

## Build Order

1. The module, the Retired Skill entry, the regeneration sequence of API
   Contract 3 and the three test changes.
2. The ADR-0206 notice, the glossary sentence, the guide, the Roundfix
   skill reference, its version record and the Behavior Surface Record
   (depends on: 1, so the documents describe delivered behavior).
3. Final QA gate (depends on: 2).

No step needs another Spec's code. Specs 0366 to 0369 merge first and may
move `CONTEXT.md`, `docs/references/behavior-surfaces.json`,
`skills/testdata/owned-skill-versions.json` and the Roundfix skill's version
fields; task_02 regenerates the two records by their commands on its own
base.

## Risks & Considerations

- The bump moves every pinned catalog skill changed since `c8b2e85`:
  `context7-cli` in all five snapshots, `ui-ux-pro-max`, `vite` and
  `onioncry` in the two TypeScript snapshots, and `fastapi` and `pdf` in the
  `python` snapshot,
  so adopters see those as Setup Snapshot Differences after the release, the
  routine effect ADR-0360 states.
- A TypeScript adopter that edited `inngest` keeps it beside the seven until
  it deletes it by hand; the report names the paths.
- An adopter without network keeps `inngest` deleted and the seven missing
  until its next update with network; Doctor names them (ADR-0364).

## Glossary

- changes: **Retired Skill** — a skill the catalog replaces with skills it names becomes a Retired Skill, and the module that required it requires the replacements.

## Decisions

- `inngest` becomes a Retired Skill and the module requires the seven. See
  ADR-0380.
- The snapshots move only through the Catalog Skill Bump at `e45466d2`, from
  the operator's clone, followed by the parity transform, the digests and a
  managed refresh built after them.
- Two implementation Tasks and the QA gate, in series.
