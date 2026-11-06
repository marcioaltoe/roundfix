// Suite: expected workflow triggers.
// Invariant: only workflows guaranteed for every pull request into the base are expected.
// Boundary IN: workflow YAML bytes and the literal base branch.
// Boundary OUT: GitHub trigger evaluation, repository files and commands.
package delivery

import (
	"reflect"
	"testing"
)

func TestExpectedWorkflowsFollowTheTriggerRule(t *testing.T) {
	t.Parallel()
	files := map[string][]byte{
		"a.yml": []byte("on: {pull_request: {branches: [main], types: [opened, reopened, edited, synchronize]}}"),
		"b.yml": []byte("on: {pull_request: {branches: [main], paths: ['x/**']}}"),
		"c.yml": []byte("on: [push, pull_request]"),
		"d.yml": []byte("on: {push: {branches: [main]}}"),
		"e.yml": []byte("\"on\":\n  pull_request:\n"),
		"f.yml": []byte("on: {pull_request: {types: [opened, reopened, synchronize], branches: [main]}, push: {branches: [main]}}"),
		"g.yml": []byte("on: pull_request"),
		"h.yml": []byte("on:\n  push:\n"),
		"i.yml": []byte("on: {pull_request: {branches: ['releases/**']}}"),
		"j.yml": []byte("on: {pull_request: {types: [opened]}}"),
		"k.yml": []byte("on: {pull_request_target: {branches: main}}"),
		"l.yml": []byte("on: {pull_request: {branches-ignore: [docs]}}"),
		"m.yml": []byte("on: {pull_request: {branches: [main, '!main-old']}}"),
		"n.yml": []byte("on: workflow_dispatch"),
		"o.yml": []byte("on: ["),
		"p.yml": []byte("on: {pull_request: {paths-ignore: ['docs/**']}}"),
	}
	want := []string{"a.yml", "c.yml", "e.yml", "f.yml", "g.yml", "h.yml", "k.yml"}
	if got := ExpectedWorkflows(files, "main"); !reflect.DeepEqual(got, want) {
		t.Fatalf("ExpectedWorkflows = %v, want %v (b.yml has paths: ['x/**'])", got, want)
	}
}

func TestExpectedWorkflowsSkipUnsupportedTriggers(t *testing.T) {
	t.Parallel()
	files := map[string][]byte{
		"empty.yml":            []byte("on: {pull_request: {}}"),
		"target.yml":           []byte("on: pull_request_target"),
		"push.yml":             []byte("on: push"),
		"missing.yml":          []byte("name: CI"),
		"root.yml":             []byte("- pull_request"),
		"tags.yml":             []byte("on: {pull_request: {tags: [v1]}}"),
		"false.yml":            []byte("on: {pull_request: false}"),
		"invalid-branches.yml": []byte("on: {pull_request: {branches: {main: true}}}"),
		"invalid-entry.yml":    []byte("on: {pull_request: {branches: [main, {}]}}"),
		"invalid-types.yml":    []byte("on: {pull_request: {types: synchronize}}"),
		"null.yml":             []byte("on:"),
		"empty-file.yml":       nil,
	}
	want := []string{"empty.yml", "push.yml", "target.yml"}
	if got := ExpectedWorkflows(files, "main"); !reflect.DeepEqual(got, want) {
		t.Fatalf("ExpectedWorkflows = %v, want %v", got, want)
	}
}
