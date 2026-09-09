package spec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const legacyAuthorizationDirectory = "docs/workflow/authorizations"

// AuthorizationRole declares which bounded record format a caller is asking
// the reader to interpret.
type AuthorizationRole string

const (
	AuthorizationRoleSpec   AuthorizationRole = "spec"
	AuthorizationRoleLegacy AuthorizationRole = "legacy"
)

// AuthorizationOutcome separates an operative grant from a proved refusal
// and from a source the reader could not inspect.
type AuthorizationOutcome string

const (
	AuthorizationGranted    AuthorizationOutcome = "granted"
	AuthorizationRefused    AuthorizationOutcome = "refused"
	AuthorizationUnresolved AuthorizationOutcome = "unresolved"
)

// AuthorizationStatus is the approval state declared by a record.
type AuthorizationStatus string

const (
	AuthorizationStatusApproved  AuthorizationStatus = "approved"
	AuthorizationStatusProposed  AuthorizationStatus = "proposed"
	AuthorizationStatusWithdrawn AuthorizationStatus = "withdrawn"
)

// AuthorizationOperation is one action a grant can permit. The vocabulary is
// deliberately closed so adding a new delivery action cannot create authority.
type AuthorizationOperation string

const (
	AuthorizationOperationImplement   AuthorizationOperation = "implement"
	AuthorizationOperationCommit      AuthorizationOperation = "commit"
	AuthorizationOperationPush        AuthorizationOperation = "push"
	AuthorizationOperationPullRequest AuthorizationOperation = "pull_request"
	AuthorizationOperationMerge       AuthorizationOperation = "merge"
	AuthorizationOperationRelease     AuthorizationOperation = "release"
)

var authorizationOperations = [...]AuthorizationOperation{
	AuthorizationOperationImplement,
	AuthorizationOperationCommit,
	AuthorizationOperationPush,
	AuthorizationOperationPullRequest,
	AuthorizationOperationMerge,
	AuthorizationOperationRelease,
}

// AllAuthorizationOperations returns the complete operation vocabulary in its
// declared order.
func AllAuthorizationOperations() []AuthorizationOperation {
	return append([]AuthorizationOperation(nil), authorizationOperations[:]...)
}

// AuthorizationRegeneration is one command and any outputs a historical
// record explicitly enumerated. A nil Outputs slice leaves output ownership to
// the repository's regeneration declaration.
type AuthorizationRegeneration struct {
	Command string
	Outputs []string
}

// AuthorizationSource identifies the record bytes the reader inspected. A
// non-empty Revision is the resolved commit rather than the caller's ref.
type AuthorizationSource struct {
	Path     string
	Revision string
}

// AuthorizationRecord is the inspectable projection of one authorization
// record, including non-granting proposals and refused records.
type AuthorizationRecord struct {
	Role          AuthorizationRole
	Status        AuthorizationStatus
	GrantedAt     time.Time
	Action        string
	Consuming     []string
	Paths         []string
	Operations    []AuthorizationOperation
	Regenerations []AuthorizationRegeneration
	Source        AuthorizationSource
}

// NamesSpec reports whether the declared consuming field names specSlug. The
// numeric form is accepted only for prose-era legacy records.
func (record AuthorizationRecord) NamesSpec(specSlug string) bool {
	for _, consuming := range record.Consuming {
		if consuming == specSlug {
			return true
		}
		if record.Role != AuthorizationRoleLegacy {
			continue
		}
		number, _, hasSuffix := strings.Cut(specSlug, "-")
		if hasSuffix && consuming == number {
			return true
		}
	}
	return false
}

// AuthorizationReasonCode gives callers a stable reason category while Field,
// Value, and Detail retain the exact withholding evidence.
type AuthorizationReasonCode string

const (
	AuthorizationReasonRole                AuthorizationReasonCode = "role"
	AuthorizationReasonMalformedRecord     AuthorizationReasonCode = "malformed_record"
	AuthorizationReasonStatus              AuthorizationReasonCode = "status"
	AuthorizationReasonGranted             AuthorizationReasonCode = "granted"
	AuthorizationReasonAction              AuthorizationReasonCode = "action"
	AuthorizationReasonConsuming           AuthorizationReasonCode = "consuming"
	AuthorizationReasonPaths               AuthorizationReasonCode = "paths"
	AuthorizationReasonOperations          AuthorizationReasonCode = "operations"
	AuthorizationReasonRegeneration        AuthorizationReasonCode = "sanctioned_regeneration"
	AuthorizationReasonContradictory       AuthorizationReasonCode = "contradictory"
	AuthorizationReasonUnreadableRecord    AuthorizationReasonCode = "unreadable_record"
	AuthorizationReasonUnavailableRevision AuthorizationReasonCode = "unavailable_revision"
)

// AuthorizationReason describes why a record did not grant. Field and Value
// are machine-readable evidence; Detail is suitable for a diagnostic.
type AuthorizationReason struct {
	Code   AuthorizationReasonCode
	Field  string
	Value  string
	Detail string
}

// AuthorizationResolution is the complete answer for one record and asking
// Spec. Record remains populated for refusals so proposals stay inspectable.
type AuthorizationResolution struct {
	Outcome AuthorizationOutcome
	Record  AuthorizationRecord
	Reason  AuthorizationReason
}

// Permits reports whether an operative grant explicitly lists operation.
func (resolution AuthorizationResolution) Permits(operation AuthorizationOperation) bool {
	if resolution.Outcome != AuthorizationGranted || !allowedAuthorizationOperation(operation) {
		return false
	}
	for _, permitted := range resolution.Record.Operations {
		if permitted == operation {
			return true
		}
	}
	return false
}

// AuthorizationReadRequest names one working-tree or committed record. An
// empty Revision reads the working tree. An empty AskingSpec validates the
// record's intrinsic grant without narrowing it to one consumer.
type AuthorizationReadRequest struct {
	RepoRoot   string
	RecordPath string
	Revision   string
	Role       AuthorizationRole
	AskingSpec string
}

// ReadAuthorization reads and classifies one authorization record. Evidence
// failures are returned as unresolved results; malformed or non-operative
// records are refusals.
func ReadAuthorization(ctx context.Context, request AuthorizationReadRequest) AuthorizationResolution {
	root, err := filepath.Abs(request.RepoRoot)
	if err != nil {
		return unresolvedAuthorization(
			AuthorizationRecord{Role: request.Role},
			AuthorizationReasonUnreadableRecord,
			"record_path",
			request.RecordPath,
			fmt.Sprintf("resolve repository root: %v", err),
		)
	}
	recordPath, ok := cleanAuthorizationSourcePath(request.RecordPath)
	if !ok {
		return unresolvedAuthorization(
			AuthorizationRecord{Role: request.Role, Source: AuthorizationSource{Path: request.RecordPath}},
			AuthorizationReasonUnreadableRecord,
			"record_path",
			request.RecordPath,
			fmt.Sprintf("authorization record path %q is not repository-relative", request.RecordPath),
		)
	}
	record := AuthorizationRecord{
		Role:   request.Role,
		Source: AuthorizationSource{Path: recordPath},
	}
	if request.Role != AuthorizationRoleSpec && request.Role != AuthorizationRoleLegacy {
		return refusedAuthorization(
			record,
			AuthorizationReasonRole,
			"role",
			string(request.Role),
			fmt.Sprintf("authorization role %q is unsupported", request.Role),
		)
	}
	if request.Role == AuthorizationRoleLegacy && !pathWithinAuthorizationDirectory(recordPath, legacyAuthorizationDirectory) {
		return refusedAuthorization(
			record,
			AuthorizationReasonRole,
			"role",
			string(request.Role),
			fmt.Sprintf("legacy authorization record %q is outside %s", recordPath, legacyAuthorizationDirectory),
		)
	}

	var content []byte
	if strings.TrimSpace(request.Revision) == "" {
		content, err = readWorkingAuthorization(root, recordPath)
		if err != nil {
			return unresolvedAuthorization(
				record,
				AuthorizationReasonUnreadableRecord,
				"record_path",
				recordPath,
				fmt.Sprintf("read authorization record %q: %v", recordPath, err),
			)
		}
	} else {
		var revision string
		revision, err = resolveAuthorizationRevision(ctx, root, request.Revision)
		if err != nil {
			return unresolvedAuthorization(
				record,
				AuthorizationReasonUnavailableRevision,
				"revision",
				request.Revision,
				fmt.Sprintf("resolve authorization revision %q: %v", request.Revision, err),
			)
		}
		record.Source.Revision = revision
		content, err = readAuthorizationRevision(ctx, root, revision, recordPath)
		if err != nil {
			return unresolvedAuthorization(
				record,
				AuthorizationReasonUnreadableRecord,
				"record_path",
				recordPath,
				fmt.Sprintf("read authorization record %q at %s: %v", recordPath, revision, err),
			)
		}
	}

	parsed, reason := parseAuthorizationRecord(request.Role, record.Source, content)
	if reason.Code != "" {
		return AuthorizationResolution{Outcome: AuthorizationRefused, Record: parsed, Reason: reason}
	}
	if request.AskingSpec != "" && !parsed.NamesSpec(request.AskingSpec) {
		return refusedAuthorization(
			parsed,
			AuthorizationReasonConsuming,
			"consuming",
			request.AskingSpec,
			fmt.Sprintf("consuming does not name asking Spec %q", request.AskingSpec),
		)
	}
	if reason, err := validateAuthorizationPaths(ctx, root, parsed.Source.Revision, parsed.Paths); err != nil {
		return unresolvedAuthorization(
			parsed,
			AuthorizationReasonUnreadableRecord,
			"paths",
			"",
			fmt.Sprintf("inspect bounded authorization paths: %v", err),
		)
	} else if reason.Code != "" {
		return AuthorizationResolution{Outcome: AuthorizationRefused, Record: parsed, Reason: reason}
	}
	return AuthorizationResolution{Outcome: AuthorizationGranted, Record: parsed}
}

func parseAuthorizationRecord(
	role AuthorizationRole,
	source AuthorizationSource,
	content []byte,
) (AuthorizationRecord, AuthorizationReason) {
	record := AuthorizationRecord{Role: role, Source: source}
	normalized := bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
	frontmatter, body, err := splitFrontmatter(normalized)
	if err == nil {
		return parseAuthorizationFrontmatter(role, source, frontmatter, body)
	}
	if role != AuthorizationRoleLegacy {
		return record, authorizationReason(
			AuthorizationReasonMalformedRecord,
			"frontmatter",
			"",
			fmt.Sprintf("authorization record has invalid frontmatter: %v", err),
		)
	}
	return parseLegacyAuthorization(source, normalized)
}

type authorizationFrontmatterPresence struct {
	status     bool
	granted    bool
	grantedNil bool
}

func parseAuthorizationFrontmatter(
	role AuthorizationRole,
	source AuthorizationSource,
	frontmatter []byte,
	body []byte,
) (AuthorizationRecord, AuthorizationReason) {
	record := AuthorizationRecord{Role: role, Source: source}
	fields, err := authorizationFrontmatterMapping(frontmatter)
	if err != nil {
		return record, authorizationReason(
			AuthorizationReasonMalformedRecord,
			"frontmatter",
			"",
			fmt.Sprintf("decode authorization frontmatter: %v", err),
		)
	}
	presence := authorizationFrontmatterPresence{}

	if node, ok := fields["status"]; ok {
		presence.status = true
		value, err := authorizationScalar(node)
		if err != nil {
			return record, fieldShapeReason("status", err)
		}
		record.Status = AuthorizationStatus(strings.TrimSpace(value))
	} else if role == AuthorizationRoleLegacy {
		record.Status = AuthorizationStatusApproved
	}
	if node, ok := fields["granted"]; ok {
		presence.granted = true
		if node.Tag == "!!null" {
			presence.grantedNil = true
		} else {
			value, err := authorizationDateScalar(node)
			if err != nil {
				return record, fieldShapeReason("granted", err)
			}
			grantedAt, err := time.Parse("2006-01-02", strings.TrimSpace(value))
			if err != nil {
				return record, authorizationReason(
					AuthorizationReasonGranted,
					"granted",
					value,
					fmt.Sprintf("granted value %q is not a YYYY-MM-DD date", value),
				)
			}
			record.GrantedAt = grantedAt
		}
	}
	if node, ok := fields["action"]; ok {
		value, err := authorizationScalar(node)
		if err != nil {
			return record, fieldShapeReason("action", err)
		}
		record.Action = strings.TrimSpace(value)
	}
	if node, ok := fields["consuming"]; ok {
		values, err := authorizationStringValues(node)
		if err != nil {
			return record, fieldShapeReason("consuming", err)
		}
		if role == AuthorizationRoleLegacy && node.Kind == yaml.ScalarNode && len(values) == 1 {
			values = splitLegacyConsumers(values[0])
		}
		record.Consuming = trimNonEmptyAuthorizationValues(values)
	}
	if node, ok := fields["paths"]; ok {
		if node.Tag != "!!null" {
			values, err := authorizationStringSequence(node)
			if err != nil {
				return record, fieldShapeReason("paths", err)
			}
			record.Paths = append([]string(nil), values...)
		}
	}
	if node, ok := fields["operations"]; ok && node.Tag != "!!null" {
		values, err := authorizationStringSequence(node)
		if err != nil {
			return record, fieldShapeReason("operations", err)
		}
		seen := make(map[AuthorizationOperation]struct{}, len(values))
		for _, value := range values {
			operation := AuthorizationOperation(value)
			if strings.TrimSpace(value) != value || !allowedAuthorizationOperation(operation) {
				return record, authorizationReason(
					AuthorizationReasonOperations,
					"operations",
					value,
					fmt.Sprintf("operations contains unsupported token %q", value),
				)
			}
			if _, duplicate := seen[operation]; duplicate {
				return record, authorizationReason(
					AuthorizationReasonOperations,
					"operations",
					value,
					fmt.Sprintf("operations repeats token %q", value),
				)
			}
			seen[operation] = struct{}{}
			record.Operations = append(record.Operations, operation)
		}
	}
	regenerations, reason := parseAuthorizationRegenerations(body)
	record.Regenerations = regenerations
	if reason.Code != "" {
		return record, reason
	}

	return classifyAuthorizationRecord(record, presence)
}

func classifyAuthorizationRecord(
	record AuthorizationRecord,
	presence authorizationFrontmatterPresence,
) (AuthorizationRecord, AuthorizationReason) {
	if !presence.status && record.Role == AuthorizationRoleSpec {
		return record, authorizationReason(
			AuthorizationReasonStatus,
			"status",
			"",
			"status is required for a Spec authorization record",
		)
	}
	switch record.Status {
	case AuthorizationStatusProposed:
		if !record.GrantedAt.IsZero() {
			return record, authorizationReason(
				AuthorizationReasonContradictory,
				"status,granted",
				record.GrantedAt.Format("2006-01-02"),
				"status proposed contradicts a non-null granted date",
			)
		}
		return record, authorizationReason(
			AuthorizationReasonStatus,
			"status",
			string(record.Status),
			"status proposed is not an operative grant",
		)
	case AuthorizationStatusWithdrawn:
		return record, authorizationReason(
			AuthorizationReasonStatus,
			"status",
			string(record.Status),
			"status withdrawn is not an operative grant",
		)
	case AuthorizationStatusApproved:
	default:
		return record, authorizationReason(
			AuthorizationReasonStatus,
			"status",
			string(record.Status),
			fmt.Sprintf("status %q is not approved, proposed, or withdrawn", record.Status),
		)
	}
	if !presence.granted || presence.grantedNil || record.GrantedAt.IsZero() {
		return record, authorizationReason(
			AuthorizationReasonGranted,
			"granted",
			"",
			"granted must contain a YYYY-MM-DD date for an approved record",
		)
	}
	if record.Action == "" {
		return record, authorizationReason(
			AuthorizationReasonAction,
			"action",
			"",
			"action is required for an approved record",
		)
	}
	if len(record.Consuming) == 0 {
		return record, authorizationReason(
			AuthorizationReasonConsuming,
			"consuming",
			"",
			"consuming must name at least one Spec or direct action",
		)
	}
	if len(record.Paths) == 0 {
		return record, authorizationReason(
			AuthorizationReasonPaths,
			"paths",
			"",
			"paths must contain at least one exact repository-relative path",
		)
	}
	return record, AuthorizationReason{}
}

var (
	legacyFullSpecPattern = regexp.MustCompile(`\b[0-9]{4}-[a-z0-9][a-z0-9-]*\b`)
	legacySpecIDPattern   = regexp.MustCompile(`(^|[^0-9])([0-9]{4})([^0-9]|$)`)
	legacyBacktickPattern = regexp.MustCompile("`([^`]+)`")
)

func parseLegacyAuthorization(
	source AuthorizationSource,
	content []byte,
) (AuthorizationRecord, AuthorizationReason) {
	record := AuthorizationRecord{
		Role:   AuthorizationRoleLegacy,
		Status: AuthorizationStatusApproved,
		Source: source,
	}
	base := path.Base(source.Path)
	if len(base) >= len("2006-01-02") {
		record.GrantedAt, _ = time.Parse("2006-01-02", base[:len("2006-01-02")])
	}
	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 {
		record.Action = legacyAuthorizationAction(lines[0])
	}
	record.Consuming = legacyAuthorizationConsumers(lines)
	record.Paths = legacyAuthorizationPaths(lines)
	regenerations, reason := parseAuthorizationRegenerations(content)
	record.Regenerations = regenerations
	if reason.Code != "" {
		return record, reason
	}
	return classifyAuthorizationRecord(record, authorizationFrontmatterPresence{
		status:  true,
		granted: !record.GrantedAt.IsZero(),
	})
}

func legacyAuthorizationAction(title string) string {
	title = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(title), "#"))
	if _, action, ok := strings.Cut(title, "—"); ok {
		title = strings.TrimSpace(action)
	}
	if before, _, ok := strings.Cut(title, " ("); ok {
		title = strings.TrimSpace(before)
	}
	return title
}

func legacyAuthorizationConsumers(lines []string) []string {
	var consumers []string
	seen := make(map[string]struct{})
	inConsuming := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			heading := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")))
			inConsuming = heading == "consuming spec" ||
				heading == "consuming specs" ||
				heading == "consuming work" ||
				heading == "which spec uses what"
			continue
		}
		if !inConsuming {
			continue
		}
		remaining := line
		for _, consumer := range legacyFullSpecPattern.FindAllString(line, -1) {
			appendUniqueAuthorizationValue(&consumers, seen, consumer)
			remaining = strings.ReplaceAll(remaining, consumer, strings.Repeat(" ", len(consumer)))
		}
		for _, match := range legacySpecIDPattern.FindAllStringSubmatch(remaining, -1) {
			appendUniqueAuthorizationValue(&consumers, seen, match[2])
		}
		if strings.Contains(strings.ToLower(line), "applied directly") {
			appendUniqueAuthorizationValue(&consumers, seen, "direct")
		}
	}
	return consumers
}

func legacyAuthorizationPaths(lines []string) []string {
	var paths []string
	seen := make(map[string]struct{})
	inPaths := false
	lastSourceBaselineRoot := ""
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			heading := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")))
			inPaths = heading == "authorized paths" || heading == "bounded files" || heading == "authorized scope"
			continue
		}
		if !inPaths {
			continue
		}
		for _, match := range legacyBacktickPattern.FindAllStringSubmatch(line, -1) {
			candidate := match[1]
			if candidate == "manifest.json" &&
				lastSourceBaselineRoot != "" &&
				strings.Contains(line, "same Source Baseline") {
				candidate = path.Join(lastSourceBaselineRoot, candidate)
			}
			if !looksLikeLegacyRepositoryPath(candidate) {
				continue
			}
			appendUniqueAuthorizationValue(&paths, seen, candidate)
			if corpusIndex := strings.Index(candidate, "/corpus/"); corpusIndex >= 0 &&
				strings.Contains(candidate[:corpusIndex], "/source-baselines/") {
				lastSourceBaselineRoot = candidate[:corpusIndex]
			}
		}
	}
	return paths
}

func looksLikeLegacyRepositoryPath(value string) bool {
	if strings.ContainsAny(value, " \t\r\n") || strings.HasPrefix(value, "make ") {
		return false
	}
	return value == "Makefile" ||
		strings.Contains(value, "/") ||
		strings.HasPrefix(value, ".")
}

func splitLegacyConsumers(value string) []string {
	parts := strings.Split(value, ",")
	consumers := make([]string, 0, len(parts))
	for _, part := range parts {
		if consumer := strings.TrimSpace(part); consumer != "" {
			consumers = append(consumers, consumer)
		}
	}
	return consumers
}

func trimNonEmptyAuthorizationValues(values []string) []string {
	trimmed := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			trimmed = append(trimmed, value)
		}
	}
	return trimmed
}

func appendUniqueAuthorizationValue(values *[]string, seen map[string]struct{}, value string) {
	if _, duplicate := seen[value]; duplicate {
		return
	}
	seen[value] = struct{}{}
	*values = append(*values, value)
}

func authorizationFrontmatterMapping(frontmatter []byte) (map[string]*yaml.Node, error) {
	var document yaml.Node
	if err := yaml.Unmarshal(frontmatter, &document); err != nil {
		return nil, err
	}
	if document.Kind != yaml.DocumentNode || len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return nil, errors.New("frontmatter must be a YAML mapping")
	}
	mapping := document.Content[0]
	fields := make(map[string]*yaml.Node, len(mapping.Content)/2)
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		key := mapping.Content[index]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
			return nil, errors.New("frontmatter keys must be strings")
		}
		if _, duplicate := fields[key.Value]; duplicate {
			return nil, fmt.Errorf("frontmatter field %q is declared more than once", key.Value)
		}
		fields[key.Value] = mapping.Content[index+1]
	}
	return fields, nil
}

func authorizationScalar(node *yaml.Node) (string, error) {
	if node.Tag == "!!null" {
		return "", errors.New("must not be null")
	}
	if node.Kind != yaml.ScalarNode || node.Tag != "!!str" {
		return "", errors.New("must be a string")
	}
	return node.Value, nil
}

func authorizationDateScalar(node *yaml.Node) (string, error) {
	if node.Kind != yaml.ScalarNode || (node.Tag != "!!str" && node.Tag != "!!timestamp") {
		return "", errors.New("must be a date string")
	}
	return node.Value, nil
}

func authorizationStringValues(node *yaml.Node) ([]string, error) {
	if node.Tag == "!!null" {
		return nil, nil
	}
	if node.Kind == yaml.ScalarNode {
		value, err := authorizationScalar(node)
		if err != nil {
			return nil, err
		}
		return []string{value}, nil
	}
	return authorizationStringSequence(node)
}

func authorizationStringSequence(node *yaml.Node) ([]string, error) {
	if node.Kind != yaml.SequenceNode {
		return nil, errors.New("must be a list")
	}
	values := make([]string, 0, len(node.Content))
	for _, item := range node.Content {
		value, err := authorizationScalar(item)
		if err != nil {
			return nil, fmt.Errorf("list item: %w", err)
		}
		values = append(values, value)
	}
	return values, nil
}

func fieldShapeReason(field string, err error) AuthorizationReason {
	return authorizationReason(
		AuthorizationReasonMalformedRecord,
		field,
		"",
		fmt.Sprintf("%s %v", field, err),
	)
}

func parseAuthorizationRegenerations(content []byte) ([]AuthorizationRegeneration, AuthorizationReason) {
	lines := strings.Split(string(content), "\n")
	inSection := false
	var regenerations []AuthorizationRegeneration
	for index := 0; index < len(lines); index++ {
		trimmed := strings.TrimSpace(lines[index])
		if strings.HasPrefix(trimmed, "## ") {
			inSection = strings.EqualFold(
				strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")),
				"Sanctioned regeneration",
			)
			continue
		}
		if !inSection || trimmed != "```yaml" {
			continue
		}
		start := index + 1
		for index++; index < len(lines) && strings.TrimSpace(lines[index]) != "```"; index++ {
		}
		if index >= len(lines) {
			return regenerations, authorizationReason(
				AuthorizationReasonRegeneration,
				"sanctioned_regeneration",
				"",
				"sanctioned regeneration YAML block has no closing fence",
			)
		}
		regeneration, err := decodeAuthorizationRegeneration(strings.Join(lines[start:index], "\n"))
		if err != nil {
			return regenerations, authorizationReason(
				AuthorizationReasonRegeneration,
				"sanctioned_regeneration",
				"",
				fmt.Sprintf("decode sanctioned regeneration: %v", err),
			)
		}
		regenerations = append(regenerations, regeneration)
	}
	return regenerations, AuthorizationReason{}
}

func decodeAuthorizationRegeneration(content string) (AuthorizationRegeneration, error) {
	var decoded struct {
		Command string    `yaml:"command"`
		Outputs *[]string `yaml:"outputs"`
	}
	decoder := yaml.NewDecoder(strings.NewReader(content))
	decoder.KnownFields(true)
	if err := decoder.Decode(&decoded); err != nil {
		return AuthorizationRegeneration{}, err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return AuthorizationRegeneration{}, errors.New("multiple YAML documents")
		}
		return AuthorizationRegeneration{}, err
	}
	decoded.Command = strings.TrimSpace(decoded.Command)
	if decoded.Command == "" {
		return AuthorizationRegeneration{}, errors.New("command is required")
	}
	regeneration := AuthorizationRegeneration{Command: decoded.Command}
	if decoded.Outputs == nil {
		return regeneration, nil
	}
	if len(*decoded.Outputs) == 0 {
		return AuthorizationRegeneration{}, errors.New("outputs must not be empty when declared")
	}
	regeneration.Outputs = append([]string(nil), (*decoded.Outputs)...)
	return regeneration, nil
}

func validateAuthorizationPaths(
	ctx context.Context,
	repoRoot string,
	revision string,
	paths []string,
) (AuthorizationReason, error) {
	seen := make(map[string]struct{}, len(paths))
	for _, declared := range paths {
		if detail := invalidAuthorizationPath(declared); detail != "" {
			return authorizationReason(
				AuthorizationReasonPaths,
				"paths",
				declared,
				fmt.Sprintf("path %q %s", declared, detail),
			), nil
		}
		if _, duplicate := seen[declared]; duplicate {
			return authorizationReason(
				AuthorizationReasonPaths,
				"paths",
				declared,
				fmt.Sprintf("path %q is declared more than once", declared),
			), nil
		}
		seen[declared] = struct{}{}

		var symlink string
		var err error
		if revision == "" {
			symlink, err = authorizationFilesystemSymlink(repoRoot, declared)
		} else {
			symlink, err = authorizationGitSymlink(ctx, repoRoot, revision, declared)
		}
		if err != nil {
			return AuthorizationReason{}, err
		}
		if symlink != "" {
			return authorizationReason(
				AuthorizationReasonPaths,
				"paths",
				declared,
				fmt.Sprintf("path %q resolves through symlink %q", declared, symlink),
			), nil
		}
	}
	return AuthorizationReason{}, nil
}

func invalidAuthorizationPath(value string) string {
	if value == "" {
		return "is empty"
	}
	if strings.TrimSpace(value) != value {
		return "has leading or trailing whitespace"
	}
	if filepath.IsAbs(value) || path.IsAbs(value) || filepath.VolumeName(value) != "" {
		return "is absolute"
	}
	if strings.ContainsAny(value, "\\\x00\r\n") {
		return "contains an unsupported separator or control character"
	}
	if strings.ContainsAny(value, "*?[{}<>") {
		return "contains a glob or placeholder"
	}
	components := strings.Split(value, "/")
	for _, component := range components {
		if component == ".." {
			return "traverses upward"
		}
		if component == "" || component == "." {
			return "is not in canonical repository-relative form"
		}
	}
	if clean := path.Clean(value); clean != value || clean == "." {
		return "is not in canonical repository-relative form"
	}
	return ""
}

func authorizationFilesystemSymlink(repoRoot string, relative string) (string, error) {
	current := repoRoot
	for _, component := range strings.Split(relative, "/") {
		current = filepath.Join(current, filepath.FromSlash(component))
		info, err := os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		if err != nil {
			return "", fmt.Errorf("stat %q: %w", filepath.ToSlash(current), err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			relativeSymlink, relErr := filepath.Rel(repoRoot, current)
			if relErr != nil {
				return "", fmt.Errorf("make symlink %q relative to repository: %w", current, relErr)
			}
			return filepath.ToSlash(relativeSymlink), nil
		}
	}
	return "", nil
}

func authorizationGitSymlink(ctx context.Context, repoRoot, revision, relative string) (string, error) {
	components := strings.Split(relative, "/")
	for index := range components {
		candidate := strings.Join(components[:index+1], "/")
		output, err := runAuthorizationGit(ctx, repoRoot, "ls-tree", revision, "--", candidate)
		if err != nil {
			return "", err
		}
		if len(output) == 0 {
			return "", nil
		}
		metadata, _, found := bytes.Cut(output, []byte{'\t'})
		if !found {
			return "", fmt.Errorf("parse Git tree entry for %q", candidate)
		}
		fields := strings.Fields(string(metadata))
		if len(fields) < 1 {
			return "", fmt.Errorf("parse Git tree mode for %q", candidate)
		}
		if fields[0] == "120000" {
			return candidate, nil
		}
	}
	return "", nil
}

func allowedAuthorizationOperation(operation AuthorizationOperation) bool {
	for _, allowed := range authorizationOperations {
		if operation == allowed {
			return true
		}
	}
	return false
}

func readWorkingAuthorization(repoRoot, recordPath string) ([]byte, error) {
	if symlink, err := authorizationFilesystemSymlink(repoRoot, recordPath); err != nil {
		return nil, err
	} else if symlink != "" {
		return nil, fmt.Errorf("record resolves through symlink %q", symlink)
	}
	absPath := filepath.Join(repoRoot, filepath.FromSlash(recordPath))
	info, err := os.Lstat(absPath)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errors.New("record is not a regular file")
	}
	return os.ReadFile(absPath)
}

func resolveAuthorizationRevision(ctx context.Context, repoRoot, revision string) (string, error) {
	if strings.ContainsAny(revision, "\x00\r\n") {
		return "", errors.New("revision contains an unsupported control character")
	}
	output, err := runAuthorizationGit(
		ctx,
		repoRoot,
		"rev-parse",
		"--verify",
		"--end-of-options",
		strings.TrimSpace(revision)+"^{commit}",
	)
	if err != nil {
		return "", err
	}
	resolved := strings.TrimSpace(string(output))
	if resolved == "" || strings.Contains(resolved, "\n") {
		return "", fmt.Errorf("Git returned invalid revision %q", resolved)
	}
	return resolved, nil
}

func readAuthorizationRevision(
	ctx context.Context,
	repoRoot string,
	revision string,
	recordPath string,
) ([]byte, error) {
	entry, err := runAuthorizationGit(ctx, repoRoot, "ls-tree", revision, "--", recordPath)
	if err != nil {
		return nil, err
	}
	metadata, _, found := bytes.Cut(entry, []byte{'\t'})
	if !found {
		return nil, errors.New("record is absent from the revision")
	}
	fields := strings.Fields(string(metadata))
	if len(fields) < 2 || fields[0] == "120000" || fields[1] != "blob" {
		return nil, errors.New("record is not a regular Git blob")
	}
	return runAuthorizationGit(ctx, repoRoot, "cat-file", "blob", revision+":"+recordPath)
}

func runAuthorizationGit(ctx context.Context, repoRoot string, args ...string) ([]byte, error) {
	gitArgs := append([]string{"-C", repoRoot, "-c", "core.fsmonitor=false"}, args...)
	command := exec.CommandContext(ctx, "git", gitArgs...)
	command.Env = authorizationGitEnvironment()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return nil, fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), detail, err)
	}
	return append([]byte(nil), stdout.Bytes()...), nil
}

func authorizationGitEnvironment() []string {
	environment := make([]string, 0, len(os.Environ())+1)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "GIT_OPTIONAL_LOCKS=") {
			continue
		}
		environment = append(environment, entry)
	}
	return append(environment, "GIT_OPTIONAL_LOCKS=0")
}

func cleanAuthorizationSourcePath(value string) (string, bool) {
	if invalidAuthorizationPath(value) != "" || strings.Contains(value, ":") {
		return "", false
	}
	return value, true
}

func pathWithinAuthorizationDirectory(candidate, root string) bool {
	return candidate != root && strings.HasPrefix(candidate, root+"/")
}

func authorizationReason(
	code AuthorizationReasonCode,
	field string,
	value string,
	detail string,
) AuthorizationReason {
	return AuthorizationReason{Code: code, Field: field, Value: value, Detail: detail}
}

func refusedAuthorization(
	record AuthorizationRecord,
	code AuthorizationReasonCode,
	field string,
	value string,
	detail string,
) AuthorizationResolution {
	return AuthorizationResolution{
		Outcome: AuthorizationRefused,
		Record:  record,
		Reason:  authorizationReason(code, field, value, detail),
	}
}

func unresolvedAuthorization(
	record AuthorizationRecord,
	code AuthorizationReasonCode,
	field string,
	value string,
	detail string,
) AuthorizationResolution {
	return AuthorizationResolution{
		Outcome: AuthorizationUnresolved,
		Record:  record,
		Reason:  authorizationReason(code, field, value, detail),
	}
}
