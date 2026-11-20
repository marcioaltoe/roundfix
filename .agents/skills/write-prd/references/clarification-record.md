# Clarification Record

A PRD's `## Decisions` section records who settled its questions. When a person
answers a question, write one dated line for that question and its answer:

```markdown
- YYYY-MM-DD: <question> — <answer>
```

Use the actual calendar date of the answer. Grilling stays the interactive
entry when a person drives the authoring; use answers already given in the
conversation rather than interviewing them again.

When no person was present, name the source the author decided from:

```markdown
Autonomous route: no person present; decisions taken from <source>.
```

Replace `<source>` with the source of those decisions. Both shapes may appear
together when the author also settles questions from sources after a person's
answers. Do not invent an interview or a dated answer for an unattended cycle.

The Spec Consistency Check reports `SC-DECISIONS-UNRECORDED` as a gap from the
PRD stage when `## Decisions` is absent or contains neither shape. A dated line
must begin `- YYYY-MM-DD: ` followed by non-whitespace text and have a valid
calendar date; the autonomous-route line must have a non-blank source and end
with a period. The check reads the record's shape, never whether a person was
really present. Its fix names the write-prd rule "Clarification Record".

The horizon is the oldest commit that added
`.agents/skills/write-prd/references/clarification-record.md`. A PRD whose
adding commit is at or after that commit is held, as is an uncommitted PRD.
A PRD committed before the horizon or a repository without the guide skips the
code with `a PRD committed at or after .agents/skills/write-prd/references/clarification-record.md`;
a missing PRD skips it with its path. `--strict` promotes the gap to an error.

## Delegated Decisions

When a session adopts the option an author recommended for an Open Authoring
Question, with nobody present, the adopted option is a Delegated Decision. It
is a third entry of the record and never one of the two shapes above: do not
write it as a dated answer of a person, and do not write it under the
autonomous-route line.

```markdown
- delegated: YYYY-MM-DD, <prd|techspec|tasks>
  - question: <the question, in one line>
  - options: <each option, separated by semicolons>
  - adopted: <the adopted option> — <the author's reason>
  - under: <Project Config | the user's answer of YYYY-MM-DD>
  - answer: none yet
```

The date is the day the option was adopted and the phase is the Authoring
Phase that asked. The five field lines are indented two spaces, each appears
once and each fits one line. `under` names the grant: `Project Config` when
the project's Authoring Question Preference is `follow`, or `the user's answer
of YYYY-MM-DD` when the user answered before the flow that the session follows
the author. A question of the TechSpec or the Tasks phase is recorded here, in
the PRD.

When the user later answers the question, replace `none yet`:

```markdown
  - answer: YYYY-MM-DD, kept — <the user's words>
  - answer: YYYY-MM-DD, changed — <the user's words>
```

Write the answer into the entry only while no Run and no delivery of the Spec
has started. The status view of `roundfix follow-up --status` lists the
entries of every active Spec.

`SC-DECISIONS-UNRECORDED` also reports, as a gap, an entry that begins
`- delegated:` and cannot be read: a head without a valid date and phase, a
missing or repeated field, an `adopted` line without the option, an em dash
and the reason, an `under` line that names neither grant, or an `answer` line
in neither form. A Delegated Decision never satisfies the rule above: the PRD
still needs a dated answer or the autonomous-route line.
