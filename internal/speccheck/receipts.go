package speccheck

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Receipt is one source-and-quote pair an artifact wrote.
type Receipt struct {
	Artifact string
	Line     int
	Source   string
	Quote    string
}

// ReceiptedClaim pairs an attribution with receipts naming its record in the same paragraph.
type ReceiptedClaim struct {
	Claim    Claim
	Receipts []Receipt
}

var receiptPattern = regexp.MustCompile("(?:\\bADR-[0-9]{4}\\b|`[^`\\s\\p{Z}\\x{0085}]+`):[\\s\\p{Z}\\x{0085}]+\"[^\"]*\"")
var receiptDecisionPattern = regexp.MustCompile(`^ADR-[0-9]{4}$`)

// Receipts reads written receipts outside fenced blocks.
func Receipts(artifact string, content []byte) []Receipt {
	var receipts []Receipt
	walkCitationParagraphs(content, func(line int, paragraph string) {
		receipts = append(receipts, receiptsInParagraph(artifact, line, paragraph)...)
	})
	return receipts
}

func receiptsInParagraph(artifact string, line int, paragraph string) []Receipt {
	var receipts []Receipt
	for _, match := range receiptPattern.FindAllStringIndex(paragraph, -1) {
		text := paragraph[match[0]:match[1]]
		colon := strings.Index(text, ":")
		// A path token may contain a colon; its closing backtick defines the source.
		if strings.HasPrefix(text, "`") {
			colon = strings.Index(text[1:], "`") + 2
		}
		source := strings.Trim(text[:colon], "`")
		if strings.HasPrefix(text, "`") && !strings.Contains(source, "/") && (filepath.Ext(source) == "" || filepath.Ext(source) == ".") {
			continue
		}
		quoted := strings.TrimSpace(text[colon+1:])
		if strings.Count(text, "\n") > 1 {
			continue
		}
		receipts = append(receipts, Receipt{Artifact: artifact, Line: line + strings.Count(paragraph[:match[0]], "\n"), Source: source, Quote: normalizeReceiptText(quoted[1 : len(quoted)-1])})
	}
	return receipts
}

// ReceiptedClaims keeps the attribution grammar and pairs only paragraph-local receipts.
func ReceiptedClaims(artifact string, content []byte) []ReceiptedClaim {
	var claims []ReceiptedClaim
	walkCitationParagraphs(content, func(line int, paragraph string) {
		receipts := receiptsInParagraph(artifact, line, paragraph)
		for _, claim := range citationClaimsInParagraph(artifact, line, paragraph) {
			paired := ReceiptedClaim{Claim: claim}
			for _, receipt := range receipts {
				if receipt.Source == claim.Target {
					paired.Receipts = append(paired.Receipts, receipt)
				}
			}
			claims = append(claims, paired)
		}
	})
	return claims
}

func normalizeReceiptText(text string) string { return strings.Join(strings.Fields(text), " ") }

// ProveReceipt resolves a source and proves a case-sensitive contiguous quote.
func ProveReceipt(repoRoot string, receipt Receipt) (sourcePath string, proven bool, reason string, err error) {
	quote := normalizeReceiptText(receipt.Quote)
	if len(strings.Fields(quote)) < 3 {
		return "", false, "a receipt needs at least three words", nil
	}
	sourcePath, resolved, err := resolveReceiptSource(repoRoot, receipt.Source)
	if err != nil {
		return "", false, "", err
	}
	if !resolved {
		return "", false, "does not resolve to a file", nil
	}
	content, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(sourcePath)))
	if err != nil {
		return sourcePath, false, "", fmt.Errorf("read receipt source %q: %w", sourcePath, err)
	}
	if !strings.Contains(normalizeReceiptText(string(content)), quote) {
		return sourcePath, false, "does not contain that text", nil
	}
	return sourcePath, true, "", nil
}

func resolveReceiptSource(repoRoot, source string) (string, bool, error) {
	if !receiptDecisionPattern.MatchString(source) {
		return receiptFilePath(repoRoot, source)
	}
	number := strings.TrimPrefix(source, "ADR-")
	for _, directory := range []string{"docs/adr", "docs/history/adr"} {
		// Check the directory components too; Glob must never traverse a symlink.
		if _, ok, err := receiptPathInfo(repoRoot, directory); err != nil {
			return "", false, err
		} else if !ok {
			continue
		}
		matches, err := filepath.Glob(filepath.Join(repoRoot, filepath.FromSlash(directory), number+"-*.md"))
		if err != nil {
			return "", false, fmt.Errorf("resolve receipt record %q: %w", source, err)
		}
		if len(matches) == 0 {
			continue
		}
		if len(matches) != 1 {
			return "", false, nil
		}
		relative, err := filepath.Rel(repoRoot, matches[0])
		if err != nil {
			return "", false, err
		}
		return receiptFilePath(repoRoot, filepath.ToSlash(relative))
	}
	return "", false, nil
}

func receiptFilePath(repoRoot, relative string) (string, bool, error) {
	info, ok, err := receiptPathInfo(repoRoot, relative)
	if !ok || err != nil {
		return "", false, err
	}
	if !info.Mode().IsRegular() {
		return "", false, nil
	}
	return relative, true, nil
}

func receiptPathInfo(repoRoot, relative string) (os.FileInfo, bool, error) {
	path := filepath.FromSlash(relative)
	if path == "." || path == "" || filepath.IsAbs(path) || filepath.Clean(path) != path || path == ".." || strings.HasPrefix(path, ".."+string(filepath.Separator)) {
		return nil, false, nil
	}
	current := filepath.Clean(repoRoot)
	components := strings.Split(path, string(filepath.Separator))
	for index, component := range components {
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		if err != nil {
			return nil, false, fmt.Errorf("inspect receipt path %q: %w", relative, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil, false, nil
		}
		if index == len(components)-1 {
			return info, true, nil
		}
		if !info.IsDir() {
			return nil, false, nil
		}
	}
	return nil, false, nil
}

func detectReceipts(result *Result, repoRoot string, horizon contractHorizon, artifacts []string) error {
	if !horizon.held {
		addSkip(result, CodeReceiptMissing, horizon.missing)
	}
	for _, path := range artifacts {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read receipt artifact %q: %w", path, err)
		}
		artifact := artifactDisplayPath(repoRoot, path)
		for _, receipt := range Receipts(artifact, content) {
			sourcePath, proven, reason, err := ProveReceipt(repoRoot, receipt)
			if err != nil {
				return err
			}
			if proven {
				continue
			}
			source := sourcePath
			if source == "" {
				source = receipt.Source
			}
			where := []Location{{Path: artifact, Line: receipt.Line}}
			if sourcePath != "" {
				where = append(where, Location{Path: sourcePath, Line: 1})
			}
			result.Findings = append(result.Findings, Finding{Code: CodeReceiptUnproven, Severity: SeverityError,
				Summary: artifact + " quotes " + receipt.Source + " as " + strconv.Quote(receipt.Quote) + ", but " + source + " " + reason,
				Where:   where, Fix: "Copy the passage verbatim from " + source + " into the receipt, or remove the receipt."})
		}
		if !horizon.held {
			continue
		}
		var uncovered []Claim
		for _, paired := range ReceiptedClaims(artifact, content) {
			if len(paired.Receipts) == 0 {
				uncovered = append(uncovered, paired.Claim)
			}
		}
		if len(uncovered) == 0 {
			continue
		}
		resolved, err := resolveCitationClaims(repoRoot, uncovered)
		if err != nil {
			return err
		}
		for _, resolution := range resolved {
			claim := resolution.claim
			result.Findings = append(result.Findings, Finding{Code: CodeReceiptMissing, Severity: SeverityGap,
				Summary: artifact + " attributes " + strconv.Quote(claim.sentence) + " to " + claim.Target + " without a receipt",
				Where:   []Location{{Path: artifact, Line: claim.Line}, {Path: resolution.record.DisplayPath, Line: resolution.record.TitleLine}},
				Fix:     "Add " + claim.Target + ": \"<verbatim passage>\" to the same paragraph in " + artifact + "."})
		}
	}
	return nil
}
