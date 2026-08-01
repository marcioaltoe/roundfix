#!/usr/bin/env ruby

require "yaml"

root = File.expand_path("../../../../../..", __dir__)
workflow = YAML.load_file(File.join(root, ".github/workflows/release.yml"))
publish = workflow.fetch("jobs").fetch("release").fetch("steps")
  .find { |step| step["name"] == "Publish to npm" }
env = publish.fetch("env")
script = publish.fetch("run")

puts "SECRET_EXPRESSION_SCOPE=#{env["NPM_TOKEN"] == "${{ secrets.NPM_TOKEN }}" ? "publish-step" : "other"}"
puts "AUTH_TOKEN_ASSIGNMENTS=#{script.scan("NODE_AUTH_TOKEN").length}"
puts "SCRIPT_SECRET_EXPRESSIONS=#{script.scan("${{ secrets.NPM_TOKEN }}").length}"
puts "FALLBACK_SHELL_REFERENCES=#{script.scan('$NPM_TOKEN').length}"

if env["NPM_TOKEN"] == "${{ secrets.NPM_TOKEN }}"
  warn "FAIL: Task 04 requires the token secret expression to exist only inside the fallback branch, but the current workflow injects it into the entire Publish to npm step environment"
  exit 1
end

puts "PASS: retained token is unavailable outside the fallback branch"
