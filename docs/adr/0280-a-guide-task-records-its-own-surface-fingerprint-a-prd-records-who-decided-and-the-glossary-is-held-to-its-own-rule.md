---
status: accepted
created_at: 2026-10-16T00:00:00Z
updated_at: 2026-10-16T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A guide Task records its own surface fingerprint, a PRD records who decided, and the glossary is held to its own rule

Three authoring guards had a rule without a check. On 2026-10-14 a Task of
Spec 0268 edited three user guides while its authorization's Limits forbade
`docs/references/behavior-surfaces.json`; its Verification failed twice and the
Run ended Unresolved, although `roundfix spec check --strict` was clean on the
plan. The `roundfix` skill stated the QA secret refusal in a form the Daemon
does not print. ADR-0244 holds every Spec to the glossary, while nothing held
`CONTEXT.md` to itself, and a congruence audit on 2026-10-14 found sixteen
undefined phrases and three closed lists that had drifted for months. On the
same day the maintainer noted that "a skill `grilling` não é usada
corretamente na criação de PRDs": PRDs of unattended cycles were written
without an interview by construction, and nothing in them said so. On
2026-10-16 the maintainer approved cycle AH, whose item C is this decision.

**A guide Task records its own fingerprint.** The Behavior Surface Record holds
the bytes of every user guide page, and its check runs in the repository
Verification at every Task's settlement. A pending non-QA Task that declares a
page under `docs/user-guide/` therefore declares the record itself and names
its `-record-skill-coverage` step in a requirement; another Task of the Spec
cannot cover it, because the first Task settles before the second runs. The
Spec Consistency Check reports such a Task, and a Spec whose authorization
Limits name the record, as the gap `SC-SURFACE-RECORD-UNTASKED`, with no
horizon, as ADR-0265 holds structural proof gaps. A change to a command's help
is not inferred from the files a Task declares, because one source file holds
the usage text of many commands; the record's own check keeps that case.

**A PRD records who decided.** The `## Decisions` section of a PRD is its
Clarification Record: one dated line for each question a person answered, or
the line `Autonomous route: no person present; decisions taken from <source>.`
The check reports a PRD with neither as the gap `SC-DECISIONS-UNRECORDED`. Like
the glossary declaration of ADR-0244, the rule holds only PRDs committed at or
after the commit that adds its guide to the `write-prd` skill: a PRD already
approved cannot be asked again who answered it, and inventing a date would
record a falsehood. The check reads the shape of the record, never whether a
person was really present.

**The glossary is held to its own rule.** A Repository Contract Test reads
`CONTEXT.md` and fails when a capitalized phrase in a definition is neither a
defined term, a plural of one, nor an allowlisted name; when a link names a
file that is not an accepted ADR; when a mapped closed list differs from the
exported constants it mirrors; when an `_Avoid_` line is not a plain
comma-separated list; or when a sentence of more than sixty characters appears
in two entries. The allowlist and the constant mapping live in that test, and
each allowlisted phrase carries its reason. An `_Avoid_` item may name another
defined term, because a neighbouring term is what a reader confuses.

## Considered Options

- Letting another Task of the Spec record the fingerprint, as the Backlog Entry
  first proposed, was rejected: the guide Task's own settlement runs the
  record's check and fails before that Task runs.
- Inferring a help change from a declared command source was rejected: it
  would flag every Task that touches a command's file without changing its
  usage text.
- A Clarification Record held for every active Spec, as ADR-0265 holds its
  proof sections, was rejected: those sections can be written at any time,
  while who answered a question cannot be reconstructed after the PRD was
  approved.
- Requiring `grilling` before every PRD was rejected: the Baseline does not
  wait for a person on glossary work, and an unattended cycle has nobody to
  interview; the record states the route instead.
- Holding the glossary through `roundfix spec check` was rejected: the
  glossary is not a Spec artifact, and adopters' glossaries carry no closed
  list mapped to Roundfix's constants.
- Forbidding defined terms in `_Avoid_` lines was rejected: it would remove 77
  warnings about neighbouring terms that readers confuse.
