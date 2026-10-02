# Executed commands

Worktree: /Users/marcio/.roundfix/worktrees/roundfix-deliver-0216-baseline-wording-left-after-the-stack-wave-20312d12b7b8aa4b-2dda714a-63478fa1/run_20261002T151444Z_5b33df78773c15eb

```sh
GOCACHE=/private/tmp/roundfix-0216-qa-cache rtk proxy make build
rtk proxy ./bin/roundfix --version
rtk proxy ./bin/roundfix spec check 0216-baseline-wording-left-after-the-stack-wave --strict
GOCACHE=/private/tmp/roundfix-0216-qa-cache rtk proxy go test ./internal/baseline -count=1 -v -run '^(TestTheBackendGuideScopesTheBucketProhibition|TestTheBucketClauseReplacesTheClauseAdoptersHold|TestStandardTypeScriptStructuralClauseRetention|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestReadoptionCompatibilityMaintainedFixture|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestTheBackendAndFrontendGuidesNameTheirDeclaredWorkspace|TestAGuideWithoutADeclaredWorkspaceKeepsItsGenericScope|TestEveryBuiltInWorkspaceBindsAGuideItsProfileRenders|TestTheBackendAndFrontendGuidesSayWhatTheyGovern|TestTheComposedProfileTakesTheComposedSetup|TestTheComposedProfilePlanConverges|TestFormatterComposition|TestEveryBuiltInProfileTellsTheAgentToAskForAPersonOnlySkill|TestTheSkillGuidesStateThePersonOnlySkillClause)$' > /private/tmp/roundfix-0216-qa-d6a1wga5/focused-tests-rerun.log 2>&1
rtk proxy python3 /private/tmp/roundfix-0216-qa-d6a1wga5/flows.py
rtk proxy python3 /private/tmp/roundfix-0216-qa-d6a1wga5/audit.py
```

The focused selection has no trailing space in the pattern. Full child commands, paths, separate streams and exits for CLI flows, regeneration and negative owners are in flows.log. The copied Python sources are non-compiled .py.txt evidence; copy into a disposable scratch directory as .py to replay, adjusting scratch paths. The base-build command ran in the disposable base checkout at dd0e5de9. The historical binary was used solely for base adoption setup. Every audited public CLI invocation ran ./bin/roundfix from the audited worktree.

The changed-path runner uses exported speccheck.GovernedPath from the current source copy. Its non-compiled source is governed-paths.go.txt; copied into <scratch-current>/qahelper/main.go, it ran `GOCACHE=/private/tmp/roundfix-0216-qa-cache rtk proxy go run ./qahelper` inside that copy, using changed-paths.json and granted-paths.json in the scratch parent. Exit 0.

The promise/glossary assertions counted every numbered metric and API contract, checked consuming Task References and metric Coverage Map, inspected both candidate phrases against existing glossary entries and ADR-0222, and asserted the glossary/agent guides contain no link to this Spec. Exact resulting mappings are in promise-glossary.log.

Outside reads used `git -C /Users/marcio/dev/skills show b3c45a4:<skill>/SKILL.md` (resolving the pinned tree path), and directory listings of the five named /Users/marcio/dev/<adopter>/packages paths plus read-only rev-parse HEAD. Web reads used the PRD's two official URLs on 2026-10-02; the OpenAI URL redirects to https://learn.chatgpt.com/docs/build-skills. No outside repository was written.

Repository Verification was not executed: the Daemon supplied `make verify-changed`, pass, exit 0, diagnostics removed on success.
