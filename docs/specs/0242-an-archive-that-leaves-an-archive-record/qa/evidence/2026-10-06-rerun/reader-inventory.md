# Audited reader inventory

111 history-string paths and 70 archive-root Go paths; union 136 paths. New helper archive_reader.go is the planned Git recovery seam, covered by squash/leftover regression tests. New record tests and authoring artifacts do not add unrelated runtime readers.

| Path | Disposition | Named regression or reason |
| --- | --- | --- |
| .agents/skills/archive-spec/SKILL.md | changed with regression | Guidance/digest/link source validated by TestArchiveRecordClausesAreAppended, TestSettlementGuidanceIsOneTable, mirror comparisons, baseline-digests or pinned corpus gates; guidance corrections verified in this rerun. |
| .agents/skills/brainstorming/SKILL.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| .agents/skills/roundfix/references/archive.md | changed with regression | Guidance/digest/link source validated by TestArchiveRecordClausesAreAppended, TestSettlementGuidanceIsOneTable, mirror comparisons, baseline-digests or pinned corpus gates; guidance corrections verified in this rerun. |
| .agents/skills/roundfix/references/spec.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| .agents/skills/write-idea/SKILL.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| .agents/skills/write-prd/SKILL.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| .agents/skills/write-tasks/SKILL.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| .agents/skills/write-techspec/SKILL.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| .coderabbit.yaml | unaffected | Path filter/export scope preserved; .secondbrain-export migration is owned by Spec 0243. |
| .secondbrain-export | unaffected | Path filter/export scope preserved; .secondbrain-export migration is owned by Spec 0243. |
| CHANGELOG.md | unaffected | Historical provenance paths resolve in Git; no live archive-folder assertion. |
| docs/adr/0120-retired-documentation-lives-under-one-history-root.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0138-spec-run-git-policy-commit-per-task-push-only-at-clean.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0139-one-active-run-per-work-target-keyed-in-the-run-database.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0140-an-agent-selection-is-proven-as-the-exact-advertised-tuple.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0141-review-runs-execute-in-a-clean-tracked-user-checkout.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0142-head-bound-review-source-evidence-decides-the-watch-outcome.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0143-release-planning-is-read-only-including-its-reset-mode.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0144-setup-is-declarative-manifest-identity-and-declared-decision-effects.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0145-a-managed-region-is-trusted-by-its-marker-and-refreshed-when-unrecorded.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0146-release-publication-is-all-or-nothing-with-a-bounded-token-fallback.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0147-the-adapters-refusal-fires-first-and-the-catalogue-is-the-net.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0149-one-regeneration-declaration-the-grant-names-the-command-the-tree-names-outputs.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0163-review-artifact-retirement-reads-recorded-evidence.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0173-a-baseline-plan-reports-the-citations-its-history-relocations-break.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0230-an-archived-spec-keeps-its-relative-links.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md | new authoring declaration | This Spec/ADR describes the reader contract; not another runtime reader. |
| docs/agents/docs-layout.md | changed with regression | Guidance/digest/link source validated by TestArchiveRecordClausesAreAppended, TestSettlementGuidanceIsOneTable, mirror comparisons, baseline-digests or pinned corpus gates; guidance corrections verified in this rerun. |
| docs/agents/specific-repository.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| docs/references/archived-evidence-measurement.md | unaffected | Historical provenance paths resolve in Git; no live archive-folder assertion. |
| docs/references/grok-through-cursor-measurement.md | unaffected | Historical provenance paths resolve in Git; no live archive-folder assertion. |
| docs/references/jev-router-measurement-2026-10-02.md | unaffected | Historical provenance paths resolve in Git; no live archive-folder assertion. |
| docs/references/jev-router-measurement-2026-10-04.md | unaffected | Historical provenance paths resolve in Git; no live archive-folder assertion. |
| docs/specs/0242-an-archive-that-leaves-an-archive-record/_authorization.md | new authoring declaration | This Spec/ADR describes the reader contract; not another runtime reader. |
| docs/specs/0242-an-archive-that-leaves-an-archive-record/_prd.md | new authoring declaration | This Spec/ADR describes the reader contract; not another runtime reader. |
| docs/specs/0242-an-archive-that-leaves-an-archive-record/_techspec.md | new authoring declaration | This Spec/ADR describes the reader contract; not another runtime reader. |
| docs/specs/0242-an-archive-that-leaves-an-archive-record/references/2026-10-06-history-keeps-only-what-the-secondbrain-needs.md | new authoring declaration | This Spec/ADR describes the reader contract; not another runtime reader. |
| docs/specs/0242-an-archive-that-leaves-an-archive-record/task_01.md | new authoring declaration | This Spec/ADR describes the reader contract; not another runtime reader. |
| docs/specs/0242-an-archive-that-leaves-an-archive-record/task_02.md | new authoring declaration | This Spec/ADR describes the reader contract; not another runtime reader. |
| docs/specs/0242-an-archive-that-leaves-an-archive-record/task_03.md | new authoring declaration | This Spec/ADR describes the reader contract; not another runtime reader. |
| docs/specs/0242-an-archive-that-leaves-an-archive-record/task_04.md | new authoring declaration | This Spec/ADR describes the reader contract; not another runtime reader. |
| docs/specs/0242-an-archive-that-leaves-an-archive-record/task_05.md | new authoring declaration | This Spec/ADR describes the reader contract; not another runtime reader. |
| docs/user-guide/commands/archive.md | changed with regression | Guidance/digest/link source validated by TestArchiveRecordClausesAreAppended, TestSettlementGuidanceIsOneTable, mirror comparisons, baseline-digests or pinned corpus gates; guidance corrections verified in this rerun. |
| docs/user-guide/commands/spec.md | unaffected | User guide archive-root prose or supported history commands; legacy links checked with pinned-history source. |
| docs/user-guide/context-driven-development.md | changed with regression | Guidance/digest/link source validated by TestArchiveRecordClausesAreAppended, TestSettlementGuidanceIsOneTable, mirror comparisons, baseline-digests or pinned corpus gates; guidance corrections verified in this rerun. |
| internal/authorization/authorization_test.go | changed with regression | TestAuthorizationParseCharacterization, TestAuthorizationParsesPreconditionRepairs; pinned Git corpus or disposable fixtures |
| internal/baseline/apply_test.go | unaffected | TestApplyCanonicalizesLegacyRepositoryRules, TestRepositoryRuleBlockPreservesRepositoryEditAndEmptyReapply; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md | changed with regression | Guidance/digest/link source validated by TestArchiveRecordClausesAreAppended, TestSettlementGuidanceIsOneTable, mirror comparisons, baseline-digests or pinned corpus gates; guidance corrections verified in this rerun. |
| internal/baseline/assets/modules/context-workflow.json | changed with regression | Guidance/digest/link source validated by TestArchiveRecordClausesAreAppended, TestSettlementGuidanceIsOneTable, mirror comparisons, baseline-digests or pinned corpus gates; guidance corrections verified in this rerun. |
| internal/baseline/assets/modules/spec-workflow.json | changed with regression | Guidance/digest/link source validated by TestArchiveRecordClausesAreAppended, TestSettlementGuidanceIsOneTable, mirror comparisons, baseline-digests or pinned corpus gates; guidance corrections verified in this rerun. |
| internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md | unaffected | Frozen source-baseline or receipt-characterization fixture; bytes remain intentionally pinned, not read from the live archive root. |
| internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md | unaffected | Frozen source-baseline or receipt-characterization fixture; bytes remain intentionally pinned, not read from the live archive root. |
| internal/baseline/derived_ownership_test.go | changed with regression | TestOutputsForCommand, TestOutputsForSkillsSyncListsOwnedSkillMirrors; pinned Git corpus or disposable fixtures |
| internal/baseline/history_citations_test.go | unaffected | TestRelocationCitationsReportEachCitationForm, TestRelocationCitationsReportARelocatedFilesOutwardLink; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/baseline/history_layout.go | unaffected | Relocates arbitrary retired files; record filenames need no Spec-folder read. |
| internal/baseline/history_layout_test.go | unaffected | TestDiscoverHistoryLayoutLegacyArchiveShapes, TestDiscoverHistoryLayoutCurrentLayoutReportsNothing; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/baseline/plan_test.go | unaffected | TestPlanDocumentStrictCodecs, TestHistoryRelocationPlanCarriesOrderedIdentitiesOutsideRenderedCarriers; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/baseline/repository_test.go | unaffected | TestInventoryWalkIgnoresTransientErrorsInsideExcludedTrees, TestInventoryBudget; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/baseline/transaction_test.go | unaffected | TestTransactionStagesBeforeMutation, TestTransactionRollback; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/archive.go | changed with regression | TestArchiveCommandLeavesTheArchiveRecord / TestArchivePlanWritesNothing |
| internal/cli/archive_record_readers_test.go | changed with regression | TestInspectItemReadsTheArchiveRecord, TestPrerequisitesCountAnArchiveRecord, TestReviewCorrectionAllowsTheArchiveRecordPath, TestReviewParksForACorrectiveSpecOnAnArchiveRecord; pinned Git corpus or disposable fixtures |
| internal/cli/archive_test.go | changed with regression | TestArchiveAcceptsARecordedSupersession, TestRunArchiveMovesCompletedSpecAndStampsMetadata; pinned Git corpus or disposable fixtures |
| internal/cli/baseline_apply_test.go | unaffected | TestBaselineApplyCommand, TestBaselineApplyTextReportsHistoryMoves; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/baseline_plan_test.go | unaffected | TestBaselinePlanPreflightJSONActionRequired, TestBaselinePlanPreflightText; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/baseline_update_test.go | unaffected | TestBaselineUpdateAppliesManifestPlanAndReportsJSON, TestBaselineUpdateTextReportsHistoryMoves; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/deliver.go | changed with regression | TestInspectItemReadsTheArchiveRecord / TestPrerequisitesCountAnArchiveRecord |
| internal/cli/deliver_review_correction.go | changed with regression | TestReviewCorrectionAllowsTheArchiveRecordPath |
| internal/cli/deliver_review_correction_test.go | unaffected | TestACorrectionOutsideTheArchivedSpecIsRefused; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/deliver_workflow.go | changed with regression | TestDeliveryAcceptsAnExactRetirement / TestDeliveryRefusesAnInexactRetirement |
| internal/cli/reconcile_item_branch_test.go | unaffected | TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec, TestReconcileKeepsTheItemBranchOfALiveQueueItem; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/review.go | changed with regression | TestReviewDiffOmitsTheRemovedSpecFolder / TestReviewParksForACorrectiveSpecOnAnArchiveRecord |
| internal/cli/review_archived_spec_test.go | unaffected | TestReviewRecordsNoArchivedSpecForAnActiveSpec, TestDeliveryReviewResultCarriesArchivedSpecs; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/review_convention_validator_test.go | unaffected | TestReviewValidatorDismissesADaemonSettlementAsC2, TestReviewValidatorCannotDismissOutsideTheConventionRegion; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/review_override_convention_test.go | unaffected | TestOverrideArchiveIsEligibleForConventionC5, TestArchiveWithoutOverrideIsNotEligibleForC5; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/review_scope_test.go | unaffected | TestReviewOmitsQAEvidenceAndUpstreamSkills, TestReviewOmitsASkillTheLockDroppedAtTheHead; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/review_test.go | unaffected | TestReviewReadsAnArchivedSpecFromTheCandidate; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/runs_causes.go | changed with regression | TestRunCausesReadAnArchivedTaskGraphFromGit / TestRunCausesReportAnArchiveRecordWithoutItsRevision |
| internal/cli/runs_causes_test.go | unaffected | TestRunsCausesPrintsAClassifiedWindow, TestRunsCausesPrintsAnEmptyWindow; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/cli/spec_audit_record_test.go | changed with regression | TestSpecAuditAcceptsAnArchiveRecordSlug; pinned Git corpus or disposable fixtures |
| internal/cli/spec_check.go | changed with regression | TestArchiveCommandLeavesTheArchiveRecord / TestDeliveryArchiveStageCommitsTheArchiveRecord / TestInspectItemReadsTheArchiveRecord |
| internal/cli/supersede.go | changed with regression | TestSupersedeAcceptsAnArchiveRecord |
| internal/docscontract/publicdocs_test.go | unaffected | TestBaselineDocumentationContract, TestGuidanceCompositionDocumentation; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/gittest/pinned_history_test.go | changed with regression | TestPinnedHistoryMaterializesThePinnedPaths; pinned Git corpus or disposable fixtures |
| internal/judge/grouping_test.go | changed with regression | TestGroupingQuestionIsTheMeasuredOne; pinned Git corpus or disposable fixtures |
| internal/judge/questions_test.go | changed with regression | TestQuestionFileLoads; pinned Git corpus or disposable fixtures |
| internal/judge/task_acceptance_measure_test.go | unaffected | TestTaskAcceptanceQuestionFileLoads, TestTaskAcceptanceResolvesConfiguredAndArchivedRoot; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/spec/archive.go | changed with regression | TestArchiveWritesTheArchiveRecordAndRemovesTheSpecFolder |
| internal/spec/archive_layout_characterization_test.go | changed with regression | TestArchiveLayoutCharacterizationRecordsEveryRetiredFamily, TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095; pinned Git corpus or disposable fixtures |
| internal/spec/archive_links_test.go | changed with regression | TestArchiveKeepsALinkWhoseTargetWasAlreadyArchived; pinned Git corpus or disposable fixtures |
| internal/spec/archive_record_test.go | changed with regression | TestArchiveWritesTheArchiveRecordAndRemovesTheSpecFolder, TestArchiveRecordKeepsEveryOverrideField, TestArchiveRecordOfASupersededSpec, TestArchiveRecordStaysWithinTheTargetSize; pinned Git corpus or disposable fixtures |
| internal/spec/archive_test.go | changed with regression | TestArchivedPassCorpusRemainsArchiveEligible, TestArchivedQAOverrideCorpusIncludesFailedSpec; pinned Git corpus or disposable fixtures |
| internal/spec/qa_frontmatter_test.go | changed with regression | TestReadQAReportRefusesAnEmptyFrontMatter, TestReadQAReportRefusesAMissingOpeningLine; pinned Git corpus or disposable fixtures |
| internal/spec/qa_test.go | changed with regression | TestArchivedQAReportCorpusRemainsReadable; pinned Git corpus or disposable fixtures |
| internal/spec/spec_test.go | changed with regression | TestQAGateLegacyArchivedManifestsLoadUnchanged, TestRepositorySpecCorpusStillLoads, TestListActiveFiltersInactiveArchivedAndNonSpecDirectories; pinned Git corpus or disposable fixtures |
| internal/specaudit/archive_record_test.go | changed with regression | TestAuditClaimsTheArchiveRecord; pinned Git corpus or disposable fixtures |
| internal/specaudit/audit.go | changed with regression | TestAuditClaimsTheArchiveRecord |
| internal/speccheck/active_spec_paths.go | unaffected | Excludes archive/history roots when detecting active path pins; folder content is not read. |
| internal/speccheck/active_spec_paths_test.go | unaffected | TestActiveSpecPathPinIsAnError, TestActiveSpecPathPinSkipsMarkdownSpecRootsAndIgnoredFiles; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/speccheck/archive_record_test.go | changed with regression | TestArchivedSlugsIncludeArchiveRecords, TestTaskContextResolvesThroughTheArchiveRecord; pinned Git corpus or disposable fixtures |
| internal/speccheck/backlog.go | changed with regression | TestArchivedSlugsIncludeArchiveRecords / TestTaskContextResolvesThroughTheArchiveRecord |
| internal/speccheck/backlog_deferred_test.go | unaffected | TestADeferredBacklogEntryIsTerminal, TestADeferredBacklogEntryLeftActiveIsReported; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/speccheck/backlog_test.go | unaffected | TestCheckBacklogUnmoved, TestTerminalBacklogEntryLeftActiveIsUnmoved; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/speccheck/citations.go | changed with regression | TestArchivedSlugsIncludeArchiveRecords / TestTaskContextResolvesThroughTheArchiveRecord |
| internal/speccheck/citations_test.go | unaffected | TestArchivedFindingClosesWithReasonAndEvidence, TestArchivedFindingWithOnlyClosureReasonIsRefused, TestArchivedFindingWithOnlyClosureEvidenceIsRefused, TestArchivedFindingWithABlankClosureReasonIsRefused; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/speccheck/constraints_characterization_test.go | unaffected | TestConstraintReaderCharacterizesGrantCitation, TestConstraintsResolveSpecContainedRecord; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/speccheck/governed_repocontract_test.go | changed with regression | TestGovernedSetCoversOwnedShippedTemplates, TestGovernedSetOnlyGrows; pinned Git corpus or disposable fixtures |
| internal/speccheck/mechanical.go | unaffected | Finds historical authorization paths at committed ancestry before archive; no current-folder prerequisite. |
| internal/speccheck/mechanical_test.go | unaffected | TestGateAcceptsItsOwnDeclaredTerm, TestGateRefusalNamesThePreconditionThatStoppedIt; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/speccheck/receipt_characterization_test.go | changed with regression | TestReceiptCharacterizationKeepsEveryParsedClaim; pinned Git corpus or disposable fixtures |
| internal/speccheck/receipts.go | unaffected | Reads archived ADRs, not Spec folders; ADR family is untouched. |
| internal/speccheck/receipts_test.go | unaffected | TestReceiptsResolveArchivedAndInactiveRecords; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/speccheck/testdata/receipt-characterization/0181-prd.md | unaffected | Frozen source-baseline or receipt-characterization fixture; bytes remain intentionally pinned, not read from the live archive root. |
| internal/speccheck/testdata/receipt-characterization/claims-golden.json | unaffected | Frozen source-baseline or receipt-characterization fixture; bytes remain intentionally pinned, not read from the live archive root. |
| internal/speccheck/transcripts_test.go | changed with regression | TestAHeldTechSpecWithoutSurfaceTranscriptsIsAGap, TestSurfaceTranscriptsNoneWithAReasonIsAccepted; pinned Git corpus or disposable fixtures |
| internal/suiteguardcontract/archive_record_test.go | changed with regression | TestSanctionedRegenerationsReadArchiveRecords; pinned Git corpus or disposable fixtures |
| internal/suiteguardcontract/regeneration.go | changed with regression | TestSanctionedRegenerationsReadArchiveRecords |
| internal/suiteguardcontract/regeneration_test.go | unaffected | TestSanctionedRegenerationReadsArchivedSpecGrants; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/worktree/archive_reader.go | changed with regression | TestMergeEvidenceFromTheArchiveRecordAfterASquash / TestLeftoversReadTaskContextAtTheSourceRevision |
| internal/worktree/archive_record_test.go | changed with regression | TestMergeEvidenceFromTheArchiveRecordAfterASquash, TestTaskCompletionFollowsTheArchiveRecord, TestQASupersessionReadsTheArchiveRecord; pinned Git corpus or disposable fixtures |
| internal/worktree/merged_head.go | changed with regression | TestTaskCompletionFollowsTheArchiveRecord / TestQASupersessionReadsTheArchiveRecord |
| internal/worktree/merged_head_test.go | unaffected | TestMergedHeadDefaultBranchReleasesAnArchivedSpecRunAfterMainMoved; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/worktree/merged_spec_absent_target_test.go | unaffected | TestApplyReleasesAMergedRunWhoseTargetBranchIsGone, TestApplyKeepsAnUndeclaredLeftoverWhenTheTargetIsGone; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/worktree/merged_spec_leftovers_test.go | unaffected | TestArchivedSpecLeftoverScopeKinds, TestArchivedSpecLeftoversRequireTheArchivedPRD, TestArchivedSpecLeftoverRenameKeepsAnOutsideSource, TestArchivedSpecDirtyPathsKeepLiteralFilenames; existing temporary fixtures or family/layout checks, no new record-only dependency |
| internal/worktree/worktree.go | changed with regression | TestLeftoversKeepTheWorktreeWithoutTheSourceRevision |
| internal/worktree/worktree_test.go | unaffected | TestInspectTerminalRunRequiresArchivedEvidence, TestSupersedingQAReportRecognisesAnArchivedCopy, TestReconcileFallbackRequiresArchivedEvidence, TestPruneTerminalReportRequiresArchivedEvidence; existing temporary fixtures or family/layout checks, no new record-only dependency |
| skills/archive-spec/SKILL.md | changed with regression | Guidance/digest/link source validated by TestArchiveRecordClausesAreAppended, TestSettlementGuidanceIsOneTable, mirror comparisons, baseline-digests or pinned corpus gates; guidance corrections verified in this rerun. |
| skills/baseline_skill_contract_test.go | unaffected | TestBaselineSkillContract, TestNoPythonBaselineRuntime; existing temporary fixtures or family/layout checks, no new record-only dependency |
| skills/brainstorming/SKILL.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| skills/owned_skill_edit_repocontract_test.go | changed with regression | TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical; pinned Git corpus or disposable fixtures |
| skills/roundfix/references/archive.md | changed with regression | Guidance/digest/link source validated by TestArchiveRecordClausesAreAppended, TestSettlementGuidanceIsOneTable, mirror comparisons, baseline-digests or pinned corpus gates; guidance corrections verified in this rerun. |
| skills/roundfix/references/spec.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| skills/write-idea/SKILL.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| skills/write-prd/SKILL.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| skills/write-tasks/SKILL.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |
| skills/write-techspec/SKILL.md | unaffected | History-root policy or fixtures retain their meaning for records and earlier folders; inventory records this category. |

New scan paths: ['docs/specs/0242-an-archive-that-leaves-an-archive-record/task_07.md', 'internal/speccheck/archive_license_git_test.go']

| docs/specs/0242-an-archive-that-leaves-an-archive-record/task_07.md | declared authoring document | Records the corrective scope; no runtime reader. |
| internal/speccheck/archive_license_git_test.go | changed with regression | TestArchiveLicenseResolvesThroughTheRecordOrGit passes; record/Git/fallback/never-archived fixtures. |
