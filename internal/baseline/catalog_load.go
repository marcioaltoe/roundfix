package baseline

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

type catalogLoader struct {
	source      fs.FS
	assets      map[string][]byte
	documents   map[string]document
	diagnostics []Diagnostic
}

func newCatalogLoader(source fs.FS) *catalogLoader {
	return &catalogLoader{
		source:    source,
		assets:    make(map[string][]byte),
		documents: make(map[string]document),
	}
}

func (l *catalogLoader) load() *Catalog {
	l.readAssets()

	catalog := &Catalog{
		assets:         l.assets,
		profiles:       l.readCollection("profiles", profileSchemas, "profile", true),
		modules:        l.readCollection("modules", moduleSchemas, "module", true),
		setups:         l.readCollection("setups", setupSchemas, "setup", true),
		transitions:    l.readCollection("retention", transitionSchemas, "transition", false),
		orderedModules: make(map[string][]string),
		semanticOwners: make(map[string]SemanticOwner),
	}
	catalog.decisions = l.readIndexed(
		"decisions.json",
		"decisions",
		[]string{"setup-context-driven/decisions-v1"},
		"decision",
	)
	catalog.templates = l.readIndexed(
		"templates/index.json",
		"templates",
		[]string{"setup-context-driven/templates-v1"},
		"template",
	)

	l.validateRequiredDocuments()
	l.validateSchemaFields(catalog)
	l.validateVersions("module", catalog.modules)
	l.validateVersions("profile", catalog.profiles)
	l.validateVersions("setup", catalog.setups)
	l.validateVersions("transition", catalog.transitions)
	l.validateVersions("decision", catalog.decisions)
	l.validateVersions("template", catalog.templates)
	l.validateTemplates(catalog)
	l.validateModules(catalog)
	l.validateModuleCycles(catalog)
	l.validateProfiles(catalog)
	l.validateProfileFormatters(catalog)
	l.validateGuidanceComposition(catalog)
	l.validateDecisionEffects(catalog)
	l.validateSetups(catalog)
	l.validateSkillActivations(catalog)
	l.validateTransitions(catalog)
	l.validateSourceBaselines(catalog)

	if len(l.diagnostics) == 0 {
		if err := catalog.finishIdentity(); err != nil {
			l.add("catalog.identity.invalid", "", err.Error())
		}
	}
	return catalog
}

func (l *catalogLoader) validateSchemaFields(catalog *Catalog) {
	l.allowDocumentFields("contract-v1.json", "schemaVersion", "version", "contracts", "diagnostics")
	l.allowDocumentFields("coverage.json", "schemaVersion", "version", "coverage")
	l.allowDocumentFields("decisions.json", "schemaVersion", "version", "decisions")
	l.allowDocumentFields(
		"lock-hash-compatibility-v1.json",
		"schemaVersion",
		"version",
		"expectedSha256",
		"files",
	)
	l.allowDocumentFields("skill-activations.json", "schemaVersion", "version", "bundles", "activations")
	l.allowDocumentFields("source-baselines/index.json", "schemaVersion", "version", "baselines")
	l.allowDocumentFields("templates/index.json", "schemaVersion", "version", "templates")

	moduleFields := []string{
		"schemaVersion", "id", "version", "kind", "title", "dependsOn",
		"conflictsWith", "rootBlocks", "supportingGuides", "rules",
		"requiredSkills", "skillDispatch", "requiredDecisions",
		"repositoryExtensions", "instructionHierarchy", "semanticOwners",
	}
	for id, doc := range catalog.modules {
		l.allowFields("catalog.module.field.unknown", id, doc, moduleFields...)
	}
	profileFields := []string{
		"schemaVersion", "id", "version", "markerVersion", "title", "setup",
		"entryDecisions", "modules", "requiredRules", "formatter", "stack",
		"workspaces", "optionalModules", "architecture", "httpContract",
		"capabilitySets", "capabilities", "activationBundles", "verification",
	}
	for id, doc := range catalog.profiles {
		l.allowFields("catalog.profile.field.unknown", id, doc, profileFields...)
	}
	for id, doc := range catalog.setups {
		l.allowFields(
			"catalog.setup.field.unknown",
			id,
			doc,
			"schemaVersion",
			"id",
			"version",
			"source",
			"digest",
			"skills",
			"activationBundles",
		)
	}
	for id, doc := range catalog.transitions {
		l.allowFields(
			"catalog.transition.field.unknown",
			id,
			doc,
			"schemaVersion",
			"id",
			"version",
			"fromBaseline",
			"toBaseline",
			"legacyManifestFingerprints",
			"priorClauses",
			"mappings",
		)
	}
	for id, doc := range catalog.decisions {
		l.allowFields(
			"catalog.decision.field.unknown",
			id,
			doc,
			"id",
			"version",
			"type",
			"default",
			"summary",
			"effects",
			"values",
			"modes",
			"suggestion",
			"requiresCapabilities",
		)
	}
	for id, doc := range catalog.templates {
		l.allowFields(
			"catalog.template.field.unknown",
			id,
			doc,
			"id",
			"version",
			"kind",
			"path",
			"tokens",
		)
	}
}

func (l *catalogLoader) allowDocumentFields(assetPath string, fields ...string) {
	if doc := l.documents[assetPath]; doc != nil {
		l.allowFields("catalog.document.field.unknown", assetPath, doc, fields...)
	}
}

func (l *catalogLoader) allowFields(code, owner string, doc document, fields ...string) {
	allowed := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		allowed[field] = struct{}{}
	}
	for field := range doc {
		if _, ok := allowed[field]; !ok {
			l.add(code, owner, field)
		}
	}
}

func (l *catalogLoader) readAssets() {
	err := fs.WalkDir(l.source, ".", func(assetPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if !safeRelative(assetPath) {
			l.add("catalog.asset.path.invalid", assetPath, "")
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect %s: %w", assetPath, err)
		}
		if !info.Mode().IsRegular() {
			l.add("catalog.asset.kind.invalid", assetPath, info.Mode().String())
			return nil
		}
		data, err := fs.ReadFile(l.source, assetPath)
		if err != nil {
			return fmt.Errorf("read %s: %w", assetPath, err)
		}
		l.assets[assetPath] = append([]byte(nil), data...)
		if path.Ext(assetPath) == ".json" {
			doc, diagnostics := decodeDocument(data, assetPath)
			l.diagnostics = append(l.diagnostics, diagnostics...)
			if doc != nil {
				l.documents[assetPath] = doc
			}
		}
		return nil
	})
	if err != nil {
		l.add("catalog.assets.read.failed", "", err.Error())
	}
}

func (l *catalogLoader) validateRequiredDocuments() {
	required := map[string]string{
		"contract-v1.json":     "setup-context-driven/assets-v1",
		"decisions.json":       "setup-context-driven/decisions-v1",
		"templates/index.json": "setup-context-driven/templates-v1",
	}
	for assetPath, schema := range required {
		doc := l.documents[assetPath]
		if doc == nil {
			l.add("catalog.asset.required.missing", assetPath, "")
			continue
		}
		if got, _ := stringValue(doc, "schemaVersion"); got != schema {
			l.add("catalog.schema.invalid", assetPath, fmt.Sprintf("got %q, want %q", got, schema))
		}
	}
	optional := map[string][]string{
		"coverage.json": {
			"setup-context-driven/coverage-v1",
			"setup-context-driven/coverage-v2",
		},
		"lock-hash-compatibility-v1.json": {
			"setup-context-driven/external-lock-hash-compatibility-v1",
		},
		"skill-activations.json": {
			"setup-context-driven/skill-activations-v1",
		},
		"source-baselines/index.json": {
			"setup-context-driven/source-baseline-index/0.0.1",
		},
	}
	for assetPath, schemas := range optional {
		if doc := l.documents[assetPath]; doc != nil {
			l.validateSchema(assetPath, doc, schemas)
		}
	}
}

func (l *catalogLoader) readCollection(
	dir string,
	schemas []string,
	kind string,
	required bool,
) map[string]document {
	result := make(map[string]document)
	prefix := dir + "/"
	for _, assetPath := range sortedKeys(l.documents) {
		if !strings.HasPrefix(assetPath, prefix) ||
			strings.Contains(strings.TrimPrefix(assetPath, prefix), "/") {
			continue
		}
		doc := l.documents[assetPath]
		l.validateSchema(assetPath, doc, schemas)
		id, ok := stringValue(doc, "id")
		if !ok {
			l.add("catalog."+kind+".id.missing", assetPath, "")
			continue
		}
		if want := id + ".json"; path.Base(assetPath) != want {
			l.add("catalog."+kind+".filename.mismatch", assetPath, want)
		}
		if _, exists := result[id]; exists {
			l.add("catalog."+kind+".id.duplicate", assetPath, id)
		}
		result[id] = doc
	}
	if required && len(result) == 0 {
		l.add("catalog."+kind+".collection.empty", dir, "")
	}
	return result
}

func (l *catalogLoader) readIndexed(
	assetPath string,
	key string,
	schemas []string,
	kind string,
) map[string]document {
	result := make(map[string]document)
	doc := l.documents[assetPath]
	if doc == nil {
		l.add("catalog.asset.required.missing", assetPath, "")
		return result
	}
	l.validateSchema(assetPath, doc, schemas)
	items, ok := objectList(doc[key])
	if !ok {
		l.add("catalog."+kind+".collection.invalid", assetPath, key)
		return result
	}
	for index, item := range items {
		id, ok := stringValue(item, "id")
		if !ok {
			l.add("catalog."+kind+".id.missing", assetPath, fmt.Sprintf("index %d", index))
			continue
		}
		if _, exists := result[id]; exists {
			l.add("catalog."+kind+".id.duplicate", assetPath, id)
		}
		result[id] = item
	}
	return result
}

func (l *catalogLoader) validateSchema(assetPath string, doc document, accepted []string) {
	schema, ok := stringValue(doc, "schemaVersion")
	if !ok {
		l.add("catalog.schema.missing", assetPath, "")
		return
	}
	if !containsString(accepted, schema) {
		l.add("catalog.schema.invalid", assetPath, schema)
	}
}

func (l *catalogLoader) validateVersions(kind string, values map[string]document) {
	for id, doc := range values {
		schema, _ := stringValue(doc, "schemaVersion")
		if schema == "setup-context-driven/profile/0.0.1" ||
			schema == "setup-context-driven/setup-snapshot/0.0.1" {
			if version, ok := doc["version"].(string); !ok || version != "0.0.1" {
				l.add("catalog."+kind+".version.invalid", id, "")
			}
			continue
		}
		if version, ok := integerValue(doc, "version"); !ok || version < 1 {
			l.add("catalog."+kind+".version.invalid", id, "")
		}
	}
}

func (l *catalogLoader) add(code, assetPath, info string) {
	l.diagnostics = append(l.diagnostics, Diagnostic{Code: code, Path: assetPath, Info: info})
}

func containsString(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}

var (
	moduleSchemas = []string{
		"setup-context-driven/module-v1",
		"setup-context-driven/module-v2",
		"setup-context-driven/module-v3",
	}
	profileSchemas = []string{
		"setup-context-driven/profile-v1",
		"setup-context-driven/profile-v2",
		"setup-context-driven/profile-v3",
		"setup-context-driven/profile/0.0.1",
	}
	setupSchemas = []string{
		"setup-context-driven/setup-snapshot-v1",
		"setup-context-driven/setup-snapshot-v2",
		"setup-context-driven/setup-snapshot/0.0.1",
	}
	transitionSchemas = []string{"setup-context-driven/upgrade-transition-v1"}
)

func uniqueStrings(values []string) bool {
	sorted := append([]string(nil), values...)
	sort.Strings(sorted)
	for index := 1; index < len(sorted); index++ {
		if sorted[index-1] == sorted[index] {
			return false
		}
	}
	return true
}
