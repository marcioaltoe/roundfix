---
task: task_06
spec: 0078-roundfix-asks-for-the-review
status: completed
type: qa
complexity: medium
---

# Task 06: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_05 settles `completed`.

This Spec exists because a correctness fix left the loop unable to advance, so
this gate's most valuable rows are the ones proving the loop advances now and
that it consumes exactly what the configuration says it will.

## Requirements

1. MUST run only after task_05 settles `completed`.
2. MUST observe that a Round which pushes publishes exactly one review request,
   including the Round whose Final Push is followed by the artifact-only docs
   commit.
3. MUST observe that the same head asked twice publishes once, rather than
   accepting the Task's claim.
4. MUST observe all four rows of the Preflight coherence table, including both
   refusals and their exit codes.
5. MUST confirm `fetch` publishes no request under any configuration.
6. MUST confirm no automatic retry, backoff, or capacity wait was introduced,
   since all are explicitly out of scope.
7. MUST confirm the repository's own committed configuration pair is coherent,
   because an incoherent pair strands every later Run.
8. MUST classify any finding by user impact and record typed blocked-row counts.

## Acceptance Criteria

- [ ] The gate runs only after task_05 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the one-request-per-Round observation end to end.
- [ ] The report records both Preflight refusals observed independently.

## Verification

- `test -f docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/qa-report-2026-08-05.md`
  — expected: exit 0; this Spec's dated QA report exists.
- `rtk ruby -e 'require "yaml"; require "date"; def metadata(text, path); match = text.match(/\A---\n(.*?)\n---\n/m); raise "#{path}: missing frontmatter" unless match; YAML.safe_load(match[1], permitted_classes: [Date], aliases: false); end; report_path = "docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/qa-report-2026-08-05.md"; task_path = "docs/specs/_archived/0078-roundfix-asks-for-the-review/task_05.md"; report = File.read(report_path); meta = metadata(report, report_path); raise "#{report_path}: spec" unless meta["spec"] == "0078-roundfix-asks-for-the-review"; raise "#{report_path}: verdict" unless %w[pass partial fail].include?(meta["verdict"]); %w[rows_blocked_environment rows_blocked_finding].each { |key| raise "#{report_path}: #{key}" unless meta[key].is_a?(Integer) && meta[key] >= 0 }; task = metadata(File.read(task_path), task_path); raise "#{task_path}: status" unless task["status"] == "completed"; results = report[/^## Results\n(.*?)(?=^## )/m, 1] or raise "#{report_path}: results"; required = { "R05" => ["Four-row Preflight table and both exit-2 refusals", "both equal pairs refused"], "R06" => ["One request after a Round", "one new-head request"], "R07" => ["Artifact-only descendant produces no second request", "one request for fix head"], "R10" => ["remains read-only under every pair", "zero requests"], "R11" => ["Explicit refusal ends Review Skipped with no retry", "no retry/second request"], "R12" => ["Request count is bounded by Round cap", "once inside a Round"] }; required.each { |id, (story, evidence)| line = results.lines.find { |candidate| candidate.start_with?("| #{id} |") }; raise "#{report_path}: missing #{id}" unless line; cells = line.split("|", -1)[1...-1].map(&:strip); raise "#{report_path}: malformed #{id}" unless cells.length == 5; raise "#{report_path}: #{id} story" unless cells[1].include?(story); raise "#{report_path}: #{id} status" unless cells[3] == "pass" || cells[3].match?(/\Ablocked \(environment: .+\)\z/); raise "#{report_path}: #{id} evidence" unless cells[4].include?(evidence) }'`
  — expected: exit 0; Task 05 is completed, the report carries this Spec's
  verdict and typed blocked-row counts, and Results rows R05–R07 and R10–R12
  record the two Preflight refusals, one-request-per-Round behavior, artifact
  descendant behavior, `fetch` exemption, and no-retry result.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0054, ADR-0080, ADR-0091.
