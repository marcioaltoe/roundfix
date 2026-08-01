#!/usr/bin/env ruby

require "yaml"

root = File.expand_path("../../../../../..", __dir__)
workflow = YAML.load_file(File.join(root, ".github/workflows/release.yml"))
runbook = File.read(File.join(root, "docs/user-guide/release-runbook.md"))

scripts = workflow.fetch("jobs").fetch("release").fetch("steps")
  .map { |step| step["run"] }.compact
  .join("\n")
runtime_prefixes = scripts.scan(/::(?:error|warning)::([a-z]+):/).flatten.uniq.sort
documented_prefixes = runbook.scan(/^\| `([a-z]+:)` \|/).flatten
  .map { |prefix| prefix.delete_suffix(":") }.uniq.sort
missing = runtime_prefixes - documented_prefixes

puts "WORKFLOW_PREFIXES=#{runtime_prefixes.join(",")}"
puts "RUNBOOK_PREFIXES=#{documented_prefixes.join(",")}"

if missing.any?
  warn "FAIL: runbook omits workflow failure prefix(es): #{missing.map { |prefix| "#{prefix}:" }.join(", ")}"
  exit 1
end

puts "PASS: runbook documents every workflow failure prefix"
