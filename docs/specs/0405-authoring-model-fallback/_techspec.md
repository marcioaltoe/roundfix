---
spec: 0405-authoring-model-fallback
prd: _prd.md
created: 2026-11-17
---

# An Authoring Phase falls back to its declared fallback when the declared model's usage limit runs out — Technical Spec

## Executive Summary

The design adds one key to the Authoring Model Declaration (`fallback`, with `model` and `reasoning_effort`). It adds a new entry kind to the Phase Record, `fallback`, under the present schema. Two subcommands are added: `roundfix phase fallback` decides and records an Authoring Fallback from the error an Author Subagent ended on, and `roundfix phase model` prints the model and effort a phase's Author Subagent starts on. The four authoring skills and the Roundfix skill's phase reference tell the session to run them and pass what they print. Nothing is generated: the model and effort passed when a subagent starts override its file. `roundfix phase start`, the installation of skills and subagents, and the Run's fallbacks keep their contracts.

The primary trade-off: Roundfix recognises a usage limit by the text Claude Code writes, and a model's family by a name contained in a model name. It needs no list of models and no access to the harness, so it ships no model name (ADR-0401). In exchange, a new wording or a model name without its family stops the phase as today, never on the wrong model.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; a fallback entry carries the phase, the Spec slug, the repository path and the repository identity that Spec 0404 records in a phase entry. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added; neither command reads a remote, starts a harness or asks a judge. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the PRD's row stands, with ADR-0427 and ADR-0428 as the decisions this design implements and ADR-0401, ADR-0407, ADR-0414, ADR-0408, ADR-0392, ADR-0050, ADR-0256 and ADR-0349 as the ones it obeys or amends. This stage adds ADR-0429 (ADR-0429: "A refusal exits `3`"). ADR-0405, ADR-0410, ADR-0412 and ADR-0415 are listed as not applicable, as the PRD's row says. The phase command keeps its contract (ADR-0414: "records each phase as before and reads no earlier entry"). Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-17, quoted in the authorization record, which keeps the grant of every owned skill. The design changes no Makefile, workflow, Go module, lint or formatter file, no Baseline module and no file under `docs/agents/`, and no Task changes `.roundfixrc.yml`. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0405-authoring-model-fallback/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/phase.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## System Architecture

No package is new. Each subsection below is a component, and the Coverage Map names components by these titles.

### Fallback declaration

`internal/config/authoring.go` gains `AuthoringFallback` and a `Fallback` field on the entries of `idea`, `prd`, `techspec` and `tasks`. `validateAuthoringEntry` accepts the `fallback` key on those entries and refuses it on `research`. The installation's declaration (`authoringDeclaration` in `internal/cli/authoring_models.go`) ignores the field, so headers, subagent files and the three checks do not change. The configuration guide and the Roundfix skill's setup reference, which describe the declaration, describe the key.

### Limit reader

`internal/phaserecord/limit.go` holds `ReadUsageLimit`, a pure function that finds the limit line in an error and classifies it (ADR-0429). It holds no model name.

### Fallback entries

`internal/phaserecord/fallback.go` writes and reads the `fallback` kind of line in the day files. The writer is `Recorder.RecordFallback`, which appends like `StartPhase`. The reader is `Fallbacks`, which reads every day file and returns the fallback entries oldest first. `Phases` and the readers Spec 0404 adds keep skipping a kind they do not know.

### Phase subcommands

`internal/cli/phase.go` dispatches `start`, `model` and `fallback`. `model` and `fallback` resolve the Git root as `start` does, compute the repository identity as Spec 0404's phase entry does, and load Project Config with `loadCommandConfig`. They then decide from the entry of their phase. The phase usage and the root help in `internal/cli/cli.go` name both.

### Skill and guide texts

- **Authoring skills:** each of the four authoring skills' `## Declared model` section changes its "You activated this skill yourself" bullet (API Contract 6).
- **Roundfix skill:** its phase reference describes both commands.
- **User guide:** `docs/user-guide/commands/phase.md` describes both commands, `docs/user-guide/configuration.md` describes the key, and the Phase Record section of `docs/user-guide/usage.md` describes the entry.
- **Glossary and ADRs:** `CONTEXT.md` adds **Authoring Fallback** and revises three terms. ADR-0401 and ADR-0407 carry their amendment notices.

## Implementation Design

### Interfaces

```go
// internal/config
type AuthoringFallback struct {
	Model           string `yaml:"model"`
	ReasoningEffort string `yaml:"reasoning_effort"`
}
// AuthoringEntry gains: Fallback *AuthoringFallback `yaml:"fallback"`

// internal/phaserecord
type LimitKind string // "model", "shared", "other", "none"
type UsageLimit struct {
	Kind LimitKind
	Name string // as written, e.g. "Fable", "weekly", "monthly spend"
	Line string // the matched line, one line, at most 200 characters
}
func ReadUsageLimit(message string) UsageLimit
```

```go
// internal/phaserecord
type FallbackEntry struct {
	Phase                                       Phase
	Spec, Repository, RepositoryID              string
	Harness, Session                            string
	DeclaredModel, DeclaredEffort, Model, Effort string
	Limit, Message                              string
}
type RecordedFallback struct {
	Time time.Time
	FallbackEntry
}
func (r *Recorder) RecordFallback(ctx context.Context, entry FallbackEntry) error
func Fallbacks(homeDir string) ([]RecordedFallback, int, error)
```

Invariants:

```text
1. ReadUsageLimit takes the first line, in message order, that holds "hit your <name> limit" after "You've" or "You’ve", or that begins "<name> limit reached"; letter case is ignored.
2. A name of one word other than session or weekly is Kind model; session or weekly is Kind shared; a name of two or more words is Kind other; no such line is Kind none with empty Name and Line.
3. Line is the matched line with control characters replaced by spaces, trimmed, cut to 200 characters on a rune boundary.
4. A name matches a model when the model, in lower case, contains the name, in lower case.
5. A fallback entry is written only by phase fallback, only when it prints a model, and before it prints.
6. phase model writes nothing; phase fallback writes at most one line; phase start is unchanged.
7. A fallback entry belongs to a Spec in a repository when its spec equals the slug and its repository_id equals the current one, or, when either identity is empty, its repository path equals the current Git root.
8. Readers of phase entries skip a line whose kind is fallback without counting it as invalid.
```

### Data Models

**Project Config.** Under `skills.authoring.claude.<phase>`, for `idea`, `prd`, `techspec` and `tasks`:

```yaml
fallback:
  model: opus            # required; the authoring model name pattern
  reasoning_effort: high # optional; low, medium, high, xhigh or max
```

Validation, in order, each a configuration error that names the path:

- `fallback` on `research`: `skills.authoring.claude.research.fallback is not a supported key; expected model, reasoning_effort or tools`.
- Not a mapping: `<path>.fallback must be a mapping`.
- An unknown key: `<path>.fallback.<key> is not a supported key; expected model or reasoning_effort`.
- A malformed model or effort: the messages of the entry's own keys, with the `.fallback` path.
- No model: `<path>.fallback declares no model; set model`.
- No entry model: `<path>.fallback needs <path>.model`.
- The same model: `<path>.fallback.model "<m>" is the declared model; a fallback must name another model`, compared ignoring letter case.

The refusal of an unknown entry key now reads `expected model, reasoning_effort or fallback`.

**Fallback entry.** One line in the day file of its UTC date, under `roundfix/phase-record/v1`, with its keys in this order:

```json
{"schema":"roundfix/phase-record/v1","kind":"fallback","time":"2026-11-21T02:05:11Z","phase":"techspec","spec":"0412-example","repository":"/repo","repository_id":"sha256:5f0c…","harness":"claude","session":"3f9a1c20","declared_model":"fable","declared_effort":"xhigh","model":"opus","effort":"xhigh","limit":"Fable","message":"You've hit your Fable limit · resets Nov 22 at 8pm"}
```

`harness` and `session` come from the environment as for a phase entry, and are empty without one. An empty effort is the empty string. A line that does not decode, or holds a phase outside the four, counts as skipped.

### API Contracts

1. API Contract: `roundfix phase fallback <phase> --spec <slug> --error <message>`.
   - **Usage errors (exit `2`, nothing written).** The phase must be `idea`, `prd`, `techspec` or `tasks`, `--spec` must match the phase command's slug pattern, and `--error` must be non-empty and at most 4,096 bytes. The working directory must be inside a Git repository, and Project Config must load. A refusal prints `roundfix: phase fallback failed: <reason>` and `Run 'roundfix phase --help' for usage.` on standard error.
   - **Decision.** With the entry of the phase in the Authoring Model Declaration, the checks run in this order, and the first that fails gives the refusal line, exit `3`, with nothing written:
     1. `ReadUsageLimit` reads Kind `none`: `the error is not a usage limit`.
     2. Kind `shared`: `a <name> limit is shared by every model`.
     3. Kind `other`: `a <name> limit is not a model's usage limit`.
     4. No entry model: `<phase> declares no model`.
     5. The name does not match the declared model: `the <Name> limit is not that of the declared model <model>`.
     6. No fallback: `<phase> declares no fallback`.
     7. The fallback model matches the name: `the fallback <model> is a <Name> model too`.
   - **Refusal output.** The refusal is one line on standard output: `No fallback for <phase> of <slug>: <reason>.`
   - **Success.** Otherwise the command appends the fallback entry and prints one line, exit `0`: `Start roundfix-<phase>-author with model <m> and effort <e>, the fallback of <declared>: "<Line>"`. Without a fallback effort, `and effort <e>` reads `and no effort`. A write failure prints `roundfix: phase fallback failed: <error>` on standard error, exit `1`.
2. API Contract: `roundfix phase model <phase> --spec <slug>`.
   - **Usage errors.** The same as API Contract 1, without `--error`.
   - **No declared model.** When the entry of the phase declares no model, the command prints `No Author Subagent model for <phase>: the Authoring Model Declaration declares no model for it.`, exit `0`.
   - **Reading the record.** Otherwise it reads `Fallbacks`. An unreadable record directory prints `roundfix: phase model failed: <error>`, exit `1`; skipped lines are ignored.
   - **On the fallback.** It names the fallback when all three hold:
     - an entry of this Spec in this repository (invariant 7) has a `limit` that matches this phase's declared model (invariant 4);
     - this phase declares a fallback;
     - the fallback's model does not match that `limit`.

     The line is `Start roundfix-<phase>-author with model <m> and effort <e>, the fallback of <declared>: <phase0> of <slug> fell back when the <limit> limit ran out.`, where `<phase0>` is the phase of the oldest such entry.
   - **On the declared model.** Otherwise the line is `Start roundfix-<phase>-author with no model and no effort: the declared model <declared>.`, exit `0`.
3. API Contract: `roundfix phase --help` and the root help.
   - **Phase usage.** The usage of `roundfix phase` adds the two usage lines, two command lines (`model     Print the model and effort a phase's Author Subagent starts on` and `fallback  Decide and record an Authoring Fallback from the error an Author Subagent ended on`), the phase note `model and fallback take idea, prd, techspec or tasks`, and the option `--error`.
   - **Root help.** The root help adds the same two usage lines after the `phase start` line and reads `phase      Record a Supervisor Phase, or decide an Authoring Fallback`.
   - **Spec 0404's lines.** The lines Spec 0404 adds for `phase start` stay as it leaves them.
4. API Contract: the Phase Record.
   - The fallback entry of Data Models joins the day files and their retention.
   - `phaserecord.Phases` and every reader of phase entries skip it (invariant 8).
   - The GC count of day files is unchanged, because it counts files, not lines.
5. API Contract: Project Config, as Data Models says. A Roundfix older than the release refuses a Project Config that holds the key, as ADR-0408 records ("refuses a key it does not know"); the configuration guide says so.
6. API Contract: the authoring skills. In each of `write-idea`, `write-prd`, `write-techspec` and `write-tasks`, the bullet that starts the Author Subagent says, for its own phase:
   - **Before starting.** Run `roundfix phase model <phase> --spec <slug>`. Start the subagent with the model and effort it prints, or passing no model and no effort when it prints the declared model or no Author Subagent model.
   - **On an error.** When the subagent ends on an error, run `roundfix phase fallback <phase> --spec <slug> --error "<the line of the error that names the limit>"`.
   - **When it prints a model.** Tell the user which phase switched, from which model to which, and the limit line. Without waiting, start the same subagent with that model and effort, the same brief, the output it returned and the files it wrote, and say in the brief that it runs on the fallback.
   - **When it prints no fallback, or the subagent cannot start.** Report the error and author here only if the user says so.
   - **An older Roundfix.** When the installed Roundfix does not know the commands, start passing no model and no effort and, on an error, report it as before.
   - **Unchanged.** The rest of each section, and the fix or refactor route of `write-techspec` and `write-tasks`, stay as they are.
7. API Contract: the Roundfix skill's phase reference gains an `Authoring Fallback` section with the two commands, their exit codes and the rule of API Contract 6. The Roundfix skill names `roundfix phase model` and `roundfix phase fallback` where it names the phase command.
8. API Contract: the glossary and the ADRs.
   - **Glossary.** `CONTEXT.md` adds **Authoring Fallback** and revises **Authoring Model Declaration** (it may declare a fallback per phase), **Author Subagent** (it is started with the printed model and effort when its phase falls back) and **Phase Record** (it holds fallback entries).
   - **ADR-0401** carries `Amended by ADR-0427 (2026-11-17)` with what changed.
   - **ADR-0407** carries `Amended by ADR-0428 (2026-11-17)` with what changed.

### Surface Transcripts

The transcripts run in a Git repository whose Project Config declares `idea`, `prd` and `techspec` as `fable`/`xhigh` with the fallback `opus`/`xhigh`, and `tasks` as `opus`/`high` with the fallback `sonnet`/`high`. `CLAUDE_CODE_SESSION_ID` is set, and the home is temporary. Transcripts 3 to 5 run after transcript 2, in the same home.

1. Surface Transcript: the phase usage names both subcommands.

   ```transcript
   $ roundfix phase --help
   stdout:
   Usage:
   ...
     roundfix phase model <phase> --spec <slug>
     roundfix phase fallback <phase> --spec <slug> --error <message>
   ...
     model     Print the model and effort a phase's Author Subagent starts on
     fallback  Decide and record an Authoring Fallback from the error an Author Subagent ended on
   ...
   stderr:
   exit: 0
   ```

2. Surface Transcript: a model's own limit falls back and is recorded.

   ```transcript
   $ roundfix phase fallback idea --spec 0412-example --error "You've hit your Fable limit · resets Nov 22 at 8pm"
   stdout:
   Start roundfix-idea-author with model opus and effort xhigh, the fallback of fable: "You've hit your Fable limit · resets Nov 22 at 8pm"
   stderr:
   exit: 0
   ```

3. Surface Transcript: a later phase on the limited model starts on its fallback.

   ```transcript
   $ roundfix phase model prd --spec 0412-example
   stdout:
   Start roundfix-prd-author with model opus and effort xhigh, the fallback of fable: idea of 0412-example fell back when the Fable limit ran out.
   stderr:
   exit: 0
   ```

4. Surface Transcript: a later phase on another model starts as declared.

   ```transcript
   $ roundfix phase model tasks --spec 0412-example
   stdout:
   Start roundfix-tasks-author with no model and no effort: the declared model opus.
   stderr:
   exit: 0
   ```

5. Surface Transcript: the next Spec starts on the declared model.

   ```transcript
   $ roundfix phase model prd --spec 0413-another
   stdout:
   Start roundfix-prd-author with no model and no effort: the declared model fable.
   stderr:
   exit: 0
   ```

6. Surface Transcript: a weekly limit stops the phase.

   ```transcript
   $ roundfix phase fallback techspec --spec 0412-example --error "You've hit your weekly limit · resets Mon 12:00am"
   stdout:
   No fallback for techspec of 0412-example: a weekly limit is shared by every model.
   stderr:
   exit: 3
   ```

7. Surface Transcript: a limit that is not the declared model's stops the phase.

   ```transcript
   $ roundfix phase fallback tasks --spec 0412-example --error "You've hit your Fable limit · resets Nov 22 at 8pm"
   stdout:
   No fallback for tasks of 0412-example: the Fable limit is not that of the declared model opus.
   stderr:
   exit: 3
   ```

8. Surface Transcript: a phase outside the four is a usage error.

   ```transcript
   $ roundfix phase model delivery --spec 0412-example
   stdout:
   stderr:
   roundfix: phase model failed: phase "delivery" has no Author Subagent; use idea, prd, techspec or tasks
   Run 'roundfix phase --help' for usage.
   exit: 2
   ```

## Coverage Map

- Goal 1 → Phase subcommands (API Contracts 1 and 2), Skill and guide texts (API Contract 6).
- Goal 2 → Phase subcommands: the success line of API Contract 1 names the phase's Author Subagent, the declared and fallback models and the limit line; Skill and guide texts: API Contract 6 has the session tell the user that line's content before it restarts, every time.
- Goal 3 → Phase subcommands (API Contract 2), Fallback entries.
- Goal 4 → Fallback entries (API Contract 4).
- Goal 5 → Fallback declaration (API Contract 5), Skill and guide texts: the declared line of API Contract 2.
- Story 1 → Fallback declaration.
- Story 2 → Limit reader, Phase subcommands (API Contract 1).
- Story 3 → Skill and guide texts (API Contract 6).
- Stories 4 and 5 → Phase subcommands (API Contract 2), Fallback entries.
- Story 6 → Limit reader, Phase subcommands: the refusals of API Contract 1.
- Story 7 → Fallback declaration, Phase subcommands: the declared and no-model lines of API Contract 2.
- Core Feature 1 → Fallback declaration (API Contract 5); Skill and guide texts, the configuration guide.
- Core Feature 2 → Limit reader, Phase subcommands (API Contract 1), Fallback entries.
- Core Feature 3 → Skill and guide texts (API Contract 6).
- Core Feature 4 → Phase subcommands (API Contract 2), Skill and guide texts (API Contract 6).
- Core Feature 5 → Fallback entries (API Contract 4).
- Core Feature 6 → Phase subcommands: refusal 7 of API Contract 1 for a fallback's own limit; Skill and guide texts, the typed-skill sentence of the phase reference.
- Core Feature 7 → Skill and guide texts (API Contracts 7 and 8).
- Success Metrics 1 and 2 → Phase subcommands, Fallback entries, Skill and guide texts.
- Success Metric 3 → Phase subcommands (API Contract 2, Surface Transcript 5).
- Success Metric 4 → Fallback declaration, Skill and guide texts: the declared line and the older-Roundfix sentence.

Deviations from the PRD:

- Core Feature 6 — a fallback's own limit is refused because the fallback model then matches the limit's name (refusal 5, the limit is not that of the declared model). The command tells the case apart by name, and keeps no count of earlier switches.
- Core Feature 5 — the restarted Author Subagent's phase entry carries the model it runs on only once Spec 0404 is delivered. This Spec requires 0404, so it adds no flag of its own to `roundfix phase start`.

## Absence Evidence

1. The `fallback` key of the declaration — command: `grep -rn "fallback" internal/config/authoring.go`; result: no match.
2. The limit reader — command: `grep -rn "hit your\|limit reached\|UsageLimit" internal --include=*.go`; result: no match outside tests of other packages.
3. The fallback entry kind — command: `grep -rn '"fallback"' internal/phaserecord internal/cli/phase.go`; result: no match; `Phases` skips every kind but `phase`.
4. The two subcommands — command: `grep -n '"model"\|"fallback"' internal/cli/phase.go`; result: no match; `parsePhaseStart` refuses every subcommand but `start`.
5. The glossary term — command: `grep -c "Authoring Fallback" CONTEXT.md`; result: 0.

## ADR Alternatives

- ADR-0427: the design is its decision; the harness's fallback chain, frontmatter fields, a Baseline rule, generated fallback subagents, a question before switching and a chain are the options it rejected.
- ADR-0428: the design is its decision; retrying the declared model, moving every later phase, the reset time and a session-only memory are the options it rejected.
- ADR-0429: the design is its decision; skill prose, flags on `phase start`, a shipped list of families, JSON output and exit `1` for a refusal are the options it rejected.
- ADR-0401: amended by ADR-0427 so the session passes the fallback's model and effort; it still rejects shipped model names and keeps both mechanisms, and no option it rejected is taken.
- ADR-0407: amended by ADR-0428 with a new entry kind; the record keeps no prompt text, token count or price, and no line is rewritten.
- ADR-0414: the phase command still records each phase as before and reads no earlier entry; `phase model` reads entries in a subcommand of its own.
- ADR-0408: the key reaches this repository's Project Config only after a release, in the operator's step; no option it rejected is taken.
- ADR-0392: an older binary refuses a Project Config that holds the key, as it does every unknown key; the guide says so.
- ADR-0405: not applicable — a user-level skill copy is reported as before.
- ADR-0410: not applicable — the status view keeps reading phase entries only, and a fallback entry is not one.
- ADR-0412: not applicable — planning references and the release preparation keep their release impact.
- ADR-0415: not applicable — the recommended declaration of the Doctor Command and the Update Check gains no fallback in this Spec.
- ADR-0050: not applicable — the fallbacks of a Run, and their rule that work already begun is not switched, are untouched.
- ADR-0256: the skills and the phase reference describe both commands in the same change.
- ADR-0349: the authorization record carries the Release Ceiling `minor` for a new key and two command paths.

## Integration Points

- Claude Code, through text only. The session copies the limit line of the error into `--error`, and passes the printed model and effort when it starts the Author Subagent. Roundfix starts no harness and reads no transcript.
- Project Config, through `loadCommandConfig`.
- The Phase Record, through `internal/phaserecord`.
- The repository identity, through the function Spec 0404 exports from `internal/baseline`.

## Testing Approach

- **Existing seams.** The tests use the seams that already exist:
  - the command entry `runCLI` with a temporary home and a temporary Git repository holding a `.roundfixrc.yml`;
  - the injected environment lookup for `CLAUDE_CODE_SESSION_ID`;
  - the injected clock `phaseNow`;
  - table tests on `internal/config` and `internal/phaserecord`.
- **No new seam**, and every new test calls `t.Parallel()`.
- **Unit tests.**
  - `ReadUsageLimit` over the messages of the harness reference, both apostrophes, letter case, multi-line errors, control characters and a 300-character line.
  - The validation of every refusal in Data Models.
  - `RecordFallback`, `Fallbacks` and the skip of the new kind by `Phases`.
- **Command tests.**
  - Each refusal of API Contract 1 writes no file.
  - The success line and the entry are written.
  - `phase model` follows invariants 4 and 7, across two Specs and two repositories.
  - Surface Transcripts 1 to 8 are asserted line by line.
- **Skill tests.** A test under `skills/` asserts the paragraph of API Contract 6 in the four authoring skills and the section of API Contract 7.
- **Tests this change moves**, measured in the prototype of the references (`references/2026-11-17-harness-and-test-measurements.md`), which the Tasks declare:
  - `TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow/phase_field` in `./internal/config`;
  - `TestPhaseStartSurfaceTranscripts/transcript_1` in `./internal/cli`, and, once Spec 0404 is delivered, its test of the full `roundfix phase --help`;
  - in `./internal/docscontract`, `TestEveryCommandIsNamedInTheRoundfixSkill`, `TestEveryCommandIsNamedInTheUserGuide`, `TestAnUndocumentedCommandIsReported` and `TestTheSkillCoverageMapIsCurrent`, which the documentation and the regeneration satisfy. `TestAnUndocumentedCommandIsReported` keeps expecting `reopen` only.
- **No test** depends on a release tag, a ref a fresh checkout lacks, the network, a real harness or the real home.

## Build Order

1. Fallback declaration: the key, its validation and the configuration guide.
2. Limit reader, Fallback entries and Phase subcommands: `phase fallback`, `phase model`, the helps, the command guide, the Phase Record section, and the coverage and surface records (depends on: 1).
3. Skill and guide texts: the four authoring skills, the Roundfix skill and its phase reference, the glossary and the ADR notices, with the regeneration (depends on: 2).
4. The QA gate (depends on: 1, 2, 3).

## Risks & Considerations

- **The text of a limit can change.** A new wording reads as "not a usage limit" and stops the phase as today. The limit reader's tests pin the wordings of Claude Code 2.1.326.
- **The fallback must be an alias.** The Agent tool takes an alias for the model at start; a full model identifier is refused by the start, and the skill then reports the error.
- **Untrusted text.** The error is text from the harness. Only the matched line is kept, as one line of at most 200 characters, and it is printed in quotes. The skill text says it is data.
- **Two machines.** A Spec authored on two machines falls back only where the record holds the entry (ADR-0428).
- **Order of Specs.** Specs 0401 to 0404 change the same help, the same phase command and the same skill sections first. Each Task runs its packages again on its own base.

## Glossary

None.

## Decisions

- Two subcommands of `roundfix phase`, a limit read from text, and exit `3` for a refusal. See ADR-0429.
- No JSON output: the session reads one line, and the exit code carries the decision (ADR-0429).
- `phase model` prints "no Author Subagent model" rather than failing when a phase declares no model, so a skill can always run it before starting.
- A fallback entry stays in the record when the restart fails, so the Spec keeps its fallback (ADR-0429, Consequences).
- The installation ignores the new key: no header, subagent file or check changes.
