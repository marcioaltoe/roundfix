package spec

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// DeliveryType is the Conventional Commits type a Spec delivers as.
type DeliveryType string

const (
	DeliveryTypeFeat     DeliveryType = "feat"
	DeliveryTypeFix      DeliveryType = "fix"
	DeliveryTypePerf     DeliveryType = "perf"
	DeliveryTypeRefactor DeliveryType = "refactor"
)

// DeliveryTypes returns the supported set in this order.
func DeliveryTypes() []DeliveryType {
	return []DeliveryType{DeliveryTypeFeat, DeliveryTypeFix, DeliveryTypePerf, DeliveryTypeRefactor}
}

// UnsupportedDeliveryTypeError names a PRD type outside the supported set.
type UnsupportedDeliveryTypeError struct{ Value string }

func (e *UnsupportedDeliveryTypeError) Error() string {
	return fmt.Sprintf("PRD type %q is not one of feat, fix, perf, refactor", e.Value)
}

// ParsePRDDeliveryType reads the optional type key of a PRD's frontmatter.
// present is false when the key is absent; an unsupported or non-string
// value returns UnsupportedDeliveryTypeError.
func ParsePRDDeliveryType(prd []byte) (deliveryType DeliveryType, present bool, err error) {
	if !bytes.HasPrefix(prd, []byte("---\n")) {
		return "", false, nil
	}
	front, _, err := splitFrontmatter(prd)
	if err != nil {
		return "", false, fmt.Errorf("read PRD type: %w", err)
	}
	var fields map[string]yaml.Node
	if err := yaml.Unmarshal(front, &fields); err != nil {
		return "", false, fmt.Errorf("parse PRD type: %w", err)
	}
	node, present := fields["type"]
	if !present {
		return "", false, nil
	}
	if node.Kind == yaml.ScalarNode && node.Tag == "!!str" {
		for _, supported := range DeliveryTypes() {
			if DeliveryType(node.Value) == supported {
				return supported, true, nil
			}
		}
	}
	value := node.Value
	if node.Kind != yaml.ScalarNode {
		encoded, err := yaml.Marshal(&node)
		if err != nil {
			return "", true, fmt.Errorf("read PRD type value: %w", err)
		}
		value = strings.TrimSpace(string(encoded))
	}
	return "", true, &UnsupportedDeliveryTypeError{Value: value}
}
