<!-- setup-context-driven:begin id=guide.external-triage version=0.0.1 -->

# External triage

- **stop-and-ask**: Stop and ask the maintainer for the forge label of a triage state the repository has not mapped; never invent a label.

- **mandatory**: Use this workflow only for issues and pull requests managed in an external forge. Classify each item as a bug or an enhancement, and state its user-visible problem and next action in English, before changing its labels or status.

- **mandatory**: Start every comment posted on the forge with a sentence stating that an AI agent generated it.

- **mandatory**: Move each item through the triage states needs-triage, needs-info, ready, and wontfix, applying the forge label the repository maps to each state. The repository records that mapping in this guide, outside its setup markers, when it enables external triage.

- **mandatory**: When an item lacks what triage needs, mark it needs-info, ask the reporter specific questions that name each missing fact, and stop triaging it until the reporter answers.

- **mandatory**: Triage an external pull request as an issue with code attached: classify it and move it through the same states before reviewing or merging its code.

- **mandatory**: Record each wontfix decision as a declined Backlog Entry whose reason cites the forge item, and answer a repeated request with that recorded decision instead of deciding it again.

- **mandatory**: Route accepted work into the repository's local Spec workflow; a forge label or status never stands in for a Task's status.

<!-- setup-context-driven:end id=guide.external-triage -->
