# Command evidence — 2026-07-25

## Build and repositories

- Roundfix source:
  `10720c8d1954b6e2e21e22ec880945a28e4721e3`
- Built CLI:
  `rtk go build -buildvcs=false -o /private/tmp/roundfix-spec0048-qa ./cmd/roundfix`
- Fluxus source:
  `1aeed7e8370c3d14137c42b0c789dcbe3bd1ba3b`
- Fluxus source checkout was clean. All Baseline writes occurred in
  `--no-hardlinks` clones under
  `/private/tmp/roundfix-spec0048-fluxus.KKj2C6/`.

## Project Constraint audit

The active PRD and TechSpec both describe identifier strategy,
authentication and HTTP, active ADR obligations, and tooling authority with
applicability reasons. The TechSpec cites operative sources under
`docs/agents/`; the PRD cites none. This fails the active-Spec source contract.

Task 06 commit `7cef20d988efb96f99bc4bcd8f61452f70fa610b`
changed `internal/cli/baseline_release_gate_test.go` and introduced the exact
Oxfmt version pin:

```go
const baselineReleaseOxfmtVersion = "0.59.0"
```

The same change generates a fixture `package.json` with that dependency pin.
The PRD and TechSpec both state that this Spec authorizes no tooling
configuration, script, ignore file, plugin declaration, or version pin. No
bounded authorization names the changed test file. The audit therefore fails
before accepting the Task's reported evidence.

Task 08 commit `10720c8d1954b6e2e21e22ec880945a28e4721e3`
changes only `task_08.md` from `status: pending` to `status: completed`.
The Task has no `## Result`, and all five subtasks and six acceptance
checkboxes remain open.

## Static and focused gates

The unmodified repository gate failed for an environment reason:

```text
rtk make verify
Go test: 2291 passed, 2 failed, 1 skipped in 22 packages
TestGuidanceCompositionJourney/standard-typescript-monorepo:
bun install v1.3.14
error: bun is unable to write files to tempdir: PermissionDenied
exit 2
```

The exact focused journey reproduced the same denial after setting
`TMPDIR=/private/tmp`. The fixture creates its own temp directory inside the
writable worktree, but Bun still returns `PermissionDenied`. Per the sandbox
contract, the gate was not retry-looped.

Independent focused checks:

```text
Project decision, rendering, and tooling-clause tests: 51 passed
Project Constraint, documentation, asset, and skill tests: 35 passed
Task 08 verification regex: 25 passed
roundfix skills check: passed for all 14 repository-owned skills
make skills-sync-check: 4 passed
baseline assets sync --check --format json: ok=true, 0 findings
project-specific branding search: no matches
git diff --check: passed
```

The Task 08 verification regex matches `TestBaselineReleaseGate`; no
`TestProjectDecisionJourney`, `TestProjectConstraintJourney`, or
`TestToolingAuthorizationJourney` function exists in `internal/` or
`skills/`.

## Public automation behavior

Against the clean Fluxus update clone, the built CLI exited `3`, emitted one
`roundfix/baseline-result/v1` document, named all unresolved values including
`identifier.strategy` and `auth.provider`, emitted no partial Plan, and left
the clone clean:

```text
required Baseline decisions are missing: auth.provider, autonomous.enabled,
domain.layout, http.contract, identifier.strategy, language.generated,
repository.extension.enabled, secondbrain.enabled, spec.scaffold,
triage.external, verification.gate
```

The completed Fluxus Decision Document is
`fluxus-update-decisions.json`. Automation then reached the expected
preservation-classification boundary without mutating the clone.

## Public human update failure

The built root Baseline command ran in a PTY against the same clean Fluxus
update clone. The maintainer selected Preservation, reused the TypeScript
Profile, and kept every displayed default.

The new identifier prompt was:

```text
Keep suggested UUID version 7:
identifier.strategy={"kind":"uuid-v7"}
```

The workflow first kept the persisted HTTP exception whose reason starts with
`Preserve the provider-owned session...`. It then displayed a Better Auth
default whose reason starts with `Session, OAuth redirect...`. Accepting both
defaults ended the workflow:

```text
roundfix: baseline failed: normalize Baseline decisions:
project decisions "auth.provider" and "http.contract":
"auth.provider" conflicts with "http.contract" for owner "Better Auth"
and scope "/api/auth/*"
exit 1
```

No Fluxus clone byte changed. The update journey cannot reach a Plan without
manually changing one of the two product-proposed values.

## Public Fluxus greenfield journey

Automation transformed the same reviewed Fluxus decisions only by changing
`preservation.mode` to `greenfield`. Planning exited `0`:

```text
schema: roundfix/baseline-plan/v1
Plan Digest: sha256:e7ee8a09b41270b0820ba4b29c828ebfe64ce88f041baad867a70ae2f77ce6ef
File changes: 13
Managed entries: 26
```

The normalized Plan contains `identifier.strategy={"kind":"uuid-v7"}`,
`auth.provider`, and one identical Better Auth exception inside
`http.contract`. Applying an all-zero digest exited `3` and left the clone
clean. Applying the exact digest exited `0` with `state: verified`.

Fresh reads of generated public guidance confirmed:

- `docs/agents/domain.md` scopes UUID version 7 to new project-owned Internal
  Identifiers and preserves provider, protocol, natural, and business
  identifiers.
- `docs/agents/backend.md` contains the complete ordered HTTP contract and
  Better Auth provider rule.
- `docs/agents/agent-instructions.md` contains the universal tooling
  authorization clause.
- `docs/agents/spec-routing.md` requires all four Project Constraint areas and
  operative `docs/agents/` sources.

`rtk bun run fmt` passed on 589 files. The first Fluxus `rtk make verify`
reported missing local `turbo`; `rtk bun install --frozen-lockfile` restored
1,203 lockfile-defined packages, and the exact Verification rerun passed.

Fresh planning then exited `0` with zero file changes and Plan Digest
`sha256:c3ea85f164a850ded9a3fe390eaace29bdf920c704c417781f99f269c9516d03`.
Exact-digest reapply exited `0` with:

```text
state: verified
message: approved Baseline Plan is already applied and verified
verifiedPostimages: 14
```
