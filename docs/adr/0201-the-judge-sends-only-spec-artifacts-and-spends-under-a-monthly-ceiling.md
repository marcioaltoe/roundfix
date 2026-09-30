---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The judge sends only Spec artifacts and spends under a monthly ceiling

ADR-0200 adds a command that sends repository text to a paid external model.
The maintainer authorized sending this repository's Spec artifacts only, with
the key read from the environment and a spending ceiling of US$5 per calendar
month. A boundary kept by convention would hold only until the first
convenient exception, so the command keeps it by construction:

- **What leaves the machine.** A request's state is built only from text the
  command read through one reader that accepts the Spec's own PRD and TechSpec
  and this repository's ADRs, as regular files inside their roots. Source code,
  diffs, QA evidence, Run state, other repositories and Secondbrain content
  cannot reach a request, because nothing else can produce the value a request
  is built from. An artifact that is not English is not sent, because the
  thresholds were measured on English text only.
- **The key.** It is read from `TYPESAFE_API_KEY` in the environment and
  nowhere else, sent only in the request's authorization header, and never
  written to output, a log or a file. Without it the command skips and sends
  nothing.
- **Every call is recorded.** Each request appends one line to a Judge Log in
  Roundfix Home, one file per UTC calendar month, with the state hash, the
  question, the answer, its probabilities and confidence, the latency, the
  tokens, the model that answered and the cost at the pinned version's price.
- **The ceiling is read from that record.** Before each request the command
  sums the month's recorded cost across every repository. At or above US$5 it
  sends nothing more and reports the skip. An unreadable record also stops
  requests, because the spend below the ceiling cannot be shown.

## Consequences

- The ceiling bounds the machine's spend, not one repository's, since the key
  and the bill belong to one account.
- Two commands started at the same moment can together pass the ceiling by
  their own requests, a few cents at most.
- Findings, Backlog Entries and Task files are within the authorization but
  not read, because neither adopted judgment needs them.
