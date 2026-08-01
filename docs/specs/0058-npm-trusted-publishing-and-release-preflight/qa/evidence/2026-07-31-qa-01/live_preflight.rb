#!/usr/bin/env ruby

require "json"
require "open3"
require "tmpdir"
require "yaml"

root = File.expand_path("../../../../../..", __dir__)
workflow = YAML.load_file(File.join(root, ".github/workflows/release.yml"))
steps = workflow.fetch("jobs").fetch("release").fetch("steps")
script = steps.find { |step| step["name"] == "Publication preflight" }.fetch("run")
version = JSON.parse(File.read(File.join(root, "dist/npm/roundfix/package.json"))).fetch("version")

Dir.mktmpdir("roundfix-live-preflight-") do |dir|
  summary = File.join(dir, "summary.md")
  stdout, stderr, status = Open3.capture3(
    {
      "GITHUB_STEP_SUMMARY" => summary,
      "RUNNER_TEMP" => dir,
      "TARGET_VERSION" => "v#{version}",
    },
    "bash", "-c", script, chdir: root,
  )

  puts stdout
  warn stderr unless stderr.empty?
  puts "--- job summary ---"
  puts File.read(summary) if File.exist?(summary)
  puts "PREFLIGHT_EXIT=#{status.exitstatus}"
  exit(status.exitstatus.between?(0, 1) ? 0 : status.exitstatus)
end
