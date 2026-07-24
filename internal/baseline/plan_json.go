package baseline

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// MarshalPlanDocument writes the deterministic portable representation.
func MarshalPlanDocument(document PlanDocument) ([]byte, error) {
	if err := ValidatePlanDocument(document); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode Baseline Plan: %w", err)
	}
	return append(data, '\n'), nil
}

// ParsePlanDocument rejects duplicate keys, unknown fields, trailing values,
// projection mismatches, and digest mismatches.
func ParsePlanDocument(data []byte) (PlanDocument, error) {
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return PlanDocument{}, fmt.Errorf("decode Baseline Plan: %w", err)
	}
	var document PlanDocument
	if err := decodeStrictJSON(data, &document); err != nil {
		return PlanDocument{}, fmt.Errorf("decode Baseline Plan: %w", err)
	}
	if err := ValidatePlanDocument(document); err != nil {
		return PlanDocument{}, err
	}
	return document, nil
}

// MarshalResult writes one strict Baseline automation result.
func MarshalResult(result Result) ([]byte, error) {
	if err := validateResult(result); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode Baseline result: %w", err)
	}
	return append(data, '\n'), nil
}

// ParseResult strictly parses roundfix/baseline-result/v1.
func ParseResult(data []byte) (Result, error) {
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return Result{}, fmt.Errorf("decode Baseline result: %w", err)
	}
	var result Result
	if err := decodeStrictJSON(data, &result); err != nil {
		return Result{}, fmt.Errorf("decode Baseline result: %w", err)
	}
	if err := validateResult(result); err != nil {
		return Result{}, err
	}
	return result, nil
}

func validateResult(result Result) error {
	if result.SchemaVersion != ResultSchemaVersion {
		return fmt.Errorf("unsupported Baseline result schema %q", result.SchemaVersion)
	}
	if result.Operation == "" || result.State == "" {
		return errors.New("Baseline result requires operation and state")
	}
	if result.VerifiedPostimages == nil || result.Warnings == nil || result.Recommendations == nil {
		return errors.New("Baseline result verifiedPostimages, warnings, and recommendations must be arrays")
	}
	return nil
}

func decodeStrictJSON(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values are not allowed")
		}
		return err
	}
	return nil
}

func rejectDuplicateJSONKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	var walk func() error
	walk = func() error {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		switch delimiter := token.(type) {
		case json.Delim:
			switch delimiter {
			case '{':
				seen := make(map[string]struct{})
				for decoder.More() {
					keyToken, err := decoder.Token()
					if err != nil {
						return err
					}
					key, ok := keyToken.(string)
					if !ok {
						return errors.New("JSON object key is not a string")
					}
					if _, duplicate := seen[key]; duplicate {
						return fmt.Errorf("duplicate JSON key %q", key)
					}
					seen[key] = struct{}{}
					if err := walk(); err != nil {
						return err
					}
				}
				_, err = decoder.Token()
				return err
			case '[':
				for decoder.More() {
					if err := walk(); err != nil {
						return err
					}
				}
				_, err = decoder.Token()
				return err
			default:
				return fmt.Errorf("unexpected JSON delimiter %q", delimiter)
			}
		default:
			return nil
		}
	}
	if err := walk(); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values are not allowed")
		}
		return err
	}
	return nil
}
