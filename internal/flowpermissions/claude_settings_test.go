package flowpermissions

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestEditClaudeSettings(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, before, after string }{
		{"absent_file_is_created", "", "{\n  \"permissions\": {\n    \"allow\": [\n      \"new\"\n    ],\n    \"ask\": [\n      \"guard\"\n    ]\n  }\n}\n"},
		{"empty_object_gains_permissions", "{}", "{\n  \"permissions\": {\n    \"allow\": [\n      \"new\"\n    ],\n    \"ask\": [\n      \"guard\"\n    ]\n  }\n}"},
		{"missing_ask_list_is_created", "{\n  \"permissions\": {\n    \"allow\": [\"user\"]\n  }\n}", "{\n  \"permissions\": {\n    \"allow\": [\"user\", \"new\"],\n    \"ask\": [\n      \"guard\"\n    ]\n  }\n}"},
		{"user_rules_kept", `{"permissions":{"allow":["user", "user2"],"ask":["user-ask"]}}`, `{"permissions":{"allow":["user", "user2", "new"],"ask":["user-ask", "guard"]}}`},
		{"other_keys_kept", ` {"env": { "X":"a\\b" }, "permissions":{"allow":["user"],"ask":["user-ask"],"deny":["no"]},"x":1e2} `, ` {"env": { "X":"a\\b" }, "permissions":{"allow":["user", "new"],"ask":["user-ask", "guard"],"deny":["no"]},"x":1e2} `},
		{"tab_indentation_kept", "{\n\t\"permissions\": {\n\t\t\"allow\": [\n\t\t\t\"user\"\n\t\t],\n\t\t\"ask\": [\"user-ask\"]\n\t}\n}", "{\n\t\"permissions\": {\n\t\t\"allow\": [\n\t\t\t\"user\",\n\t\t\t\"new\"\n\t\t],\n\t\t\"ask\": [\"user-ask\", \"guard\"]\n\t}\n}"},
		{"one_line_lists", `{"permissions":{"allow":["a","b"],"ask":["c","d"]}}`, `{"permissions":{"allow":["a","b","new"],"ask":["c","d","guard"]}}`},
		{"no_final_newline_kept", `{"permissions":{"allow":["a"],"ask":["b"]}}`, `{"permissions":{"allow":["a", "new"],"ask":["b", "guard"]}}`},
		{"windows_line_ends_kept", "{\r\n  \"permissions\": {\r\n    \"allow\": [\r\n      \"a\"\r\n    ],\r\n    \"ask\": [\r\n      \"b\"\r\n    ]\r\n  }\r\n}\r\n", "{\r\n  \"permissions\": {\r\n    \"allow\": [\r\n      \"a\",\r\n      \"new\"\r\n    ],\r\n    \"ask\": [\r\n      \"b\",\r\n      \"guard\"\r\n    ]\r\n  }\r\n}\r\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var before []byte
			if tc.before != "" {
				before = []byte(tc.before)
			}
			got, err := EditClaudeSettings(before, []string{"new"}, []string{"guard"}, nil)
			if err != nil || string(got) != tc.after {
				t.Fatalf("edit: %v\ngot  %q\nwant %q", err, got, tc.after)
			}
			again, err := EditClaudeSettings(got, []string{"new", "new"}, []string{"guard"}, nil)
			if err != nil || !bytes.Equal(again, got) {
				t.Fatalf("second edit: %q, %v", again, err)
			}
			if strings.Contains(tc.before, `"ask"`) {
				restored, err := EditClaudeSettings(got, nil, nil, []string{"new", "guard"})
				if err != nil || !bytes.Equal(restored, before) {
					t.Fatalf("round trip: %q, %v; want %q", restored, err, before)
				}
			}
		})
	}
	for _, tc := range []struct {
		name, before, after string
		add, ask, remove    []string
	}{
		{"second_apply_changes_nothing", `{"permissions":{"allow":["new"],"ask":["guard"]}}`, `{"permissions":{"allow":["new"],"ask":["guard"]}}`, []string{"new"}, []string{"guard"}, nil},
		{"removal_restores_first_bytes", `{"permissions":{"allow":["user", "new"],"ask":["user-ask", "guard"]}}`, `{"permissions":{"allow":["user"],"ask":["user-ask"]}}`, nil, nil, []string{"new", "guard"}},
		{"removal_keeps_user_rules", `{"permissions":{"allow":["new", "u", "new", "v", "new"],"ask":["guard", "u"]},"other":true}`, `{"permissions":{"allow":["u", "v"],"ask":["u"]},"other":true}`, nil, nil, []string{"new", "guard"}},
		{"emptied_list_is_brackets", "{\"permissions\":{\"allow\":[\n  \"new\"\n],\"ask\":[ \"guard\" ]}}", `{"permissions":{"allow":[],"ask":[]}}`, nil, nil, []string{"new", "guard"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := EditClaudeSettings([]byte(tc.before), tc.add, tc.ask, tc.remove)
			if err != nil || string(got) != tc.after {
				t.Fatalf("got %q, %v; want %q", got, err, tc.after)
			}
		})
	}
	for _, tc := range []struct {
		name, data string
		want       error
	}{
		{"comment_is_not_strict_json", "{ // comment\n}", ErrNotStrictJSON},
		{"trailing_comma_is_not_strict_json", `{"permissions":{"allow":["x",]}}`, ErrNotStrictJSON},
		{"allow_of_another_type", `{"permissions":{"allow":null}}`, ErrUnexpectedShape},
		{"duplicate_key", `{"permissions":{},"permissions":{}}`, ErrUnexpectedShape},
		{"duplicate_permission_key", `{"permissions":{"allow":[],"allow":[]}}`, ErrUnexpectedShape},
		{"non_object", `[]`, ErrUnexpectedShape},
		{"permissions_null", `{"permissions":null}`, ErrUnexpectedShape},
		{"ask_of_another_type", `{"permissions":{"ask":true}}`, ErrUnexpectedShape},
		{"non_string_element", `{"permissions":{"allow":[null]}}`, ErrUnexpectedShape},
		{"trailing_document", `{} {}`, ErrNotStrictJSON},
		{"empty_existing_file", "", ErrNotStrictJSON},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := ReadClaudeSettings([]byte(tc.data)); !errors.Is(err, tc.want) {
				t.Fatalf("read: %v; want %v", err, tc.want)
			}
			got, err := EditClaudeSettings([]byte(tc.data), []string{"new"}, nil, nil)
			if !errors.Is(err, tc.want) || got != nil {
				t.Fatalf("edit returned %q, %v; want nil, %v", got, err, tc.want)
			}
		})
	}
}

func TestClaudeSettingsSpanEdges(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, before, after string
		allow, ask, remove  []string
	}{
		{"multiple_additions_copy_last_separator", "{\"permissions\":{\"allow\":[\"a\", \"b\",\n\t\"c\"],\"ask\":[]}}", "{\"permissions\":{\"allow\":[\"a\", \"b\",\n\t\"c\",\n\t\"d\",\n\t\"e\"],\"ask\":[]}}", []string{"d", "e"}, nil, nil},
		{"escaped_rule_already_present", `{"permissions":{"allow":["\u0061"],"ask":[]}}`, `{"permissions":{"allow":["\u0061"],"ask":[]}}`, []string{"a"}, nil, nil},
		{"missing_allow_keeps_ask_bytes", "{\n  \"permissions\": {\n    \"ask\": [ \"user\" ]\n  }\n}", "{\n  \"permissions\": {\n    \"ask\": [ \"user\" ],\n    \"allow\": [\n      \"new\"\n    ]\n  }\n}", []string{"new"}, nil, nil},
		{"empty_existing_lists_round_trip", `{"permissions":{"allow":[],"ask":[]}}`, `{"permissions":{"allow":["a", "b"],"ask":["guard"]}}`, []string{"a", "b"}, []string{"guard"}, nil},
		{"no_op_keeps_all_bytes", " {\n \"permissions\":{\"allow\": [ ]}} ", " {\n \"permissions\":{\"allow\": [ ]}} ", nil, nil, []string{"absent"}},
		{"removal_keeps_escaped_user_text", `{"permissions":{"allow":["gone", "\u0061", "gone"],"ask":[]}}`, `{"permissions":{"allow":["\u0061"],"ask":[]}}`, nil, nil, []string{"gone"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before := []byte(tc.before)
			original := bytes.Clone(before)
			got, err := EditClaudeSettings(before, tc.allow, tc.ask, tc.remove)
			if err != nil || string(got) != tc.after {
				t.Fatalf("%q, %v; want %q", got, err, tc.after)
			}
			if !bytes.Equal(before, original) {
				t.Fatal("editor changed input slice")
			}
			if tc.name == "empty_existing_lists_round_trip" {
				restored, err := EditClaudeSettings(got, nil, nil, []string{"a", "b", "guard"})
				if err != nil || !bytes.Equal(restored, before) {
					t.Fatalf("empty lists round trip: %q, %v", restored, err)
				}
			}
		})
	}
	t.Run("absent_removal_is_no_op", func(t *testing.T) {
		got, err := EditClaudeSettings(nil, nil, nil, []string{"gone"})
		if err != nil || got != nil {
			t.Fatalf("%q, %v", got, err)
		}
	})
}
