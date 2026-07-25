# Command evidence — build `863af11`

## Environment

- Roundfix source: `863af1100a82c166a314577d7f8de8a50d803db6`
- Built CLI:
  `rtk go build -buildvcs=false -o /private/tmp/roundfix-spec0048-qa-rerun.VyLHao/roundfix ./cmd/roundfix`
- Fluxus source: `1aeed7e8370c3d14137c42b0c789dcbe3bd1ba3b`
- Disposable repositories:
  `/private/tmp/roundfix-spec0048-qa-rerun.VyLHao/fluxus-greenfield`
  and
  `/private/tmp/roundfix-spec0048-qa-rerun.VyLHao/fluxus-update`
- The Fluxus clones used dependency trees from the clean source checkout at
  the same commit because the child-Agent sandbox cannot download packages.
  Root dependencies were read through a symlink; workspace dependency trees
  were copied so Vite could write its fixture temp files. The source Oxfmt
  manifest reports `0.59.0`, matching the lockfile-owned checkout.

## Project Constraint audit

All eight Tasks are `completed` and carry a `## Result`. The active PRD and
TechSpec each account for identifier strategy, authentication and HTTP, active
ADR obligations, and tooling authority with applicability reasons and
operative `docs/agents/` sources.

Daemon-owned Task paths were resolved with
`rtk git diff-tree --no-commit-id --name-only -r <commit>`. Task 06 commit
`7cef20d` temporarily added a test-only Oxfmt constant. Reopened Task 08 commit
`863af11` removed it and reads the existing maintained Profile version instead.
The cumulative Spec diff from base `0ce9144` contains no new protected tooling
configuration, script, ignore file, plugin declaration, or version pin:

- `internal/cli/baseline_release_gate_test.go` removes the pre-existing
  hard-coded `"oxfmt":"0.59.0"` fixture literal and resolves the version from
  the maintained Profile.
- The maintained Profile already owned Oxfmt `0.59.0` at the Spec base; the
  current diff changes its golden digest and selected decisions/rules, not its
  formatter version.
- The current worktree delta contains only this QA report and its evidence.

## Roundfix static and focused gates

The exact full gate is environment-blocked:

```text
rtk make verify
Go test: 2340 passed, 4 failed, 1 skipped in 22 packages
TestGuidanceCompositionJourney/standard-typescript-monorepo:
error: FailedToOpenSocket downloading package manifest oxfmt
TestProjectDecisionJourney/affected_Profile_apply_audit_and_empty_reapply:
error: ConnectionRefused downloading package manifest oxfmt
```

The failing fixture creates an isolated Bun home and cache, then installs
Oxfmt. The sandbox denies its package-manifest request. Per the daemon QA
contract, the exact gate was not rerun through a network bypass.

Independent fresh checks:

```text
Task 01 focused checks: 18 passed
Task 02 focused checks: 9 passed
Task 03 focused checks: 10 passed
Task 04 focused checks: 15 passed
Task 05 and 06 focused checks: 12 passed
Task 07 focused checks: 23 passed
Task 08 engine checks: 59 passed
Task 08 authoring/tooling skill checks: 10 passed
Task 08 public CLI checks not requiring download: 5 passed
roundfix skills check: all 14 repository-owned skills passed
git diff --check: passed
```

`rtk proxy go test ... -list` reported the required named tests:

```text
TestProjectDecisionJourneyEngine
TestToolingAuthorizationJourneyCoreClause
TestBaselineReleaseGate
TestProjectDecisionJourney
TestProjectConstraintJourney
TestToolingAuthorizationJourney
```

The canonical asset audit passed:

```json
{
  "ok": true,
  "summary": {
    "errors": 0,
    "decisions": 0,
    "warnings": 0,
    "info": 0
  },
  "plannedChanges": []
}
```

## Public fail-closed behavior

Planning against the clean Fluxus update clone without decisions exited `3`,
returned `roundfix/baseline-result/v1` with `state: action_required`, named all
11 unresolved stable IDs including `identifier.strategy` and `auth.provider`,
emitted no Plan, and left the clone clean.

A wrong all-zero Plan Digest exited `3`, named the expected digest, and left
the greenfield clone clean.

A complete Decision Document whose `auth.provider` rationale conflicted with
`http.contract` exited `2` before a Plan, named both decision IDs plus owner
`Better Auth` and scope `/api/auth/*`, and left the conflict clone clean.

## Fluxus greenfield journey

The fresh automation Plan completed:

```text
schemaVersion: roundfix/baseline-plan/v1
Plan Digest: sha256:e7ee8a09b41270b0820ba4b29c828ebfe64ce88f041baad867a70ae2f77ce6ef
File changes: 13
Managed entries: 26
```

The Plan contains `identifier.strategy={"kind":"uuid-v7"}`, the complete
Better Auth provider decision, and one identical provider exception in
`http.contract`. Exact-digest apply returned `state: verified`.

Fresh reads confirmed:

```text
docs/agents/domain.md:
Use UUID version 7 for new project-owned Internal Identifiers only.
Preserve external provider identifiers, protocol identifiers, natural keys,
and business codes according to their source contracts.

docs/agents/backend.md:
Application HTTP mode: Post-only.
Better Auth owns GET and POST for /api/auth/* with the confirmed session,
OAuth redirect, callback, and related protocol rationale.

docs/agents/agent-instructions.md:
The universal tooling clause names every protected mutation verb, tooling
category, artifact category, express authorization, and non-authorizing input.
```

`rtk bun run fmt` completed on 589 files. With the same-commit dependencies
pre-provisioned, `rtk make verify` passed. A fresh post-Verification Plan had
zero file changes and 14 postimages. Exact-digest reapply returned:

```text
state: verified
message: approved Baseline Plan is already applied and verified
Plan Digest: sha256:c3ea85f164a850ded9a3fe390eaace29bdf920c704c417781f99f269c9516d03
verifiedPostimages: 14
```

## Fluxus update journey

The public root command ran in a PTY. The maintainer selected Preservation,
reused the TypeScript Profile, and accepted every displayed default. Prompts
6–7 displayed the exact same persisted Better Auth rationale:

```text
Preserve the provider-owned session, OAuth redirect, callback, and related
protocol semantics.
```

Unlike the prior build, the flow reached the complete Change Plan without
manual decision repair. Exact confirmation applied digest
`sha256:76097b9d63f1b29547a426ea10d8bf76b3ed5ad60402062a0dab3b35e4b51df2`
and reported 15 verified postimages.

`rtk bun run fmt` completed on 589 files and `rtk make verify` passed with the
same-commit dependencies pre-provisioned. A second public root run accepted all
persisted defaults, displayed no file changes, and exact confirmation returned:

```text
Result: approved Baseline Plan is already applied and verified
Plan Digest: sha256:b17834ab2c8fc9a04cd21360774282ad9de9b1abb69b31dbd6fb538856a44b1c
Verified postimages: 14
```
