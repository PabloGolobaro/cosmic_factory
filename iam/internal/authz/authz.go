// Package authz — встроенный PDP на базе OPA: оценивает rego-политики из policy/.
package authz

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/open-policy-agent/opa/v1/rego"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
)

const (
	decisionQuery = "data.authz.allow"
	policyPath    = "policy/authz.rego"
)

//go:embed policy/authz.rego
var policy string

// Engine оценивает политику доступа. Запрос компилируется один раз при создании.
type Engine struct {
	query rego.PreparedEvalQuery
}

// New компилирует встроенную политику.
func New(ctx context.Context) (*Engine, error) {
	query, err := rego.New(
		rego.Query(decisionQuery),
		rego.Module(policyPath, policy),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("компиляция политик: %w", err)
	}

	return &Engine{query: query}, nil
}

// Allow возвращает решение политики для переданного input.
func (e *Engine) Allow(ctx context.Context, in model.AuthzInput) (bool, error) {
	rs, err := e.query.Eval(ctx, rego.EvalInput(in))
	if err != nil {
		return false, fmt.Errorf("вычисление политики: %w", err)
	}

	return rs.Allowed(), nil
}
