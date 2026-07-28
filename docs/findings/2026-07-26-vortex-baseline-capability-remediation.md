---
status: done
created_at: 2026-07-26
updated_at: 2026-07-28
---

# Context-Driven Baseline — capability divergences do not carry enough evidence to be remediated (2026-07-26)

A Baseline alignment against `standard-typescript-monorepo` in
`/Users/marcio/dev/vortex` returned `action_required` with two blocking and
five advisory divergences. Every divergence was individually correct, but the
report did not carry enough information to act on any of them. Resolving the
two blocking items required reading
`internal/baseline/assets/profiles/standard-typescript-monorepo.json` and
`internal/baseline/profile_alignment.go` to learn what each probe actually
inspects. That source reading, not the remediation itself, was the dominant
cost of the session.

The three divergence classes in this run needed three unrelated responses: one
needed a documentation file, one needed a real dependency migration, and five
needed nothing at all. Prompt 16 offered three whole-repository options, none
of which was "fix the repository and re-run".

This report corroborates and extends
[the fluxus profile-refresh finding](2026-07-26-baseline-profile-refresh-retention-gap.md).
Findings 4, 5, and 6 of that report (symlink capability discovery, advisory
presentation, workspace Verification mapping) reproduced here in a second
repository and a different environment. Only the new observations are
expanded below; section 6 of this report closes the open `rg` classification
question that report left unresolved.

## Session evidence

- Roundfix `0.0.1 (8bcd61f, built 2026-07-26 10:10:16 -0300)`, macOS 25.5.0.
- Repository `/Users/marcio/dev/vortex`, profile `standard-typescript-monorepo`,
  state `action_required`.
- Blocking: `capability.stack.logtape` (`capability.required.missing`),
  `capability.stack.postgresql` (`capability.contract.missing`).
- Advisory: `capability.optional.docker`, `capability.rg`, `capability.rtk`,
  `verification.build`, `verification.workspace`.
- Repository facts at alignment time: backend used `pino ^10.3.1` and
  `winston ^3.19.0` behind a `Logger` port with three adapters; `postgres
  ^3.4.9` with Drizzle and a `postgres:18-alpine` compose service; no
  `DATABASE.md` and no `docs/architecture/` directory; root `package.json`
  declared `fmt`, `lint`, `test`, and `typecheck` but neither `build` nor
  `verify`; `make verify` existed and was the selected `verification.gate`.
- `roundfix baseline plan --profile standard-typescript-monorepo --format json`
  exited `3` on missing decisions and never reached capability evaluation.

## 1. Divergence output does not name the probe it evaluated

- Symptom / evidence: `capability.stack.logtape` reported "LogTape has no
  compatible local evidence" and `capability.stack.postgresql` reported
  "PostgreSQL implementation evidence was found, but the required repository
  contract is absent". Neither message states which paths were read or what
  content was expected. The PostgreSQL next action did list three accepted
  paths; the LogTape next action ("Add compatible local LogTape evidence")
  named none.
- Root cause: The Profile's `capabilities[].probe` object carries the complete
  decision input — `kind`, `paths`, and `contains` for `declared-file`, and
  `executable` for `executable` — but alignment renders only the capability
  title and a generic next action. The probe is fully deterministic and
  entirely unexposed.
- Action / suggestion: Render the evaluated probe with every unsatisfied
  divergence. For `declared-file`, list each inspected path with its state
  (`absent`, `present, no match`, `unreadable`) and the expected content. For
  `executable`, print the inspected PATH candidate or state that no candidate
  was found. A maintainer should never need to read catalog assets or Go
  source to learn what a probe wanted.

Shape:

```text
blocking capability.stack.logtape (capability.required.missing)
  probe: declared-file contains "@logtape/logtape"
    package.json                     present, no match
    packages/frontend/package.json   present, no match
    packages/backend/package.json    present, no match
```

## 2. A required stack capability blocks a repository holding a deliberate equivalent

- Symptom / evidence: `capability.stack.logtape` is `strength: required`, so a
  repository that logs through pino and winston behind its own `Logger` port
  cannot reach a ready Baseline. The divergence gives no signal that LogTape
  is a house-stack choice rather than a Baseline mechanism, and no signal of
  what the alternative costs.
- Root cause: The Profile encodes an opinionated stack, which is correct, but
  required stack capabilities are reported with the same vocabulary as
  structural prerequisites. "No compatible local evidence" describes a missing
  file; it does not say "this Profile selects LogTape as the logging stack and
  your repository selects something else".
- Action / suggestion: For a required *stack* capability, state the selected
  technology, the reason the Profile requires it, and both resolutions
  explicitly: adopt it, or remove it through a reviewed repository-owned
  adaptation. Where the catalog can determine that the removal cascades to no
  decision — `capability.stack.logtape` is referenced by no entry in
  `assets/decisions.json` — say so, because that is exactly the fact that
  makes the adaptation cheap and safe.

A second, separable point: the probe is a substring match for
`"@logtape/logtape"` in a `package.json`. It proves declaration, never use. A
repository can satisfy this required capability by adding an unused
dependency. If declaration is the intended bar, name the capability for what
it checks. If use is the intended bar, the probe needs source evidence, as
`capability.stack.postgresql` already demands a contract artifact.

## 3. The contract-artifact probe is the right idea, applied to exactly one capability

- Symptom / evidence: `capability.stack.postgresql` is the only capability
  whose probe reads a documentation artifact (`DATABASE.md`,
  `docs/architecture/database.json`, or `docs/architecture/postgresql.json`
  containing `PostgreSQL`) rather than a dependency string. Its message is the
  most informative in the whole report because it separates implementation
  evidence from contract evidence. The repository had none of the three paths
  and no `docs/architecture/` directory at all.
- Root cause: The two-part implementation-plus-contract check exists in the
  catalog for one capability with no stated rule for when it applies. From the
  outside it reads as an inconsistency rather than a deliberate distinction.
- Action / suggestion: Two changes, independently useful.
  1. State the rule. If durable stores earn a contract requirement because
     schema, migration, and identifier conventions must be written down, say
     that in the capability explanation, and apply it to any future capability
     that meets the rule.
  2. Offer the artifact. A `capability.contract.missing` divergence knows the
     accepted paths and the required content. Baseline should be able to
     propose a contract stub at the first accepted path as a normal Change
     Plan entry, subject to the same digest confirmation as any other write.
     Remediation took a hand-written `DATABASE.md`; a stub with the required
     headings would have made the gap self-closing.

## 4. Prompt 16 has no option for "remediate in the repository and re-run"

- Symptom / evidence: The prompt offered `1. Change Baseline Profile`,
  `2. Create a reviewed repository-owned Profile adaptation`, and `3. Decline
  without writing`. The correct response for `capability.stack.postgresql` was
  to write one file and re-run. The correct response for
  `capability.stack.logtape` was a dependency migration and re-run. The
  correct response for all five advisories was to do nothing. None of the
  three options expresses any of that, so the session had to exit through
  option 3 and act outside the tool.
- Root cause: The prompt resolves divergence *as a set*, with profile-scoped
  choices only. A mixed set whose members need repository-scoped, unrelated
  responses has no representation.
- Action / suggestion: Add an explicit fourth option — exit without writing,
  print the per-divergence remediation, and name the exact command to re-run.
  This is behaviourally close to option 3, but it is the difference between
  "declined" and "paused for repository work", and only one of those is what
  happened.

A related constraint worth documenting on the prompt itself:
`NewProfileAdaptationDraft` rejects a draft with no removals
(`custom.profile.adaptation.removal.required`), so option 2 cannot be used to
remap Verification roles or to acknowledge advisories. Option 2 is a removal
mechanism only.

## 5. There is no read-only way to re-check capabilities after remediation

- Symptom / evidence:
  `roundfix baseline plan --repo . --profile standard-typescript-monorepo --format json`
  exited `3` with `required Baseline decisions are missing: auth.provider,
  autonomous.enabled, domain.layout, http.contract, identifier.strategy,
  language.generated, repository.extension.enabled, secondbrain.enabled,
  spec.scaffold, triage.external, verification.gate`. Decision resolution runs
  before capability evaluation, so confirming that the two blocking probes now
  pass required reproducing the probe logic by hand against the catalog
  instead of asking Roundfix.
- Root cause: Capability evidence depends only on the Profile and bounded
  repository facts. Decisions gate the *plan*, not the probes, but the command
  ordering makes the two inseparable.
- Action / suggestion: Expose capability alignment without decisions — a
  `roundfix baseline capabilities --profile <id>` read-only report, or a
  `--capabilities-only` mode on `plan` that skips decision resolution and
  emits the capability outcomes with the probe detail from finding 1. This is
  the natural verification loop for an agent or a CI check: remediate,
  re-check, repeat. Today that loop does not exist.

## 6. Symlink capability discovery — corroboration, and the `rg` case resolved

[Finding 4 of the fluxus report](2026-07-26-baseline-profile-refresh-retention-gap.md)
recorded that `lookPathWithoutExecution` rejects executable symlinks and could
not classify the `rg` warning from the interactive transcript. This session
reproduced the defect in a second repository and classified all three
capabilities by replaying the probe's exact `os.Lstat`-plus-`IsRegular` logic
over the live PATH:

```text
rg       candidate=<none on PATH>            kind=-        accepted=false
rtk      candidate=/opt/homebrew/bin/rtk     kind=symlink  accepted=false
docker   candidate=/usr/local/bin/docker     kind=symlink  accepted=false
git      candidate=/opt/homebrew/bin/git     kind=symlink  accepted=false
bun      candidate=/Users/marcio/.bun/bin/bun kind=regular accepted=true
```

- `rtk` (0.43.0) and Docker (29.6.2) were installed and working; both are
  false negatives, matching the fluxus evidence.
- `rg` is a **true** negative in this environment. There is no `rg` regular
  file anywhere on PATH. The interactive `rg` is a shell function injected by
  the coding harness, which a `PATH` scan correctly cannot see. Finding 4's
  open question can be closed: the fix must not assume every rejected
  executable capability is a symlink case.
- `git` is included only as a severity signal. Roundfix does not probe it
  today, but Homebrew's `git` is a symlink, so the same probe would reject it
  if a `capability.git` were ever added. On macOS the symlink case is the
  norm, not the exception.

This corroboration should raise the priority of finding 4 rather than open a
separate work item.

## 7. A clean adoption warns about the thirteen files it just authored

- Symptom / evidence: After a `verified` greenfield apply in
  `/Users/marcio/dev/vortex`, an immediate re-plan with identical inputs correctly
  reported `fileChanges: 0` — the adoption is idempotent — but emitted thirteen
  warnings, one per generated guide:

  ```text
  baseline.inventory.nested-carrier-conflict  docs/agents/agent-instructions.md
  baseline.inventory.nested-carrier-conflict  docs/agents/skill-dispatch.md
  baseline.inventory.nested-carrier-conflict  docs/agents/setup-context.json
  … 10 more, covering every file the apply wrote
  ```

  `docs/agents/specific-repository.md` is the only file under that directory that does
  not warn.
- Root cause: The inventory walker in
  [`repository.go:405`](../internal/baseline/repository.go) warns for every carrier with
  `scope == "nested"` unless `isRecognizedRepositoryRuleCarrier` matches, and that
  predicate recognizes only the repository-rules path. Setup-owned guides are nested by
  construction, carry `setup-context-driven:begin` markers, and are listed in the
  Setup Manifest and `managedEntries` of the very plan that wrote them — but the walker
  does not consult either, so it cannot tell its own output from unreviewed
  repository-authored instructions.
- Action / suggestion: Treat a nested carrier as setup-owned when the Setup Manifest
  claims it and its managed markers verify, and warn only for nested carriers that are
  unclaimed or whose markers fail verification. The current behaviour trains maintainers
  to dismiss a thirteen-line warning block after every successful adoption, which is
  exactly how a real conflict gets missed.

This compounds finding 5: with no read-only capability re-check, the idempotent re-plan
is the natural way to confirm an adoption landed — and it is the path that emits the
spurious warnings.

## Recommended implementation order

1. Fix `lookPathWithoutExecution` to resolve a bounded symlink chain
   (fluxus finding 4). Highest severity, now reproduced twice, and it
   silently mis-reports on any Homebrew or Docker Desktop machine.
2. Render probe detail on every unsatisfied divergence (finding 1). This is
   the single change that would have removed most of this session's cost, and
   it also supplies the diagnostic text finding 4 needs.
3. Add a read-only capability re-check (finding 5). Without it, items 1 and 2
   improve the diagnosis but not the remediation loop.
4. Add the "remediate and re-run" outcome to Prompt 16 and document the
   removal-only constraint on the adaptation option (finding 4 above).
5. Clarify required stack capabilities and their cascade (finding 2).
6. Decide whether the contract-artifact probe is a rule or an exception, and
   offer the stub either way (finding 3).

## What worked — keep

- `capability.contract.missing` as a distinct code from
  `capability.required.missing`. Separating "you do not run PostgreSQL" from
  "you run PostgreSQL but have not written the contract" is the most useful
  distinction in the report, and it named its three accepted paths.
- Blocking versus advisory strength was classified correctly for all seven
  divergences. The problem in fluxus finding 5 is the presentation of that
  classification, not the classification itself.
- `baseline plan` refused cleanly with exit `3` and a complete list of missing
  decisions, and emitted no partial plan.
- The format-role alias fallback in `resolveProfileVerificationCommand`
  silently resolved the Profile's `bun run format` against the repository's
  `fmt` script, so `verification.format` never appeared as a divergence. That
  is exactly the mapping behaviour fluxus finding 6 asks for on the `build`
  and `workspace` roles.

## Suggested acceptance checks

- A `declared-file` divergence renders every inspected path with its state and
  the expected content; an `executable` divergence renders the inspected PATH
  candidate or an explicit "no candidate found".
- Executable discovery accepts a relative symlink, an absolute symlink, and a
  direct regular executable; it rejects a broken symlink, a symlink cycle, and
  a non-executable target, each with a distinct diagnostic.
- Capability alignment is obtainable for a Profile with zero decisions
  supplied, and its outcomes match the outcomes produced by a full plan for
  the same repository state.
- A capability whose removal cascades to no decision reports that fact in its
  divergence.
- Prompt 16 exposes an outcome that writes nothing, prints per-divergence
  remediation, and names the re-run command; selecting it is distinguishable
  in the journal from a decline.
- A `capability.contract.missing` divergence can produce a contract stub as a
  Change Plan entry gated by the normal Plan Digest confirmation.
- An idempotent re-plan immediately after a verified apply reports zero file changes and
  zero warnings; a nested carrier warns only when the Setup Manifest does not claim it or
  its managed markers fail verification.

## Addendum — 2026-07-28 — Routed to Spec 0057

All seven findings, including the symlink-discovery corroboration and the
resolved `rg` classification, are owned by
[Spec 0057 — Baseline capability evidence and retention](../specs/0057-baseline-capability-evidence-and-retention/_prd.md).
