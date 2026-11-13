---
spec: 0388-a-completion-notice-skips-a-record-the-merge-adopted
prd: _prd.md
created: 2026-11-13
---

# A completion notice skips a record the merge adopted — Technical Spec

## Executive Summary

The defect is skill text. The closure bullet of `## Completion notices` in the
roundfix skill's `references/release.md` treats every path its deletion
listing prints as a closed record, and that listing prints every adopted
record too. This Spec rewrites the bullet to set adopted records aside, adds
one paragraph that tells how to recognize an adopted record with Git alone,
restores `memory_message_send` in the send paragraph, syncs the shipped
mirror, raises and records the skill's version, and pins every guarantee of the
section: phrases for each sentence that promises something, and an executable
check that runs the section's own Git commands on a temporary repository. No
Go production code, Baseline module, guide, glossary or ADR changes.

The trade-off accepted: the rule reads what the merge adds under the Spec Root
instead of the deleted record's own frontmatter, so it needs two more Git
commands per merge. The frontmatter at `<merge commit>^` is the cheaper read,
but the characterization shows it still holds `open` or `pending` for an
adopted and a closed record alike, so it cannot decide.

## Project Constraints

- Identifier strategy: not applicable — no identifier, record or schema is
  minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  network call is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0396 applies unchanged; its
  cited ADR-0383, ADR-0390, ADR-0330 and ADR-0319 do not apply, as
  `_prd.md` → Project Constraints states. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  answers of 2026-11-13 quoted in the authorization record; the Governed
  Paths are the `roundfix` skill's `SKILL.md` pair and its canonical release
  reference. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0388-a-completion-notice-skips-a-record-the-merge-adopted/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/release.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-13 at `f953f77f`.

1. `## Completion notices` runs from line 117 to the end of
   `.agents/skills/roundfix/references/release.md`, byte-identical to
   `skills/roundfix/references/release.md`. Its closure bullet (lines 145 to
   153) lists the records a merge deletes with
   `git diff --name-only --diff-filter=D <merge commit>^ <merge commit> -- docs/backlog docs/findings`
   and sends one notice for each Notice Request each holds.
2. The adoption contract of `docs/agents/docs-layout.md` moves an adopted
   Backlog Entry into `<spec root>/<slug>/references/` with `status: promoted`
   and `spec` set; a Finding becomes `done` (or `partial`) when its Spec is
   linked, and has no `spec` field. Squash merges land the status change and
   the move in one commit, so the record at `<merge commit>^` reads `open` or
   `pending`.
3. In this repository's planning merges `73421db2`, `a99c17b2`, `dfb4dbac`
   and `62bc941a`, the closure command lists each adopted entry, while an
   unlimited `git diff --name-status` pairs it as `R094`-`R096`. Git limits
   the diff to the pathspec before it looks for renames, so no setting pairs a
   move whose destination lies outside it.
4. `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest`
   (`skills/roundfix_modes_test.go`) reads both layouts (`canonical`,
   `embedded`) and pins 20 phrases of the section; none excludes an adopted
   record, and none pins the request's home, the no-release condition, the
   delivery wait, the cancel's return of sources, the send channel, or the
   report fallback.
5. `roundfix` is at 0.1.228 in both `SKILL.md` copies; the record command
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   raises it and appends the digest.

## System Architecture

- `.agents/skills/roundfix/references/release.md` and its mirror: the closure
  bullet, one new paragraph after the list of events, and the send paragraph.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` and
  `skills/testdata/owned-skill-versions.json`: the version raise, written by
  the record command.
- `skills/roundfix_modes_test.go`: more phrases in the existing test, the new
  test and its helpers, in the same file and package (`package skills`), with
  `roundfix/internal/gittest` for the temporary repository.

## Implementation Design

### Interfaces

The closure bullet and the paragraph after it replace lines 145 to 153 of the
canonical reference, byte for byte (the bullet keeps its two-space
indentation; one blank line separates them):

```markdown
- **A record is closed without a Spec.** Every terminal Backlog Entry or
  Finding is deleted in the change that records its disposition, so after
  each merge list the records it deletes with
  `git diff --name-only --diff-filter=D <merge commit>^ <merge commit> -- docs/backlog docs/findings`,
  read each one with `git show <merge commit>^:<path>`, set aside each record
  the merge adopted, as the next paragraph says, and send one notice for each
  Notice Request each other record holds, naming the merge commit: what was
  delivered when its work shipped without a Spec or, when the work is closed
  without being delivered, such as a declined Backlog Entry, the notice says
  so instead.

A record the merge adopted is not closed: its notice follows the Spec that
adopted it, and the closure event sends nothing for it. The adoption contract
of `docs/agents/docs-layout.md` moves an adopted Backlog Entry or Finding into
a Spec's `references/` folder, and a diff limited to `docs/backlog` and
`docs/findings` cannot pair that move with its destination under any rename
setting, so the listing of deleted records names every adopted record too.
The record read at `<merge commit>^` still carries the status it had before
the merge, so it cannot tell an adopted record from a closed one. List the
files the merge adds under the Spec root with
`git diff --name-only --diff-filter=A <merge commit>^ <merge commit> -- <spec root>`
and the files there that mention a record with
`git grep -l -F -e <file name> <merge commit> -- <spec root>`. A listed record
was adopted when the merge adds a file of the same name under a Spec's
`references/` folder whose frontmatter, read with
`git show <merge commit>:<added path>`, records the adoption: `status: promoted`
with `spec` naming that Spec for a Backlog Entry, or `status: done` or
`status: partial` for a Finding. It was also adopted when a PRD's `## Sources`
entry names it at the merge commit. When the merge adds a file of the same
name whose frontmatter cannot be read or records another status, and no PRD's
`## Sources` entry names the record, the event sends no notice for it, and the
session names the record and the added file and asks the user.
```

The send paragraph keeps every word after "Send each notice once" and becomes,
wrapped at 78 columns:

```markdown
Send each notice once with `memory_message_send`, from the session that
follows the event and never from a Daemon-assigned Run, and only to the
project that asked: a completion notice is not a release note. Each notice
names the record that kept its request, so two requests from the same project
receive two notices. A tag and a merge commit each occur once, so each request
falls due at exactly one event, and a session sends notices only for the event
it follows. Before sending, list the outbox with `memory_message_list` and
`box: "outbox"`: a pending notice to the same project that names the same
record and the same event was already sent, and is not sent again. A session
that cannot tell whether a notice for that request already went, because an
earlier session for the same event may have sent it and the recipient may have
consumed it, asks the user instead of sending a second one. When the session
exposes no ai-memory tools, put each notice's text in the report for the user
to deliver.
```

The test helpers, in `skills/roundfix_modes_test.go`:

```go
// every inline-code span of the section that starts with "git ", whitespace-normalized
func completionNoticeGitCommands(section string) []string
// runs the section command starting with prefix in root under diff.renames=<renames>,
// placeholders replaced; exit 1 with no output (git grep, no match) is empty
func completionNoticeRun(t *testing.T, root, renames string, commands []string, prefix string, values map[string]string) []string
// applies the closure rule with the section's commands: adopted, closed, asks the user
func completionNoticeClassify(t *testing.T, root, mergeCommit, renames string, commands []string) (adopted, closed, ask []string)
// a repository whose last commit on main (squashed, or a merge commit of two branch commits) adopts, closes and moves records
func completionNoticeFixture(t *testing.T, mergeShape bool) (root, mergeCommit string)
```

### Data Models

None. The version record keeps its schema; one `roundfix` entry is appended
by the record command.

### API Contracts

1. API Contract: the section text. In both copies of `references/release.md`
   the closure bullet, the paragraph after it and the send paragraph read as
   Interfaces writes them; every other byte of the file equals the Delivery
   Base; the two copies are byte-identical.
2. API Contract: the guarantee pins.
   `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest` keeps its
   20 phrases and adds the 16 of Testing Approach, one for each guarantee of
   the inventory outside the adoption rule;
   `TestACompletionNoticeSkipsARecordTheMergeAdopted` pins the 13 phrases of
   the adoption rule.
3. API Contract: the executable check.
   `TestACompletionNoticeSkipsARecordTheMergeAdopted/<layout>/<shape>`, for
   layouts `canonical` and `embedded` and shapes `squash` and `merge`, builds
   the fixture of Testing Approach, runs the deletion listing, the addition
   listing and the `git grep` command exactly as the section writes them, with
   `<merge commit>`, `<spec root>` (`docs/specs`) and `<file name>` replaced,
   under `diff.renames` `false`, `true` and `copies`, and reports
   `diff.renames=<v>: <class> records are <got>, want <want>` for each class
   that differs. A command the section no longer writes fails the subtest with
   `completion notices prescribe no command starting "<prefix>"`.
4. API Contract: the skill version. Both `SKILL.md` copies of `roundfix` carry
   the next version written by the record command, recorded with the folder
   digest; no version or digest is typed by hand.

### Surface Transcripts

None. No command, flag, exit code, message or stream placement changes; skill
text, its version and tests change.

### Guarantee inventory

Each guarantee of the section after this Spec and the phrase that pins it (E =
the existing test, N = the new test; "kept" = pinned at `f953f77f`):

1. One notice per Notice Request under the continuity rule — E: "receives one completion notice under the cross-project continuity rule of `docs/agents/agent-instructions.md`".
2. The request lives in its record — E: "Until then the request lives in its record: …adopts it." (whole sentence).
3. The adopting PRD carries the request — E kept: "carries the request into its `## Sources` entry"; E: "so the Spec's own record keeps it after the Backlog Entry or Finding moves into the Spec".
4. Several Specs name each other — E kept.
5. One event, found with Git alone — E: "Each notice falls due at one event, and the session that follows that event finds the requests due with Git alone".
6. Release event: commands and content — E kept (three phrases); E: "send one notice for each Notice Request its `## Sources` records".
7. No-release event: condition — E: "When the release plan for the range that holds the merge answers `no_release`, or the maintainer decides that no release carries it"; commands and naming kept.
8. Closure event: retirement premise — E: "Every terminal Backlog Entry or Finding is deleted in the change that records its disposition"; commands kept; E: "send one notice for each Notice Request each other record holds, naming the merge commit"; E: "what was delivered when its work shipped without a Spec"; E kept: "the notice says so instead".
9. Adoption rule — N: the 13 phrases of Testing Approach, and the executable check.
10. Delivery waits for the last Spec — E kept; E: "while one of them remains there at the new tag or merge commit, the event sends nothing for that request".
11. Grouping by message — E kept (two); E: "entries in several Specs that the same event archives and that record the same message are one request".
12. Cancel — E kept; E: "the cancel returns each adopted Backlog Entry and Finding to its folder, and the request stays there".
13. Channel and sender — E: "Send each notice once with `memory_message_send`, from the session that follows the event"; E kept: Daemon and recipient.
14. One notice per record, one event per request — E kept (two); E: "a session sends notices only for the event it follows".
15. Outbox and doubt — E kept (three); E: "a pending notice to the same project that names the same record and the same event was already sent, and is not sent again".
16. No ai-memory — E: "When the session exposes no ai-memory tools, put each notice's text in the report for the user to deliver."

### Invariants

```text
1. No record the merge adopted receives a notice at the closure event, under any rename setting and either merge shape.
2. A record the merge deleted with no adopting copy and no PRD naming it is a closure and receives its notice.
3. An adopting copy whose frontmatter does not confirm the adoption, with no PRD naming the record, sends nothing and asks the user.
4. Every guarantee of the inventory has its own pinned phrase in both layouts.
```

### Sabotages

Each runs on a copy of the audited head and names the test that observes it;
the prototype at `f953f77f` failed each as listed.

1. Delete the sentence "A record the merge adopted is not closed: …sends nothing for it." from the canonical copy only: `TestACompletionNoticeSkipsARecordTheMergeAdopted/canonical` fails.
2. The same in the shipped copy only: `TestACompletionNoticeSkipsARecordTheMergeAdopted/embedded` fails.
3. Narrow the addition listing's pathspec to `-- <spec root>/*/_prd.md` in the canonical copy: `TestACompletionNoticeSkipsARecordTheMergeAdopted/canonical/squash` and `/canonical/merge` fail with the unclear entry classified as closed.
4. Put the v0.143.0 text of the canonical reference back: `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest/canonical` and `TestACompletionNoticeSkipsARecordTheMergeAdopted/canonical` fail.
5. Drop "with `memory_message_send`" from the send paragraph: `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest/canonical` fails.
6. Delete "It was also adopted when a PRD's `## Sources` entry names it at the merge commit.": `TestACompletionNoticeSkipsARecordTheMergeAdopted/canonical` fails.
7. Shorten the frontmatter values to "`status: promoted`" alone: `TestACompletionNoticeSkipsARecordTheMergeAdopted/canonical` fails.

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1 → API Contract 1, Invariants 1 and 2, task_01
- Goal 2, Core Feature 1, Success Metric 2 → API Contract 3, Invariant 3, task_01
- Goal 3, Core Features 2 and 3, Success Metric 3 → API Contract 2, Guarantee inventory, Invariant 4, Sabotages, task_01
- Goal 4, Core Feature 4, Success Metric 4 → API Contract 4, task_01; task_02 also records the PRD's Unreachable Acceptance as one declared row

## Absence Evidence

1. The executable check of the section's Git commands — command: `grep -rn -e "diff-filter" -e "completionNotice" skills/*.go`; result: at `f953f77f` the only matches are the two `--diff-filter` phrases of `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest`, which compare text; no test runs a command the release reference prescribes.
2. The adoption rule in the section — command: `grep -n -i -e "adopted" -e "promoted" .agents/skills/roundfix/references/release.md`; result: at `f953f77f` the only match is line 163, "the cancel returns each adopted Backlog Entry and Finding to its folder"; nothing sets an adopted record aside at the closure event.

## ADR Alternatives

- ADR-0396: the design keeps its decision and takes none of its rejected options; the notice stays a session's job, read from the adopting PRD's `## Sources`, with no field a command carries.
- ADR-0383: not applicable — it decides when a message expects a reply, and this Spec changes no message clause; it is listed only because ADR-0396 cites it.
- ADR-0390: not applicable — it decides the precedence of blocking messages and the answer to direct questions, which this Spec leaves unchanged; it is listed only because ADR-0396 cites it.
- ADR-0330: not applicable — it decides that ai-memory carries messages, and this Spec adds no channel; it is listed only because ADR-0396 cites it.
- ADR-0319: not applicable — it decides the follow-up order, and this Spec changes no ordering; it is listed only because ADR-0396 cites it.

## Integration Points

None. The new test runs the local `git` (a Base Test Tool) in `t.TempDir()`
through `roundfix/internal/gittest`, with the isolated configuration
`gittest.IsolatedEnv` gives; no network, provider or ai-memory.

## Testing Approach

- `skills/roundfix_modes_test.go`, `package skills`:
  - `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest` adds the 16
    phrases the inventory marks "E:" (not "E kept"), in the order of the
    inventory.
  - `TestACompletionNoticeSkipsARecordTheMergeAdopted`, subtests per layout,
    each first checking the 13 adoption phrases: "set aside each record the
    merge adopted, as the next paragraph says"; "A record the merge adopted is
    not closed: its notice follows the Spec that adopted it, and the closure
    event sends nothing for it."; "cannot pair that move with its destination
    under any rename setting, so the listing of deleted records names every
    adopted record too"; "The record read at `<merge commit>^` still carries
    the status it had before the merge, so it cannot tell an adopted record
    from a closed one."; the addition command; the `git grep` command;
    "git show <merge commit>:<added path>"; "`status: promoted` with `spec`
    naming that Spec for a Backlog Entry, or `status: done` or
    `status: partial` for a Finding"; "It was also adopted when a PRD's
    `## Sources` entry names it at the merge commit."; "the event sends no
    notice for it, and the session names the record and the added file and
    asks the user"; "moves an adopted Backlog Entry or Finding into a Spec's
    `references/` folder"; "A listed record was adopted when the merge adds a
    file of the same name under a Spec's `references/` folder whose
    frontmatter"; "When the merge adds a file of the same name whose
    frontmatter cannot be read or records another status, and no PRD's
    `## Sources` entry names the record". Then subtests `squash` and `merge`.
  - The fixture: on `main`, `docs/backlog/2026-01-01-adopted-entry.md`,
    `-02-declined-entry.md`, `-04-unconfirmed-entry.md`,
    `-05-named-entry.md` (all `status: open`) and
    `docs/findings/2026-01-03-adopted-finding.md` (`pending`); then one change
    that writes `docs/specs/0001-example/_prd.md` whose `## Sources` names
    entries 01, 03 and 05, moves 01 into `references/` as `promoted` with
    `spec: 0001-example`, moves 03 as `done`, deletes 02 and 05, and moves 04
    still `open`. Expected: adopted 01, 05 and 03 (in listing order: backlog
    01, backlog 05, findings 03); closed 02; asks the user 04.
- Unchanged and run: `go test ./skills` (mirror, version record, version
  raise, other mode contracts), `go test -tags repocontract ./skills`, and the
  documentation contracts with `-tags docscontract`.
- The file compiles for Windows (`GOOS=windows go vet ./skills`): it uses
  `os/exec` and no platform call.

## Build Order

1. The section text, its sync, the version record and the tests (task_01; no
   dependency).
2. QA gate (task_02) (depends on: 1).

No Spec is listed in `requires`; Spec 0389 delivers after this one and also
raises the `roundfix` version.

## Risks & Considerations

- The rule asks the session to run two more commands per merge and read a few
  frontmatter blocks; a merge adopts a handful of records at most.
- A Spec Root in another repository leaves no adopting copy in this one; the
  rule then relies on the PRD's `## Sources`, which that layout also keeps
  elsewhere. Stated as a non-goal in the PRD.
- The record command chooses the next version from the tree it runs on; if
  Spec 0389 raised `roundfix` first, this Spec raises past it.
- Adopters keep the v0.143.0 text until their next `roundfix baseline update`
  after v0.143.1.

## Glossary

None.

## Decisions

- Recognize adoption by the copy the merge adds and its frontmatter, or by a
  PRD's `## Sources`; never by the frontmatter at `<merge commit>^`.
- Ask the user, never send, when a copy exists but does not confirm the
  adoption.
- Keep the v0.143.0 deletion listing; restore `memory_message_send` and round
  1's `## Sources` route; keep round 2's cancel rule.
- Pin each guarantee with its own phrase and run the section's commands in a
  test, in the existing contract file. No ADR: ADR-0396 already decides where
  an adopted request lives.
