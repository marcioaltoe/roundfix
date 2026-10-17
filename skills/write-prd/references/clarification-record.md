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
