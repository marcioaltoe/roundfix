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
raise "prior probe boundary not found" if source == File.read(prior)

eval(source, TOPLEVEL_BINDING, prior)
