#!/usr/bin/env ruby

require "open3"

root = File.expand_path("../../../../../..", __dir__)

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

load File.expand_path("../2026-08-01-qa-02/remediation_scope_probe.rb", __dir__)

expected = {
  "04d5bbb" => ["docs/user-guide/release-runbook.md"],
  "83fbf6b" => [".github/workflows/release.yml"],
  "171f6a3" => [
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/reviews/round-001/issue_001.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/reviews/round-001/issue_002.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/reviews/round-001/round.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/reviews/round-002/issue_001.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/reviews/round-002/issue_002.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/reviews/round-002/issue_003.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/reviews/round-002/issue_004.md",
    "docs/specs/0058-npm-trusted-publishing-and-release-preflight/reviews/round-002/round.md",
  ],
}
expected.each do |commit, paths|
  actual = git_at(root, "diff-tree", "--no-commit-id", "--name-only", "-r", commit)
    .lines.map(&:strip).reject(&:empty?).sort
  assert_current(actual == paths.sort, "#{commit} changes only its consequent-fix or review-artifact scope")
end

%w[397227ff 4493add 04d5bbb 83fbf6b 171f6a3].each_cons(2) do |older, newer|
  _stdout, _stderr, status = Open3.capture3(
    "git", "-c", "core.fsmonitor=false", "merge-base", "--is-ancestor",
    older, newer, chdir: root,
  )
  assert_current(status.success?, "#{older} precedes #{newer}")
end

review_files = Dir.glob(File.join(root, "docs/specs/0058-npm-trusted-publishing-and-release-preflight/reviews/round-002/issue_*.md"))
statuses = review_files.map { |path| File.read(path)[/^status: (\S+)/, 1] }
assert_current(statuses.sort == %w[invalid invalid resolved resolved], "review artifacts settle every fetched round-002 issue")

worktree_paths = git_at(root, "status", "--porcelain").lines.map { |line| line[3..].strip }
assert_current(
  worktree_paths.all? { |path| path.start_with?("docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/") },
  "current delta contains only this QA report and evidence",
)

puts "SUMMARY: current protected-tooling and review-artifact chronology passes"
