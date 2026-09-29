package input

// AuthorizeInput — запрос на проверку прав владельца сессии.
// OwnerUUID пустой, если у действия нет конкретного ресурса.
type AuthorizeInput struct {
	SessionUUID string
	Action      string
	OwnerUUID   string
}
