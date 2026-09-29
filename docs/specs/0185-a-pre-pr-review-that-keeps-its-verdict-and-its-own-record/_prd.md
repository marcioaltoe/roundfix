---
spec: 0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record
status: active
created: 2026-09-29
surfaces: [backend, cli, docs]
---

# A pre-PR review that keeps its verdict and its own record

`roundfix review` (Specs 0153, 0160 and 0179) runs the configured reviewer
over the candidate and records the result. On 2026-09-29 it lost a real verdict
and let one checkout's review replace another's.

- **A verdict after progress text is unclassifiable.** In
  `~/dev/roundfix-wave5` at head `3e77ed70` the review ended `blocked`, exit 2,
  with `unclassifiable agent output: neither a no-findings nor findings verdict
  is present`. The answer file held five well-formed findings. The first line
  read `…looking for concrete correctness, regression, or security
  issues.Findings:`. The Agent runner appends every `agent_message_chunk` of
  the session to one string with no boundary. A progress message and the final
  answer therefore reach the classifier as one line, and no line starts with
  `Findings:`. Spec 0160 already accepts a sentence before the verdict, but
  only on its own line. The same concatenation feeds the sealed prompts whose
  output the Baseline parses as JSON. A full review round was spent, and the
  findings had to be read from the answer file by hand.
- **One record is shared by every checkout.** The record
  `pre-pr-review.json` and the answer `pre-pr-review-answer.txt` live at the
  repository's Artifact Directory. Since Spec 0162 every checkout of a
  repository shares it: the main checkout, linked worktrees and each Delivery
  Queue item worktree. On 2026-09-29 reviews from `~/dev/roundfix` and
  `~/dev/roundfix-wave5` replaced each other's record. The record and answer
  for head `3e77ed70` are gone; the file now names `roundfix-wave5` at
  `da3c2eb7`. A replaced record can no longer be reused, so its findings are
  reviewed again. It can no longer be disposed from its own checkout, which
  refuses a record that belongs to another repository. Its answer evidence is
  lost. The Delivery Queue reads its verdict from the review command's own
  standard output, so the verdict it acts on is not replaced. Its reuse, its
  dispositions from the item worktree and its answer evidence are exposed in
  the same way.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. A checkout's record
  directory is named by a digest of its path, which is not an identity any
  other surface reads. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the local
  Agent runtime only; no credential is read and no network call is added.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0174 (this Spec) reads the review
  verdict from the reviewer's final message and keeps a record per checkout,
  and is recorded before implementation. ADR-0153 keeps the Pre-PR Review
  Policy explicit and ADR-0151 keeps the configured review profile; neither
  changes. ADR-0165 parks a blocking review after archive from the record the
  command prints, which is unchanged. ADR-0169 (Spec 0182) makes the merge
  base the record's `baseCommit` and reuse key, and this Spec keeps that key.
  ADR-0017 drives ACP Runtimes through acpx, whose stream the runner now reads
  with boundaries. ADR-0020 lets a parsed prompt result outrank the acpx exit
  code, and a final message changes none of that. This Spec's gate is bound by
  ADR-0080, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155 and ADR-0156.
  ADR-0093 checks Spec consistency by citation, and ADR-0094 makes that check
  artifact-presence-aware. ADR-0166 (Spec 0181) has the Daemon record the
  paths a Task changed without declaring them; every Task here declares the
  paths it edits, so it changes nothing here. ADR-0167 (Spec 0181) keeps the
  pre-PR Pull Request row from deciding a qualifying partial; this Spec's gate
  aims at `pass`, so it does not apply. ADR-0168 (Spec 0181) opens a
  related-ADR gap only for ADRs that predate the Spec, which narrows this
  Spec's own check and holds. ADR-0170 (Spec 0182) governs Task
  Carry-Forward, and ADR-0171, ADR-0172 and ADR-0173 (Specs 0183 and 0184)
  govern retention, Doctor storage and Baseline relocation citations; this
  Spec touches none of them. ADR-0097 cites ADR-0080 but carries a QA row
  forward, not a review record, so it does not apply. All the others hold. ADR-0176 (Spec 0181) narrows only which Spec text the citation checks read, and this Spec does not depend on it, so it does not apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer approved adding this Spec to
  the second queue on 2026-09-29, answering a structured question with "Nova
  Spec na Onda 5"; the Roundfix skill files ride the standing grant of
  2026-09-18 for keeping the shipped skills true to the CLI, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A reviewer's verdict given in its final message is classified, whatever
  progress messages came before it in the session.
- The classifier's strictness inside that final message does not change.
- A sealed prompt's output is its final message.
- Each checkout keeps its own pre-PR review record and answer, and a review in
  one checkout never replaces another checkout's.
- Reuse and dispositions keep working in the checkout that ran the review.

## Core Features

1. **Agent messages keep their boundaries.** When two message chunks both
   carry a `messageId`, a different one starts a new message and the same one
   continues it, even across other updates. Without one, the Agent runner
   starts a new message when a message chunk follows a thought, a tool call or
   a plan. It reports the messages in order, and the joined text separates them
   with a blank line.
2. **The verdict is read from the final message.** The Pre-PR Review Command
   classifies the last non-blank message with today's rules. The answer file
   keeps every message. A progress message is no longer part of the verdict.
   Both verdicts together, or a no-findings verdict beside other content in the
   final message, still block.
3. **A sealed prompt parses its final message.** The sealed prompt's output is
   its last non-blank message. The output cap still counts every message, and a
   tool call is still refused.
4. **A review record per checkout.** The record and the answer live in a
   directory of the Artifact Directory named from the checkout's path. The
   review, its reuse and `roundfix review dispose` read only the current
   checkout's directory. The disposition ledger stays shared, because each
   disposition names its repository and reviewed head.

## Non-Goals / Out of Scope

- Loosening the classifier into substring matching, or accepting a verdict
  anywhere but at the start of a line in the final message.
- Changing the reviewer prompt, the Pre-PR Review Policy, the profiles, the
  merge-base contract of Spec 0182, or how the Delivery Queue reads a verdict.
- Migrating, reading or deleting a record written at the old shared location.
- Pruning record directories of removed checkouts; each holds two small files.
- Changing the disposition ledger's location or format.

## Success Metrics

1. A fake ACP stream with a progress message, a tool call, then
   `Findings:` and one list item reports two messages. The review classifies
   `findings` with one finding, and the answer file keeps both messages.
2. The same stream with distinct `messageId` values and no intervening update
   also reports two messages. Chunks sharing one `messageId` across a tool call
   stay one message, and two chunks with neither boundary stay one message.
3. A final message `No findings.` after a progress message classifies
   `reviewed`. A final message carrying both verdicts still blocks.
4. A sealed stream with a commentary message, a thought, then a JSON message
   yields exactly the JSON.
5. Two checkouts of one repository, each running `roundfix review` with a fake
   reviewer, keep two records. Each checkout reuses and disposes its own, and
   neither run changes the other's record or answer.

## Recorded limits

- Two message chunks with no intervening update and no `messageId` remain one
  message, as the Agent Client Protocol leaves them. codex-acp 2.0.0 sends a
  `messageId` on every agent message chunk.
- A record at the old shared location is ignored. The first review in each
  checkout after the upgrade runs fresh instead of reusing it.
- Record directories of removed checkouts stay until someone deletes them.

## Decisions

- **Fix the stream, not the classifier.** The boundary exists in the protocol,
  and the runner dropped it. A classifier that searches for `Findings:` inside
  a line would accept verdict-shaped prose. See ADR-0174.
- **The final message is the answer.** The reviewer's commentary is not its
  verdict, as Codex's own non-interactive surface treats the last message.
  See ADR-0174.
- **Per checkout, not per head.** `roundfix review dispose --fixed-by` runs at a
  later head than the one reviewed, so the reader locates the record by
  checkout and still checks its head. See ADR-0174.
- **Ignore the old location.** Reading it back would reintroduce the shared
  record for one more review, and a fresh review costs one run.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in
the Task Graph. The negative cases carry the weight:

- chunks without a boundary that are split anyway;
- a final message with both verdicts that passes;
- a sealed prompt that parses commentary;
- a review in one checkout that changes another checkout's record.

A happy-path test would miss each of these.

The outside-evidence rows rest on sources this Spec did not produce:

- The Agent Client Protocol's message-id RFD. It defines `messageId` on
  `agent_message_chunk` as stable across a message's chunks. It shows that
  consecutive chunks without one are ambiguous, while a tool call between them
  "definitely" starts a new message
  (<https://github.com/agentclientprotocol/agent-client-protocol/blob/main/docs/rfds/message-id.mdx>).
- The installed codex-acp 2.0.0 adapter. Its
  `/opt/homebrew/lib/node_modules/@agentclientprotocol/codex-acp/dist/index.js`
  builds every agent message chunk with `createAgentTextMessageChunk(item.text,
  item.id, …)`, so each message carries its own `messageId`.
- Codex's non-interactive documentation. It says `--output-last-message`
  "writes the final message to the file"
  (<https://learn.chatgpt.com/docs/non-interactive-mode.md>).
- The live shared record at
  `~/.roundfix/artifacts/339f8dac2b687a04/pre-pr-review.json`, written by other
  sessions on 2026-09-29. It now names `roundfix-wave5` at `da3c2eb7`, while
  the dispositions ledger beside it still carries dispositions for `fe517d2b`
  and `ccd5c097` from `~/dev/roundfix`.

## Research basis

The two adopted Backlog Entries are indexed in
[references/_index.md](references/_index.md). The Secondbrain was consulted
through `wiki/index.md` and
`qmd query "roundfix pre-PR review verdict classification agent message chunks"`
and
`qmd query "roundfix artifact root shared by worktrees repository key review record"`.
It returned this repository's own mirrors: Specs 0160, 0162 and 0179, and the
commands guide. It also returned the concept page
`wiki/concepts/revisao-de-codigo-em-workflows-agenticos.md`, which records
Codex's `--output-last-message` as part of its non-interactive review surface.
That supports reading the final message. Exa located the Codex documentation
above. Context7 supplied the ACP message-id RFD, which settles the boundary
rule of ADR-0174.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
