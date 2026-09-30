---
spec: 0205-an-advisory-judge-for-spec-authoring
prd: _prd.md
created: 2026-09-30
---

# An advisory judge for Spec authoring — Technical Spec

## Executive Summary

A new package, `internal/judge`, plans the judgments for one Spec from its
PRD, its TechSpec and the ADRs they cite, asks each one of the TypeSafe
evaluation endpoint, records every call in a monthly Judge Log, and compares
each answer with the measured threshold. A new subcommand,
`roundfix spec judge`, prints what it raises. The `write-prd` and
`write-techspec` skills run it after the Spec Consistency Check.

The trade-off this design accepts is fidelity to the measurement over
elegance. The claim extraction, the section selection and every truncation
limit reproduce the 2026-09-30 benchmark, including its own attribution verb
list, rather than reuse the Spec Consistency Check's parser. A threshold holds
only for inputs shaped like the ones it was measured on, and the two parsers
disagree on which sentences are claims. The cost is a second, small parser
whose rules live in the question file beside the thresholds they serve.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; a request is
  named by the SHA-256 of its state and a Judge Log file by its UTC month.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — one HTTPS `POST` per judgment to the
  endpoint in the question file, with `Authorization: Bearer <key>` and
  `Content-Type: application/json` and no other header. The key comes only
  from `TYPESAFE_API_KEY` in the command's environment and is written nowhere.
  The transport honors the standard proxy environment variables. Request and
  response fields are those of the published API reference. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0200 and ADR-0201 (this Spec),
  ADR-0116, ADR-0183, ADR-0176, ADR-0035, ADR-0089, ADR-0187, ADR-0189,
  ADR-0193, ADR-0081, ADR-0149, ADR-0182, ADR-0178, ADR-0166 and ADR-0167 hold
  as the PRD states. ADR-0184 has this TechSpec state the new command as
  Surface Transcripts. The gate is bound by ADR-0080, ADR-0088, ADR-0091,
  ADR-0096, ADR-0104, ADR-0117, ADR-0155 and ADR-0156, and by ADR-0093 and
  ADR-0094. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the skill files, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/write-prd/SKILL.md`, `skills/write-prd/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`, `skills/write-techspec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`. Sanctioned regeneration:
  `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

One new package, `internal/judge`. It is the first Roundfix code that makes a
network request, and it imports nothing from `internal/speccheck`, so the
Spec Consistency Check keeps no network dependency.

| Concern | Owner | File |
| --- | --- | --- |
| Questions, thresholds, rules, model, price, ceiling | `questions.json`, `Load` | `internal/judge/questions.json` (new), `internal/judge/questions.go` (new) |
| The only reader of Spec artifacts | `Source`, `readSpecArtifact`, `readADR` | `internal/judge/source.go` (new) |
| The language gate | `isEnglish` | `internal/judge/language.go` (new) |
| Claims and goal pairs | `PlanSpec` | `internal/judge/pairs.go` (new) |
| One request and its answer | `client.ask` | `internal/judge/client.go` (new) |
| The Judge Log and the month's spend | `judgeLog` | `internal/judge/log.go` (new) |
| Asking a plan, thresholds, stopping | `Run` | `internal/judge/judge.go` (new) |
| The command | `runSpecJudgeCommand` | `internal/cli/spec_judge.go` (new), `internal/cli/spec_check.go`, `internal/cli/cli.go` |

## Implementation Design

### Questions and thresholds

`internal/judge/questions.json` is embedded in the binary and is the one
reviewable source of every question, threshold and rule below. Its bytes are
exactly the block that follows. The question texts, criteria and thresholds
are the measurement's; the extraction patterns and limits are the ones its
datasets were built with.

```json
{
  "schema": "roundfix/judge-questions/v1",
  "model": "jev-1.13.0",
  "endpoint": "https://api.typesafe.ai/v1/systemone",
  "usd_per_million_input_tokens": 0.042,
  "monthly_ceiling_usd": 5.0,
  "language": {
    "english_words": ["the", "of", "and", "to", "in", "is", "that", "it", "for", "with", "are", "be", "by", "this", "which", "when", "not", "from", "or", "an", "on", "as", "at", "was", "were", "has", "have", "will", "must"],
    "portuguese_words": ["de", "que", "não", "uma", "para", "com", "os", "das", "dos", "em", "é", "se", "por", "mais", "na", "no", "ao", "pela", "pelo", "como", "mas", "são", "está", "foi", "ser", "também", "sua", "seu"],
    "min_english_share": 0.08,
    "max_portuguese_share": 0.04
  },
  "judgments": {
    "citation-support": {
      "question_id": "relation",
      "question": {
        "type": "choice",
        "instructions": "`claim` is a sentence from a specification. It says what 'the cited decision' does, keeps or requires. `decision_record` is the full text of one architecture decision record. How does `decision_record` relate to what `claim` attributes to the cited decision?",
        "criteria": {
          "supports": "The decision record states what the claim attributes to the cited decision, or directly implies it",
          "contradicts": "The decision record states the opposite of what the claim attributes to the cited decision",
          "says_nothing": "The decision record is about something else or does not address what the claim attributes to the cited decision"
        }
      },
      "raise_when_choice_is_not": "supports",
      "raise_min_confidence": 0.8,
      "attribution_pattern": "\\bADR-(\\d{4})\\s+(?:already\\s+)?(?:keeps?|makes?|applies|apply|owns?|checks?|requires?|accepts?|places?|has|preserves?|gives?|surfaces?|separates?|sanctions?|lets?|allows?|governs?|protects?|establish(?:es|ed)?|puts?|says?|records?|defines?|forbids?|binds?|treats?|moves?|reads?|runs?|refuses?|limits?|sets?)\\b",
      "claim_min_chars": 60,
      "claim_max_chars": 500,
      "decision_record_max_chars": 7000
    },
    "goal-mechanism": {
      "question_id": "delivers_goal",
      "question": {
        "type": "noul",
        "instructions": "`goal` is one goal of a product requirements document. `section` is one section of the technical specification written for it, titled `section_title`. Does `section` describe a mechanism that would deliver `goal`?",
        "criteria": {
          "true": "The section describes concrete behavior, components or steps whose effect is the goal",
          "false": "The section is about a different goal, or only restates the goal without a mechanism"
        }
      },
      "raise_when_noul_below": 0.3,
      "coverage_line_pattern": "^(?:- |\\| )(?:PRD )?Goals? (\\d+)(?:\\s*[–-]\\s*(\\d+))?[^→|>]*(?:→|->|\\|)\\s*(.+)$",
      "generic_section_title_pattern": "(?i)^(coverage map|api contracts?|testing approach|build order|project constraints|executive summary|risks?( &| and)? considerations?|risks?|integration points|data models?|references?|research basis|system architecture|implementation design|open questions?|decisions?|overview|vocabulary contract)$",
      "section_title_min_chars": 8,
      "section_min_chars": 200,
      "section_max_chars": 3500
    }
  }
}
```

A change to any value in this file is a change to the judge's contract. A new
model version needs a new measurement before `model` or a threshold moves.

### Interfaces

```go
// internal/judge/questions.go
//go:embed questions.json
var questionFile []byte

type Questions struct {
	Model, Endpoint                               string
	USDPerMillionInputTokens, MonthlyCeilingUSD   float64
	Language                                      LanguageGate
	Citation                                      CitationJudgment // "citation-support"
	Goal                                          GoalJudgment     // "goal-mechanism"
}

// Load parses the embedded file and compiles its patterns; it is the
// package's only source of texts, thresholds and limits.
func Load() (Questions, error)
```

```go
// internal/judge/source.go
// Source is text the judge may send. Its fields are unexported, and only
// the two readers below build one.
type Source struct{ path, text string }

// readSpecArtifact reads "_prd.md" or "_techspec.md" of one Spec directory.
func readSpecArtifact(specDir, name string) (Source, error)

// readADR reads docs/adr/<number>-*.md; ok is false unless it is accepted.
func readADR(repoRoot, number string) (src Source, ok bool, err error)
```

```go
// internal/judge/pairs.go and judge.go
type Stage string // "" (both), "prd", "techspec"

func PlanSpec(q Questions, repoRoot, specDir string, stage Stage) (Plan, error)

type Request struct {
	RepoRoot, SpecDir, Spec string
	Stage                   Stage
	Key                     string // TYPESAFE_API_KEY; "" asks nothing
	HomeDir                 string // Judge Log: <HomeDir>/.roundfix/judge
	Transport               http.RoundTripper
	Now                     func() time.Time
}

func Run(ctx context.Context, q Questions, req Request) (Report, error)
```

`Plan` holds the pending judgments, the judgments skipped while planning and
the skipped artifacts. A pending judgment carries its kind, artifact, line,
target (`ADR-NNNN` or `Goal N → <section title>`), text and an unexported
state, which only the planning code builds, and only from `Source` values.
`Report` holds every judgment with its outcome (`advisory`, `clear` or
`skipped`), reason, answer and answering model, plus the calls, input tokens,
cost, the month's cost, the ceiling, and the run-level `Skipped` or `Stopped`
reason. `Run` returns an error only when the Spec's PRD cannot be read.

### The only readers

- `readSpecArtifact` accepts only `_prd.md` and `_techspec.md`, directly in
  the Spec directory, as a regular file (`Lstat`, never following a symbolic
  link) of at most 1 MiB.
- `readADR` accepts only a regular file directly in `<repo>/docs/adr` whose
  name starts with the four-digit number and a hyphen and ends in `.md`, of at
  most 1 MiB. It is judged only when it is an accepted ADR in this
  repository's sense (`docs/agents/docs-layout.md`): its front matter status
  is `accepted`, or it is a legacy ADR with no lifecycle front matter whose
  body does not mark it inactive. An ADR whose status is `proposed`,
  `rejected`, `deprecated` or `superseded` is never read, so no draft or
  retired decision is ever sent.
- A refused PRD or TechSpec is reported as a skipped artifact with reason
  `not a regular file in the Spec directory`; a refused ADR skips the claim.
- The request code accepts only a pending judgment, and a pending judgment's
  state holds only strings cut from `Source` text plus the fixed phrase
  `the cited decision`.

### Citation claims

For `_prd.md` (stage `prd` or both) and `_techspec.md` (stage `techspec` or
both), after the front matter is removed:

1. A line starting with three backticks toggles a fence and ends the
   paragraph; lines inside a fence are skipped. A blank line, or a line whose
   trimmed text starts with `#` or `|`, ends the paragraph and is skipped. A
   line whose trimmed text starts with `- ` or digits and `. ` ends the
   paragraph and starts a new one without that marker.
2. A paragraph's lines are joined with single spaces. It is split into
   sentences after `.`, `!` or `?` followed by whitespace when the next
   character is an ASCII capital, a backtick, `(` or `[`.
3. A sentence is a claim when it holds exactly one `ADR-` token followed by
   four digits, `attribution_pattern` matches it, and `readADR` accepts that
   number.
4. The claim text replaces the token with `the cited decision`, removes every
   `(Spec NNNN)`, collapses whitespace and capitalizes its first character. A
   claim outside `claim_min_chars` to `claim_max_chars` is skipped with reason
   `claim length outside the measured 60-500 characters`.
5. The state is `{"claim": <claim>, "decision_record": <ADR text after its
   front matter, trimmed, cut at decision_record_max_chars>}`.
6. The line is the first line at or after the paragraph's start that holds
   the claim's ADR token.

Characters are Unicode code points. An ADR that fails the language gate
skips the claim with reason `decision record is not English`.

### Goal pairs

For stage `techspec` or both, when `_techspec.md` exists:

1. Goal N is the Nth item of the PRD's `## Goals` section: a line starting
   with `- ` or digits and `. ` starts an item, following non-blank lines
   continue it, and whitespace is collapsed.
2. In the TechSpec's `## Coverage Map` section, a line that starts with
   whitespace continues the previous line. Each resulting line, trimmed, that
   matches `coverage_line_pattern` names goals `lo` to `hi` and a target text.
3. A section is each `## ` or `### ` heading with its body up to the next
   heading of the same or a higher level. Its normalized title has backticks
   and asterisks removed, whitespace collapsed and letters lowered. It
   qualifies when that title does not match `generic_section_title_pattern`
   and its trimmed body holds at least `section_min_chars` characters.
4. A qualifying section is named by the line when its normalized title holds
   at least `section_title_min_chars` characters and occurs in the normalized
   target text. The longest title wins.
5. The state is `{"goal": <goal>, "section_title": <title as written>,
   "section": <trimmed body cut at section_max_chars>}`, one per goal in the
   range, at the Coverage Map line's number.
6. A goal outside the PRD's list is skipped with reason `Goal N is not in the
   PRD`; a line with no named qualifying section, with reason `no named section
   qualifies`.

### The language gate

`isEnglish(text)` splits the text after its front matter into maximal runs
of Unicode letters, lowercased. It is true when the share of words in
`english_words` is at least `min_english_share` and the share in
`portuguese_words` is below `max_portuguese_share`. A PRD or TechSpec that
fails is a skipped artifact with reason `not English`, and nothing from it is
planned. The goals come from the PRD, so a non-English PRD plans no goal pair.

### Asking

For each pending judgment in plan order, with identical states asked once:

1. **Before any request.** No key: the run is skipped with
   `TYPESAFE_API_KEY is not set`. An unreadable Judge Log: skipped with
   `judge log unreadable: <error>`. The month's cost at or above the ceiling:
   skipped with `monthly ceiling reached (US$<spent> of US$5.00)`. A skipped
   run sends nothing.
2. **Before each request.** The month's cost plus this run's cost at or above
   the ceiling stops the run with the same ceiling reason.
3. **The request.** `POST` to `endpoint` with body
   `{"state": <state>, "model": <model>, "questions": {<question_id>: <question>}}`,
   the state encoded without HTML escaping. Each attempt has 30 seconds. A
   `429` or `529` is retried at most twice, after the `Retry-After` seconds
   when present, capped at 10, or else after 1.5 and 3 seconds.
4. **The answer.** A `200` whose `model` differs from the pinned `model` is
   skipped with `answered by <model>, thresholds belong to jev-1.13.0`. A
   `200` without a well-typed answer under the question ID is skipped with
   `unreadable answer`. Otherwise citation support is `advisory` when `choice`
   is not `supports` and `confidence` is at least `raise_min_confidence`, and
   goal to mechanism is `advisory` when `noul` is below
   `raise_when_noul_below`; any other answer is `clear`.
5. **Failures.** A `401` or `403` stops the run with `key refused (HTTP <n>)`.
   Any other `4xx` skips that judgment with `request refused (HTTP <n>)` and
   continues. A `5xx`, a `429` or `529` after its retries, a timeout or a
   network error stops the run with `service unavailable (<detail>)`.
6. **The record.** Every request appends one Judge Log line before the next
   request. A failed append stops the run with
   `judge log not writable: <error>`, after the answer in hand is reported.

A stopped run counts every judgment not asked as skipped. No failure changes
the exit code.

### Data Models

The Judge Log is `<home>/.roundfix/judge/<YYYY-MM>.jsonl`, named by the UTC
month of the call, created `0700` for the directory and `0600` for the file,
opened for append. One line per request:

```json
{"schema":"roundfix/judge-log/v1","time":"2026-10-01T12:00:00Z","repository":"/work/repo","spec":"0300-example","judgment":"citation-support","artifact":"docs/specs/0300-example/_techspec.md","line":31,"target":"ADR-0035","state_hash":"9fd1a51820e9a8e9","question_id":"relation","requested_model":"jev-1.13.0","model":"jev-1.13.0","answer":"says_nothing","probabilities":{"supports":0.02,"contradicts":0.05,"says_nothing":0.93},"confidence":0.93,"noul":null,"latency_ms":348,"input_tokens":842,"output_tokens":20,"cost_usd":0.000035364,"status":200,"attempts":1,"error":"","outcome":"advisory"}
```

- `state_hash` is the first 16 hex digits of the SHA-256 of the encoded
  state. `cost_usd` is `input_tokens` times `usd_per_million_input_tokens`
  divided by one million; a request with no usage costs `0`.
- A goal line has `"answer":null`, `"probabilities":null`,
  `"confidence":null` and a `noul`. A failed request has `"model":""`, a zero
  or HTTP `status`, its `error` and `"outcome":"skipped"`.
- The month's cost is the sum of `cost_usd` over the current month's file. A
  missing file is `0`; a line that does not parse makes the log unreadable.
- The key is in no field, and the repository path is the Git root.

### Surface Transcripts

The fixture for Transcripts 1 to 5 is a repository whose Spec
`0300-example` is English, whose `docs/adr/0035-*.md` is an accepted ADR
about the Spec Root, and whose TechSpec attributes a Run Event retention rule
to it. The fake transport answers each request with `input_tokens` 842.

1. Surface Transcript: two raised judgments at the TechSpec stage.

   ```transcript
   $ roundfix spec judge 0300-example --stage techspec
   stdout:
   advisory citation-support docs/specs/0300-example/_techspec.md:31 ADR-0035 says_nothing at confidence 0.93: The cited decision keeps every Run Event for ninety days after the Run ends.
   advisory goal-mechanism docs/specs/0300-example/_techspec.md:74 Goal 2 → The widget cache: P(delivers) 0.12
   Judge: 2 advisory, 3 clear, 0 skipped; 5 call(s), 4210 input tokens, US$0.0002; month US$0.0002 of US$5.00; model jev-1.13.0
   stderr:
   exit: 0
   ```

2. Surface Transcript: no key in the environment.

   ```transcript
   $ roundfix spec judge 0300-example
   stdout:
   Judge: skipped: TYPESAFE_API_KEY is not set; 5 judgment(s) not asked
   stderr:
   exit: 0
   ```

3. Surface Transcript: the month's Judge Log already holds US$5.0003.

   ```transcript
   $ roundfix spec judge 0300-example
   stdout:
   Judge: skipped: monthly ceiling reached (US$5.0003 of US$5.00); 5 judgment(s) not asked
   stderr:
   exit: 0
   ```

4. Surface Transcript: a Spec written in Portuguese.

   ```transcript
   $ roundfix spec judge 0301-exemplo --stage prd
   stdout:
   skipped docs/specs/0301-exemplo/_prd.md: not English
   Judge: 0 advisory, 0 clear, 0 skipped; 0 call(s), 0 input tokens, US$0.0000; month US$0.0000 of US$5.00; model jev-1.13.0
   stderr:
   exit: 0
   ```

5. Surface Transcript: the service fails on the second request.

   ```transcript
   $ roundfix spec judge 0300-example --stage techspec
   stdout:
   advisory citation-support docs/specs/0300-example/_techspec.md:31 ADR-0035 says_nothing at confidence 0.93: The cited decision keeps every Run Event for ninety days after the Run ends.
   Judge: 1 advisory, 0 clear, 4 skipped; 2 call(s), 842 input tokens, US$0.0000; month US$0.0000 of US$5.00; model jev-1.13.0; stopped: service unavailable (HTTP 503)
   stderr:
   exit: 0
   ```

6. Surface Transcript: an unknown Spec.

   ```transcript
   $ roundfix spec judge 0999-missing
   stdout:
   stderr:
   roundfix: spec judge failed: unknown active Spec slug "0999-missing"
   Run 'roundfix spec judge --help' for usage.
   exit: 2
   ```

A skipped judgment prints
`skipped <judgment> <artifact>:<line> <target>: <reason>`. A clear judgment
prints nothing in text form.

### API Contracts

1. API Contract: `roundfix spec judge <slug> [--stage <prd|techspec>] [--format <text|json>]`
   — judges one active Spec in the resolved Spec Root, Surface Transcripts 1
   to 6. Exit `0` whenever it ran; exit `2` for an unknown flag, an unknown
   stage or format, an unknown active Spec, a missing `_prd.md`, or
   `--stage techspec` without `_techspec.md`. It writes only the Judge Log.
2. API Contract: `--format json` prints one document with schema
   `roundfix/spec-judge/v1` and the fields `spec`, `model`, `skipped`,
   `stopped` (each a reason or null), `artifacts_skipped` (each `artifact` and
   `reason`), `judgments` (each `kind`, `artifact`, `line`, `target`, `text`,
   `section_title`, `outcome`, `reason`, `answer`, `probabilities`,
   `confidence`, `noul` and `model`, with null where not applicable), `calls`,
   `input_tokens`, `cost_usd`, `month_cost_usd` and `month_ceiling_usd`. It
   lists clear judgments too.
3. API Contract: the TypeSafe request and response — the body and headers of
   Asking step 3, read back as the API reference defines: `model`,
   `answers.<question_id>` with `type` and `choice`, `probabilities` and
   `confidence`, or `noul`, and `usage.input_tokens` and
   `usage.output_tokens`.
4. API Contract: the Judge Log line of Data Models, schema
   `roundfix/judge-log/v1`, one per request.
5. API Contract: help — `roundfix spec --help` and the top-level help name
   `roundfix spec judge <slug> [--stage <prd|techspec>] [--format <text|json>]`,
   and `roundfix spec judge --help` names `TYPESAFE_API_KEY`, the monthly
   ceiling, the Judge Log and that the command never fails for a judgment.

## Coverage Map

- Goal 1 → Citation claims; Goal pairs; Asking; API Contract 1.
- Goal 2 → Asking; API Contract 1.
- Goal 3 → The only readers; The language gate; API Contract 3.
- Goal 4 → Data Models; Asking; API Contract 4.
- Goal 5 → Questions and thresholds.
- User Story 1 → Citation claims; Surface Transcript 1.
- User Story 2 → Goal pairs; Surface Transcript 1.
- User Story 3 → Asking; Surface Transcripts 2, 3, 4 and 5.
- User Story 4 → Data Models; API Contract 4.
- User Story 5 → The only readers; Testing Approach 3.
- Core Feature 1 → API Contracts 1 and 2.
- Core Feature 2 → Citation claims; Asking.
- Core Feature 3 → Goal pairs; Asking.
- Core Feature 4 → Questions and thresholds.
- Core Feature 5 → Asking.
- Core Feature 6 → Asking; API Contract 1.
- Core Feature 7 → The only readers; The language gate.
- Core Feature 8 → Asking; Data Models.
- Core Feature 9 → Data Models; API Contract 4.
- Core Feature 10 → Asking; Data Models.
- Core Feature 11 → API Contract 5; Testing Approach 6.
- Success Metric 1 → Testing Approach 5.
- Success Metric 2 → Testing Approaches 4 and 5.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 1.
- Success Metric 6 → Testing Approach 2.

## Integration Points

- **TypeSafe.** `POST https://api.typesafe.ai/v1/systemone`, direct, no
  gateway. Tests and the QA gate never reach it: package and command tests
  inject an `http.RoundTripper` through the existing command dependencies,
  and the QA gate runs the binary only where nothing may be sent (no key, a
  reached ceiling, a non-English Spec or an unknown Spec), with the proxy
  variables pointed at a closed local port.
- **Roundfix Home.** The Judge Log under `.roundfix/judge/`. Nothing else in
  Roundfix Home is read or written, and the Run Database is not opened.
- **The authoring skills.** `write-prd` and `write-techspec` call the
  command; they never parse its output beyond reading it.

## Testing Approach

1. **The question file.** Task 01's Verification requires the embedded file
   to be byte-identical to the block under Questions and thresholds, and a
   test requires that it loads, compiles its patterns and pins `jev-1.13.0`.
2. **Planning.** New `internal/judge/pairs_test.go` and
   `internal/judge/language_test.go` build fixture Specs in temporary
   directories and require every rule of Citation claims, Goal pairs and The
   language gate, one case per rule, including an English and a Portuguese
   fixture.
3. **The boundary.** New `internal/judge/source_test.go` requires refusal of
   a symbolic-link PRD, a PRD over 1 MiB, a file outside `docs/adr`, and a
   proposed ADR. New `TestRequestsCarryOnlySpecArtifactText` in
   `internal/judge/judge_test.go` captures every request of a fixture run
   whose repository also holds a Go source file and a findings file carrying
   sentinel strings, and requires that every string in every state, split at
   `the cited decision`, occurs in the PRD, TechSpec or an ADR, that no
   sentinel and no key appears in any body, and that the only headers are the
   two named.
4. **Asking.** New `internal/judge/client_test.go`, `log_test.go` and
   `judge_test.go` require: the request body and headers; each threshold on
   both sides of its boundary; an answer from another model never raised or
   cleared; the retries with `Retry-After: 0`; each stop and skip reason of
   Asking; no request after a stop, at the ceiling or with an unreadable log;
   and one Judge Log line per request with every field and a `0600` file.
5. **The command.** New `internal/cli/spec_judge_test.go` reproduces Surface
   Transcripts 1 to 6 and API Contracts 2 and 5 with a fake transport, a fixed
   clock and a temporary Roundfix Home, and passes the key in the command's
   own environment.
6. **Documents.** Task 03's Verification requires the command in the
   command reference and the Roundfix Skill, and Task 04's requires the new
   headings of `write-prd` and `write-techspec`.

Each new gate is proved to fail: the Task that adds it records, in its
Result, the sabotage it applied and the failing test, then restores the code.

## Build Order

1. The question file, the readers, the language gate and planning, task_01
   (depends on: none).
2. The request, the Judge Log, the ceiling and `Run`, task_02 (depends on: 1).
   `Run` asks the plan task_01 builds.
3. The `spec judge` command, its help, the command reference and the Roundfix
   Skill, task_03 (depends on: 2).
4. The `write-prd` and `write-techspec` headings, task_04 (depends on: 3).
   They call the command task_03 ships, and both Tasks record skill versions
   in `skills/testdata/owned-skill-versions.json`.
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

This Spec is delivered after Spec 0194; task_03 edits the `spec` reference
and command files it creates.

## Risks & Considerations

- **A live call from a test.** A developer's shell may export the key. Every
  command test sets the key in the command's own environment and injects the
  transport, and the QA gate points the proxy variables at a closed port.
- **Parser drift.** The extraction reproduces the measurement, not the Spec
  Consistency Check. A change to either parser is not a change to the other.
- **Portuguese words in English text.** An English Spec quoting Portuguese
  stays above the English share; the measured margin is 0.103 against 0.08.
- **Rate limits.** The published limits are adjusting. A `429` is retried
  twice and then stops the run, which is a skip.
- **Help text.** Scripts that read `roundfix spec --help` see one more line.

## Decisions

- Advisory, pinned and fail-open. See ADR-0200.
- Spec artifacts only, key from the environment, every call logged, a
  monthly ceiling. See ADR-0201.
- The measurement's own extraction, kept in the question file.
- One request per judgment, as measured, with identical states asked once.
- A command, not a skill or an agent.
