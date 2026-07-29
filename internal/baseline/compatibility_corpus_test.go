package baseline

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"roundfix/skills"
)

// Suite: Baseline compatibility corpus
// Invariant: every maintained pre-cutover contract remains classified, content-addressed, and assigned to Go.
// Boundary IN: the frozen parity matrix, fixtures, blobs, and artifact manifest.
// Boundary OUT: individual Go destination behavior, owned by the package tests named by each matrix destination.

//go:embed testdata/parity-corpus/v1
var baselineCompatibilityCorpus embed.FS

const baselineCompatibilityRoot = "testdata/parity-corpus/v1"

type baselineCompatibilityManifest struct {
	SchemaVersion string `json:"schemaVersion"`
	FrozenDate    string `json:"frozenDate"`
	SourceSuite   struct {
		Root            string `json:"root"`
		TestCount       int    `json:"testCount"`
		InventoryDigest string `json:"inventoryDigest"`
	} `json:"sourceSuite"`
	Profiles        []string `json:"profiles"`
	AdoptionStates  []string `json:"adoptionStates"`
	FixtureIDs      []string `json:"fixtureIds"`
	RetiredBehavior []string `json:"retiredBehavior"`
	Artifacts       []struct {
		Path   string `json:"path"`
		SHA256 string `json:"sha256"`
		Bytes  int    `json:"bytes"`
	} `json:"artifacts"`
}

type baselineCompatibilityMatrix struct {
	SchemaVersion        string         `json:"schemaVersion"`
	SourceTestCount      int            `json:"sourceTestCount"`
	RowCount             int            `json:"rowCount"`
	Classifications      []string       `json:"classifications"`
	ClassificationCounts map[string]int `json:"classificationCounts"`
	Rows                 []struct {
		Behavior           string   `json:"behavior"`
		Classification     string   `json:"classification"`
		ContractDimensions []string `json:"contractDimensions"`
		FixtureIDs         []string `json:"fixtureIds"`
		GoDestination      string   `json:"goDestination"`
		ID                 string   `json:"id"`
		Rationale          *string  `json:"rationale"`
		SourceKind         string   `json:"sourceKind"`
		SourcePath         string   `json:"sourcePath"`
		SourceSuite        *string  `json:"sourceSuite"`
		SourceTest         *string  `json:"sourceTest"`
	} `json:"rows"`
}

type baselineCompatibilityBlobs struct {
	SchemaVersion string            `json:"schemaVersion"`
	Encoding      string            `json:"encoding"`
	Blobs         map[string]string `json:"blobs"`
}

type baselineCompatibilityFixture struct {
	SchemaVersion       string          `json:"schemaVersion"`
	ID                  string          `json:"id"`
	Classification      string          `json:"classification"`
	Profile             *string         `json:"profile"`
	AdoptionState       string          `json:"adoptionState"`
	ContractAreas       []string        `json:"contractAreas"`
	SourceTests         []string        `json:"sourceTests"`
	PlanDigest          *string         `json:"planDigest"`
	FileIdentities      json.RawMessage `json:"fileIdentities"`
	Input               json.RawMessage `json:"input"`
	ManagedEntryLedger  json.RawMessage `json:"managedEntryLedger"`
	Manifest            json.RawMessage `json:"manifest"`
	NormalizedOutput    json.RawMessage `json:"normalizedOutput"`
	PlannedByteSequence json.RawMessage `json:"plannedByteSequence"`
	PostState           json.RawMessage `json:"postState"`
	Rationale           json.RawMessage `json:"rationale"`
	Refusal             json.RawMessage `json:"refusal"`
	RepositoryPreimage  json.RawMessage `json:"repositoryPreimage"`
	Rollback            json.RawMessage `json:"rollback"`
}

func TestBaselineCompatibilityCorpus(t *testing.T) {
	if *updateBaselineDigests {
		regenerateBaselineCompatibilityCorpus(t)
		return
	}
	manifest := readBaselineCompatibilityJSON[baselineCompatibilityManifest](t, "manifest.json")
	if manifest.SchemaVersion != "setup-context-driven/parity-corpus/v1" {
		t.Fatalf("manifest schema = %q", manifest.SchemaVersion)
	}
	if manifest.FrozenDate != "2026-07-23" {
		t.Fatalf("manifest frozen date = %q", manifest.FrozenDate)
	}
	if manifest.SourceSuite.Root != "tests" || manifest.SourceSuite.TestCount != 240 {
		t.Fatalf("source suite = %#v", manifest.SourceSuite)
	}
	assertBaselineCompatibilityDigest(t, manifest.SourceSuite.InventoryDigest)
	if len(manifest.RetiredBehavior) != 0 {
		t.Fatalf("retired behavior must remain empty, got %v", manifest.RetiredBehavior)
	}

	assertBaselineCompatibilityArtifacts(t, manifest)

	matrix := readBaselineCompatibilityJSON[baselineCompatibilityMatrix](t, "matrix.json")
	if matrix.SchemaVersion != "setup-context-driven/parity-matrix/v1" {
		t.Fatalf("matrix schema = %q", matrix.SchemaVersion)
	}
	if matrix.SourceTestCount != manifest.SourceSuite.TestCount ||
		matrix.RowCount != len(matrix.Rows) ||
		matrix.RowCount != 243 {
		t.Fatalf(
			"matrix counts source=%d rows=%d actual=%d",
			matrix.SourceTestCount,
			matrix.RowCount,
			len(matrix.Rows),
		)
	}

	wantClassifications := map[string]int{
		"exact":          163,
		"semantic":       63,
		"designed-delta": 16,
		"ancillary":      1,
		"retired":        0,
	}
	if fmt.Sprint(matrix.ClassificationCounts) != fmt.Sprint(wantClassifications) {
		t.Fatalf("classification counts = %v, want %v", matrix.ClassificationCounts, wantClassifications)
	}
	if got := sortedBaselineCompatibilityStrings(matrix.Classifications); fmt.Sprint(got) !=
		fmt.Sprint(sortedBaselineCompatibilityKeys(wantClassifications)) {
		t.Fatalf("classifications = %v", matrix.Classifications)
	}

	fixtureIDs := baselineCompatibilitySet(t, "fixture id", manifest.FixtureIDs)
	rowIDs := make(map[string]struct{}, len(matrix.Rows))
	destinations := make(map[string]struct{})
	actualClassifications := make(map[string]int, len(wantClassifications))
	for classification := range wantClassifications {
		actualClassifications[classification] = 0
	}
	pythonRows := 0
	for _, row := range matrix.Rows {
		if strings.TrimSpace(row.ID) == "" || strings.TrimSpace(row.Behavior) == "" {
			t.Fatalf("matrix row has empty identity or behavior: %#v", row)
		}
		if _, duplicate := rowIDs[row.ID]; duplicate {
			t.Fatalf("duplicate matrix row %q", row.ID)
		}
		rowIDs[row.ID] = struct{}{}
		if _, ok := wantClassifications[row.Classification]; !ok {
			t.Fatalf("row %q has unknown classification %q", row.ID, row.Classification)
		}
		actualClassifications[row.Classification]++
		if strings.TrimSpace(row.GoDestination) == "" {
			t.Fatalf("row %q has no Go destination", row.ID)
		}
		destinations[row.GoDestination] = struct{}{}
		if len(row.ContractDimensions) == 0 {
			t.Fatalf("row %q has no contract dimensions", row.ID)
		}
		for _, fixtureID := range row.FixtureIDs {
			if _, ok := fixtureIDs[fixtureID]; !ok {
				t.Fatalf("row %q references unknown fixture %q", row.ID, fixtureID)
			}
		}
		if row.Classification == "designed-delta" || row.Classification == "retired" {
			if row.Rationale == nil || strings.TrimSpace(*row.Rationale) == "" {
				t.Fatalf("row %q requires a rationale", row.ID)
			}
		}
		if row.SourceKind == "python-test" {
			pythonRows++
			if !strings.HasPrefix(row.SourcePath, "tests/") ||
				row.SourceSuite == nil || strings.TrimSpace(*row.SourceSuite) == "" ||
				row.SourceTest == nil || strings.TrimSpace(*row.SourceTest) == "" {
				t.Fatalf("row %q lost its frozen source evidence", row.ID)
			}
		}
	}
	if fmt.Sprint(actualClassifications) != fmt.Sprint(wantClassifications) {
		t.Fatalf("row classifications = %v, want %v", actualClassifications, wantClassifications)
	}
	if pythonRows != manifest.SourceSuite.TestCount {
		t.Fatalf("frozen source rows = %d, want %d", pythonRows, manifest.SourceSuite.TestCount)
	}
	if len(destinations) != 26 {
		t.Fatalf("Go destination count = %d, want 26", len(destinations))
	}

	assertBaselineCompatibilityFixtures(t, manifest, fixtureIDs)
	assertBaselineCompatibilityBlobs(t)
}

func regenerateBaselineCompatibilityCorpus(t *testing.T) {
	t.Helper()
	const fixtureRelative = "fixtures/asset-sync.json"
	fixturePath := filepath.FromSlash(path.Join(baselineCompatibilityRoot, fixtureRelative))
	fixtureBytes, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read compatibility fixture for regeneration: %v", err)
	}
	fixture := decodeBaselineCompatibilityDocument(t, fixtureRelative, fixtureBytes)
	setupDigests := regenerateBaselineCompatibilitySetups(t, fixture)
	updateBaselineCompatibilityLedger(t, fixture, setupDigests)
	updatedFixture, err := json.MarshalIndent(fixture, "", "  ")
	if err != nil {
		t.Fatalf("encode regenerated compatibility fixture: %v", err)
	}
	updatedFixture = append(updatedFixture, '\n')
	writeBaselineDerivedArtifact(t, fixturePath, updatedFixture)

	manifestPath := filepath.FromSlash(path.Join(baselineCompatibilityRoot, "manifest.json"))
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read compatibility manifest for regeneration: %v", err)
	}
	manifest := decodeBaselineCompatibilityDocument(t, "manifest.json", manifestBytes)
	frozenBefore := baselineCompatibilityFrozenFields(t, manifestBytes)
	artifacts, ok := manifest["artifacts"].([]any)
	if !ok {
		t.Fatal("compatibility manifest artifacts are not an array")
	}
	updated := false
	sum := sha256.Sum256(updatedFixture)
	for _, rawArtifact := range artifacts {
		artifact, ok := rawArtifact.(map[string]any)
		if !ok {
			t.Fatal("compatibility manifest artifact is not an object")
		}
		if artifact["path"] != fixtureRelative {
			continue
		}
		artifact["bytes"] = json.Number(fmt.Sprint(len(updatedFixture)))
		artifact["sha256"] = hex.EncodeToString(sum[:])
		updated = true
		break
	}
	if !updated {
		t.Fatalf("compatibility manifest has no artifact %q", fixtureRelative)
	}
	updatedManifest, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("encode regenerated compatibility manifest: %v", err)
	}
	updatedManifest = append(updatedManifest, '\n')
	frozenAfter := baselineCompatibilityFrozenFields(t, updatedManifest)
	for key, before := range frozenBefore {
		if !bytes.Equal(before, frozenAfter[key]) {
			t.Fatalf("compatibility manifest frozen field %q changed during regeneration", key)
		}
	}
	writeBaselineDerivedArtifact(t, manifestPath, updatedManifest)
}

func regenerateBaselineCompatibilitySetups(
	t *testing.T,
	fixture map[string]any,
) map[string]string {
	t.Helper()
	manifest, ok := fixture["manifest"].(map[string]any)
	if !ok {
		t.Fatal("asset-sync fixture manifest is not an object")
	}
	rawSetups, ok := manifest["setups"].([]any)
	if !ok {
		t.Fatal("asset-sync fixture setups are not an array")
	}
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	digests := make(map[string]string, len(rawSetups))
	for _, rawSetup := range rawSetups {
		setup, ok := rawSetup.(map[string]any)
		if !ok {
			t.Fatal("asset-sync fixture setup is not an object")
		}
		setupID, ok := setup["id"].(string)
		if !ok || setupID == "" {
			t.Fatal("asset-sync fixture setup has no id")
		}
		rawSkills, ok := setup["skills"].([]any)
		if !ok {
			t.Fatalf("asset-sync fixture setup %q skills are not an array", setupID)
		}
		for _, rawSkill := range rawSkills {
			skill, ok := rawSkill.(map[string]any)
			if !ok {
				t.Fatalf("asset-sync fixture setup %q skill is not an object", setupID)
			}
			source, ok := skill["source"].(map[string]any)
			if !ok || source["type"] != "repo" {
				continue
			}
			name, ok := skill["name"].(string)
			if !ok || name == "" || source["name"] != "roundfix" {
				t.Fatalf("asset-sync fixture setup %q has invalid repo skill", setupID)
			}
			digest, err := skills.SkillFolderHash(
				t.Context(),
				filepath.Join(repoRoot, ".agents", "skills", name),
			)
			if err != nil {
				t.Fatalf("hash canonical skill %q: %v", name, err)
			}
			skill["contentDigest"] = digest
		}
		var digestPayload any = rawSkills
		if bundles, exists := setup["activationBundles"]; exists {
			digestPayload = map[string]any{
				"activationBundles": bundles,
				"skills":            rawSkills,
			}
		}
		digest, err := canonicalSHA256(digestPayload)
		if err != nil {
			t.Fatalf("compute asset-sync fixture setup %q digest: %v", setupID, err)
		}
		setup["digest"] = digest
		digests[setupID] = digest
	}
	return digests
}

func updateBaselineCompatibilityLedger(
	t *testing.T,
	fixture map[string]any,
	setupDigests map[string]string,
) {
	t.Helper()
	ledger, ok := fixture["managedEntryLedger"].([]any)
	if !ok {
		t.Fatal("asset-sync fixture managed entry ledger is not an array")
	}
	for _, rawEntry := range ledger {
		entry, ok := rawEntry.(map[string]any)
		if !ok {
			t.Fatal("asset-sync fixture managed entry is not an object")
		}
		id, _ := entry["id"].(string)
		digest, ok := setupDigests[id]
		if !ok {
			t.Fatalf("asset-sync fixture managed entry %q has no setup", id)
		}
		entry["digest"] = digest
	}
}

func decodeBaselineCompatibilityDocument(
	t *testing.T,
	name string,
	data []byte,
) map[string]any {
	t.Helper()
	var document map[string]any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&document); err != nil {
		t.Fatalf("decode compatibility document %q: %v", name, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		t.Fatalf("compatibility document %q has trailing JSON", name)
	}
	return document
}

func baselineCompatibilityFrozenFields(t *testing.T, data []byte) map[string][]byte {
	t.Helper()
	var document map[string]json.RawMessage
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatalf("decode compatibility manifest frozen fields: %v", err)
	}
	result := make(map[string][]byte)
	for _, key := range []string{
		"frozenDate",
		"sourceSuite",
		"profiles",
		"adoptionStates",
		"fixtureIds",
		"retiredBehavior",
	} {
		value, ok := document[key]
		if !ok {
			t.Fatalf("compatibility manifest frozen field %q is missing", key)
		}
		result[key] = append([]byte(nil), value...)
	}
	return result
}

func assertBaselineCompatibilityArtifacts(t *testing.T, manifest baselineCompatibilityManifest) {
	t.Helper()
	declared := make(map[string]struct{}, len(manifest.Artifacts))
	for _, artifact := range manifest.Artifacts {
		if _, duplicate := declared[artifact.Path]; duplicate {
			t.Fatalf("duplicate manifest artifact %q", artifact.Path)
		}
		declared[artifact.Path] = struct{}{}
		data := readBaselineCompatibilityFile(t, artifact.Path)
		sum := sha256.Sum256(data)
		if len(data) != artifact.Bytes || hex.EncodeToString(sum[:]) != artifact.SHA256 {
			t.Fatalf(
				"artifact %q identity = bytes:%d sha256:%x, want bytes:%d sha256:%s; %s",
				artifact.Path,
				len(data),
				sum,
				artifact.Bytes,
				artifact.SHA256,
				baselineDigestRegenerationHint,
			)
		}
	}

	var actual []string
	if err := fs.WalkDir(baselineCompatibilityCorpus, baselineCompatibilityRoot, func(filePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filePath == path.Join(baselineCompatibilityRoot, "manifest.json") {
			return nil
		}
		actual = append(actual, strings.TrimPrefix(filePath, baselineCompatibilityRoot+"/"))
		return nil
	}); err != nil {
		t.Fatalf("walk compatibility corpus: %v", err)
	}
	sort.Strings(actual)
	if fmt.Sprint(actual) != fmt.Sprint(sortedBaselineCompatibilityKeys(declared)) {
		t.Fatalf("manifest artifacts = %v, actual = %v", sortedBaselineCompatibilityKeys(declared), actual)
	}
}

func assertBaselineCompatibilityFixtures(
	t *testing.T,
	manifest baselineCompatibilityManifest,
	fixtureIDs map[string]struct{},
) {
	t.Helper()
	profiles := baselineCompatibilitySet(t, "profile", manifest.Profiles)
	adoptionStates := baselineCompatibilitySet(t, "adoption state", manifest.AdoptionStates)
	observedProfiles := make(map[string]struct{})
	observedStates := make(map[string]struct{})
	for fixtureID := range fixtureIDs {
		fixture := readBaselineCompatibilityJSON[baselineCompatibilityFixture](
			t,
			path.Join("fixtures", fixtureID+".json"),
		)
		if fixture.SchemaVersion != "setup-context-driven/parity-fixture/v1" ||
			fixture.ID != fixtureID ||
			fixture.Classification != "exact" {
			t.Fatalf("fixture %q identity = %#v", fixtureID, fixture)
		}
		if _, ok := adoptionStates[fixture.AdoptionState]; !ok {
			t.Fatalf("fixture %q has unknown adoption state %q", fixtureID, fixture.AdoptionState)
		}
		observedStates[fixture.AdoptionState] = struct{}{}
		if fixture.Profile != nil {
			if _, ok := profiles[*fixture.Profile]; !ok {
				t.Fatalf("fixture %q has unknown profile %q", fixtureID, *fixture.Profile)
			}
			observedProfiles[*fixture.Profile] = struct{}{}
		}
		if len(fixture.ContractAreas) == 0 || len(fixture.SourceTests) == 0 {
			t.Fatalf("fixture %q lost contract or source evidence", fixtureID)
		}
		if fixture.PlanDigest != nil {
			assertBaselineCompatibilityHex(t, *fixture.PlanDigest)
		}
	}
	if fmt.Sprint(sortedBaselineCompatibilityKeys(observedProfiles)) !=
		fmt.Sprint(sortedBaselineCompatibilityKeys(profiles)) {
		t.Fatalf("observed profiles = %v, want %v", observedProfiles, profiles)
	}
	if fmt.Sprint(sortedBaselineCompatibilityKeys(observedStates)) !=
		fmt.Sprint(sortedBaselineCompatibilityKeys(adoptionStates)) {
		t.Fatalf("observed adoption states = %v, want %v", observedStates, adoptionStates)
	}
}

func assertBaselineCompatibilityBlobs(t *testing.T) {
	t.Helper()
	blobs := readBaselineCompatibilityJSON[baselineCompatibilityBlobs](t, "blobs.json")
	if blobs.SchemaVersion != "setup-context-driven/parity-blobs/v1" || blobs.Encoding != "base64" {
		t.Fatalf("blob contract = schema:%q encoding:%q", blobs.SchemaVersion, blobs.Encoding)
	}
	if len(blobs.Blobs) != 175 {
		t.Fatalf("blob count = %d, want 175", len(blobs.Blobs))
	}
	for digest, encoded := range blobs.Blobs {
		assertBaselineCompatibilityDigest(t, digest)
		data, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			t.Fatalf("decode blob %q: %v", digest, err)
		}
		sum := sha256.Sum256(data)
		if digest != "sha256:"+hex.EncodeToString(sum[:]) {
			t.Fatalf("blob %q bytes hash to sha256:%x", digest, sum)
		}
	}
}

func readBaselineCompatibilityJSON[T any](t *testing.T, relative string) T {
	t.Helper()
	var value T
	data := readBaselineCompatibilityFile(t, relative)
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		t.Fatalf("decode compatibility corpus %q: %v", relative, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		t.Fatalf("compatibility corpus %q has trailing JSON", relative)
	}
	return value
}

func readBaselineCompatibilityFile(t *testing.T, relative string) []byte {
	t.Helper()
	data, err := baselineCompatibilityCorpus.ReadFile(path.Join(baselineCompatibilityRoot, relative))
	if err != nil {
		t.Fatalf("read compatibility corpus %q: %v", relative, err)
	}
	return data
}

func baselineCompatibilitySet(t *testing.T, kind string, values []string) map[string]struct{} {
	t.Helper()
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			t.Fatalf("%s is empty", kind)
		}
		if _, duplicate := result[value]; duplicate {
			t.Fatalf("duplicate %s %q", kind, value)
		}
		result[value] = struct{}{}
	}
	return result
}

func assertBaselineCompatibilityDigest(t *testing.T, digest string) {
	t.Helper()
	if !strings.HasPrefix(digest, "sha256:") {
		t.Fatalf("digest %q has no sha256 domain", digest)
	}
	assertBaselineCompatibilityHex(t, strings.TrimPrefix(digest, "sha256:"))
}

func assertBaselineCompatibilityHex(t *testing.T, value string) {
	t.Helper()
	if len(value) != sha256.Size*2 {
		t.Fatalf("sha256 value %q has length %d", value, len(value))
	}
	if _, err := hex.DecodeString(value); err != nil {
		t.Fatalf("sha256 value %q is not lowercase hexadecimal: %v", value, err)
	}
	if strings.ToLower(value) != value {
		t.Fatalf("sha256 value %q is not lowercase", value)
	}
}

func sortedBaselineCompatibilityStrings(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	return result
}

func sortedBaselineCompatibilityKeys[T any](values map[string]T) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
