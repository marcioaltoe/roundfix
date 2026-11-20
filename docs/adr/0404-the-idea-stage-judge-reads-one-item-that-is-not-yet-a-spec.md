---
status: accepted
knowledge: true
created_at: 2026-11-14T15:00:00Z
updated_at: 2026-11-14T16:00:00Z
deprecated_at: null
superseded_by: null
---

# The Idea-stage judge reads one item that is not yet a Spec

The Idea-stage judgment of
[ADR-0403](0403-a-feature-is-judged-for-the-idea-stage-before-its-prd-and-the-judgment-never-refuses.md)
runs before any PRD exists, so it cannot read a Spec. It reads one item the
Follow-up Command already lists, or one request of the session handed in as a
file, sends the state the follow-up judge already sends, asks three measured
questions on the pinned model, and keeps the wide undecided band, because the
measurement was too small to fix a threshold.

## Context

Every other judgment of `roundfix spec judge` needs a PRD.
[ADR-0201](0201-the-judge-sends-only-spec-artifacts-and-spends-under-a-monthly-ceiling.md)
keeps the data boundary by construction;
[ADR-0209](0209-the-judge-suggests-which-sources-share-a-spec.md) let Findings
and Backlog Entries reach a request; and
[ADR-0319](0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md)
let the judge read a Follow-up Item, "one item's kind, type, origin, title and
an excerpt of its text", whatever its language, after a secret scan.

The maintainer's first boundary for this judgment, set with the Idea on
2026-11-14, was only what the follow-up judge already sends. A request the
user types in the session has no file, so under that boundary the agent would
answer the criteria alone for it, and that is how the request behind this
decision arrived. Asked the same day whether the judge may receive such a
request, the maintainer answered "Pode, com varredura (Recommended)": "O
agente grava o pedido num arquivo entregue como uma mensagem; vai ao juiz só o
nome do projeto, a primeira linha como título e os primeiros 4.000 caracteres,
depois da mesma varredura de segredos. Os pedidos que nascem na sessão passam
a ser julgados como os outros."

On 2026-11-14 the questions were measured on this repository's own history,
with `jev-1.13.0` through TypeSafe: 38 Backlog Entries typed `feat`, labeled by
what happened to each (23 delivered as asked, 15 reshaped, split or declined),
20 fix entries and 10 Findings as controls, six controls the author wrote and
16 Portuguese versions paired with their English text; 270 requests, US$0.0146.
`docs/references/idea-stage-judgment-measurement.md` holds the set, the answers
and the measures. `shape_open` ranked the feature items at AUROC 0.75 (95%
interval 0.58 to 0.90), `unverified_assumption` at 0.66 (0.48 to 0.83) and
`market_question` at 0.81 (0.65 to 0.94), against 0.68 for the length of the
text alone. At yes from 0.7 and no up to 0.3, eight of the twelve items read
as indicated were in the reshaped, split or declined class, one of the eight
read as not indicated was, and 18 of 38 were undecided. Portuguese moved the
answers by 0.06 or less on average and one verdict in 16 by one step. Two
passes of the same model moved six of 90 verdicts by one step. `jev-1.14.0`
gave a different verdict in nine of 90 cases. No request of a session was
measured.

## Decision

**Two inputs.** `roundfix spec judge --stage idea` takes exactly one of them.
With `--item <reference>` it resolves the item through the Follow-up Command's
own inventory: an open Backlog Entry, a note in `docs/_inbox/` or an entry of
a messages file the caller hands in; an unknown reference and an active Spec
are usage errors. With `--request <path>` it reads one file that holds a
request the user made in the session, which the session wrote outside the
repository and removes afterwards. Nothing else can become a request: no
other flag, file or standard input adds text, and a Spec artifact and source
code cannot be sent.

**The boundary widens by one kind of text.** By the maintainer's answer above,
a request of the session may leave the machine. Of the request file the judge
receives only the first line, as the title, cut to 200 characters, and the
first 4,000 characters, with the project's name, after the same secret scan.
This amends ADR-0201, whose requests held only repository artifacts and whose
later decisions added Findings, Backlog Entries and Follow-up Items.

**The state is the follow-up judge's.** The project's name and the item's
kind, type, origin, title and the first 4,000 characters of its text, after
the secret scan; an item that matches a secret pattern is not sent. A request
of the session has the kind and the type `request` and the project's name as
its origin. Items are sent whatever their language. The recipients, the keys,
the Judge Log and the monthly ceiling stay as ADR-0201,
[ADR-0231](0231-the-jev-ceiling-is-a-user-config-value.md) and
[ADR-0239](0239-each-openrouter-stage-reads-its-own-roundfix-key-first.md)
decide, and the Run Marker as
[ADR-0283](0283-a-run-marker-keeps-the-judge-silent-inside-a-run.md) decides.
The Judge Log keeps a hash of the state, never its text.

**Less than the boundary allows.** A Finding, an item typed `fix`, `perf` or
`refactor`, and a Waiting Item are not sent and are reported as not judged: in
the measurement nine of 20 fix entries and seven of 10 Findings read as
indicated, because the fix of a defect is often not designed yet, and
[ADR-0353](0353-a-waiting-item-is-held-back-from-the-follow-up-list-unless-all-is-asked.md)
keeps a Waiting Item from the judge.

**Three questions, asked together.** `shape_open`, `unverified_assumption` and
`market_question`, worded as measured, in one request. Three other questions
were measured and dropped: two describe the kind of request and not what the
Idea stage studies, and one read as yes the requests that only ask for a
measurement.

**The starting thresholds, and the pinned model.** An answer is yes at 0.70 or
more and no at 0.30 or less; anything between is unsure. These are TypeSafe's
starting points, kept because 15 positive and 23 negative cases cannot tell
one cut point from its neighbours. The answers count only from the pinned Jev
1.13 model of ADR-0200; an answer from any other model is reported as not
judged.

**Fail-open, exit 0.** When the judge does not answer, the command says why,
prints the three criteria with their yes and no meanings for the agent to
answer, and exits `0`. It changes no other command's output or exit code.

## Considered Options

- **Keep the judge to inventoried items.** The first boundary. Rejected by the
  maintainer's answer: requests born in the session would be the only ones
  judged by the agent alone.
- **Send the whole request file, or hand it in through the messages schema.**
  Rejected: the answer bounds what leaves the machine to the first line and
  the first 4,000 characters, and a request has no sender, id or date to
  invent.
- **A language gate, as ADR-0201 has for Spec artifacts.** Rejected: the
  measurement found no shift in Portuguese larger than the shift between two
  passes of the same model, and this fleet's requests arrive in Portuguese.
- **Tune the thresholds on the measured set.** Rejected: no cut point is
  distinguishable from its neighbours on 38 cases, and TypeSafe's own guidance
  reads "Avoid fitting thresholds closely to a handful of test cases."
- **Add the question about whether the work is worth building.** Rejected for
  now: it raised recall on the measured set and also read as yes the items
  that are themselves measurements; its interval includes 0.5.
- **Measure and pin `jev-1.14.0` for this judgment.** Rejected: a second pin
  for one judgment, and nine of 90 verdicts differ; moving the pin means
  measuring every judgment again.
- **Judge Findings and fix entries too.** Rejected: the Idea stage answers a
  product question, and the questions read a correction as an open shape.

## Consequences

- A request of the session reaches the judge with no record in the
  repository; the record of the judgment, in the Idea or the PRD, names it
  `the request in the session`.
- The thresholds were measured on Backlog Entries. A request in the user's
  own words, often short and in Portuguese, is judged with the same starting
  points until the records show otherwise.
- The undecided band is wide: on the measured history about half of the
  feature items go to the user.
- A verdict near a cut point can change between two runs; the output always
  shows the probabilities, and the rule treats the verdict as advice.
- A request costs about US$0.00005 at the measured size, under the same
  monthly ceiling.
- The record each feature Spec keeps is what a later measurement needs to fix
  the thresholds.
