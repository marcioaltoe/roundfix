<!-- source-baseline-entry: clause.external-triage.classify-before-labelling -->
- MUST use this workflow only for items in an external forge, and classify each item as a bug or an enhancement, in English, before changing its labels or status.
<!-- /source-baseline-entry: clause.external-triage.classify-before-labelling -->

<!-- source-baseline-entry: clause.external-triage.move-through-mapped-states -->
- MUST move each item through the states needs-triage, needs-info, ready, and wontfix with the forge labels the repository maps to them.
<!-- /source-baseline-entry: clause.external-triage.move-through-mapped-states -->

<!-- source-baseline-entry: clause.external-triage.ask-for-an-unmapped-label -->
- MUST stop and ask the maintainer for the forge label of an unmapped triage state; never invent a label.
<!-- /source-baseline-entry: clause.external-triage.ask-for-an-unmapped-label -->

<!-- source-baseline-entry: clause.external-triage.needs-info-asks-and-stops -->
- MUST mark an item needs-info, ask the reporter specific questions, and stop until the reporter answers.
<!-- /source-baseline-entry: clause.external-triage.needs-info-asks-and-stops -->

<!-- source-baseline-entry: clause.external-triage.disclose-ai-authorship -->
- MUST start every forge comment with a sentence stating that an AI agent generated it.
<!-- /source-baseline-entry: clause.external-triage.disclose-ai-authorship -->

<!-- source-baseline-entry: clause.external-triage.record-wontfix-decisions -->
- MUST record each wontfix decision as a declined Backlog Entry and answer a repeated request with it.
<!-- /source-baseline-entry: clause.external-triage.record-wontfix-decisions -->

<!-- source-baseline-entry: clause.external-triage.pull-request-is-an-issue-with-code -->
- MUST triage an external pull request as an issue with code attached, through the same states.
<!-- /source-baseline-entry: clause.external-triage.pull-request-is-an-issue-with-code -->

<!-- source-baseline-entry: clause.external-triage.route-accepted-work-to-specs -->
- MUST route accepted work into the local Spec workflow; a forge label or status never stands in for a Task's status.
<!-- /source-baseline-entry: clause.external-triage.route-accepted-work-to-specs -->

<!-- source-baseline-entry: rule.external-triage -->
Use this workflow only for issues managed in an external forge. Classify the user-visible problem and next action in English before changing labels or status, and route implementation through the repository's local Spec policy.
<!-- /source-baseline-entry: rule.external-triage -->
