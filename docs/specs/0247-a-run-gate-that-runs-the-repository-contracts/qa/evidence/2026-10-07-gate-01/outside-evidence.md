# Independent acceptance evidence

All sources freshly observed for this rerun; none credited from the previous report.

## Microsoft

Source: [Microsoft](https://devblogs.microsoft.com/devops/accelerated-continuous-testing-with-test-impact-analysis-part-3/)

Fresh web retrieval read the fallback, Controlling TIA and Trusting TIA sections. Unknown file types cause a full-test fallback; path filters can constrain non-source impacts. Comparing selected and full runs validates selection. These published observations support conservative fallback and preserving the complete gate; they do not prove this implementation or its time budget.

## Datadog

Source: [Datadog](https://docs.datadoghq.com/tests/test_impact_analysis/)

Direct web open failed with unsupported text/markdown; sandbox curl-cffi was denied network access. A fresh domain-restricted web search returned the published page (with a tracking query parameter) and independently its Getting Started and How It Works pages. Read the tracked-files section: Makefiles and dependency manifests can impact tests; changes in tracked files select all tests. The selection example independently confirms this behavior. This supports module-wide fallback and explicit non-source inputs. No service API was invoked.

## CircleCI

Source: [CircleCI](https://circleci.com/docs/guides/test/set-up-test-impact-analysis/)

Fresh web retrieval read Full test run paths and Test selection rules. Dependency manifests and CI configuration can select every test; explicit rules cover non-source files or always select specific tests. These published observations support always/relevant directives and module-wide fallback. No CircleCI operation was invoked.

Datadog independent confirmation: https://docs.datadoghq.com/tests/test_impact_analysis/how_it_works/ and https://docs.datadoghq.com/getting_started/test_impact_analysis/.
