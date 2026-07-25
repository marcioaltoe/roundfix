package baseline

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

const (
	catalogSchema = "roundfix/baseline-catalog/v1"
	digestDomain  = catalogSchema + "\x00"
)

//go:embed assets
var embeddedCatalogFS embed.FS

var embeddedAssets = mustSub(embeddedCatalogFS, "assets")

// Diagnostic identifies one closed-catalog validation failure.
type Diagnostic struct {
	Code string
	Path string
	Info string
}

func (d Diagnostic) String() string {
	message := d.Code
	if d.Path != "" {
		message += ": " + d.Path
	}
	if d.Info != "" {
		message += ": " + d.Info
	}
	return message
}

// ValidationError reports every deterministic validation failure found while
// loading a catalog.
type ValidationError struct {
	Diagnostics []Diagnostic
}

func (e *ValidationError) Error() string {
	if len(e.Diagnostics) == 0 {
		return "validate Baseline catalog"
	}
	messages := make([]string, len(e.Diagnostics))
	for index, diagnostic := range e.Diagnostics {
		messages[index] = diagnostic.String()
	}
	return "validate Baseline catalog: " + strings.Join(messages, "; ")
}

// Has reports whether the validation result contains code.
func (e *ValidationError) Has(code string) bool {
	for _, diagnostic := range e.Diagnostics {
		if diagnostic.Code == code {
			return true
		}
	}
	return false
}

// Asset is one immutable catalog file.
type Asset struct {
	Path string
	Data []byte
}

// Entry is one immutable, canonically serialized catalog declaration.
type Entry struct {
	ID            string
	SchemaVersion string
	Data          []byte
}

// InstructionHierarchyLevel is one precedence tier in the generated root
// instruction map.
type InstructionHierarchyLevel struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	RootBlocks []string `json:"rootBlocks"`
}

// SemanticOwner is one active managed guide that owns the named policy
// classifications.
type SemanticOwner struct {
	ManagedID       string   `json:"managedId"`
	Path            string   `json:"path"`
	Module          string   `json:"module"`
	Title           string   `json:"title"`
	Classifications []string `json:"classifications"`
}

// SemanticOwnerRegistry contains only semantic destinations active in one
// resolved Baseline Profile.
type SemanticOwnerRegistry map[string]SemanticOwner

// Catalog is the validated, deterministic Baseline catalog authority.
//
// Its methods return copies so callers cannot mutate catalog state.
type Catalog struct {
	assets               map[string][]byte
	profiles             map[string]document
	modules              map[string]document
	decisions            map[string]document
	templates            map[string]document
	setups               map[string]document
	transitions          map[string]document
	orderedModules       map[string][]string
	instructionHierarchy []InstructionHierarchyLevel
	semanticOwners       map[string]SemanticOwner
	normalized           []byte
	digest               string
}

// LoadEmbeddedCatalog loads the catalog compiled into the Roundfix binary.
func LoadEmbeddedCatalog() (*Catalog, error) {
	catalog, err := LoadCatalog(embeddedAssets)
	if err != nil {
		return nil, fmt.Errorf("load embedded Baseline catalog: %w", err)
	}
	return catalog, nil
}

// LoadCatalog validates a catalog filesystem rooted at the assets directory.
//
// This boundary is exported for repository-owned profile validation and
// maintainer asset checks. It performs no filesystem writes or network access.
func LoadCatalog(assetsFS fs.FS) (*Catalog, error) {
	loader := newCatalogLoader(assetsFS)
	catalog := loader.load()
	if len(loader.diagnostics) != 0 {
		sort.SliceStable(loader.diagnostics, func(i, j int) bool {
			left, right := loader.diagnostics[i], loader.diagnostics[j]
			if left.Code != right.Code {
				return left.Code < right.Code
			}
			if left.Path != right.Path {
				return left.Path < right.Path
			}
			return left.Info < right.Info
		})
		return nil, &ValidationError{Diagnostics: loader.diagnostics}
	}
	return catalog, nil
}

// ProfileIDs returns the built-in Baseline Profile IDs in lexical order.
func (c *Catalog) ProfileIDs() []string {
	return sortedKeys(c.profiles)
}

// ModuleIDs returns the embedded module IDs in lexical order.
func (c *Catalog) ModuleIDs() []string {
	return sortedKeys(c.modules)
}

// DecisionIDs returns the embedded decision IDs in lexical order.
func (c *Catalog) DecisionIDs() []string {
	return sortedKeys(c.decisions)
}

// TemplateIDs returns the embedded template IDs in lexical order.
func (c *Catalog) TemplateIDs() []string {
	return sortedKeys(c.templates)
}

// SetupIDs returns the embedded Setup Snapshot IDs in lexical order.
func (c *Catalog) SetupIDs() []string {
	return sortedKeys(c.setups)
}

// TransitionIDs returns the embedded Upgrade Retention transition IDs.
func (c *Catalog) TransitionIDs() []string {
	return sortedKeys(c.transitions)
}

// OrderedModules returns the profile's deterministic dependency order.
func (c *Catalog) OrderedModules(profileID string) ([]string, bool) {
	modules, ok := c.orderedModules[profileID]
	return append([]string(nil), modules...), ok
}

// InstructionHierarchy returns the catalog's confirmed root precedence order.
func (c *Catalog) InstructionHierarchy() []InstructionHierarchyLevel {
	levels := make([]InstructionHierarchyLevel, len(c.instructionHierarchy))
	for index, level := range c.instructionHierarchy {
		levels[index] = level
		levels[index].RootBlocks = append([]string(nil), level.RootBlocks...)
	}
	return levels
}

// SemanticOwnerRegistry derives semantic destinations from the intersection
// of active modules and active managed artifacts.
func (c *Catalog) SemanticOwnerRegistry(
	activeModules []string,
	activeArtifacts []string,
) SemanticOwnerRegistry {
	modules := stringSet(activeModules)
	artifacts := stringSet(activeArtifacts)
	registry := make(SemanticOwnerRegistry)
	for managedID, owner := range c.semanticOwners {
		if _, ok := modules[owner.Module]; !ok {
			continue
		}
		if _, ok := artifacts[managedID]; !ok {
			continue
		}
		owner.Classifications = append([]string(nil), owner.Classifications...)
		registry[managedID] = owner
	}
	return registry
}

// Asset returns an immutable copy of one catalog file.
func (c *Catalog) Asset(assetPath string) (Asset, bool) {
	data, ok := c.assets[assetPath]
	return Asset{Path: assetPath, Data: append([]byte(nil), data...)}, ok
}

// Profile returns one built-in Baseline Profile.
func (c *Catalog) Profile(id string) (Entry, bool) {
	return catalogEntry(id, c.profiles)
}

// Module returns one embedded Baseline module.
func (c *Catalog) Module(id string) (Entry, bool) {
	return catalogEntry(id, c.modules)
}

// Decision returns one embedded decision declaration.
func (c *Catalog) Decision(id string) (Entry, bool) {
	return catalogEntry(id, c.decisions)
}

// Template returns one embedded template declaration.
func (c *Catalog) Template(id string) (Entry, bool) {
	return catalogEntry(id, c.templates)
}

// TemplateContent returns the immutable bytes selected by a template.
func (c *Catalog) TemplateContent(id string) (Asset, bool) {
	template, ok := c.templates[id]
	if !ok {
		return Asset{}, false
	}
	templatePath, ok := stringValue(template, "path")
	if !ok {
		return Asset{}, false
	}
	return c.Asset(path.Join("templates", templatePath))
}

// Setup returns one embedded Setup Snapshot.
func (c *Catalog) Setup(id string) (Entry, bool) {
	return catalogEntry(id, c.setups)
}

// RetentionTransition returns one embedded Upgrade Retention transition.
func (c *Catalog) RetentionTransition(id string) (Entry, bool) {
	return catalogEntry(id, c.transitions)
}

// Normalized returns the canonical catalog identity document.
func (c *Catalog) Normalized() []byte {
	return append([]byte(nil), c.normalized...)
}

// Digest returns the domain-separated SHA-256 catalog identity.
func (c *Catalog) Digest() string {
	return c.digest
}

type normalizedCatalog struct {
	SchemaVersion string         `json:"schemaVersion"`
	Files         []fileIdentity `json:"files"`
	Profiles      []string       `json:"profiles"`
	Modules       []string       `json:"modules"`
	Decisions     []string       `json:"decisions"`
	Templates     []string       `json:"templates"`
	Setups        []string       `json:"setups"`
	Transitions   []string       `json:"retentionTransitions"`
}

type fileIdentity struct {
	Path   string `json:"path"`
	Bytes  int    `json:"bytes"`
	Digest string `json:"digest"`
}

func (c *Catalog) finishIdentity() error {
	paths := sortedKeys(c.assets)
	files := make([]fileIdentity, 0, len(paths))
	for _, assetPath := range paths {
		sum := sha256.Sum256(c.assets[assetPath])
		files = append(files, fileIdentity{
			Path:   assetPath,
			Bytes:  len(c.assets[assetPath]),
			Digest: "sha256:" + hex.EncodeToString(sum[:]),
		})
	}
	normalized := normalizedCatalog{
		SchemaVersion: catalogSchema,
		Files:         files,
		Profiles:      c.ProfileIDs(),
		Modules:       c.ModuleIDs(),
		Decisions:     c.DecisionIDs(),
		Templates:     c.TemplateIDs(),
		Setups:        c.SetupIDs(),
		Transitions:   c.TransitionIDs(),
	}
	data, err := json.Marshal(normalized)
	if err != nil {
		return fmt.Errorf("serialize normalized Baseline catalog: %w", err)
	}
	data = append(data, '\n')
	sum := sha256.Sum256(append([]byte(digestDomain), data...))
	c.normalized = data
	c.digest = "sha256:" + hex.EncodeToString(sum[:])
	return nil
}

func mustSub(root fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(root, dir)
	if err != nil {
		panic(fmt.Sprintf("open embedded Baseline catalog %q: %v", dir, err))
	}
	return sub
}

func sortedKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func catalogEntry(id string, values map[string]document) (Entry, bool) {
	doc, ok := values[id]
	if !ok {
		return Entry{}, false
	}
	data, err := json.Marshal(doc)
	if err != nil {
		return Entry{}, false
	}
	schema, _ := stringValue(doc, "schemaVersion")
	return Entry{ID: id, SchemaVersion: schema, Data: data}, true
}

func safeRelative(assetPath string) bool {
	return assetPath != "" &&
		assetPath != "." &&
		path.Clean(assetPath) == assetPath &&
		!path.IsAbs(assetPath) &&
		!strings.ContainsAny(assetPath, "\\\x00") &&
		!isDrivePath(assetPath) &&
		!containsParent(assetPath)
}

func containsParent(assetPath string) bool {
	for _, part := range strings.Split(assetPath, "/") {
		if part == ".." {
			return true
		}
	}
	return false
}

func isDrivePath(assetPath string) bool {
	return len(assetPath) >= 2 &&
		((assetPath[0] >= 'A' && assetPath[0] <= 'Z') ||
			(assetPath[0] >= 'a' && assetPath[0] <= 'z')) &&
		assetPath[1] == ':'
}
