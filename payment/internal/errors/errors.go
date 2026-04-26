package errs

import "errors"

var (

	// Ошибки валидации.
	ErrInvalidUUID          = errors.New("неверный формат UUID")
	ErrInvalidPaymentMethod = errors.New("неверный метод оплаты")
)
