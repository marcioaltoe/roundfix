#!/usr/bin/env ruby

require "json"
require "open3"
require "tmpdir"
require "yaml"

ROOT = File.expand_path("../../../../../..", __dir__)
WORKFLOW = File.join(ROOT, ".github/workflows/release.yml")
SENTINEL = "qa-secret-sentinel"

def assert(condition, message)
  raise message unless condition

  puts "PASS: #{message}"
end

workflow_text = File.read(WORKFLOW)
workflow = YAML.load_file(WORKFLOW)
publish = workflow.fetch("jobs").fetch("release").fetch("steps")
  .find { |step| step["name"] == "Publish to npm" }
env = publish.fetch("env")
script = publish.fetch("run")

assert(env.fetch("NPM_FALLBACK_TOKEN") == "${{ secrets.NPM_TOKEN }}", "secret is mapped through the step environment")
assert(!script.include?("${{ secrets."), "parsed run script contains no secret expression")
assert(workflow_text.scan("secrets.NPM_TOKEN").length == 1, "workflow has one retained-token secret reference")
assert(script.index("unset NPM_FALLBACK_TOKEN") < script.index("npm publish"), "exported token copy is removed before the first publish")
assert(script.include?('NODE_AUTH_TOKEN="$npm_fallback_token" npm publish'), "fallback command receives the non-exported shell value")
assert(!script.match?(/(?:echo|printf).*NPM_FALLBACK_TOKEN|(?:echo|printf).*npm_fallback_token/), "token variables are absent from output commands")

platforms = JSON.parse(File.read(File.join(ROOT, "dist/npm/platforms.json")))
release_set = platforms.map { |row| row.fetch("package") } + ["roundfix"]

npm_stub = <<~'BASH'
  #!/usr/bin/env bash
  set -euo pipefail
  node_state=absent
  fallback_env_state=absent
  shell_state=absent
  if [ -n "${NODE_AUTH_TOKEN:-}" ]; then
    node_state=present
    if [ "$NODE_AUTH_TOKEN" != "qa-secret-sentinel" ]; then exit 70; fi
  fi
  if [ -n "${NPM_FALLBACK_TOKEN:-}" ]; then fallback_env_state=present; fi
  if [ -n "${npm_fallback_token:-}" ]; then shell_state=present; fi
  printf '%s|node=%s|fallback_env=%s|shell=%s\n' "$PWD" "$node_state" "$fallback_env_state" "$shell_state" >> "$QA_NPM_LOG"
  if [ ! -f "$QA_FAILED_ONCE" ] && [ "$node_state" = absent ]; then
    case "$QA_SCENARIO" in
      auth)
        : > "$QA_FAILED_ONCE"
        printf 'npm ERR! code ENEEDAUTH\n' >&2
        exit 1
        ;;
      network)
        : > "$QA_FAILED_ONCE"
        printf 'npm ERR! network timeout while contacting registry\n' >&2
        exit 1
        ;;
    esac
  fi
BASH

def run_publish(script, npm_stub, release_set, scenario, fallback)
  Dir.mktmpdir("roundfix-secret-boundary-") do |dir|
    bin = File.join(dir, "bin")
    Dir.mkdir(bin)
    npm = File.join(bin, "npm")
    File.write(npm, npm_stub)
    File.chmod(0o755, npm)
    File.write(File.join(dir, "roundfix-release-set"), release_set.join("\n") + "\n")
    summary = File.join(dir, "summary.md")
    log = File.join(dir, "npm.log")
    stdout, stderr, status = Open3.capture3(
      {
        "PATH" => "#{bin}:#{ENV.fetch("PATH")}",
        "FALLBACK_WINDOW" => fallback,
        "GITHUB_STEP_SUMMARY" => summary,
        "NPM_FALLBACK_TOKEN" => SENTINEL,
        "NODE_AUTH_TOKEN" => nil,
        "QA_FAILED_ONCE" => File.join(dir, "failed-once"),
        "QA_NPM_LOG" => log,
        "QA_SCENARIO" => scenario,
        "RUNNER_TEMP" => dir,
      },
      "bash", "-c", script, chdir: ROOT,
    )
    summary_text = File.exist?(summary) ? File.read(summary) : ""
    log_text = File.exist?(log) ? File.read(log) : ""
    files = Dir.glob(File.join(dir, "**", "*"), File::FNM_DOTMATCH)
      .select { |path| File.file?(path) && path != npm }
    persisted = files.any? { |path| File.binread(path).include?(SENTINEL) }
    [stdout, stderr, status.exitstatus, summary_text, log_text, persisted]
  end
end

stdout, stderr, status, summary, log, persisted = run_publish(script, npm_stub, release_set, "success", "1")
assert(status == 0, "all-OIDC path completes")
assert(log.lines.length == 6, "all six coordinates attempt OIDC")
assert(log.lines.all? { |line| line.include?("node=absent|fallback_env=absent|shell=absent") }, "every OIDC npm process inherits no token variable")
assert(summary.include?("No coordinates required"), "all-OIDC path persists an empty fallback record")
assert(!(stdout + stderr + summary).include?(SENTINEL) && !persisted, "all-OIDC path exposes or persists no token value")

stdout, stderr, status, summary, log, persisted = run_publish(script, npm_stub, release_set, "auth", "1")
assert(status == 0, "open-window authentication failure completes through fallback")
assert(log.lines.length == 7 && log.lines[0].include?("node=absent") && log.lines[1].include?("node=present"), "only the fallback retry receives the token")
assert(log.lines.drop(2).all? { |line| line.include?("node=absent|fallback_env=absent|shell=absent") }, "later OIDC attempts remain token-free")
assert(summary.include?(release_set.first) && !summary.include?(SENTINEL), "fallback summary records only the coordinate")
assert(!(stdout + stderr).include?(SENTINEL) && !persisted, "fallback path exposes or persists no token value")

stdout, stderr, status, summary, log, persisted = run_publish(script, npm_stub, release_set, "auth", "0")
assert(status == 1 && (stdout + stderr).include?("::error::identity: #{release_set.first}"), "closed window reports identity and fails")
assert(log.lines.length == 1 && summary.empty?, "closed window performs no retry or completed summary")
assert(!(stdout + stderr).include?(SENTINEL) && !persisted, "closed path exposes or persists no token value")

stdout, stderr, status, summary, log, persisted = run_publish(script, npm_stub, release_set, "network", "1")
assert(status == 1 && (stdout + stderr).include?("::error::publish: #{release_set.first}"), "non-authentication failure uses publish prefix")
assert((stdout + stderr).include?("network timeout") && log.lines.length == 1 && summary.empty?, "non-authentication failure surfaces cause without retry")
assert(!(stdout + stderr).include?(SENTINEL) && !persisted, "non-authentication path exposes or persists no token value")

expected_dirs = platforms.map do |row|
  File.join(ROOT, "dist/npm/packages/cli-#{row.fetch("package").sub("@roundfix/cli-", "")}")
end + [File.join(ROOT, "dist/npm/roundfix")]
actual_dirs = run_publish(script, npm_stub, release_set, "success", "1")[4]
  .lines.map { |line| line.split("|", 2).first }
assert(actual_dirs == expected_dirs, "publication order remains five platform packages then launcher")

puts "SUMMARY: retained token is confined to the fallback command"
