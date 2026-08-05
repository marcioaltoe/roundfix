# Repository configuration coherence

The committed Project Config has `review_source.request_review: true` and
`watch.max_rounds: 2`. The committed Review Source configuration has
`reviews.auto_review.enabled: false` and the explicit non-empty
`description_keyword: "coderabbit:review"`.

The built public command:

`rtk bin/roundfix watch --source coderabbit --pr 123 --head-repo
marcioaltoe/roundfix --head-branch ma/0078-roundfix-asks-for-the-review
--artifact-dir /private/tmp/roundfix-qa-0078.b1jXZM/repository-artifacts
--until-clean --no-input`

passed review-request coherence and stopped at the next Preflight check because
the retained per-Run branch has no push upstream. It reported that it created
no Run, fetched no Review Source issue, started no Agent, committed nothing,
and pushed nothing. Pair `false/true` therefore reaches the post-coherence
boundary through the real CLI.

`rtk bin/roundfix doctor` loaded the effective configuration but could not
complete unrelated machine checks: the Claude adapter lineage was unproven and
the QA sandbox denied Codex state initialization under `~/.codex`. This is an
environment parity deviation, not evidence against the configuration pair;
the feature-specific public Preflight above is the criterion's observable.
