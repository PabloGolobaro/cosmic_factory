// Package authz содержит общие для сервисов имена действий, проверяемых политикой OPA в IAM.
// Значения должны совпадать со строками в iam/internal/authz/policy/authz.rego.
package authz

const (
	ActionOrderCreate = "order:create"
	ActionOrderRead   = "order:read"
	ActionOrderPay    = "order:pay"
	ActionOrderCancel = "order:cancel"
)
