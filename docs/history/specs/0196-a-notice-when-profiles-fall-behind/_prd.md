---
spec: 0196-a-notice-when-profiles-fall-behind
status: archived
created: 2026-09-30
surfaces: [backend, cli, docs]
archived: "2026-10-01"
source_slug: 0196-a-notice-when-profiles-fall-behind
---


# A notice when profiles fall behind the recommendation

Spec 0189 gives Roundfix one dated Recommended Profile per Agent Work Category,
and its built-in profiles follow that snapshot with each release. A configured
profile does not follow anything:

- On 2026-09-30 this repository's own Project Config still led with
  `gpt-5.6-sol`, a day after `gpt-6.1-sol` replaced it as the Codex workhorse.
  Nothing in Roundfix said so. The maintainer found out from a provider page
  and changed the profiles by hand (#292).
- Adopting a recommendation means writing a profile fragment by hand and
  passing it to `roundfix profiles configure`.
- A repository that keeps another model on purpose has no way to say so, so
  any reminder would repeat forever.

This Spec adds the comparison, a notice where a person will read it, one
command to adopt the recommendation, and a way to record a deliberate
difference. It is delivered after Spec 0189, because it reads the Recommended
Profile that Spec adds.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. A Profile Deviation
  is named by the snapshot date it was declared against. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the comparison reads local
  configuration and static data in the binary. No credential is read and no
  network call is added. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0181 (this Spec) says where the
  comparison runs, what it costs and how a deliberate difference is declared.
  ADR-0180 supplies the Recommended Profile it compares with. ADR-0049 keeps a
  profile atomic, so adoption replaces a whole profile and a replaced profile
  loses its deviation. ADR-0107 proves only the categories a configuration
  defines, and the comparison covers the same categories. ADR-0140 proves an
  Agent Selection as the exact advertised tuple, and adoption writes no tuple
  without that proof. This Spec's gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096,
  ADR-0104, ADR-0117, ADR-0155 and ADR-0156, and ADR-0093 and ADR-0094 check
  its consistency by citation and artifact presence. ADR-0166 records a Task's
  undeclared paths, ADR-0167 keeps the pre-PR Pull Request row from deciding a
  qualifying partial, and ADR-0176 reads citations only from authored text.
  All hold.
  ADR-0097 cites ADR-0080 but carries a QA row forward on unmoved evidence.
  ADR-0168 cites ADR-0093 but sets the horizon of the related-ADR check.
  ADR-0182 (Spec 0190) cites ADR-0096 but moves mechanical facts to Task
  settlement. ADR-0183 and ADR-0184 (Spec 0191) cite ADR-0093 and ADR-0156 but
  add receipts for attributed claims and transcripts for command surfaces.
  This Spec changes none of the five, so they do not apply.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("considere autorizado a ajustar todas as skills se necessário"),
  recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A maintainer who upgrades Roundfix is told which configured profiles differ
  from the recommendation the release ships.
- One command shows each difference, and one command adopts the recommendation
  after proving it.
- A deliberate difference is recorded once, with its reason, and stays quiet
  until the recommendation changes.
- The comparison never delays, blocks or fails any command that does work.

## User Stories

1. As a maintainer who upgrades Roundfix, I want to be told when my configured
   profiles differ from the recommendation the new release ships, so that a
   retired model does not stop my Runs by surprise.
2. As a maintainer, I want one command that shows each difference and one that
   adopts the recommendation after proving it, so that I do not write profile
   fragments by hand.
3. As a maintainer who keeps another model on purpose, I want to record that
   choice with its reason, so that Roundfix stops reminding me until the
   recommendation changes.
4. As an operator of unattended Runs, I want the comparison to stay out of
   Run start and Preflight, so that a recommendation never delays or blocks
   work.

## Core Features

1. **A profile can carry a Profile Deviation.** A configured profile accepts
   one more key, `deviation`, with the snapshot date it was declared against
   (`from`) and a `reason`. Both are required, and a malformed deviation is
   refused with the path of the offending key. `roundfix profiles configure`
   writes a deviation that its fragment carries, and a profile replaced
   without one loses it.
2. **`roundfix profiles check` reports the comparison.** It compares the
   effective profile of every configured category with the Recommended
   Profile and reports each as `current`, `differs` or `pinned`. A category is
   `pinned` when it differs and its deviation names the snapshot the binary
   ships. The command is read-only and offline, prints text or `--json`, and
   exits `0` whether or not a category differs.
3. **`roundfix profiles show` states the status of each category.** Each
   category gains its recommendation status and, when it has one, its
   deviation.
4. **`roundfix profiles check --apply` adopts the recommendation.** It writes
   the Recommended Profile for every differing category to the named scope,
   through the same proof, preview and confirmation as
   `roundfix profiles configure`. A pinned category is skipped. `--dry-run`
   writes nothing.
5. **Doctor reports the comparison in one line that never fails.** The line
   `recommendations:` follows `profiles:` and is `ok`, `found`, or `skipped`
   with the error text when the comparison itself fails. It never changes
   Doctor's exit code.
6. **`roundfix upgrade` ends with the notice.** Every invocation that reports
   a release outcome writes the comparison to standard error. Standard output
   and the exit code do not change. After a release is installed, the notice
   comes from the installed executable. When the comparison cannot run, one
   line says so and the exit code still does not change.
7. **The comparison stays off the hot paths.** Only `profiles check`,
   `profiles show`, `doctor` and `upgrade` reach it, and a test fails when
   another command does.

## Declared breaks

- A configuration that uses `deviation` is rejected by a Roundfix older than
  this Spec, because profile keys are strict.
- `roundfix upgrade` writes to standard error on the four outcomes that used
  to leave it empty: no release published, already current, upgrade available
  and upgraded.
- `roundfix doctor` prints one more line.
- `roundfix profiles show` prints one more line per category, two when the
  category has a deviation. Its JSON gains `recommendation_status` and, when
  present, `deviation`. The schema stays `roundfix/profiles/v2`, because both
  fields are additions.
- `roundfix profiles configure` prints a deviation line in its preview and a
  `deviation` field in its JSON, only for a profile that carries one.

## Non-Goals / Out of Scope

- Running the comparison when a Run starts, in Preflight, in the Daemon, or in
  `implement`, `watch`, `review` or `deliver`.
- Failing, blocking or changing the exit code of any command because a
  profile differs.
- Writing a profile without `--apply`, or without the proof and confirmation
  `profiles configure` already requires.
- Comparing optional categories a configuration does not define. They inherit
  `general` and are reported as inherited.
- Fetching a newer recommendation from the network. The snapshot is the one in
  the binary.
- A notice on the commands that already report version freshness. That report
  stays as it is.
- Changing this repository's `.roundfixrc.yml`. After Spec 0189 its profiles
  equal the recommendation.

## Success Metrics

1. Success Metric: with a Project Config whose `backend` profile differs from
   the recommendation, `roundfix profiles check` exits `0`, names `backend`
   with the configured and the recommended profile, and its summary counts one
   difference.
2. Success Metric: a deviation whose `from` equals the shipped snapshot makes
   the category `pinned` and removes it from the difference count. The same
   deviation with another date leaves the category `differs`.
3. Success Metric: `roundfix profiles check --apply --scope project --yes`
   proves and writes the Recommended Profile for each differing category and
   leaves a pinned category's bytes untouched. A following
   `roundfix profiles check` counts no difference.
4. Success Metric: for each release outcome of `roundfix upgrade`, standard
   output and the exit code are byte-identical to the characterization
   recorded before the change, and standard error carries the notice.
5. Success Metric: after an install, the notice is what the installed
   executable's `profiles check` printed. When that command fails, `upgrade`
   still exits `0` and standard error carries one line that says the
   comparison did not run.
6. Success Metric: adding a call to the comparison from `implement` fails a
   test.

## Recorded limits

- An upgrade performed by an executable older than this Spec prints no notice.
  The notice starts with the upgrade after that one.
- A deviation lasts one snapshot. A release with a new snapshot date reports
  the category again, even when the recommendation for that category did not
  change.
- `--apply --scope user` cannot change a category the Project Config defines.
  It names the category and leaves it.
- `--apply` opens disposable Agent Sessions to prove what it writes, as
  `profiles configure` does. `profiles check` without `--apply` opens none.
- The comparison is exact: the same Preferred Selection and the same Fallback
  Chain, in order. A profile that adds a third fallback differs.

## Decisions

- **Where the comparison runs, and what it may do.** See ADR-0181.
- **A deviation is dated, not described.** `from` is the snapshot date. It is
  short to write by hand, and it makes the choice be decided again once per
  snapshot.
- **One renderer.** `profiles check`, the `upgrade` notice and the Doctor line
  report the same counts from the same comparison, so they cannot disagree.
- **Adoption reuses `profiles configure`.** `--apply` builds the fragment and
  hands it to the existing write, so proof, preview, confirmation, rollback
  and exit codes are the existing ones.
- **The notice goes to standard error.** Standard output of `upgrade` is one
  line that scripts read.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- This repository's history: the commit before #292 (`9e439dbb^`) holds a
  `.roundfixrc.yml` whose implementation profiles lead with `gpt-5.6-sol`, and
  #292 replaced them by hand on 2026-09-30. OpenAI's announcement of
  `gpt-6.1-sol` (<https://openai.com/index/introducing-gpt-6-1-sol/>) is dated
  2026-09-29. No Roundfix command reported the difference in between.
- Published practice for the adoption flow: the Qtile project documents a
  configuration migration command that shows each change, asks before it
  writes, accepts `--yes` to skip the question, and has a mode that only lists
  what would change
  (<https://docs.qtile.org/en/stable/manual/commands/shell/qtile-migrate.html>,
  read 2026-09-30). `--apply` follows the same shape with Roundfix's own proof
  in front of the write.
- Measured on a scratch copy of `9e439dbb`: a new Doctor line invalidates five
  tests in `internal/cli/doctor_test.go` and none elsewhere; a notice on
  standard error invalidates two tests in `internal/cli/upgrade_test.go`; a
  new line and a new JSON field in `profiles show` invalidate none.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and the query `qmd query
"aviso após upgrade quando a configuração diverge da recomendação; desvio
deliberado declarado com motivo; adoção explícita" --all --files --min-score
0.3`. Its hits were Specs and QA Reports of other projects, none about a
configuration notice, so no Secondbrain file informs this design and none is
cited. The query recorded in Spec 0189 for the same subject had the same
result. Exa was asked for published documentation of command-line tools that
report an outdated configuration without failing and offer an explicit command
to adopt the change. The Qtile page above is the primary source it found for
the adoption flow. The other results were small projects whose documents agree
that warnings belong on standard error and that exit codes are the contract;
they are not cited as evidence.

## Technical candidate

The [_techspec.md](_techspec.md) records the data shapes, the command
transcripts, coverage and build order.
