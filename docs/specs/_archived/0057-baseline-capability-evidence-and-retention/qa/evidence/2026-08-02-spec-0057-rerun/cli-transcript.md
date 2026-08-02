# Spec 0057 QA command evidence

Build under test: `981fc945f8433356f988bbb627a08fed65b90e46`.

All commands were run from the Run Worktree unless a scratch repository is
named. Scratch repositories were under `/private/tmp` and used local-only Git
history. No commit was created in the Run Worktree.

## Static and focused verification

```text
$ rtk make verify
3056 Go tests passed in 24 packages
4 Skill tests passed
roundfix skills check passed
exit 0

$ rtk proxy env GOCACHE=/private/tmp/roundfix-0057-qa-gocache go test ./internal/baseline ./internal/cli -run 'Test(SameIdentityDriftRequiresRetention|ReadyPlanNeverCarriesEmptyLedger|ClauseDeltaRendersBeforeLedger|ExecutableCandidateResolution|ExecutableCandidateNeverExecutes|DivergenceCarriesProbeEvidence|DivergenceRendersProbe|DivergenceGroupsByRequirement|PortableVerificationRoleMapping|CapabilityRecheck|CapabilityRecheckMatchesFullPlan|DivergencePromptRemediateOutcome|DivergencePromptJournalsDistinctly|CarrierClassification|UnclassifiableCarrierStillWarns|ResultStatusMatrix|CompletionLanguageRequiresRetention|BaselinePlanCharacterization|BaselineMacroJourneysPublicCLI)' -count=1
ok roundfix/internal/baseline 20.793s
ok roundfix/internal/cli 14.411s
exit 0
```

## Public capability re-check

```text
$ ./bin/roundfix baseline capabilities check --repo . --format json
state: ready
profile: go-cli-tui
context7: satisfied
exa-web-search: satisfied
firecrawl: advisory, missing
rg: satisfied
rtk: satisfied; sourcePath=/opt/homebrew/bin/rtk
exit 0

$ ./bin/roundfix baseline capabilities check --repo . --format text
Baseline capability re-check: ready
Profile: go-cli-tui
Advisory divergences:
- capability.firecrawl ...
  This advisory does not block readiness or apply.
exit 0
```

The only Run Worktree status entry after the read-only command was this QA
report. A Git-initialized scratch repository with no Setup Manifest returned
exit 2 and named `no resolvable Baseline Profile`.

The `standard-typescript-monorepo` re-check returned exit 3 and rendered
Blocking, Advisory, and Informational sections. Its JSON divergences grouped as
20 blocking, 1 advisory, and 1 informational item. Declared-file evidence
included inspected path, state, and expected content; the selected stack named
Better Auth, repository remediation, Profile adaptation, and the decision
cascade.

## Executable discovery and finding reproduction

A scratch PATH contained the chain `rtk -> rtk-hop -> rtk-target`. The target
would create `/private/tmp/roundfix-0057-executed-marker` if invoked.

```text
$ PATH=/private/tmp/roundfix-0057-fake-bin:/usr/bin:/bin <absolute-bin>/roundfix baseline capabilities check --profile go-cli-tui --repo . --format json
rtk: satisfied
sourcePath: /private/tmp/roundfix-0057-fake-bin/rtk
exit 0

$ test ! -e /private/tmp/roundfix-0057-executed-marker
exit 0
```

After changing the candidate to each invalid state, JSON reported distinct
`broken-link`, `link-cycle`, and `not-executable` details. For the broken-link
case, the corresponding public text output was:

```text
Advisory divergences:
- capability.rtk (capability.evidence.insufficient): rtk local evidence is present but insufficient.
  This advisory does not block readiness or apply.
  Next action: Install rtk and ensure it is available on PATH.
```

The text omitted both the inspected candidate path and `broken-link`, although
both values were present in the same command's JSON result.

## Public plan, apply, status, and fail-closed checks

Scratch repository `/private/tmp/roundfix-0057-no-profile.Q56u2f` was planned
with `go-cli-tui`. The approved plan digest was
`sha256:572ab8c6105349467055e0ef996fd2446e62417e57c1a7650299eb95d67d5706`.
Apply and exact reapply both exited 0. The first result reported approved
postimages and Profile alignment verified, with semantic retention, repository
Verification, and idempotence not run. Exact reapply reported idempotence
verified and still did not print completion language because retention was not
run. Re-plan produced zero file changes.

Adding unmanaged `nested/AGENTS.md` produced one warning that bytes inside the
managed markers are managed and bytes outside are preserved. A clean re-plan
had no managed-guide warning.

After adding a declared `make qa-marker` target, plan JSON mapped it to role
`repository-gate`, recorded the exact command and declaration path, and left
`/private/tmp/roundfix-0057-verification-executed-marker` absent. Selecting
`make missing` returned action-required exit 3.

A `rust-cli` transition from the applied Go profile produced a ready plan with
six file changes. Apply exited 0 and the Setup Manifest recorded `rust-cli` and
`baseline.rust-cli-0.0.1`. An all-zero confirmation digest was rejected with
exit 3. Editing `AGENTS.md` after planning made the approved plan stale and was
also rejected with exit 3.

## Interactive remediation outcome

In scratch repository `/private/tmp/roundfix-0057-remediate.jhNcZj`, the real
terminal workflow selected `standard-typescript-monorepo` and then outcome 3,
`Remediate in repository and re-run`. It exited 3, printed remediation for each
divergence and the exact command:

```text
roundfix baseline capabilities check --profile standard-typescript-monorepo --repo /private/tmp/roundfix-0057-remediate.jhNcZj
```

The structured result category was `remediation`, explicitly said that no
bytes were written, and the scratch Git status remained clean.

## Environment limits

The public binary contains one immutable embedded catalog generation. A source
fixture from an older corpus was rejected as an invalid strict Setup Manifest,
so a genuine previously-applied/current-catalog same-identity drift cannot be
manufactured without a second build generation. The focused tests above are
the supervised equivalent evidence for retention, clause ordering, and the
retention-dependent completion state.

The QA prompt states that no Pull Request is open for the Spec target branch;
PR approval, checks, review-thread, Merge-Ready, and ancestry evidence is
therefore unavailable.
