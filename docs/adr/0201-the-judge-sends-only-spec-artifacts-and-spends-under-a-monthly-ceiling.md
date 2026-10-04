---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-10-04T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The judge sends only Spec artifacts and spends under a monthly ceiling

ADR-0200 adds a command that sends repository text to a paid external model.
The maintainer authorized sending this repository's Spec artifacts only, with
the key read from the environment and a spending ceiling of US$5 per calendar
month. On 2026-10-01 the maintainer named OpenRouter as a recipient of the
same data: "Podemos consumir através do openrouter. Estou muito disposto a
seguir por esse caminho", and asked for a key of Roundfix's own there:
"Gosto de ter chaves separadas para determinar o real custo de cada
trabalho/projeto no openrouter". A boundary kept by convention would hold only until the first
convenient exception, so the command keeps it by construction:

- **What leaves the machine.** A request's state is built only from text the
  command read through one reader that accepts the Spec's own PRD and TechSpec
  and this repository's ADRs, as regular files inside their roots. Source code,
  diffs, QA evidence, Run state, other repositories and Secondbrain content
  cannot reach a request, because nothing else can produce the value a request
  is built from. An artifact that is not English is not sent, because the
  thresholds were measured on English text only.
- **Who receives it.** Two recipients, chosen once per run. When
  `ROUNDFIX_OPENROUTER_API_KEY` is set, every request goes to OpenRouter's
  System One API (`https://openrouter.ai/api/v1/systemone`), which forwards it
  to TypeSafe; otherwise, when `ROUNDFIX_TYPESAFE_API_KEY` is set, to TypeSafe directly
  (`https://api.typesafe.ai/v1/systemone`). No other service receives a
  request, and a run never switches recipient.
- **The keys.** Each is read from its variable in the environment and
  nowhere else, sent only in the authorization header of its own endpoint,
  and never written to output, a log or a file. The generic
  `OPENROUTER_API_KEY` is never read, so Roundfix's Jev spend stays on its own
  OpenRouter key. Without either key the command skips, names
  `ROUNDFIX_OPENROUTER_API_KEY` as the key to set, and sends nothing.
- **Every call is recorded.** Each request appends one line to a Judge Log in
  Roundfix Home, one file per UTC calendar month, with the state hash, the
  question, the recipient, the answer, its probabilities and confidence, the
  latency, the tokens, the model that answered and the cost: OpenRouter's
  reported `usage.cost` when the response carries it, or else the input
  tokens at the pinned version's price.
- **The ceiling is read from that record.** Before each request the command
  sums the month's recorded cost across every repository and both
  recipients. At or above US$5 it
  sends nothing more and reports the skip. An unreadable record also stops
  requests, because the spend below the ceiling cannot be shown.

## Consequences

- The ceiling bounds the machine's Jev spend, not one repository's or one
  account's: it sums the OpenRouter and the TypeSafe bills together.
- OpenRouter sees the same Spec artifacts TypeSafe sees, and adds a hop to
  every request.
- Two commands started at the same moment can together pass the ceiling by
  their own requests, a few cents at most.
- Findings, Backlog Entries and Task files are within the authorization but
  not read, because neither adopted judgment needs them.

**Key scope (2026-10-01).** The maintainer named Roundfix's keys for every OpenRouter and TypeSafe use, not the judge alone: "ROUNDFIX_OPENROUTER_API_KEY e ROUNDFIX_TYPESAFE_API_KEY … a key de openrouter para o roundfix vai ser para todo uso do openrouter no roundfix e do typesafe a mesma coisa e não só o jev." Any later Roundfix feature that calls OpenRouter or TypeSafe reads these two names, never the generic `OPENROUTER_API_KEY` or `TYPESAFE_API_KEY`, so each project's spend stays attributable.

**Ceiling value (2026-10-04).** Superseded in part by ADR-0231: the ceiling is the User Config value `jev.monthly_ceiling_usd`, and the US$5 above is its default when that value is unset. Every other part of this decision stands.
