<!-- setup-context-driven:begin id=guide.external-triage version=0.0.1 -->

# External triage

- Use this workflow only for issues managed in an external forge. Classify the user-visible problem and next action in English before changing labels or status, and route implementation through the repository's local Spec policy.

<!-- setup-context-driven:end id=guide.external-triage -->

<!-- roundfix:repository-rule:begin id=rule.164b03244b3e089313c3cd9b379f4cce885e87b985e19ef76fb2ee8bff5ec7e6 -->
### Triage labels

These labels apply only to external GitHub issues. Internal Spec Tasks use
`pending | in_progress | completed | failed` in task-file frontmatter.

| Canonical role | GitHub label | Meaning |
| --- | --- | --- |
| `needs-triage` | `needs-triage` | A maintainer must evaluate the issue |
| `needs-info` | `needs-info` | Waiting for information from the reporter |
| `ready-for-agent` | `ready-for-agent` | Fully specified and ready for an autonomous Agent |
| `ready-for-human` | `ready-for-human` | Requires human implementation |
| `wontfix` | `wontfix` | Will not be actioned |

Always use the canonical label string. Do not substitute generic GitHub labels
such as `question`, `invalid`, or `help wanted` for these workflow roles.

<!-- roundfix:repository-rule:end id=rule.164b03244b3e089313c3cd9b379f4cce885e87b985e19ef76fb2ee8bff5ec7e6 -->
