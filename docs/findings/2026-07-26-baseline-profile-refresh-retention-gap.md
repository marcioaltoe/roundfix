---
status: pending
created_at: 2026-07-26
updated_at: 2026-07-26
---

# Context-Driven Baseline — profile refresh applied without semantic retention accounting (2026-07-26)

A live Baseline update in `/Users/marcio/dev/fluxus` correctly detected that
the existing Setup Manifest referenced a changed Baseline Profile, but the
approved Change Plan contained zero Upgrade Retention Contract entries while
replacing managed agent guidance. Apply verified every postimage, so the
transaction was internally consistent, but that success did not prove that
the previous Normative Clauses remained represented or were explicitly
rejected.

This report follows the earlier
[live-adoption finding](../docs/findings/2026-07-23-setup-context-driven-adoption-process-improvements.md)
and
[greenfield acceptance target](../docs/findings/2026-07-24-greenfield-agent-guidance-acceptance-target.md).
It records a profile-refresh path that those reports did not cover: the
Baseline identity stayed the same while its embedded Profile and catalog
digests changed.

## Session evidence

- Command: `roundfix baseline` from `/Users/marcio/dev/fluxus` on `main`.
- Initial state:
  `incompatible — the existing Setup Manifest references an unavailable or changed Baseline Profile`.
- Instruction-preservation choice: `Preservation`.
- Baseline Profile: `standard-typescript-monorepo`.
- Repository extension remained disabled:
  `repository.extension.enabled=false`.
- Previous Profile digest:
  `sha256:c6e3009029dfdead6762af19b8eb94b1a6668e9e718d835992e9473daa6a9882`.
- Applied Profile digest:
  `sha256:b75ff11efd71c241ed1f521c6576beb61c24b99d327f4ba4e8bd722ddd2e2d19`.
- Previous catalog digest:
  `sha256:aee2a3768e7ef9d3e2dd6a29c60d73daf05b055f5e1683ed40b9c6509c7f9401`.
- Applied catalog digest:
  `sha256:f082de03ea52e1b17e75daf0d65ef2d535f8d8d36951fe08dd435732e4008c1e`.
- Approved Plan Digest:
  `sha256:eca3b27d5c08742edcbfe4382bced8daaa44f139686e6b9417ff3cb92f851fc9`.
- Apply result: `Baseline apply: verified`, with 15 verified postimages.
- Change Plan: 13 file paths, including the immutable root backup
  `AGENTS.d199b5dadd21d8187af411a1111578d4231043c3d9354893aaaef6e1131d2a3f.md`.
- Upgrade Retention Contract ledger: `0 entries`.
- The managed-file diff contained 244 additions and 378 deletions across
  `AGENTS.md` and 11 `docs/agents/` files.
- Apply did not run the formatter, build, tests, or selected repository
  Verification. It reported them as recommendations.

## 1. The unchanged Baseline identity bypasses retention when Profile semantics change

- Symptom / evidence: The interactive workflow declared the Setup Manifest
  incompatible because the Profile was unavailable or changed. The same plan
  then emitted an empty Upgrade Retention Contract ledger and allowed managed
  guidance to change.
- Root cause: `resolvePlanRetention` in
  [`internal/baseline/plan.go`](../internal/baseline/plan.go) treats the
  Manifest as retention-compatible when its schema, version, and
  `generator.baseline` match the selected Profile identifier. It returns
  before calling `currentSetupManifestProfileIsValid`, which is the check that
  compares the catalog and Profile digests. Both the previous and current
  manifests use
  `baseline.standard-typescript-monorepo-0.0.1`, so changed Profile semantics
  do not require a maintained transition.
- Action / suggestion: Make full Setup Manifest validity the first
  compatibility test. A matching Baseline identifier must not bypass
  retention when `profileDigest`, `catalogDigest`, managed clause identities,
  or managed artifact digests changed. Require a transition keyed by the
  source tuple `(generator.baseline, profileDigest, catalogDigest)` or stop
  planning with an action-required result.

Proposed invariant:

> An incompatible Setup Manifest may produce a ready update plan only when
> every previous managed Normative Clause has one explicit retention
> disposition.

Add a regression fixture where the Baseline identifier remains unchanged,
the embedded Profile digest changes, and one managed clause disappears. The
plan must either contain complete retention accounting or exit action-required;
an empty retention ledger must never be ready.

## 2. Preservation does not account for the managed nested guides being replaced

- Symptom / evidence: The operator selected Preservation, but the plan asked
  for no clause classifications and reported zero retained entries. It also
  repeated `baseline.inventory.nested-carrier-conflict` for every generated
  guide under `docs/agents/`, including files the same plan updated.
- Root cause: [`internal/baseline/repository.go`](../internal/baseline/repository.go)
  classifies every path under `docs/agents/` as an instruction carrier, then
  warns for every nested carrier except the recognized repository-extension
  paths. The preservation tests explicitly keep nested sources out of the
  Source Baseline. The inventory does not distinguish:
  - a current Manifest-owned generated guide;
  - a stale Manifest-owned generated guide;
  - a repository extension;
  - an unmanaged nested instruction carrier.
- Action / suggestion: Classify carriers against the existing Setup Manifest
  before emitting warnings. Do not emit a nested-carrier conflict for a
  current managed artifact whose digest matches the Manifest. Treat a stale
  managed artifact as upgrade input and require clause-level retention. Keep
  the current warning for unmanaged nested carriers and repository-authored
  bytes outside setup-owned markers.

The warning text must also describe the real write boundary. A file that
contains a managed block can change while its repository-authored bytes remain
untouched. Replace `nested carrier remains unchanged` with a diagnostic that
states which bytes are managed, which bytes are preserved, and which bytes
were excluded from retention accounting.

## 3. File-level hashes hide removed or weakened Normative Clauses

- Symptom / evidence: The consolidated review showed exact before/after file
  hashes and managed entry identifiers, but it did not show the semantic
  effect of the update. The resulting diff removed or compressed rules for:
  - backend dependency direction, thin HTTP handlers, and persistence
    boundaries;
  - frontend system boundaries and side-effect-free public entries;
  - detailed Task ownership, QA readiness, and archive sequencing;
  - Secondbrain degraded-mode behavior;
  - several documentation lifecycle details.
- Root cause: Managed entry accounting proves byte ownership and transaction
  integrity, but it does not compare the prior and next clause catalogs for
  semantic continuity. The current Change Plan can therefore be exact without
  making rule removal visible at the level the maintainer must approve.
- Action / suggestion: Add a `Normative Clause changes` section before final
  confirmation. For every previous clause, show one of:
  `retained`, `moved`, `replaced`, `repository-document`,
  `repository-extension`, `reasoned-rejection`, or `unaccounted`.
  Include enforcement-strength changes and the new semantic owner.

Final confirmation must report:

```text
Prior clauses: <count>
Accounted clauses: <count>
Unaccounted clauses: <count>
Repository-specific classifications requiring review: <count>
```

Do not offer Apply while `Unaccounted clauses` is greater than zero.

## 4. Executable capability discovery rejects common package-manager symlinks

- Symptom / evidence: Profile alignment reported missing `rtk` and Docker
  evidence. In the same environment, `rtk` resolved to
  `/opt/homebrew/bin/rtk` and Docker to `/usr/local/bin/docker`. Both paths are
  executable symlinks:
  - `/opt/homebrew/bin/rtk -> ../Cellar/rtk/0.43.0/bin/rtk`;
  - `/usr/local/bin/docker -> /Applications/Docker.app/Contents/Resources/bin/docker`.
- Root cause:
  [`lookPathWithoutExecution`](../internal/baseline/profile_alignment.go)
  uses `os.Lstat` and accepts only a regular executable file. It rejects the
  symlink itself without inspecting a bounded symlink chain or reporting why
  the candidate was insufficient. Homebrew and Docker Desktop commonly expose
  executables through symlinks.
- Action / suggestion: Discover the PATH candidate without executing it, then
  resolve its symlink chain and validate that the final target is a regular
  executable file. Reject cycles, missing targets, non-executable targets, and
  unsupported file types with an explicit evidence diagnostic. Record both
  the PATH candidate and resolved target.

Add coverage for:

- a direct regular executable;
- a relative symlink to a regular executable;
- an absolute symlink to a regular executable;
- a broken symlink;
- a symlink cycle;
- a regular non-executable target.

The `rg` warning could not be classified from the interactive transcript
alone because the user's terminal PATH may differ from the Agent environment.
Capability output must expose the inspected PATH candidate or state that no
candidate was found, so maintainers can distinguish missing software from a
rejected candidate.

## 5. Optional and recommended divergences read like required remediation

- Symptom / evidence: Profile alignment was `ready`, but each optional or
  recommended divergence used imperative next-action text. This made Docker,
  Inngest, `rg`, and `rtk` appear to be unfinished adoption work even though
  none blocked planning or Apply.
- Root cause: The text presentation does not lead with requirement strength
  and completion impact. `capability.optional.missing` and
  `capability.recommended.missing` are technically distinct, but their
  remediation language resembles a blocking diagnostic.
- Action / suggestion: Group divergences under `Blocking`, `Advisory`, and
  `Informational`. For every advisory, render `Does not block Baseline
  readiness or Apply` before any optional next action. Do not recommend adding
  Inngest evidence unless the Profile module is selected for repository use.

## 6. The workspace Verification expectation duplicates a valid selected gate

- Symptom / evidence: The Profile expected `bun run verify`, while Fluxus
  selected and declared `rtk make verify`. The Setup Manifest correctly marked
  the selected gate as repository-executable and the portable workspace
  command as unresolved, but the interactive result still presented the
  latter as a divergence.
- Root cause: The Profile owns a fixed workspace command, while the
  `verification.gate` decision is stored as a separate role. There is no typed
  mapping that lets the selected repository gate satisfy the portable
  workspace Verification role.
- Action / suggestion: Let the maintainer map each portable Verification role
  to a declared repository command during alignment. If the selected
  repository gate intentionally serves the workspace role, persist that
  mapping instead of requiring a duplicate package script. Keep unmapped
  profile expectations advisory and never present them as commands Baseline
  executed.

## 7. Apply success does not distinguish byte verification from semantic readiness

- Symptom / evidence: The final headline was `Baseline apply: verified`,
  followed by repeated retention warnings and unrun repository Verification
  recommendations. The headline can be read as proof that the update is fully
  accepted even though it proves only that approved postimages were written
  and verified.
- Root cause: One success word covers transaction integrity, Profile
  alignment, semantic retention, idempotence, and repository Verification,
  although these are separate states.
- Action / suggestion: Render a final status matrix:

```text
Approved postimages: verified
Semantic retention: verified | action required
Profile alignment: ready | blocked, with advisory count
Repository Verification: passed | failed | not run
Idempotence check: passed | not run
```

Use `Baseline update complete` only when semantic retention is verified and
the completion contract's idempotence step has passed. Keep repository
Verification as an explicit external result because Baseline must not execute
it during Apply.

## Recommended implementation order

1. Block same-identity Profile drift when retention accounting is absent.
2. Add a regression fixture for changed Profile and catalog digests with an
   unchanged Baseline identifier.
3. Separate managed, stale-managed, repository-owned, and unmanaged nested
   carriers.
4. Add the clause-level semantic delta to interactive and JSON plans.
5. Accept safe executable symlinks and expose rejected-candidate evidence.
6. Clarify advisory presentation and Verification role mapping.
7. Split the final result into transaction, retention, alignment,
   Verification, and idempotence states.

## What worked — keep

- The interactive workflow detected the changed Profile instead of silently
  treating the old digest as current.
- Final mutation required approval of the exact Plan Digest.
- Apply verified exact postimages and created the immutable root carrier
  backup.
- File changes appeared before the complete managed-entry ledger.
- Optional and recommended capability gaps did not block a ready Profile.
- Formatter, build, tests, and repository Verification remained
  recommendations; Apply did not execute repository commands.
- The repeated nested-carrier warnings exposed the unsafe area, even though
  their classification and wording need correction.

## Suggested acceptance checks

- Given an existing Manifest with the same Baseline identifier and an older
  Profile digest, when a managed clause changes, then a plan with empty
  retention exits action-required.
- Given a Manifest-owned guide whose digest matches, profile alignment emits
  no nested-carrier conflict for that guide.
- Given an unmanaged nested instruction carrier, planning keeps it unchanged
  and emits one path-specific conflict warning.
- Given `rtk` or Docker installed through a valid executable symlink,
  capability discovery reports satisfied evidence without executing the
  binary.
- Given a selected repository Verification gate mapped to the workspace role,
  the portable `bun run verify` expectation does not remain an unresolved
  divergence.
- Given successful Apply with repository Verification and idempotence not run,
  the result reports verified postimages without claiming the Baseline update
  complete.
