# Review-request coherence table

Build: `bdf6ff8d4d680188a97986ee1340ab9ff052a499`.

The built CLI ran from four clean scratch Git repositories under
`/private/tmp/roundfix-qa-0078.b1jXZM/`. Each invocation used an explicit
writable Artifact Directory and `--no-input`.

| pushTriggersReview | request_review | Public `resolve` observation | Public `watch` observation |
| --- | --- | --- | --- |
| true | false | Passed coherence and stopped later because the scratch branch has no push upstream. | Covered by the same Preflight predicate and the fresh table test. |
| false | true | Passed coherence and stopped later because the scratch branch has no push upstream. | The repository's own built `watch` also passed coherence and stopped later at the Run Worktree's absent upstream. |
| false | false | Exit 2; named `.coderabbit.yaml`, both resolved booleans, `pushTriggersReview=false`, the stranded Run, and `set review_source.request_review to true in Project Config`; declared no side effects. | Exit 2 with the same row-specific evidence and no-side-effects guarantee. |
| true | true | Exit 2; named `.coderabbit.yaml`, both resolved booleans, `pushTriggersReview=true`, duplicate review after every push, and `set review_source.request_review to false in Project Config`; declared no side effects. | Exit 2 with the same row-specific evidence and no-side-effects guarantee. |

Fresh independent command:

`rtk go test ./internal/preflight ./internal/cli -count=1 -run
'TestRunEnforcesReviewRequestCoherence|TestRunExemptsFetchFromReviewRequestCoherence|TestRunResolveReviewRequestCoherenceRefusalExitsTwoBeforeRunCreation'
-v`

Result: exit 0, 13 tests passed across 2 packages. This exercises all four
predicate rows, the `fetch` exemption, both operational commands' refusal
paths, public exit 2, and absence of Run Database creation.
