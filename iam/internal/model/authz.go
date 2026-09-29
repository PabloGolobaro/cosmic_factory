package model

// AuthzInput — входные данные для политики OPA (input в rego).
type AuthzInput struct {
	Subject  AuthzSubject  `json:"subject"`
	Action   string        `json:"action"`
	Resource AuthzResource `json:"resource"`
}

// AuthzSubject — кто выполняет действие.
type AuthzSubject struct {
	UUID string `json:"uuid"`
	Role Role   `json:"role"`
}

// AuthzResource — над чем выполняется действие.
// OwnerUUID пустой, если у действия нет конкретного ресурса (например, создание).
type AuthzResource struct {
	OwnerUUID string `json:"owner_uuid"`
}

// AuthzDecision — результат авторизации.
type AuthzDecision struct {
	Allowed bool
	User    User
}
