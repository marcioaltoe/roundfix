#!/usr/bin/env ruby

require "open3"

root = File.expand_path("../../../../../..", __dir__)
spec = File.join(root, "docs/specs/0058-npm-trusted-publishing-and-release-preflight")

def git_at(root, *args)
  stdout, stderr, status = Open3.capture3(
    "git", "-c", "core.fsmonitor=false", *args, chdir: root,
  )
  raise "git #{args.join(" ")} failed: #{stderr}" unless status.success?

  stdout
end

def assert_current(condition, message)
  raise message unless condition

  puts "PASS: #{message}"
end

load File.expand_path("../2026-08-01-qa-03/current_scope_probe.rb", __dir__)

task08 = File.read(File.join(spec, "task_08.md"))
assert_current(task08.match?(/^status: completed$/), "task_08 is completed")
assert_current(task08.include?("## Result") && task08.include?("Acceptance evidence"), "task_08 Result carries acceptance evidence")

expected = {
  "a8276a4" => [
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/_tasks.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/task_08.md",
  ],
  "e45dd37" => [
    ".github/workflows/release.yml",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/task_08.md",
  ],
}
expected.each do |commit, paths|
  actual = git_at(root, "diff-tree", "--no-commit-id", "--name-only", "-r", commit)
    .lines.map(&:strip).reject(&:empty?).sort
  assert_current(actual == paths.sort, "#{commit} changes only its Task-graph or authorized tooling slice")
end

%w[397227ff a8276a4 e45dd37].each_cons(2) do |older, newer|
  _stdout, _stderr, status = Open3.capture3(
    "git", "-c", "core.fsmonitor=false", "merge-base", "--is-ancestor",
    older, newer, chdir: root,
  )
  assert_current(status.success?, "#{older} precedes #{newer}")
end

authorization = git_at(root, "show", "397227ff:docs/specs/0058-npm-trusted-publishing-and-release-preflight/_prd.md")
assert_current(authorization.include?('authorizes changes to exactly `.github/workflows/release.yml`'), "express authorization predates task_08")

worktree_paths = git_at(root, "status", "--porcelain").lines.map { |line| line[3..].strip }
assert_current(
  worktree_paths.all? { |path| path.start_with?("docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/") },
  "current delta contains only this QA report and evidence",
)

puts "SUMMARY: all eight Tasks and protected-tooling chronology pass"
