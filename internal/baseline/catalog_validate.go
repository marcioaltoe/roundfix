package baseline

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strings"
)

var (
	templateToken  = regexp.MustCompile(`\{\{\s*([A-Za-z0-9_.-]+)\s*\}\}`)
	lowerSHA256    = regexp.MustCompile(`^[0-9a-f]{64}$`)
	immutableGitID = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)
)

var confirmedInstructionHierarchy = []string{
	"universal",
	"context",
	"spec",
	"autonomous",
	"stack",
	"surface",
	"optional-knowledge",
	"repository-specific",
}

func (l *catalogLoader) validateTemplates(catalog *Catalog) {
	for templateID, template := range catalog.templates {
		requireFields(l, "template", templateID, template, "id", "version", "kind", "path")
		templatePath, ok := stringValue(template, "path")
		if !ok || !safeRelative(templatePath) {
			l.add("catalog.template.path.invalid", templateID, "")
			continue
		}
		fullPath := path.Join("templates", templatePath)
		content, ok := catalog.assets[fullPath]
		if !ok {
			l.add("catalog.template.file.missing", templateID, templatePath)
			continue
		}
		declared, ok := stringList(template["tokens"])
		if !ok {
			l.add("catalog.template.tokens.invalid", templateID, "")
			continue
		}
		if !uniqueStrings(declared) {
			l.add("catalog.template.token.duplicate", templateID, "")
		}
		declaredSet := stringSet(declared)
		for _, match := range templateToken.FindAllSubmatch(content, -1) {
			token := string(match[1])
			if _, ok := declaredSet[token]; !ok {
				l.add("catalog.template.token.undeclared", templateID, token)
			}
		}
	}
}

func (l *catalogLoader) validateModules(catalog *Catalog) {
	decisions := stringSet(catalog.DecisionIDs())
	templates := catalog.templates
	coverage := make(map[string]document)
	if l.documents["coverage.json"] != nil {
		coverage = l.indexSimpleIDs("coverage.json", "coverage", "coverage")
	} else {
		for moduleID, module := range catalog.modules {
			schema, _ := stringValue(module, "schemaVersion")
			if schema == "setup-context-driven/module-v2" ||
				schema == "setup-context-driven/module-v3" {
				l.add("catalog.coverage.missing", moduleID, "")
			}
		}
	}
	seenRules := make(map[string]string)
	seenBlocks := make(map[string]string)
	seenGuides := make(map[string]string)
	seenClauses := make(map[string]string)
	seenDispatch := make(map[string]string)
	seenTriggers := make(map[string]string)
	seenReferences := make(map[string]string)
	seenExtensions := make(map[string]string)

	for moduleID, module := range catalog.modules {
		requireFields(
			l,
			"module",
			moduleID,
			module,
			"id",
			"version",
			"dependsOn",
			"conflictsWith",
			"rootBlocks",
			"supportingGuides",
			"rules",
			"requiredSkills",
			"requiredDecisions",
		)
		for _, field := range []string{"dependsOn", "conflictsWith", "requiredSkills", "requiredDecisions"} {
			values, ok := stringList(module[field])
			if !ok {
				l.add("catalog.module."+field+".invalid", moduleID, "")
				continue
			}
			if !uniqueStrings(values) {
				l.add("catalog.module."+field+".duplicate", moduleID, "")
			}
		}
		for _, dependency := range stringsOrEmpty(module["dependsOn"]) {
			if _, ok := catalog.modules[dependency]; !ok {
				l.add("catalog.module.dependency.unknown", moduleID, dependency)
			}
		}
		for _, conflict := range stringsOrEmpty(module["conflictsWith"]) {
			if _, ok := catalog.modules[conflict]; !ok {
				l.add("catalog.module.conflict.unknown", moduleID, conflict)
			}
		}
		for _, decision := range stringsOrEmpty(module["requiredDecisions"]) {
			if _, ok := decisions[decision]; !ok {
				l.add("catalog.module.decision.unknown", moduleID, decision)
			}
		}

		for _, rule := range objectsOrDiagnostic(l, module["rules"], "catalog.module.rules.invalid", moduleID) {
			ruleID := validateOwnedID(l, "rule", moduleID, rule, seenRules)
			if ruleID == "" {
				continue
			}
			if version, ok := integerValue(rule, "version"); !ok || version < 1 {
				l.add("catalog.rule.version.invalid", ruleID, "")
			}
			schema, _ := stringValue(module, "schemaVersion")
			if schema == "setup-context-driven/module-v2" ||
				schema == "setup-context-driven/module-v3" {
				coverageIDs, ok := stringList(rule["coverage"])
				if !ok || len(coverageIDs) == 0 || !uniqueStrings(coverageIDs) {
					l.add("catalog.rule.coverage.invalid", ruleID, "")
				}
				for _, coverageID := range coverageIDs {
					if _, ok := coverage[coverageID]; !ok {
						l.add("catalog.rule.coverage.unknown", ruleID, coverageID)
					}
				}
			}
			clauses, ok := objectList(rule["clauses"])
			if schema == "setup-context-driven/module-v3" {
				if !ok || len(clauses) == 0 {
					l.add("catalog.rule.clauses.invalid", ruleID, "")
				}
			}
			for _, clause := range clauses {
				clauseID := validateOwnedID(l, "clause", ruleID, clause, seenClauses)
				enforcement, _ := stringValue(clause, "enforcement")
				if !containsString([]string{"mandatory", "prohibited", "stop-and-ask"}, enforcement) {
					l.add("catalog.clause.enforcement.invalid", clauseID, enforcement)
				}
				if guidance, _ := stringValue(clause, "guidance"); strings.TrimSpace(guidance) == "" {
					l.add("catalog.clause.guidance.invalid", clauseID, "")
				}
			}
		}

		l.validateArtifacts(moduleID, module, "rootBlocks", "root-block", templates, seenBlocks, seenReferences)
		l.validateArtifacts(moduleID, module, "supportingGuides", "guide", templates, seenGuides, seenReferences)

		for _, dispatch := range objectsOrEmpty(module["skillDispatch"]) {
			dispatchID, ok := stringValue(dispatch, "skill")
			if !ok {
				dispatchID, ok = stringValue(dispatch, "id")
			}
			if !ok {
				l.add("catalog.skill.dispatch.id.missing", moduleID, "")
				continue
			}
			if owner, exists := seenDispatch[dispatchID]; exists {
				l.add("catalog.skill.dispatch.id.duplicate", dispatchID, owner+", "+moduleID)
			}
			seenDispatch[dispatchID] = moduleID
			triggers := objectsOrEmpty(dispatch["triggers"])
			if when, _ := stringValue(dispatch, "when"); strings.TrimSpace(when) == "" && len(triggers) == 0 {
				l.add("catalog.skill.dispatch.when.invalid", dispatchID, "")
			}
			for _, trigger := range triggers {
				triggerID := validateOwnedID(l, "skill.trigger", dispatchID, trigger, seenTriggers)
				if when, _ := stringValue(trigger, "when"); strings.TrimSpace(when) == "" {
					l.add("catalog.skill.trigger.when.invalid", triggerID, "")
				}
			}
		}

		for _, extension := range objectsOrEmpty(module["repositoryExtensions"]) {
			extensionID := validateOwnedID(l, "repository.extension", moduleID, extension, seenExtensions)
			targetPath, ok := stringValue(extension, "path")
			if !ok || !safeRelative(targetPath) {
				l.add("catalog.repository.extension.path.invalid", extensionID, targetPath)
			}
			templateID, _ := stringValue(extension, "template")
			if template, ok := templates[templateID]; !ok {
				l.add("catalog.repository.extension.template.unknown", extensionID, templateID)
			} else if kind, _ := stringValue(template, "kind"); kind != "repository-extension" {
				l.add("catalog.repository.extension.template.kind", extensionID, templateID)
			}
			decisionID, _ := stringValue(extension, "decision")
			if _, ok := decisions[decisionID]; !ok {
				l.add("catalog.repository.extension.decision.unknown", extensionID, decisionID)
			}
		}
	}
}

func (l *catalogLoader) validateArtifacts(
	moduleID string,
	module document,
	field string,
	expectedKind string,
	templates map[string]document,
	seen map[string]string,
	seenReferences map[string]string,
) {
	artifacts := objectsOrDiagnostic(l, module[field], "catalog.module."+field+".invalid", moduleID)
	for _, artifact := range artifacts {
		kind := strings.TrimSuffix(expectedKind, "-block")
		artifactID := validateOwnedID(l, kind, moduleID, artifact, seen)
		if artifactID == "" {
			continue
		}
		templateID, ok := stringValue(artifact, "template")
		if !ok {
			l.add("catalog.template.reference.missing", artifactID, "")
		} else if template, exists := templates[templateID]; !exists {
			l.add("catalog.template.reference.unknown", artifactID, templateID)
		} else if got, _ := stringValue(template, "kind"); got != expectedKind {
			l.add("catalog.template.kind.mismatch", artifactID, fmt.Sprintf("%s != %s", got, expectedKind))
		}
		if field == "supportingGuides" {
			targetPath, ok := stringValue(artifact, "path")
			if !ok || !safeRelative(targetPath) {
				l.add("catalog.guide.path.invalid", artifactID, targetPath)
			}
		}
		for _, reference := range objectsOrEmpty(artifact["references"]) {
			referenceID := validateOwnedID(l, "reference", artifactID, reference, seenReferences)
			token, _ := stringValue(reference, "token")
			templateTokens := stringsOrEmpty(templates[templateID]["tokens"])
			if !containsString(templateTokens, token) {
				l.add("catalog.reference.token.unknown", referenceID, token)
			}
		}
	}
}

func (l *catalogLoader) validateModuleCycles(catalog *Catalog) {
	state := make(map[string]int)
	var stack []string
	var visit func(string)
	visit = func(moduleID string) {
		switch state[moduleID] {
		case 1:
			start := 0
			for index, id := range stack {
				if id == moduleID {
					start = index
					break
				}
			}
			cycle := append(append([]string(nil), stack[start:]...), moduleID)
			l.add("catalog.module.dependency.cycle", moduleID, strings.Join(cycle, " -> "))
			return
		case 2:
			return
		}
		state[moduleID] = 1
		stack = append(stack, moduleID)
		for _, dependency := range stringsOrEmpty(catalog.modules[moduleID]["dependsOn"]) {
			if _, ok := catalog.modules[dependency]; ok {
				visit(dependency)
			}
		}
		stack = stack[:len(stack)-1]
		state[moduleID] = 2
	}
	for _, moduleID := range catalog.ModuleIDs() {
		visit(moduleID)
	}
}

func (l *catalogLoader) validateProfiles(catalog *Catalog) {
	allRules := make(map[string]string)
	for moduleID, module := range catalog.modules {
		for _, rule := range objectsOrEmpty(module["rules"]) {
			if id, ok := stringValue(rule, "id"); ok {
				allRules[id] = moduleID
			}
		}
	}

	for profileID, profile := range catalog.profiles {
		requireFields(l, "profile", profileID, profile, "id", "version", "setup", "entryDecisions", "modules")
		setupID, _ := stringValue(profile, "setup")
		setup, setupExists := catalog.setups[setupID]
		if !setupExists {
			l.add("catalog.profile.setup.unknown", profileID, setupID)
		}
		selected, ok := stringList(profile["modules"])
		if !ok {
			l.add("catalog.profile.modules.invalid", profileID, "")
			continue
		}
		if !uniqueStrings(selected) {
			l.add("catalog.profile.module.duplicate", profileID, "")
		}
		selectedSet := stringSet(selected)
		position := make(map[string]int, len(selected))
		for index, moduleID := range selected {
			position[moduleID] = index
		}
		for index, moduleID := range selected {
			module, exists := catalog.modules[moduleID]
			if !exists {
				l.add("catalog.profile.module.unknown", profileID, moduleID)
				continue
			}
			for _, dependency := range stringsOrEmpty(module["dependsOn"]) {
				if _, included := selectedSet[dependency]; !included {
					l.add("catalog.profile.dependency.missing", profileID, moduleID+" -> "+dependency)
				} else if position[dependency] >= index {
					l.add("catalog.profile.dependency.order", profileID, moduleID+" -> "+dependency)
				}
			}
			for _, conflict := range stringsOrEmpty(module["conflictsWith"]) {
				if _, included := selectedSet[conflict]; included {
					l.add("catalog.profile.module.conflict", profileID, moduleID+" -> "+conflict)
				}
			}
		}
		catalog.orderedModules[profileID] = append([]string(nil), selected...)

		entryDecisions, ok := stringList(profile["entryDecisions"])
		if !ok || !uniqueStrings(entryDecisions) {
			l.add("catalog.profile.entryDecisions.invalid", profileID, "")
		}
		for _, decisionID := range entryDecisions {
			if _, exists := catalog.decisions[decisionID]; !exists {
				l.add("catalog.profile.entryDecision.unknown", profileID, decisionID)
			}
		}
		if err := validateConditionalProfileDecisions(
			entryDecisions,
			profileCapabilityIDs(profile),
			catalog,
		); err != nil {
			l.add("catalog.profile.decision.capability.invalid", profileID, err.Error())
		}
		for _, ruleID := range stringsOrEmpty(profile["requiredRules"]) {
			moduleID, exists := allRules[ruleID]
			if !exists {
				l.add("catalog.profile.rule.unknown", profileID, ruleID)
			} else if _, included := selectedSet[moduleID]; !included {
				l.add("catalog.profile.rule.module.missing", profileID, ruleID)
			}
		}
		if setupExists {
			setupSkills := make(map[string]struct{})
			for _, skill := range objectsOrEmpty(setup["skills"]) {
				if name, ok := stringValue(skill, "name"); ok {
					setupSkills[name] = struct{}{}
				}
			}
			for _, moduleID := range selected {
				for _, skill := range stringsOrEmpty(catalog.modules[moduleID]["requiredSkills"]) {
					if _, exists := setupSkills[skill]; !exists {
						l.add("catalog.profile.skill.outside-setup", profileID, skill)
					}
				}
			}
		}
		l.validateProfileCapabilities(profileID, profile)
	}
}

func (l *catalogLoader) validateProfileCapabilities(profileID string, profile document) {
	capabilities, present := profile["capabilities"]
	if !present {
		return
	}
	seen := make(map[string]string)
	for _, capability := range objectsOrDiagnostic(l, capabilities, "catalog.profile.capabilities.invalid", profileID) {
		capabilityID := validateOwnedID(l, "capability", profileID, capability, seen)
		for _, field := range []string{"title", "category", "strength"} {
			if _, ok := stringValue(capability, field); !ok {
				l.add("catalog.capability.field.invalid", capabilityID, field)
			}
		}
		if _, ok := objectValue(capability["probe"]); !ok {
			l.add("catalog.capability.probe.invalid", capabilityID, "")
		}
	}
}

func (l *catalogLoader) validateProfileFormatters(catalog *Catalog) {
	for profileID, profile := range catalog.profiles {
		schema, _ := stringValue(profile, "schemaVersion")
		if schema != "setup-context-driven/profile/0.0.1" {
			continue
		}
		formatter, ok := objectValue(profile["formatter"])
		if !ok {
			l.add("catalog.profile.formatter.invalid", profileID, "formatter must be an object")
			continue
		}
		kind, kindOK := stringValue(formatter, "kind")
		switch kind {
		case "none":
			l.allowFields(
				"catalog.profile.formatter.field.unknown",
				profileID,
				formatter,
				"kind",
			)
		case "selected":
			l.allowFields(
				"catalog.profile.formatter.field.unknown",
				profileID,
				formatter,
				"kind",
				"id",
				"version",
				"fixturePaths",
				"goldenDigest",
			)
			formatterID, formatterIDOK := stringValue(formatter, "id")
			version, versionOK := stringValue(formatter, "version")
			fixturePaths, fixturePathsOK := stringList(formatter["fixturePaths"])
			goldenDigest, digestOK := stringValue(formatter, "goldenDigest")
			if !formatterIDOK || !versionOK || strings.TrimSpace(formatterID) == "" ||
				strings.TrimSpace(version) == "" {
				l.add("catalog.profile.formatter.identity.invalid", profileID, "")
			}
			if !fixturePathsOK || len(fixturePaths) == 0 || !uniqueStrings(fixturePaths) {
				l.add("catalog.profile.formatter.fixtures.invalid", profileID, "")
				continue
			}
			fixtures := make(map[string][]byte, len(fixturePaths))
			prefix := path.Join("formatter-fixtures", profileID, "golden") + "/"
			for _, fixturePath := range fixturePaths {
				if !safeRelative(fixturePath) ||
					!strings.HasPrefix(fixturePath, prefix) ||
					path.Ext(fixturePath) != ".md" {
					l.add(
						"catalog.profile.formatter.fixture.path.invalid",
						profileID,
						fixturePath,
					)
					continue
				}
				content, exists := catalog.assets[fixturePath]
				if !exists {
					l.add(
						"catalog.profile.formatter.fixture.missing",
						profileID,
						fixturePath,
					)
					continue
				}
				fixtures[strings.TrimPrefix(fixturePath, prefix)] = content
			}
			if !digestOK || !lowerSHA256.MatchString(goldenDigest) {
				l.add("catalog.profile.formatter.goldenDigest.invalid", profileID, "")
				continue
			}
			if len(fixtures) != len(fixturePaths) {
				continue
			}
			if observed := portableFileDigest(fixtures); observed != goldenDigest {
				l.add(
					"catalog.profile.formatter.goldenDigest.mismatch",
					profileID,
					fmt.Sprintf("got %s, want %s", goldenDigest, observed),
				)
			}
		default:
			if !kindOK {
				kind = ""
			}
			l.add("catalog.profile.formatter.kind.invalid", profileID, kind)
		}
	}
}

func portableFileDigest(files map[string][]byte) string {
	digest := sha256.New()
	for _, relative := range sortedKeys(files) {
		writeLengthPrefixed(digest, relative)
		var length [8]byte
		binary.BigEndian.PutUint64(length[:], uint64(len(files[relative])))
		_, _ = digest.Write(length[:])
		_, _ = digest.Write(files[relative])
	}
	return hex.EncodeToString(digest.Sum(nil))
}

func (l *catalogLoader) validateGuidanceComposition(catalog *Catalog) {
	core, ok := catalog.modules["core"]
	if !ok {
		return
	}
	schema, _ := stringValue(core, "schemaVersion")
	if schema != "setup-context-driven/module-v3" {
		return
	}
	for moduleID, module := range catalog.modules {
		if moduleID == "core" {
			continue
		}
		for _, field := range []string{"instructionHierarchy", "semanticOwners"} {
			if _, exists := module[field]; exists {
				l.add("catalog.guidance.owner.invalid", moduleID, field)
			}
		}
	}

	rootModules := make(map[string]string)
	moduleRoots := make(map[string][]string)
	guides := make(map[string]SemanticOwner)
	guidePaths := make(map[string][]string)
	for moduleID, module := range catalog.modules {
		for _, root := range objectsOrEmpty(module["rootBlocks"]) {
			rootID, _ := stringValue(root, "id")
			rootModules[rootID] = moduleID
			moduleRoots[moduleID] = append(moduleRoots[moduleID], rootID)
		}
		for _, guide := range objectsOrEmpty(module["supportingGuides"]) {
			managedID, _ := stringValue(guide, "id")
			targetPath, _ := stringValue(guide, "path")
			guides[managedID] = SemanticOwner{
				ManagedID: managedID,
				Path:      targetPath,
				Module:    moduleID,
			}
			guidePaths[targetPath] = append(guidePaths[targetPath], managedID)
		}
	}

	hierarchy, hierarchyOK := objectValue(core["instructionHierarchy"])
	if !hierarchyOK {
		l.add("catalog.instruction-hierarchy.invalid", "core", "")
	} else {
		l.allowFields(
			"catalog.instruction-hierarchy.field.unknown",
			"core",
			hierarchy,
			"narrowerPolicy",
			"levels",
		)
		policy, _ := stringValue(hierarchy, "narrowerPolicy")
		if policy != "strengthen-only" {
			l.add("catalog.instruction-hierarchy.policy.invalid", "core", policy)
		}
		levels := objectsOrDiagnostic(
			l,
			hierarchy["levels"],
			"catalog.instruction-hierarchy.levels.invalid",
			"core",
		)
		if len(levels) != len(confirmedInstructionHierarchy) {
			l.add(
				"catalog.instruction-hierarchy.order.invalid",
				"core",
				fmt.Sprintf("%d levels", len(levels)),
			)
		}
		seenRoots := make(map[string]string)
		rootPosition := make(map[string]int)
		validLevels := make([]InstructionHierarchyLevel, 0, len(levels))
		position := 0
		for index, level := range levels {
			levelPath := fmt.Sprintf("core.levels[%d]", index)
			l.allowFields(
				"catalog.instruction-hierarchy.level.field.unknown",
				levelPath,
				level,
				"id",
				"title",
				"rootBlocks",
			)
			levelID, _ := stringValue(level, "id")
			if index >= len(confirmedInstructionHierarchy) ||
				levelID != confirmedInstructionHierarchy[index] {
				l.add("catalog.instruction-hierarchy.order.invalid", levelPath, levelID)
			}
			title, titleOK := stringValue(level, "title")
			if !titleOK {
				l.add("catalog.instruction-hierarchy.title.invalid", levelPath, "")
			}
			rootBlocks, rootsOK := stringList(level["rootBlocks"])
			if !rootsOK || len(rootBlocks) == 0 || !uniqueStrings(rootBlocks) {
				l.add("catalog.instruction-hierarchy.roots.invalid", levelPath, "")
			}
			for _, rootID := range rootBlocks {
				_, exists := rootModules[rootID]
				if !exists {
					l.add("catalog.instruction-hierarchy.root.unknown", levelPath, rootID)
					continue
				}
				if previous, duplicate := seenRoots[rootID]; duplicate {
					l.add(
						"catalog.instruction-hierarchy.root.duplicate",
						rootID,
						previous+", "+levelID,
					)
				}
				seenRoots[rootID] = levelID
				rootPosition[rootID] = position
				position++
			}
			validLevels = append(validLevels, InstructionHierarchyLevel{
				ID:         levelID,
				Title:      title,
				RootBlocks: append([]string(nil), rootBlocks...),
			})
		}
		for rootID := range rootModules {
			if _, exists := seenRoots[rootID]; !exists {
				l.add("catalog.instruction-hierarchy.root.missing", rootID, "")
			}
		}
		for rootID, moduleID := range rootModules {
			current, currentExists := rootPosition[rootID]
			if !currentExists {
				continue
			}
			for _, dependency := range stringsOrEmpty(catalog.modules[moduleID]["dependsOn"]) {
				for _, dependencyRoot := range moduleRoots[dependency] {
					if dependencyPosition, exists := rootPosition[dependencyRoot]; exists &&
						dependencyPosition >= current {
						l.add(
							"catalog.instruction-hierarchy.dependency.order",
							rootID,
							dependencyRoot,
						)
					}
				}
			}
		}
		catalog.instructionHierarchy = validLevels
	}

	owners := objectsOrDiagnostic(
		l,
		core["semanticOwners"],
		"catalog.semantic-owners.invalid",
		"core",
	)
	seenOwners := make(map[string]string)
	seenClassifications := make(map[string]string)
	for index, declaration := range owners {
		ownerPath := fmt.Sprintf("core.semanticOwners[%d]", index)
		l.allowFields(
			"catalog.semantic-owner.field.unknown",
			ownerPath,
			declaration,
			"managedId",
			"title",
			"classifications",
		)
		managedID, managedIDOK := stringValue(declaration, "managedId")
		if !managedIDOK {
			l.add("catalog.semantic-owner.managed-id.invalid", ownerPath, "")
			continue
		}
		if previous, duplicate := seenOwners[managedID]; duplicate {
			l.add(
				"catalog.semantic-owner.managed-id.duplicate",
				managedID,
				previous+", "+ownerPath,
			)
		}
		seenOwners[managedID] = ownerPath
		owner, exists := guides[managedID]
		if !exists {
			l.add("catalog.semantic-owner.destination.unknown", ownerPath, managedID)
			continue
		}
		title, titleOK := stringValue(declaration, "title")
		if !titleOK {
			l.add("catalog.semantic-owner.title.invalid", managedID, "")
		}
		classifications, classificationsOK := stringList(declaration["classifications"])
		if !classificationsOK || len(classifications) == 0 ||
			!uniqueStrings(classifications) {
			l.add("catalog.semantic-owner.classifications.invalid", managedID, "")
		}
		for _, classification := range classifications {
			if previous, duplicate := seenClassifications[classification]; duplicate {
				l.add(
					"catalog.semantic-owner.classification.duplicate",
					classification,
					previous+", "+managedID,
				)
			}
			seenClassifications[classification] = managedID
		}
		owner.Title = title
		owner.Classifications = append([]string(nil), classifications...)
		catalog.semanticOwners[managedID] = owner
	}
	for managedID := range guides {
		if _, exists := seenOwners[managedID]; !exists {
			l.add("catalog.semantic-owner.destination.missing", managedID, "")
		}
	}

	l.validateRootGuidePointers(catalog, guides, guidePaths)
	l.validateClauseWeakening(catalog)
}

func (l *catalogLoader) validateRootGuidePointers(
	catalog *Catalog,
	guides map[string]SemanticOwner,
	guidePaths map[string][]string,
) {
	pointers := make(map[string]int)
	for _, module := range catalog.modules {
		for _, root := range objectsOrEmpty(module["rootBlocks"]) {
			rootID, _ := stringValue(root, "id")
			templateID, _ := stringValue(root, "template")
			template, templateExists := catalog.templates[templateID]
			if !templateExists {
				continue
			}
			templatePath, _ := stringValue(template, "path")
			templateBytes, templateBytesExist := catalog.assets[path.Join("templates", templatePath)]
			if !templateBytesExist {
				continue
			}
			for _, reference := range objectsOrEmpty(root["references"]) {
				token, _ := stringValue(reference, "token")
				if occurrences := bytes.Count(
					templateBytes,
					[]byte("{{"+token+"}}"),
				); occurrences != 1 {
					l.add(
						"catalog.instruction-hierarchy.pointer.occurrence",
						rootID,
						fmt.Sprintf("%s=%d", token, occurrences),
					)
				}
				managedID, managed := stringValue(reference, "managedId")
				if !managed {
					continue
				}
				owner, exists := guides[managedID]
				if !exists {
					l.add(
						"catalog.instruction-hierarchy.pointer.destination.unknown",
						rootID,
						managedID,
					)
					continue
				}
				pointers[owner.Path]++
			}
		}
	}
	for targetPath, managedIDs := range guidePaths {
		switch pointers[targetPath] {
		case 1:
		case 0:
			l.add(
				"catalog.instruction-hierarchy.pointer.missing",
				targetPath,
				strings.Join(managedIDs, ", "),
			)
		default:
			l.add(
				"catalog.instruction-hierarchy.pointer.duplicate",
				targetPath,
				fmt.Sprintf("%d pointers", pointers[targetPath]),
			)
		}
	}
}

func (l *catalogLoader) validateClauseWeakening(catalog *Catalog) {
	protected := stringSet(catalog.DecisionIDs())
	for _, rule := range objectsOrEmpty(catalog.modules["core"]["rules"]) {
		for _, clause := range objectsOrEmpty(rule["clauses"]) {
			if clauseID, ok := stringValue(clause, "id"); ok {
				protected[clauseID] = struct{}{}
			}
		}
	}
	for moduleID, module := range catalog.modules {
		if moduleID == "core" {
			continue
		}
		for _, rule := range objectsOrEmpty(module["rules"]) {
			for _, clause := range objectsOrEmpty(rule["clauses"]) {
				if _, declared := clause["weakens"]; !declared {
					continue
				}
				clauseID, _ := stringValue(clause, "id")
				targets, ok := stringList(clause["weakens"])
				if !ok || len(targets) == 0 {
					l.add("catalog.clause.weakening.invalid", clauseID, "")
					continue
				}
				for _, target := range targets {
					if _, protectedTarget := protected[target]; protectedTarget {
						l.add("catalog.clause.weakening.prohibited", clauseID, target)
					} else {
						l.add("catalog.clause.weakening.target.unknown", clauseID, target)
					}
				}
			}
		}
	}
}

func (l *catalogLoader) validateDecisionEffects(catalog *Catalog) {
	artifacts := catalogArtifacts(catalog)
	graph := make(map[string][]string)
	allowedFields := stringSet([]string{
		"when",
		"activateModules",
		"requireDecisions",
		"includeArtifacts",
		"excludeArtifacts",
		"selectTemplates",
		"renderBindings",
	})
	for decisionID, decision := range catalog.decisions {
		decisionType, _ := stringValue(decision, "type")
		if !containsString(
			[]string{
				"auth-provider",
				"boolean",
				"enum",
				"http-contract",
				"identifier-strategy",
				"string",
			},
			decisionType,
		) {
			l.add("catalog.decision.type.invalid", decisionID, decisionType)
		}
		if value, ok := decision["default"]; ok {
			if err := validateDecisionValue(decision, value); err != nil {
				l.add("catalog.decision.default.invalid", decisionID, err.Error())
			}
		}
		if value, ok := decision["suggestion"]; ok {
			if err := validateDecisionValue(decision, value); err != nil {
				l.add("catalog.decision.suggestion.invalid", decisionID, err.Error())
			}
		}
		requiredCapabilities, capabilitiesOK := stringList(decision["requiresCapabilities"])
		if _, declared := decision["requiresCapabilities"]; declared {
			if !capabilitiesOK ||
				len(requiredCapabilities) == 0 ||
				!uniqueStrings(requiredCapabilities) {
				l.add("catalog.decision.capabilities.invalid", decisionID, "")
			}
			knownCapabilities := catalogCapabilityIDs(catalog)
			for _, capabilityID := range requiredCapabilities {
				if _, exists := knownCapabilities[capabilityID]; !exists {
					l.add("catalog.decision.capability.unknown", decisionID, capabilityID)
				}
			}
		}
		effects, ok := objectList(decision["effects"])
		if !ok || len(effects) == 0 {
			l.add("catalog.decision.effects.invalid", decisionID, "")
			continue
		}
		for index, effect := range effects {
			for field := range effect {
				if _, ok := allowedFields[field]; !ok {
					l.add("catalog.decision.effect.field.unknown", decisionID, field)
				}
			}
			if _, ok := objectValue(effect["when"]); !ok {
				l.add("catalog.decision.effect.condition.invalid", decisionID, fmt.Sprintf("index %d", index))
			}
			for _, moduleID := range stringsOrEmpty(effect["activateModules"]) {
				if _, ok := catalog.modules[moduleID]; !ok {
					l.add("catalog.decision.effect.module.unknown", decisionID, moduleID)
				}
			}
			for _, requiredID := range stringsOrEmpty(effect["requireDecisions"]) {
				if _, ok := catalog.decisions[requiredID]; !ok {
					l.add("catalog.decision.effect.decision.unknown", decisionID, requiredID)
				} else {
					graph[decisionID] = append(graph[decisionID], requiredID)
				}
			}
			for _, field := range []string{"includeArtifacts", "excludeArtifacts"} {
				for _, artifactID := range stringsOrEmpty(effect[field]) {
					if _, ok := artifacts[artifactID]; !ok {
						l.add("catalog.decision.effect.artifact.unknown", decisionID, artifactID)
					}
				}
			}
			for _, selection := range objectsOrEmpty(effect["selectTemplates"]) {
				artifactID, _ := stringValue(selection, "artifact")
				templateID, _ := stringValue(selection, "template")
				if _, ok := artifacts[artifactID]; !ok {
					l.add("catalog.decision.effect.artifact.unknown", decisionID, artifactID)
				}
				if _, ok := catalog.templates[templateID]; !ok {
					l.add("catalog.decision.effect.template.unknown", decisionID, templateID)
				}
			}
			for _, binding := range objectsOrEmpty(effect["renderBindings"]) {
				artifactID, _ := stringValue(binding, "artifact")
				templateID, _ := stringValue(binding, "template")
				token, _ := stringValue(binding, "token")
				if _, ok := artifacts[artifactID]; !ok {
					l.add("catalog.decision.effect.artifact.unknown", decisionID, artifactID)
				}
				template, ok := catalog.templates[templateID]
				if !ok {
					l.add("catalog.decision.effect.template.unknown", decisionID, templateID)
				} else if !containsString(stringsOrEmpty(template["tokens"]), token) {
					l.add("catalog.decision.effect.token.unknown", decisionID, token)
				}
			}
		}
	}
	l.validateDecisionCycles(graph)
}

func (l *catalogLoader) validateDecisionCycles(graph map[string][]string) {
	state := make(map[string]int)
	var visit func(string)
	visit = func(id string) {
		if state[id] == 1 {
			l.add("catalog.decision.dependency.cycle", id, "")
			return
		}
		if state[id] == 2 {
			return
		}
		state[id] = 1
		for _, dependency := range graph[id] {
			visit(dependency)
		}
		state[id] = 2
	}
	for id := range graph {
		visit(id)
	}
}

func (l *catalogLoader) validateSetups(catalog *Catalog) {
	for setupID, setup := range catalog.setups {
		requireFields(l, "setup", setupID, setup, "id", "version", "source", "digest", "skills")
		skills, ok := objectList(setup["skills"])
		if !ok {
			l.add("catalog.setup.skills.invalid", setupID, "")
			continue
		}
		seenNames := make(map[string]string)
		normalized := make([]document, 0, len(skills))
		for _, skill := range skills {
			name := validateOwnedID(l, "setup.skill", setupID, document{"id": skill["name"]}, seenNames)
			targetPath, ok := stringValue(skill, "path")
			if !ok || !safeRelative(targetPath) {
				l.add("catalog.setup.skill.path.invalid", setupID, name)
			}
			source, ok := objectValue(skill["source"])
			if !ok {
				l.add("catalog.setup.skill.source.invalid", setupID, name)
				continue
			}
			sourceType, _ := stringValue(source, "type")
			switch sourceType {
			case "github":
				ref, _ := stringValue(source, "ref")
				sourcePath, sourcePathOK := stringValue(source, "path")
				repository, repositoryOK := stringValue(source, "repository")
				treeDigest, digestOK := stringValue(skill, "treeDigest")
				if !repositoryOK || !strings.Contains(repository, "/") {
					l.add("catalog.setup.skill.repository.invalid", setupID, name)
				}
				if !immutableGitID.MatchString(ref) {
					l.add("catalog.setup.skill.ref.mutable", setupID, name)
				}
				if !sourcePathOK || !safeRelative(sourcePath) {
					l.add("catalog.setup.skill.source.path.invalid", setupID, name)
				}
				if !digestOK || !lowerSHA256.MatchString(treeDigest) {
					l.add("catalog.setup.skill.treeDigest.invalid", setupID, name)
				}
			case "repo":
				digest, ok := stringValue(skill, "contentDigest")
				if !ok || !lowerSHA256.MatchString(digest) {
					l.add("catalog.setup.skill.contentDigest.invalid", setupID, name)
				}
			default:
				l.add("catalog.setup.skill.source.type.unknown", setupID, name)
			}
			normalized = append(normalized, skill)
		}

		var payload any = normalized
		if bundles, exists := setup["activationBundles"]; exists {
			payload = document{"activationBundles": bundles, "skills": toAnySlice(normalized)}
		}
		expected, err := canonicalSHA256(payload)
		if err != nil {
			l.add("catalog.setup.digest.invalid", setupID, err.Error())
			continue
		}
		declared, _ := stringValue(setup, "digest")
		if declared != expected {
			l.add("catalog.setup.digest.mismatch", setupID, fmt.Sprintf("got %s, want %s", declared, expected))
		}
	}
}

func (l *catalogLoader) validateSkillActivations(catalog *Catalog) {
	doc := l.documents["skill-activations.json"]
	if doc == nil {
		return
	}
	bundles := objectsOrDiagnostic(l, doc["bundles"], "catalog.skill.bundle.collection.invalid", "skill-activations.json")
	seenBundles := make(map[string]string)
	bundleSkills := make(map[string][]string)
	for _, bundle := range bundles {
		id := validateOwnedID(l, "skill.bundle", "skill-activations.json", bundle, seenBundles)
		skills, ok := stringList(bundle["skills"])
		if !ok || len(skills) == 0 || !uniqueStrings(skills) {
			l.add("catalog.skill.bundle.skills.invalid", id, "")
		}
		bundleSkills[id] = skills
	}
	seenTriggers := make(map[string]string)
	for _, activation := range objectsOrDiagnostic(l, doc["activations"], "catalog.skill.activation.collection.invalid", "skill-activations.json") {
		id := validateOwnedID(l, "skill.activation", "skill-activations.json", activation, seenTriggers)
		owner, _ := stringValue(activation, "owner")
		if _, ok := catalog.modules[owner]; !ok {
			l.add("catalog.skill.activation.owner.unknown", id, owner)
		}
		bundle, _ := stringValue(activation, "bundle")
		if _, ok := bundleSkills[bundle]; !ok {
			l.add("catalog.skill.activation.bundle.unknown", id, bundle)
		}
	}
}

func (l *catalogLoader) validateTransitions(catalog *Catalog) {
	clauses := make(map[string]string)
	extensions := make(map[string]struct{})
	for _, module := range catalog.modules {
		for _, rule := range objectsOrEmpty(module["rules"]) {
			for _, clause := range objectsOrEmpty(rule["clauses"]) {
				if id, ok := stringValue(clause, "id"); ok {
					enforcement, _ := stringValue(clause, "enforcement")
					clauses[id] = enforcement
				}
			}
		}
		for _, extension := range objectsOrEmpty(module["repositoryExtensions"]) {
			if id, ok := stringValue(extension, "id"); ok {
				extensions[id] = struct{}{}
			}
		}
	}
	for transitionID, transition := range catalog.transitions {
		requireFields(
			l,
			"transition",
			transitionID,
			transition,
			"id",
			"version",
			"fromBaseline",
			"toBaseline",
			"legacyManifestFingerprints",
			"priorClauses",
			"mappings",
		)
		prior := make(map[string]string)
		for _, clause := range objectsOrDiagnostic(l, transition["priorClauses"], "catalog.transition.priorClauses.invalid", transitionID) {
			id, _ := stringValue(clause, "id")
			if _, exists := prior[id]; exists {
				l.add("catalog.transition.priorClause.duplicate", transitionID, id)
			}
			enforcement, _ := stringValue(clause, "enforcement")
			prior[id] = enforcement
			digest, _ := stringValue(clause, "guidanceDigest")
			if !lowerSHA256.MatchString(digest) {
				l.add("catalog.transition.guidanceDigest.invalid", transitionID, id)
			}
		}
		mapped := make(map[string]struct{})
		for _, mapping := range objectsOrDiagnostic(l, transition["mappings"], "catalog.transition.mappings.invalid", transitionID) {
			from, _ := stringValue(mapping, "fromClause")
			if _, exists := mapped[from]; exists {
				l.add("catalog.transition.mapping.duplicate", transitionID, from)
			}
			mapped[from] = struct{}{}
			if _, exists := prior[from]; !exists {
				l.add("catalog.transition.mapping.source.unknown", transitionID, from)
			}
			disposition, _ := stringValue(mapping, "disposition")
			targets := stringsOrEmpty(mapping["targets"])
			if disposition == "rejected" && len(targets) != 0 {
				l.add("catalog.transition.rejected.targets.invalid", transitionID, from)
			}
			if disposition != "rejected" && len(targets) == 0 {
				l.add("catalog.transition.accepted.targets.missing", transitionID, from)
			}
			for _, target := range targets {
				enforcement, clauseExists := clauses[target]
				_, extensionExists := extensions[target]
				if !clauseExists && !extensionExists {
					l.add("catalog.transition.target.unknown", transitionID, target)
				} else if clauseExists && enforcement != prior[from] {
					l.add("catalog.transition.target.enforcement.mismatch", transitionID, target)
				}
			}
		}
		for priorID := range prior {
			if _, ok := mapped[priorID]; !ok {
				l.add("catalog.transition.mapping.incomplete", transitionID, priorID)
			}
		}
	}
}

func (l *catalogLoader) validateSourceBaselines(catalog *Catalog) {
	index := l.documents["source-baselines/index.json"]
	if index == nil {
		return
	}
	seen := make(map[string]string)
	for _, baseline := range objectsOrDiagnostic(l, index["baselines"], "catalog.sourceBaseline.collection.invalid", "source-baselines/index.json") {
		baselineID := validateOwnedID(l, "sourceBaseline", "source-baselines/index.json", baseline, seen)
		profileID, _ := stringValue(baseline, "profile")
		if _, ok := catalog.profiles[profileID]; !ok {
			l.add("catalog.sourceBaseline.profile.unknown", baselineID, profileID)
		}
		baselinePath, ok := stringValue(baseline, "path")
		if !ok || !safeRelative(baselinePath) {
			l.add("catalog.sourceBaseline.path.invalid", baselineID, baselinePath)
			continue
		}
		docPath := path.Join("source-baselines", baselinePath, "baseline.json")
		doc := l.documents[docPath]
		if doc == nil {
			l.add("catalog.sourceBaseline.document.missing", baselineID, docPath)
			continue
		}
		if got, _ := stringValue(doc, "id"); got != baselineID {
			l.add("catalog.sourceBaseline.id.mismatch", baselineID, got)
		}
		if got, _ := stringValue(doc, "profile"); got != profileID {
			l.add("catalog.sourceBaseline.profile.mismatch", baselineID, got)
		}
		manifest, _ := stringValue(doc, "manifest")
		if _, ok := catalog.assets[path.Join("source-baselines", baselinePath, manifest)]; !ok {
			l.add("catalog.sourceBaseline.manifest.missing", baselineID, manifest)
		}
		corpus, _ := stringValue(doc, "corpus")
		prefix := path.Join("source-baselines", baselinePath, corpus) + "/"
		foundCorpus := false
		for assetPath := range catalog.assets {
			if strings.HasPrefix(assetPath, prefix) {
				foundCorpus = true
				break
			}
		}
		if !foundCorpus {
			l.add("catalog.sourceBaseline.corpus.missing", baselineID, corpus)
		}
		baseline, err := catalog.SourceBaseline(baselineID)
		if err != nil {
			l.add("catalog.sourceBaseline.integrity.invalid", baselineID, err.Error())
			continue
		}
		l.validateSourceBaselineCoverage(catalog, baseline)
	}
}

func (l *catalogLoader) validateSourceBaselineCoverage(
	catalog *Catalog,
	baseline SourceBaseline,
) {
	profile := catalog.profiles[baseline.Identity.Profile]
	selectedModules := stringSet(stringsOrEmpty(profile["modules"]))
	entryIDs := make(map[string]struct{}, len(baseline.Entries))
	carriers := make(map[string]struct{}, len(baseline.Entries))
	for _, entry := range baseline.Entries {
		entryIDs[entry.ID] = struct{}{}
		carriers[entry.Carrier] = struct{}{}
	}

	rules := make(map[string]document)
	for moduleID := range selectedModules {
		for _, rule := range objectsOrEmpty(catalog.modules[moduleID]["rules"]) {
			ruleID, _ := stringValue(rule, "id")
			rules[ruleID] = rule
		}
	}
	for _, ruleID := range stringsOrEmpty(profile["requiredRules"]) {
		rule := rules[ruleID]
		clauses := objectsOrEmpty(rule["clauses"])
		if len(clauses) == 0 {
			if _, ok := entryIDs[ruleID]; !ok {
				l.add(
					"catalog.sourceBaseline.required-rule.missing",
					baseline.Identity.ID,
					ruleID,
				)
			}
			continue
		}
		for _, clause := range clauses {
			clauseID, _ := stringValue(clause, "id")
			if _, ok := entryIDs[clauseID]; !ok {
				l.add(
					"catalog.sourceBaseline.required-clause.missing",
					baseline.Identity.ID,
					clauseID,
				)
			}
		}
	}

	for managedID, owner := range catalog.semanticOwners {
		if _, active := selectedModules[owner.Module]; !active ||
			owner.Module == repositoryExtensionModuleID {
			continue
		}
		if _, ok := carriers[owner.Path]; !ok {
			l.add(
				"catalog.sourceBaseline.semantic-destination.missing",
				baseline.Identity.ID,
				managedID,
			)
		}
	}
}

func (l *catalogLoader) indexSimpleIDs(assetPath, field, kind string) map[string]document {
	result := make(map[string]document)
	doc := l.documents[assetPath]
	items, ok := objectList(doc[field])
	if !ok {
		l.add("catalog."+kind+".collection.invalid", assetPath, field)
		return result
	}
	for _, item := range items {
		id, ok := stringValue(item, "id")
		if !ok {
			l.add("catalog."+kind+".id.missing", assetPath, "")
			continue
		}
		if _, exists := result[id]; exists {
			l.add("catalog."+kind+".id.duplicate", assetPath, id)
		}
		result[id] = item
	}
	return result
}

func catalogArtifacts(catalog *Catalog) map[string]struct{} {
	result := make(map[string]struct{})
	for _, module := range catalog.modules {
		for _, field := range []string{"rootBlocks", "supportingGuides"} {
			for _, artifact := range objectsOrEmpty(module[field]) {
				id, _ := stringValue(artifact, "id")
				result[id] = struct{}{}
			}
		}
	}
	return result
}

func requireFields(l *catalogLoader, kind, id string, doc document, fields ...string) {
	for _, field := range fields {
		if _, ok := doc[field]; !ok {
			l.add("catalog."+kind+".field.missing", id, field)
		}
	}
}

func validateOwnedID(
	l *catalogLoader,
	kind string,
	owner string,
	doc document,
	seen map[string]string,
) string {
	id, ok := stringValue(doc, "id")
	if !ok {
		l.add("catalog."+kind+".id.missing", owner, "")
		return ""
	}
	if previous, exists := seen[id]; exists {
		l.add("catalog."+kind+".id.duplicate", id, previous+", "+owner)
	}
	seen[id] = owner
	return id
}

func objectsOrDiagnostic(l *catalogLoader, value any, code, owner string) []document {
	objects, ok := objectList(value)
	if !ok {
		l.add(code, owner, "")
		return nil
	}
	return objects
}

func objectsOrEmpty(value any) []document {
	objects, _ := objectList(value)
	return objects
}

func stringsOrEmpty(value any) []string {
	values, _ := stringList(value)
	return values
}

func stringSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

func canonicalSHA256(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("serialize canonical digest input: %w", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func toAnySlice(values []document) []any {
	result := make([]any, len(values))
	for index := range values {
		result[index] = values[index]
	}
	return result
}
