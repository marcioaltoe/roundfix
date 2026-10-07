# Independent acceptance evidence

The shell HTTP attempts to all three hosts returned CONNECT tunnel 403. The
web retrieval tool then reached the published sources; this is observed outside
evidence, not a substitution from Spec prose.

## Microsoft

Source: https://devblogs.microsoft.com/devops/accelerated-continuous-testing-with-test-impact-analysis-part-3/

Read the fallback, controlling TIA and trust sections. Unknown file types
trigger the complete test set, with explicit path filters available. The article
also recommends comparing selected and complete runs. This supports conservative
fallback and retaining the complete repository gate. It does not independently
prove this implementation or its 20-second budget.

## Datadog

Source: https://docs.datadoghq.com/tests/test_impact_analysis/
Independent confirmation: https://docs.datadoghq.com/tests/test_impact_analysis/how_it_works/

Read the indexed tracked-files section and selection example: dependency
manifests and Makefiles can affect tests without being code coverage inputs;
a tracked-file change runs the whole test set. This supports module-wide
fallback and explicit declarations for non-source inputs. No Datadog API ran.

## CircleCI

Source: https://circleci.com/docs/guides/test/set-up-test-impact-analysis/

Read Full test run paths and Test selection rules. The guide includes dependency
manifests and CI configuration as inputs that can select every test. Explicit
selection rules connect non-source changes to tests and allow tests to run
on every change. This supports always/relevant declarations. No CircleCI
service operation ran.
