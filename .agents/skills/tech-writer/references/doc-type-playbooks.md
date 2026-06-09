# Doc-type playbooks

Everyday document types: READMEs, docs pages, tracker issues and bug reports, PR descriptions, review comments, and async updates. Read the section for the type being written; skip the rest. PRDs, tech specs, ADRs, task files, and QA artifacts have their own playbooks — see the router in SKILL.md.

## README

Adapt as needed; not every project needs every section.

```markdown
# project-name

One sentence: what it does and for whom. One sentence: the concrete win.

## Install

Copy-pasteable command(s) for the main platform(s).

## Usage

The single most common task, shown end to end.

## Configuration (if applicable)

Table or list of options: name, type, default, what it changes.

## Contributing (for OSS)

How to set up dev environment, run tests, and submit changes.

## License
```

## Issues and bug reports

Action-forward title: component, then behavior — "Login button: disable on Safari 17 beta", not "Some login thing". One idea per issue; search for duplicates first.

Bug report body:

1. **Expected vs actual** — two bullets.
2. **Reproduction steps** — numbered, specific. A bug that cannot be reproduced cannot be fixed.
3. **Environment** — OS, browser, version.
4. **Evidence** — screenshot, log excerpt, failing test.
5. **Acceptance criteria** — pass/fail definition of done.

Add scope constraints when needed to prevent yak-shaving ("no new dependencies", "keep p95 under 200ms").

Copy this skeleton literally — plain-text field labels, no bold prefixes:

```markdown
# <component>: <observed behavior> (<conditions>)

Expected: <from the source>
Actual: <from the source>

## Reproduction steps

1. <step>

## Environment

- <only environments the source states>

## Evidence

- <links and traces from the source>

Hypothesis: <only if the evidence suggests a cause — never assert a root cause as fact. Omit this line if no evidence points anywhere.>
```

Include acceptance criteria only when the source provides them. If it does not, omit that section entirely — never invent pass conditions, browser lists, or run counts. For QA-pipeline bug reports with impact/severity metadata, use the QA playbook instead.

## PR descriptions

First line: a standalone imperative summary of what changed ("Delete the FizzBuzz RPC and replace callers with FooBar"). It becomes the version-history line people search.

Body answers WHY, not what — the diff already shows what. Include: context, decisions not visible in the code, trade-offs, known shortcomings, links to the bug/benchmark/design. Call out risky areas so reviewers know where to focus. Add screenshots for UI changes and testing instructions when non-obvious. "Fix bug" and "Phase 1" are not descriptions. Re-read the description before merging — it drifts during review.

## Code review comments

Comment on the code, never the developer. Explain why, not just what. Label severity so authors know what is mandatory: "Nit:", "Optional:", "FYI:". Point out problems; let the author choose the fix unless a specific fix matters. Praise good work — it teaches as much as criticism. Do not block on personal preference that no style guide backs. Approve when the change improves overall code health, not when it is perfect.

## Async updates and team posts

TL;DR first. Frontload all context — assume the reader has none; one self-contained message resolves in one reply instead of three follow-ups. Never send a bare "hello" or "got a sec?".

Status update shape: STATUS (what shipped, linked) / BLOCKERS (what, plus the specific unblocking ask) / ASKS (decision needed, owner, deadline) / FYI. Scannable in 30 seconds.

Set explicit response expectations: "FYI only", "need a decision by Thursday EOD", "blocking — please respond within 24h". Put decisions in a durable, searchable, linkable place — chat threads are not a decision record.

## Docs and quickstarts: pick the mode

Before writing, name which Diátaxis mode the page is, and do not mix modes on one page:

- **Tutorial** — learning by doing; a guaranteed-success path for a beginner.
- **How-to guide** — steps to solve one real task for someone who knows the basics.
- **Reference** — complete, accurate description of the machinery; neutral tone, tables welcome.
- **Explanation** — background and reasoning; the only place for "why it works this way".

Be prescriptive: recommend one path instead of enumerating every option. Keep code blocks shorter than one screen; never include shell prompt characters that break copy-paste.

When the source states behavior qualitatively, reuse the source's word and stop. For `Backoff: "exponential"` with `BaseDelay: 100ms`, write exactly this much: "Retries use exponential backoff starting from `BaseDelay`." Writing "doubles", "grows as 2^n", or "100ms, 200ms, 400ms" when the source gave no multiplier is fabrication — the actual implementation may use any growth factor, and your invented numbers will be wrong in print. The same applies to statistics: quote p95 as p95; never reinterpret a percentile as an average, median, or "half of all requests". In field and option tables, the description cell restates the source's comment and stops — it never adds when the field applies ("after the first call", "before the first retry"), what the opposite value does, or why someone would set it.