#!/usr/bin/env ruby

require "json"
require "open3"
require "yaml"

ROOT = File.expand_path("../../../../../..", __dir__)
SPEC = File.join(ROOT, "docs/specs/0058-npm-trusted-publishing-and-release-preflight")

def assert(condition, message)
  raise message unless condition

  puts "PASS: #{message}"
end

def git(*args)
  stdout, stderr, status = Open3.capture3("git", "-c", "core.fsmonitor=false", *args, chdir: ROOT)
  raise "git #{args.join(" ")} failed: #{stderr}" unless status.success?

  stdout
end

tasks = (1..5).map { |number| File.read(File.join(SPEC, format("task_%02d.md", number))) }
assert(tasks.all? { |task| task.match?(/^status: completed$/) }, "all five Task files are completed")
assert(tasks.all? { |task| task.include?("## Result") && task.include?("Acceptance evidence") }, "all Task Results carry acceptance evidence")

prd = File.read(File.join(SPEC, "_prd.md"))
techspec = File.read(File.join(SPEC, "_techspec.md"))
constraint_names = ["Identifier strategy", "Authentication and HTTP", "Active ADR obligations", "Tooling authority"]
[prd, techspec].each_with_index do |artifact, index|
  label = index.zero? ? "PRD" : "TechSpec"
  assert(constraint_names.all? { |name| artifact.include?(name) }, "#{label} accounts for every Project Constraint category")
  assert(artifact.include?("docs/agents/domain.md") && artifact.include?("docs/agents/agent-instructions.md"), "#{label} cites operative agent guidance")
end
assert(prd.include?('authorizes changes to exactly `.github/workflows/release.yml`'), "PRD expressly authorizes the exact protected workflow")
assert(techspec.include?('bounded files: `.github/workflows/release.yml`'), "TechSpec repeats the exact protected-tooling boundary")

authorization_prd = git("show", "397227ff:docs/specs/0058-npm-trusted-publishing-and-release-preflight/_prd.md")
assert(authorization_prd.include?('authorizes changes to exactly `.github/workflows/release.yml`'), "authorization exists in commit 397227ff")

ancestry = %w[397227ff 4951413 119cf10 21bc4bf 8d14a67 b0052e9 47de307 b411b30]
ancestry.each_cons(2) do |older, newer|
  _stdout, _stderr, status = Open3.capture3("git", "-c", "core.fsmonitor=false", "merge-base", "--is-ancestor", older, newer, chdir: ROOT)
  assert(status.success?, "#{older} precedes #{newer}")
end

expected_task_paths = {
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
}
expected_task_paths.each do |commit, expected|
  actual = git("diff-tree", "--no-commit-id", "--name-only", "-r", commit).lines.map(&:strip).reject(&:empty?).sort
  assert(actual == expected.sort, "#{commit} changes only its authorized Task slice")
end

implementation_commits = git("rev-list", "--reverse", "119cf10..b411b30").lines.map(&:strip)
assert(
  implementation_commits.map { |commit| commit[0, 7] } == %w[21bc4bf 8d14a67 b0052e9 47de307 b411b30],
  "no prerequisite or consequent fix is folded around the five Task commits",
)

changed_paths = git("diff", "--name-only", "119cf10..b411b30").lines.map(&:strip).reject(&:empty?)
assert(changed_paths.none? { |path| path.start_with?("dist/npm/", "cmd/", "internal/") }, "package manifests, command code, and Upgrade Command code are unchanged")
assert(
  changed_paths.all? do |path|
    path == ".github/workflows/release.yml" ||
      path == "CONTEXT.md" ||
      path == "docs/user-guide/release-runbook.md" ||
      path.start_with?("docs/specs/0058-npm-trusted-publishing-and-release-preflight/")
  end,
  "assembled implementation has no out-of-scope path",
)

final_workflow = YAML.load_file(File.join(ROOT, ".github/workflows/release.yml"))
baseline_workflow = YAML.load(git("show", "119cf10:.github/workflows/release.yml"))
final_steps = final_workflow.fetch("jobs").fetch("release").fetch("steps")
baseline_steps = baseline_workflow.fetch("jobs").fetch("release").fetch("steps")
find_step = ->(steps, name) { steps.find { |step| step["name"] == name } }
assert(find_step.call(final_steps, "Verify gate") == find_step.call(baseline_steps, "Verify gate"), "Verification step is byte-equivalent to the pre-Task workflow")
assert(find_step.call(final_steps, "Cross-compile and stage").fetch("run") == find_step.call(baseline_steps, "Cross-compile and stage").fetch("run"), "cross-compilation script is byte-equivalent to the pre-Task workflow")
assert(find_step.call(final_steps, "GitHub Release").fetch("run") == find_step.call(baseline_steps, "GitHub Release").fetch("run"), "GitHub Release and changelog script is byte-equivalent to the pre-Task workflow")

platforms = JSON.parse(File.read(File.join(ROOT, "dist/npm/platforms.json")))
coordinates = platforms.map { |row| row.fetch("package") } + ["roundfix"]
assets = platforms.map { |row| row.fetch("asset") }
assert(coordinates.length == 6 && coordinates.uniq.length == 6, "release package names remain one launcher plus five unique platforms")
assert(assets.length == 5 && assets.uniq.length == 5, "platform asset names remain five unique Upgrade Command inputs")

runbook = File.read(File.join(ROOT, "docs/user-guide/release-runbook.md"))
coordinates.each do |coordinate|
  assert(runbook.include?("`#{coordinate}`"), "runbook names trusted publisher coordinate #{coordinate}")
end
assert(runbook.scan("| `marcioaltoe` | `roundfix` | `release.yml` |").length == 6, "runbook binds all six coordinates to the same owner, repository, and workflow")
assert(runbook.match?(/only when `npm publish`\s+attempts the OIDC exchange/), "runbook states trusted-publisher validation is publish-time only")
assert(runbook.include?("manual `workflow_dispatch` trigger") && runbook.include?("This rehearsal is publish-free"), "runbook documents the publish-free rehearsal")
assert(runbook.match?(/one complete tagged\s+release whose fallback\s+record is empty/), "runbook defines the fallback exit evidence")
assert(runbook.include?("NPM_TRUSTED_PUBLISHING_FALLBACK") && runbook.match?(/Remove the\s+`NPM_TOKEN` repository secret/), "runbook defines the fallback switch and repository-secret removal")
%w[registry: undetermined: identity: runtime:].each do |prefix|
  assert(runbook.include?("`#{prefix}`"), "runbook documents #{prefix} recovery vocabulary")
end

context = File.read(File.join(ROOT, "CONTEXT.md"))
assert(context.include?("**Release Set**:") && context.match?(/launcher package and five platform package coordinates/), "glossary defines Release Set")
assert(context.include?("**Publication Preflight**:") && context.include?("read-only eligibility check"), "glossary defines Publication Preflight")

preflight_script = find_step.call(final_steps, "Publication preflight").fetch("run")
assert(preflight_script.scan(/\bcurl\b/).length == 1 && !preflight_script.match?(/\bsleep\b/), "preflight has no automatic cooldown retry")
assert(changed_paths.none? { |path| path.include?("release_plan") }, "Release Plan Command is unchanged")

unless runbook.match?(/disallow.{0,40}token|token.{0,40}disallow/i)
  warn "FAIL: runbook never instructs the maintainer to disallow token publication for the owned packages"
  exit 1
end

puts "PASS: runbook covers the PRD token-publication disablement step"
puts "SUMMARY: docs and scope probe completed with all assertions passing"
