package baseline

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

const (
	authProviderDecisionID = "auth.provider"
	httpContractDecisionID = "http.contract"
)

// HTTPContract is the normalized repository-owned HTTP Contract Decision.
type HTTPContract struct {
	Mode       string          `json:"mode"`
	Exceptions []HTTPException `json:"exceptions"`
	Source     map[string]any  `json:"source,omitempty"`
}

func normalizeProjectDecisions(
	input []DecisionValue,
	catalog *Catalog,
) ([]DecisionValue, error) {
	normalized := cloneDecisionValues(input)
	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i].ID < normalized[j].ID
	})

	httpIndex := decisionValueIndex(normalized, httpContractDecisionID)
	authIndex := decisionValueIndex(normalized, authProviderDecisionID)
	if httpIndex < 0 && authIndex < 0 {
		return normalized, nil
	}
	if catalog == nil {
		return nil, errors.New("normalize project decisions: catalog is required")
	}

	var contract HTTPContract
	if httpIndex >= 0 {
		declaration, ok := catalog.decisions[httpContractDecisionID]
		if !ok {
			return nil, projectDecisionError(errors.New("HTTP Contract Decision declaration is missing"))
		}
		var err error
		contract, err = normalizeHTTPContract(normalized[httpIndex].Value, declaration)
		if err != nil {
			return nil, projectDecisionError(fmt.Errorf("normalize %q: %w", httpContractDecisionID, err))
		}
	}

	var provider *AuthProviderDecision
	if authIndex >= 0 {
		if httpIndex < 0 {
			return nil, projectDecisionError(errors.New("authentication provider requires an HTTP Contract Decision"))
		}
		value, err := normalizeAuthProviderDecision(normalized[authIndex].Value)
		if err != nil {
			return nil, projectDecisionError(fmt.Errorf("normalize %q: %w", authProviderDecisionID, err))
		}
		provider = &value
	}

	if httpIndex >= 0 {
		derived, err := deriveHTTPContract(contract, provider)
		if err != nil {
			return nil, projectDecisionError(err)
		}
		normalized[httpIndex].Value = httpContractValue(derived)
	}
	if authIndex >= 0 {
		normalized[authIndex].Value = authProviderValue(*provider)
	}
	return normalized, nil
}

func projectDecisionError(err error) error {
	return fmt.Errorf(
		"project decisions %q and %q: %w",
		authProviderDecisionID,
		httpContractDecisionID,
		err,
	)
}

func decisionValueIndex(decisions []DecisionValue, id string) int {
	for index := range decisions {
		if decisions[index].ID == id {
			return index
		}
	}
	return -1
}

func normalizeHTTPContract(value any, declaration document) (HTTPContract, error) {
	object, ok := objectValue(value)
	if !ok {
		return HTTPContract{}, errors.New("must be an object")
	}
	if !hasOnlyFields(object, "mode", "exceptions", "source") {
		return HTTPContract{}, errors.New("requires only mode, exceptions, and optional source")
	}
	mode, ok := object["mode"].(string)
	if !ok || strings.TrimSpace(mode) == "" {
		return HTTPContract{}, errors.New("mode is required")
	}
	mode = strings.TrimSpace(mode)
	if !slices.Contains(stringsOrEmpty(declaration["modes"]), mode) {
		return HTTPContract{}, fmt.Errorf("mode %q is not allowed", mode)
	}

	exceptions, err := normalizeHTTPExceptions(object["exceptions"])
	if err != nil {
		return HTTPContract{}, err
	}
	contract := HTTPContract{Mode: mode, Exceptions: exceptions}
	if source, exists := object["source"]; exists {
		normalizedSource, ok := objectValue(source)
		if !ok {
			return HTTPContract{}, errors.New("source must be an object")
		}
		contract.Source = cloneJSONValue(normalizedSource).(map[string]any)
	}
	return contract, nil
}

func normalizeHTTPExceptions(value any) ([]HTTPException, error) {
	if value == nil {
		return []HTTPException{}, nil
	}
	raw, ok := value.([]any)
	if !ok {
		return nil, errors.New("exceptions must be an array")
	}
	normalized := make([]HTTPException, 0, len(raw))
	indexByKey := make(map[string]int, len(raw))
	for index, item := range raw {
		exception, err := normalizeHTTPException(item)
		if err != nil {
			return nil, fmt.Errorf("exception %d: %w", index, err)
		}
		key := normalizedHTTPExceptionKey(exception)
		if existingIndex, duplicate := indexByKey[key]; duplicate {
			if equalHTTPExceptions(normalized[existingIndex], exception) {
				continue
			}
			return nil, fmt.Errorf(
				"conflicting exceptions for owner %q and scope %q",
				exception.Owner,
				exception.Scope,
			)
		}
		indexByKey[key] = len(normalized)
		normalized = append(normalized, exception)
	}
	sortHTTPExceptions(normalized)
	return normalized, nil
}

func normalizeHTTPException(value any) (HTTPException, error) {
	object, ok := objectValue(value)
	if !ok {
		return HTTPException{}, errors.New("must be an object")
	}
	if !hasExactFields(object, "scope", "methods", "owner", "reason") {
		return HTTPException{}, errors.New("requires exactly scope, methods, owner, and reason")
	}
	scope, ok := object["scope"].(string)
	if !ok || strings.TrimSpace(scope) == "" {
		return HTTPException{}, errors.New("scope must be a non-empty string")
	}
	owner, ok := object["owner"].(string)
	if !ok || strings.TrimSpace(owner) == "" {
		return HTTPException{}, errors.New("owner must be a non-empty string")
	}
	reason, ok := object["reason"].(string)
	if !ok || strings.TrimSpace(reason) == "" {
		return HTTPException{}, errors.New("reason must be a non-empty string")
	}
	methods, ok := decisionStringList(object["methods"])
	if !ok || len(methods) == 0 {
		return HTTPException{}, errors.New("methods must be a non-empty string array")
	}
	normalizedMethods := make([]string, len(methods))
	seenMethods := make(map[string]struct{}, len(methods))
	for index, method := range methods {
		normalizedMethod := strings.ToUpper(strings.TrimSpace(method))
		switch normalizedMethod {
		case "GET", "POST":
		default:
			return HTTPException{}, fmt.Errorf("method %q is not allowed", method)
		}
		if _, duplicate := seenMethods[normalizedMethod]; duplicate {
			return HTTPException{}, fmt.Errorf("methods contain duplicate %q", normalizedMethod)
		}
		seenMethods[normalizedMethod] = struct{}{}
		normalizedMethods[index] = normalizedMethod
	}
	sort.Strings(normalizedMethods)
	return HTTPException{
		Scope:   strings.TrimSpace(scope),
		Methods: normalizedMethods,
		Owner:   strings.TrimSpace(owner),
		Reason:  strings.TrimSpace(reason),
	}, nil
}

func normalizeAuthProviderDecision(value any) (AuthProviderDecision, error) {
	object, ok := objectValue(value)
	if !ok {
		return AuthProviderDecision{}, errors.New("must be an object")
	}
	if !hasExactFields(object, "kind", "routeException") {
		return AuthProviderDecision{}, errors.New("better-auth requires exactly kind and routeException")
	}
	kind, ok := object["kind"].(string)
	if !ok || strings.TrimSpace(kind) != "better-auth" {
		return AuthProviderDecision{}, fmt.Errorf("kind %q is not allowed", kind)
	}
	exception, err := normalizeHTTPException(object["routeException"])
	if err != nil {
		return AuthProviderDecision{}, fmt.Errorf("routeException %w", err)
	}
	if exception.Owner != "Better Auth" {
		return AuthProviderDecision{}, errors.New(`routeException owner must be "Better Auth"`)
	}
	return AuthProviderDecision{Kind: "better-auth", RouteException: exception}, nil
}

func deriveHTTPContract(
	contract HTTPContract,
	provider *AuthProviderDecision,
) (HTTPContract, error) {
	if provider == nil {
		return contract, nil
	}
	key := normalizedHTTPExceptionKey(provider.RouteException)
	for _, exception := range contract.Exceptions {
		if normalizedHTTPExceptionKey(exception) != key {
			continue
		}
		if equalHTTPExceptions(exception, provider.RouteException) {
			return contract, nil
		}
		return HTTPContract{}, fmt.Errorf(
			"%q conflicts with %q for owner %q and scope %q",
			authProviderDecisionID,
			httpContractDecisionID,
			provider.RouteException.Owner,
			provider.RouteException.Scope,
		)
	}
	contract.Exceptions = append(contract.Exceptions, provider.RouteException)
	sortHTTPExceptions(contract.Exceptions)
	return contract, nil
}

func normalizedHTTPExceptionKey(exception HTTPException) string {
	return strings.ToLower(strings.TrimSpace(exception.Owner)) +
		"\x00" +
		strings.TrimSpace(exception.Scope)
}

func equalHTTPExceptions(first, second HTTPException) bool {
	return first.Scope == second.Scope &&
		first.Owner == second.Owner &&
		first.Reason == second.Reason &&
		slices.Equal(first.Methods, second.Methods)
}

func sortHTTPExceptions(exceptions []HTTPException) {
	sort.Slice(exceptions, func(i, j int) bool {
		if exceptions[i].Scope != exceptions[j].Scope {
			return exceptions[i].Scope < exceptions[j].Scope
		}
		if exceptions[i].Owner != exceptions[j].Owner {
			return exceptions[i].Owner < exceptions[j].Owner
		}
		return strings.Join(exceptions[i].Methods, "\x00") <
			strings.Join(exceptions[j].Methods, "\x00")
	})
}

func httpContractValue(contract HTTPContract) map[string]any {
	exceptions := make([]any, len(contract.Exceptions))
	for index, exception := range contract.Exceptions {
		exceptions[index] = httpExceptionValue(exception)
	}
	value := map[string]any{
		"mode":       contract.Mode,
		"exceptions": exceptions,
	}
	if contract.Source != nil {
		value["source"] = cloneJSONValue(contract.Source)
	}
	return value
}

func authProviderValue(provider AuthProviderDecision) map[string]any {
	return map[string]any{
		"kind":           provider.Kind,
		"routeException": httpExceptionValue(provider.RouteException),
	}
}

func httpExceptionValue(exception HTTPException) map[string]any {
	methods := make([]any, len(exception.Methods))
	for index, method := range exception.Methods {
		methods[index] = method
	}
	return map[string]any{
		"scope":   exception.Scope,
		"methods": methods,
		"owner":   exception.Owner,
		"reason":  exception.Reason,
	}
}

func hasOnlyFields(object document, fields ...string) bool {
	allowed := stringSet(fields)
	for field := range object {
		if _, ok := allowed[field]; !ok {
			return false
		}
	}
	return true
}
