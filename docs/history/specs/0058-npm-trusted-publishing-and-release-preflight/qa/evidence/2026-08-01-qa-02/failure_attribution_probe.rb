#!/usr/bin/env ruby

require "open3"
require "tmpdir"
require "yaml"

root = File.expand_path("../../../../../..", __dir__)
workflow = YAML.load_file(File.join(root, ".github/workflows/release.yml"))
steps = workflow.fetch("jobs").fetch("release").fetch("steps")
publish = steps.find { |step| step["name"] == "Publish to npm" }.fetch("run")
publish = publish.gsub('${{ secrets.NPM_TOKEN }}', 'qa-secret-sentinel')

Dir.mktmpdir("roundfix-attribution-") do |dir|
  bin_dir = File.join(dir, "bin")
  Dir.mkdir(bin_dir)
  npm = File.join(bin_dir, "npm")
  File.write(
    npm,
    <<~'BASH',
      #!/usr/bin/env bash
      printf 'npm ERR! network timeout while contacting registry\n' >&2
      exit 1
    BASH
  )
  File.chmod(0o755, npm)
  File.write(File.join(dir, "roundfix-release-set"), "@roundfix/cli-darwin-arm64\nroundfix\n")
  summary = File.join(dir, "summary.md")

  stdout, stderr, status = Open3.capture3(
    {
      "FALLBACK_WINDOW" => "1",
      "GITHUB_STEP_SUMMARY" => summary,
      "NODE_AUTH_TOKEN" => nil,
      "PATH" => "#{bin_dir}:#{ENV.fetch("PATH")}",
      "RUNNER_TEMP" => dir,
    },
    "bash", "-c", publish, chdir: root,
  )

  combined = stdout + stderr
  puts combined
  puts "PUBLISH_EXIT=#{status.exitstatus}"

  misclassified = combined.include?("npm ERR! network timeout") &&
    combined.include?("::warning::identity: @roundfix/cli-darwin-arm64") &&
    combined.include?("retrying with the bounded token fallback")

  if misclassified
    warn "FAIL: a registry transport failure is classified as identity and retried with the token"
    exit 1
  end

  puts "PASS: non-identity publish failure was not classified as identity"
end
