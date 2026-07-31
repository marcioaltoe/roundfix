#!/usr/bin/env ruby

require "json"
require "open3"
require "tmpdir"
require "yaml"

ROOT = File.expand_path("../../../../../..", __dir__)
WORKFLOW_PATH = File.join(ROOT, ".github/workflows/release.yml")
FIXTURE_DIR = File.join(
  ROOT,
  "docs/specs/0058-npm-trusted-publishing-and-release-preflight/fixtures",
)

def assert(condition, message)
  raise message unless condition

  puts "PASS: #{message}"
end

def run_bash(script, env = {})
  Dir.mktmpdir("roundfix-qa-") do |dir|
    script_path = File.join(dir, "probe.sh")
    File.write(script_path, script)
    File.chmod(0o755, script_path)
    stdout, stderr, status = Open3.capture3(env, "bash", script_path, chdir: ROOT)
    return [stdout, stderr, status.exitstatus]
  end
end

workflow_text = File.read(WORKFLOW_PATH)
workflow = YAML.load_file(WORKFLOW_PATH)
steps = workflow.fetch("jobs").fetch("release").fetch("steps")
named_steps = steps.each_with_object({}) do |step, result|
  result[step["name"]] = step if step["name"]
end

expected_names = [
  "Guard Trusted Publishing runtime",
  "Validate tag",
  "Verify gate",
  "Publication preflight",
  "Cross-compile and stage",
  "Publish to npm",
  "GitHub Release",
]
assert(expected_names.all? { |name| named_steps.key?(name) }, "all release stages exist")
assert(
  expected_names.map { |name| steps.index(named_steps.fetch(name)) } ==
    expected_names.map { |name| steps.index(named_steps.fetch(name)) }.sort,
  "Verification, preflight, build, publish, and GitHub Release retain required order",
)
assert(workflow_text.include?('node-version: "24"'), "workflow requests Node 24")
assert(workflow_text.include?('registry-url: "https://registry.npmjs.org"'), "npm registry URL is unchanged")
assert(workflow_text.include?("id-token: write"), "workflow grants OIDC id-token permission")
assert(workflow_text.include?("contents: write"), "workflow retains GitHub Release permission")
assert(workflow_text.match?(/push:\n\s+tags:\n\s+- "v\*"/), "v* tag trigger remains declared")
assert(workflow_text.include?("workflow_dispatch:"), "manual preflight trigger is declared")
assert(workflow_text.match?(/version:\n\s+description:.*\n\s+required: true/), "manual trigger requires a version")

mutating_steps = ["Cross-compile and stage", "Publish to npm", "GitHub Release"]
mutating_steps.each do |name|
  assert(
    named_steps.fetch(name).fetch("if").include?("github.event_name == 'push'"),
    "#{name} excludes workflow_dispatch",
  )
end

runtime_script = named_steps.fetch("Guard Trusted Publishing runtime").fetch("run")
Dir.mktmpdir("roundfix-runtime-") do |dir|
  npm_stub = File.join(dir, "npm")
  File.write(npm_stub, "#!/usr/bin/env bash\nprintf '%s\\n' \"$QA_NPM_VERSION\"\n")
  File.chmod(0o755, npm_stub)
  [
    ["11.5.0", 1, "runtime:"],
    ["11.5.1", 0, "satisfies"],
    ["12.0.0", 0, "satisfies"],
  ].each do |version, expected_exit, marker|
    stdout, stderr, status = Open3.capture3(
      {"PATH" => "#{dir}:#{ENV.fetch("PATH")}", "QA_NPM_VERSION" => version},
      "bash",
      "-c",
      runtime_script,
      chdir: ROOT,
    )
    combined = stdout + stderr
    assert(status.exitstatus == expected_exit && combined.include?(marker), "runtime guard classifies npm #{version}")
  end
end

validate_script = named_steps.fetch("Validate tag").fetch("run")
owned_version = JSON.parse(File.read(File.join(ROOT, "dist/npm/roundfix/package.json"))).fetch("version")
[
  ["push", "v#{owned_version}", 0, "releasing #{owned_version}"],
  ["workflow_dispatch", "v#{owned_version}", 0, "preflighting #{owned_version}"],
  ["workflow_dispatch", "invalid", 1, "input invalid is not a semver version"],
  ["workflow_dispatch", "v99.99.99", 1, "does not match the checked-in Roundfix version"],
].each do |event, target, expected_exit, marker|
  stdout, stderr, status = run_bash(
    validate_script,
    {"GITHUB_EVENT_NAME" => event, "TARGET_VERSION" => target},
  )
  assert(status == expected_exit && (stdout + stderr).include?(marker), "validation handles #{event} target #{target}")
end

preflight_script = named_steps.fetch("Publication preflight").fetch("run")
classifier = preflight_script[/^\s*CLASSIFY_JQ='(.*)'$/, 1]
assert(!classifier.nil?, "workflow exposes one single-line CLASSIFY_JQ program")

fixture_matrix = {
  "version-used.json" => ["used", 0],
  "version-unpublished-single.json" => ["used", 0],
  "package-cooldown.json" => ["cooldown", 0],
  "eligible.json" => ["eligible", 0],
  "malformed.json" => [nil, 1],
}
fixture_matrix.each do |fixture, (expected, expected_failure)|
  stdout, _stderr, status = Open3.capture3(
    "jq", "-r", "--arg", "tag", "0.0.2", classifier, File.join(FIXTURE_DIR, fixture),
  )
  if expected_failure == 1
    assert(status.exitstatus != 0, "classifier rejects #{fixture}")
  else
    assert(status.exitstatus == 0 && stdout.strip == expected, "classifier resolves #{fixture} as #{expected}")
  end
end

platforms = JSON.parse(File.read(File.join(ROOT, "dist/npm/platforms.json")))
expected_release_set = platforms.map { |row| row.fetch("package") } + ["roundfix"]
assert(expected_release_set.length == 6 && expected_release_set.uniq.length == 6, "manifest defines five unique platform coordinates plus launcher")

curl_stub = <<~'BASH'
  #!/usr/bin/env bash
  set -euo pipefail
  output=""
  while [ "$#" -gt 0 ]; do
    case "$1" in
      --output) output="$2"; shift 2 ;;
      *) shift ;;
    esac
  done
  count=0
  if [ -f "$QA_COUNT_FILE" ]; then count="$(cat "$QA_COUNT_FILE")"; fi
  count=$((count + 1))
  printf '%s\n' "$count" > "$QA_COUNT_FILE"
  case "$QA_SCENARIO" in
    eligible) fixture="eligible.json"; status=200 ;;
    used) if [ "$count" -eq 1 ]; then fixture="version-used.json"; else fixture="eligible.json"; fi; status=200 ;;
    cooldown) if [ "$count" -eq 1 ]; then fixture="package-cooldown.json"; else fixture="eligible.json"; fi; status=200 ;;
    malformed) if [ "$count" -eq 1 ]; then fixture="malformed.json"; else fixture="eligible.json"; fi; status=200 ;;
    absent) fixture="eligible.json"; if [ "$count" -eq 1 ]; then status=404; else status=200; fi ;;
    http) fixture="eligible.json"; if [ "$count" -eq 1 ]; then status=503; else status=200; fi ;;
    transport) if [ "$count" -eq 1 ]; then exit 7; else fixture="eligible.json"; status=200; fi ;;
    mixed)
      if [ "$count" -eq 1 ]; then fixture="version-used.json"
      elif [ "$count" -eq 2 ]; then fixture="package-cooldown.json"
      else fixture="eligible.json"
      fi
      status=200
      ;;
    *) exit 64 ;;
  esac
  cp "$QA_FIXTURE_DIR/$fixture" "$output"
  printf '%s' "$status"
BASH

def run_preflight(script, curl_stub, scenario, fixture_dir, target)
  Dir.mktmpdir("roundfix-preflight-") do |dir|
    bin_dir = File.join(dir, "bin")
    Dir.mkdir(bin_dir)
    curl_path = File.join(bin_dir, "curl")
    File.write(curl_path, curl_stub)
    File.chmod(0o755, curl_path)
    summary = File.join(dir, "summary.md")
    count = File.join(dir, "count")
    stdout, stderr, status = Open3.capture3(
      {
        "PATH" => "#{bin_dir}:#{ENV.fetch("PATH")}",
        "GITHUB_STEP_SUMMARY" => summary,
        "QA_COUNT_FILE" => count,
        "QA_FIXTURE_DIR" => fixture_dir,
        "QA_SCENARIO" => scenario,
        "RUNNER_TEMP" => dir,
        "TARGET_VERSION" => "v#{target}",
      },
      "bash", "-c", script, chdir: ROOT,
    )
    summary_text = File.exist?(summary) ? File.read(summary) : ""
    release_set = File.read(File.join(dir, "roundfix-release-set")).lines.map(&:strip)
    return [stdout, stderr, status.exitstatus, summary_text, release_set]
  end
end

preflight_expectations = {
  "eligible" => [0, "eligible"],
  "used" => [1, "::error::registry:"],
  "cooldown" => [1, "in cooldown until 2026-07-26T17:14:13Z"],
  "malformed" => [1, "::error::undetermined:"],
  "absent" => [1, " is absent"],
  "http" => [1, "registry returned HTTP 503"],
  "transport" => [1, "registry transport failed"],
}
preflight_expectations.each do |scenario, (expected_exit, marker)|
  stdout, stderr, status, summary, release_set = run_preflight(
    preflight_script, curl_stub, scenario, FIXTURE_DIR, "0.0.2",
  )
  combined = stdout + stderr
  assert(status == expected_exit && combined.include?(marker), "preflight handles #{scenario} registry response")
  assert(release_set == expected_release_set, "preflight derives exact Release Set for #{scenario}")
  assert(summary.scan(/^\| `.*` \| `0\.0\.2` \|/).length == 6, "preflight summarizes all six coordinates for #{scenario}")
end

stdout, stderr, status, _summary, _set = run_preflight(
  preflight_script, curl_stub, "mixed", FIXTURE_DIR, "0.0.2",
)
combined = stdout + stderr
assert(status == 1 && combined.scan("::error::registry:").length == 2, "preflight reports every blocked coordinate before stopping")
assert(combined.include?("already used") && combined.include?("in cooldown until 2026-07-26T17:14:13Z"), "multi-coordinate errors preserve cause and cooldown expiry")

publish_script = named_steps.fetch("Publish to npm").fetch("run")
assert(publish_script.scan("NODE_AUTH_TOKEN").length == 1, "token environment name occurs only in fallback retry")
assert(!publish_script.match?(/echo .*NPM_TOKEN|echo .*NODE_AUTH_TOKEN/), "publish script has no token-output command")
assert(named_steps.fetch("Publish to npm").fetch("env").fetch("FALLBACK_WINDOW").include?("vars."), "fallback switch comes from repository configuration")

safe_publish_script = publish_script.gsub('${{ secrets.NPM_TOKEN }}', 'qa-secret-sentinel')
npm_stub = <<~'BASH'
  #!/usr/bin/env bash
  set -euo pipefail
  token_state="oidc"
  if [ -n "${NODE_AUTH_TOKEN:-}" ]; then token_state="token"; fi
  printf '%s|%s\n' "$PWD" "$token_state" >> "$QA_NPM_LOG"
  if [ "$QA_PUBLISH_SCENARIO" = "first-oidc-fails" ] && [ "$token_state" = "oidc" ] && [ ! -f "$QA_FAILED_ONCE" ]; then
    : > "$QA_FAILED_ONCE"
    printf 'simulated OIDC failure\n' >&2
    exit 1
  fi
  exit 0
BASH

def run_publish(script, npm_stub, release_set, scenario, fallback)
  Dir.mktmpdir("roundfix-publish-") do |dir|
    bin_dir = File.join(dir, "bin")
    Dir.mkdir(bin_dir)
    npm_path = File.join(bin_dir, "npm")
    File.write(npm_path, npm_stub)
    File.chmod(0o755, npm_path)
    File.write(File.join(dir, "roundfix-release-set"), release_set.join("\n") + "\n")
    summary = File.join(dir, "summary.md")
    npm_log = File.join(dir, "npm.log")
    stdout, stderr, status = Open3.capture3(
      {
        "PATH" => "#{bin_dir}:#{ENV.fetch("PATH")}",
        "FALLBACK_WINDOW" => fallback,
        "GITHUB_STEP_SUMMARY" => summary,
        "NODE_AUTH_TOKEN" => nil,
        "QA_FAILED_ONCE" => File.join(dir, "failed-once"),
        "QA_NPM_LOG" => npm_log,
        "QA_PUBLISH_SCENARIO" => scenario,
        "RUNNER_TEMP" => dir,
      },
      "bash", "-c", script, chdir: ROOT,
    )
    summary_text = File.exist?(summary) ? File.read(summary) : ""
    npm_log_text = File.exist?(npm_log) ? File.read(npm_log) : ""
    return [stdout, stderr, status.exitstatus, summary_text, npm_log_text]
  end
end

stdout, stderr, status, summary, npm_log = run_publish(
  safe_publish_script, npm_stub, expected_release_set, "all-oidc", "1",
)
assert(status == 0 && summary.include?("No coordinates required"), "all-OIDC publish records an empty fallback summary")
assert(npm_log.lines.length == 6 && npm_log.lines.all? { |line| line.include?("|oidc") }, "all six coordinates attempt OIDC without token")

stdout, stderr, status, summary, npm_log = run_publish(
  safe_publish_script, npm_stub, expected_release_set, "first-oidc-fails", "1",
)
assert(status == 0 && summary.include?(expected_release_set.first), "enabled fallback records the failed coordinate and completes")
assert(npm_log.lines.length == 7 && npm_log.lines[0].include?("|oidc") && npm_log.lines[1].include?("|token"), "enabled fallback retries only after OIDC failure")
assert(!(stdout + stderr + summary).include?("qa-secret-sentinel"), "fallback never exposes the token sentinel")

stdout, stderr, status, summary, npm_log = run_publish(
  safe_publish_script, npm_stub, expected_release_set, "first-oidc-fails", "0",
)
assert(status == 1 && (stdout + stderr).include?("::error::identity: #{expected_release_set.first}"), "closed fallback names the failed coordinate and stops")
assert(summary.empty? && npm_log.lines.length == 1, "closed fallback stops before later coordinates and summary completion")

ordered_dirs = expected_release_set.map do |coordinate|
  if coordinate == "roundfix"
    File.join(ROOT, "dist/npm/roundfix")
  else
    File.join(ROOT, "dist/npm/packages/cli-#{coordinate.sub("@roundfix/cli-", "")}")
  end
end
actual_dirs = run_publish(safe_publish_script, npm_stub, expected_release_set, "all-oidc", "1")[4]
  .lines.map { |line| line.split("|", 2).first }
assert(actual_dirs == ordered_dirs, "publish order is five platform packages followed by launcher")

puts "SUMMARY: workflow probe completed with all assertions passing"
