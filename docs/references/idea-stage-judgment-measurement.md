# Idea-stage judgment measurement

Recorded by the author of Spec 0393 on 2026-11-14, while writing its PRD. It measures the yes or no questions that judge whether a feature needs the Idea stage before its PRD, on a labeled set built from this repository's own Backlog Entries, with controls. It replaces nothing: the four-case pilot of the same day stays the operator's record.

## What was sent, and to whom

- Recipient: TypeSafe, `https://api.typesafe.ai/v1/systemone`, with the key read from `ROUNDFIX_TYPESAFE_API_KEY` (the name only; no value was printed or stored).
- Models: `jev-1.13.0` twice (the version the judge pins) and `jev-1.14.0` once (requested as `jev-latest`). Every response reported the model named here.
- Text: only Backlog Entries and Findings of this repository, read from Git history at the commit that added each one, and controls and translations the author wrote. No message of another project, no Spec artifact, no source code and no request typed in a session was sent. Five feature entries that reproduce another project's message, and two that quote a relayed request, were left out for that reason.
- Secret scan: every case, and the whole source file of every repository case, passed `internal/secretscan` before any request (90 cases and 68 files scanned, none flagged).
- State: the shape the follow-up judge already sends, one item per request:

```json
{"project": "roundfix", "item": {"kind": "backlog", "type": "feat", "origin": "roundfix", "title": "<first heading>", "text": "<body without front matter, cut to 4000 characters>"}}
```

## The questions

Six Noul questions were asked together in each request. The first five are the pilot's questions, reworded for the `item` and `project` state; the sixth was added before any answer was read.

| id | instructions | yes means | no means |
| --- | --- | --- | --- |
| `shape_open` | `item` is one pending request addressed to the software project `project`. Is the right shape of the solution for `item` still open? | At least two materially different designs could satisfy the item, and the item does not settle which one to build. | The item names one concrete design, or only one reasonable design exists. |
| `unverified_assumption` | `item` is one pending request addressed to the software project `project`. Does `item` rest on an assumption about its users, the business or the market that it gives no evidence for? | A claim about what users need, will adopt or will pay for is taken for granted, and the work would be useless if that claim were wrong. | The need is shown by evidence in the item itself, such as a measured defect, a reported failure or a decision the product owner already made. |
| `market_question` | `item` is one pending request addressed to the software project `project`. Does deciding whether to build `item` depend on a market or competitive question? | The decision needs to know what other tools already offer, how the product would be positioned, or whether people outside the team would adopt it. | The decision depends only on the product's own behavior and its current users. |
| `new_capability` | `item` is one pending request addressed to the software project `project`. Does `item` ask for a kind of capability that `project` does not have today? | It adds a new capability, surface or concept to the product. | It corrects or adjusts behavior the product already has. |
| `behavior_specified` | `item` is one pending request addressed to the software project `project`. Does `item` already specify the behavior to build precisely enough to write its requirements without first asking whether to build it or in what shape? | The inputs, the outputs and the rules of the wanted behavior are stated, and what remains is detail. | What exactly should be built, for whom, or whether it is worth building is still an open question. |
| `worth_unsettled` | `item` is one pending request addressed to the software project `project`. Is it still unsettled whether `item` is worth building? | The item states a hypothesis, asks to try or measure something first, or leaves open whether the benefit justifies its cost. | The item states a need that a decision, a measurement or a reported failure already settles. |

## The labeled set

- **Feature items (38).** Every entry typed `feat` or `feature` that Git history holds under the backlog or among a Spec's adopted references (45), less the seven left out above. Each is labeled by what happened to it, read from its last state in Git, the Archive Record of the Spec that adopted it and that Spec's Non-Goals at the record's `source_revision`: `as asked` (23) when one Spec delivered the mechanism the entry's Shape names; `reshaped` (5) when the delivering Spec or an ADR rejected that mechanism, or the Spec was replaced before or right after delivery; `split` (4) when the entry was delivered across more than one Spec; `declined` (6) when it ended without a consuming Spec. The last three are the positive class (15): the cases where a study of whether and in what shape to build would have had something to change.
- **Fix controls (20)** and **Finding controls (10).** Backlog Entries typed `fix` and Findings, drawn with a fixed seed among those that quote no one. They are outside the judged population and show what the questions answer when asked about a correction.
- **Author controls (12).** Six short requests the author wrote, three expected to be indicated and three not, each in English and in Portuguese.
- **Translations (10).** The author's Portuguese translations of ten feature items, five from the positive class and five `as asked`, paired with their English originals.

The labels are the author's reading of the records, one reader, made before the answers were computed. Two are debatable and are kept as labeled: F05 asks for a measurement and was closed after the measurement said not to adopt the feature, and F20 was delivered as asked and replaced within a week.

## Run accounting

| measure | value |
| --- | ---: |
| requests | 270 (90 cases, three passes) |
| input tokens | 347,715 |
| cost at US$0.042 per million input tokens | US$0.0146 |
| input tokens per request, six questions | 994 to 2105, median 1197 |
| cost per request | US$0.000054 |

## Results on `jev-1.13.0`, first pass

### How well each question ranks the feature items

AUROC of the positive class against `as asked`, over the 38 feature items, with a 95% bootstrap interval (4,000 resamples of the cases, seed 393).

| score | AUROC | 95% interval | mean, positive | mean, as asked |
| --- | ---: | --- | ---: | ---: |
| `shape_open` | 0.75 | [0.58, 0.90] | 0.66 | 0.43 |
| `unverified_assumption` | 0.66 | [0.48, 0.83] | 0.41 | 0.32 |
| `market_question` | 0.81 | [0.65, 0.94] | 0.12 | 0.08 |
| `new_capability` | 0.76 | [0.60, 0.90] | 0.91 | 0.71 |
| not behavior_specified | 0.75 | [0.57, 0.91] | 0.68 | 0.45 |
| `worth_unsettled` | 0.69 | [0.49, 0.86] | 0.63 | 0.42 |
| highest of the three kept | 0.74 | [0.57, 0.89] | 0.67 | 0.48 |
| text length | 0.68 | [0.48, 0.87] | 1971 | 1269 |

### The verdict at 0.7 and 0.3

The verdict is computed from `shape_open`, `unverified_assumption` and `market_question`: indicated when any answer is 0.7 or more, not indicated when every answer is 0.3 or less, undecided otherwise.

| what happened | indicated | undecided | not indicated |
| --- | ---: | ---: | ---: |
| as asked | 4 | 12 | 7 |
| reshaped | 2 | 3 | 0 |
| split | 2 | 1 | 1 |
| declined | 4 | 2 | 0 |
| fix controls | 9 | 6 | 5 |
| Finding controls | 7 | 2 | 1 |

- Of the 12 feature items read as indicated, 8 are in the positive class: 67%, with a 95% Wilson interval of 39% to 86%.
- Of the 8 read as not indicated, 1 is in the positive class (F38, a profile later split in two Specs): 12%, interval 2% to 47%.
- 18 of 38 (47%) are undecided, 6 of them in the positive class.

Other cut points on the same answers, for comparison only:

| yes at | no at | indicated (positive + as asked) | undecided | not indicated (positive + as asked) |
| ---: | ---: | --- | ---: | --- |
| 0.7 | 0.3 | 8 + 4 | 18 | 1 + 7 |
| 0.8 | 0.2 | 4 + 4 | 28 | 0 + 2 |
| 0.6 | 0.4 | 9 + 5 | 12 | 2 + 10 |
| 0.7 | 0.5 | 8 + 4 | 9 | 3 + 14 |

### The author's controls

All six controls read the same in both languages: the three expected to be indicated are indicated (the loose product idea, the open shape and the unchecked market claim), the two specified features are not indicated, and the defect report is undecided, because `shape_open` answers 0.42 in English and 0.56 in Portuguese for a defect whose fix is not described.

### Portuguese against English

Sixteen pairs, the same text in both languages.

| question | mean absolute difference | largest difference |
| --- | ---: | ---: |
| `shape_open` | 0.060 | 0.19 |
| `unverified_assumption` | 0.024 | 0.09 |
| `market_question` | 0.014 | 0.07 |

The verdict changed in 1 of 16 pairs: F07, whose `shape_open` moved from 0.65 to 0.71 and so from undecided to indicated. No pair moved between indicated and not indicated.

### Two passes of the same model

The second `jev-1.13.0` pass over the same 90 states differed from the first by 0.009 on average and by at most 0.07 on the three kept questions. That was enough to change the verdict of 6 cases (F02, F03, F09, X06, F07-pt, F12-pt), each with an answer within 0.07 of a cut point, always between undecided and a neighbour.

### `jev-1.14.0`

On the same states, `jev-1.14.0` ranks the feature items at AUROC 0.69 [0.52, 0.85] for `shape_open`, 0.62 [0.43, 0.80] for `unverified_assumption` and 0.84 [0.70, 0.96] for `market_question`, and gives a different verdict from `jev-1.13.0` in 9 of 90 cases. Its answers are not interchangeable with the pinned model's.

## Reading

- **Questions kept: `shape_open`, `unverified_assumption` and `market_question`.** They are the three criteria the maintainer named. `shape_open` carries most of the usable signal. `market_question` ranks the items but never passes 0.31 on the feature items; it reaches 0.93 or more on the two controls that hold a market claim. `unverified_assumption` is weak here, with an interval that includes 0.5.
- **Questions dropped: `new_capability`, `behavior_specified` and `worth_unsettled`.** The first two rank as well as `shape_open` but say what kind of request it is, not what the Idea stage would study. `worth_unsettled` reads five of the six declined items as yes, and also five items delivered as asked, among them the two that only ask for a measurement; its interval includes 0.5.
- **Thresholds: not fixed by this measurement.** With 15 positive and 23 negative cases every interval is about 0.3 wide or wider, the text's length alone ranks the items at 0.68, and no cut point in the table is distinguishable from its neighbours. The documentation's starting points, yes at 0.7 or more and no at 0.3 or less, stay, with the wide undecided band going to the user.
- **Portuguese: judged as English is.** The mean shift on the kept questions is 0.06 or less and one verdict in 16 moved by one step, about what two passes of the same model move.
- **A correction is not a feature.** Nine of 20 fix entries and seven of 10 Findings read as indicated, because the fix of a defect is often not designed yet. The judgment applies to features only.
- **The model stays `jev-1.13`.** The thresholds and these numbers belong to it; moving the pin means measuring again.

## Limits

- One repository, one labeler, 38 feature items, labels read from what happened and not from an experiment: an item delivered as asked may still have been the wrong shape, and a study might not have changed a reshaped one.
- The feature items are Backlog Entries a maintainer or an operator had already written in English; a message or a note may read differently, and none was measured, because no message of another project may be sent for a measurement.
- Longer entries are more often in the positive class, and the questions were not separated from length.
- The questions were chosen on the same set they are reported on.

## Index of the cases

| id | group | language | source | label | basis |
| --- | --- | --- | --- | --- | --- |
| F00 | feature | en | `docs/backlog/2026-08-06-atomic-inbox-capture-helper.md@142e3232` | declined | deferred without a consuming Spec; deleted from the history root by the 0272 delivery |
| F01 | feature | en | `docs/backlog/2026-08-08-record-usage-per-agent-session.md@1a31c965` | declined | deferred without a consuming Spec; deleted from the history root by the 0272 delivery |
| F02 | feature | en | `docs/backlog/2026-08-12-speak-acp-directly-instead-of-through-acpx.md@385b5f85` | declined | deferred without a consuming Spec; deleted from the history root by the 0272 delivery |
| F03 | feature | en | `docs/backlog/2026-09-08-support-claude-fable-5-1-and-gpt-6-astra.md@7e67308b` | split | adopted by portfolio Spec 0123, retired on 2026-09-25 with its features shipped through other Specs (commit 595876fc) |
| F04 | feature | en | `docs/backlog/2026-09-25-baseline-history-moves-need-separate-approval.md@7c06573f` | reshaped | Spec 0184 Non-Goals reject the separate approval the entry asks for (ADR-0173) and deliver the citation-impact report instead |
| F05 | feature | en | `docs/backlog/2026-09-25-jev-gateway-baseline-before-routing.md@595876fc` | declined | closed done without a Spec: measured, routing not adopted, front closed on 2026-09-29 |
| F06 | feature | en | `docs/backlog/2026-09-25-post-merge-cleanup-of-runs-and-branches.md@7c06573f` | as asked | Spec 0175 delivers the automatic release after a delivery merge and the manual path |
| F07 | feature | en | `docs/backlog/2026-09-25-storage-reclaim-notice.md@7c06573f` | reshaped | Spec 0183 Non-Goals reject the storage line in runs list and Doctor that the entry's Shape names |
| F08 | feature | en | `docs/backlog/2026-09-28-blocking-review-after-archive-parks-publication.md@b92aefda` | as asked | Spec 0179 |
| F09 | feature | en | `docs/backlog/2026-09-28-implement-spec-skill-delegates-to-roundfix.md@b92aefda` | as asked | Spec 0180 |
| F10 | feature | en | `docs/backlog/2026-09-28-queue-preparation-inventory.md@b92aefda` | as asked | Spec 0180; the Secondbrain inbox left out |
| F11 | feature | en | `docs/backlog/2026-09-28-queue-wide-limits-and-one-pending-question.md@b92aefda` | as asked | Spec 0180; spend and concurrency limits left out, the deadline and the single pending question delivered |
| F12 | feature | en | `docs/backlog/2026-09-28-revalidate-later-specs-when-prerequisites-change.md@b92aefda` | as asked | Spec 0180 |
| F13 | feature | en | `docs/backlog/2026-09-28-review-findings-get-evidence-backed-dispositions.md@b92aefda` | as asked | Spec 0179 |
| F14 | feature | en | `docs/backlog/2026-09-28-unselected-review-providers-receive-no-request.md@b92aefda` | as asked | Spec 0179 |
| F15 | feature | en | `docs/backlog/2026-09-30-measure-why-corrective-tasks-happen.md@d6528a08` | as asked | Spec 0214 |
| F16 | feature | en | `docs/backlog/2026-09-30-re-measure-the-task-lint-judgment-on-a-cleaner-label.md@c49a0128` | as asked | Spec 0214 |
| F17 | feature | en | `docs/backlog/2026-09-30-release-plan-does-not-report-skill-and-guide-checks.md@d6528a08` | as asked | Spec 0223 |
| F18 | feature | en | `docs/backlog/2026-09-30-run-grok-through-the-cursor-runtime.md@30e8504f` | as asked | Spec 0217 |
| F19 | feature | en | `docs/backlog/2026-09-30-stack-modules-need-rules-with-force.md@d6528a08` | split | Spec 0206 Non-Goals hand the composed profile and the frontend layout to Spec 0207 and leave five listed items open |
| F20 | feature | en | `docs/backlog/2026-10-01-try-the-jev-router-for-cheap-agent-categories.md@da4c76d1` | reshaped | Spec 0218 delivered the router; ADR-0235 and the Light Tier of Spec 0233 replaced it within a week |
| F21 | feature | en | `docs/backlog/2026-10-03-a-grok-fallback-for-docs-and-chore.md@23000806` | declined | declined on 2026-10-04: the cursor runtime stays opt-in only |
| F22 | feature | en | `docs/backlog/2026-10-03-branch-prefix-is-a-required-decision.md@03ea3c52` | as asked | Spec 0223 |
| F23 | feature | en | `docs/backlog/2026-10-05-a-judge-assigned-model-tier-per-task.md@e69ac275` | reshaped | Spec 0233 delivers a Light Tier by declared complexity; its Non-Goals keep the judge's tier advisory and unread at dispatch |
| F24 | feature | en | `docs/backlog/2026-10-06-history-keeps-only-what-the-secondbrain-needs.md@4d6fd152` | split | delivered across Specs 0242 and 0243 and the later history Specs |
| F25 | feature | en | `docs/backlog/2026-10-12-a-decomposed-spec-cannot-be-cancelled.md@7a8b7f70` | as asked | Spec 0261; the maintainer's decision was already in the entry |
| F26 | feature | en | `docs/backlog/2026-10-14-promoted-copies-keep-prose-paths-to-unpromoted-folder-files.md@d0155acf` | as asked | Spec 0272 |
| F27 | feature | en | `docs/backlog/2026-10-14-spec-check-does-not-flag-a-guide-task-without-the-coverage-record.md@d0155acf` | as asked | Spec 0274 |
| F28 | feature | en | `docs/backlog/2026-10-16-spec-check-does-not-flag-an-owned-skill-edit-without-its-version-record.md@123f6524` | as asked | Spec 0276 |
| F29 | feature | en | `docs/backlog/2026-11-01-a-spec-authorization-carries-the-release-decision.md@619cff66` | as asked | Spec 0338; the maintainer's decision was already in the entry |
| F30 | feature | en | `docs/backlog/2026-11-03-the-backlog-and-findings-contracts-name-the-waiting-declaration.md@c2ae4621` | as asked | Spec 0346 |
| F31 | feature | en | `docs/backlog/2026-11-10-an-adr-marked-as-not-applying-in-portuguese-still-seeds-the-related-adrs.md@04ece683` | as asked | Spec 0379 takes the first of the three shapes the entry lists |
| F32 | feature | en | `docs/backlog/2026-11-12-roundfix-installs-and-updates-through-a-homebrew-tap.md@2e99852c` | declined | declined on 2026-11-13 after the maintainer's study review (commit 4cb2e08b) |
| F33 | feature | en | `docs/specs/0128-release-planning-with-bare-stable-tags/references/2026-09-08-plan-releases-with-bare-stable-tags.md@6b8ea487` | reshaped | its Spec 0128 was retired as superseded by Spec 0147 before delivery (commit 7806acfa) |
| F34 | feature | en | `docs/specs/0195-owned-skills-and-a-release-step-that-follow-the-bundle/references/2026-09-30-a-release-does-not-check-skills-and-guides.md@d6528a08` | as asked | Spec 0195 |
| F35 | feature | en | `docs/specs/0229-a-jev-router-that-checks-credit-and-names-its-model/references/2026-10-04-jev-router-cost-controls.md@1ab1fa1b` | as asked | Spec 0229 |
| F36 | feature | en | `docs/specs/0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd/references/2026-10-16-the-glossary-is-not-held-to-the-rule-it-applies-to-specs.md@123f6524` | as asked | Spec 0274 |
| F38 | feature | en | `docs/specs/0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff/references/2026-10-22-a-baseline-profile-for-python-repositories.md@56dd1a76` | split | delivered across Specs 0309 and 0310, with the profile renamed and amended during planning |
| X01 | fix-control | en | `docs/backlog/2026-09-25-fallback-proof-timeout-blocks-dispatch.md@7c06573f` | fix | a Backlog Entry typed fix, outside the judged population |
| X02 | fix-control | en | `docs/backlog/2026-09-25-skills-lock-preimage-taken-after-the-fetch.md@7c06573f` | fix | a Backlog Entry typed fix, outside the judged population |
| X03 | fix-control | en | `docs/backlog/2026-10-05-the-judge-log-under-counts-router-spend.md@7eabb720` | fix | a Backlog Entry typed fix, outside the judged population |
| X04 | fix-control | en | `docs/backlog/2026-11-10-the-catalog-skill-bump-leaves-a-stale-catalog-digest-in-the-setup-manifest.md@04ece683` | fix | a Backlog Entry typed fix, outside the judged population |
| X05 | fix-control | en | `docs/backlog/2026-10-14-a-dismissal-after-a-committed-fix-needs-a-reset.md@d0155acf` | fix | a Backlog Entry typed fix, outside the judged population |
| X06 | fix-control | en | `docs/backlog/2026-10-12-the-baseline-hardcodes-the-better-auth-base-path.md@7a8b7f70` | fix | a Backlog Entry typed fix, outside the judged population |
| X07 | fix-control | en | `docs/backlog/2026-09-08-capture-inbox-contract-matches-the-fleet-door.md@7e67308b` | fix | a Backlog Entry typed fix, outside the judged population |
| X08 | fix-control | en | `docs/backlog/2026-08-14-the-archive-destination-still-depends-on-which-spec-root-is-configured.md@6a27b556` | fix | a Backlog Entry typed fix, outside the judged population |
| X09 | fix-control | en | `docs/backlog/2026-09-08-completed-work-has-a-terminal-state-without-a-spec.md@7e67308b` | fix | a Backlog Entry typed fix, outside the judged population |
| X10 | fix-control | en | `docs/backlog/2026-09-08-secondbrain-capture-has-one-durability-owner.md@7e67308b` | fix | a Backlog Entry typed fix, outside the judged population |
| X11 | fix-control | en | `docs/backlog/2026-10-07-history-sanitize-refuses-older-archived-specs.md@405274c6` | fix | a Backlog Entry typed fix, outside the judged population |
| X12 | fix-control | en | `docs/backlog/2026-09-25-read-only-commands-and-schema-versions.md@7c06573f` | fix | a Backlog Entry typed fix, outside the judged population |
| X13 | fix-control | en | `docs/backlog/2026-09-30-an-installed-skill-that-matches-its-lock-can-trail-the-snapshot.md@30e8504f` | fix | a Backlog Entry typed fix, outside the judged population |
| X14 | fix-control | en | `docs/backlog/2026-10-16-a-qa-gate-archive-plan-contacts-the-judge-an-offline-qa-forbids.md@0422580f` | fix | a Backlog Entry typed fix, outside the judged population |
| X15 | fix-control | en | `docs/backlog/2026-09-25-help-token-accepted-anywhere.md@7c06573f` | fix | a Backlog Entry typed fix, outside the judged population |
| X16 | fix-control | en | `docs/backlog/2026-09-25-declared-none-must-fit-one-line.md@7c06573f` | fix | a Backlog Entry typed fix, outside the judged population |
| X17 | fix-control | en | `docs/backlog/2026-10-16-review-convention-c2-asks-a-qa-task-for-settlement-sections.md@0422580f` | fix | a Backlog Entry typed fix, outside the judged population |
| X18 | fix-control | en | `docs/backlog/2026-10-20-runs-causes-reads-no-diagnostic-of-a-run-that-recorded-no-artifact-root.md@510ed75e` | fix | a Backlog Entry typed fix, outside the judged population |
| X19 | fix-control | en | `docs/backlog/2026-11-13-a-completion-notice-reads-an-adopted-record-as-closed.md@f953f77f` | fix | a Backlog Entry typed fix, outside the judged population |
| X20 | fix-control | en | `docs/backlog/2026-10-12-an-archive-record-points-at-a-branch-commit-a-squash-merge-drops.md@7a8b7f70` | fix | a Backlog Entry typed fix, outside the judged population |
| N01 | finding-control | en | `docs/findings/2026-07-23-setup-context-driven-adoption-process-improvements.md@4f3fd527` | finding | a Finding, outside the judged population |
| N02 | finding-control | en | `docs/findings/2026-09-15-operational-command-tests-reach-github.md@4b86285f` | finding | a Finding, outside the judged population |
| N03 | finding-control | en | `docs/findings/_archived/2026-08-05-contract-seams-between-daemon-gate-and-archive.md@142e3232` | finding | a Finding, outside the judged population |
| N04 | finding-control | en | `docs/findings/2026-08-06-rollup-baseline-and-derived-tooling.md@142e3232` | finding | a Finding, outside the judged population |
| N05 | finding-control | en | `docs/findings/2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md@385b5f85` | finding | a Finding, outside the judged population |
| N06 | finding-control | en | `docs/findings/2026-07-24-greenfield-agent-guidance-acceptance-target.md@4f3fd527` | finding | a Finding, outside the judged population |
| N07 | finding-control | en | `docs/findings/2026-08-14-preflight-starves-when-the-machine-is-busy.md@6a27b556` | finding | a Finding, outside the judged population |
| N08 | finding-control | en | `docs/findings/2026-09-14-the-mechanical-stage-audits-only-the-current-runs-task-commits.md@7a15e970` | finding | a Finding, outside the judged population |
| N09 | finding-control | en | `docs/findings/2026-08-03-gate-and-review-rounds-need-a-convergence-rule.md@77309d2b` | finding | a Finding, outside the judged population |
| N10 | finding-control | en | `docs/findings/2026-09-25-an-opt-in-jev-judgment-layer-is-proposed.md@7c06573f` | finding | a Finding, outside the judged population |
| C1-en | author-control | en | written by the author | not expected | a control written by the author for this measurement |
| C1-pt | author-control | pt | written by the author | not expected | the author's Portuguese version of the same control |
| C2-en | author-control | en | written by the author | expected | a control written by the author for this measurement |
| C2-pt | author-control | pt | written by the author | expected | the author's Portuguese version of the same control |
| C3-en | author-control | en | written by the author | not expected | a control written by the author for this measurement |
| C3-pt | author-control | pt | written by the author | not expected | the author's Portuguese version of the same control |
| C4-en | author-control | en | written by the author | expected | a control written by the author for this measurement |
| C4-pt | author-control | pt | written by the author | expected | the author's Portuguese version of the same control |
| C5-en | author-control | en | written by the author | expected | a control written by the author for this measurement |
| C5-pt | author-control | pt | written by the author | expected | the author's Portuguese version of the same control |
| C6-en | author-control | en | written by the author | not expected | a control written by the author for this measurement |
| C6-pt | author-control | pt | written by the author | not expected | the author's Portuguese version of the same control |
| F07-pt | feature-translation | pt | written by the author | reshaped | the author's Portuguese translation of docs/backlog/2026-09-25-storage-reclaim-notice.md@7c06573f |
| F04-pt | feature-translation | pt | written by the author | reshaped | the author's Portuguese translation of docs/backlog/2026-09-25-baseline-history-moves-need-separate-approval.md@7c06573f |
| F32-pt | feature-translation | pt | written by the author | declined | the author's Portuguese translation of docs/backlog/2026-11-12-roundfix-installs-and-updates-through-a-homebrew-tap.md@2e99852c |
| F00-pt | feature-translation | pt | written by the author | declined | the author's Portuguese translation of docs/backlog/2026-08-06-atomic-inbox-capture-helper.md@142e3232 |
| F05-pt | feature-translation | pt | written by the author | declined | the author's Portuguese translation of docs/backlog/2026-09-25-jev-gateway-baseline-before-routing.md@595876fc |
| F09-pt | feature-translation | pt | written by the author | as asked | the author's Portuguese translation of docs/backlog/2026-09-28-implement-spec-skill-delegates-to-roundfix.md@b92aefda |
| F08-pt | feature-translation | pt | written by the author | as asked | the author's Portuguese translation of docs/backlog/2026-09-28-blocking-review-after-archive-parks-publication.md@b92aefda |
| F11-pt | feature-translation | pt | written by the author | as asked | the author's Portuguese translation of docs/backlog/2026-09-28-queue-wide-limits-and-one-pending-question.md@b92aefda |
| F12-pt | feature-translation | pt | written by the author | as asked | the author's Portuguese translation of docs/backlog/2026-09-28-revalidate-later-specs-when-prerequisites-change.md@b92aefda |
| F10-pt | feature-translation | pt | written by the author | as asked | the author's Portuguese translation of docs/backlog/2026-09-28-queue-preparation-inventory.md@b92aefda |

## Answers

`jev-1.13.0`, first pass, all six questions, then the verdict from the three kept questions; then the three kept questions on the second pass and on `jev-1.14.0`.

| id | shape_open | unverified_assumption | market_question | new_capability | behavior_specified | worth_unsettled | verdict | second pass | jev-1.14.0 |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- | --- |
| F00 | 0.92 | 0.64 | 0.12 | 0.94 | 0.11 | 0.83 | indicated | 0.93 0.63 0.11 | 0.93 0.60 0.09 |
| F01 | 0.92 | 0.31 | 0.08 | 0.97 | 0.12 | 0.65 | indicated | 0.93 0.29 0.08 | 0.95 0.17 0.06 |
| F02 | 0.71 | 0.22 | 0.13 | 0.91 | 0.06 | 0.92 | indicated | 0.69 0.23 0.13 | 0.50 0.13 0.08 |
| F03 | 0.71 | 0.29 | 0.13 | 0.88 | 0.48 | 0.71 | indicated | 0.66 0.28 0.13 | 0.61 0.20 0.11 |
| F04 | 0.93 | 0.34 | 0.11 | 0.87 | 0.13 | 0.58 | indicated | 0.93 0.33 0.11 | 0.97 0.17 0.09 |
| F05 | 0.42 | 0.50 | 0.14 | 0.91 | 0.35 | 0.92 | undecided | 0.44 0.46 0.14 | 0.34 0.44 0.14 |
| F06 | 0.29 | 0.26 | 0.07 | 0.90 | 0.70 | 0.26 | not indicated | 0.27 0.27 0.07 | 0.23 0.18 0.04 |
| F07 | 0.65 | 0.38 | 0.10 | 0.94 | 0.40 | 0.43 | undecided | 0.64 0.38 0.10 | 0.66 0.38 0.09 |
| F08 | 0.24 | 0.43 | 0.10 | 0.86 | 0.50 | 0.44 | undecided | 0.25 0.45 0.10 | 0.25 0.39 0.06 |
| F09 | 0.25 | 0.33 | 0.06 | 0.39 | 0.53 | 0.31 | undecided | 0.24 0.30 0.06 | 0.20 0.32 0.03 |
| F10 | 0.43 | 0.57 | 0.12 | 0.92 | 0.38 | 0.46 | undecided | 0.44 0.55 0.11 | 0.50 0.59 0.10 |
| F11 | 0.44 | 0.50 | 0.09 | 0.88 | 0.35 | 0.38 | undecided | 0.48 0.52 0.09 | 0.56 0.52 0.06 |
| F12 | 0.25 | 0.27 | 0.07 | 0.72 | 0.66 | 0.35 | not indicated | 0.27 0.27 0.07 | 0.21 0.15 0.04 |
| F13 | 0.85 | 0.55 | 0.09 | 0.93 | 0.23 | 0.48 | indicated | 0.85 0.52 0.08 | 0.87 0.53 0.05 |
| F14 | 0.59 | 0.37 | 0.09 | 0.26 | 0.46 | 0.30 | undecided | 0.61 0.36 0.09 | 0.63 0.27 0.06 |
| F15 | 0.26 | 0.31 | 0.08 | 0.95 | 0.63 | 0.86 | undecided | 0.27 0.31 0.07 | 0.31 0.18 0.05 |
| F16 | 0.40 | 0.44 | 0.05 | 0.61 | 0.20 | 0.95 | undecided | 0.42 0.47 0.05 | 0.34 0.53 0.03 |
| F17 | 0.29 | 0.18 | 0.06 | 0.91 | 0.80 | 0.26 | not indicated | 0.27 0.18 0.07 | 0.27 0.14 0.03 |
| F18 | 0.85 | 0.72 | 0.26 | 0.97 | 0.13 | 0.92 | indicated | 0.87 0.71 0.25 | 0.81 0.74 0.23 |
| F19 | 0.94 | 0.46 | 0.15 | 0.93 | 0.10 | 0.47 | indicated | 0.94 0.46 0.14 | 0.98 0.36 0.17 |
| F20 | 0.68 | 0.70 | 0.11 | 0.87 | 0.09 | 0.95 | indicated | 0.70 0.69 0.10 | 0.50 0.70 0.12 |
| F21 | 0.79 | 0.73 | 0.10 | 0.81 | 0.18 | 0.94 | indicated | 0.79 0.73 0.11 | 0.70 0.79 0.07 |
| F22 | 0.27 | 0.18 | 0.06 | 0.71 | 0.84 | 0.18 | not indicated | 0.30 0.19 0.06 | 0.17 0.10 0.03 |
| F23 | 0.53 | 0.32 | 0.07 | 0.94 | 0.20 | 0.95 | undecided | 0.49 0.30 0.07 | 0.27 0.28 0.05 |
| F24 | 0.36 | 0.19 | 0.05 | 0.92 | 0.70 | 0.15 | undecided | 0.32 0.20 0.06 | 0.35 0.17 0.04 |
| F25 | 0.15 | 0.16 | 0.07 | 0.96 | 0.84 | 0.12 | not indicated | 0.14 0.17 0.08 | 0.06 0.06 0.03 |
| F26 | 0.90 | 0.14 | 0.05 | 0.22 | 0.16 | 0.64 | indicated | 0.88 0.14 0.05 | 0.88 0.08 0.03 |
| F27 | 0.33 | 0.07 | 0.04 | 0.20 | 0.85 | 0.11 | undecided | 0.33 0.07 0.04 | 0.36 0.04 0.02 |
| F28 | 0.56 | 0.10 | 0.04 | 0.39 | 0.79 | 0.17 | undecided | 0.57 0.10 0.04 | 0.49 0.06 0.02 |
| F29 | 0.16 | 0.13 | 0.06 | 0.76 | 0.84 | 0.12 | not indicated | 0.14 0.13 0.06 | 0.19 0.07 0.04 |
| F30 | 0.16 | 0.56 | 0.11 | 0.58 | 0.76 | 0.48 | undecided | 0.16 0.57 0.11 | 0.09 0.56 0.08 |
| F31 | 0.94 | 0.48 | 0.10 | 0.65 | 0.13 | 0.87 | indicated | 0.94 0.45 0.10 | 0.99 0.44 0.07 |
| F32 | 0.58 | 0.55 | 0.31 | 0.93 | 0.28 | 0.74 | undecided | 0.60 0.55 0.32 | 0.42 0.20 0.34 |
| F33 | 0.55 | 0.30 | 0.10 | 0.89 | 0.73 | 0.11 | undecided | 0.58 0.29 0.10 | 0.66 0.14 0.05 |
| F34 | 0.66 | 0.11 | 0.06 | 0.81 | 0.55 | 0.09 | undecided | 0.67 0.10 0.06 | 0.69 0.04 0.03 |
| F35 | 0.46 | 0.29 | 0.09 | 0.89 | 0.56 | 0.74 | undecided | 0.46 0.29 0.08 | 0.49 0.16 0.05 |
| F36 | 0.23 | 0.13 | 0.05 | 0.88 | 0.81 | 0.24 | not indicated | 0.23 0.13 0.05 | 0.19 0.06 0.02 |
| F38 | 0.17 | 0.28 | 0.10 | 0.95 | 0.90 | 0.13 | not indicated | 0.17 0.28 0.10 | 0.13 0.18 0.07 |
| X01 | 0.28 | 0.11 | 0.05 | 0.12 | 0.73 | 0.13 | not indicated | 0.29 0.12 0.05 | 0.21 0.08 0.03 |
| X02 | 0.83 | 0.10 | 0.05 | 0.09 | 0.69 | 0.23 | indicated | 0.84 0.09 0.05 | 0.88 0.10 0.04 |
| X03 | 0.93 | 0.08 | 0.05 | 0.18 | 0.38 | 0.47 | indicated | 0.93 0.08 0.05 | 0.98 0.05 0.03 |
| X04 | 0.84 | 0.14 | 0.04 | 0.09 | 0.48 | 0.54 | indicated | 0.85 0.15 0.04 | 0.94 0.10 0.02 |
| X05 | 0.79 | 0.12 | 0.05 | 0.21 | 0.62 | 0.19 | indicated | 0.78 0.12 0.05 | 0.76 0.10 0.03 |
| X06 | 0.30 | 0.11 | 0.05 | 0.41 | 0.61 | 0.28 | not indicated | 0.32 0.11 0.05 | 0.24 0.07 0.04 |
| X07 | 0.67 | 0.16 | 0.03 | 0.07 | 0.39 | 0.28 | undecided | 0.67 0.15 0.03 | 0.62 0.15 0.02 |
| X08 | 0.81 | 0.14 | 0.05 | 0.18 | 0.53 | 0.27 | indicated | 0.81 0.16 0.05 | 0.89 0.12 0.04 |
| X09 | 0.76 | 0.15 | 0.05 | 0.45 | 0.40 | 0.34 | indicated | 0.75 0.15 0.05 | 0.84 0.13 0.03 |
| X10 | 0.86 | 0.15 | 0.04 | 0.12 | 0.24 | 0.52 | indicated | 0.87 0.15 0.04 | 0.94 0.12 0.02 |
| X11 | 0.40 | 0.16 | 0.05 | 0.41 | 0.82 | 0.18 | undecided | 0.33 0.16 0.05 | 0.37 0.10 0.04 |
| X12 | 0.58 | 0.15 | 0.05 | 0.62 | 0.69 | 0.13 | undecided | 0.59 0.17 0.05 | 0.41 0.14 0.05 |
| X13 | 0.67 | 0.10 | 0.06 | 0.22 | 0.74 | 0.22 | undecided | 0.69 0.11 0.06 | 0.68 0.09 0.05 |
| X14 | 0.88 | 0.10 | 0.04 | 0.21 | 0.43 | 0.19 | indicated | 0.87 0.10 0.04 | 0.95 0.05 0.02 |
| X15 | 0.85 | 0.09 | 0.06 | 0.07 | 0.67 | 0.18 | indicated | 0.85 0.10 0.05 | 0.87 0.09 0.04 |
| X16 | 0.44 | 0.17 | 0.05 | 0.14 | 0.81 | 0.20 | undecided | 0.43 0.16 0.05 | 0.43 0.14 0.04 |
| X17 | 0.27 | 0.07 | 0.03 | 0.08 | 0.79 | 0.14 | not indicated | 0.28 0.07 0.03 | 0.28 0.05 0.02 |
| X18 | 0.34 | 0.12 | 0.04 | 0.11 | 0.72 | 0.46 | undecided | 0.31 0.11 0.04 | 0.17 0.08 0.02 |
| X19 | 0.22 | 0.07 | 0.03 | 0.14 | 0.88 | 0.12 | not indicated | 0.23 0.08 0.04 | 0.16 0.03 0.02 |
| X20 | 0.20 | 0.11 | 0.05 | 0.62 | 0.82 | 0.13 | not indicated | 0.21 0.09 0.05 | 0.14 0.06 0.03 |
| N01 | 0.73 | 0.15 | 0.08 | 0.61 | 0.48 | 0.31 | indicated | 0.73 0.15 0.08 | 0.75 0.15 0.07 |
| N02 | 0.85 | 0.07 | 0.03 | 0.17 | 0.22 | 0.37 | indicated | 0.85 0.08 0.04 | 0.91 0.05 0.02 |
| N03 | 0.85 | 0.11 | 0.06 | 0.39 | 0.39 | 0.29 | indicated | 0.85 0.12 0.06 | 0.87 0.09 0.05 |
| N04 | 0.82 | 0.33 | 0.08 | 0.77 | 0.24 | 0.42 | indicated | 0.83 0.33 0.08 | 0.83 0.31 0.06 |
| N05 | 0.58 | 0.10 | 0.06 | 0.18 | 0.46 | 0.18 | undecided | 0.61 0.11 0.06 | 0.74 0.07 0.04 |
| N06 | 0.28 | 0.25 | 0.05 | 0.67 | 0.67 | 0.15 | not indicated | 0.30 0.24 0.05 | 0.27 0.15 0.03 |
| N07 | 0.92 | 0.10 | 0.05 | 0.48 | 0.14 | 0.86 | indicated | 0.92 0.09 0.06 | 0.96 0.06 0.03 |
| N08 | 0.90 | 0.09 | 0.03 | 0.38 | 0.12 | 0.38 | indicated | 0.90 0.09 0.03 | 0.93 0.08 0.02 |
| N09 | 0.35 | 0.13 | 0.06 | 0.77 | 0.50 | 0.45 | undecided | 0.34 0.13 0.06 | 0.36 0.08 0.04 |
| N10 | 0.76 | 0.53 | 0.16 | 0.89 | 0.06 | 0.92 | indicated | 0.75 0.52 0.16 | 0.69 0.51 0.12 |
| C1-en | 0.42 | 0.06 | 0.04 | 0.04 | 0.86 | 0.09 | undecided | 0.45 0.06 0.04 | 0.35 0.06 0.04 |
| C1-pt | 0.56 | 0.06 | 0.04 | 0.04 | 0.86 | 0.09 | undecided | 0.50 0.05 0.04 | 0.45 0.05 0.04 |
| C2-en | 0.97 | 0.95 | 0.93 | 0.95 | 0.03 | 0.92 | indicated | 0.97 0.95 0.93 | 0.98 0.98 0.97 |
| C2-pt | 0.97 | 0.95 | 0.93 | 0.95 | 0.04 | 0.91 | indicated | 0.97 0.95 0.92 | 0.99 0.98 0.97 |
| C3-en | 0.12 | 0.17 | 0.07 | 0.96 | 0.90 | 0.10 | not indicated | 0.12 0.19 0.07 | 0.06 0.11 0.04 |
| C3-pt | 0.15 | 0.20 | 0.07 | 0.96 | 0.89 | 0.10 | not indicated | 0.15 0.22 0.07 | 0.07 0.11 0.04 |
| C4-en | 0.98 | 0.84 | 0.17 | 0.95 | 0.04 | 0.79 | indicated | 0.98 0.83 0.17 | 0.99 0.89 0.16 |
| C4-pt | 0.98 | 0.78 | 0.16 | 0.95 | 0.05 | 0.78 | indicated | 0.98 0.79 0.16 | 0.99 0.80 0.13 |
| C5-en | 0.64 | 0.96 | 0.95 | 0.95 | 0.07 | 0.92 | indicated | 0.69 0.96 0.94 | 0.61 0.99 0.98 |
| C5-pt | 0.75 | 0.96 | 0.96 | 0.95 | 0.08 | 0.91 | indicated | 0.77 0.96 0.96 | 0.72 0.99 0.98 |
| C6-en | 0.08 | 0.17 | 0.08 | 0.91 | 0.85 | 0.15 | not indicated | 0.08 0.17 0.08 | 0.05 0.11 0.06 |
| C6-pt | 0.09 | 0.16 | 0.07 | 0.90 | 0.86 | 0.14 | not indicated | 0.08 0.14 0.08 | 0.05 0.09 0.05 |
| F07-pt | 0.71 | 0.32 | 0.09 | 0.93 | 0.50 | 0.35 | indicated | 0.69 0.31 0.08 | 0.78 0.30 0.06 |
| F04-pt | 0.93 | 0.31 | 0.11 | 0.86 | 0.14 | 0.52 | indicated | 0.93 0.36 0.11 | 0.98 0.18 0.09 |
| F32-pt | 0.42 | 0.58 | 0.24 | 0.94 | 0.24 | 0.81 | undecided | 0.44 0.59 0.26 | 0.30 0.33 0.24 |
| F00-pt | 0.93 | 0.73 | 0.10 | 0.94 | 0.12 | 0.82 | indicated | 0.94 0.74 0.10 | 0.97 0.86 0.06 |
| F05-pt | 0.51 | 0.50 | 0.15 | 0.91 | 0.34 | 0.92 | undecided | 0.49 0.48 0.15 | 0.44 0.49 0.10 |
| F09-pt | 0.28 | 0.34 | 0.07 | 0.38 | 0.52 | 0.31 | undecided | 0.28 0.33 0.07 | 0.19 0.30 0.04 |
| F08-pt | 0.27 | 0.42 | 0.09 | 0.83 | 0.42 | 0.36 | undecided | 0.28 0.37 0.09 | 0.23 0.27 0.06 |
| F11-pt | 0.63 | 0.50 | 0.12 | 0.87 | 0.35 | 0.37 | undecided | 0.65 0.48 0.12 | 0.79 0.37 0.07 |
| F12-pt | 0.30 | 0.30 | 0.08 | 0.81 | 0.64 | 0.27 | not indicated | 0.31 0.30 0.08 | 0.21 0.19 0.04 |
| F10-pt | 0.48 | 0.54 | 0.10 | 0.91 | 0.44 | 0.42 | undecided | 0.47 0.53 0.10 | 0.49 0.57 0.05 |

## The author's controls and translations

Each text is one item's body; its first heading is the item's title.

### C1-en

````markdown
# The status command prints the table header twice on an empty list

## Symptom

With no Runs recorded, `status` prints the column header two times and then exits 0. Reproduced on the current release with an empty database.

## Where

The table printer used by `status`.

## Expected

One header line, then the line that says no Run was found. A regression test covers the empty list.

## Evidence

Reproduced three times on a clean install; the second header comes from the footer branch.
````

### C1-pt

````markdown
# O comando status imprime o cabeçalho da tabela duas vezes com a lista vazia

## Sintoma

Sem nenhuma Run registrada, o `status` imprime o cabeçalho das colunas duas vezes e depois sai com 0. Reproduzido na versão atual com o banco de dados vazio.

## Onde

O impressor de tabelas usado pelo `status`.

## Esperado

Uma linha de cabeçalho e depois a linha que diz que nenhuma Run foi encontrada. Um teste de regressão cobre a lista vazia.

## Evidência

Reproduzido três vezes numa instalação limpa; o segundo cabeçalho vem do ramo do rodapé.
````

### C2-en

````markdown
# Offer a hosted dashboard with team accounts and paid plans

## Opportunity

Maybe the product could offer a hosted web dashboard where managers follow the agents' work from the browser, with team accounts and paid plans.

## Value

Companies would probably pay for it, and it could become the main product line. Nobody has asked for it yet.

## Shape

Not decided: a hosted service, a self-hosted server or a static report are all possible.
````

### C2-pt

````markdown
# Oferecer um painel hospedado com contas de equipe e planos pagos

## Oportunidade

Talvez o produto possa oferecer um painel web hospedado em que os gestores acompanham o trabalho dos agentes pelo navegador, com contas de equipe e planos pagos.

## Valor

As empresas provavelmente pagariam por isso, e ele poderia virar a principal linha de produto. Ninguém pediu isso ainda.

## Forma

Não decidida: um serviço hospedado, um servidor instalado pelo próprio cliente ou um relatório estático são todos possíveis.
````

### C3-en

````markdown
# The list command gains a --since flag

## Opportunity

Operators filter the Run list by hand with grep on the date column.

## Value

The maintainer decided on the flag after three operators reported the same workaround.

## Shape

`list --since <YYYY-MM-DD>` keeps only the Runs started on or after that date, in text and JSON; an invalid date exits 2 with a usage error; without the flag the output keeps its bytes.
````

### C3-pt

````markdown
# O comando list ganha a flag --since

## Oportunidade

Os operadores filtram a lista de Runs à mão, com grep na coluna de data.

## Valor

O mantenedor decidiu pela flag depois que três operadores relataram o mesmo contorno.

## Forma

`list --since <AAAA-MM-DD>` mantém só as Runs iniciadas nessa data ou depois dela, em texto e em JSON; uma data inválida sai com 2 e erro de uso; sem a flag a saída mantém os mesmos bytes.
````

### C4-en

````markdown
# Teams can see each other's Run results

## Opportunity

People in the same team cannot see what a colleague's Runs did without asking for a screenshot.

## Value

Shared results would cut the back-and-forth in reviews.

## Shape

Open. It could be a web page each Run publishes, a chat bot that posts a summary, an export file committed to the repository, or a shared database; each has a different cost and privacy profile, and nothing has been chosen.
````

### C4-pt

````markdown
# As equipes podem ver os resultados das Runs umas das outras

## Oportunidade

Pessoas da mesma equipe não conseguem ver o que as Runs de um colega fizeram sem pedir uma captura de tela.

## Valor

Resultados compartilhados reduziriam o vai e vem nas revisões.

## Forma

Em aberto. Pode ser uma página web que cada Run publica, um bot de chat que posta um resumo, um arquivo de exportação gravado no repositório ou um banco de dados compartilhado; cada um tem custo e perfil de privacidade diferentes, e nada foi escolhido.
````

### C5-en

````markdown
# A marketplace of community plugins

## Opportunity

Other developer tools grew through plugin marketplaces, so this one should have a marketplace where the community publishes and sells plugins.

## Value

A marketplace would attract outside contributors and set the product apart from its competitors. No user has been asked, and no competitor's marketplace has been studied.

## Shape

A public registry with ratings and paid listings.
````

### C5-pt

````markdown
# Um marketplace de plugins da comunidade

## Oportunidade

Outras ferramentas de desenvolvimento cresceram com marketplaces de plugins, então esta deveria ter um marketplace em que a comunidade publica e vende plugins.

## Valor

Um marketplace atrairia colaboradores de fora e diferenciaria o produto dos concorrentes. Nenhum usuário foi consultado, e o marketplace de nenhum concorrente foi estudado.

## Forma

Um registro público com avaliações e anúncios pagos.
````

### C6-en

````markdown
# The cache directory follows a configured size bound

## Opportunity

The build cache reached 26 GB and filled the disk twice in one day, measured with `du` on both occasions.

## Value

A bounded cache stops the disk from filling during unattended work.

## Shape

Decided by the maintainer: a `cache.max_bytes` key in the user configuration, default 10 GB; before each build the oldest entries are removed until the cache fits; the command prints one line naming the bytes it removed.
````

### C6-pt

````markdown
# O diretório de cache segue um limite de tamanho configurado

## Oportunidade

O cache de build chegou a 26 GB e encheu o disco duas vezes no mesmo dia, medido com `du` nas duas ocasiões.

## Valor

Um cache limitado impede que o disco encha durante o trabalho sem supervisão.

## Forma

Decidido pelo mantenedor: uma chave `cache.max_bytes` na configuração do usuário, com padrão de 10 GB; antes de cada build as entradas mais antigas são removidas até o cache caber; o comando imprime uma linha com os bytes que removeu.
````

### F07-pt

````markdown
# Avisar o operador quando o armazenamento das Runs pode ser recuperado

## Oportunidade

Os artefatos de Run e o Run Database crescem a cada Run; só o `gc` informa os bytes recuperáveis, e só quando uma pessoa o executa (4,4 GB com 619 MB recuperáveis depois de um dia).

## Valor

O operador fica sabendo do espaço recuperável onde ele já olha, sem nenhuma exclusão automática.

## Forma

Reaproveitar a medição de armazenamento e imprimir uma linha de limite no `runs list` ou no Doctor, apontando para o `roundfix gc`.

Evidência: secondbrain `inbox/roundfix/2026-09-17-os-artefatos-de-run-crescem-e-so-um-humano-recupera.md`.
````

### F04-pt

````markdown
# Aprovar as movimentações de histórico da Baseline separadamente das atualizações gerenciadas

## Oportunidade

O `baseline update` amarra uma atualização de bytes gerenciados e a realocação de ADRs do próprio repositório para o histórico sob um único Plan Digest, de modo que um adotante não consegue aceitar uma sem a outra; o Fiscus adiou um plano inteiro por causa disso.

## Valor

Os adotantes aceitam as atualizações gerenciadas sem realocar decisões que não revisaram, e as citações de caminho deixam de quebrar em silêncio.

## Forma

Aprovação (ou digest) separada para as `HistoryMoves`, ou pelo menos mostrar o impacto delas nas citações antes da aprovação; precisa de uma decisão de desenho diante da ADR-0073 e da ADR-0103.

Evidência: secondbrain `inbox/roundfix/2026-09-17-um-digest-so-prende-refresh-gerenciado-e-mudanca-de-documento-do-repositorio.md`.
````

### F32-pt

````markdown
# O Roundfix é instalado e atualizado por um tap do Homebrew

## Oportunidade

Registrado a pedido do mantenedor, repassado pelo jarvis na mensagem do ai-memory `01a1d27a-e25a-7180-86e3-69dc7cfa3ff3` de 2026-11-12 (22:35 BRT); a ideia estava no Todoist do mantenedor desde 2026-10-20. Um tap do Homebrew do próprio Roundfix, com uma fórmula, permitiria a um usuário de macOS ou Linux instalar com `brew install` e atualizar com `brew upgrade`, como o OpenCode 2 faz com o tap dele.

## Valor

As máquinas do mantenedor atualizam todo o conjunto de ferramentas pelo Homebrew e pelo `dev-upgrade`; uma fórmula traria o Roundfix para essa rotina, no lugar da instalação separada pelo npm.

## Forma

Um repositório de tap com uma fórmula que baixa o artefato de release construído pelo workflow de release existente para a plataforma, atualizada pelo workflow de release a cada tag. Não vinculante; ainda sem Spec.
````

### F00-pt

````markdown
# Adicionar um auxiliar atômico para a captura durável no inbox

## Oportunidade

As sessões da frota poderiam capturar uma Inbox Entry válida pelo contrato com uma única ação observável, em vez de criar, preparar, revisar, gravar e cronometrar a entrada à mão. Esta intenção chegou por `inbox/roundfix/2026-08-06-automate-inbox-capture-durability.md` no Secondbrain.

## Valor

Uma ação atômica de captura reduziria o risco de perder a meta de durabilidade de um minuto e deixaria o caminho, o commit do brain e o tempo decorrido disponíveis para a sessão que captura, sem reconstruí-los a partir do histórico do shell.

## Forma

Um desenho futuro poderia validar o contrato da Inbox Entry, criar só a entrada de destino, gravar só esse caminho relativo ao brain e informar o caminho, o hash do commit e os segundos entre a captura e a durabilidade. Esta forma não é vinculante; a Spec responsável deve escolher se isso pertence a uma skill, a um auxiliar ou a uma superfície do produto.
````

### F05-pt

````markdown
# Medir o uso de tokens dos agentes pelo jev-gateway antes de ligar o roteamento

## Oportunidade

O Roundfix não sabe dizer quantas requisições e quantos tokens uma Task custa por runtime. O jev-gateway fica entre o agente ACP e o provedor dele; com `--routing off` ele só conta, e com o roteamento ligado ele envia as decisões de escolha de ferramenta ao Jev. Um benchmark da comunidade relata grandes economias em correções de bugs com GPT e regressões em modelos Claude, então o ganho depende do runtime.

## Valor

Uma linha de base medida por runtime e por tipo de Task, e depois uma decisão sobre o roteamento baseada em evidência. O roteamento envia esquemas de ferramentas e contexto à TypeSafe, então também precisa de uma decisão de compartilhamento de dados.

## Forma

Rodar o runtime codex pelo gateway em modo de linha de base, sem mudar a configuração pessoal de agentes do mantenedor, registrar as contagens por Run e comparar com o roteamento ligado só depois que o trabalho de testes determinísticos (D3) remover as falhas por carga. Mostrar a comparação ao mantenedor antes de adotar o roteamento.

Evidência: secondbrain `raw/roundfix/2026-09-25-sequencia-de-eficiencia.md` e `raw/web/2026-09-25-exa-roundfix-eficiencia-e-custo.json`; https://github.com/vinilana/jev-gateway.
````

### F09-pt

````markdown
# O ponto de entrada implement-spec delega ao Roundfix

## Oportunidade

Core Feature 8 da Spec 0127: a skill própria `implement-spec` ainda carrega a sua própria orientação de loop, em vez de delegar a implementação e a Verification ao `roundfix implement`/`deliver`.

## Valor

Um único loop de implementação; o Supervisor nunca escreve código nem testes.

## Forma

Reescrever a skill para preparar a fila e passar o trabalho ao Roundfix, sem copiar um script de janela.

Evidência: trazido da Spec 0127 quando essa Spec de portfólio foi aposentada em 2026-09-28 (as funcionalidades entregues saíram nas Specs que a substituição dela nomeia).
````

### F08-pt

````markdown
# Uma revisão bloqueante depois do arquivamento estaciona a publicação e pede uma nova Spec

## Oportunidade

Core Feature 9 da Spec 0127: quando uma revisão pré-PR bloqueia depois que a Spec foi arquivada no branch, o caminho de correção é improvisado.

## Valor

A Spec arquivada fica intacta, a publicação estaciona, e a correção é escrita como uma nova Spec, com autoridade e portão próprios; o candidato é revisado de novo depois que ela é arquivada.

## Forma

Um motivo de estacionamento `review-blocked` que nomeia a Spec corretiva necessária; nenhum orçamento de ciclos a autoriza em silêncio.

Evidência: trazido da Spec 0127 quando essa Spec de portfólio foi aposentada em 2026-09-28 (as funcionalidades entregues saíram nas Specs que a substituição dela nomeia).
````

### F11-pt

````markdown
# Limites para a fila inteira e uma única pergunta pendente

## Oportunidade

Core Features 4 e 5 da Spec 0127: a aprovação e os limites acompanham cada Spec, mas os limites de tempo, gasto, concorrência e correção para a fila inteira não são decididos nem aplicados, e nada garante que só uma pergunta ao mantenedor fique pendente (o silêncio ou uma opção recomendada nunca pode contar como resposta).

## Valor

Uma fila sem supervisão para num limite explícito, e não num limite inferido, e o mantenedor enfrenta uma decisão de cada vez.

## Forma

Declarar os limites da fila no registro da fila; parar as ações dependentes quando uma autorização faltar ou mudar; um único registro durável de pergunta pendente.

Evidência: trazido da Spec 0127 quando essa Spec de portfólio foi aposentada em 2026-09-28 (as funcionalidades entregues saíram nas Specs que a substituição dela nomeia).
````

### F12-pt

````markdown
# As Specs seguintes da fila são revalidadas quando um pré-requisito muda as premissas delas

## Oportunidade

Core Feature 6 da Spec 0127: depois que um item anterior da fila entra no main, os itens seguintes rodam sobre o novo main sem conferir se as premissas com que foram escritos (arquivos tocados, assinaturas, caminhos governados) ainda valem — a Onda 2 esbarrou nisso quando a Spec 0175 foi escrita antes de a 0173 mudar o código de entrega que ela edita.

## Valor

Uma contradição, uma nova exigência de autoridade ou uma nova ação irreversível vira um bloqueio durável, e não uma decisão improvisada dentro de uma Run.

## Forma

Rodar de novo o `spec check` e a auditoria de caminhos governados para cada item seguinte no início dele; estacionar com um bloqueio nomeado em caso de falha.

Evidência: trazido da Spec 0127 quando essa Spec de portfólio foi aposentada em 2026-09-28 (as funcionalidades entregues saíram nas Specs que a substituição dela nomeia).
````

### F10-pt

````markdown
# A preparação da fila inventaria o que está aprovado para rodar

## Oportunidade

Core Feature 1 da Spec 0127: o `deliver start` recebe slugs; nada inventaria as Specs ativas, a intenção do Backlog, os Findings e o Inbox, registra as dependências e as decisões pendentes e nomeia quais Specs estão de fato aprovadas.

## Valor

O operador (ou o Supervisor) vê uma única fila preparada, com os pré-requisitos e as decisões em aberto, antes de qualquer coisa rodar; um inventário nunca é autoridade de implementação.

## Forma

Um `deliver plan` só de leitura (ou `deliver start --dry-run`) sobre as raízes do ciclo de vida, reaproveitando o `spec check` e a resolução de autorizações.

Evidência: trazido da Spec 0127 (`docs/history/specs/0127-durable-unattended-spec-workflow/_prd.md`) quando essa Spec de portfólio foi aposentada em 2026-09-28 (as funcionalidades entregues saíram nas Specs que a substituição dela nomeia).
````
