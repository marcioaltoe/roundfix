---
status: accepted
created_at: 2026-09-29T00:00:00Z
updated_at: 2026-09-29T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A pre-PR review reads the final message and keeps a record per checkout

The Pre-PR Review Command classified the concatenation of every Agent message chunk in the review session. A
progress message and the final answer therefore reached the classifier joined
without a boundary. On 2026-09-29 a well-formed five-finding answer began with
`…regression, or security issues.Findings:`, and the review was recorded
`blocked: unclassifiable agent output`. The same concatenation fed the sealed
prompts whose output the Baseline parses as JSON.

The review record and its answer file also lived at the repository's Artifact
Directory. Since the durable repository key, every checkout of a repository
shares that directory: the main checkout, linked worktrees and every Delivery
Queue item worktree. A review run in one checkout replaced the record and the
answer of another. The replaced record could no longer be reused or disposed
from its own checkout, and its answer evidence was lost.

Now the Agent runner keeps message boundaries, as the Agent Client Protocol
defines them. When two message chunks both carry a `messageId`, the identifiers
decide. Without one, a message chunk that follows a thought, a tool call or a
plan starts a new message. Codex's non-interactive surface likewise writes only
the agent's last message. The reviewer's verdict is read from the last non-blank
message, and the answer file keeps every message separated by a blank line. A
sealed prompt's output is its last non-blank message. Each checkout keeps its
own review record and answer under the Artifact Directory, in a directory named
from the checkout path. The disposition ledger stays shared, because every
disposition already names its repository and reviewed head.

## Consequences

The classifier's strictness is unchanged inside the final message. Both
verdicts together, or a no-findings verdict beside other content, still block.
A progress message is no longer part of the verdict, so a review that ends in
`No findings.` after commentary passes. Two chunks with no intervening update
and no `messageId` remain one message, as the protocol leaves them. A record
written at the old shared location is neither read nor removed, so the first
review in each checkout after the upgrade runs fresh. The delivery engine keeps
reading the record from the command's own output, and a reader never takes a
record from another checkout.
