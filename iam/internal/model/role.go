package model

// Role — роль пользователя, определяет набор разрешённых действий в политике OPA.
type Role string

const (
	RoleClient  Role = "client"
	RoleManager Role = "manager"
)
