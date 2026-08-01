#!/usr/bin/env ruby

prior = File.expand_path(
  "../2026-07-31-qa-01/workflow_probe.rb",
  __dir__,
)
source = File.read(prior)
source = source.sub(
  "printf 'simulated OIDC failure\\n' >&2",
  "printf 'npm ERR! code ENEEDAUTH\\n' >&2",
)
source = source.sub(
  '        "NODE_AUTH_TOKEN" => nil,',
  "        \"NODE_AUTH_TOKEN\" => nil,\n" \
    "        \"NPM_FALLBACK_TOKEN\" => \"qa-secret-sentinel\",",
)

unless source.include?("npm ERR! code ENEEDAUTH") &&
    source.include?('"NPM_FALLBACK_TOKEN" => "qa-secret-sentinel"')
  raise "prior probe boundaries not found"
end

eval(source, TOPLEVEL_BINDING, prior)

base_text, base_stderr, base_status = Open3.capture3(
  "git", "-c", "core.fsmonitor=false", "show",
  "21bc4bf^:.github/workflows/release.yml", chdir: ROOT,
)
raise "pre-Spec workflow read failed: #{base_stderr}" unless base_status.success?

base_workflow = YAML.load(base_text)
base_steps = base_workflow.fetch("jobs").fetch("release").fetch("steps")
base_named = base_steps.each_with_object({}) do |step, result|
  result[step["name"]] = step if step["name"]
end
current_workflow = YAML.load_file(WORKFLOW_PATH)
current_steps = current_workflow.fetch("jobs").fetch("release").fetch("steps")
current_named = current_steps.each_with_object({}) do |step, result|
  result[step["name"]] = step if step["name"]
end

["Verify gate", "Cross-compile and stage", "GitHub Release"].each do |name|
  assert(
    current_named.fetch(name).fetch("run") == base_named.fetch(name).fetch("run"),
    "#{name} script remains byte-identical to the pre-Spec workflow",
  )
end
assert(
  current_steps.index(current_named.fetch("GitHub Release")) >
    current_steps.index(current_named.fetch("Publish to npm")),
  "GitHub Release remains after npm publication",
)
