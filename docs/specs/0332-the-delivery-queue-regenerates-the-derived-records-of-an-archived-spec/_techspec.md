---
spec: 0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec
prd: _prd.md
created: 2026-10-30
---

# The Delivery Queue regenerates the derived records of an archived Spec — Technical Spec

## Executive Summary

The queue's derived merge is not at fault, and neither is the suite guard's
reader. The guard already honors an Archive Record's `regeneration` list. The
coverage record writer is at fault: `TestCoverageEquivalence`
(`internal/spec/coverage_test.go`) rewrites the record under
`-update-coverage-record` and never calls
`suiteguard.DeclareSanctionedRegeneration`. The guard exempts a write only
under a declared command (ADR-0149), so it refuses every coverage re-record,
whatever record sanctions it. The fix adds that declaration and then proves it
twice:

- The regeneration contract runs every Project Config regeneration in a
  repository copy. That copy is sanctioned only by an archived item's record,
  and the contract no longer injects a declaration of its own.
- A queue-level fixture drives `ResolveConflict` over an archived item whose
  Pull Request conflicts only in the coverage record.

The trade-off is that both proofs run real `go test` commands in repository
copies, at about 10–60 s each. A fake would not exercise the guard, and the
guard is where the defect lived.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path or credential;
  all work is offline in temporary repositories. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0149, ADR-0192, ADR-0081 and
  ADR-0345 govern the design. ADR-0233, ADR-0250, ADR-0256, ADR-0197,
  ADR-0272, ADR-0279 and ADR-0286 are
  listed because they cite one of these, and they do not apply.
  `_prd.md` → Project Constraints states why each applies or does not, and
  ADR Alternatives below relates each one. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AS decisions of 2026-10-30 quoted in the authorization record; no
  Makefile, workflow, module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec/_authorization.md`;
  bounded files:
  `internal/spec/coverage_test.go`,
  `docs/agents/specific-repository.md`,
  `docs/references/coverage-record.json`.

## Current behavior

Read on 2026-10-30 at `35de7f74`; `_prd.md` → Validation records the runs.

1. `commandDeliveryWorkflow.ResolveConflict` (`internal/cli/deliver_workflow.go`)
   merges the fetched default branch with `--no-commit`. When every conflicted
   path matches a Derived Path Declaration, it takes the default side, runs
   each matched `regenerate` through `daemon.ExecVerifier` with the log
   `delivery/<slug>/derived-regeneration-<n>.log`, and on a non-zero exit
   returns `regenerate derived paths: <err>`. Its deferred cleanup restores
   the worktree and aborts the merge.
2. `suiteguard.run` (`internal/suiteguard/suiteguard.go`) fingerprints the
   repository before and after a guarded package's tests.
   `removeSanctionedRegeneration` drops only the changed paths that some
   declaration lists for the command the process declared, and keeps every
   change when no command was declared.
3. `suiteguardcontract.readSanctionedRegenerations`
   (`internal/suiteguardcontract/regeneration.go`) reads four sources:
   granted `_authorization.md` records under `docs/specs/`, legacy records,
   archived legacy folders, and, through `readArchiveRecordRegenerations`,
   the `regeneration` list of every `docs/history/specs/*.md` Archive Record.
4. `.roundfixrc.yml` declares seven regenerations. Their writers fall into
   three groups:

   | declaration | writer | guarded | names its command |
   | --- | --- | --- | --- |
   | module record, profile digest, formatter goldens | `internal/baseline` (`TestEveryBaselineModuleVersionIsRecorded`, then `make baseline-digests`) | yes | yes: `recordModuleVersionsCommand`, `baselineDigestRegenerationCommand` |
   | `docs/agents/`, `setup-context.json` | `go run ./cmd/roundfix baseline update` | no (not a test) | not needed |
   | owned-skill version record | `./skills` | no | not needed |
   | coverage record | `internal/spec` `TestCoverageEquivalence` | yes | **no** |
   | Behavior Surface record | `./internal/docscontract` | no | not needed |

   The skill coverage map (`docs/references/skill-coverage.json`) is
   hand-maintained and no regeneration writes it.
5. `TestRegenerationIsDeclared` (`internal/config/regeneration_declared_test.go`)
   runs every regeneration in a copy. Its fixture writes a legacy
   authorization and an `internal/spec/regeneration_fixture_test.go` whose
   `init` declares the coverage command, so the contract passes over the
   defect. The copy's coverage record is never stale, so the guarded write is
   never exercised.

## System Architecture

No production file changes. Each component is listed with its change:

- `internal/spec/coverage_test.go`: the coverage writer names its command
  (task_01).
- `internal/config/regeneration_declared_test.go`: the fixture models an
  archived item and makes the coverage record stale. It loses the injected
  declaration and gains an unsanctioned case (task_01).
- `internal/cli/deliver_derived_archived_item_test.go` (new): the queue-level
  archived-item fixture (task_02).
- `docs/agents/specific-repository.md` and `CONTEXT.md`: the rule and the
  revised **Archive Record** (task_01).
- `docs/references/skill-coverage.json`: the `config keys` entry gets a
  dated Coverage Review, because a test file in `internal/config` changes and
  no configuration key does. The Behavior Surface fingerprint excludes test
  files, so `docs/references/behavior-surfaces.json` does not move (task_01).
- `docs/references/coverage-record.json`: re-recorded with its command after
  task_02 adds two top-level tests.

## Implementation Design

### Interfaces

```go
// internal/spec/coverage_test.go
const coverageRecordCommand = "go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1"

// in TestCoverageEquivalence, before writeCoverageRecord:
if *updateCoverageRecord {
	suiteguard.DeclareSanctionedRegeneration(coverageRecordCommand)
	// ...existing write...
}

// internal/config/regeneration_declared_test.go
func writeArchivedItemRecord(t *testing.T, repository string, sanctionCoverage bool)
func staleCoverageRecord(t *testing.T, repository string) (restoredTest string)
```

### Data Models

The fixture Archive Record is `docs/history/specs/0001-the-archived-item.md`.
Its front matter carries `schema: roundfix/archive-record/v1`, the identity
fields of a passing record, and a `regeneration` list. That list holds:

- the module record command, with `outputs` for
  `internal/baseline/assets/modules/context-workflow.json` and
  `internal/baseline/module-versions.json` (config fixture only);
- `make baseline-digests`, with no outputs, so that the ownership records
  resolve them (config fixture only);
- the coverage command, with the output `docs/references/coverage-record.json`,
  unless the case under test removes it.

### API Contracts

1. API Contract: the coverage writer's declaration. When
   `-update-coverage-record` is set, `TestCoverageEquivalence` calls
   `suiteguard.DeclareSanctionedRegeneration(coverageRecordCommand)` before it
   writes. Without the flag it declares nothing and writes nothing. The
   constant equals the `.roundfixrc.yml` regenerate string for
   `docs/references/coverage-record.json` byte for byte. The config contract
   proves the equality by running that string under a record that lists only
   it.
2. API Contract: the regeneration contract's fixture. In the repository copy
   `TestRegenerationIsDeclared` removes three sources of sanction:
   - every `docs/history/specs/*.md` record;
   - the `docs/specs/` tree;
   - `docs/workflow/authorizations/`.

   It writes the archived item's record (Data Models, coverage included). It
   edits the module as today and replaces one entry of the copied coverage
   record with `TestRegenerationFixtureStaleEntry`. It commits, then runs
   every declaration as today. Every command must exit 0. The changed paths
   must include `docs/references/coverage-record.json`, and the record must no
   longer contain the stale name. No file is written under `internal/spec/`.
   The ablation subtests keep their bytes.
3. API Contract: the unsanctioned case. The subtest
   `an_unsanctioned_coverage_write_is_refused` builds a second copy whose
   record omits the coverage command. It runs only the coverage declaration
   and requires a non-nil error and the output line
   `suiteguard: modified: docs/references/coverage-record.json`.
4. API Contract: the queue-level fixture. A helper copies the repository into
   a temporary `item` repository from inside the test process. It reads the
   root found by walking up to `.roundfixrc.yml` and skips `.git`, `.gocache`,
   `bin`, `docs/` and `.agents/`. It then copies
   `docs/references/coverage-record.json` alone. The fixture's history is:
   - base: a commit that adds `docs/specs/0001-the-archived-item/_prd.md`,
     pushed to a bare `origin`;
   - item branch: a commit that replaces one record entry with
     `TestItemSideEntry`, then an archive commit that removes the Spec folder
     and adds the record (Data Models, coverage command only, or none);
   - default branch: a commit that replaces the same entry with
     `TestDefaultSideEntry`, pushed.

   The workflow uses `preflight.ExecGitRunner` and `roundconfig.Builtin()`,
   with the artifact directory and home under the test's temporary directory.
   It keeps the inherited `GOCACHE` and `GOMODCACHE`.
5. API Contract: resolution. With the sanction, `ResolveConflict` returns:
   - a new `Head`, no `SourcePaths`, and `Regenerated` equal to the
     `.roundfixrc.yml` coverage command;
   - a committed record that contains the original entry and neither fake
     entry, with no conflict markers;
   - a merge commit whose parents are the item head and `origin/main`, whose
     message carries `Roundfix-Delivery: derived-merge`;
   - a clean worktree.
6. API Contract: refusal. Without the sanction, `ResolveConflict` returns an
   error containing `regenerate derived paths`, HEAD stays at the item head,
   the worktree is clean, and no merge is in progress. One
   `derived-regeneration-*.log` under `delivery/<slug>/` contains
   `suiteguard: modified: docs/references/coverage-record.json`.

### Surface Transcripts

None. No command, flag, message or exit code changes; the queue's existing
park reason and log name keep their bytes.

### Invariants

```text
1. A run of a guarded package without a record flag declares no command, so any repository write it makes is refused.
2. A declared command's write is exempt only for outputs a record in the tree lists for that command; every other changed path is refused.
3. No Project Config declaration and no queue step sanctions a write or disables the suite guard.
4. The queue's merge, park reasons, commit subject and log names keep their bytes.
5. No fixture injects a sanctioned-regeneration declaration into a copied package.
```

## Coverage Map

- Goal 1, Core Feature 1, Success Metric 1 → `coverageRecordCommand` and its
  declaration in `TestCoverageEquivalence` (API Contract 1), proved by
  `TestRegenerationIsDeclared` (API Contract 2).
- Goal 2, Success Metric 2 → `an_unsanctioned_coverage_write_is_refused`
  (API Contract 3), Invariants 1–2.
- Goal 1, Core Feature 3, Success Metric 3 →
  `internal/cli/deliver_derived_archived_item_test.go` (API Contracts 4–6).
- Goal 3, Core Feature 2, Success Metric 4 → the Current behavior table and
  `TestRegenerationIsDeclared` running all seven declarations (API Contract 2).
- Core Feature 4 → `docs/agents/specific-repository.md`, **Archive Record**
  in `CONTEXT.md` (task_01).
- `_prd.md` → Skills, `config keys` unchanged → the Coverage Review in
  `docs/references/skill-coverage.json` (task_01).
- Success Metric 5 → declared unreachable in `_prd.md`.

Deviations from the PRD: None.

## Absence Evidence

1. The coverage writer's declaration — command: `rg -n "DeclareSanctionedRegeneration" internal/spec`; result: no match.
2. A regeneration contract case where no record sanctions the command — command: `rg -n "unsanctioned|staleCoverageRecord|archive-record" internal/config/regeneration_declared_test.go`; result: no match.
3. A queue-level fixture over an archived item that uses a real guarded command — command: `rg -ln "update-coverage-record" internal/cli --glob '*_test.go'`; result: no match; the conflict tests use shell `regenerate` commands only.
4. The repository rule naming the declaration — command: `rg -n "DeclareSanctionedRegeneration" docs/agents`; result: no match.

## ADR Alternatives

- ADR-0149: followed and refined; the guard's single declaration and its exemption rule are kept, and the private exemption list it rejected is not reintroduced; the writer now supplies the command it needs.
- ADR-0192: followed; a Derived Path Declaration still grants no authority, and the option of sanctioning its paths is rejected again in ADR-0345.
- ADR-0081: followed; authorization stays drawn around the causing edit, and no per-output grant is added.
- ADR-0345: implemented as recorded, including its rejection of restoring the authorization at `source_revision`, sanctioning queue declarations and skipping the guard.
- ADR-0233: untouched; line-scoped skill version resolution is not changed.
- ADR-0250: untouched; the coverage record's format and platforms are not changed.
- ADR-0256: untouched; the release skill wait is not changed, and the `config keys` Coverage Review follows its Skills Declaration.
- ADR-0197: untouched; review lineage rounds are not changed.
- ADR-0272: untouched; review records and corrections are not changed.
- ADR-0279: untouched; status and dismissals are not changed.
- ADR-0286: untouched; Correction Review is not changed.

## Integration Points

Local `git` and `go` only. No network: the fixtures use a bare local `origin`,
and the nested `go test` reads the inherited module cache.

## Testing Approach

All fixtures are temporary repositories under `t.TempDir()`, made by
in-process copies so the Go test cache records their inputs. A copy made by a
child process would hide a later change in the writer from a cached result.
No test reads or writes the real `~/.roundfix`, sets `HOME` for the nested
`go`, or reaches the network.

- task_01: `internal/config/regeneration_declared_test.go`
  `TestRegenerationIsDeclared`, whose new subtest is
  `an_unsanctioned_coverage_write_is_refused`. It is the class contract: a
  guarded writer in `.roundfixrc.yml` that names no command fails it once its
  record is stale, as the coverage writer did.
  Declared break: the fixture no longer writes
  `docs/workflow/authorizations/regeneration-fixture.md` or
  `internal/spec/regeneration_fixture_test.go`.
- task_02: `internal/cli/deliver_derived_archived_item_test.go`
  `TestAnArchivedItemsDerivedRecordConflictIsMergedAndRegenerated` and
  `TestAnArchivedItemWithoutTheSanctionKeepsItsDerivedRecordConflict`. The
  existing conflict tests (`TestADerivedConflictIsMergedAndRegenerated`,
  `TestConflictRecoveryAbortsAFailedRegeneration`) run in the same
  Verification.
- Sabotage proofs belong to the QA gate. Each one is observable through the
  named case:
  - removing the declaration fails `TestRegenerationIsDeclared` on the
    coverage command's exit;
  - injecting it again in the fixture makes the unsanctioned subtest fail;
  - letting `removeSanctionedRegeneration` exempt a write when no command was
    declared makes the unsanctioned subtest and the refusal test fail.

## Build Order

1. The coverage writer names its command; the regeneration contract models
   an archived item and its unsanctioned case; the rule and glossary
   (task_01).
2. The queue-level archived-item fixture and the coverage record re-record
   (task_02) (depends on: 1 — the resolution test passes only once the
   writer declares its command).
3. QA gate (task_03) (depends on: 1, 2).

## Risks & Considerations

- Cost: each repository copy runs real `go test`. The config contract already
  takes about 58 s and gains one coverage-only copy. The two queue tests run
  in parallel and reuse the build cache (about 3 s per nested command,
  measured).
- A Run worktree holds this Spec's own `_authorization.md`, which sanctions
  the coverage command. Both fixtures remove or skip `docs/`, so that record
  cannot mask the unsanctioned case.
- The guard still honors every record in the tree (ADR-0345, Consequences).
  The fixtures prove that the item's record suffices, not that it is the only
  record honored.

## Glossary

- changes: **Archive Record**

## Decisions

- The writer names its command; the guard's rule is unchanged. See ADR-0345.
- The archived item's record is the sanction after archive. See ADR-0345.
- In-process repository copies, so the test cache tracks the writer.
