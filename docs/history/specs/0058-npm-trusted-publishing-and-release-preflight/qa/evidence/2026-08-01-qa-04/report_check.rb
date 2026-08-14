#!/usr/bin/env ruby

require "yaml"

report = File.expand_path("../../qa-report-2026-08-01-04.md", __dir__)
text = File.read(report)
frontmatter = YAML.load(text.split(/^---\s*$/, 3).fetch(1))

raise "report is not closed partial" unless frontmatter["status"] == "closed" && frontmatter["verdict"] == "partial"
raise "blocked counts mismatch" unless frontmatter["rows_blocked_environment"] == 2 && frontmatter["rows_blocked_finding"] == 0

rows = text.lines.select { |line| line.match?(/^\| \d{2} \|/) }
statuses = rows.map { |line| line.split("|")[-2].strip }
raise "row count mismatch: #{rows.length}" unless rows.length == 29
raise "pass count mismatch" unless statuses.count("pass") == 27
raise "environment-block count mismatch" unless statuses.count { |status| status.start_with?("blocked (environment:") } == 2
raise "nonterminal or unexpected row status" unless statuses.all? { |status| status == "pass" || status.start_with?("blocked (environment:") }

spec_root = File.expand_path("../../..", __dir__)
evidence_paths = text.scan(/`(qa\/evidence\/2026-08-01-qa-04\/[^`]+)`/).flatten.uniq
raise "no evidence paths found" if evidence_paths.empty?
missing = evidence_paths.reject { |path| File.exist?(File.join(spec_root, path)) }
raise "missing evidence: #{missing.join(", ")}" unless missing.empty?

raise "final verdict is not actionable" unless text.match?(/## Final verdict\n\nPartial:.*Integrate and push/m)

puts "PASS: report is closed partial with 29 terminal rows"
puts "PASS: 27 pass and 2 environment-blocked rows match frontmatter"
puts "PASS: every run-local evidence path resolves"
