package baseline

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"sort"
	"strings"
)

const (
	sourceBaselineSchema         = "setup-context-driven/source-baseline/0.0.1"
	sourceBaselineManifestSchema = "setup-context-driven/source-baseline-manifest/0.0.1"
	sourceAccountingSchema       = "setup-context-driven/source-accounting/0.0.1"
	upgradeTransitionSchema      = "setup-context-driven/upgrade-transition-v1"
)

var sourceMachinePath = regexp.MustCompile(
	"(?i)(?:^|[\\s`\"'])/(?:users|home)/|(?:^|[\\s`\"'])[a-z]:\\\\users\\\\",
)

// SourceBaselineIdentity identifies one immutable maintained corpus.
type SourceBaselineIdentity struct {
	SchemaVersion       string   `json:"schemaVersion"`
	Version             string   `json:"version"`
	ID                  string   `json:"id"`
	Profile             string   `json:"profile"`
	Corpus              string   `json:"corpus"`
	Manifest            string   `json:"manifest"`
	EntryCount          int      `json:"entryCount"`
	CorpusDigest        string   `json:"corpusDigest"`
	ManifestDigest      string   `json:"manifestDigest"`
	DeniedProjectTokens []string `json:"deniedProjectTokens"`
}

// SourceBaselineEntry independently accounts for one maintained Normative
// Clause, recommendation, or Operational Contract.
type SourceBaselineEntry struct {
	ID          string  `json:"id"`
	Path        string  `json:"path"`
	Kind        string  `json:"kind"`
	Enforcement string  `json:"enforcement"`
	Carrier     string  `json:"carrier"`
	Structure   *string `json:"structure"`
	Start       int     `json:"start"`
	End         int     `json:"end"`
	Digest      string  `json:"digest"`
}

// SourceBaseline is one fully validated maintained Source Baseline.
type SourceBaseline struct {
	Identity   SourceBaselineIdentity  `json:"identity"`
	Entries    []SourceBaselineEntry   `json:"entries"`
	Accounting []SourceAccountingEntry `json:"accounting"`
}

// SourceAccountingEntry records how one earlier source entry is retained by a
// maintained Source Baseline.
type SourceAccountingEntry struct {
	SourceBaseline string   `json:"sourceBaseline"`
	SourceEntry    string   `json:"sourceEntry"`
	Disposition    string   `json:"disposition"`
	Targets        []string `json:"targets"`
	Reason         string   `json:"reason"`
}

type sourceAccountingDocument struct {
	SchemaVersion string                  `json:"schemaVersion"`
	Version       string                  `json:"version"`
	Baseline      string                  `json:"baseline"`
	Entries       []SourceAccountingEntry `json:"entries"`
}

type sourceBaselineIndex struct {
	SchemaVersion string                      `json:"schemaVersion"`
	Version       string                      `json:"version"`
	Baselines     []sourceBaselineIndexRecord `json:"baselines"`
}

type sourceBaselineIndexRecord struct {
	ID             string   `json:"id"`
	Path           string   `json:"path"`
	Profile        string   `json:"profile"`
	EntryIDs       []string `json:"entryIds"`
	EntryCount     int      `json:"entryCount"`
	CorpusDigest   string   `json:"corpusDigest"`
	ManifestDigest string   `json:"manifestDigest"`
}

type sourceBaselineManifest struct {
	SchemaVersion string                `json:"schemaVersion"`
	Version       string                `json:"version"`
	ID            string                `json:"id"`
	Profile       string                `json:"profile"`
	Entries       []SourceBaselineEntry `json:"entries"`
}

// UpgradeRetentionClause is one prior Normative Clause requiring independent
// accounting.
type UpgradeRetentionClause struct {
	ID             string `json:"id"`
	Enforcement    string `json:"enforcement"`
	Carrier        string `json:"carrier"`
	GuidanceDigest string `json:"guidanceDigest"`
}

// UpgradeRetentionAccounting is one exact prior-clause disposition.
type UpgradeRetentionAccounting struct {
	FromClause  string   `json:"fromClause"`
	Disposition string   `json:"disposition"`
	Targets     []string `json:"targets"`
	Reason      string   `json:"reason"`
}

// UpgradeRetentionContract is one maintained clause-level transition.
type UpgradeRetentionContract struct {
	SchemaVersion              string                       `json:"schemaVersion"`
	ID                         string                       `json:"id"`
	Version                    int                          `json:"version"`
	FromBaseline               string                       `json:"fromBaseline"`
	ToBaseline                 string                       `json:"toBaseline"`
	LegacyManifestFingerprints []string                     `json:"legacyManifestFingerprints"`
	PriorClauses               []UpgradeRetentionClause     `json:"priorClauses"`
	Accounting                 []UpgradeRetentionAccounting `json:"mappings"`
}

// SourceBaseline loads and verifies one maintained Source Baseline from the
// immutable catalog.
func (c *Catalog) SourceBaseline(id string) (SourceBaseline, error) {
	indexAsset, ok := c.Asset("source-baselines/index.json")
	if !ok {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q: index is missing", id)
	}
	var index sourceBaselineIndex
	if err := strictJSON(indexAsset.Data, &index); err != nil {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline index: %w", err)
	}
	var record *sourceBaselineIndexRecord
	for position := range index.Baselines {
		if index.Baselines[position].ID == id {
			record = &index.Baselines[position]
			break
		}
	}
	if record == nil {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q: unknown id", id)
	}
	if !safeRelative(record.Path) {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q: unsafe catalog path %q", id, record.Path)
	}

	prefix := path.Join("source-baselines", record.Path)
	identityAsset, ok := c.Asset(path.Join(prefix, "baseline.json"))
	if !ok {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q: identity is missing", id)
	}
	var identity SourceBaselineIdentity
	if err := strictJSON(identityAsset.Data, &identity); err != nil {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q identity: %w", id, err)
	}
	if identity.SchemaVersion != sourceBaselineSchema ||
		identity.ID != id ||
		identity.Profile != record.Profile ||
		identity.Manifest == "" ||
		identity.Corpus == "" {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q: identity headers disagree", id)
	}

	manifestPath := path.Join(prefix, identity.Manifest)
	manifestAsset, ok := c.Asset(manifestPath)
	if !ok {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q: manifest is missing", id)
	}
	var manifest sourceBaselineManifest
	if err := strictJSON(manifestAsset.Data, &manifest); err != nil {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q manifest: %w", id, err)
	}
	if manifest.SchemaVersion != sourceBaselineManifestSchema ||
		manifest.ID != id ||
		manifest.Profile != identity.Profile {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q: manifest headers disagree", id)
	}
	manifestSum := sha256.Sum256(manifestAsset.Data)
	manifestDigest := hex.EncodeToString(manifestSum[:])
	corpusDigest, err := c.sourceCorpusDigest(prefix, identity.Corpus)
	if err != nil {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q corpus: %w", id, err)
	}
	if identity.EntryCount != len(manifest.Entries) ||
		record.EntryCount != len(manifest.Entries) ||
		len(record.EntryIDs) != len(manifest.Entries) ||
		identity.ManifestDigest != manifestDigest ||
		record.ManifestDigest != manifestDigest ||
		identity.CorpusDigest != corpusDigest ||
		record.CorpusDigest != corpusDigest {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q: counts or digests disagree", id)
	}
	for index, entry := range manifest.Entries {
		if record.EntryIDs[index] != entry.ID {
			return SourceBaseline{}, fmt.Errorf("load Source Baseline %q: entry order differs at %d", id, index)
		}
		if err := c.validateSourceBaselineEntry(prefix, entry); err != nil {
			return SourceBaseline{}, fmt.Errorf("load Source Baseline %q: %w", id, err)
		}
	}
	entryIDs := make(map[string]struct{}, len(manifest.Entries))
	for _, entry := range manifest.Entries {
		entryIDs[entry.ID] = struct{}{}
	}
	accounting, err := c.loadSourceAccounting(prefix, id, entryIDs)
	if err != nil {
		return SourceBaseline{}, err
	}
	if err := c.validateSourceCorpusPolicy(
		prefix,
		identity.Corpus,
		identity.DeniedProjectTokens,
	); err != nil {
		return SourceBaseline{}, fmt.Errorf("load Source Baseline %q: %w", id, err)
	}
	return SourceBaseline{
		Identity:   identity,
		Entries:    append([]SourceBaselineEntry(nil), manifest.Entries...),
		Accounting: accounting,
	}, nil
}

func (c *Catalog) loadSourceAccounting(
	prefix, baselineID string,
	currentEntries map[string]struct{},
) ([]SourceAccountingEntry, error) {
	asset, ok := c.Asset(path.Join(prefix, "accounting.json"))
	if !ok {
		return []SourceAccountingEntry{}, nil
	}
	var document sourceAccountingDocument
	if err := strictJSON(asset.Data, &document); err != nil {
		return nil, fmt.Errorf("load Source Baseline %q accounting: %w", baselineID, err)
	}
	if document.SchemaVersion != sourceAccountingSchema ||
		document.Baseline != baselineID {
		return nil, fmt.Errorf("load Source Baseline %q: accounting headers disagree", baselineID)
	}
	seen := make(map[string]struct{}, len(document.Entries))
	for _, entry := range document.Entries {
		key := entry.SourceBaseline + "\x00" + entry.SourceEntry
		if entry.SourceBaseline == "" ||
			entry.SourceEntry == "" ||
			strings.TrimSpace(entry.Reason) == "" {
			return nil, fmt.Errorf("load Source Baseline %q: invalid accounting for %q", baselineID, entry.SourceEntry)
		}
		if _, duplicate := seen[key]; duplicate {
			return nil, fmt.Errorf("load Source Baseline %q: duplicate accounting for %q", baselineID, entry.SourceEntry)
		}
		seen[key] = struct{}{}
		if entry.Disposition == "rejected" {
			if len(entry.Targets) != 0 {
				return nil, fmt.Errorf("load Source Baseline %q: rejected accounting for %q has targets", baselineID, entry.SourceEntry)
			}
		} else if len(entry.Targets) == 0 {
			return nil, fmt.Errorf("load Source Baseline %q: accounting for %q has no target", baselineID, entry.SourceEntry)
		}
		for _, target := range entry.Targets {
			if _, ok := currentEntries[target]; !ok {
				return nil, fmt.Errorf(
					"load Source Baseline %q: accounting for %q targets unknown entry %q",
					baselineID,
					entry.SourceEntry,
					target,
				)
			}
		}
	}
	return append([]SourceAccountingEntry(nil), document.Entries...), nil
}

func (c *Catalog) validateSourceCorpusPolicy(
	prefix, corpusDir string,
	deniedProjectTokens []string,
) error {
	if len(deniedProjectTokens) == 0 || !uniqueStrings(deniedProjectTokens) {
		return errors.New("Source Baseline denied project tokens are invalid")
	}
	for _, token := range deniedProjectTokens {
		if strings.TrimSpace(token) == "" || strings.ToLower(token) != token {
			return errors.New("Source Baseline denied project tokens are invalid")
		}
	}

	corpusPrefix := path.Join(prefix, corpusDir) + "/"
	for _, assetPath := range sortedKeys(c.assets) {
		if !strings.HasPrefix(assetPath, corpusPrefix) {
			continue
		}
		text := string(c.assets[assetPath])
		folded := strings.ToLower(text)
		for _, token := range deniedProjectTokens {
			if strings.Contains(folded, token) {
				return fmt.Errorf(
					"Source Baseline corpus %q contains denied project token %q",
					strings.TrimPrefix(assetPath, corpusPrefix),
					token,
				)
			}
		}
		if strings.Contains(folded, "<!-- setup-context-driven:") {
			return fmt.Errorf(
				"Source Baseline corpus %q contains a generated managed marker",
				strings.TrimPrefix(assetPath, corpusPrefix),
			)
		}
		if sourceMachinePath.MatchString(text) {
			return fmt.Errorf(
				"Source Baseline corpus %q contains a machine-specific path",
				strings.TrimPrefix(assetPath, corpusPrefix),
			)
		}
	}
	return nil
}

func (c *Catalog) sourceCorpusDigest(prefix, corpusDir string) (string, error) {
	corpusPrefix := path.Join(prefix, corpusDir) + "/"
	var assetPaths []string
	for assetPath := range c.assets {
		if strings.HasPrefix(assetPath, corpusPrefix) {
			assetPaths = append(assetPaths, assetPath)
		}
	}
	if len(assetPaths) == 0 {
		return "", fmt.Errorf("no corpus files under %q", corpusDir)
	}
	sort.Strings(assetPaths)
	digest := sha256.New()
	for _, assetPath := range assetPaths {
		relative := strings.TrimPrefix(assetPath, corpusPrefix)
		data := c.assets[assetPath]
		writeLengthPrefixed(digest, relative)
		var length [8]byte
		binary.BigEndian.PutUint64(length[:], uint64(len(data)))
		_, _ = digest.Write(length[:])
		_, _ = digest.Write(data)
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func (c *Catalog) validateSourceBaselineEntry(prefix string, entry SourceBaselineEntry) error {
	if !safeRelative(entry.Path) ||
		entry.End <= entry.Start ||
		entry.Start < 0 ||
		!isRawSHA256(entry.Digest) {
		return fmt.Errorf("Source Baseline Entry %q has invalid evidence", entry.ID)
	}
	asset, ok := c.Asset(path.Join(prefix, entry.Path))
	if !ok || entry.End > len(asset.Data) {
		return fmt.Errorf("Source Baseline Entry %q source range is unavailable", entry.ID)
	}
	sum := sha256.Sum256(asset.Data[entry.Start:entry.End])
	if hex.EncodeToString(sum[:]) != entry.Digest {
		return fmt.Errorf("Source Baseline Entry %q digest is stale", entry.ID)
	}
	return nil
}

// UpgradeRetentionContract loads one exact clause-level transition from the
// immutable catalog.
func (c *Catalog) UpgradeRetentionContract(id string) (UpgradeRetentionContract, error) {
	entry, ok := c.RetentionTransition(id)
	if !ok {
		return UpgradeRetentionContract{}, fmt.Errorf("load Upgrade Retention Contract %q: unknown id", id)
	}
	var contract UpgradeRetentionContract
	if err := strictJSON(entry.Data, &contract); err != nil {
		return UpgradeRetentionContract{}, fmt.Errorf("load Upgrade Retention Contract %q: %w", id, err)
	}
	if contract.SchemaVersion != upgradeTransitionSchema ||
		contract.ID != id ||
		len(contract.PriorClauses) == 0 {
		return UpgradeRetentionContract{}, fmt.Errorf("load Upgrade Retention Contract %q: invalid headers", id)
	}
	prior := make(map[string]struct{}, len(contract.PriorClauses))
	for _, clause := range contract.PriorClauses {
		if _, duplicate := prior[clause.ID]; duplicate ||
			!isRawSHA256(clause.GuidanceDigest) {
			return UpgradeRetentionContract{}, fmt.Errorf("load Upgrade Retention Contract %q: invalid prior clause %q", id, clause.ID)
		}
		prior[clause.ID] = struct{}{}
	}
	accounted := make(map[string]struct{}, len(contract.Accounting))
	for _, accounting := range contract.Accounting {
		if _, exists := prior[accounting.FromClause]; !exists {
			return UpgradeRetentionContract{}, fmt.Errorf(
				"load Upgrade Retention Contract %q: unknown source clause %q",
				id,
				accounting.FromClause,
			)
		}
		if _, duplicate := accounted[accounting.FromClause]; duplicate {
			return UpgradeRetentionContract{}, fmt.Errorf(
				"load Upgrade Retention Contract %q: duplicate accounting for %q",
				id,
				accounting.FromClause,
			)
		}
		accounted[accounting.FromClause] = struct{}{}
		if accounting.Disposition == "rejected" {
			if len(accounting.Targets) != 0 || strings.TrimSpace(accounting.Reason) == "" {
				return UpgradeRetentionContract{}, fmt.Errorf(
					"load Upgrade Retention Contract %q: invalid rejection for %q",
					id,
					accounting.FromClause,
				)
			}
		} else if len(accounting.Targets) == 0 {
			return UpgradeRetentionContract{}, fmt.Errorf(
				"load Upgrade Retention Contract %q: accounting for %q has no target",
				id,
				accounting.FromClause,
			)
		}
	}
	if len(accounted) != len(prior) {
		return UpgradeRetentionContract{}, fmt.Errorf("load Upgrade Retention Contract %q: clause accounting is incomplete", id)
	}
	return contract, nil
}

func strictJSON(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		if err == nil {
			return fmt.Errorf("unexpected trailing JSON value")
		}
		return err
	}
	return nil
}
