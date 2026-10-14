---
spec: 0267-no-secret-and-no-spec-content-leaves-the-machine
status: active
created: 2026-10-13
surfaces: [cli]
---

# No secret and no Spec content leaves the machine

This is an engineering-framed minimal PRD. The Spec fixes two defects in how
a Run treats credentials and enters through the TechSpec as a bug fix
(`docs/agents/spec-routing.md`).

Two adopted Backlog Entries record the defects, both from the Delivery Queue
item of Spec 0259 on 2026-10-12 and 2026-10-13:

- [QA evidence records environment values, secrets included](references/2026-10-13-qa-evidence-records-environment-values.md):
  three QA Runs committed the output of `make -p` as evidence, which held a
  third-party API key; GitHub push protection refused the push, and the
  branch was rebuilt without `qa/evidence`.
- [A QA Run can send Spec content to the judge provider](references/2026-10-13-a-qa-run-can-send-spec-content-to-the-judge-provider.md):
  the QA ran `roundfix archive --plan`, and because the QA environment
  inherited the judge's keys, Archive Advice sent the Spec's files to the
  judge provider without an authorization for that data to leave the machine.

Both are one class: what a Run's commands can read from the environment
leaves the machine. A characterization on 2026-10-13 at `9963ffd3`, through a
`go test -overlay` probe that wrote nothing to the repository, measured that
`acpxCommandEnv`, which builds the environment of every Agent Session process,
passes `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`, `ROUNDFIX_OPENROUTER_API_KEY` and
`ROUNDFIX_TYPESAFE_API_KEY` through, and that a Verification command run by
`daemon.ExecVerifier` sees `ROUNDFIX_TYPESAFE_API_KEY` when the Daemon's
environment holds it. The QA Report commit (`commitQAReport` in
`internal/daemon/task_engine.go`), the Delivery Queue push
(`GitHubCLI.PushBranch` in `internal/delivery/github.go`) and the Final Push
(`GitPusher.Push` in `internal/daemon/daemon.go`) read no file content before
they commit or push.

The judge reads exactly three variables, as `judge.Questions.KeyVariables`
lists them from `internal/judge/questions.json` and
`internal/openrouterkey`: the two Roundfix OpenRouter variables and
`ROUNDFIX_TYPESAFE_API_KEY`. The generic `OPENROUTER_API_KEY` and
`TYPESAFE_API_KEY` are never read by the judge, so they are not Judge Key
Variables.

## Project Constraints

- Identifier strategy: not applicable. No identifier, schema or stored format
  gains a key; the new park blocker is a string in the existing blocker field.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable. The change governs which processes
  receive the judge's credentials and what a push may carry, and adds no
  endpoint, credential or transport. Credentials stay in the environment, read
  only by the judge (`docs/agents/agent-instructions.md`: "Never read, print,
  commit, or generate secrets"). Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0273, recorded by this Spec, decides
  the Secret Scan and the withheld Judge Key Variables. ADR-0015 says the
  Daemon "commits the QA Report in its own commit either way" (ADR-0015:
  "commits the QA Report in its own commit either way"); ADR-0273 refines it
  for a report that carries a Secret Match. ADR-0201 keeps each judge key read from the
  environment only (ADR-0201: "Each is read from its variable in the
  environment and nowhere else"), which is why withholding the variables
  keeps the judge from every process of a Run. ADR-0239 names each
  stage's variables in one list, which this Spec reads instead of copying.
  ADR-0238 gives the Light Tier its key by a variable reference; this Spec
  moves that reference to the implementation Stage Key. ADR-0138 owns the
  Final Push, which gains the scan. ADR-0256: the Tasks change the
  `deliver` and configuration user-guide surfaces, so task_04 changes the
  covering skills. ADR-0189 and ADR-0233 govern the owned skills' version
  raise. ADR-0272 refines ADR-0233's retry rule for a
  parked Delivery Queue item; this Spec adds the `secret-found` park, which
  returns through `roundfix deliver retry` and leaves that rule unchanged. ADR-0244: the Spec declares its terms below. ADR-0259 and ADR-0271
  govern the new tests: `t.Parallel()` first or a named Sequential reason,
  and no real-time allowance for what a test does not assert. ADR-0265: the
  TechSpec carries its absence evidence, ADR alternatives and coverage
  cross-read. ADR-0200, ADR-0208, ADR-0209 and ADR-0231 do not
  apply: the judge's advisory role, source grouping, grouping suggestions and
  ceiling are unchanged; only which processes hold the judge's keys changes.
  ADR-0235 does not apply: no model selection or subscription rule changes.
  ADR-0250 does not apply: no Baseline module version or coverage record
  changes. ADR-0257 does not apply: the Daemon stays the only full gate inside
  a Run, and only the environment of its Verification commands changes.
  ADR-0258: the QA Task's rows declare their inputs, the scope row a
  `commit_range`. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. No lint, formatter, test-runner, Makefile or
  CI configuration changes. The Governed Paths the Tasks touch, measured with
  `speccheck.GovernedPath` through a `go test -overlay` probe that wrote
  nothing, are owned skill files, granted by the cycle AG approval of
  2026-10-13. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0267-no-secret-and-no-spec-content-leaves-the-machine/_authorization.md`.
  Bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`.

## Prerequisites

None that the Verification needs. Spec 0266 is delivered before this one and
edits `internal/cli/review.go` and the review correction proof of the
Delivery Queue retry; this Spec changes neither, and touches the Delivery
Queue only in `internal/delivery/engine.go`'s push step and in
`internal/delivery/park_class.go`. The operator rebases this plan on the
default branch after Spec 0266 merges and runs the strict check again.

## Goals

- QA evidence records environment variable names, never values, and never
  persists the output of a command that prints the environment.
- A QA Report whose files carry a Secret Match is never committed, and the
  Run says which file, line and pattern kind refused it, never the value.
- No push Roundfix makes, from the Delivery Queue or as a Final Push, sends a
  commit that adds a Secret Match.
- No Agent Session and no Verification command of a Run receives a Judge Key
  Variable, so no Run reaches the judge provider.

## User Stories

1. As the maintainer, I want a QA Run that captured a secret to stop before
   it commits, so that the value never reaches the branch and push
   protection never has to catch it.
2. As an operator reading a refusal, I want the file, the line and the kind
   of secret, so that I can remove it without the refusal printing the value
   again.
3. As the maintainer, I want the Delivery Queue to refuse a push that GitHub
   would block, so that the item parks with a next action instead of failing
   at the remote.
4. As the maintainer, I want a Run unable to send Spec content to the judge,
   so that data leaves the machine only where an authorization says so.
5. As a QA Agent, I want the QA contract to say how to record an environment
   fact, so that I never persist a dump to prove one.

## Core Features

1. The QA contract, in the `qa-gate` skill and the Daemon's QA prompt,
   requires recording environment variable names and whether they are set,
   never their values, and forbids persisting the output of an
   environment-printing command (`env`, `printenv`, `set`, `export -p`,
   `declare -x`, `make -p` and `make --print-data-base`).
2. A **Secret Scan** checks text against one reviewed set of **Secret
   Patterns** and reports each **Secret Match** as path, line and kind,
   never the value.
3. The QA Report commit runs the Secret Scan over the files it would stage;
   a Secret Match settles the QA Task `failed` and creates no commit.
4. The Delivery Queue push and the Final Push scan the lines added by every
   commit the push would send and refuse on a Secret Match; the Delivery
   Queue parks the item `secret-found` with a next action.
5. The **Judge Key Variables** are withheld from every Agent Session and every
   Verification command; a Light Tier session receives its OpenRouter key
   under the implementation Stage Key name.
6. The `roundfix` skill, the user guides and `CONTEXT.md` describe the scan,
   the park and the withheld variables.

## Non-Goals / Out of Scope

- Pushes Roundfix does not make: an operator's own `git push` and the plan
  PRs of the Supervisor.
- Git hooks and the QA format command, which are repository configuration;
  they keep the Run's environment (ADR-0273, Consequences).
- Redacting values, rewriting history, or removing a secret for the operator.
- How a Delivery Retry treats an item branch the operator rewrote; the
  existing Delivery Retry contract applies.
- A Spec authorization that opens the judge to a Run (ADR-0273).
- The generic `OPENROUTER_API_KEY` and `TYPESAFE_API_KEY`, which no Roundfix
  stage reads; the Light Tier already removes the generic OpenRouter variable.
- Detecting secrets that match no Secret Pattern; the set is a reviewed list,
  not a proof of absence.

## Success Metrics

1. Success Metric: Agent Session and Verification environments carry 0 Judge
   Key Variables, against 3 and 1 measured at `9963ffd3`, with a control
   variable still present.
2. Success Metric: a QA step whose evidence file adds a runtime-built sample
   of each Secret Pattern kind produces 0 QA Report commits and one refusal
   per file naming path, line and kind, with the sample in 0 bytes of output,
   Run Events or the Task file.
3. Success Metric: a push whose branch adds a sample in one commit and
   removes it in the next is refused before `git push` in both push paths,
   and the remote branch is absent afterwards.
4. Success Metric: the Secret Scan over every file `git ls-files` lists at the
   audited head reports 0 Secret Matches, so the set refuses none of the
   repository's own text.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not write,
found and read through Exa during authoring:

- GitHub Docs, "Working with push protection from the command line"
  (https://docs.github.com/en/code-security/how-tos/secure-your-secrets/work-with-leak-prevention/push-protection-on-the-command-line):
  "To resolve a blocked push, you must remove the secret from all of the
  commits it appears in." This is why the push scan reads every commit the
  push would send, not the net difference.
- GNU Make manual, "Variables from the Environment" and "Summary of Options"
  (https://www.gnu.org/software/make/manual/make.html): "Every environment
  variable that make sees when it starts up is transformed into a make
  variable with the same name and value." and, for `-p`/`--print-data-base`,
  "Print the data base (rules and variable values) that results from reading
  the makefiles". Together they explain why `make -p` output carries every
  environment value and is an environment-printing command.

## Research record

- Secondbrain: read `wiki/index.md`, then ran
  `qmd query "secret scanning QA evidence environment variables push protection" --all --files --min-score 0.3`
  and
  `qmd query "agent sandbox credentials environment leak API key third party provider" --all --files --min-score 0.3`.
  The first returned the `qa-gate` skill mirror (0.80), the Roundfix decision mirror on
  environment-blocked QA rows (0.62) and conexus Spec 0015's
  `projects/conexus/mirror/docs/history/specs/0015-sessao-autenticada-para-o-gate-de-qa/task_04.md`
  (0.62), which kept QA credentials out of documents by naming constants and
  never values, and whose QA confirmed it by reading the diff. That practice is
  Core Feature 1; this Spec adds the mechanical refusal the diff reading
  lacked. The second query returned conexus sandbox-validation notes about
  provider credentials, with no decision on withholding variables. No prior
  decision covers a Run's environment.
- Exa: searches found the GitHub push protection pages and the GNU Make
  passages quoted under Acceptance evidence. The GitHub pages also say push
  protection displays "Up to five detected secrets" at a time; the Secret Scan
  reports every match, because the operator must remove them all.

## Glossary

- adds: **Secret Pattern**
- adds: **Secret Match**
- adds: **Secret Scan**
- adds: **Judge Key Variables**

## Decisions

- Withhold the keys rather than teach `archive --plan` to detect a Run; see
  ADR-0273.
- Refuse rather than redact, and scan every pushed commit; see ADR-0273.
- No Spec authorization reopens the judge inside a Run; see ADR-0273.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the Secret Scan, the two refusal
points, the withheld variables and the build order.
