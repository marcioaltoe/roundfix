---
spec: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
prd: _prd.md
created: 2026-11-03
---

# The Spec skills announce the Spec Events the user defines — Technical Spec

## Executive Summary

Five owned skill surfaces perform the seven moments of a Spec's life, and none
of them tells anything outside the repository. The design adds one section,
`## Spec Event Announcement`, to each: `write-prd`, `write-techspec`,
`write-tasks` and `archive-spec` `SKILL.md`, and the Roundfix skill's
`references/implement.md`. Every section carries the same generic paragraph
(API Contract 1), which makes the step exist only when the user's own
instructions define its destination and format, and one paragraph naming the
events that surface announces (API Contracts 2 to 6). `CONTEXT.md` gains
**Spec Event** and **Spec Event Announcement**, the context-driven development
guide gains a short section with the owner of each event, and a new
skill-contract test holds the five surfaces to the wording and every owned
skill, canonical and mirror, free of the reported assistant's, tracker's,
channel's and block's names. No Go production code, Baseline module, guide
under `docs/agents/` or configuration changes.

The trade-off accepted: the skills name the facts an announcement carries but
not its words, so a user who wants announcements must write the destination
and the format into their own instructions; until then the step does
nothing.

## Project Constraints

- Identifier strategy: not applicable — nothing is minted; skill text, two
  glossary entries, one guide section and one test file. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path or credential;
  Roundfix's code sends nothing. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0354, written by this Spec,
  ADR-0354: "A session announces only when the user's own instructions,
  global or for the repository, define both". The owned-skill version rule and
  the `### QA settlement` rule are repository rules; the design keeps both.
  Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AW answers of 2026-11-03 quoted in the authorization record. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0344-the-spec-skills-announce-the-spec-events-the-user-defines/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Current behavior

Read at `4e61bf6d`.

- `write-prd` creates the Spec folder in its step 5 and writes `_prd.md` in
  step 6; its last sections are `## Delivery Type and release decision`,
  `## Advisory judgment`, `## Anti-patterns` and `## References`.
- `write-techspec` mints the folder and a minimal `_prd.md` only on its
  correction route (`## Preconditions — two entry modes`); its sections end
  with the same three after `## Delivery Type and release decision`.
- `write-tasks` reads `_authorization.md` in its Project Constraint preflight
  and writes the graph in `### 3. Write the files`; `## Release decision`
  precedes `## Process`.
- `archive-spec` drives `roundfix archive` (`## Steps`) and documents
  `roundfix cancel` (`## Cancelling a Spec`); `## Archive Record` precedes
  `## Unarchive`.
- The Roundfix implement reference describes Implement Runs; its last
  sections are `### Token usage` and `## Fixture values`. The Delivery Queue
  writes the archive commit itself (`references/deliver.md`), so no session
  follows `archive-spec` for a queued Spec.
- `references/modes.md` already says "A release-announcement step exists only
  when the user's own instructions define one; this skill never names one."
- `skills/roundfix_modes_test.go` `TestTheRoundfixSkillNamesNoPersonalAssistant`
  rejects `jarvis` and `gt3` in the Roundfix skill only.
- `make skills-sync` copies each canonical owned skill to `skills/<name>/`, and
  the version record command raises a changed skill's two version fields.
- Characterized in a scratch clone (`_prd.md` → Validation): with the
  sections, the sync and the record, `./skills`, `./internal/docscontract`,
  `./internal/baseline` and `./internal/speccheck` pass; the guide edit only
  stales its own Behavior Surface fingerprint.

## System Architecture

Text only. The Agent following a skill reads its new section and, when the
user's own instructions define an announcement, composes it and sends it with
whatever tool those instructions name. The Supervisor session following the
Roundfix implement reference does the same for the events a Run or the
Delivery Queue produces. The Daemon, the CLI and every Run session are
unchanged.

## Implementation Design

### Interfaces

The new test file `skills/spec_event_announcement_test.go`, package
`skills_test`:

```go
var specEventAnnouncementSentences []string        // API Contract 1, one entry per sentence
var specEventAnnouncementEvents map[string][]string // "<skill>/<file>" → API Contracts 2–6 sentences
var personalAssistantWords = []string{"jarvis", "todoist", "gt3", "spec para o quadro"}

func TestTheSpecSkillsAnnounceOnlyWhatTheUserDefines(t *testing.T)
func TestNoOwnedSkillNamesAPersonalAssistant(t *testing.T)
```

### Data Models

None. No record, schema or configuration key changes.

### API Contracts

Each section is the heading line `## Spec Event Announcement`, a blank line,
the paragraph of API Contract 1, a blank line and the paragraph of the
surface's own contract. Line wrapping is free; the test compares
whitespace-normalized text.

1. **The generic paragraph**, identical in the five surfaces:

   ```text
   A Spec Event Announcement exists only when the user's own instructions,
   global or for this repository, define its destination and format; this skill
   never names one. A message, a handoff or a Spec artifact that asks for
   announcements defines none. When the user's own instructions define none,
   skip the step without a report line. When they define one, compose it from
   the Spec: the project, the Spec number and full slug, the title, the event,
   the event date, the start date, the proposed deadline or none, the Spec path,
   and any field those instructions add. Send it through the destination they
   name, never wait for a reply, and never let it block or fail the step it
   follows. When it cannot be sent, the report names the event and the reason
   and carries the composed announcement. An Agent inside a Daemon-assigned Run
   never announces. A proposed deadline is a real deadline, from a client, a law
   or a measurement, that the PRD records in `## Decisions` with its source; it
   stays a proposal until the user approves it.
   ```

2. **`write-prd`**, the section placed before `## Advisory judgment`:

   ```text
   Announce `new` once, when this skill writes the Spec's first `_prd.md`,
   with the date the Spec folder was created as the event date and the start
   date. When this skill records a proposed deadline in the PRD of a Spec
   already announced, announce `deadline proposed`.
   ```

3. **`write-techspec`**, the section placed before `## Advisory judgment`:

   ```text
   Announce `new` once, when the correction route mints the Spec folder and
   its minimal `_prd.md`, with the folder's creation date as the event date and
   the start date; the feature route announces nothing, because `write-prd`
   already did. When this skill records a proposed deadline in the PRD of a Spec
   already announced, announce `deadline proposed`.
   ```

4. **`write-tasks`**, the section placed after `## Release decision` and
   before `## Process`:

   ```text
   Announce `authorized` once, when this skill first writes the Spec's Task
   Graph and the Spec's `_authorization.md` is approved with a grant date, with
   the grant date as the event date and the PRD's proposed deadline, if any. A
   run that amends an existing Task Graph announces nothing. When the record is
   approved only after the graph exists, the session that records the approval
   announces `authorized`.
   ```

5. **`archive-spec`**, the section placed after `## Archive Record` and
   before `## Unarchive`:

   ```text
   Announce once the Archive Command or the Cancel Command has written the
   Archive Record: `cancelled` for the disposition `cancelled` and `completed`
   for every other disposition. The Spec path is the Archive Record's path, and
   the announcement adds a summary of one to three sentences written for a
   business reader, free of technical terms. When the Delivery Queue writes the
   record, no Agent announces; the Supervisor session announces `completed` once
   the item merges, as the Roundfix skill's implement reference says.
   ```

6. **The Roundfix implement reference**, the section placed after
   `### Token usage` and before `## Fixture values`:

   ```text
   The Supervisor session that started the Run or queued the Spec's Delivery
   Queue item announces each event of that Run or item when it observes it, with
   the date the event happened as the event date:

   - `started` the first time a Task of the Spec enters `in_progress`, read from
     the Run Event Stream's `task-status` category or the Task files; a later
     Run, a retry or a carry-forward of the same Spec never announces `started`
     again.
   - `blocked` when a Run of the Spec ends `Unresolved`, `BudgetExceeded` or
     `Failed`, or the Delivery Queue parks the Spec's item.
   - `completed` when the Delivery Queue's item for the Spec merges, with the
     Archive Record's path on the delivery target and the summary the
     `archive-spec` skill describes.

   Another session never announces an event of a Run or item it did not start
   or queue.
   ```

7. **Glossary.** `CONTEXT.md` gains, right after the **Ignored Spec File**
   entry, through `domain-modeling`:

   ```text
   **Spec Event**:
   One of the seven moments in a Spec's life that a Spec Event Announcement reports: `new` when its first PRD is written, `authorized` when its first Task Graph is written for an approved authorization record, `deadline proposed` when a real deadline is added to a Spec already announced, `started` when its first Task enters `in_progress`, `blocked` when a Run of the Spec ends `Unresolved`, `BudgetExceeded` or `Failed` or the Delivery Queue parks its item, `completed` when its Archive Record is written or the Delivery Queue merges the Spec it archived, and `cancelled` when the Cancel Command writes its record ([ADR-0354](docs/adr/0354-the-spec-skills-announce-spec-events-only-through-the-users-own-instructions.md)).
   _Avoid_: Spec status, milestone, Run Event

   **Spec Event Announcement**:
   The message a session following `write-prd`, `write-techspec`, `write-tasks`, `archive-spec` or the Roundfix skill's implement reference sends at a Spec Event, only when the user's own instructions, global or for the repository, define its destination and format. It carries the project, the Spec number and full slug, the title, the event, the event date, the start date, the proposed deadline and the Spec path, never waits for a reply and never blocks or fails its step; no owned skill names a destination, and an Agent inside a Daemon-assigned Run never sends one ([ADR-0354](docs/adr/0354-the-spec-skills-announce-spec-events-only-through-the-users-own-instructions.md)).
   _Avoid_: Notification, release announcement, Spec report
   ```

8. **Guide.** `docs/user-guide/context-driven-development.md` gains, right
   before `### QA report reachability and verdicts`, a section
   `### Spec Event Announcements` with one paragraph and one table:

   ```text
   A user who follows Specs on a board, a tracker or an assistant can have the
   Spec skills announce each Spec Event as it happens. The skills never name a
   destination: the user's own instructions, global or for the repository,
   define the destination and the format, and without them no skill announces
   anything or adds a report line. Each announcement carries the project, the
   Spec number and full slug, the title, the event, the event date, the start
   date, the proposed deadline or none, and the Spec path, plus any field the
   instructions add; `completed` and `cancelled` add a summary of one to three
   sentences for a business reader. An announcement never waits for a reply and
   never blocks the step it follows, and an Agent inside a Daemon-assigned Run
   never announces ([ADR-0354](../adr/0354-the-spec-skills-announce-spec-events-only-through-the-users-own-instructions.md)).

   | Spec Event | Announced by | When |
   | --- | --- | --- |
   | `new` | `write-prd`, or `write-techspec` on its correction route | the Spec's first `_prd.md` is written |
   | `authorized` | `write-tasks` | the first Task Graph is written for an approved `_authorization.md` |
   | `deadline proposed` | `write-prd` or `write-techspec` | a real deadline is recorded in the PRD of a Spec already announced |
   | `started` | the Supervisor session that started the Run or queued the item | the first Task of the Spec enters `in_progress` |
   | `blocked` | the same Supervisor session | a Run of the Spec ends `Unresolved`, `BudgetExceeded` or `Failed`, or the Delivery Queue parks its item |
   | `completed` | `archive-spec`, or the Supervisor session when the Delivery Queue archives | the Archive Record is written, or the queue item merges |
   | `cancelled` | `archive-spec` | the Cancel Command writes the Archive Record |
   ```

9. **The contract test.**
   `TestTheSpecSkillsAnnounceOnlyWhatTheUserDefines` reads, under both
   `.agents/skills/` and `skills/`, the five files `write-prd/SKILL.md`,
   `write-techspec/SKILL.md`, `write-tasks/SKILL.md`,
   `archive-spec/SKILL.md` and `roundfix/references/implement.md`; requires
   exactly one line `## Spec Event Announcement`; takes the section up to the
   next `## ` heading, normalizes whitespace, and requires every sentence of
   API Contract 1 and of that file's contract, failing with the path and the
   missing sentence. `TestNoOwnedSkillNamesAPersonalAssistant` reads the
   `skills` keys of `skills/testdata/owned-skill-versions.json` (failing when
   none), walks every file of each owned skill under both trees, and fails
   with the repository-relative path and the word for each file whose
   lower-cased text contains a word of `personalAssistantWords`. Both call
   `t.Parallel()` first and start no process.

### Surface Transcripts

None. No command output changes.

### Invariants

```text
1. Every section of every owned skill other than the new Spec Event Announcement sections keeps its text; only the version lines of the five changed SKILL.md pairs move.
2. No `### QA settlement` section changes.
3. No owned skill file, canonical or mirror, contains jarvis, todoist, gt3 or "spec para o quadro" in any letter case.
4. Each mirror equals its canonical copy.
5. No Go production file, Baseline module, docs/agents guide or configuration file changes.
```

### Declared breaks

None. The new sections add text, the version record gains one entry per
changed skill as every owned-skill edit does, and the Behavior Surface Record
changes only the fingerprint of the context-driven development guide. No
existing test, golden or fixture pins the changed files' bytes outside the
version and fingerprint records (`grep -rn 'Spec Event' skills internal`
finds nothing).

## Coverage Map

- Goal 1, User Story 1, Core Features 1 and 2, Success Metric 1 → API Contracts 1 to 6 under Implementation Design, checked by `TestTheSpecSkillsAnnounceOnlyWhatTheUserDefines` (API Contract 9).
- Goal 2, User Story 3, Core Feature 3, Success Metric 2 → API Contract 9 under Implementation Design, `TestNoOwnedSkillNamesAPersonalAssistant`, and Invariant 3.
- Goal 3 → API Contract 1 under Implementation Design: the announcement never waits for a reply and never blocks or fails its step.
- Goal 4, User Story 2 → API Contract 1 under Implementation Design: no step and no report line without a definition; Invariant 1.
- Core Feature 4, Success Metric 3 → API Contracts 7 and 8 under Implementation Design.
- Success Metric 4 → Invariants 1, 2 and 4, the version record step of Build Order 1 to 3, and the skills and documentation contract suites of Testing Approach.

Deviations from the PRD: None.

## Absence Evidence

1. A Spec Event Announcement section — command: `grep -rn 'Spec Event' .agents/skills skills CONTEXT.md docs/user-guide docs/agents`; result: no match, exit 1.
2. A guard over every owned skill — command: `grep -ln -i 'todoist' skills/*_test.go`; result: no match, exit 1; the only guard, in `skills/roundfix_modes_test.go`, covers the Roundfix skill and two words.
3. The contract test file — command: `ls skills/spec_event_announcement_test.go`; result: no such file.

## ADR Alternatives

- ADR-0354: written by this Spec; its rejected alternatives (the assistant's
  block in the skills, a Baseline clause or configuration key, announcing from
  inside a Run, a report line without a definition) stay rejected: the
  sections name no destination, add no clause or key, forbid the Run Agent
  from announcing and stay silent without a definition.

## Integration Points

None. No code sends anything; the user's own instructions name the tool.

## Testing Approach

All tests read this repository's files; none starts a process or touches the
network.

- `skills/spec_event_announcement_test.go`: the two tests of API Contract 9.
- Existing suites: `go test ./skills` (mirrors, version record, the
  Roundfix layout and wording tests) and
  `go test -tags docscontract ./internal/docscontract` (Behavior Surface
  Record, guide contracts).
- Characterization before the change: the new test fails for each of the ten
  files with "has 0 Spec Event Announcement headings, want 1", and the guard
  passes.
- Sabotages, each on its own temporary copy of the audited head, each run in
  the scratch prototype while authoring:
  1. remove the sentence "An Agent inside a Daemon-assigned Run never
     announces." from the canonical `archive-spec` `SKILL.md`:
     `TestTheSpecSkillsAnnounceOnlyWhatTheUserDefines/.agents/skills/archive-spec/SKILL.md`
     fails naming the sentence;
  2. append "Send it to Jarvis." to the canonical `write-prd` PRD template
     (a file other than `SKILL.md`): `TestNoOwnedSkillNamesAPersonalAssistant`
     fails naming `.agents/skills/write-prd/references/prd-template.md` and
     `jarvis`;
  3. append "Track it in Todoist." to the mirror `skills/qa-gate/SKILL.md`
     (an owned skill this Spec does not edit, in the mirror tree, mixed
     case): the guard fails naming that path and `todoist`;
  4. append a second `## Spec Event Announcement` heading to the mirror
     `skills/roundfix/references/implement.md`: the first test fails with
     "has 2 Spec Event Announcement headings, want 1".

## Build Order

1. `CONTEXT.md` (API Contract 7) and the `write-prd` and `write-techspec`
   sections (API Contracts 1 to 3), with the mirrors synced and the versions
   recorded last (task_01).
2. The `write-tasks` and `archive-spec` sections (API Contracts 1, 4 and 5),
   mirrors and versions (task_02) (depends on: 1).
3. The Roundfix implement reference section (API Contracts 1 and 6), the
   Roundfix `SKILL.md` version, and the guide section (API Contract 8), with
   the Behavior Surface Record re-recorded after the last guide edit
   (task_03) (depends on: 2).
4. The contract test (API Contract 9) (task_04) (depends on: 3).
5. QA gate (task_05) (depends on: 4).

Each step depends on the previous because all edit
`skills/testdata/owned-skill-versions.json` (1 to 3) or read every section
(4). No Verification depends on another Spec's code, so `_tasks.md` has no
`requires`.

Shared files with the Specs authored in parallel (lesson of parallel
authoring): Spec 0343 edits the Roundfix skill's `references/modes.md` and
`references/follow-up.md`, raises the Roundfix `SKILL.md` version fields and
may edit `skills/roundfix_modes_test.go`; Spec 0342 may edit Roundfix
references, `.roundfixrc` documentation and user guides. All three may touch
`skills/testdata/owned-skill-versions.json`, the Roundfix `SKILL.md` version
lines, `CONTEXT.md` and `docs/references/behavior-surfaces.json`. This Spec
delivers last: its Runs regenerate the version record and the fingerprint
with their commands on top of what the other two merged, it adds its
glossary entries after **Ignored Spec File**, away from other entries, and it
touches neither `modes.md`, `follow-up.md` nor `roundfix_modes_test.go`.

## Risks & Considerations

- A Claude Agent inside a Run may load the user's global instructions and
  try to announce. Every section says it never does, and the operator's
  proposed global section repeats it; the Run's sandbox has no network.
- Two Supervisor sessions could both observe one Run. Only the session that
  started the Run or queued the item announces, so the event is sent once.
- An event no supervising session observed is lost (ADR-0354
  Consequences); reconstructing past events is out of scope.
- Spec 0343 may move the Roundfix version lines first; the record command
  raises from the highest recorded version, so no manual merge is needed.

## Glossary

- adds: **Spec Event** — one of the seven moments in a Spec's life that an announcement reports
- adds: **Spec Event Announcement** — the message a Spec skill's session sends at a Spec Event, only when the user's own instructions define it
- not a term: **Ignored Spec File** — existing glossary term, unchanged

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AW brief for Spec 0344 of 2026-11-03 and the maintainer's answers quoted in `_prd.md` → Decisions.

- The section lives under its own heading in each surface, never inside
  `### QA settlement`, and carries one shared paragraph so a single test
  list proves all five.
- The guard is a new test over every owned skill and both trees rather than
  a widened `TestTheRoundfixSkillNamesNoPersonalAssistant`, so this Spec does
  not edit a file Spec 0343 may change.
- The guide section goes into the context-driven development guide, which
  describes the Spec pipeline these skills implement; no command guide
  changes, because no command changes.
- No Baseline clause: the step is a user preference read at run time, not a
  repository rule (ADR-0354).
