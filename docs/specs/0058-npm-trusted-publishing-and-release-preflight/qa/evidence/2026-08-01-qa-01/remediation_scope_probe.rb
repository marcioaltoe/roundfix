#!/usr/bin/env ruby

require "json"
require "open3"

ROOT = File.expand_path("../../../../../..", __dir__)
SPEC = File.join(ROOT, "docs/specs/0058-npm-trusted-publishing-and-release-preflight")

def assert(condition, message)
  raise message unless condition

  puts "PASS: #{message}"
end

def git(*args)
  stdout, stderr, status = Open3.capture3(
    "git", "-c", "core.fsmonitor=false", *args, chdir: ROOT,
  )
  raise "git #{args.join(" ")} failed: #{stderr}" unless status.success?

  stdout
end

tasks = (1..7).map do |number|
  File.read(File.join(SPEC, format("task_%02d.md", number)))
end
assert(tasks.all? { |task| task.match?(/^status: completed$/) }, "all seven Task files are completed")
assert(tasks.all? { |task| task.include?("## Result") && task.include?("Acceptance evidence") }, "all seven Task Results carry acceptance evidence")

prd = File.read(File.join(SPEC, "_prd.md"))
techspec = File.read(File.join(SPEC, "_techspec.md"))
constraints = ["Identifier strategy", "Authentication and HTTP", "Active ADR obligations", "Tooling authority"]
[prd, techspec].each_with_index do |artifact, index|
  label = index.zero? ? "PRD" : "TechSpec"
  assert(constraints.all? { |name| artifact.include?(name) }, "#{label} accounts for all Project Constraints")
  assert(artifact.include?("docs/agents/domain.md"), "#{label} cites domain guidance")
  assert(artifact.include?("docs/agents/agent-instructions.md"), "#{label} cites tooling and authentication guidance")
end
assert(prd.include?('authorizes changes to exactly `.github/workflows/release.yml`'), "PRD expressly authorizes the exact protected workflow")
assert(techspec.include?('bounded files: `.github/workflows/release.yml`'), "TechSpec repeats the exact protected-tooling boundary")

authorization = git("show", "397227ff:docs/specs/0058-npm-trusted-publishing-and-release-preflight/_prd.md")
assert(authorization.include?('authorizes changes to exactly `.github/workflows/release.yml`'), "authorization exists before every tooling Task")

ancestry = %w[
  397227ff 4951413 119cf10 21bc4bf 8d14a67 b0052e9 47de307 b411b30
  b0a219a 551433d ab34e03 4493add
]
ancestry.each_cons(2) do |older, newer|
  _stdout, _stderr, status = Open3.capture3(
    "git", "-c", "core.fsmonitor=false", "merge-base", "--is-ancestor",
    older, newer, chdir: ROOT,
  )
  assert(status.success?, "#{older} precedes #{newer}")
end

expected = {
  "21bc4bf" => [
    ".github/workflows/release.yml",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/task_01.md",
  ],
  "8d14a67" => [
    ".github/workflows/release.yml",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/fixtures/eligible.json",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/fixtures/malformed.json",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/fixtures/package-cooldown.json",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/fixtures/version-unpublished-single.json",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/fixtures/version-used.json",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/task_02.md",
  ],
  "b0052e9" => [
    ".github/workflows/release.yml",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/task_03.md",
  ],
  "47de307" => [
    ".github/workflows/release.yml",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/task_04.md",
  ],
  "b411b30" => [
    "CONTEXT.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/task_05.md",
    "docs/user-guide/release-runbook.md",
  ],
  "551433d" => [
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/_prd.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/_tasks.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/_techspec.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/task_06.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/task_07.md",
  ],
  "ab34e03" => [
    ".github/workflows/release.yml",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/task_06.md",
  ],
  "4493add" => [
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/task_07.md",
    "docs/user-guide/release-runbook.md",
  ],
}
expected.each do |commit, paths|
  actual = git("diff-tree", "--no-commit-id", "--name-only", "-r", commit)
    .lines.map(&:strip).reject(&:empty?).sort
  assert(actual == paths.sort, "#{commit} changes only its bounded Task or remediation slice")
end

remediation_commits = git("rev-list", "--reverse", "b0a219a..4493add")
  .lines.map { |line| line.strip[0, 7] }
assert(remediation_commits == %w[551433d ab34e03 4493add], "remediation has separate Spec, tooling Task, and docs Task commits")

assert(
  prd.match?(/Ownership and trusted-publisher\s+configuration are explicitly\s+out of the preflight's reach/),
  "PRD no longer promises an impossible identity preflight",
)
assert(techspec.include?("identity cannot be preflighted at all"), "TechSpec matches the amended identity boundary")

platforms = JSON.parse(File.read(File.join(ROOT, "dist/npm/platforms.json")))
coordinates = platforms.map { |row| row.fetch("package") } + ["roundfix"]
runbook = File.read(File.join(ROOT, "docs/user-guide/release-runbook.md"))
closing = runbook[/^## Closing the fallback window\n(.*?)(?=^## |\z)/m, 1]
assert(!closing.nil?, "runbook has a closing procedure")
coordinates.each do |coordinate|
  assert(closing.include?("`#{coordinate}`"), "closing procedure names #{coordinate}")
end
assert(closing.match?(/fallback\s+record is empty/i), "registry shutdown follows empty fallback proof")
assert(closing.match?(/disallow token publication/i), "closing procedure disables token publication on npm")
assert(closing.match?(/reopen each package.*confirm/im), "closing procedure independently confirms all package settings")
assert(closing.include?("NPM_TRUSTED_PUBLISHING_FALLBACK"), "repository variable removal remains documented")
assert(closing.include?("NPM_TOKEN"), "repository secret removal remains documented")
assert(closing.include?("fallback branch"), "workflow fallback removal remains documented")

puts "SUMMARY: remediation and protected-tooling scope probe completed with all assertions passing"
