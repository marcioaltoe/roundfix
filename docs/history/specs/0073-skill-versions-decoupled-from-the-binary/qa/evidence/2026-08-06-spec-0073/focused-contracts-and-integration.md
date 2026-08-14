# Focused contracts and integration

The focused cross-package command exited 0:

```text
rtk proxy env GOCACHE=/private/tmp/roundfix-qa-0073-focused-cache go test ./skills ./internal/baseline ./internal/cli -count=1 -run '^(TestReadinessComparesDeclaredVersionToMinimum|TestCatalogRejectsMissingOwnedSkillMinimum|TestCatalogRegenerationMode|TestRegenerationBreaksGoldenDigestCycle|TestCatalogDigestExcludesOwnedSkillContent|TestCatalogCompatibility|TestBaselineCompatibilityCorpus|TestReadoptionCompatibilityMaintainedFixture|TestCharacterizationCorporaDoNotRecordOwnedSkillDigests|TestRepositoryReadinessNeverComparesThirdPartySkillVersions|TestOwnedSkillBundleReadinessKeepsStatesDistinct|TestRunDoctorRepositorySkillReadiness|TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness)$' -v
```

It passed all selected tests and their state-specific subtests. The evidence
covers equal/newer/below/unversioned comparison, unreachable-source error
identity, shared Doctor/Skills Check results, third-party call-ledger
exclusion, owned-content-independent catalog identity, characterization-corpus
digest absence, strict generated-guide loading, regeneration-only relaxation,
the pre-Spec compatibility corpus, and readoption compatibility.

The Spec's defining macro did not pass:

```text
rtk proxy env GOCACHE=/private/tmp/roundfix-qa-0073-integration-cache go test -tags=integration ./skills -run '^TestOwnedSkillEditLeavesMakeVerifyGreen$' -count=1 -v
exit 1
```

Its isolated tracked repository edited owned Skill bytes and launched the exact
`make verify` target. That nested gate failed in an unrelated daemon test:

```text
TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement
expected zero Verification commands for queued Task, total starts=2
```

The focused failing test then passed 20 isolated repetitions, and the complete
`internal/daemon` package passed 10 repetitions with `-parallel=16`. This makes
the observed failure load-sensitive and nondeterministic; it does not prove an
environment cause and does not erase the failed owned-Skill edit journey.

