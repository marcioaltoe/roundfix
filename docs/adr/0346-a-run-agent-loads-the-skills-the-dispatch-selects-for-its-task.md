---
status: accepted
created_at: 2026-10-30T18:00:00Z
updated_at: 2026-10-30T18:00:00Z
deprecated_at: null
superseded_by: null
---

# A Run Agent loads the skills the dispatch selects for its Task

The Baseline renders `docs/agents/skill-dispatch.md`, whose triggers and
activation bundles make skills mandatory for matching work, but the Task
prompt a Run Agent receives never names that guide, and a runtime keeps only
each skill's description in context until something invokes the skill. Over
60 days on the operator's machine, 108 Claude sessions of frontend Tasks of a
TypeScript adopter loaded `implement-task` in 21 and every frontend skill in
none: the mandatory frontend bundle had no effect.

**The Task prompt and `implement-task` carry the load rule.** The Daemon adds
`docs/agents/skill-dispatch.md` to the Spec Context Bundle when the
repository has it. The Task prompt's execution invariants tell the Agent,
before it edits anything, to read that guide, select only the triggers and
activation bundles that match its Task (its Task Type, the paths it declares
and the work it asks for), load each selected skill, and name them in
`## Result` on one `Skills loaded:` line, which the glossary calls the
Loaded Skill List. A skill only a person can start is not loaded in a Run and
is named as skipped. The `implement-task` skill states the same rule for a
standalone execution.

**The QA gate follows the same rule.** The QA contract tells the QA Agent to
load the skills of the QA activation and of each trigger the Spec's declared
surfaces match, and only those. The QA Report needs no new field: the gate's
evidence stays its matrix.

**The guide is a carry-forward input.** A Task Carry-Forward already treats
the root instructions and `implement-task` as inputs that declare the Task
(ADR-0284). The dispatch guide now shapes the same handoff, so a change to it
since settlement moves the Task's inputs.

**The Daemon does not check the list.** The Loaded Skill List is the Agent's
own claim. A Daemon check could prove only that the line was written, never
that a skill was loaded, so it would be a gate accepted on text. The evidence
that the rule works is the sessions' recorded tool calls, measured again
after the release.

Three alternatives were rejected. Injecting the selected skills' bodies into
the prompt would make the Daemon interpret free-text triggers, tie the prompt
to one runtime's skill layout and spend context on every Task whatever it
needs. Telling the Agent to load every installed or dispatched skill would
fill the context with skills that match nothing; the load rule loads only
what matches. Making the Loaded Skill List a Verification gate would refuse
correct work over prose (ADR-0057 keeps the Daemon the only owner of the
Task's verdict, decided by the declared Verification).
