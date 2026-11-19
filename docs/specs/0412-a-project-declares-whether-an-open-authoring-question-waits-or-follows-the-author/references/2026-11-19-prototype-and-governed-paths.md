# Characterization, prototype of the texts, and Governed Paths

Recorded by the author of the PRD and the TechSpec on 2026-11-19, on the Mac mini, from this repository at `894710c3` (the release of v0.148.2). It answers four questions: what the surfaces this Spec changes do today, whether the texts of the TechSpec pass the repository's contracts, which files the change and its regeneration rewrite, and which of those files are Governed Paths. The Go changes of the TechSpec were not prototyped.

## 1. What the surfaces do today

| Subject | Where | What it does |
| --- | --- | --- |
| The Clarification Record's check | `internal/speccheck/clarification_record.go`, `detectClarificationRecord` | Inside the `## Decisions` section it looks for a line that matches `^- (\d{4}-\d{2}-\d{2}): \S` with a valid date, or `^Autonomous route: no person present; decisions taken from \S.*\.$`. Either one returns clean. It reads no other line, so an entry in another form moves no result. Outside tests, no other code of `internal/` matches `## Decisions`. |
| The code set of the check | `internal/speccheck/coherence.go`, `internal/docscontract/testdata/corpus-golden.json`, `internal/spec/archive_layout_characterization_test.go` | A new code needs a row in the detector table and in both corpus pins. The three files are Governed Paths. |
| A key only a project sets | `internal/config/authoring.go`, `prepareAuthoringConfig` | `delivery.require_defect_class` and `skills` are removed from User Config before validation, each with one warning that names Project Config. |
| A key read where a branch cannot change it | `internal/cli/correction_signs.go`, `internal/config/delivery.go`, `internal/preflight/default_branch_head.go` | `roundfix deliver status` resolves the default branch head from the local refs, the remote-tracking ref first, with `preflight.ResolveDefaultBranchHead`, and reads the configuration at that commit with `roundconfig.DeliveryConfigAtCommit`, which runs `git ls-tree` and `git show` and resolves the bytes as `Load` does. Neither fetches. |
| An unknown key | `internal/config/config.go`, the strict decode | A Project Config with `authoring:` is refused today: `field authoring not found in type config.configOverlay` (run on 2026-11-19 with the binary built at the base). |
| The default configuration | `defaultConfigYAML` | It holds neither `skills` nor `delivery.require_defect_class`. |
| The status view | `internal/cli/follow_up_status.go` | It builds one document with the fields `schema` to `maintainer`, prints `Waiting on the maintainer` and then `Sources:`, lists the active Specs with `spec.ListActiveDetailed` and loads each one's Task Graph. It reads no PRD. It already receives a Git runner as a parameter. |
| Tests that pin the status text | `internal/cli/follow_up_status_queue_test.go`, `internal/cli/follow_up_status_test.go` | `TestTheStatusViewShowsTheSpecQueueTable`, `TestTheStatusViewSurfaceTranscripts` with `statusQueueSurface1` and `statusQueueSurface2`, and `TestFollowUpStatusNamesTheManualClassification` assert the lines around `Waiting on the maintainer` and `Sources:`. |
| The field lists in the guides | `internal/cli/live_documents_docscontract_test.go`, `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes` | The sentence "The status report carries these JSON fields:" in the `follow-up` guide and in both copies of the Roundfix skill's `follow-up` reference must name the JSON fields of the status document, in order. |
| The Waiting Declaration | `internal/cli/follow_up_inventory.go` | Read only on a Backlog Entry and a Finding. |

The status view of an empty repository on `main`, with a temporary `HOME` and no remote, printed these lines with the binary built at the base. The TechSpec's Surface Transcript 1 is these lines with the two new ones before `Sources:`.

```text
Status for repo in <root>
In progress: 0
Pending Specs: 0
Spec queue: 0
Next release: unavailable: resolve release range: find latest stable release tag: no stable release tag; create or pass a stable vMAJOR.MINOR.PATCH base tag
Blockers: 0
Waiting on the maintainer: 0
Sources:
  runs: checked 0
  delivery-queue: checked 0
  pull-requests: unavailable: <the forge tool's message>
  git: checked 0
  specs: checked 0
  release-plan: unavailable: <the same reason as the release line>
  backlog: checked 0
  phase-record: checked 0
  delivery-history: checked 0
```

## 2. What the prototype applied

On a scratch worktree of `894710c3`, outside the planning clone and removed afterwards, a script applied the exact texts of the TechSpec's API Contracts 4 to 14:

- the rule `rule.spec.open-authoring-question` with its two clauses in `internal/baseline/assets/modules/spec-workflow.json`, the guide's rule list and guide version, and two lines in the closed list of `internal/baseline/clause_characterization_test.go`;
- the section and the stage paragraph in the canonical `write-prd`, `write-techspec` and `write-tasks` skills, the section of the Clarification Record guide and the lines of the PRD template;
- the three texts of the Roundfix skill's `modes` reference, the paragraph of its `follow-up` reference and the bullet of its `spec` reference;
- the four glossary entries and the two revised definitions, the four user-guide texts and the notice of ADR-0280.

Three things differ from the TechSpec. The field-list sentences that follow the Go structs were not applied, because the structs did not change. The glossary gained a one-line stand-in for the definition of Rule Challenge, which the Spec delivered before this one adds. The five new test files were not written.

## 3. The commands and the files each one rewrote

`make baseline-digests` first ran while the skill mirrors still differed from their canonical files and stopped at `./skills:TestAuthorialSkillSync` with exit 2. After the mirrors were written, each command below exited 0, in this order.

| command | files it rewrote |
| --- | --- |
| `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1` | `internal/baseline/module-versions.json`, and the version line of `internal/baseline/assets/modules/spec-workflow.json` |
| `make skills-sync` | the mirrors under `skills/` of the eight canonical skill files edited |
| `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` | `skills/testdata/owned-skill-versions.json` and the version fields of both copies of the `SKILL.md` of `roundfix`, `write-prd`, `write-techspec` and `write-tasks` |
| `make baseline-digests` | `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/testdata/catalog.digest`, `internal/baseline/testdata/catalog.normalized.json`, `internal/baseline/testdata/catalog.diagnostics.golden.json`, and under `internal/baseline/testdata/plan-characterization/`: `advisory-only-divergences.golden.json`, `clean-adoption.golden.json`, `idempotent-replan-after-verified-apply.golden.json` and `same-baseline-changed-profile-and-catalog-digests.golden.json` |
| `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text` | `docs/agents/spec-routing.md` and `docs/agents/setup-context.json`; a second run printed `File changes: 0` |
| `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` | `docs/references/behavior-surfaces.json` |

The change is 40 files before the Go work and the new tests. The derived Baseline files are the same set the prototype of a new rule in the same module measured for the Spec delivered before this one. No file under `internal/baseline/assets/source-baselines/` or `internal/baseline/assets/setups/` changed.

## 4. Tests on the prototype

Each ran with `GOMAXPROCS=2`, after the commands above.

- `go test ./internal/baseline -count=1`: exit 0, with the two clause ids added to the closed list `characterizedBaselineForce`. It includes `TestNoTwoBaselineClausesShareText`, `TestShippedGuidanceCitesNoRepositoryRecord` and the tests that pin the question clause and the route clauses.
- `go test ./skills -count=1`: exit 0. No sentence an existing test pins changed, and the `--follow-up` section kept eight steps.
- `go test -count=1 -tags docscontract ./internal/docscontract`: the first run exited 1 with one failure the texts caused. `TestProfilesDocumentationContractMatchesPublicGuidance` forbids the word "recommendation" in both copies of the `write-tasks` skill, and the shared section used it twice. The two sentences were reworded to "the recommended option" and "the option you recommend yourself", and the second run exited 0. That run held this Spec's PRD and TechSpec in the active corpus, so `TestCheckCorpusGolden` and `TestCheckActiveCorpusHasNoErrors` passed with them, and the glossary contracts passed with the four entries and the two revised definitions.
- `go test ./internal/speccheck ./internal/spec -count=1`: exit 0.
- `go test ./internal/cli -run 'Baseline|GuidanceComposition|FollowUpGuides|LiveDocuments' -count=1`: exit 0.

Not run on the prototype: `make verify-changed`, and any test of the Go changes, which do not exist yet.

## 5. Governed Paths

Measured on the planning clone with `speccheck.GovernedPath`, in a `go test -overlay` probe of `internal/speccheck` that wrote nothing to the repository.

Governed, and changed by this Spec:

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`
- `.agents/skills/roundfix/references/follow-up.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/spec.md`
- `.agents/skills/write-prd/SKILL.md`, `skills/write-prd/SKILL.md`
- `.agents/skills/write-prd/references/clarification-record.md`
- `.agents/skills/write-prd/references/prd-template.md`, `skills/write-prd/references/prd-template.md`
- `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`
- `.agents/skills/write-techspec/SKILL.md`, `skills/write-techspec/SKILL.md`
- `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`
- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- `internal/baseline/assets/modules/spec-workflow.json`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`

Not governed, and changed or created by this Spec: under `internal/config/`, `config.go`, `authoring.go` and the new `authoring_questions_test.go`; under `internal/cli/`, `follow_up_status.go`, the new `follow_up_status_authoring.go` and its test, `follow_up_status_test.go`, `follow_up_status_queue_test.go` and `live_documents_docscontract_test.go`; under `internal/speccheck/`, `clarification_record.go` and the new `delegated_decision.go` and its test; `internal/baseline/module-versions.json`, `internal/baseline/clause_characterization_test.go`, the new `internal/baseline/open_authoring_question_clauses_test.go`, the three catalog files and the four plan goldens; the mirrors `skills/roundfix/references/follow-up.md`, `skills/roundfix/references/modes.md`, `skills/roundfix/references/spec.md` and `skills/write-prd/references/clarification-record.md`; `skills/testdata/owned-skill-versions.json` and the new `skills/open_authoring_question_test.go`; `CONTEXT.md`, the four user-guide pages, `docs/references/behavior-surfaces.json`, ADR-0442 and ADR-0280.

Controls, governed and not changed: `.roundfixrc.yml`, the `Makefile`, `go.mod`, `.github/workflows/ci-verify.yml`, `docs/agents/agent-instructions.md`, `docs/agents/autonomous-work.md`, `internal/speccheck/coherence.go`, `internal/docscontract/testdata/corpus-golden.json`, `internal/spec/archive_layout_characterization_test.go`, `docs/references/coverage-record.json`, `.agents/skills/write-idea/SKILL.md` and both copies of the TechSpec template. Controls measured not governed and not changed: `internal/config/delivery.go`, `internal/cli/follow_up.go`, `internal/docscontract/corpus_test.go`, `docs/references/skill-coverage.json` and `AGENTS.md`.
