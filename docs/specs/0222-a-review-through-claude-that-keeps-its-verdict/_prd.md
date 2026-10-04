---
spec: 0222-a-review-through-claude-that-keeps-its-verdict
status: active
created: 2026-10-04
surfaces: [backend, cli, docs]
---

# A review through Claude that keeps its verdict

On 2026-10-02 two repositories that select `pre_pr_review.provider: claude`
with the `review` profile `claude / opus / high` could not complete the pre-PR
review the repository policy requires
([the adopted Backlog Entry](references/2026-10-03-the-review-through-claude-blocks-on-a-refused-permission.md)).
`roundfix review` blocked as a `review transport anomaly` because acpx exited
`5` after the turn had already delivered its verdict: the reviewer had asked to
run a terminal command and the read-only session refused it. The block hid a
clean `No findings.` in one repository and seven real findings in four rounds
in the other, which got no finding ids and so could not be disposed. In the
second repository a large candidate also overflowed the 1,000,000-token Claude
context, and the command printed `agent/protocol error` while the real reason,
`Prompt is too long`, stayed in the answer file. Spec 0212's bound counts the
diff in bytes and ignores the Spec context the command adds and the window of
the provider it sends to.

This is a bug fix: the review keeps its contract and its policy. It mints this
minimal PRD for the downstream artifact contract; the design lives in the
[_techspec.md](_techspec.md).

## Prerequisites

None. Specs 0223 and 0224 are authored in the same cycle; if either also
raises the Roundfix Skill's version, the operator orders the queue so that the
later Spec raises it from the earlier one's value.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; finding
  ids keep the existing `F<n>` sequence and the review record gains two
  optional fields with the existing camel-case JSON names. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the review still drives the
  locally installed acpx; no request, credential or forge read is added, and
  no test reaches a provider or the network. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0227 (this Spec) decides that a
  read-only review turn that ended after a refused permission keeps its
  verdict and that the Claude review prompt is bounded in tokens. ADR-0020:
  "a subsequent nonzero acpx exit is classified as teardown noise", the
  precedent this Spec applies to the review for the permission exit only.
  ADR-0174: "A pre-PR review reads the final message", so the verdict and the
  prompt-too-long line are read from the final answer. ADR-0196's anchor
  validation and ADR-0197's two-round lineage apply unchanged to the findings
  this Spec keeps. ADR-0153: "Pre-PR review is an explicit provider policy", which this Spec
  keeps. ADR-0151: "Configured review profiles select the independent
  reviewer", and the profile stays as it is.
  ADR-0169: "The pre-PR review diffs the candidate from its merge base", and
  that diff is what the bound measures. ADR-0187 splits the Roundfix Skill by command and ADR-0189 ties an
  owned skill's version to its content, so the skill edit raises the version.
  ADR-0184: "A TechSpec now declares numbered Surface Transcripts", applied to
  the Claude bound's refusal. The gate is bound by ADR-0080, ADR-0088,
  ADR-0091, ADR-0104, ADR-0155, ADR-0156 and ADR-0167, and ADR-0093, ADR-0117,
  ADR-0168, ADR-0176 and ADR-0183 check this Spec's consistency by citation and
  receipt. ADR-0178 authorizes each Task commit by its grant, ADR-0182 runs
  Settlement Checks before it, and ADR-0166 records undeclared paths; every
  Task declares its paths. ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry, ADR-0165 cites ADR-0153 but decides how a blocking review after archive parks publication, and ADR-0192 cites ADR-0178 but decides how a conflict confined to declared derived paths is resolved; ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when it is observed again and its evidence snapshot; this Spec changes none of them, so none applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and its
  `SKILL.md` mirror are Governed Paths, and the maintainer authorized skill
  edits ("considere autorizado a ajustar todas as skills se necessário") and,
  on 2026-10-04, this Spec. No other Governed Path changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0222-a-review-through-claude-that-keeps-its-verdict/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A review through Claude whose reviewer tried a refused command is classified
  by its verdict, and its findings can be disposed.
- An operator whose review overflowed the provider's context reads that reason
  in the command's output, not a protocol error.
- The Claude review prompt is bounded by the provider's context window in
  tokens, counting everything the command sends.
- A review through Codex, and every read-write Agent Session, behaves as
  before.

## User Stories

1. As an operator reviewing through Claude, I want a verdict the reviewer
   delivered after the session refused a command to count, so that a clean
   review lets me publish and a review with findings gives me ids to dispose.
2. As an operator, I want the review record to say the session refused a
   permission, so that I can tell the reviewer worked without a tool it asked
   for.
3. As an operator whose candidate is too large for the reviewer, I want the
   command to say `Prompt is too long` with the runtime's numbers, so that I
   know to split the candidate instead of retrying.
4. As an operator reviewing through Claude, I want an oversized prompt refused
   before any provider call, with its estimated tokens and the budget, so that
   I do not spend a review on a request that cannot fit.

## Core Features

1. **A refused permission does not discard a delivered turn.** When a
   read-only Agent Session's turn ends with `end_turn` and acpx exits `5`, the
   runner keeps the result, marks that the session refused a permission and
   publishes the permission-denied status. Every other exit after a parsed
   result remains a transport anomaly, and read-write sessions are unchanged
   (ADR-0227).
2. **The review classifies that answer by its verdict.** `roundfix review`
   classifies the answer as it classifies any other, so a findings answer gets
   finding ids and the record carries `permissionRefused: true`; a block for
   any other reason names the refusal.
3. **The overflow reason reaches the command.** When the final answer carries
   a line starting `Prompt is too long`, the review blocks with the reason
   `review prompt too long: ` followed by that line, whatever the exit.
4. **A Claude prompt bound in tokens.** For the `claude` provider the command
   estimates the tokens of the whole prompt, at two bytes per token, and
   refuses before readiness and any provider call when the estimate exceeds
   500,000 tokens, half of the 1,000,000-token window; the record carries the
   estimate. The `codex` provider keeps its 917,504-byte diff bound.
5. **The skill and the guide say so.** The Roundfix Skill's `review`
   reference and the `review` command guide describe the refused-permission
   classification, the record field, the prompt-too-long reason and the
   Claude bound.

## Non-Goals / Out of Scope

- Denying the terminal tool up front: acpx 0.19.4 cannot pass the Claude
  adapter's tool restrictions (ADR-0227). A later acpx that can is a follow-up.
- Changing the review prompt's instructions, the verdict grammar, the finding
  validation or the two-round lineage.
- Changing the Codex bound, the Codex review path, or any read-write session's
  handling of an acpx exit.
- Predicting the files a `@path` mention makes Claude Code attach; the reserve
  in the bound covers them, and the prompt-too-long reason covers what remains.
- Splitting an oversized candidate into several reviews.

## Success Metrics

1. Success Metric: a fake acpx that prints a parsed `end_turn` result and
   exits `5` in a read-only session yields a result with no transport anomaly
   and the refusal marked; with exit `1`, or with a stop reason other than
   `end_turn`, the anomaly stays.
2. Success Metric: a review whose runner reports a refused permission and a
   findings answer exits `1` with findings `F1..Fn`, and `roundfix review
   dispose` accepts `F1`; a `No findings.` answer exits `0`.
3. Success Metric: a review whose answer starts with `Prompt is too long`
   blocks with that line in its reason, both after a runtime failure and after
   a parsed result.
4. Success Metric: a Claude review whose prompt estimate exceeds 500,000
   tokens exits `2` before any runner call; a prompt within the budget whose
   diff exceeds 917,504 bytes reaches the runner; a Codex review with that
   diff is refused by its byte bound as before.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The fiscus and oraculum reports of 2026-10-02, now in the Secondbrain at
  `inbox/roundfix/_triaged/2026-10-02-review-claude-bloqueia-quando-o-revisor-pede-terminal.md`
  and
  `inbox/roundfix/_triaged/2026-10-02-review-claude-perde-achados-e-estoura-o-prompt-no-oraculum.md`,
  including the runtime's own message `Prompt is too long · the request is
  ~1294791 tokens (limit 1000000)`.
- Three real reviews through `claude / opus / high` during authoring on
  2026-10-04, in a disposable clone of this repository with the 0.33.0-era
  binary at `6ea9e1e9`: one where the reviewer ran `sed` and `grep` through
  its terminal without any permission request and the review succeeded; one
  where it ran `go vet`, the session answered `session/request_permission`
  with `reject`, the turn ended `end_turn` with one finding, and the command
  blocked with exit code `5` as a transport anomaly; and one whose diff
  `@path`-mentioned 40 files of 3,387,660 bytes, whose first context reading
  was 82,632 tokens against about 25,000 for the other two.
- The installed acpx 0.19.4, whose `applyPermissionExitCode` sets
  `PERMISSION_DENIED` (`5`) when a turn's permission requests were all denied
  or cancelled and none approved, and the installed Claude adapter 0.85.0,
  which reported `size: 1000000` for `claude-opus-5-5`.
- Anthropic's Agent SDK permissions guide
  (<https://code.claude.com/docs/en/agent-sdk/permissions>): `allowedTools`
  pre-approves the listed tools, and "Tools not listed here still exist";
  only `disallowedTools` removes a tool. Anthropic's token-counting guide
  (<https://platform.claude.com/docs/en/build-with-claude/token-counting>):
  models from Opus 4.7 use a tokenizer that produces about 30 percent more
  tokens for the same text.

## Decisions

- The permission exit is classified only for a read-only session whose turn
  ended `end_turn`; see ADR-0227.
- The Claude bound is half the window in estimated tokens over the whole
  prompt; Codex keeps its measured byte bound; see ADR-0227.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
