# Triage labels

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
