#!/usr/bin/env ruby

require "yaml"

prior = File.expand_path("../2026-07-31-qa-01/workflow_probe.rb", __dir__)
source = File.read(prior)
source = source.sub(
  "printf 'simulated OIDC failure\\n' >&2",
  "printf 'npm ERR! code ENEEDAUTH\\n' >&2",
)
source = source.sub(
  '        "NODE_AUTH_TOKEN" => nil,',
  "        \"NODE_AUTH_TOKEN\" => nil,\n        \"NPM_TOKEN\" => \"qa-secret-sentinel\",",
)
raise "prior probe boundaries not found" unless source.include?("npm ERR! code ENEEDAUTH") && source.include?('"NPM_TOKEN" => "qa-secret-sentinel"')

eval(source, TOPLEVEL_BINDING, prior)

workflow_path = File.expand_path("../../../../../../.github/workflows/release.yml", __dir__)
workflow = YAML.load_file(workflow_path)
raise "workflow-level permissions remain" if workflow.key?("permissions")

release = workflow.fetch("jobs").fetch("release")
expected_permissions = {"id-token" => "write", "contents" => "write"}
raise "release-job permission scope mismatch" unless release.fetch("permissions") == expected_permissions

publish = release.fetch("steps").find { |step| step["name"] == "Publish to npm" }
raise "step environment does not map the npm secret" unless publish.fetch("env").fetch("NPM_TOKEN") == "${{ secrets.NPM_TOKEN }}"
raise "generated shell still interpolates the secret expression" if publish.fetch("run").include?("${{ secrets.NPM_TOKEN }}")
raise "fallback does not consume the mapped shell variable" unless publish.fetch("run").include?('NODE_AUTH_TOKEN="$NPM_TOKEN"')

puts "PASS: review repair scopes permissions to the release job"
puts "PASS: review repair maps the secret through step env and keeps the generated shell expression-free"
